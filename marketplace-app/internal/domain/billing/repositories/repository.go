package repositories

import (
	"errors"

	"github.com/ZXstrike/shared/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetBalanceByUserIDAndType(userID, walletType string) (float64, error)
	TopUpBalanceByUserIDAndType(userID, walletType string, amount float64) error
	DeductBalanceByUserIDAndType(userID, walletType string, amount float64) error
	GetPaymentHistory(userID string) ([]models.PaymentTransaction, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetBalanceByUserIDAndType(userID, walletType string) (float64, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ? AND wallet_type = ?", userID, walletType).First(&wallet).Error; err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func (r *repository) TopUpBalanceByUserIDAndType(userID, walletType string, amount float64) error {
	if amount <= 0 {
		return errors.New("top-up amount must be positive")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		// Find the wallet to ensure it exists before trying to update.
		if err := tx.Where("user_id = ? AND wallet_type = ?", userID, walletType).First(&wallet).Error; err != nil {
			return err // Returns gorm.ErrRecordNotFound if wallet doesn't exist
		}

		// Atomically update the balance.
		result := tx.Model(&wallet).Update("balance", gorm.Expr("balance + ?", amount))
		if result.Error != nil {
			return result.Error
		}

		// Log this transaction in a payment history table
		transaction := models.PaymentTransaction{
			UserID:      userID,
			Amount:      amount,
			Description: "Top-up to " + walletType + " wallet",
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) DeductBalanceByUserIDAndType(userID, walletType string, amount float64) error {
	if amount <= 0 {
		return errors.New("deduction amount must be positive")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Atomically find and update the wallet if the balance is sufficient.
		result := tx.Model(&models.Wallet{}).
			Where("user_id = ? AND wallet_type = ? AND balance >= ?", userID, walletType, amount).
			Update("balance", gorm.Expr("balance - ?", amount))

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("insufficient balance or wallet not found")
		}

		// Log the deduction
		transaction := models.PaymentTransaction{
			UserID:      userID,
			Amount:      -amount, // Store deduction as a negative amount
			Description: "Deduction from " + walletType + " wallet",
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) GetPaymentHistory(userID string) ([]models.PaymentTransaction, error) {
	var payments []models.PaymentTransaction
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}
