package repositories

import (
	"context"

	"github.com/ZXstrike/shared/pkg/models" // Ensure this package provides the new User model definition
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetRoleByName(ctx context.Context, s string) (*models.Role, error)
	CreateRole(ctx context.Context, role *models.Role) error
}
type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) GetRoleByName(ctx context.Context, s string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).First(&role, "name = ?", s).Error
	if err != nil {
		return nil, err // Return nil for role if error occurs (e.g., gorm.ErrRecordNotFound)
	}
	return &role, nil
}

func (r *repository) CreateRole(ctx context.Context, role *models.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *repository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err // Return nil for user if error occurs (e.g., gorm.ErrRecordNotFound)
	}
	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err // Return nil for user if error occurs
	}
	return &user, nil
}
