package service

import (
	"crypto/ecdsa"

	"github.com/ZXstrike/marketplace-app/internal/domain/subscription/repositories"
	"github.com/ZXstrike/shared/pkg/models"
)

type Service interface {
	SubscribeToAPI(userID string, apiVersionID string) error
	UnsubscribeFromAPI(userID string, apiVersionID string) error
	GetSubscription(userID string, apiVersionID string) (*models.Subscription, error)
	GetSubscriptionsByUserID(userID string) ([]models.Subscription, error)
}

type service struct {
	repo       repositories.Repository
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func New(repo repositories.Repository, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) Service {
	return &service{repo, privateKey, publicKey}
}

func (s *service) SubscribeToAPI(userID string, apiVersionID string) error {
	// Implementation for subscribing to an API
	return s.repo.SubscribeToAPI(userID, apiVersionID)
}

func (s *service) UnsubscribeFromAPI(userID string, subscriptionID string) error {
	// Implementation for unsubscribing from an API
	return s.repo.UnsubscribeFromAPI(userID, subscriptionID)
}

func (s *service) GetSubscription(userID string, subscriptionID string) (*models.Subscription, error) {
	// Implementation for getting a subscription
	return s.repo.GetSubscription(userID, subscriptionID)
}

func (s *service) GetSubscriptionsByUserID(userID string) ([]models.Subscription, error) {
	// Implementation for getting subscriptions by user ID
	return s.repo.GetSubscriptionsByUserID(userID)
}
