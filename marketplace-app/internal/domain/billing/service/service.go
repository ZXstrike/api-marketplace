package service

import (
	"crypto/ecdsa"

	"github.com/ZXstrike/marketplace-app/internal/domain/billing/repositories"
	"github.com/ZXstrike/shared/pkg/models"
)

type Service interface {
	TopUp(userID, walletType string, amount float64) error
	GetBillingInfo(userID, walletType string) (map[string]interface{}, error)
	ProcessPayment(userID, walletType string, amount float64) error
	GetPaymentHistory(userID string) ([]models.PaymentTransaction, error)
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

func (s *service) TopUp(userID, walletType string, amount float64) error {
	return s.r.TopUpBalanceByUserIDAndType(userID, walletType, amount)
}

func (s *service) GetBillingInfo(userID, walletType string) (map[string]interface{}, error) {
	balance, err := s.r.GetBalanceByUserIDAndType(userID, walletType)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user_id":     userID,
		"wallet_type": walletType,
		"balance":     balance,
	}, nil
}

func (s *service) ProcessPayment(userID, walletType string, amount float64) error {
	return s.r.DeductBalanceByUserIDAndType(userID, walletType, amount)
}

func (s *service) GetPaymentHistory(userID string) ([]models.PaymentTransaction, error) {
	return s.r.GetPaymentHistory(userID)
}
