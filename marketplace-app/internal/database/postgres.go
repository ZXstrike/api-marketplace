package database

import (
	"fmt"
	"log"

	"github.com/ZXstrike/internal/config"
	"github.com/ZXstrike/internal/models"
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

	migration(db)

	PostgresDB = db

	return db, nil
}

func migration(db *gorm.DB) {
	// Automatically migrate your schema, to keep your database up to date.
	err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
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
		// &models.RateLimitCounter{},  // optional: drop if using pure Redis
	)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("✅ database migrated successfully")
}
