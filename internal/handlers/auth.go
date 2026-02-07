package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nodeloc/git-store/internal/config"
	"github.com/nodeloc/git-store/internal/models"
	"github.com/nodeloc/git-store/internal/services"
	"github.com/nodeloc/git-store/internal/utils"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db             *gorm.DB
	config         *config.Config
	githubOAuthSvc *services.GitHubOAuthService
	stateStore     map[string]bool // In production, use Redis
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:             db,
		config:         cfg,
		githubOAuthSvc: services.NewGitHubOAuthService(cfg),
		stateStore:     make(map[string]bool),
	}
}

func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	// Generate state token
	state := generateStateToken()
	h.stateStore[state] = true

	// Get GitHub OAuth URL
	authURL := h.githubOAuthSvc.GetAuthURL(state)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	// Verify state
	if !h.stateStore[state] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state token"})
		return
	}
	delete(h.stateStore, state)

	ctx := context.Background()

	// Exchange code for token
	token, err := h.githubOAuthSvc.ExchangeCode(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code"})
		return
	}

	// Get GitHub user info
	githubUser, err := h.githubOAuthSvc.GetUserInfo(ctx, token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Get user email
	emails, err := h.githubOAuthSvc.GetUserEmails(ctx, token.AccessToken)
	if err != nil || len(emails) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user email"})
		return
	}

	var primaryEmail string
	for _, email := range emails {
		if email.GetPrimary() {
			primaryEmail = email.GetEmail()
			break
		}
	}
	if primaryEmail == "" {
		primaryEmail = emails[0].GetEmail()
	}

	// Find or create user
	var user models.User
	err = h.db.Where("email = ?", primaryEmail).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		// Create new user
		user = models.User{
			Email:     primaryEmail,
			Name:      githubUser.GetName(),
			AvatarURL: githubUser.GetAvatarURL(),
			Role:      "user",
		}

		// Check if this is the admin user
		if primaryEmail == h.config.AdminEmail ||
			(h.config.AdminGitHubID != "" && h.config.AdminGitHubID == fmt.Sprintf("%d", githubUser.GetID())) {
			user.Role = "admin"
		}

		if err := h.db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	} else {
		// User exists, check and update admin role if needed
		shouldBeAdmin := primaryEmail == h.config.AdminEmail ||
			(h.config.AdminGitHubID != "" && h.config.AdminGitHubID == fmt.Sprintf("%d", githubUser.GetID()))
		
		if shouldBeAdmin && user.Role != "admin" {
			user.Role = "admin"
			h.db.Save(&user)
		} else if !shouldBeAdmin && user.Role == "admin" {
			// Optional: remove admin role if no longer in config
			// user.Role = "user"
			// h.db.Save(&user)
		}
	}

	// Create or update GitHub account
	var githubAccount models.GitHubAccount
	err = h.db.Where("git_hub_user_id = ?", githubUser.GetID()).First(&githubAccount).Error
	if err == gorm.ErrRecordNotFound {
		githubID := githubUser.GetID()
		githubAccount = models.GitHubAccount{
			UserID:       user.ID,
			GitHubUserID: &githubID,
			AccountType:  "user",
			Login:        githubUser.GetLogin(),
			AccessToken:  token.AccessToken,
		}
		if err := h.db.Create(&githubAccount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create GitHub account"})
			return
		}
	} else {
		// Update access token
		githubAccount.AccessToken = token.AccessToken
		h.db.Save(&githubAccount)
	}

	// Generate JWT token
	jwtToken, err := utils.GenerateJWT(user.ID, user.Email, user.Role, h.config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Redirect to frontend with token
	frontendURL := h.config.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3001"
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/callback?token=%s", frontendURL, jwtToken))
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	// Get user_id from context (set by AuthMiddleware)
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// userIDValue is uuid.UUID type from middleware
	uid, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var user models.User
	if err := h.db.Preload("GitHubAccounts").First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// In a stateless JWT system, logout is handled client-side
	// The client should delete the JWT token
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func (h *AuthHandler) GetGitHubAccounts(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var accounts []models.GitHubAccount
	if err := h.db.Where("user_id = ?", uid).Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get GitHub accounts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}

func generateStateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
