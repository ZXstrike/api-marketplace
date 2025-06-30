package service

import (
	"context"
	"crypto/ecdsa"
	"mime/multipart"

	"github.com/ZXstrike/marketplace-app/internal/domain/user/repositories"
	"github.com/ZXstrike/shared/pkg/models"
)

type Service interface {
	GetUserProfile(ctx context.Context, id string) (*models.User, error)
	UpdateUserProfile(ctx context.Context, id string, description string) error
	ChangeUserPassword(ctx context.Context, id string, oldPass string, newPass string) error
	UpdateUserProfilePicture(ctx context.Context, id string, file *multipart.FileHeader) (string, error)
}

type service struct {
	repo       repositories.Repository
	fileRepo   repositories.FileRepository
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func New(repo repositories.Repository, fileRepo repositories.FileRepository, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) Service {
	return &service{repo, fileRepo, privateKey, publicKey}
}

func (s *service) GetUserProfile(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)

}

func (s *service) UpdateUserProfile(ctx context.Context, id string, description string) error {
	return s.repo.Update(ctx, id, description)
}

func (s *service) ChangeUserPassword(ctx context.Context, id string, oldPass string, newPass string) error {
	return s.repo.ChangePassword(ctx, id, oldPass, newPass)
}

func (s *service) UpdateUserProfilePicture(ctx context.Context, id string, file *multipart.FileHeader) (string, error) {

	publicURL, err := s.fileRepo.SaveProfile(file, id)
	if err != nil {
		return "", err
	}

	err = s.repo.UpdateProfilePicture(ctx, id, publicURL)
	if err != nil {
		return "", err
	}

	return publicURL, nil
}
