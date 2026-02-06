package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/nodeloc/git-store/internal/models"
)

type PageHandler struct {
	db *gorm.DB
}

func NewPageHandler(db *gorm.DB) *PageHandler {
	return &PageHandler{db: db}
}

// GetPublicPages returns all published pages
func (h *PageHandler) GetPublicPages(c *gin.Context) {
	var pages []models.Page
	if err := h.db.Where("status = ?", "published").Order("sort_order ASC, title ASC").Find(&pages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages"})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// GetPublicPageBySlug returns a single published page by slug
func (h *PageHandler) GetPublicPageBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var page models.Page
	if err := h.db.Where("slug = ? AND status = ?", slug, "published").First(&page).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch page"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// GetAdminPages returns all pages for admin (including drafts)
func (h *PageHandler) GetAdminPages(c *gin.Context) {
	var pages []models.Page
	if err := h.db.Order("sort_order ASC, title ASC").Find(&pages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages"})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// GetAdminPageByID returns a single page by ID (for editing)
func (h *PageHandler) GetAdminPageByID(c *gin.Context) {
	pageID := c.Param("id")

	id, err := uuid.Parse(pageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	var page models.Page
	if err := h.db.First(&page, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch page"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// CreatePage creates a new page
func (h *PageHandler) CreatePage(c *gin.Context) {
	var input struct {
		Slug      string `json:"slug" binding:"required"`
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Status    string `json:"status"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if slug already exists
	var existingPage models.Page
	if err := h.db.Where("slug = ?", input.Slug).First(&existingPage).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page with this slug already exists"})
		return
	}

	status := input.Status
	if status == "" {
		status = "draft"
	}

	page := models.Page{
		Slug:      input.Slug,
		Title:     input.Title,
		Content:   input.Content,
		Status:    status,
		SortOrder: input.SortOrder,
	}

	if err := h.db.Create(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create page"})
		return
	}

	c.JSON(http.StatusCreated, page)
}

// UpdatePage updates an existing page
func (h *PageHandler) UpdatePage(c *gin.Context) {
	pageID := c.Param("id")

	id, err := uuid.Parse(pageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	var page models.Page
	if err := h.db.First(&page, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch page"})
		return
	}

	var input struct {
		Slug      string `json:"slug"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		Status    string `json:"status"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if slug is being changed and if it already exists
	if input.Slug != "" && input.Slug != page.Slug {
		var existingPage models.Page
		if err := h.db.Where("slug = ? AND id != ?", input.Slug, id).First(&existingPage).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Page with this slug already exists"})
			return
		}
		page.Slug = input.Slug
	}

	if input.Title != "" {
		page.Title = input.Title
	}
	if input.Content != "" {
		page.Content = input.Content
	}
	if input.Status != "" {
		page.Status = input.Status
	}
	page.SortOrder = input.SortOrder
	page.UpdatedAt = time.Now()

	if err := h.db.Save(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// DeletePage deletes a page
func (h *PageHandler) DeletePage(c *gin.Context) {
	pageID := c.Param("id")

	id, err := uuid.Parse(pageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	var page models.Page
	if err := h.db.First(&page, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch page"})
		return
	}

	if err := h.db.Delete(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}
