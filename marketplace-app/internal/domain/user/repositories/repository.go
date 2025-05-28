package repositories

import (
	"context"

	"github.com/ZXstrike/shared/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, id string, description string) error
	ChangePassword(ctx context.Context, id string, newPass string) error
	UpdateProfilePicture(ctx context.Context, id string, newProfilePicture string) error
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	return &user, err
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	return &user, err
}

func (r *repository) Update(ctx context.Context, id string, description string) error {
	var user models.User
	err := r.db.WithContext(ctx).Model(&user).Where("id = ?", id).Update("description", description).Error
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(user).Error

}

func (r *repository) ChangePassword(ctx context.Context, id string, newPass string) error {
	var user models.User
	err := r.db.WithContext(ctx).Model(&user).Where("id = ?", id).Update("password", newPass).Error
	return err
}

func (r *repository) UpdateProfilePicture(ctx context.Context, id string, newProfilePicture string) error {
	var user models.User
	err := r.db.WithContext(ctx).Model(&user).Where("id = ?", id).Update("profile_picture", newProfilePicture).Error
	return err
}
