package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc/git-store/internal/config"
	"github.com/nodeloc/git-store/internal/handlers"
	"github.com/nodeloc/git-store/internal/middleware"
	"github.com/nodeloc/git-store/internal/services"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Middleware
	r.Use(middleware.CORSMiddleware())

	// Serve static files (uploaded images)
	r.Static("/uploads", "./uploads")

	// Initialize GitHub service (using Personal Access Token)
	var githubSvc *services.GitHubService
	if cfg.GitHubAdminToken != "" {
		githubSvc = services.NewGitHubService(cfg)
		log.Println("GitHub Service initialized successfully")
	} else {
		log.Println("Warning: GitHub Admin Token not configured, repository access features will be disabled")
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	pluginHandler := handlers.NewPluginHandler(db, cfg)
	orderHandler := handlers.NewOrderHandler(db, cfg)
	paymentHandler := handlers.NewPaymentHandler(db, cfg)
	licenseHandler := handlers.NewLicenseHandler(db, cfg)
	tutorialHandler := handlers.NewTutorialHandler(db, cfg)
	categoryHandler := handlers.NewCategoryHandler(db)
	pageHandler := handlers.NewPageHandler(db)
	adminHandler := handlers.NewAdminHandler(db, cfg, githubSvc)
	dashboardHandler := handlers.NewDashboardHandler(db, cfg)
	githubWebhookHandler := handlers.NewGitHubWebhookHandler(db, cfg)
	uploadHandler := handlers.NewUploadHandler("./uploads")
	configHandler := handlers.NewConfigHandler(cfg)

	// Dev auth handler (only in development)
	var devAuthHandler *handlers.DevAuthHandler
	if cfg.AppEnv == "development" {
		devAuthHandler = handlers.NewDevAuthHandler(db, cfg)
	}

	// Public routes
	api := r.Group("/api")
	{
		// Health check (support both GET and HEAD)
		healthCheck := func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		}
		api.GET("/health", healthCheck)
		api.HEAD("/health", healthCheck)

		// Public config
		api.GET("/config", configHandler.GetPublicConfig)

		// Dev login (only in development)
		if cfg.AppEnv == "development" && devAuthHandler != nil {
			api.GET("/dev/login", devAuthHandler.DevLogin)
		}

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.GET("/github", authHandler.GitHubLogin)
			auth.GET("/github/callback", authHandler.GitHubCallback)
			auth.GET("/me", middleware.AuthMiddleware(cfg), authHandler.GetMe)
			auth.POST("/logout", middleware.AuthMiddleware(cfg), authHandler.Logout)
		}

		// Public plugin routes
		plugins := api.Group("/plugins")
		{
			plugins.GET("", pluginHandler.ListPlugins)
			plugins.GET("/id/:id", pluginHandler.GetPluginByID)
			plugins.GET("/:slug", pluginHandler.GetPlugin)
		}

		// Public tutorial routes
		tutorials := api.Group("/tutorials")
		{
			tutorials.GET("/public", tutorialHandler.ListPublicTutorials)
			tutorials.GET("/:slug", tutorialHandler.GetTutorial)
		}

		// Public categories routes
		api.GET("/categories", categoryHandler.GetCategories)

		// Public page routes
		pages := api.Group("/pages")
		{
			pages.GET("", pageHandler.GetPublicPages)
			pages.GET("/:slug", pageHandler.GetPublicPageBySlug)
		}

		// Public settings route (for site name, etc.)
		api.GET("/settings/public", adminHandler.GetPublicSettings)

		// Payment webhook routes (no auth)
		webhooks := api.Group("/webhooks")
		{
			webhooks.POST("/stripe", paymentHandler.StripeWebhook)
			webhooks.POST("/paypal", paymentHandler.PayPalWebhook)
			webhooks.POST("/alipay", paymentHandler.AlipayNotify)
			webhooks.GET("/alipay", paymentHandler.AlipayNotify) // 易支付使用 GET 请求
			webhooks.POST("/github", githubWebhookHandler.HandleGitHubAppWebhook)
		}

		// Public license verification API (no auth required)
		api.GET("/licenses/:id/verify", licenseHandler.VerifyLicense)
	}

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// User routes
		user := protected.Group("/user")
		{
			user.GET("/licenses", licenseHandler.GetUserLicenses)
			user.GET("/orders", orderHandler.GetUserOrders)
			user.GET("/github-accounts", authHandler.GetGitHubAccounts)
			user.GET("/github-app/status", githubWebhookHandler.GetInstallationStatus)
		}

		// Order routes
		orders := protected.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrder)
		}

		// Payment routes
		payments := protected.Group("/payments")
		{
			payments.POST("/stripe/create-intent", paymentHandler.CreateStripePaymentIntent)
			payments.POST("/paypal/create-order", paymentHandler.CreatePayPalOrder)
			payments.POST("/paypal/capture-order", paymentHandler.CapturePayPalOrder)
			payments.POST("/alipay/create", paymentHandler.CreateAlipayPayment)
		}

		// License routes
		licenses := protected.Group("/licenses")
		{
			licenses.GET("/:id", licenseHandler.GetLicense)
			licenses.POST("/:id/renew", licenseHandler.RenewLicense)
			licenses.GET("/:id/history", licenseHandler.GetLicenseHistory)
		}

		// Tutorial routes (protected)
		tutorialsProtected := protected.Group("/tutorials")
		{
			tutorialsProtected.GET("", tutorialHandler.ListTutorials)
		}
	}

	// Admin routes (require admin role)
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg))
	admin.Use(middleware.AdminMiddleware())
	{
		// Plugin management
		adminPlugins := admin.Group("/plugins")
		{
			adminPlugins.GET("", adminHandler.ListAllPlugins)
			adminPlugins.POST("", adminHandler.CreatePlugin)
			adminPlugins.GET("/:id", adminHandler.GetPluginByID)
			adminPlugins.PUT("/:id", adminHandler.UpdatePlugin)
			adminPlugins.DELETE("/:id", adminHandler.DeletePlugin)
			adminPlugins.POST("/sync-repos", adminHandler.SyncGitHubRepos)
		}

		// GitHub integration
		adminGitHub := admin.Group("/github")
		{
			adminGitHub.GET("/repositories", adminHandler.ListGitHubRepos)
		}

		// Order management
		adminOrders := admin.Group("/orders")
		{
			adminOrders.GET("", adminHandler.ListAllOrders)
			adminOrders.GET("/:id", adminHandler.GetOrderByID)
			adminOrders.PUT("/:id/status", adminHandler.UpdateOrderPaymentStatus)
			adminOrders.POST("/:id/refund", adminHandler.RefundOrder)
		}

		// License management
		adminLicenses := admin.Group("/licenses")
		{
			adminLicenses.GET("", adminHandler.ListAllLicenses)
			adminLicenses.GET("/:id", adminHandler.GetLicenseByID)
			adminLicenses.POST("/:id/revoke", adminHandler.RevokeLicense)
			adminLicenses.POST("/:id/extend", adminHandler.ExtendLicense)
		}

		// Tutorial management
		adminTutorials := admin.Group("/tutorials")
		{
			adminTutorials.GET("", adminHandler.ListAllTutorials)
			adminTutorials.POST("", adminHandler.CreateTutorial)
			adminTutorials.GET("/:id", adminHandler.GetTutorialByID)
			adminTutorials.PUT("/:id", adminHandler.UpdateTutorial)
			adminTutorials.DELETE("/:id", adminHandler.DeleteTutorial)
		}

		// Categories management
		adminCategories := admin.Group("/categories")
		{
			adminCategories.GET("", categoryHandler.GetAllCategories)
			adminCategories.POST("", categoryHandler.CreateCategory)
			adminCategories.PUT("/:id", categoryHandler.UpdateCategory)
			adminCategories.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		// Pages management
		adminPages := admin.Group("/pages")
		{
			adminPages.GET("", pageHandler.GetAdminPages)
			adminPages.GET("/:id", pageHandler.GetAdminPageByID)
			adminPages.POST("", pageHandler.CreatePage)
			adminPages.PUT("/:id", pageHandler.UpdatePage)
			adminPages.DELETE("/:id", pageHandler.DeletePage)
		}

		// Image upload
		admin.POST("/upload/image", uploadHandler.UploadImage)

		// Statistics & Dashboard
		adminStats := admin.Group("/statistics")
		{
			adminStats.GET("/dashboard", dashboardHandler.GetDashboardStats)
			adminStats.GET("/revenue", dashboardHandler.GetRevenueStats)
			adminStats.GET("/users", dashboardHandler.GetUserStats)
			adminStats.GET("/plugins", dashboardHandler.GetPluginStats)
		}

		// System settings
		adminSettings := admin.Group("/settings")
		{
			adminSettings.GET("", adminHandler.GetSettings)
			adminSettings.PUT("", adminHandler.UpdateSettings)
		}

		// Exchange rates management
		adminExchangeRates := admin.Group("/exchange-rates")
		{
			adminExchangeRates.GET("", adminHandler.GetExchangeRates)
			adminExchangeRates.POST("/update", adminHandler.UpdateExchangeRates)
		}

		// User management
		adminUsers := admin.Group("/users")
		{
			adminUsers.GET("", adminHandler.ListAllUsers)
			adminUsers.GET("/:id", adminHandler.GetUserByID)
			adminUsers.PUT("/:id", adminHandler.UpdateUser)
			adminUsers.DELETE("/:id", adminHandler.DeleteUser)
		}
	}
}
