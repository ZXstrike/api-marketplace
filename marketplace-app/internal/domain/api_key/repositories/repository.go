package repositories

import (
	"github.com/ZXstrike/shared/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	// CORRECTED a more logical and consistent signature.
	CreateAPIKey(prefix string, apiKeyHash []byte, subscriptionID string) error
	GetSubscriptionAPIKeys(subscriptionID string) ([]models.APIKey, error)
	DeleteAPIKey(apiKeyID string) error
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateAPIKey(prefix string, apiKeyHash []byte, subscriptionID string) error {
	// The check for subscription existence here is good, but it might be redundant
	// if the service layer already does it. However, it provides good data integrity.
	var subscription models.Subscription
	if err := r.db.Where("id = ?", subscriptionID).First(&subscription).Error; err != nil {
		return err // Subscription not found
	}

	newAPIKey := models.APIKey{
		SubscriptionID: subscriptionID,
		// No need to set the Subscription struct itself, GORM handles it via the foreign key.
		KeyValueHash: apiKeyHash,
		KeyPrefix:    prefix, // FIXED: Added the missing prefix assignment.
	}

	if err := r.db.Create(&newAPIKey).Error; err != nil {
		return err
	}

	return nil
}

// GetSubscriptionAPIKeys and DeleteAPIKey are fine as they were.
func (r *repository) GetSubscriptionAPIKeys(subscriptionID string) ([]models.APIKey, error) {
	var apiKeys []models.APIKey
	if err := r.db.Where("subscription_id = ?", subscriptionID).Find(&apiKeys).Error; err != nil {
		return nil, err
	}
	return apiKeys, nil
}

func (r *repository) DeleteAPIKey(apiKeyId string) error {
	// Using .Delete with a model and primary key is more idiomatic in GORM.
	if err := r.db.Delete(&models.APIKey{}, "id = ?", apiKeyId).Error; err != nil {
		return err
	}
	return nil
}
