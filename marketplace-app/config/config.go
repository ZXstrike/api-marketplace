package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type AppConfig struct {
	ServerPort     string
	PrivateKey     *ecdsa.PrivateKey
	PublicKey      *ecdsa.PublicKey
	PostgresConfig PostgresConfig
	RedisConfig    RedisConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Database int
}

var Config *AppConfig

func LoadConfig() {
	// Load configuration from environment variables or a config file
	// For simplicity, we are hardcoding the values here

	// Load the private and public keys here
	privateKey, publicKey, err := LoadECDSAKeys()

	if err != nil {
		panic("Failed to load ECDSA keys: " + err.Error())
	}

	Config = &AppConfig{
		ServerPort: os.Getenv("SERVER_PORT"),
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		PostgresConfig: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_DB"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DATABASE"),
		},
	}
}

func LoadECDSAKeys() (privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, err error) {
	privateKeyPEM, err := loadECDSAPrivateKey(os.Getenv("ECDSA_PRIVATE_KEY"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load ECDSA private key: %w", err)
	}

	publicKeyPEM, err := loadECDSAPublicKey(os.Getenv("ECDSA_PUBLIC_KEY"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load ECDSA public key: %w", err)
	}

	return privateKeyPEM, publicKeyPEM, nil
}

func loadECDSAPrivateKey(key string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("error decoding private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}

	return privateKey, nil
}

func loadECDSAPublicKey(key string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("error decoding public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}

	return publicKey.(*ecdsa.PublicKey), nil
}
