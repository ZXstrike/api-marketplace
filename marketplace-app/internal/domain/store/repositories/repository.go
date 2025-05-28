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
	if err := r.db.WithContext(ctx).Preload("Roles").Where("roles.name = ?", "store_owner").Find(&users).Error; err != nil {
		return nil, err
	}

	var stores []models.User
	for _, user := range users {
		for _, role := range user.Roles {
			if role.Name == "store_owner" {
				stores = append(stores, user)
				break
			}
		}
	}

	return stores, nil
}

func (r *repository) UpdateStore(ctx context.Context, user_id string, description string) error {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", user_id).Error; err != nil {
		return err
	}

	user.Description = description
	return r.db.WithContext(ctx).Save(&user).Error
}
