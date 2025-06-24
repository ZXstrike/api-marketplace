package repositories

import (
	"context"

	"github.com/ZXstrike/shared/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateStore(ctx context.Context, user_id string) error
	GetStoreByUserID(ctx context.Context, user_id string) (*models.User, error)
	GetStoreByUsername(ctx context.Context, username string) (*models.User, error)
	GetAllStores(ctx context.Context) ([]models.User, error)
	UpdateStore(ctx context.Context, user_id string, description string) error
	GetStoreApis(ctx context.Context, user_id string) ([]models.API, error)
	GetApiVersionsSubsCount(ctx context.Context, apiID string) (int, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateStore(ctx context.Context, user_id string) error {
	var user models.User

	if err := r.db.WithContext(ctx).First(&user, "id = ?", user_id).Error; err != nil {
		return err
	}

	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, "name = ?", "store_owner").Error; err != nil {
		if r.db.WithContext(ctx).Where("name = ?", "store_owner").First(&role).RowsAffected == 0 {
			// Create the role if it doesn't exist
			role = models.Role{
				Name:        "store_owner",
				Description: "Owner of the store",
			}
			if err := r.db.WithContext(ctx).Create(&role).Error; err != nil {
				return err
			}
		}
	}

	userRole := models.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
	}

	return r.db.WithContext(ctx).Create(&userRole).Error

}

func (r *repository) GetStoreByUserID(ctx context.Context, user_id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&user, "id = ?", user_id).Error; err != nil {
		return nil, err
	}

	for _, role := range user.Roles {
		if role.Name == "store_owner" {
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
		if role.Name == "store_owner" {
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
		Where("roles.name = ?", "store_owner").
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

func (r *repository) GetStoreApis(ctx context.Context, user_id string) ([]models.API, error) {
	var apis []models.API
	if err := r.db.Where("provider_id = ?", user_id).
		Find(&apis).Error; err != nil {
		return nil, err
	}

	return apis, nil
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
