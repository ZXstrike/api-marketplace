package models

import (
	"time"

	"gorm.io/gorm"
)

// Base model with common timestamps & soft‑delete
type Base struct {
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

// User represents any account (consumer, provider, admin)
type User struct {
	Base
	ID                string               `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username          string               `json:"username" gorm:"uniqueIndex;not null"`
	Email             string               `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash      string               `json:"-" gorm:"column:password_hash;not null"`
	Description       string               `json:"description" gorm:"type:text"`
	ProfilePictureURL string               `json:"profile_picture_url" gorm:"type:varchar(255)"`
	AccountBalance    float64              `json:"-" gorm:"type:decimal(12,2);not null;default:0"`
	Roles             []Role               `json:"roles,omitempty" gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`
	APIs              []API                `json:"apis,omitempty" gorm:"foreignKey:ProviderID;constraint:OnDelete:CASCADE"`
	Subscriptions     []Subscription       `json:"subscriptions,omitempty" gorm:"foreignKey:ConsumerUserID;constraint:OnDelete:CASCADE"`
	Payments          []PaymentTransaction `json:"payments,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Payouts           []ProviderPayout     `json:"payouts,omitempty" gorm:"foreignKey:ProviderID;constraint:OnDelete:CASCADE"`
}

// Role represents a permission set
type Role struct {
	Base
	ID          string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description" gorm:"type:text"`
	Users       []User `json:"users,omitempty" gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`
}

// Join table: user ↔ role
type UserRole struct {
	UserID     string    `json:"user_id" gorm:"primaryKey;type:uuid"`
	RoleID     string    `json:"role_id" gorm:"primaryKey;type:uuid"`
	AssignedAt time.Time `json:"assigned_at" gorm:"autoCreateTime"`
}

// Category for grouping APIs
type Category struct {
	Base
	ID   string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name string `json:"name" gorm:"uniqueIndex;not null"`
	Slug string `json:"slug" gorm:"uniqueIndex;not null"`
	APIs []API  `json:"apis,omitempty" gorm:"many2many:api_categories;constraint:OnDelete:CASCADE"`
}

// API product published by a provider
type API struct {
	Base
	ID          string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProviderID  string       `json:"provider_id" gorm:"type:uuid;not null;index"`
	Provider    User         `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
	Name        string       `json:"name" gorm:"not null"`
	Description string       `json:"description" gorm:"type:text"`
	Categories  []Category   `json:"categories,omitempty" gorm:"many2many:api_categories;constraint:OnDelete:CASCADE"`
	Versions    []APIVersion `json:"versions,omitempty" gorm:"foreignKey:APIID;constraint:OnDelete:CASCADE"`
}

// Join table: api ↔ category
type APICategory struct {
	APIID      string `json:"api_id" gorm:"primaryKey;type:uuid"`
	CategoryID string `json:"category_id" gorm:"primaryKey;type:uuid"`
}

// A released version of an API
type APIVersion struct {
	Base
	ID            string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	APIID         string         `json:"api_id" gorm:"type:uuid;not null;index"`
	API           API            `json:"api,omitempty" gorm:"foreignKey:APIID"`
	VersionString string         `json:"version_string" gorm:"not null"`
	PricePerCall  float64        `json:"price_per_call" gorm:"type:decimal(12,6);not null;default:0"`
	Endpoints     []Endpoint     `json:"endpoints,omitempty" gorm:"foreignKey:APIVersionID;constraint:OnDelete:CASCADE"`
	Subscriptions []Subscription `json:"subscriptions,omitempty" gorm:"foreignKey:APIVersionID;constraint:OnDelete:CASCADE"`
}

// An individual route on an API version, now with docs
type Endpoint struct {
	Base
	ID            string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	APIVersionID  string     `json:"api_version_id" gorm:"type:uuid;not null;index"`
	APIVersion    APIVersion `json:"api_version,omitempty" gorm:"foreignKey:APIVersionID"`
	HTTPMethod    string     `json:"http_method" gorm:"size:10;not null"`
	Path          string     `json:"path" gorm:"not null"`
	Documentation string     `json:"documentation" gorm:"type:text"` // ← per‑endpoint docs
}

// A consumer’s subscription to one API version
type Subscription struct {
	Base
	ID               string               `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ConsumerUserID   string               `json:"consumer_user_id" gorm:"type:uuid;not null;index"`
	Consumer         User                 `json:"consumer,omitempty" gorm:"foreignKey:ConsumerUserID"`
	APIVersionID     string               `json:"api_version_id" gorm:"type:uuid;not null;index"`
	APIVersion       APIVersion           `json:"api_version,omitempty" gorm:"foreignKey:APIVersionID"`
	MaxMonthlyBudget float64              `json:"max_monthly_budget" gorm:"type:decimal(12,2);not null;default:0"`
	APIKeys          []APIKey             `json:"api_keys,omitempty" gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE"`
	UsageLogs        []UsageLog           `json:"usage_logs,omitempty" gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE"`
	Payments         []PaymentTransaction `json:"payments,omitempty" gorm:"foreignKey:SubscriptionID;constraint:OnDelete:SET NULL"`
	Payouts          []ProviderPayout     `json:"payouts,omitempty" gorm:"foreignKey:SubscriptionID;constraint:OnDelete:SET NULL"`
	Statements       []MonthlyStatement   `json:"statements,omitempty" gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE"`
}

// An API key issued under a subscription
type APIKey struct {
	Base
	ID             string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SubscriptionID string       `json:"subscription_id" gorm:"type:uuid;not null;index"`
	Subscription   Subscription `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
	KeyValueHash   string       `json:"key_value_hash" gorm:"not null"`
	UsageLogs      []UsageLog   `json:"usage_logs,omitempty" gorm:"foreignKey:APIKeyID;constraint:OnDelete:CASCADE"`
}

// A single API call record
type UsageLog struct {
	ID               int64        `json:"id" gorm:"primaryKey;autoIncrement"`
	SubscriptionID   string       `json:"subscription_id" gorm:"type:uuid;not null;index"`
	Subscription     Subscription `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
	APIKeyID         string       `json:"api_key_id" gorm:"type:uuid;not null;index"`
	APIKey           APIKey       `json:"api_key,omitempty" gorm:"foreignKey:APIKeyID"`
	RequestTimestamp time.Time    `json:"request_timestamp" gorm:"autoCreateTime"`
}

// Wallet transactions (top‑up, charges, payouts)
type PaymentTransaction struct {
	ID              string        `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID          string        `json:"user_id" gorm:"type:uuid;not null;index"`
	User            User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	SubscriptionID  *string       `json:"subscription_id,omitempty" gorm:"type:uuid;index"`
	Subscription    *Subscription `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
	Amount          float64       `json:"amount" gorm:"type:decimal(12,2);not null"`
	TransactionTime time.Time     `json:"transaction_time" gorm:"autoCreateTime"`
	Description     string        `json:"description" gorm:"type:text"`
}

// Pre‑computed monthly summary per subscription
type MonthlyStatement struct {
	SubscriptionID string    `json:"subscription_id" gorm:"primaryKey;type:uuid;index"`
	Year           int       `json:"year" gorm:"primaryKey"`
	Month          int       `json:"month" gorm:"primaryKey"`
	TotalCalls     int64     `json:"total_calls" gorm:"not null"`
	TotalSpent     float64   `json:"total_spent" gorm:"type:decimal(12,2);not null"`
	GeneratedAt    time.Time `json:"generated_at" gorm:"autoCreateTime"`
}

// Optional fallback: rate‑limit counters
type RateLimitCounter struct {
	APIKeyID    string    `json:"api_key_id" gorm:"primaryKey;type:uuid"`
	WindowStart time.Time `json:"window_start" gorm:"autoCreateTime"`
	CallCount   int       `json:"call_count" gorm:"not null"`
}

// Provider payout ledger
type ProviderPayout struct {
	ID             string        `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProviderID     string        `json:"provider_id" gorm:"type:uuid;not null;index"`
	Provider       User          `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
	SubscriptionID *string       `json:"subscription_id,omitempty" gorm:"type:uuid;index"`
	Subscription   *Subscription `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
	PeriodStart    time.Time     `json:"period_start" gorm:"type:date;not null"`
	PeriodEnd      time.Time     `json:"period_end" gorm:"type:date;not null"`
	GrossAmount    float64       `json:"gross_amount" gorm:"type:decimal(12,2);not null"`
	PlatformFee    float64       `json:"platform_fee" gorm:"type:decimal(12,2);not null"`
	NetAmount      float64       `json:"net_amount" gorm:"type:decimal(12,2);not null"`
	CreatedAt      time.Time     `json:"created_at" gorm:"autoCreateTime"`
}
