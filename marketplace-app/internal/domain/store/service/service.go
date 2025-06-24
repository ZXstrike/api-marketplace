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
	GetStoreApis(ctx context.Context, userID string) ([]ApisData, error)
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

type ApisData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SubsCount int    `json:"subs_count"`
}

func (s *service) GetStoreApis(ctx context.Context, userID string) ([]ApisData, error) {
	store, err := s.repo.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, nil // No store found for the user
	}

	apis, err := s.repo.GetStoreApis(ctx, store.ID)
	if err != nil {
		return nil, err
	}

	var data []ApisData

	for _, api := range apis {

		subsCount, err := s.repo.GetApiVersionsSubsCount(ctx, api.ID)
		if err != nil {
			return nil, err
		}
		data = append(data, ApisData{
			ID:        api.ID,
			Name:      api.Name,
			SubsCount: subsCount,
		})
	}

	return data, nil
}
