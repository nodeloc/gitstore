package services

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/nodeloc/git-store/internal/config"
	"github.com/nodeloc/git-store/internal/models"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type EmailService struct {
	config *config.Config
	db     *gorm.DB
	dialer *gomail.Dialer
}

type EmailData struct {
	UserName         string
	PluginName       string
	OrderNumber      string
	Amount           string
	MaintenanceUntil string
	DaysRemaining    int
	RepoURL          string
	TutorialURL      string
	RenewalURL       string
	SupportEmail     string
	SiteName         string
}

func NewEmailService(cfg *config.Config, db *gorm.DB) *EmailService {
	dialer := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)

	return &EmailService{
		config: cfg,
		db:     db,
		dialer: dialer,
	}
}

func (s *EmailService) SendEmail(to, subject, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.SMTPFromName, s.config.SMTPFrom))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *EmailService) SendPurchaseSuccessEmail(user *models.User, plugin *models.Plugin, order *models.Order, license *models.License) error {
	data := EmailData{
		UserName:         user.Name,
		PluginName:       plugin.Name,
		OrderNumber:      order.OrderNumber,
		Amount:           fmt.Sprintf("%.2f %s", order.Amount, order.Currency),
		MaintenanceUntil: license.MaintenanceUntil.Format("2006-01-02"),
		RepoURL:          plugin.GitHubRepoURL,
		TutorialURL:      fmt.Sprintf("%s/tutorials/%s", s.config.FrontendURL, plugin.Slug),
		SupportEmail:     s.config.AdminEmail,
		SiteName:         "Plugin Store",
	}

	subject := fmt.Sprintf("Purchase Successful - %s", plugin.Name)
	htmlBody, err := s.renderTemplate("purchase_success", data)
	if err != nil {
		return err
	}

	// Send email
	if err := s.SendEmail(user.Email, subject, htmlBody); err != nil {
		return err
	}

	// Log email notification
	notification := &models.EmailNotification{
		UserID:           user.ID,
		NotificationType: "purchase_success",
		Subject:          subject,
		Body:             htmlBody,
		Status:           "sent",
		SentAt:           &time.Time{},
	}
	*notification.SentAt = time.Now()

	return s.db.Create(notification).Error
}

func (s *EmailService) SendMaintenanceExpiringEmail(user *models.User, plugin *models.Plugin, license *models.License, daysRemaining int) error {
	data := EmailData{
		UserName:         user.Name,
		PluginName:       plugin.Name,
		MaintenanceUntil: license.MaintenanceUntil.Format("2006-01-02"),
		DaysRemaining:    daysRemaining,
		RepoURL:          plugin.GitHubRepoURL,
		RenewalURL:       fmt.Sprintf("%s/renew/%s", s.config.FrontendURL, license.ID),
		SupportEmail:     s.config.AdminEmail,
		SiteName:         "Plugin Store",
	}

	subject := fmt.Sprintf("Maintenance Expiring Soon - %s (%d days remaining)", plugin.Name, daysRemaining)
	htmlBody, err := s.renderTemplate("maintenance_expiring", data)
	if err != nil {
		return err
	}

	if err := s.SendEmail(user.Email, subject, htmlBody); err != nil {
		return err
	}

	notificationType := fmt.Sprintf("maintenance_expiring_%d", daysRemaining)
	notification := &models.EmailNotification{
		UserID:           user.ID,
		NotificationType: notificationType,
		Subject:          subject,
		Body:             htmlBody,
		Status:           "sent",
		SentAt:           &time.Time{},
	}
	*notification.SentAt = time.Now()

	return s.db.Create(notification).Error
}

func (s *EmailService) SendMaintenanceExpiredEmail(user *models.User, plugin *models.Plugin, license *models.License) error {
	data := EmailData{
		UserName:         user.Name,
		PluginName:       plugin.Name,
		MaintenanceUntil: license.MaintenanceUntil.Format("2006-01-02"),
		RenewalURL:       fmt.Sprintf("%s/renew/%s", s.config.FrontendURL, license.ID),
		SupportEmail:     s.config.AdminEmail,
		SiteName:         "Plugin Store",
	}

	subject := fmt.Sprintf("Maintenance Expired - %s", plugin.Name)
	htmlBody, err := s.renderTemplate("maintenance_expired", data)
	if err != nil {
		return err
	}

	if err := s.SendEmail(user.Email, subject, htmlBody); err != nil {
		return err
	}

	notification := &models.EmailNotification{
		UserID:           user.ID,
		NotificationType: "maintenance_expired",
		Subject:          subject,
		Body:             htmlBody,
		Status:           "sent",
		SentAt:           &time.Time{},
	}
	*notification.SentAt = time.Now()

	return s.db.Create(notification).Error
}

func (s *EmailService) SendRenewalSuccessEmail(user *models.User, plugin *models.Plugin, license *models.License) error {
	data := EmailData{
		UserName:         user.Name,
		PluginName:       plugin.Name,
		MaintenanceUntil: license.MaintenanceUntil.Format("2006-01-02"),
		RepoURL:          plugin.GitHubRepoURL,
		SupportEmail:     s.config.AdminEmail,
		SiteName:         "Plugin Store",
	}

	subject := fmt.Sprintf("Renewal Successful - %s", plugin.Name)
	htmlBody, err := s.renderTemplate("renewal_success", data)
	if err != nil {
		return err
	}

	if err := s.SendEmail(user.Email, subject, htmlBody); err != nil {
		return err
	}

	notification := &models.EmailNotification{
		UserID:           user.ID,
		NotificationType: "renewal_success",
		Subject:          subject,
		Body:             htmlBody,
		Status:           "sent",
		SentAt:           &time.Time{},
	}
	*notification.SentAt = time.Now()

	return s.db.Create(notification).Error
}

func (s *EmailService) renderTemplate(templateName string, data EmailData) (string, error) {
	templates := map[string]string{
		"purchase_success": `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
        .button { display: inline-block; padding: 10px 20px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Purchase Successful!</h1>
        </div>
        <div class="content">
            <p>Hi {{.UserName}},</p>
            <p>Thank you for purchasing <strong>{{.PluginName}}</strong>!</p>
            <p><strong>Order Details:</strong></p>
            <ul>
                <li>Order Number: {{.OrderNumber}}</li>
                <li>Amount: {{.Amount}}</li>
                <li>Maintenance Until: {{.MaintenanceUntil}}</li>
            </ul>
            <p><strong>Repository URL:</strong> <a href="{{.RepoURL}}">{{.RepoURL}}</a></p>
            <p><a href="{{.TutorialURL}}" class="button">View Installation Tutorial</a></p>
            <p>Your plugin is now accessible via GitHub. You'll receive update access until {{.MaintenanceUntil}}.</p>
        </div>
        <div class="footer">
            <p>Need help? Contact us at {{.SupportEmail}}</p>
            <p>&copy; 2024 {{.SiteName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,
		"maintenance_expiring": `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #FF9800; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
        .button { display: inline-block; padding: 10px 20px; background-color: #FF9800; color: white; text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Maintenance Expiring Soon</h1>
        </div>
        <div class="content">
            <p>Hi {{.UserName}},</p>
            <p>Your maintenance period for <strong>{{.PluginName}}</strong> will expire in <strong>{{.DaysRemaining}} days</strong>.</p>
            <p><strong>Expiry Date:</strong> {{.MaintenanceUntil}}</p>
            <p>After expiry, you'll no longer be able to pull updates from the repository, but the plugin will continue to work.</p>
            <p><a href="{{.RenewalURL}}" class="button">Renew Maintenance</a></p>
        </div>
        <div class="footer">
            <p>Need help? Contact us at {{.SupportEmail}}</p>
            <p>&copy; 2024 {{.SiteName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,
		"maintenance_expired": `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #F44336; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
        .button { display: inline-block; padding: 10px 20px; background-color: #F44336; color: white; text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Maintenance Expired</h1>
        </div>
        <div class="content">
            <p>Hi {{.UserName}},</p>
            <p>Your maintenance period for <strong>{{.PluginName}}</strong> has expired on {{.MaintenanceUntil}}.</p>
            <p>Your plugin will continue to work, but you can no longer pull updates from the repository.</p>
            <p>To restore update access, please renew your maintenance.</p>
            <p><a href="{{.RenewalURL}}" class="button">Renew Now</a></p>
        </div>
        <div class="footer">
            <p>Need help? Contact us at {{.SupportEmail}}</p>
            <p>&copy; 2024 {{.SiteName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,
		"renewal_success": `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Renewal Successful!</h1>
        </div>
        <div class="content">
            <p>Hi {{.UserName}},</p>
            <p>Your maintenance for <strong>{{.PluginName}}</strong> has been successfully renewed!</p>
            <p><strong>New Expiry Date:</strong> {{.MaintenanceUntil}}</p>
            <p>You can now continue receiving updates from the repository.</p>
            <p><strong>Repository URL:</strong> <a href="{{.RepoURL}}">{{.RepoURL}}</a></p>
        </div>
        <div class="footer">
            <p>Need help? Contact us at {{.SupportEmail}}</p>
            <p>&copy; 2024 {{.SiteName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,
	}

	tmplStr, ok := templates[templateName]
	if !ok {
		return "", fmt.Errorf("template not found: %s", templateName)
	}

	tmpl, err := template.New(templateName).Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
