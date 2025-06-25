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

func LoadConfig() (*AppConfig, error) {
	// Load configuration from environment variables or a config file
	// For simplicity, we are hardcoding the values here

	// Load the private and public keys here
	privateKey, publicKey, err := LoadECDSAKeys()

	if err != nil {
		panic("Failed to load ECDSA keys: " + err.Error())
	}

	Config = &AppConfig{
		ServerPort: os.Getenv("GATEWAY_PORT"),
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		PostgresConfig: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DATABASE"),
		},
	}

	Config.RedisConfig = RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Database: 0, // Default Redis database
	}

	return Config, nil
}

func LoadECDSAKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKeyPath := os.Getenv("ECDSA_PRIVATE_KEY_PATH")
	publicKeyPath := os.Getenv("ECDSA_PUBLIC_KEY_PATH")

	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read ECDSA private key file: %w", err)
	}

	publicKeyPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read ECDSA public key file: %w", err)
	}

	privateKey, err := parseECDSAPrivateKey(privateKeyPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing private key: %w", err)
	}

	publicKey, err := parseECDSAPublicKey(publicKeyPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing public key: %w", err)
	}

	return privateKey, publicKey, nil
}

func parseECDSAPrivateKey(pemBytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("invalid private key PEM")
	}
	return x509.ParseECPrivateKey(block.Bytes)
}

func parseECDSAPublicKey(pemBytes []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("invalid public key PEM")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*ecdsa.PublicKey), nil
}
