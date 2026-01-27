package config

import (
	"os"
	"strconv"
)

type Config struct {
	// App
	AppEnv      string
	AppPort     string
	AppURL      string
	FrontendURL string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// GitHub OAuth
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string

	// GitHub App
	GitHubAppID             string
	GitHubAppPrivateKeyPath string
	GitHubAppInstallationID string
	GitHubOrgName           string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// Stripe
	StripeSecretKey      string
	StripePublishableKey string
	StripeWebhookSecret  string

	// PayPal
	PayPalClientID     string
	PayPalClientSecret string
	PayPalMode         string

	// Alipay (易支付)
	AlipayPID        string // 易支付商户ID
	AlipayAPIURL     string // 易支付API地址
	AlipayPrivateKey string // 商户私钥
	AlipayPublicKey  string // 平台公钥
	AlipayNotifyURL  string // 异步通知地址
	AlipayAppID      string // 保留兼容（已废弃）

	// Email
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	SMTPFromName string

	// Cron
	CronMaintenanceCheck string

	// Admin
	AdminEmail    string
	AdminGitHubID string

	// Defaults
	DefaultMaintenanceMonths int
}

func Load() *Config {
	jwtExpiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "720"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	defaultMaintenanceMonths, _ := strconv.Atoi(getEnv("DEFAULT_MAINTENANCE_MONTHS", "12"))

	return &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		AppPort:     getEnv("APP_PORT", "8080"),
		AppURL:      getEnv("APP_URL", "http://localhost:8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "plugin_store"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubRedirectURL:  getEnv("GITHUB_REDIRECT_URL", ""),

		GitHubAppID:             getEnv("GITHUB_APP_ID", ""),
		GitHubAppPrivateKeyPath: getEnv("GITHUB_APP_PRIVATE_KEY_PATH", ""),
		GitHubAppInstallationID: getEnv("GITHUB_APP_INSTALLATION_ID", ""),
		GitHubOrgName:           getEnv("GITHUB_ORG_NAME", ""),

		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpiryHours: jwtExpiryHours,

		StripeSecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
		StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		StripeWebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),

		PayPalClientID:     getEnv("PAYPAL_CLIENT_ID", ""),
		PayPalClientSecret: getEnv("PAYPAL_CLIENT_SECRET", ""),
		PayPalMode:         getEnv("PAYPAL_MODE", "sandbox"),

		AlipayPID:        getEnv("EPAY_PID", ""),
		AlipayAPIURL:     getEnv("EPAY_API_URL", "https://p.ma3fu.com/api/pay/create"),
		AlipayPrivateKey: getEnv("EPAY_PRIVATE_KEY", ""),
		AlipayPublicKey:  getEnv("EPAY_PUBLIC_KEY", ""),
		AlipayNotifyURL:  getEnv("EPAY_NOTIFY_URL", ""),
		AlipayAppID:      getEnv("ALIPAY_APP_ID", ""),

		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     smtpPort,
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),
		SMTPFromName: getEnv("SMTP_FROM_NAME", "Plugin Store"),

		CronMaintenanceCheck: getEnv("CRON_MAINTENANCE_CHECK", "0 2 * * *"),

		AdminEmail:    getEnv("ADMIN_EMAIL", ""),
		AdminGitHubID: getEnv("ADMIN_GITHUB_ID", ""),

		DefaultMaintenanceMonths: defaultMaintenanceMonths,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
