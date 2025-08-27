package repositories

import (
	"context"
	"time"

	"github.com/ZXstrike/shared/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateStore(ctx context.Context, user_id string) error
	GetStoreByUserID(ctx context.Context, user_id string) (*models.User, error)
	GetStoreByUsername(ctx context.Context, username string) (*models.User, error)
	GetAllStores(ctx context.Context) ([]models.User, error)
	UpdateStore(ctx context.Context, user_id string, description string) error
	GetStoreApis(ctx context.Context, user_id string) ([]ApiWithRevenue, error)
	GetApiVersionsSubsCount(ctx context.Context, apiID string) (int, error)
}

type ApiWithRevenue struct {
	models.API
	TotalRevenue float64 `json:"total_revenue"`
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateStore(ctx context.Context, user_id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.First(&user, "id = ?", user_id).Error; err != nil {
			return err // User not found
		}

		// Find or create the 'provider' role
		var providerRole models.Role
		if err := tx.Where("name = ?", "provider").FirstOrCreate(&providerRole, models.Role{Name: "provider", Description: "Owner of the store"}).Error; err != nil {
			return err
		}

		// Assign the provider role to the user if they don't already have it
		var userRole models.UserRole
		if err := tx.Where("user_id = ? AND role_id = ?", user.ID, providerRole.ID).First(&userRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				userRole = models.UserRole{UserID: user.ID, RoleID: providerRole.ID}
				if err := tx.Create(&userRole).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Create a provider wallet if one doesn't already exist
		var providerWallet models.Wallet
		if err := tx.Where("user_id = ? AND wallet_type = ?", user.ID, "provider").First(&providerWallet).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				providerWallet = models.Wallet{UserID: user.ID, WalletType: "provider", Balance: 0}
				if err := tx.Create(&providerWallet).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil // Commit transaction
	})
}

func (r *repository) GetStoreByUserID(ctx context.Context, user_id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&user, "id = ?", user_id).Error; err != nil {
		return nil, err
	}

	for _, role := range user.Roles {
		if role.Name == "provider" {
			return &user, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (r *repository) GetStoreByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}

	for _, role := range user.Roles {
		if role.Name == "provider" {
			return &user, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (r *repository) GetAllStores(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).
		Distinct("users.*").
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("roles.name = ?", "provider").
		Preload("Roles").
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) UpdateStore(ctx context.Context, user_id string, description string) error {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", user_id).Error; err != nil {
		return err
	}

	user.Description = description
	return r.db.WithContext(ctx).Save(&user).Error
}

func (r *repository) GetStoreApis(ctx context.Context, user_id string) ([]ApiWithRevenue, error) {
	var apis []models.API
	if err := r.db.WithContext(ctx).Where("provider_id = ?", user_id).
		Preload("Versions").
		Find(&apis).Error; err != nil {
		return nil, err
	}

	if len(apis) == 0 {
		return []ApiWithRevenue{}, nil
	}

	apiIDs := make([]string, len(apis))
	for i, api := range apis {
		apiIDs[i] = api.ID
	}

	type apiRevenueResult struct {
		ApiID        string
		TotalRevenue float64
	}

	var revenues []apiRevenueResult
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, 0)

	if err := r.db.Model(&models.PaymentTransaction{}).
		Select("api_versions.api_id, COALESCE(SUM(payment_transactions.amount), 0) as total_revenue").
		Joins("JOIN subscriptions ON subscriptions.id = payment_transactions.subscription_id").
		Joins("JOIN api_versions ON api_versions.id = subscriptions.api_version_id").
		Where("api_versions.api_id IN ?", apiIDs).
		Where("payment_transactions.user_id = ?", user_id).
		Where("payment_transactions.transaction_time >= ? AND payment_transactions.transaction_time < ?", firstOfMonth, lastOfMonth).
		Where("payment_transactions.amount > 0"). // Only consider credits (income)
		Group("api_versions.api_id").
		Scan(&revenues).Error; err != nil {
		return nil, err
	}

	revenueMap := make(map[string]float64, len(revenues))
	for _, rev := range revenues {
		revenueMap[rev.ApiID] = rev.TotalRevenue
	}

	apisWithRevenue := make([]ApiWithRevenue, len(apis))
	for i, api := range apis {
		apisWithRevenue[i] = ApiWithRevenue{
			API:          api,
			TotalRevenue: revenueMap[api.ID],
		}
	}

	return apisWithRevenue, nil
}

func (r *repository) GetApiVersionsSubsCount(ctx context.Context, apiID string) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Subscription{}).
		Joins("JOIN api_versions ON api_versions.id = subscriptions.api_version_id").
		Where("api_versions.api_id = ?", apiID).
		Where("subscriptions.deleted_at IS NULL").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
