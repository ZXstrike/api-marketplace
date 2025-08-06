package database

import (
	"fmt"
	"log"
	"time"

	"github.com/ZXstrike/api-gateway/internal/config"
	"github.com/ZXstrike/shared/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgresDB *gorm.DB

func PostgresConnect(postgresConf *config.PostgresConfig) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		postgresConf.Host,
		postgresConf.User,
		postgresConf.Password,
		postgresConf.Database,
		postgresConf.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// --- Add Connection Pool Settings Here ---
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(25)

	// Set the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// Set the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	// --- End of Connection Pool Settings ---

	migration(db)

	PostgresDB = db

	return db, nil
}

func migration(db *gorm.DB) {
	// Automatically migrate your schema, to keep your database up to date.
	err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Wallet{},
		&models.UserRole{},
		&models.Category{},
		&models.API{},
		&models.APICategory{},
		&models.APIVersion{},
		&models.Endpoint{}, // includes Documentation field
		&models.Subscription{},
		&models.APIKey{},
		&models.UsageLog{},
		&models.PaymentTransaction{},
		&models.MonthlyStatement{},
		&models.ProviderPayout{},
		&models.PlatformFee{}, // <-- Add this line
		&models.PlatformRevenue{},
		&models.PlatformWallet{},
		&models.SystemSetting{}, // <-- Add this line
	)

	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// Create default roles if they don't exist
	models.GenerateCategories(db)

	log.Println("✅ database migrated successfully")
}
