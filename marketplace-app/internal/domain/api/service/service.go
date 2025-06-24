package service

import (
	"crypto/ecdsa"
	"errors"
	"net/url"

	"github.com/ZXstrike/marketplace-app/internal/domain/api/repositories"
	"github.com/ZXstrike/shared/pkg/models"
)

type Service interface {
	GetAPIByID(id string) (*models.API, error)
	CreateNewAPI(name string, desc string, providerId string, baseUrl string, pricePerCall float64, categories []string) (string, error)
	UpdateAPI(apiID string, name string, desc string, baseUrl string, pricePerCall float64, categories []string) error
	DeleteAPI(userId string, apiId string) error
	GetAllAPIs(page int, lenght int) ([]models.API, error)
	GetAllAPIsByUserID(userID string) ([]models.API, error)
	CreateAPIEndpoint(apiVersion string, endpoints []models.Endpoint) error
	UpdateAPIEndpoint(apiVersion string, endpoints []models.Endpoint) error
	DeleteAPIEndpoint(endpointID string) error
	GetAllEndpointsByAPIVersionID(apiVersionID string) ([]models.Endpoint, error)
	GetAllCategories() ([]models.Category, error)
}

type service struct {
	repo       repositories.Repository
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func New(repo repositories.Repository, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) Service {
	return &service{repo, privateKey, publicKey}
}
func (s *service) CreateNewAPI(name string, desc string, providerId string, baseUrl string, pricePerCall float64, categories []string) (string, error) {

	user, err := s.repo.GetUserByID(providerId)

	if err != nil {
		return "", err
	}

	var categoriesModels []models.Category
	for _, categoryName := range categories {
		category, err := s.repo.GetCategoryBySlug(categoryName)
		if err != nil {
			return "", err
		}
		if category == nil {
			return "", errors.New("category not found: " + categoryName)
		}
		categoriesModels = append(categoriesModels, *category)
	}

	if _, err := url.ParseRequestURI(baseUrl); err != nil {
		return "", errors.New("invalid base URL format")
	}

	api := models.API{
		Name:        name,
		Description: desc,
		ProviderID:  user.ID,
		Provider:    *user,
		BaseURL:     baseUrl,
		Categories:  categoriesModels,
	}

	if pricePerCall < 0 {
		return "", errors.New("price per call cannot be negative")
	}

	apiID, err := s.repo.CreateAPI(api, pricePerCall)
	if err != nil {
		return "", err
	}

	return apiID, nil
}

func (s *service) UpdateAPI(apiID string, name string, desc string, baseUrl string, pricePerCall float64, categories []string) error {
	api, err := s.repo.GetAPIByID(apiID)
	if err != nil {
		return err
	}

	if api == nil {
		return errors.New("API not found")
	}

	if _, err := url.ParseRequestURI(baseUrl); err != nil {
		return errors.New("invalid base URL format")
	}

	var categoriesModels []models.Category
	for _, categoryName := range categories {
		category, err := s.repo.GetCategoryBySlug(categoryName)
		if err != nil {
			return err
		}
		if category == nil {
			return errors.New("category not found: " + categoryName)
		}
		categoriesModels = append(categoriesModels, *category)
	}

	api.Name = name
	api.Description = desc
	api.BaseURL = baseUrl
	api.Categories = categoriesModels

	if pricePerCall < 0 {
		return errors.New("price per call cannot be negative")
	}

	return s.repo.UpdateAPI(*api)
}

func (s *service) DeleteAPI(userId string, apiId string) error {
	// Check if the user is the owner of the API
	api, err := s.repo.GetAPIByID(apiId)
	if err != nil {
		return err
	}

	if api == nil {
		return errors.New("API not found")
	}

	if api.ProviderID != userId {
		return errors.New("user does not have permission to delete this API")
	}

	return s.repo.DeleteAPI(apiId)
}

func (s *service) GetAPIByID(id string) (*models.API, error) {
	return s.repo.GetAPIByID(id)
}

func (s *service) GetAllAPIs(page int, lenght int) ([]models.API, error) {
	return s.repo.GetAllAPI(page, lenght)
}

func (s *service) GetAllAPIsByUserID(userID string) ([]models.API, error) {
	return s.repo.GetAllAPIByUserID(userID)
}

func (s *service) CreateAPIEndpoint(apiVersion string, endpoints []models.Endpoint) error {
	apiVersionData, err := s.repo.GetAPIVersionByID(apiVersion)

	if err != nil {
		return err
	}

	for _, endpoint := range endpoints {
		endpoint.APIVersionID = apiVersionData.ID
		endpoint.APIVersion = *apiVersionData

		if err := s.repo.CreateAPIEndpoint(endpoint); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) UpdateAPIEndpoint(apiVersion string, endpoints []models.Endpoint) error {

	apiVersionData, err := s.repo.GetAPIVersionByID(apiVersion)

	if err != nil {
		return err
	}

	for _, endpoint := range endpoints {

		endpointData, err := s.repo.GetAPIEndpointByID(endpoint.ID)

		if err != nil {

			endpoint.APIVersionID = apiVersionData.ID
			endpoint.APIVersion = *apiVersionData

			if err := s.repo.CreateAPIEndpoint(endpoint); err != nil {
				return err
			}
		} else {
			endpointData.HTTPMethod = endpoint.HTTPMethod
			endpointData.Path = endpoint.Path
			endpointData.Documentation = endpoint.Documentation
			if err := s.repo.UpdateAPIEndpoint(*endpointData); err != nil {
				return err
			}
		}

	}

	return nil
}

func (s *service) DeleteAPIEndpoint(endpointID string) error {
	endpoint, err := s.repo.GetAPIEndpointByID(endpointID)
	if err != nil {
		return err
	}

	if endpoint == nil {
		return errors.New("endpoint not found")
	}

	return s.repo.DeleteAPIEndpoint(endpointID)
}

func (s *service) GetAllEndpointsByAPIVersionID(apiVersionID string) ([]models.Endpoint, error) {
	apiVersion, err := s.repo.GetAPIVersionByID(apiVersionID)
	if err != nil {
		return nil, err
	}

	if apiVersion == nil {
		return nil, errors.New("API version not found")
	}

	endpoints, err := s.repo.GetAllEndpointsByAPIVersionID(apiVersionID)
	if err != nil {
		return nil, err
	}

	return endpoints, nil
}

func (s *service) GetAllCategories() ([]models.Category, error) {
	categories, err := s.repo.GetAllCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}
