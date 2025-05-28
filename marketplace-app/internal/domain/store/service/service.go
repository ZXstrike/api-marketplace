package service

import (
	"context"
	"crypto/ecdsa"

	"github.com/ZXstrike/marketplace-app/internal/domain/store/repositories"
	"github.com/ZXstrike/shared/pkg/models"
)

type Service interface {
	GetStoreByUserID(ctx context.Context, userID string) (*models.User, error)
	GetStoreByUsername(ctx context.Context, username string) (*models.User, error)
	GetAllStores(ctx context.Context) ([]models.User, error)
	CreateStore(ctx context.Context, userID string) error
	UpdateStore(ctx context.Context, userID string, description string) error
}

type service struct {
	repo       repositories.Repository
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func New(repo repositories.Repository, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) Service {
	return &service{repo, privateKey, publicKey}
}

func (s *service) GetStoreByUserID(ctx context.Context, userID string) (*models.User, error) {
	return s.repo.GetStoreByUserID(ctx, userID)
}

func (s *service) GetStoreByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetStoreByUsername(ctx, username)
}

func (s *service) GetAllStores(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAllStores(ctx)
}

func (s *service) CreateStore(ctx context.Context, userID string) error {
	return s.repo.CreateStore(ctx, userID)
}

func (s *service) UpdateStore(ctx context.Context, userID string, description string) error {
	return s.repo.UpdateStore(ctx, userID, description)
}
