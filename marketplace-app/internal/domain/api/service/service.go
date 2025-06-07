package service

import (
	"crypto/ecdsa"

	"github.com/ZXstrike/marketplace-app/internal/domain/api/repositories"
	"github.com/ZXstrike/shared/pkg/models"
)

type Service interface {
	GetAPIByID(id string) (*models.API, error)
	CreateNewAPI(name string, desc string, providerId string) error
	GetAllAPIs() ([]models.API, error)
	GetAllAPIsByUserID(userID string) ([]models.API, error)
	UpdateAPI(api models.API) error
	DeleteAPI(userId string, apiId string) error
}

type service struct {
	repo       repositories.Repository
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func New(repo repositories.Repository, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) Service {
	return &service{repo, privateKey, publicKey}
}
func (s *service) CreateNewAPI(name string, desc string, providerId string) error {

	user, err := s.repo.GetUserByID(providerId)

	if err != nil {
		return err
	}

	api := models.API{
		Name:        name,
		Description: desc,
		ProviderID:  user.ID,
		Provider:    *user,
	}

	return s.repo.Create(api)
}

func (s *service) GetAPIByID(id string) (*models.API, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetAllAPIs() ([]models.API, error) {
	return s.repo.GetAll()
}

func (s *service) GetAllAPIsByUserID(userID string) ([]models.API, error) {
	return s.repo.GetAllByUserID(userID)
}

func (s *service) UpdateAPI(api models.API) error {
	return s.repo.Update(api)
}

func (s *service) DeleteAPI(userId string, apiId string) error {
	return s.repo.Delete(apiId)
}
