package repositories

import (
	"strings"
	"time"

	"github.com/ZXstrike/shared/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetUserByID(id string) (*models.User, error)
	GetApiTotalNumbers() (int64, error)
	GetUserTotalNumbers() (int64, error)
	GetTransaction24HNumbers() (int64, error)
	GetRecentAllTopUps() ([]models.PaymentTransaction, error)
	CreateAPI(api models.API, pricePercall float64) (string, error)
	GetAPIByID(id string) (*models.API, error)
	GetAllAPI(page int, length int) ([]models.API, error)
	GetAllAPIByUserID(userID string) ([]models.API, error)
	UpdateAPI(api models.API) error
	DeleteAPI(id string) error
	GetCategoryBySlug(slug string) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	GetAPIVersionByID(id string) (*models.APIVersion, error)
	CreateAPIEndpoint(apiEndpoint models.Endpoint) error
	GetAPIEndpointByID(id string) (*models.Endpoint, error)
	UpdateAPIEndpoint(apiEndpoint models.Endpoint) error
	DeleteAPIEndpoint(id string) error
	GetAllEndpointsByAPIVersionID(apiVersionID string) ([]models.Endpoint, error)
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

	if err := r.db.Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetApiTotalNumbers() (int64, error) {
	var count int64
	if err := r.db.Model(&models.API{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) GetUserTotalNumbers() (int64, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) GetTransaction24HNumbers() (int64, error) {
	var count int64
	cutoff := time.Now().Add(-24 * time.Hour)

	tx := r.db.Model(&models.PaymentTransaction{}).
		Where("created_at >= ?", cutoff).
		Count(&count)

	if tx.Error != nil {
		// Handle missing column gracefully
		if strings.Contains(tx.Error.Error(), `column "created_at" does not exist`) {
			// Attempt to add the column (requires model definition to have CreatedAt)
			if err := r.db.Migrator().AddColumn(&models.PaymentTransaction{}, "CreatedAt"); err == nil {
				// Retry with the time filter after adding column
				if retryErr := r.db.Model(&models.PaymentTransaction{}).
					Where("created_at >= ?", cutoff).
					Count(&count).Error; retryErr == nil {
					return count, nil
				}
			}
		}
		return 0, tx.Error
	}

	return count, nil
}

func (r *repository) GetRecentAllTopUps() ([]models.PaymentTransaction, error) {
	var transactions []models.PaymentTransaction

	tx := r.db.Preload("User").Preload("Subscription").Model(&models.PaymentTransaction{}).Limit(100).Order("created_at desc").Find(&transactions)

	if tx.Error != nil {
		// Handle missing column gracefully
		if strings.Contains(tx.Error.Error(), `column "created_at" does not exist`) {
			// Attempt to add the column (requires model definition to have CreatedAt)
			if err := r.db.Migrator().AddColumn(&models.PaymentTransaction{}, "CreatedAt"); err == nil {
				// Retry with the time filter after adding column
				if retryErr := r.db.Preload("User").Preload("Subscription").Model(&models.PaymentTransaction{}).
					Limit(50).Error; retryErr == nil {

					return transactions, nil
				}
			}

		}
		return nil, tx.Error
	}

	return transactions, nil
}

func (r *repository) CreateAPI(api models.API, pricePercall float64) (string, error) {

	if err := r.db.Create(&api).Error; err != nil {
		return "", err
	}

	apiVersion := models.APIVersion{
		APIID:         api.ID,
		API:           api,
		VersionString: "v1.0.0",
		PricePerCall:  pricePercall,
	}

	if err := r.db.Create(&apiVersion).Error; err != nil {
		return "", err
	}

	if err := r.db.Model(&api).Association("Versions").Append(&apiVersion); err != nil {
		return "", err
	}

	return apiVersion.ID, nil
}

func (r *repository) GetAPIByID(id string) (*models.API, error) {
	var api models.API
	if err := r.db.First(&api, "id = ?", id).Error; err != nil {
		return nil, err
	}
	// Preload the Provider and Categories associations.
	if err := r.db.Model(&api).Preload("Provider").Preload("Categories").First(&api).Error; err != nil {
		return nil, err
	}
	// Preload the Versions association.
	if err := r.db.Model(&api).Preload("Versions").First(&api).Error; err != nil {
		return nil, err
	}

	return &api, nil
}

func (r *repository) GetAllAPI(page int, length int) ([]models.API, error) {
	var apis []models.API

	// Chain all Preload calls, then apply pagination, and finally execute Find.
	err := r.db.
		Preload("Provider").
		Preload("Categories").
		Preload("Versions.Endpoints").
		Offset((page - 1) * length).
		Limit(length).
		Find(&apis).Error

	if err != nil {
		// This will catch actual database errors.
		// GORM's Find on a slice typically doesn't return gorm.ErrRecordNotFound for an empty result set;
		// it returns an empty slice and a nil error.
		return nil, err
	}

	// This check was in your original code. It makes "no records found" an error condition.
	// If this is the desired behavior, it should be kept.
	if len(apis) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return apis, nil
}

func (r *repository) GetAllAPIByUserID(userID string) ([]models.API, error) {
	var apis []models.API
	if err := r.db.Where("provider_id = ?", userID).Find(&apis).Error; err != nil {
		return nil, err
	}
	return apis, nil
}

func (r *repository) UpdateAPI(api models.API) error {
	return r.db.Save(&api).Error
}

func (r *repository) DeleteAPI(id string) error {
	return r.db.Delete(&models.API{}, "id = ?", id).Error
}

func (r *repository) GetCategoryBySlug(slug string) (*models.Category, error) {
	var category models.Category
	if err := r.db.First(&category, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *repository) GetAPIVersionByID(id string) (*models.APIVersion, error) {
	var apiVersion models.APIVersion
	if err := r.db.First(&apiVersion, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &apiVersion, nil
}

func (r *repository) CreateAPIEndpoint(apiEndpoint models.Endpoint) error {
	return r.db.Create(&apiEndpoint).Error
}

func (r *repository) GetAPIEndpointByID(id string) (*models.Endpoint, error) {
	var apiEndpoint models.Endpoint
	if err := r.db.First(&apiEndpoint, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &apiEndpoint, nil
}

func (r *repository) UpdateAPIEndpoint(apiEndpoint models.Endpoint) error {
	return r.db.Save(&apiEndpoint).Error
}

func (r *repository) DeleteAPIEndpoint(id string) error {
	return r.db.Delete(&models.Endpoint{}, "id = ?", id).Error
}

func (r *repository) GetAllEndpointsByAPIVersionID(apiVersionID string) ([]models.Endpoint, error) {
	var endpoints []models.Endpoint
	if err := r.db.Where("api_version_id = ?", apiVersionID).Find(&endpoints).Error; err != nil {
		return nil, err
	}
	return endpoints, nil
}
