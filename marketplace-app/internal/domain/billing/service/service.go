package service

import (
	"crypto/ecdsa"

	"github.com/ZXstrike/marketplace-app/internal/domain/billing/repositories"
)

type Service interface {
	UpdateBalanceByUserID(userID string, amount float64) error
	GetBillingInfo(userID string) (map[string]interface{}, error)
	ProcessPayment(userID string, amount float64) error
	GetPaymentHistory(userID string) ([]map[string]interface{}, error)
}

type service struct {
	r          repositories.Repository
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func New(repo repositories.Repository, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) Service {
	return &service{
		r:          repo,
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (s *service) UpdateBalanceByUserID(userID string, amount float64) error {
	// This method will handle the logic to update the balance of a user.
	// Implementation will depend on the specific requirements and data structure.
	return s.r.UpdateBalanceByUserID(userID, amount)
}

func (s *service) GetBillingInfo(userID string) (map[string]interface{}, error) {
	// This method will handle the logic to get billing information for a user.
	// Implementation will depend on the specific requirements and data structure.
	balance, err := s.r.GetBalanceByUserID(userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user_id": userID,
		"balance": balance,
	}, nil
}

func (s *service) ProcessPayment(userID string, amount float64) error {
	// This method will handle the logic to process a payment for a user.
	// Implementation will depend on the specific requirements and payment gateway used.
	if err := s.r.DeductBalanceByUserID(userID, amount); err != nil {
		return err
	}
	return nil
}

func (s *service) GetPaymentHistory(userID string) ([]map[string]interface{}, error) {
	
	return []map[string]interface{}{}, nil
}
