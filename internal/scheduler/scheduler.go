package scheduler

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/nodeloc/git-store/internal/config"
	"github.com/nodeloc/git-store/internal/models"
	"github.com/nodeloc/git-store/internal/services"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Scheduler struct {
	db              *gorm.DB
	config          *config.Config
	githubSvc       *services.GitHubService
	emailSvc        *services.EmailService
	exchangeRateSvc *services.ExchangeRateService
}

// Helper function to split "owner/repo" format
func splitRepoName(repoName string) []string {
	return strings.Split(repoName, "/")
}

func SetupScheduler(c *cron.Cron, db *gorm.DB, cfg *config.Config) {
	scheduler := &Scheduler{
		db:              db,
		config:          cfg,
		emailSvc:        services.NewEmailService(cfg, db),
		exchangeRateSvc: services.NewExchangeRateService(db, cfg),
	}

	// Initialize GitHub Service
	if cfg.GitHubAdminToken != "" {
		scheduler.githubSvc = services.NewGitHubService(cfg)
		log.Println("Scheduler: GitHub Service initialized")
	} else {
		log.Println("Scheduler: GitHub Admin Token not configured, some features may not work")
	}

	// Schedule maintenance expiry check (default: daily at 2 AM)
	c.AddFunc(cfg.CronMaintenanceCheck, func() {
		log.Println("Running maintenance expiry check...")
		scheduler.CheckMaintenanceExpiry()
	})

	// Schedule daily statistics aggregation (daily at 1 AM)
	c.AddFunc("0 1 * * *", func() {
		log.Println("Running daily statistics aggregation...")
		scheduler.AggregateStatistics()
	})

	// Schedule daily exchange rate update (daily at 3 AM)
	c.AddFunc("0 3 * * *", func() {
		log.Println("Running daily exchange rate update...")
		scheduler.UpdateExchangeRates()
	})

	log.Println("Scheduler initialized")
}

// CheckMaintenanceExpiry checks for expired and expiring licenses
func (s *Scheduler) CheckMaintenanceExpiry() {
	ctx := context.Background()
	today := time.Now()

	// 1. Find expired licenses (status=active, maintenance_until < today)
	var expiredLicenses []models.License
	err := s.db.Preload("User").Preload("Plugin").Preload("GitHubAccount").
		Where("status = ? AND maintenance_until < ?", "active", today).
		Find(&expiredLicenses).Error

	if err != nil {
		log.Printf("Error finding expired licenses: %v", err)
		return
	}

	log.Printf("Found %d expired licenses", len(expiredLicenses))

	// Process each expired license
	for _, license := range expiredLicenses {
		if err := s.processExpiredLicense(ctx, &license); err != nil {
			log.Printf("Error processing expired license %s: %v", license.ID, err)
		}
	}

	// 2. Find licenses expiring in 30, 7, 1 days
	expiryWarningDays := []int{30, 7, 1}

	for _, days := range expiryWarningDays {
		targetDate := today.AddDate(0, 0, days)
		startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)

		var expiringLicenses []models.License
		err := s.db.Preload("User").Preload("Plugin").
			Where("status = ? AND maintenance_until >= ? AND maintenance_until < ?",
				"active", startOfDay, endOfDay).
			Find(&expiringLicenses).Error

		if err != nil {
			log.Printf("Error finding licenses expiring in %d days: %v", days, err)
			continue
		}

		log.Printf("Found %d licenses expiring in %d days", len(expiringLicenses), days)

		// Send warning emails
		for _, license := range expiringLicenses {
			if err := s.sendExpiryWarning(&license, days); err != nil {
				log.Printf("Error sending expiry warning for license %s: %v", license.ID, err)
			}
		}
	}
}

func (s *Scheduler) processExpiredLicense(ctx context.Context, license *models.License) error {
	log.Printf("Processing expired license: %s (Plugin: %s, User: %s)",
		license.ID, license.Plugin.Name, license.User.Email)

	// 1. Remove GitHub repository collaborator access
	if s.githubSvc != nil && license.Plugin.GitHubRepoName != "" && license.GitHubAccount.Login != "" {
		// Parse owner/repo
		repoParts := splitRepoName(license.Plugin.GitHubRepoName)
		if len(repoParts) == 2 {
			owner := repoParts[0]
			repo := repoParts[1]
			username := license.GitHubAccount.Login

			log.Printf("[License Expiry] Removing %s from %s/%s", username, owner, repo)

			err := s.githubSvc.RemoveRepositoryCollaborator(ctx, owner, repo, username)
			if err != nil {
				log.Printf("Failed to remove collaborator: %v", err)
				// Continue even if removal fails
			} else {
				log.Printf("Successfully removed %s from %s/%s", username, owner, repo)

				// Log the action
				history := models.LicenseHistory{
					LicenseID:  license.ID,
					Action:     "github_access_revoked",
					OccurredAt: time.Now(),
				}
				s.db.Create(&history)
			}
		}
	}

	// 2. Update license status to expired
	license.Status = "expired"
	if err := s.db.Save(license).Error; err != nil {
		return err
	}

	// 3. Log expiry action
	history := models.LicenseHistory{
		LicenseID:  license.ID,
		Action:     "expired",
		OccurredAt: time.Now(),
	}
	if err := s.db.Create(&history).Error; err != nil {
		log.Printf("Failed to create history record: %v", err)
	}

	// 4. Send expiry email notification
	if err := s.emailSvc.SendMaintenanceExpiredEmail(&license.User, &license.Plugin, license); err != nil {
		log.Printf("Failed to send expiry email: %v", err)
	}

	log.Printf("Successfully processed expired license %s", license.ID)
	return nil
}

func (s *Scheduler) sendExpiryWarning(license *models.License, daysRemaining int) error {
	// Check if we've already sent this warning
	notificationType := ""
	switch daysRemaining {
	case 30:
		notificationType = "maintenance_expiring_30"
	case 7:
		notificationType = "maintenance_expiring_7"
	case 1:
		notificationType = "maintenance_expiring_1"
	default:
		return nil
	}

	// Check if notification already sent
	var count int64
	s.db.Model(&models.EmailNotification{}).
		Where("user_id = ? AND notification_type = ? AND created_at > ?",
			license.UserID, notificationType, time.Now().AddDate(0, 0, -1)).
		Count(&count)

	if count > 0 {
		log.Printf("Warning email already sent for license %s", license.ID)
		return nil
	}

	// Send warning email
	return s.emailSvc.SendMaintenanceExpiringEmail(&license.User, &license.Plugin, license, daysRemaining)
}

// AggregateStatistics aggregates daily statistics
func (s *Scheduler) AggregateStatistics() {
	yesterday := time.Now().AddDate(0, 0, -1)
	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var stat models.Statistic
	stat.StatDate = startOfDay

	// Count total users
	var totalUsers int64
	s.db.Model(&models.User{}).Count(&totalUsers)
	stat.TotalUsers = int(totalUsers)

	// Count new users
	var newUsers int64
	s.db.Model(&models.User{}).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Count(&newUsers)
	stat.NewUsers = int(newUsers)

	// Count total orders
	var totalOrders int64
	s.db.Model(&models.Order{}).
		Where("payment_status = ?", "paid").
		Count(&totalOrders)
	stat.TotalOrders = int(totalOrders)

	// Count new orders
	var newOrders int64
	s.db.Model(&models.Order{}).
		Where("payment_status = ? AND paid_at >= ? AND paid_at < ?", "paid", startOfDay, endOfDay).
		Count(&newOrders)
	stat.NewOrders = int(newOrders)

	// Calculate total revenue
	var totalRevenue float64
	s.db.Model(&models.Order{}).
		Where("payment_status = ?", "paid").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalRevenue)
	stat.TotalRevenue = totalRevenue

	// Calculate daily revenue
	var dailyRevenue float64
	s.db.Model(&models.Order{}).
		Where("payment_status = ? AND paid_at >= ? AND paid_at < ?", "paid", startOfDay, endOfDay).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&dailyRevenue)
	stat.DailyRevenue = dailyRevenue

	// Count active licenses
	var activeLicenses int64
	s.db.Model(&models.License{}).
		Where("status = ?", "active").
		Count(&activeLicenses)
	stat.ActiveLicenses = int(activeLicenses)

	// Count expired licenses
	var expiredLicenses int64
	s.db.Model(&models.License{}).
		Where("status = ?", "expired").
		Count(&expiredLicenses)
	stat.ExpiredLicenses = int(expiredLicenses)

	// Save or update statistics
	var existingStat models.Statistic
	err := s.db.Where("stat_date = ?", startOfDay).First(&existingStat).Error
	if err == gorm.ErrRecordNotFound {
		// Create new statistic
		if err := s.db.Create(&stat).Error; err != nil {
			log.Printf("Error creating statistics: %v", err)
			return
		}
	} else {
		// Update existing statistic
		stat.ID = existingStat.ID
		if err := s.db.Save(&stat).Error; err != nil {
			log.Printf("Error updating statistics: %v", err)
			return
		}
	}

	log.Printf("Successfully aggregated statistics for %s", startOfDay.Format("2006-01-02"))
}

// UpdateExchangeRates 更新汇率
func (s *Scheduler) UpdateExchangeRates() {
	if err := s.exchangeRateSvc.UpdateExchangeRates(); err != nil {
		log.Printf("❌ Failed to update exchange rates: %v", err)
	}
}
