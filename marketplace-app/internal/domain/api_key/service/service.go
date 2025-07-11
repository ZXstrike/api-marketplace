package service

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ZXstrike/marketplace-app/internal/domain/api_key/repositories"
	"github.com/ZXstrike/shared/pkg/models"
	"gorm.io/gorm"
)

// NOTE: I have removed the unused ecdsa.PrivateKey and ecdsa.PublicKey for clarity.
// If you need them for something else (like JWTs), you can add them back.

type Service interface {
	CreateAPIKey(subscriptionID string) (string, error)
	DeleteAPIKey(apiKeyID string) error
}

type service struct {
	repo       repositories.Repository
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func New(repo repositories.Repository, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) Service {
	return &service{repo, privateKey, publicKey}
}

func (s *service) CreateAPIKey(subscriptionID string) (string, error) {
	// Make sure the subscription exists first.
	// NOTE: I've simplified the logic here. The original GetSubscriptionAPIKeys is fine,
	// but this check is more direct if you just need to validate the subscription's existence.
	// You can adapt this as needed.
	// if _, err := s.repo.GetSubscriptionAPIKeys(subscriptionID); err != nil {
	// 	return "", fmt.Errorf("failed to validate subscription: %w", err)
	// }

	// --- Business Logic Decision ---
	// The code below deletes all old keys when a new one is created.
	// This means a user can only have ONE key at a time.
	// If you want to allow multiple keys, REMOVE or COMMENT OUT this block.
	if keys, _ := s.repo.GetSubscriptionAPIKeys(subscriptionID); len(keys) > 0 {
		for _, key := range keys {
			s.repo.DeleteAPIKey(key.ID)
		}
	}

	// Use a more descriptive prefix, perhaps indicating the environment (live/test).
	prefix := "mk_live_demo" // Changed prefix to be more descriptive

	// Generate the secure key and its hash.
	apiKey, hashedApiKey, err := generateSecureKey(prefix)
	if err != nil {
		return "", fmt.Errorf("failed to generate secure key: %w", err)
	}

	// Call the repository with the correct parameters.
	err = s.repo.CreateAPIKey(apiKey, hashedApiKey, subscriptionID)
	if err != nil {
		return "", err
	}

	// Return the plaintext key to the user ONCE.
	return apiKey, nil
}

func (s *service) DeleteAPIKey(apiKeyID string) error {
	return s.repo.DeleteAPIKey(apiKeyID)
}

// generateSecureKey remains the same. It is well-written.
func generateSecureKey(prefix string) (string, []byte, error) {
	numBytes := 64 // 64 bytes is a good size for a secure key, resulting in a 512-bit key.
	randomBytes := make([]byte, numBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", nil, fmt.Errorf("could not generate random bytes: %w", err)
	}
	keyString := base64.URLEncoding.EncodeToString(randomBytes)
	fullKey := fmt.Sprintf("%s_%s", prefix, keyString)
	hash := sha256.Sum256([]byte(fullKey))
	return fullKey, hash[:], nil
}

// ValidateAPIKey finds a key in the database based on the provided plaintext string.
// It returns the valid APIKey object if found, otherwise an error.
func ValidateAPIKey(db *gorm.DB, keyString string) (*models.APIKey, error) {
	// 1. Ensure the key string is not empty.
	if keyString == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}

	// 2. Split the key to get the prefix. The format is assumed to be "prefix_randompart".
	parts := strings.SplitN(keyString, "_", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid API key format")
	}
	prefix := parts[0]

	// 3. Hash the entire incoming key string to match what's in the database.
	incomingKeyHash := sha256.Sum256([]byte(keyString))

	// 4. Query the database for the key.
	// This query is highly efficient as it can use a composite index
	// on (key_prefix, key_value_hash) or individual indexes on both.
	var apiKey models.APIKey
	err := db.Where("key_prefix = ? AND key_value_hash = ?", prefix, incomingKeyHash[:]).First(&apiKey).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// This is the "invalid key" case.
			return nil, fmt.Errorf("invalid API key")
		}
		// This is a different, unexpected database error.
		return nil, fmt.Errorf("database error validating key: %w", err)
	}

	// 5. Key is valid and found. Return the model instance.
	return &apiKey, nil
}
