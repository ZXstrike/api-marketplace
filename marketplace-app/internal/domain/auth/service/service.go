package service

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"time"

	"github.com/ZXstrike/marketplace-app/internal/domain/auth/repositories"
	"github.com/ZXstrike/shared/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Register(ctx context.Context, email, password, username string) error
	Login(ctx context.Context, email, password string) (string, string, *models.User, error)
	VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error)
	RefreshToken(tokenStr string) (string, error)
}

type service struct {
	repo       repositories.Repository
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func New(repo repositories.Repository, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) Service {
	return &service{repo, privateKey, publicKey}
}

func (s *service) Register(ctx context.Context, email, password, username string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := &models.User{
		Email:             email,
		Username:          username,
		PasswordHash:      string(hash),
		Description:       "",
		ProfilePictureURL: "",
	}

	var consumerRole *models.Role

	if consumerRole, err = s.repo.GetRoleByName(ctx, "consumer"); err != nil {
		if err == gorm.ErrRecordNotFound {
			s.repo.CreateRole(ctx, &models.Role{Name: "consumer", Description: "API Consumer"})
			if consumerRole, err = s.repo.GetRoleByName(ctx, "consumer"); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	user.Roles = []models.Role{*consumerRole}

	// Create a default consumer wallet for the new user
	user.Wallets = []models.Wallet{
		{WalletType: "consumer", Balance: 0},
	}

	return s.repo.Create(ctx, user)
}

func (s *service) Login(ctx context.Context, email, password string) (string, string, *models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.RegisteredClaims{
		Subject:   user.ID,
		Issuer:    "marketplace-app",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	})

	token, err := claims.SignedString(s.privateKey)
	if err != nil {
		return "", "", nil, err
	}

	return token, user.ID, user, nil
}

func (s *service) VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(*jwt.RegisteredClaims), nil
}

func (s *service) RefreshToken(tokenStr string) (string, error) {
	claims, err := s.VerifyToken(tokenStr)
	if err != nil {
		return "", err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return "", errors.New("token expired")
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	newToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	newTokenStr, err := newToken.SignedString(s.privateKey)
	if err != nil {
		return "", err
	}

	return newTokenStr, nil
}
