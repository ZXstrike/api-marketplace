package service

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ZXstrike/marketplace-app/internal/domain/billing/repositories"
	"github.com/ZXstrike/shared/pkg/models"
)

type Service interface {
	TopUp(userID, walletType string, amount float64) error
	GetBillingInfo(userID, walletType string) (map[string]interface{}, error)
	ProcessPayment(userID, walletType string, amount float64) error
	GetPaymentHistory(userID string) ([]models.PaymentTransaction, error)
	ProcessAndCreatePayout(ctx context.Context, payoutData *models.ProviderPayout) (*models.ProviderPayout, error)
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

	userBalance := 0.0
	if walletType != "admin" {
		balance, err := s.r.GetBalanceByUserIDAndType(userID, walletType)
		if err != nil {
			return nil, err
		}
		userBalance = balance
	} else {
		balance, err := s.r.GetPlatformWalletBalance()
		if err != nil {
			return nil, err
		}
		userBalance = balance
	}

	return map[string]interface{}{
		"user_id":     userID,
		"wallet_type": walletType,
		"balance":     userBalance,
	}, nil
}

func (s *service) ProcessPayment(userID, walletType string, amount float64) error {
	return s.r.DeductBalanceByUserIDAndType(userID, walletType, amount)
}

func (s *service) GetPaymentHistory(userID string) ([]models.PaymentTransaction, error) {
	return s.r.GetPaymentHistory(userID)
}

// ProcessAndCreatePayout calculates and stores a provider payout and platform revenue.
func (s *service) ProcessAndCreatePayout(ctx context.Context, payoutData *models.ProviderPayout) (*models.ProviderPayout, error) {
	// Business logic remains in the service.
	// e.g., Calculations for GrossAmount, PlatformFee, NetAmount should happen here.

	// Prepare the platform revenue record.
	revenue := &models.PlatformRevenue{
		// SourcePayoutID is set by the repository transaction.
		Amount:      payoutData.PlatformFee,
		Description: fmt.Sprintf("Platform fee from payout for provider %s", payoutData.ProviderID),
	}

	// Delegate the entire database transaction to the repository.
	if err := s.r.ProcessPayoutTransaction(ctx, payoutData, revenue); err != nil {
		return nil, fmt.Errorf("failed to process payout transaction: %w", err)
	}

	return payoutData, nil
}
