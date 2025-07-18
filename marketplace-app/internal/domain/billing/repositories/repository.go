package repositories

import (
	"github.com/ZXstrike/shared/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetBalanceByUserID(userID string) (float64, error)
	UpdateBalanceByUserID(userID string, amount float64) error
	DeductBalanceByUserID(userID string, amount float64) error
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetBalanceByUserID(userID string) (float64, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return 0, err
	}

	return user.AccountBalance, nil
}

func (r *repository) UpdateBalanceByUserID(userID string, amount float64) error {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}
	
	// Ensure the balance does not go negative
	if user.AccountBalance < 0 {
		return gorm.ErrRecordNotFound // or a custom error indicating insufficient balance
	}

	// Update the user's balance in the database
	user.AccountBalance += amount
	if user.AccountBalance < 0 {
		return gorm.ErrRecordNotFound // or a custom error indicating insufficient balance
	}

	// Save the updated user back to the database
	if err := r.db.Save(&user).Error; err != nil {
		return err
	}

	// Optionally, you can log this transaction in a payment history table
	transaction := models.PaymentTransaction{
		UserID:      userID,
		Amount:      amount,
		Description: "Top-up balance",
	}

	if err := r.db.Create(&transaction).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) DeductBalanceByUserID(userID string, amount float64) error {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	if user.AccountBalance < amount {
		return gorm.ErrRecordNotFound // or a custom error indicating insufficient balance
	}

	user.AccountBalance -= amount

	if err := r.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) GetPaymentHistory(userID string) ([]models.PaymentTransaction, error) {
	var payments []models.PaymentTransaction
	if err := r.db.Where("user_id = ?", userID).Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}
