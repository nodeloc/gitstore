package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nodeloc/git-store/internal/config"
	"github.com/nodeloc/git-store/internal/models"
	"github.com/nodeloc/git-store/internal/services"
	"github.com/nodeloc/git-store/internal/utils"
	"github.com/stripe/stripe-go/v76"
	"gorm.io/gorm"
)

// PluginHandler handles plugin-related requests
type PluginHandler struct {
	db     *gorm.DB
	config *config.Config
}

func NewPluginHandler(db *gorm.DB, cfg *config.Config) *PluginHandler {
	return &PluginHandler{db: db, config: cfg}
}

func (h *PluginHandler) ListPlugins(c *gin.Context) {
	var plugins []models.Plugin
	if err := h.db.Where("status = ?", "published").Order("created_at DESC").Find(&plugins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plugins"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"plugins": plugins})
}

func (h *PluginHandler) GetPlugin(c *gin.Context) {
	slug := c.Param("slug")
	var plugin models.Plugin
	if err := h.db.Where("slug = ? AND status = ?", slug, "published").First(&plugin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"plugin": plugin})
}

func (h *PluginHandler) GetPluginByID(c *gin.Context) {
	id := c.Param("id")
	pluginUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plugin ID"})
		return
	}

	var plugin models.Plugin
	if err := h.db.Where("id = ? AND status = ?", pluginUUID, "published").First(&plugin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"plugin": plugin})
}

// OrderHandler handles order-related requests
type OrderHandler struct {
	db     *gorm.DB
	config *config.Config
}

func NewOrderHandler(db *gorm.DB, cfg *config.Config) *OrderHandler {
	return &OrderHandler{db: db, config: cfg}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		PluginID      string `json:"plugin_id" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pluginUUID, err := uuid.Parse(req.PluginID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plugin ID"})
		return
	}

	var plugin models.Plugin
	if err := h.db.First(&plugin, "id = ?", pluginUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}

	order := models.Order{
		OrderNumber:   fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
		UserID:        userID.(uuid.UUID),
		PluginID:      pluginUUID,
		Amount:        plugin.Price,
		Currency:      plugin.Currency,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: "pending",
		Metadata:      "{}",
	}

	if err := h.db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order": order})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orderID := c.Param("id")

	var order models.Order
	if err := h.db.Preload("Plugin").Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var orders []models.Order
	if err := h.db.Preload("Plugin").Preload("License").Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	db            *gorm.DB
	config        *config.Config
	alipayService *services.AlipayService
	stripeService *services.StripeService
}

func NewPaymentHandler(db *gorm.DB, cfg *config.Config) *PaymentHandler {
	// 初始化易支付服务
	alipayService, err := services.NewAlipayService(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize Alipay service: %v", err)
	}

	// 初始化Stripe服务
	stripeService := services.NewStripeService(cfg)

	return &PaymentHandler{
		db:            db,
		config:        cfg,
		alipayService: alipayService,
		stripeService: stripeService,
	}
}

func (h *PaymentHandler) CreateStripePaymentIntent(c *gin.Context) {
	var req struct {
		OrderID string `json:"order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderUUID, err := uuid.Parse(req.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Verify order exists and belongs to user
	userID, _ := c.Get("user_id")
	var order models.Order
	if err := h.db.Preload("Plugin").Where("id = ? AND user_id = ?", orderUUID, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Check if Stripe service is available
	if h.stripeService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Stripe service is not configured",
		})
		return
	}

	// Check if Stripe is properly configured
	if h.config.StripeSecretKey == "" || h.config.StripeSecretKey == "sk_test_your_stripe_secret_key" {
		log.Printf("Stripe is not configured properly. Please set STRIPE_SECRET_KEY in .env file")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Payment service is not configured. Please contact administrator.",
		})
		return
	}

	// Create Stripe payment intent
	amountCents := int64(order.Amount * 100) // Convert to cents
	paymentReq := &services.PaymentIntentRequest{
		Amount:      amountCents,
		Currency:    "usd",
		Description: fmt.Sprintf("%s - License", order.Plugin.Name),
		Metadata: map[string]string{
			"order_id":  order.ID.String(),
			"user_id":   order.UserID.String(),
			"plugin_id": order.PluginID.String(),
		},
	}

	paymentIntent, err := h.stripeService.CreatePaymentIntent(paymentReq)
	if err != nil {
		log.Printf("Failed to create Stripe payment intent: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create payment",
			"details": "Unable to initialize payment. Please try again or contact support.",
		})
		return
	}

	// Update order with payment intent ID
	order.PaymentIntentID = paymentIntent.ID
	if err := h.db.Save(&order).Error; err != nil {
		log.Printf("Failed to update order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"client_secret":     paymentIntent.ClientSecret,
		"payment_intent_id": paymentIntent.ID,
		"order_id":          order.ID,
		"amount":            order.Amount,
	})
}

func (h *PaymentHandler) CreatePayPalOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Create PayPal order - to be implemented"})
}

func (h *PaymentHandler) CapturePayPalOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Capture PayPal order - to be implemented"})
}

func (h *PaymentHandler) CreateAlipayPayment(c *gin.Context) {
	var req struct {
		OrderID string `json:"order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderUUID, err := uuid.Parse(req.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Verify order exists and belongs to user
	userID, _ := c.Get("user_id")
	var order models.Order
	if err := h.db.Preload("Plugin").Where("id = ? AND user_id = ?", orderUUID, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// 检查易支付服务是否可用
	if h.alipayService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Alipay service is not configured",
		})
		return
	}

	// 货币转换：如果订单货币不是 CNY，需要转换
	paymentAmount := order.Amount
	if order.Currency != "CNY" {
		// 使用汇率服务转换金额
		exchangeRateSvc := services.NewExchangeRateService(h.db, h.config)
		convertedAmount, err := exchangeRateSvc.ConvertAmount(order.Amount, order.Currency, "CNY")
		if err != nil {
			log.Printf("❌ 货币转换失败 %s -> CNY: %v", order.Currency, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Currency conversion failed: %s to CNY", order.Currency),
			})
			return
		}
		paymentAmount = convertedAmount
		log.Printf("💱 货币转换: %.2f %s = %.2f CNY", order.Amount, order.Currency, paymentAmount)
	}

	// 获取客户端IP
	clientIP := c.ClientIP()
	if clientIP == "" || clientIP == "::1" {
		clientIP = "127.0.0.1"
	}

	// 创建易支付订单
	paymentReq := &services.AlipayTradeRequest{
		OutTradeNo:  order.ID.String(),
		TotalAmount: paymentAmount, // 使用转换后的 CNY 金额
		Subject:     fmt.Sprintf("%s - License", order.Plugin.Name),
		Body:        fmt.Sprintf("Order ID: %s", order.ID.String()),
		NotifyURL:   h.config.AppURL + "/api/webhooks/alipay",
		ReturnURL:   h.config.FrontendURL + "/payment/success",
		ClientIP:    clientIP,
	}

	result, err := h.alipayService.CreatePayment(paymentReq)
	if err != nil {
		log.Printf("Failed to create Alipay payment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	// 构建响应，根据返回的字段类型返回支付信息
	response := gin.H{
		"trade_no":         result.TradeNo,
		"order_id":         order.ID,
		"amount":           order.Amount,   // 原始金额
		"currency":         order.Currency, // 原始货币
		"payment_amount":   paymentAmount,  // CNY 支付金额
		"payment_currency": "CNY",
	}

	// 优先使用 PayInfo 字段（易支付实际返回的字段）
	// 根据 pay_type 判断支付方式
	payURL := result.PayInfo
	if payURL == "" {
		payURL = result.PayURL // fallback到新版字段
	}

	if payURL != "" {
		response["pay_url"] = payURL
		response["pay_type"] = result.PayType
	} else if result.QRCode != "" {
		response["pay_url"] = result.QRCode
		response["pay_type"] = "qrcode"
	} else if result.URLScheme != "" {
		response["pay_url"] = result.URLScheme
		response["pay_type"] = "urlscheme"
	}

	c.JSON(http.StatusOK, response)
}

func (h *PaymentHandler) StripeWebhook(c *gin.Context) {
	// Check if Stripe service is available
	if h.stripeService == nil {
		log.Printf("Stripe service is not configured")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Stripe service not configured"})
		return
	}

	// Read request body
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get Stripe signature header
	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		log.Printf("Missing Stripe-Signature header")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
		return
	}

	// Verify webhook signature
	event, err := h.stripeService.VerifyWebhookSignature(payload, signature)
	if err != nil {
		log.Printf("Failed to verify webhook signature: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	// Handle the event
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			log.Printf("Failed to parse payment intent: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
			return
		}

		// Get order ID from metadata
		orderIDStr, ok := paymentIntent.Metadata["order_id"]
		if !ok {
			log.Printf("Order ID not found in payment intent metadata")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID not found"})
			return
		}

		orderUUID, err := uuid.Parse(orderIDStr)
		if err != nil {
			log.Printf("Invalid order ID: %s", orderIDStr)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		// Find order
		var order models.Order
		if err := h.db.Preload("Plugin").Where("id = ?", orderUUID).First(&order).Error; err != nil {
			log.Printf("Order not found: %s", orderIDStr)
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		// Check if already processed
		if order.PaymentStatus == "paid" {
			log.Printf("Order already completed: %s", orderIDStr)
			c.JSON(http.StatusOK, gin.H{"message": "Already processed"})
			return
		}

		// Update order status
		order.PaymentStatus = "paid"
		order.PaymentMethod = "stripe"
		order.PaymentTransactionID = paymentIntent.ID
		if err := h.db.Save(&order).Error; err != nil {
			log.Printf("Failed to update order: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
			return
		}

		// Get user's first GitHub account
		var githubAccount models.GitHubAccount
		if err := h.db.Where("user_id = ?", order.UserID).First(&githubAccount).Error; err != nil {
			log.Printf("Failed to find GitHub account for user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub account not found"})
			return
		}

		// Calculate maintenance expiry
		maintenanceMonths := order.Plugin.DefaultMaintenanceMonths
		if maintenanceMonths == 0 {
			maintenanceMonths = 12
		}
		maintenanceUntil := time.Now().AddDate(0, maintenanceMonths, 0)

		// Create license
		license := models.License{
			OrderID:          order.ID,
			UserID:           order.UserID,
			PluginID:         order.PluginID,
			GitHubAccountID:  githubAccount.ID,
			LicenseType:      "permanent",
			Status:           "active",
			MaintenanceUntil: maintenanceUntil,
		}

		if err := h.db.Create(&license).Error; err != nil {
			log.Printf("Failed to create license: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create license"})
			return
		}

		log.Printf("Successfully processed Stripe payment for order: %s", orderIDStr)

	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			log.Printf("Failed to parse payment intent: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
			return
		}

		orderIDStr, ok := paymentIntent.Metadata["order_id"]
		if ok {
			orderUUID, err := uuid.Parse(orderIDStr)
			if err == nil {
				var order models.Order
				if err := h.db.Where("id = ?", orderUUID).First(&order).Error; err == nil {
					order.PaymentStatus = "failed"
					h.db.Save(&order)
				}
			}
		}

		log.Printf("Stripe payment failed for order: %s", orderIDStr)

	default:
		log.Printf("Unhandled Stripe event type: %s", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *PaymentHandler) PayPalWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PayPal webhook - to be implemented"})
}

func (h *PaymentHandler) AlipayNotify(c *gin.Context) {
	// 检查易支付服务是否可用
	if h.alipayService == nil {
		log.Printf("Alipay service is not configured")
		c.String(http.StatusServiceUnavailable, "fail")
		return
	}

	// 获取所有POST参数
	if err := c.Request.ParseForm(); err != nil {
		log.Printf("Failed to parse form: %v", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 转换为map[string]string
	params := make(map[string]string)
	for key, values := range c.Request.PostForm {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	// 验证签名
	if err := h.alipayService.VerifyNotify(params); err != nil {
		log.Printf("Failed to verify Alipay signature: %v", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 获取关键参数
	tradeStatus := params["trade_status"] // 订单状态
	outTradeNo := params["out_trade_no"]  // 商户订单号
	tradeNo := params["trade_no"]         // 平台订单号

	// 解析订单ID
	orderUUID, err := uuid.Parse(outTradeNo)
	if err != nil {
		log.Printf("Invalid order ID: %s", outTradeNo)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 查找订单
	var order models.Order
	if err := h.db.Where("id = ?", orderUUID).First(&order).Error; err != nil {
		log.Printf("Order not found: %s", outTradeNo)
		c.String(http.StatusNotFound, "fail")
		return
	}

	// 处理支付成功
	if tradeStatus == "TRADE_SUCCESS" || tradeStatus == "1" {
		// 检查订单状态，避免重复处理
		if order.PaymentStatus == "paid" {
			log.Printf("Order already completed: %s", outTradeNo)
			c.String(http.StatusOK, "success")
			return
		}

		// 更新订单状态
		order.PaymentStatus = "paid"
		order.PaymentMethod = "alipay"
		if err := h.db.Save(&order).Error; err != nil {
			log.Printf("Failed to update order: %v", err)
			c.String(http.StatusInternalServerError, "fail")
			return
		}

		// 获取用户的第一个GitHub账户
		var githubAccount models.GitHubAccount
		if err := h.db.Where("user_id = ?", order.UserID).First(&githubAccount).Error; err != nil {
			log.Printf("Failed to find GitHub account for user: %v", err)
			c.String(http.StatusInternalServerError, "fail")
			return
		}

		// 计算维护到期时间（默认12个月）
		maintenanceMonths := order.Plugin.DefaultMaintenanceMonths
		if maintenanceMonths == 0 {
			maintenanceMonths = 12
		}

		// 生成许可证
		license := models.License{
			UserID:           order.UserID,
			PluginID:         order.PluginID,
			OrderID:          order.ID,
			GitHubAccountID:  githubAccount.ID,
			LicenseType:      "permanent",
			MaintenanceUntil: utils.CalculateMaintenanceUntil(maintenanceMonths),
			Status:           "active",
		}

		if err := h.db.Create(&license).Error; err != nil {
			log.Printf("Failed to create license: %v", err)
			c.String(http.StatusInternalServerError, "fail")
			return
		}

		log.Printf("Payment successful - Order: %s, TradeNo: %s", outTradeNo, tradeNo)
		c.String(http.StatusOK, "success")
		return
	}

	// 其他状态，记录日志
	log.Printf("Payment status: %s, Order: %s, TradeNo: %s", tradeStatus, outTradeNo, tradeNo)
	c.String(http.StatusOK, "success")
}

// LicenseHandler handles license-related requests
type LicenseHandler struct {
	db     *gorm.DB
	config *config.Config
}

func NewLicenseHandler(db *gorm.DB, cfg *config.Config) *LicenseHandler {
	return &LicenseHandler{db: db, config: cfg}
}

func (h *LicenseHandler) GetUserLicenses(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var licenses []models.License
	if err := h.db.Preload("Plugin").Preload("GitHubAccount").Where("user_id = ?", userID).Order("created_at DESC").Find(&licenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch licenses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"licenses": licenses})
}

func (h *LicenseHandler) GetLicense(c *gin.Context) {
	userID, _ := c.Get("user_id")
	licenseID := c.Param("id")

	var license models.License
	if err := h.db.Preload("Plugin").Preload("GitHubAccount").Preload("History").Where("id = ? AND user_id = ?", licenseID, userID).First(&license).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"license": license})
}

func (h *LicenseHandler) RenewLicense(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Renew license - to be implemented"})
}

func (h *LicenseHandler) GetLicenseHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	licenseID := c.Param("id")

	var license models.License
	if err := h.db.Where("id = ? AND user_id = ?", licenseID, userID).First(&license).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
		return
	}

	var history []models.LicenseHistory
	if err := h.db.Where("license_id = ?", licenseID).Order("occurred_at DESC").Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

// VerifyLicense is a public API to verify license validity
// Can be called by plugins, GitHub Apps, CI/CD, etc.
func (h *LicenseHandler) VerifyLicense(c *gin.Context) {
	licenseID := c.Param("id")
	githubUsername := c.Query("github_username") // Optional: verify GitHub username matches

	var license models.License
	if err := h.db.Preload("Plugin").Preload("GitHubAccount").Preload("User").Where("id = ?", licenseID).First(&license).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"valid":   false,
			"error":   "License not found",
			"message": "Invalid license ID",
		})
		return
	}

	// Check if license is active
	if license.Status != "active" {
		c.JSON(http.StatusOK, gin.H{
			"valid":   false,
			"status":  license.Status,
			"reason":  license.RevokedReason,
			"message": "License is not active",
		})
		return
	}

	// Check maintenance period
	now := time.Now()
	maintenanceExpired := license.MaintenanceUntil.Before(now)

	// If GitHub username provided, verify it matches
	if githubUsername != "" && license.GitHubAccount.Login != githubUsername {
		c.JSON(http.StatusOK, gin.H{
			"valid":   false,
			"error":   "GitHub username mismatch",
			"message": "This license is bound to a different GitHub account",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":              true,
		"license_id":         license.ID,
		"license_type":       license.LicenseType,
		"status":             license.Status,
		"plugin_name":        license.Plugin.Name,
		"plugin_slug":        license.Plugin.Slug,
		"github_username":    license.GitHubAccount.Login,
		"maintenance_until":  license.MaintenanceUntil,
		"maintenance_active": !maintenanceExpired,
		"created_at":         license.CreatedAt,
	})
}

// TutorialHandler handles tutorial-related requests
type TutorialHandler struct {
	db     *gorm.DB
	config *config.Config
}

func NewTutorialHandler(db *gorm.DB, cfg *config.Config) *TutorialHandler {
	return &TutorialHandler{db: db, config: cfg}
}

func (h *TutorialHandler) ListPublicTutorials(c *gin.Context) {
	var tutorials []models.Tutorial
	if err := h.db.Preload("Plugin").Where("is_public = ?", true).Order("order_index ASC").Find(&tutorials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tutorials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tutorials": tutorials})
}

func (h *TutorialHandler) ListTutorials(c *gin.Context) {
	var tutorials []models.Tutorial
	if err := h.db.Preload("Plugin").Order("order_index ASC").Find(&tutorials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tutorials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tutorials": tutorials})
}

func (h *TutorialHandler) GetTutorial(c *gin.Context) {
	slug := c.Param("slug")
	var tutorial models.Tutorial
	if err := h.db.Preload("Plugin").Where("slug = ?", slug).First(&tutorial).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tutorial not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tutorial": tutorial})
}

// AdminHandler handles admin-related requests
type AdminHandler struct {
	db        *gorm.DB
	config    *config.Config
	githubSvc *services.GitHubService
}

func NewAdminHandler(db *gorm.DB, cfg *config.Config, githubSvc *services.GitHubService) *AdminHandler {
	return &AdminHandler{
		db:        db,
		config:    cfg,
		githubSvc: githubSvc,
	}
}

// ==================== Plugin Management ====================

func (h *AdminHandler) ListAllPlugins(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var plugins []models.Plugin
	var total int64

	query := h.db.Model(&models.Plugin{})

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&plugins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plugins"})
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(http.StatusOK, gin.H{
		"plugins": plugins,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *AdminHandler) CreatePlugin(c *gin.Context) {
	var req struct {
		Name                     string   `json:"name" binding:"required"`
		Slug                     string   `json:"slug" binding:"required"`
		Description              string   `json:"description"`
		LongDescription          string   `json:"long_description"`
		GitHubRepoID             int64    `json:"github_repo_id"`
		GitHubRepoURL            string   `json:"github_repo_url"`
		GitHubRepoName           string   `json:"github_repo_name"`
		Price                    float64  `json:"price"`
		Currency                 string   `json:"currency"`
		DefaultMaintenanceMonths int      `json:"default_maintenance_months"`
		Status                   string   `json:"status"`
		Category                 string   `json:"category"`
		Tags                     []string `json:"tags"`
		IconURL                  string   `json:"icon_url"`
		DemoURL                  string   `json:"demo_url"`
		DocumentationURL         string   `json:"documentation_url"`
		Version                  string   `json:"version"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values for required GitHub fields if not provided
	if req.GitHubRepoID == 0 {
		req.GitHubRepoID = 0 // Use 0 as placeholder
	}
	if req.GitHubRepoURL == "" {
		req.GitHubRepoURL = ""
	}
	if req.GitHubRepoName == "" {
		req.GitHubRepoName = ""
	}

	plugin := models.Plugin{
		Name:                     req.Name,
		Slug:                     req.Slug,
		Description:              req.Description,
		LongDescription:          req.LongDescription,
		GitHubRepoID:             req.GitHubRepoID,
		GitHubRepoURL:            req.GitHubRepoURL,
		GitHubRepoName:           req.GitHubRepoName,
		Price:                    req.Price,
		Currency:                 req.Currency,
		DefaultMaintenanceMonths: req.DefaultMaintenanceMonths,
		Status:                   req.Status,
		Category:                 req.Category,
		Tags:                     req.Tags,
		IconURL:                  req.IconURL,
		DemoURL:                  req.DemoURL,
		DocumentationURL:         req.DocumentationURL,
		Version:                  req.Version,
	}

	if plugin.Currency == "" {
		plugin.Currency = "USD"
	}
	if plugin.Status == "" {
		plugin.Status = "draft"
	}
	if plugin.DefaultMaintenanceMonths == 0 {
		plugin.DefaultMaintenanceMonths = h.config.DefaultMaintenanceMonths
	}

	if err := h.db.Create(&plugin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plugin"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"plugin": plugin})
}

func (h *AdminHandler) GetPluginByID(c *gin.Context) {
	id := c.Param("id")

	var plugin models.Plugin
	if err := h.db.First(&plugin, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plugin": plugin})
}

func (h *AdminHandler) UpdatePlugin(c *gin.Context) {
	id := c.Param("id")

	var plugin models.Plugin
	if err := h.db.First(&plugin, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}

	var req struct {
		Name                     string   `json:"name"`
		Slug                     string   `json:"slug"`
		Description              string   `json:"description"`
		LongDescription          string   `json:"long_description"`
		GitHubRepoID             int64    `json:"github_repo_id"`
		GitHubRepoURL            string   `json:"github_repo_url"`
		GitHubRepoName           string   `json:"github_repo_name"`
		Price                    float64  `json:"price"`
		Currency                 string   `json:"currency"`
		DefaultMaintenanceMonths int      `json:"default_maintenance_months"`
		Status                   string   `json:"status"`
		Category                 string   `json:"category"`
		Tags                     []string `json:"tags"`
		IconURL                  string   `json:"icon_url"`
		DemoURL                  string   `json:"demo_url"`
		DocumentationURL         string   `json:"documentation_url"`
		Version                  string   `json:"version"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"name":                       req.Name,
		"slug":                       req.Slug,
		"description":                req.Description,
		"long_description":           req.LongDescription,
		"github_repo_id":             req.GitHubRepoID,
		"github_repo_url":            req.GitHubRepoURL,
		"github_repo_name":           req.GitHubRepoName,
		"price":                      req.Price,
		"currency":                   req.Currency,
		"default_maintenance_months": req.DefaultMaintenanceMonths,
		"status":                     req.Status,
		"category":                   req.Category,
		"tags":                       req.Tags,
		"icon_url":                   req.IconURL,
		"demo_url":                   req.DemoURL,
		"documentation_url":          req.DocumentationURL,
		"version":                    req.Version,
	}

	if err := h.db.Model(&plugin).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plugin"})
		return
	}

	h.db.First(&plugin, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"plugin": plugin})
}

func (h *AdminHandler) DeletePlugin(c *gin.Context) {
	id := c.Param("id")

	var plugin models.Plugin
	if err := h.db.First(&plugin, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}

	var orderCount int64
	h.db.Model(&models.Order{}).Where("plugin_id = ?", id).Count(&orderCount)
	if orderCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete plugin with existing orders"})
		return
	}

	if err := h.db.Delete(&plugin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plugin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plugin deleted successfully"})
}

// ListGitHubRepos lists all repositories from GitHub (using configured account)
func (h *AdminHandler) ListGitHubRepos(c *gin.Context) {
	if h.githubSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "GitHub service not configured"})
		return
	}

	ctx := c.Request.Context()

	// Get owner from query param, default to "nodeloc"
	owner := c.DefaultQuery("owner", "nodeloc")

	// List repositories for the specified owner
	log.Printf("[GitHub API] Fetching repositories for owner: %s", owner)
	repos, err := h.githubSvc.ListRepositories(ctx, owner)
	if err != nil {
		log.Printf("[GitHub API] Error fetching repositories: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch GitHub repositories: " + err.Error()})
		return
	}
	log.Printf("[GitHub API] Successfully fetched %d repositories", len(repos))

	// Transform to simpler format for frontend
	type RepoInfo struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		HTMLURL     string `json:"html_url"`
		Private     bool   `json:"private"`
		Fork        bool   `json:"fork"`
		Archived    bool   `json:"archived"`
		Stars       int    `json:"stars"`
		Language    string `json:"language"`
	}

	var repoList []RepoInfo
	for _, repo := range repos {
		if repo == nil {
			continue
		}

		// Skip forked repositories
		if repo.GetFork() {
			continue
		}

		description := ""
		if repo.Description != nil {
			description = *repo.Description
		}

		language := ""
		if repo.Language != nil {
			language = *repo.Language
		}

		repoList = append(repoList, RepoInfo{
			ID:          repo.GetID(),
			Name:        repo.GetName(),
			FullName:    repo.GetFullName(),
			Description: description,
			URL:         repo.GetURL(),
			HTMLURL:     repo.GetHTMLURL(),
			Private:     repo.GetPrivate(),
			Fork:        repo.GetFork(),
			Archived:    repo.GetArchived(),
			Stars:       repo.GetStargazersCount(),
			Language:    language,
		})
	}

	c.JSON(http.StatusOK, gin.H{"repositories": repoList})
}

func (h *AdminHandler) SyncGitHubRepos(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Sync GitHub repos - requires GitHub App configuration"})
}

// ==================== Order Management ====================

func (h *AdminHandler) ListAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var orders []models.Order
	var total int64

	query := h.db.Model(&models.Order{})

	if search != "" {
		query = query.Where("order_number ILIKE ?", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("payment_status = ?", status)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	query.Count(&total)

	// 预加载关联数据
	if err := query.Preload("User").Preload("Plugin").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	// 调试：检查是否正确加载了关联数据
	if len(orders) > 0 {
		log.Printf("[Admin Debug] First order: UserID=%s, PluginID=%s, User.Email=%s, Plugin.Name=%s",
			orders[0].UserID, orders[0].PluginID, orders[0].User.Email, orders[0].Plugin.Name)
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *AdminHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	if err := h.db.Preload("User").Preload("Plugin").First(&order, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// UpdateOrderPaymentStatus updates the payment status of an order (admin only)
func (h *AdminHandler) UpdateOrderPaymentStatus(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	var req struct {
		PaymentStatus string `json:"payment_status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate payment status
	validStatuses := map[string]bool{
		"pending":  true,
		"paid":     true,
		"failed":   true,
		"refunded": true,
	}

	if !validStatuses[req.PaymentStatus] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment status"})
		return
	}

	var order models.Order
	if err := h.db.Where("id = ?", orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order"})
		return
	}

	updates := map[string]interface{}{
		"payment_status": req.PaymentStatus,
	}

	// If setting to paid, update paid_at timestamp
	if req.PaymentStatus == "paid" && order.PaidAt == nil {
		now := time.Now()
		updates["paid_at"] = &now
	}

	// If setting to refunded, update refunded_at timestamp
	if req.PaymentStatus == "refunded" && order.RefundedAt == nil {
		now := time.Now()
		updates["refunded_at"] = &now
	}

	if err := h.db.Model(&order).Updates(updates).Error; err != nil {
		log.Printf("Failed to update order payment status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order payment status"})
		return
	}

	log.Printf("[License Debug] After update - req.PaymentStatus=%s, order.ID=%s", req.PaymentStatus, order.ID)

	// If status is paid, create license if it doesn't exist
	if req.PaymentStatus == "paid" {
		log.Printf("[License Debug] Payment status is paid, checking for existing license")
		// Check if license already exists
		var existingLicense models.License
		err := h.db.Where("order_id = ?", order.ID).First(&existingLicense).Error
		log.Printf("[License Debug] License query result - err: %v", err)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			// License doesn't exist, create one

			// Get user's GitHub account (required)
			var githubAccount models.GitHubAccount
			if err := h.db.Where("user_id = ?", order.UserID).First(&githubAccount).Error; err != nil {
				log.Printf("Cannot create license: user %s has no GitHub account: %v", order.UserID, err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error":           "Cannot create license: user must login with GitHub first to link their account",
					"order_updated":   true,
					"license_created": false,
				})
				return
			}

			// Load plugin to get maintenance months
			var plugin models.Plugin
			if err := h.db.Where("id = ?", order.PluginID).First(&plugin).Error; err != nil {
				log.Printf("Failed to find plugin %s: %v", order.PluginID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Order updated but failed to find plugin"})
				return
			}

			// Calculate maintenance expiry
			maintenanceMonths := plugin.DefaultMaintenanceMonths
			if maintenanceMonths == 0 {
				maintenanceMonths = 12
			}
			maintenanceUntil := time.Now().AddDate(0, maintenanceMonths, 0)

			// Create license
			license := models.License{
				OrderID:          order.ID,
				UserID:           order.UserID,
				PluginID:         order.PluginID,
				GitHubAccountID:  githubAccount.ID,
				LicenseType:      "permanent",
				Status:           "active",
				MaintenanceUntil: maintenanceUntil,
			}

			if err := h.db.Create(&license).Error; err != nil {
				log.Printf("Failed to create license for order %s: %v", order.ID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Order updated but failed to create license"})
				return
			}

			log.Printf("Successfully created license for order %s", order.ID)

			// Grant repository access via collaborator invitation
			if h.githubSvc != nil && plugin.GitHubRepoName != "" && githubAccount.Login != "" {
				ctx := context.Background()

				// Split repo name into owner/repo
				repoParts := strings.Split(plugin.GitHubRepoName, "/")
				if len(repoParts) != 2 {
					log.Printf("[Repository Access] Invalid repo name format: %s", plugin.GitHubRepoName)
				} else {
					owner, repo := repoParts[0], repoParts[1]

					log.Printf("[Repository Access] Inviting %s as collaborator to %s with pull permission",
						githubAccount.Login, plugin.GitHubRepoName)

					err := h.githubSvc.AddRepositoryCollaborator(ctx, owner, repo, githubAccount.Login, "pull")
					if err != nil {
						log.Printf("[Repository Access] Warning: Failed to add collaborator: %v", err)
					} else {
						log.Printf("[Repository Access] Successfully invited %s as collaborator to %s",
							githubAccount.Login, plugin.GitHubRepoName)

						// Log the action
						history := models.LicenseHistory{
							LicenseID:  license.ID,
							Action:     "collaborator_invited",
							OccurredAt: time.Now(),
						}
						h.db.Create(&history)
					}
				}
			}
		}
	}

	// Reload order with associations
	if err := h.db.Preload("User").Preload("Plugin").Where("id = ?", orderID).First(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *AdminHandler) RefundOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	if err := h.db.First(&order, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.PaymentStatus != "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order is not in paid status"})
		return
	}

	now := time.Now()
	order.PaymentStatus = "refunded"
	order.RefundedAt = &now

	if err := h.db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	h.db.Model(&models.License{}).Where("order_id = ?", order.ID).Updates(map[string]interface{}{
		"status":         "revoked",
		"revoked_reason": "Order refunded",
		"revoked_at":     now,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Order refunded successfully", "order": order})
}

// ==================== License Management ====================

func (h *AdminHandler) ListAllLicenses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	pluginID := c.Query("plugin_id")
	userID := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var licenses []models.License
	var total int64

	query := h.db.Model(&models.License{}).Preload("User").Preload("Plugin").Preload("GitHubAccount")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if pluginID != "" {
		query = query.Where("plugin_id = ?", pluginID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&licenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch licenses"})
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(http.StatusOK, gin.H{
		"licenses": licenses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *AdminHandler) GetLicenseByID(c *gin.Context) {
	id := c.Param("id")

	var license models.License
	if err := h.db.Preload("User").Preload("Plugin").Preload("GitHubAccount").Preload("History").First(&license, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"license": license})
}

func (h *AdminHandler) RevokeLicense(c *gin.Context) {
	id := c.Param("id")
	userIDValue, _ := c.Get("user_id")
	adminUserID := userIDValue.(uuid.UUID)

	var license models.License
	if err := h.db.First(&license, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	license.Status = "revoked"
	license.RevokedReason = req.Reason
	license.RevokedAt = &now

	if err := h.db.Save(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke license"})
		return
	}

	history := models.LicenseHistory{
		LicenseID:   license.ID,
		Action:      "revoked",
		PerformedBy: &adminUserID,
		Metadata:    fmt.Sprintf(`{"reason": "%s"}`, req.Reason),
		OccurredAt:  now,
	}
	h.db.Create(&history)

	c.JSON(http.StatusOK, gin.H{"message": "License revoked successfully", "license": license})
}

func (h *AdminHandler) ExtendLicense(c *gin.Context) {
	id := c.Param("id")
	userIDValue, _ := c.Get("user_id")
	adminUserID := userIDValue.(uuid.UUID)

	var license models.License
	if err := h.db.First(&license, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
		return
	}

	var req struct {
		Months int `json:"months" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseDate := license.MaintenanceUntil
	if time.Now().After(baseDate) {
		baseDate = time.Now()
	}
	license.MaintenanceUntil = baseDate.AddDate(0, req.Months, 0)

	if license.Status == "expired" {
		license.Status = "active"
	}

	if err := h.db.Save(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extend license"})
		return
	}

	history := models.LicenseHistory{
		LicenseID:   license.ID,
		Action:      "renewed",
		PerformedBy: &adminUserID,
		Metadata:    fmt.Sprintf(`{"months": %d, "new_maintenance_until": "%s"}`, req.Months, license.MaintenanceUntil.Format("2006-01-02")),
		OccurredAt:  time.Now(),
	}
	h.db.Create(&history)

	c.JSON(http.StatusOK, gin.H{"message": "License extended successfully", "license": license})
}

// ==================== Tutorial Management ====================

func (h *AdminHandler) ListAllTutorials(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	pluginID := c.Query("plugin_id")
	language := c.Query("language")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var tutorials []models.Tutorial
	var total int64

	query := h.db.Model(&models.Tutorial{}).Preload("Plugin")

	if pluginID != "" {
		query = query.Where("plugin_id = ?", pluginID)
	}
	if language != "" {
		query = query.Where("language = ?", language)
	}

	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("order_index ASC, created_at DESC").Find(&tutorials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tutorials"})
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(http.StatusOK, gin.H{
		"tutorials": tutorials,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *AdminHandler) CreateTutorial(c *gin.Context) {
	var req struct {
		PluginID    *uuid.UUID `json:"plugin_id"`
		Title       string     `json:"title" binding:"required"`
		Slug        string     `json:"slug" binding:"required"`
		Content     string     `json:"content" binding:"required"`
		ContentType string     `json:"content_type"`
		OrderIndex  int        `json:"order_index"`
		IsPublic    bool       `json:"is_public"`
		Language    string     `json:"language"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tutorial := models.Tutorial{
		PluginID:    req.PluginID,
		Title:       req.Title,
		Slug:        req.Slug,
		Content:     req.Content,
		ContentType: req.ContentType,
		OrderIndex:  req.OrderIndex,
		IsPublic:    req.IsPublic,
		Language:    req.Language,
	}

	if tutorial.ContentType == "" {
		tutorial.ContentType = "markdown"
	}
	if tutorial.Language == "" {
		tutorial.Language = "en"
	}

	if err := h.db.Create(&tutorial).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tutorial"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"tutorial": tutorial})
}

func (h *AdminHandler) GetTutorialByID(c *gin.Context) {
	id := c.Param("id")

	var tutorial models.Tutorial
	if err := h.db.Preload("Plugin").First(&tutorial, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tutorial not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tutorial": tutorial})
}

func (h *AdminHandler) UpdateTutorial(c *gin.Context) {
	id := c.Param("id")

	var tutorial models.Tutorial
	if err := h.db.First(&tutorial, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tutorial not found"})
		return
	}

	var req struct {
		PluginID    *uuid.UUID `json:"plugin_id"`
		Title       string     `json:"title"`
		Slug        string     `json:"slug"`
		Content     string     `json:"content"`
		ContentType string     `json:"content_type"`
		OrderIndex  int        `json:"order_index"`
		IsPublic    bool       `json:"is_public"`
		Language    string     `json:"language"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"plugin_id":    req.PluginID,
		"title":        req.Title,
		"slug":         req.Slug,
		"content":      req.Content,
		"content_type": req.ContentType,
		"order_index":  req.OrderIndex,
		"is_public":    req.IsPublic,
		"language":     req.Language,
	}

	if err := h.db.Model(&tutorial).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tutorial"})
		return
	}

	h.db.First(&tutorial, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"tutorial": tutorial})
}

func (h *AdminHandler) DeleteTutorial(c *gin.Context) {
	id := c.Param("id")

	var tutorial models.Tutorial
	if err := h.db.First(&tutorial, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tutorial not found"})
		return
	}

	if err := h.db.Delete(&tutorial).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tutorial"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tutorial deleted successfully"})
}

// ==================== Settings Management ====================

func (h *AdminHandler) GetSettings(c *gin.Context) {
	var settings []models.SystemSetting
	if err := h.db.Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

func (h *AdminHandler) UpdateSettings(c *gin.Context) {
	var req struct {
		Settings []struct {
			Key   string `json:"key" binding:"required"`
			Value string `json:"value" binding:"required"`
		} `json:"settings" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, setting := range req.Settings {
		if err := h.db.Model(&models.SystemSetting{}).Where("key = ?", setting.Key).Updates(map[string]interface{}{
			"value":      setting.Value,
			"updated_at": time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update setting: " + setting.Key})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

// GetPublicSettings returns public settings that can be accessed without authentication
func (h *AdminHandler) GetPublicSettings(c *gin.Context) {
	var settings []models.SystemSetting
	// Only return specific public settings
	publicKeys := []string{"site_name", "site_subtitle", "logo_url"}

	if err := h.db.Where("key IN ?", publicKeys).Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// ==================== User Management ====================

func (h *AdminHandler) ListAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	role := c.Query("role")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var users []models.User
	var total int64

	query := h.db.Model(&models.User{}).Preload("GitHubAccounts")

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}

	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *AdminHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.db.Preload("GitHubAccounts").Preload("Orders").Preload("Licenses").First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.db.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		Name     string `json:"name"`
		Role     string `json:"role"`
		IsActive *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := h.db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	h.db.First(&user, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userIDValue, _ := c.Get("user_id")
	currentUserID := userIDValue.(uuid.UUID)

	if id == currentUserID.String() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete yourself"})
		return
	}

	var user models.User
	if err := h.db.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := h.db.Model(&user).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deactivated successfully"})
}

// DashboardHandler handles dashboard statistics
type DashboardHandler struct {
	db     *gorm.DB
	config *config.Config
}

func NewDashboardHandler(db *gorm.DB, cfg *config.Config) *DashboardHandler {
	return &DashboardHandler{db: db, config: cfg}
}

func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	var stats struct {
		TotalUsers     int64   `json:"total_users"`
		TotalPlugins   int64   `json:"total_plugins"`
		TotalOrders    int64   `json:"total_orders"`
		TotalRevenue   float64 `json:"total_revenue"`
		ActiveLicenses int64   `json:"active_licenses"`
		NewUsersToday  int64   `json:"new_users_today"`
		NewOrdersToday int64   `json:"new_orders_today"`
		RevenueToday   float64 `json:"revenue_today"`
	}

	h.db.Model(&models.User{}).Count(&stats.TotalUsers)
	h.db.Model(&models.Plugin{}).Where("status = ?", "published").Count(&stats.TotalPlugins)
	h.db.Model(&models.Order{}).Where("payment_status = ?", "paid").Count(&stats.TotalOrders)
	h.db.Model(&models.Order{}).Where("payment_status = ?", "paid").Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalRevenue)
	h.db.Model(&models.License{}).Where("status = ?", "active").Count(&stats.ActiveLicenses)

	today := time.Now().Truncate(24 * time.Hour)
	h.db.Model(&models.User{}).Where("created_at >= ?", today).Count(&stats.NewUsersToday)
	h.db.Model(&models.Order{}).Where("created_at >= ? AND payment_status = ?", today, "paid").Count(&stats.NewOrdersToday)
	h.db.Model(&models.Order{}).Where("created_at >= ? AND payment_status = ?", today, "paid").Select("COALESCE(SUM(amount), 0)").Scan(&stats.RevenueToday)

	c.JSON(http.StatusOK, stats)
}

func (h *DashboardHandler) GetRevenueStats(c *gin.Context) {
	period := c.DefaultQuery("period", "30d")

	var days int
	switch period {
	case "7d":
		days = 7
	case "30d":
		days = 30
	case "90d":
		days = 90
	case "1y":
		days = 365
	default:
		days = 30
	}

	startDate := time.Now().AddDate(0, 0, -days).Truncate(24 * time.Hour)

	type DailyRevenue struct {
		Date    string  `json:"date"`
		Revenue float64 `json:"revenue"`
		Orders  int     `json:"orders"`
	}

	var results []DailyRevenue
	h.db.Model(&models.Order{}).
		Select("TO_CHAR(paid_at, 'YYYY-MM-DD') as date, COALESCE(SUM(amount), 0) as revenue, COUNT(*) as orders").
		Where("payment_status = ? AND paid_at >= ?", "paid", startDate).
		Group("TO_CHAR(paid_at, 'YYYY-MM-DD')").
		Order("date ASC").
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{
		"period": period,
		"data":   results,
	})
}

func (h *DashboardHandler) GetUserStats(c *gin.Context) {
	period := c.DefaultQuery("period", "30d")

	var days int
	switch period {
	case "7d":
		days = 7
	case "30d":
		days = 30
	case "90d":
		days = 90
	case "1y":
		days = 365
	default:
		days = 30
	}

	startDate := time.Now().AddDate(0, 0, -days).Truncate(24 * time.Hour)

	type DailyUsers struct {
		Date     string `json:"date"`
		NewUsers int    `json:"new_users"`
	}

	var results []DailyUsers
	h.db.Model(&models.User{}).
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as date, COUNT(*) as new_users").
		Where("created_at >= ?", startDate).
		Group("TO_CHAR(created_at, 'YYYY-MM-DD')").
		Order("date ASC").
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{
		"period": period,
		"data":   results,
	})
}

func (h *DashboardHandler) GetPluginStats(c *gin.Context) {
	type PluginStat struct {
		ID             uuid.UUID `json:"id"`
		Name           string    `json:"name"`
		TotalSales     int64     `json:"total_sales"`
		TotalRevenue   float64   `json:"total_revenue"`
		ActiveLicenses int64     `json:"active_licenses"`
	}

	var results []PluginStat
	h.db.Model(&models.Plugin{}).
		Select(`plugins.id, plugins.name,
				COUNT(DISTINCT CASE WHEN orders.payment_status = 'paid' THEN orders.id END) as total_sales,
				COALESCE(SUM(CASE WHEN orders.payment_status = 'paid' THEN orders.amount ELSE 0 END), 0) as total_revenue,
				COUNT(DISTINCT CASE WHEN licenses.status = 'active' THEN licenses.id END) as active_licenses`).
		Joins("LEFT JOIN orders ON orders.plugin_id = plugins.id").
		Joins("LEFT JOIN licenses ON licenses.plugin_id = plugins.id").
		Group("plugins.id, plugins.name").
		Order("total_revenue DESC").
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{"plugins": results})
}

// ==================== Exchange Rate Management ====================

// GetExchangeRates 获取所有汇率
func (h *AdminHandler) GetExchangeRates(c *gin.Context) {
	var rates []models.ExchangeRate
	if err := h.db.Order("from_currency ASC, to_currency ASC").Find(&rates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exchange rates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rates": rates,
		"count": len(rates),
	})
}

// UpdateExchangeRates 手动触发更新汇率
func (h *AdminHandler) UpdateExchangeRates(c *gin.Context) {
	exchangeRateSvc := services.NewExchangeRateService(h.db, h.config)

	if err := exchangeRateSvc.UpdateExchangeRates(); err != nil {
		log.Printf("❌ Failed to update exchange rates: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update exchange rates: %v", err)})
		return
	}

	// 返回更新后的汇率
	var rates []models.ExchangeRate
	if err := h.db.Order("from_currency ASC, to_currency ASC").Find(&rates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated rates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exchange rates updated successfully",
		"rates":   rates,
		"count":   len(rates),
	})
}
