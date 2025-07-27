package models

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GenerateCategories(db *gorm.DB) error {
	categories := []string{
		"Web Development",
		"Mobile Development",
		"Data Science",
		"Machine Learning",
		"Blockchain",
		"Cloud Computing",
		"Cybersecurity",
		"DevOps",
		"Artificial Intelligence",
		"Internet of Things (IoT)",
		"Game Development",
		"Augmented Reality (AR) and Virtual Reality (VR)",
		"API Development",
		"Software Testing",
		"Database Management",
		"UI/UX Design",
		"Content Management Systems (CMS)",
		"Search Engine Optimization (SEO)",
		"Digital Marketing",
		"Business Intelligence",
		"Project Management",
		"Agile Methodologies",
		"Version Control Systems",
		"Microservices Architecture",
		"Serverless Computing",
		"Edge Computing",
	}

	var existingCategories []Category
	if err := db.Find(&existingCategories).Error; err != nil {
		return err
	}

	existingCategoryNames := make(map[string]bool)
	for _, category := range existingCategories {
		existingCategoryNames[strings.ToLower(category.Name)] = true
	}

	for _, name := range categories {
		slug := GenerateSlug(name)
		if !existingCategoryNames[strings.ToLower(name)] {
			newCategory := Category{
				Name: name,
				Slug: slug,
			}
			if err := db.Create(&newCategory).Error; err != nil {
				return err
			}
		}
	}

	log.Println("✅ categories generated successfully")
	return nil
}

func GenerateAdminUser(db *gorm.DB) error {
	// Check if admin user already exists
	var existingUser User
	if err := db.Where("username = ?", "zxsttm").First(&existingUser).Error; err == nil {
		log.Println("✅ admin user 'zxsttm' already exists")
		return nil
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	// Find or create the 'admin' role
	var adminRole Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("creating 'admin' role")
			adminRole = Role{Name: "admin", Description: "Administrator with full permissions"}
			if err := db.Create(&adminRole).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Hash the password - use a more secure password in production
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create the admin user
	adminUser := User{
		Username:     "zxsttm",
		Email:        "admin@zxsttm.com",
		PasswordHash: string(hashedPassword),
		Roles:        []Role{adminRole},
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return err
	}

	log.Println("✅ admin user 'zxsttm' generated successfully")
	return nil
}

// MigrateUserBalancesToWallets moves balances from the old user schema to the new Wallet model.
// This should be run once after the schema migration.
func MigrateUserBalancesToWallets(db *gorm.DB) error {
	// Temporary struct to read the old user data, including the soon-to-be-deprecated column.
	type OldUser struct {
		ID             string  `gorm:"primaryKey"`
		AccountBalance float64 `gorm:"column:account_balance"`
		Roles          []Role  `gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:user_id"`
	}

	log.Println("🚀 starting user balance migration to wallets...")

	// Check if the migration column still exists. If not, skip.
	if !db.Migrator().HasColumn(&User{}, "account_balance") {
		log.Println("✅ account_balance column not found, migration skipped.")
		return nil
	}

	var users []OldUser
	// Use .Preload("Roles") to get the roles for each user.
	if err := db.Table("users").Preload("Roles").Find(&users).Error; err != nil {
		return fmt.Errorf("failed to query users for migration: %w", err)
	}

	for _, user := range users {
		hasConsumerRole := false
		hasProviderRole := false
		for _, role := range user.Roles {
			if role.Name == "consumer" {
				hasConsumerRole = true
			}
			if role.Name == "provider" {
				hasProviderRole = true
			}
		}

		// Create a consumer wallet if the user has the role.
		if hasConsumerRole {
			createWalletIfNotExists(db, user.ID, "consumer", user.AccountBalance)
		}

		// Create a provider wallet if the user has the role.
		if hasProviderRole {
			// Provider wallets start at 0, as the old balance was for consumption.
			createWalletIfNotExists(db, user.ID, "provider", 0)
		}
	}

	log.Println("🎉 balance migration completed successfully.")
	return nil
}

// createWalletIfNotExists is a helper function to avoid duplicate wallet creation.
func createWalletIfNotExists(db *gorm.DB, userID, walletType string, balance float64) {
	var walletCount int64
	db.Model(&Wallet{}).Where("user_id = ? AND wallet_type = ?", userID, walletType).Count(&walletCount)

	if walletCount == 0 {
		newWallet := Wallet{
			UserID:     userID,
			WalletType: walletType,
			Balance:    balance,
		}
		if err := db.Create(&newWallet).Error; err != nil {
			log.Printf("❌ failed to create %s wallet for user %s: %v", walletType, userID, err)
		} else {
			log.Printf("✅ created %s wallet for user %s with balance %f", walletType, userID, balance)
		}
	}
}

func GenerateSlug(name string) string {
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	return slug
}
