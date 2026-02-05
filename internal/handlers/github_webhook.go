package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc/git-store/internal/config"
	"github.com/nodeloc/git-store/internal/models"
	"gorm.io/gorm"
)

type GitHubWebhookHandler struct {
	db     *gorm.DB
	config *config.Config
}

func NewGitHubWebhookHandler(db *gorm.DB, cfg *config.Config) *GitHubWebhookHandler {
	return &GitHubWebhookHandler{
		db:     db,
		config: cfg,
	}
}

// GitHub App Installation Event
type GitHubInstallationEvent struct {
	Action       string `json:"action"` // created, deleted
	Installation struct {
		ID      int64 `json:"id"`
		Account struct {
			Login string `json:"login"`
			ID    int64  `json:"id"`
			Type  string `json:"type"` // User, Organization
		} `json:"account"`
	} `json:"installation"`
	Repositories []struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repositories"`
}

// HandleGitHubAppWebhook handles GitHub App installation and other events
func (h *GitHubWebhookHandler) HandleGitHubAppWebhook(c *gin.Context) {
	// Get event type from header
	eventType := c.GetHeader("X-GitHub-Event")

	log.Printf("[GitHub Webhook] Received event: %s", eventType)

	switch eventType {
	case "installation":
		h.handleInstallationEvent(c)
	case "installation_repositories":
		h.handleInstallationRepositoriesEvent(c)
	default:
		log.Printf("[GitHub Webhook] Unhandled event type: %s", eventType)
		c.JSON(http.StatusOK, gin.H{"message": "Event received"})
	}
}

func (h *GitHubWebhookHandler) handleInstallationEvent(c *gin.Context) {
	var event GitHubInstallationEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		log.Printf("[GitHub Webhook] Failed to parse installation event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	log.Printf("[GitHub Webhook] Installation event - Action: %s, ID: %d, Account: %s (Type: %s)",
		event.Action, event.Installation.ID, event.Installation.Account.Login, event.Installation.Account.Type)

	ctx := context.Background()

	switch event.Action {
	case "created":
		// User installed the app
		h.handleInstallationCreated(ctx, &event)
	case "deleted":
		// User uninstalled the app
		h.handleInstallationDeleted(ctx, &event)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Installation event processed"})
}

func (h *GitHubWebhookHandler) handleInstallationCreated(ctx context.Context, event *GitHubInstallationEvent) {
	installationID := event.Installation.ID
	githubLogin := event.Installation.Account.Login
	accountType := event.Installation.Account.Type

	log.Printf("[GitHub App] Installation created - ID: %d, Login: %s, Type: %s", installationID, githubLogin, accountType)

	// Find GitHubAccount by login
	var githubAccount models.GitHubAccount
	err := h.db.Where("login = ?", githubLogin).First(&githubAccount).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[GitHub App] GitHub account not found for login: %s - User needs to login first", githubLogin)
			return
		}
		log.Printf("[GitHub App] Error finding GitHub account: %v", err)
		return
	}

	// Update installation_id
	githubAccount.InstallationID = &installationID
	if err := h.db.Save(&githubAccount).Error; err != nil {
		log.Printf("[GitHub App] Failed to update installation_id: %v", err)
		return
	}

	log.Printf("[GitHub App] Successfully linked installation %d to account %s", installationID, githubLogin)

	// Check if user has active licenses and grant access to repositories
	h.grantAccessToActiveLicenses(ctx, &githubAccount)
}

func (h *GitHubWebhookHandler) handleInstallationDeleted(ctx context.Context, event *GitHubInstallationEvent) {
	installationID := event.Installation.ID
	githubLogin := event.Installation.Account.Login

	log.Printf("[GitHub App] Installation deleted - ID: %d, Login: %s", installationID, githubLogin)

	// Find and update GitHubAccount
	var githubAccount models.GitHubAccount
	err := h.db.Where("login = ? AND installation_id = ?", githubLogin, installationID).First(&githubAccount).Error

	if err == nil {
		githubAccount.InstallationID = nil
		h.db.Save(&githubAccount)
		log.Printf("[GitHub App] Removed installation_id from account %s", githubLogin)
	}
}

func (h *GitHubWebhookHandler) handleInstallationRepositoriesEvent(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	action := payload["action"]
	log.Printf("[GitHub Webhook] Installation repositories event - Action: %v", action)

	c.JSON(http.StatusOK, gin.H{"message": "Installation repositories event processed"})
}

// grantAccessToActiveLicenses grants repository access for all active licenses
func (h *GitHubWebhookHandler) grantAccessToActiveLicenses(ctx context.Context, githubAccount *models.GitHubAccount) {
	if githubAccount.InstallationID == nil {
		log.Printf("[GitHub App] No installation_id for account %s", githubAccount.Login)
		return
	}

	// Find active licenses for this GitHub account
	var licenses []models.License
	err := h.db.Preload("Plugin").
		Where("git_hub_account_id = ? AND status = ?", githubAccount.ID, "active").
		Find(&licenses).Error

	if err != nil {
		log.Printf("[GitHub] Error finding active licenses: %v", err)
		return
	}

	log.Printf("[GitHub] Found %d active licenses for %s (collaborator invitations already sent)", len(licenses), githubAccount.Login)
	// Note: Collaborator invitations are sent when license is created
	// No additional action needed here
}

// GetInstallationStatus checks if user has installed the GitHub App
func (h *GitHubWebhookHandler) GetInstallationStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var githubAccount models.GitHubAccount
	err := h.db.Where("user_id = ?", userID).First(&githubAccount).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"installed": false,
			"message":   "GitHub account not linked",
		})
		return
	}

	installed := githubAccount.InstallationID != nil

	response := gin.H{
		"installed":    installed,
		"github_login": githubAccount.Login,
	}

	if installed {
		response["installation_id"] = *githubAccount.InstallationID
	}

	c.JSON(http.StatusOK, response)
}
