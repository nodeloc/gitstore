package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc/git-store/internal/config"
	"github.com/nodeloc/git-store/internal/models"
	"github.com/nodeloc/git-store/internal/utils"
	"gorm.io/gorm"
)

type DevAuthHandler struct {
	db     *gorm.DB
	config *config.Config
}

func NewDevAuthHandler(db *gorm.DB, cfg *config.Config) *DevAuthHandler {
	return &DevAuthHandler{
		db:     db,
		config: cfg,
	}
}

// DevLogin - 开发环境直接登录（仅用于测试）
func (h *DevAuthHandler) DevLogin(c *gin.Context) {
	// 只在开发环境启用
	if h.config.AppEnv == "production" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Dev login disabled in production"})
		return
	}

	email := c.DefaultQuery("email", "admin@test.com")

	var user models.User
	if err := h.db.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 生成 JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, h.config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"avatar_url": user.AvatarURL,
			"role":       user.Role,
		},
	})
}
