package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	Role      string    `gorm:"default:'user'" json:"role"` // user, admin
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	GitHubAccounts []GitHubAccount `gorm:"foreignKey:UserID" json:"github_accounts,omitempty"`
	Orders         []Order         `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	Licenses       []License       `gorm:"foreignKey:UserID" json:"licenses,omitempty"`
}

type GitHubAccount struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	GitHubUserID   *int64     `gorm:"unique" json:"github_user_id"`
	GitHubOrgID    *int64     `json:"github_org_id"`
	AccountType    string     `gorm:"not null" json:"account_type"` // user, org
	Login          string     `gorm:"not null" json:"login"`
	InstallationID *int64     `json:"installation_id"`
	AccessToken    string     `json:"-"`
	RefreshToken   string     `json:"-"`
	TokenExpiresAt *time.Time `json:"token_expires_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Licenses []License `gorm:"foreignKey:GitHubAccountID" json:"licenses,omitempty"`
}

type Plugin struct {
	ID                       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name                     string    `gorm:"not null" json:"name"`
	Slug                     string    `gorm:"unique;not null" json:"slug"`
	Description              string    `json:"description"`
	LongDescription          string    `json:"long_description"`
	GitHubRepoID             int64     `gorm:"column:github_repo_id" json:"github_repo_id"`
	GitHubRepoURL            string    `gorm:"column:github_repo_url" json:"github_repo_url"`
	GitHubRepoName           string    `gorm:"column:github_repo_name" json:"github_repo_name"`
	Price                    float64   `gorm:"type:decimal(10,2);default:0.00" json:"price"`
	Currency                 string    `gorm:"default:'USD'" json:"currency"`
	DefaultMaintenanceMonths int       `gorm:"default:12" json:"default_maintenance_months"`
	Status                   string    `gorm:"default:'draft'" json:"status"` // draft, published, archived
	Category                 string    `json:"category"`
	Tags                     []string  `gorm:"type:text[]" json:"tags"`
	IconURL                  string    `json:"icon_url"`
	DemoURL                  string    `json:"demo_url"`
	DocumentationURL         string    `json:"documentation_url"`
	Version                  string    `json:"version"`
	DownloadCount            int       `gorm:"default:0" json:"download_count"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`

	Orders    []Order    `gorm:"foreignKey:PluginID" json:"orders,omitempty"`
	Licenses  []License  `gorm:"foreignKey:PluginID" json:"licenses,omitempty"`
	Tutorials []Tutorial `gorm:"foreignKey:PluginID" json:"tutorials,omitempty"`
}

type Order struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrderNumber          string     `gorm:"unique;not null" json:"order_number"`
	UserID               uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	PluginID             uuid.UUID  `gorm:"type:uuid;not null" json:"plugin_id"`
	Amount               float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency             string     `gorm:"default:'USD'" json:"currency"`
	PaymentMethod        string     `gorm:"not null" json:"payment_method"`          // stripe, paypal, alipay
	PaymentStatus        string     `gorm:"default:'pending'" json:"payment_status"` // pending, paid, failed, refunded
	PaymentIntentID      string     `json:"payment_intent_id"`
	PaymentTransactionID string     `json:"payment_transaction_id"`
	PaidAt               *time.Time `json:"paid_at"`
	RefundedAt           *time.Time `json:"refunded_at"`
	Metadata             string     `gorm:"type:jsonb" json:"metadata"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`

	User    User     `gorm:"foreignKey:UserID" json:"user"`
	Plugin  Plugin   `gorm:"foreignKey:PluginID" json:"plugin"`
	License *License `gorm:"foreignKey:OrderID" json:"license,omitempty"`
}

type License struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	PluginID         uuid.UUID  `gorm:"type:uuid;not null" json:"plugin_id"`
	OrderID          uuid.UUID  `gorm:"type:uuid;not null" json:"order_id"`
	GitHubAccountID  uuid.UUID  `gorm:"type:uuid;column:git_hub_account_id;not null" json:"github_account_id"`
	LicenseType      string     `gorm:"default:'permanent'" json:"license_type"` // permanent, trial
	MaintenanceUntil time.Time  `gorm:"type:date;not null" json:"maintenance_until"`
	Status           string     `gorm:"default:'active'" json:"status"` // active, expired, revoked
	RevokedReason    string     `json:"revoked_reason"`
	RevokedAt        *time.Time `json:"revoked_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	User          User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Plugin        Plugin           `gorm:"foreignKey:PluginID" json:"plugin,omitempty"`
	Order         Order            `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	GitHubAccount GitHubAccount    `gorm:"foreignKey:GitHubAccountID" json:"github_account,omitempty"`
	History       []LicenseHistory `gorm:"foreignKey:LicenseID" json:"history,omitempty"`
}

type LicenseHistory struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	LicenseID   uuid.UUID  `gorm:"type:uuid;not null" json:"license_id"`
	Action      string     `gorm:"not null" json:"action"` // granted, expired, renewed, revoked, github_access_granted, github_access_revoked
	PerformedBy *uuid.UUID `gorm:"type:uuid" json:"performed_by"`
	Metadata    string     `gorm:"type:jsonb" json:"metadata"`
	OccurredAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"occurred_at"`

	License License `gorm:"foreignKey:LicenseID" json:"license,omitempty"`
}

type Tutorial struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PluginID    *uuid.UUID `gorm:"type:uuid" json:"plugin_id"`
	Title       string     `gorm:"not null" json:"title"`
	Slug        string     `gorm:"unique;not null" json:"slug"`
	Content     string     `gorm:"not null" json:"content"`
	ContentType string     `gorm:"default:'markdown'" json:"content_type"` // markdown, html
	OrderIndex  int        `gorm:"default:0" json:"order_index"`
	IsPublic    bool       `gorm:"default:false" json:"is_public"`
	Language    string     `gorm:"default:'en'" json:"language"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Plugin *Plugin `gorm:"foreignKey:PluginID" json:"plugin,omitempty"`
}

type EmailNotification struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	NotificationType string     `gorm:"not null" json:"notification_type"`
	Subject          string     `gorm:"not null" json:"subject"`
	Body             string     `gorm:"not null" json:"body"`
	SentAt           *time.Time `json:"sent_at"`
	Status           string     `gorm:"default:'pending'" json:"status"` // pending, sent, failed
	ErrorMessage     string     `json:"error_message"`
	Metadata         string     `gorm:"type:jsonb" json:"metadata"`
	CreatedAt        time.Time  `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Statistic struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	StatDate        time.Time `gorm:"type:date;unique;not null" json:"stat_date"`
	TotalUsers      int       `gorm:"default:0" json:"total_users"`
	NewUsers        int       `gorm:"default:0" json:"new_users"`
	TotalOrders     int       `gorm:"default:0" json:"total_orders"`
	NewOrders       int       `gorm:"default:0" json:"new_orders"`
	TotalRevenue    float64   `gorm:"type:decimal(12,2);default:0.00" json:"total_revenue"`
	DailyRevenue    float64   `gorm:"type:decimal(12,2);default:0.00" json:"daily_revenue"`
	ActiveLicenses  int       `gorm:"default:0" json:"active_licenses"`
	ExpiredLicenses int       `gorm:"default:0" json:"expired_licenses"`
	Metadata        string    `gorm:"type:jsonb" json:"metadata"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type SystemSetting struct {
	Key         string    `gorm:"primary_key" json:"key"`
	Value       string    `gorm:"not null" json:"value"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate hook for models to generate UUID if not provided
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (g *GitHubAccount) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

func (p *Plugin) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (l *License) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

func (lh *LicenseHistory) BeforeCreate(tx *gorm.DB) error {
	if lh.ID == uuid.Nil {
		lh.ID = uuid.New()
	}
	return nil
}

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Slug        string    `gorm:"unique;not null" json:"slug"`
	Description string    `json:"description"`
	IconURL     string    `json:"icon_url"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (t *Tutorial) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (e *EmailNotification) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

func (s *Statistic) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
