package repositories

import (
	"github.com/ZXstrike/shared/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetUserByID(id string) (*models.User, error)
	Create(api models.API) error
	GetByID(id string) (*models.API, error)
	GetAll() ([]models.API, error)
	GetAllByUserID(userID string) ([]models.API, error)
	Update(api models.API) error
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Create(api models.API) error {
	return r.db.Create(&api).Error
}

func (r *repository) GetByID(id string) (*models.API, error) {
	var api models.API
	if err := r.db.First(&api, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &api, nil
}

func (r *repository) GetAll() ([]models.API, error) {
	var apis []models.API
	if err := r.db.Find(&apis).Error; err != nil {
		return nil, err
	}
	return apis, nil
}

func (r *repository) GetAllByUserID(userID string) ([]models.API, error) {
	var apis []models.API
	if err := r.db.Where("provider_id = ?", userID).Find(&apis).Error; err != nil {
		return nil, err
	}
	return apis, nil
}

func (r *repository) Update(api models.API) error {
	return r.db.Save(&api).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&models.API{}, "id = ?", id).Error
}
