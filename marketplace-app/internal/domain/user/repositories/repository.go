package repositories

import (
	"context"

	"github.com/ZXstrike/shared/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, id string, description string) error
	ChangePassword(ctx context.Context, id string, oldPass string, newPass string) error
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
	err := r.db.WithContext(ctx).Preload("Roles").First(&user, "id = ?", id).Error
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

func (r *repository) ChangePassword(ctx context.Context, id string, oldPass string, newPass string) error {
	var user models.User

	// Fetch the user by ID
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return err
	}

	// Check if the old password matches
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPass)); err != nil {
		return err
	}

	// Generate the new password hash
	if len(newPass) < 8 {
		return gorm.ErrInvalidData // or a custom error indicating password too short
	}

	if newPass == oldPass {
		return gorm.ErrInvalidData // or a custom error indicating new password must be different
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Model(&user).Where("id = ?", id).Update("password_hash", hashPass).Error
	return err
}

func (r *repository) UpdateProfilePicture(ctx context.Context, id string, newProfilePicture string) error {
	var user models.User
	err := r.db.WithContext(ctx).Model(&user).Where("id = ?", id).Update("profile_picture_url", newProfilePicture).Error
	return err
}
