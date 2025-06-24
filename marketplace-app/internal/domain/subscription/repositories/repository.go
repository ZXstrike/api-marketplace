package repositories

import (
	"github.com/ZXstrike/shared/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	SubscribeToAPI(userID string, apiVersionID string) error
	UnsubscribeFromAPI(userID string, subscriptionID string) error
	GetSubscription(userID string, subscriptionID string) (*models.Subscription, error)
	GetSubscriptionsByUserID(userID string) ([]models.Subscription, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) SubscribeToAPI(userID string, apiVersionID string) error {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	var apiVersion models.APIVersion

	subscription := models.Subscription{
		ConsumerUserID: userID,
		Consumer:       user,
		APIVersionID:   apiVersionID,
		APIVersion:     apiVersion,
	}

	if err := r.db.Create(&subscription).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) UnsubscribeFromAPI(userID string, subscriptionID string) error {
	var subscription models.Subscription
	if err := r.db.Where("consumer_user_id = ? AND id = ?", userID, subscriptionID).First(&subscription).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&subscription).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) GetSubscription(userID string, subscriptionID string) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := r.db.Where("consumer_user_id = ? AND id = ?", userID, subscriptionID).First(&subscription).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *repository) GetSubscriptionsByUserID(userID string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := r.db.Where("consumer_user_id = ?", userID).
		Preload("APIVersion").
		Preload("APIVersion.API").
		Preload("APIKeys").
		Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}
