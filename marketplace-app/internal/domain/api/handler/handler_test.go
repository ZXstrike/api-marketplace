package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZXstrike/shared/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock type for the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) GetAPIByID(id string) (*models.API, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.API), args.Error(1)
}

func (m *MockService) GetAdminData() (map[string]interface{}, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockService) CreateNewAPI(name string, desc string, providerId string, baseUrl string, pricePerCall float64, categories []string) (string, error) {
	args := m.Called(name, desc, providerId, baseUrl, pricePerCall, categories)
	return args.String(0), args.Error(1)
}

func (m *MockService) UpdateAPI(apiID string, name string, desc string, baseUrl string, pricePerCall float64, categories []string) error {
	args := m.Called(apiID, name, desc, baseUrl, pricePerCall, categories)
	return args.Error(0)
}

func (m *MockService) DeleteAPI(userId string, apiId string) error {
	args := m.Called(userId, apiId)
	return args.Error(0)
}

func (m *MockService) GetAllAPIs(page int, length int, query string) ([]models.API, error) {
	args := m.Called(page, length, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.API), args.Error(1)
}

func (m *MockService) GetAllAPIsByUserID(userID string) ([]models.API, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.API), args.Error(1)
}

func (m *MockService) CreateAPIEndpoint(apiVersion string, endpoints []models.Endpoint) error {
	args := m.Called(apiVersion, endpoints)
	return args.Error(0)
}

func (m *MockService) UpdateAPIEndpoint(apiVersion string, endpoints []models.Endpoint) error {
	args := m.Called(apiVersion, endpoints)
	return args.Error(0)
}

func (m *MockService) DeleteAPIEndpoint(endpointID string) error {
	args := m.Called(endpointID)
	return args.Error(0)
}

func (m *MockService) GetAllEndpointsByAPIVersionID(apiVersionID string) ([]models.Endpoint, error) {
	args := m.Called(apiVersionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Endpoint), args.Error(1)
}

func (m *MockService) GetAllCategories() ([]models.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockService) GetConsumerOverview(userId string) (map[string]interface{}, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockService) GetProviderOverview(userId string) (map[string]interface{}, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestCreateNewAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		price := 10.0
		reqBody := gin.H{
			"name":           "Test API",
			"description":    "This is a test API",
			"base_url":       "http://test.com",
			"price_per_call": price,
			"categories":     []string{"test"},
		}

		mockService.On("CreateNewAPI", "Test API", "This is a test API", "user123", "http://test.com", price, []string{"test"}).Return("api123", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNewAPI(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "API created successfully", response["message"])
		assert.Equal(t, "api123", response["api_id"])
		mockService.AssertExpectations(t)
	})

	t.Run("No user authenticated", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		price := 10.0
		reqBody := gin.H{
			"name":           "Test API",
			"description":    "This is a test API",
			"base_url":       "http://test.com",
			"price_per_call": price,
			"categories":     []string{"test"},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNewAPI(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failed to create API", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		price := 10.0
		reqBody := gin.H{
			"name":           "Test API",
			"description":    "This is a test API",
			"base_url":       "http://test.com",
			"price_per_call": price,
			"categories":     []string{"test"},
		}

		mockService.On("CreateNewAPI", "Test API", "This is a test API", "user123", "http://test.com", price, []string{"test"}).Return("", errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNewAPI(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetAdminData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedData := map[string]interface{}{"users": 10.0, "apis": 5.0}
		mockService.On("GetAdminData").Return(expectedData, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/admin/data", nil)

		handler.GetAdminData(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedData, response)
		mockService.AssertExpectations(t)
	})

	t.Run("Failed to fetch admin data", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetAdminData").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/admin/data", nil)

		handler.GetAdminData(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		price := 20.0
		reqBody := gin.H{
			"name":           "Updated API",
			"description":    "This is an updated API",
			"base_url":       "http://updated.com",
			"price_per_call": price,
			"categories":     []string{"updated"},
		}

		mockService.On("UpdateAPI", "api123", "Updated API", "This is an updated API", "http://updated.com", price, []string{"updated"}).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "api123"}}

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPut, "/api/api123", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateAPI(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "API updated successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("API ID is required", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		price := 20.0
		reqBody := gin.H{
			"name":           "Updated API",
			"description":    "This is an updated API",
			"base_url":       "http://updated.com",
			"price_per_call": price,
			"categories":     []string{"updated"},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPut, "/api/", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateAPI(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to update API", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		price := 20.0
		reqBody := gin.H{
			"name":           "Updated API",
			"description":    "This is an updated API",
			"base_url":       "http://updated.com",
			"price_per_call": price,
			"categories":     []string{"updated"},
		}

		mockService.On("UpdateAPI", "api123", "Updated API", "This is an updated API", "http://updated.com", price, []string{"updated"}).Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "api123"}}

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPut, "/api/api123", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateAPI(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("DeleteAPI", "user123", "api123").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Params = gin.Params{gin.Param{Key: "id", Value: "api123"}}
		c.Request, _ = http.NewRequest(http.MethodDelete, "/api/api123", nil)

		handler.DeleteAPI(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "API deleted successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("API ID is required", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodDelete, "/api/", nil)

		handler.DeleteAPI(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("No user authenticated", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "api123"}}
		c.Request, _ = http.NewRequest(http.MethodDelete, "/api/api123", nil)

		handler.DeleteAPI(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failed to delete API", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("DeleteAPI", "user123", "api123").Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Params = gin.Params{gin.Param{Key: "id", Value: "api123"}}
		c.Request, _ = http.NewRequest(http.MethodDelete, "/api/api123", nil)

		handler.DeleteAPI(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetAllAPIs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedAPIs := []models.API{{ID: "1", Name: "Test API"}}
		mockService.On("GetAllAPIs", 1, 10, "test").Return(expectedAPIs, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/apis?page=1&length=10&query=test", nil)

		handler.GetAllAPIs(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []models.API
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedAPIs, response)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid page number", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/apis?page=abc", nil)

		handler.GetAllAPIs(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid length", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/apis?page=1&length=abc", nil)

		handler.GetAllAPIs(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to fetch APIs", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetAllAPIs", 1, 10, "").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/apis?page=1&length=10", nil)

		handler.GetAllAPIs(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetAPIByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedAPI := &models.API{ID: "api123", Name: "Test API"}
		mockService.On("GetAPIByID", "api123").Return(expectedAPI, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "api123"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/api123", nil)

		handler.GetAPIByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.API
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, *expectedAPI, response)
		mockService.AssertExpectations(t)
	})

	t.Run("API not found", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetAPIByID", "api123").Return(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "api123"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/api123", nil)

		handler.GetAPIByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Failed to fetch API", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetAPIByID", "api123").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "api123"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/api123", nil)

		handler.GetAPIByID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("API ID is required", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/", nil)

		handler.GetAPIByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetAllAPIsByUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedAPIs := []models.API{{ID: "1", Name: "Test API", ProviderID: "user123"}}
		mockService.On("GetAllAPIsByUserID", "user123").Return(expectedAPIs, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/user/apis", nil)

		handler.GetAllAPIsByUserID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []models.API
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedAPIs, response)
		mockService.AssertExpectations(t)
	})

	t.Run("No user authenticated", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/user/apis", nil)

		handler.GetAllAPIsByUserID(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failed to fetch APIs", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetAllAPIsByUserID", "user123").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/user/apis", nil)

		handler.GetAllAPIsByUserID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCreateNewAPIEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{
			"api_version_id": "v1",
			"endpoints": []gin.H{
				{
					"http_method":   "GET",
					"path":          "/test",
					"documentation": "Test endpoint",
				},
			},
		}

		endpoints := []models.Endpoint{
			{
				HTTPMethod:    "GET",
				Path:          "/test",
				Documentation: "Test endpoint",
			},
		}

		mockService.On("CreateAPIEndpoint", "v1", endpoints).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/endpoint", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNewAPIEndpoint(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "API endpoint created successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("Failed to create API endpoint", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{
			"api_version_id": "v1",
			"endpoints": []gin.H{
				{
					"http_method":   "GET",
					"path":          "/test",
					"documentation": "Test endpoint",
				},
			},
		}

		endpoints := []models.Endpoint{
			{
				HTTPMethod:    "GET",
				Path:          "/test",
				Documentation: "Test endpoint",
			},
		}

		mockService.On("CreateAPIEndpoint", "v1", endpoints).Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/endpoint", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNewAPIEndpoint(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateAPIEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{
			"api_version_id": "v1",
			"endpoints": []gin.H{
				{
					"endpoint_id":   "ep1",
					"http_method":   "PUT",
					"path":          "/test_updated",
					"documentation": "Test endpoint updated",
				},
			},
		}

		endpoints := []models.Endpoint{
			{
				ID:            "ep1",
				HTTPMethod:    "PUT",
				Path:          "/test_updated",
				Documentation: "Test endpoint updated",
			},
		}

		mockService.On("UpdateAPIEndpoint", "v1", endpoints).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPut, "/api/endpoint", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateAPIEndpoint(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "API endpoint updated successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("Failed to update API endpoint", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{
			"api_version_id": "v1",
			"endpoints": []gin.H{
				{
					"endpoint_id":   "ep1",
					"http_method":   "PUT",
					"path":          "/test_updated",
					"documentation": "Test endpoint updated",
				},
			},
		}

		endpoints := []models.Endpoint{
			{
				ID:            "ep1",
				HTTPMethod:    "PUT",
				Path:          "/test_updated",
				Documentation: "Test endpoint updated",
			},
		}

		mockService.On("UpdateAPIEndpoint", "v1", endpoints).Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPut, "/api/endpoint", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateAPIEndpoint(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteAPIEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("DeleteAPIEndpoint", "ep1").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "ep1"}}
		c.Request, _ = http.NewRequest(http.MethodDelete, "/api/endpoint/ep1", nil)

		handler.DeleteAPIEndpoint(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "API endpoint deleted successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("Endpoint ID is required", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, "/api/endpoint/", nil)

		handler.DeleteAPIEndpoint(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to delete API endpoint", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("DeleteAPIEndpoint", "ep1").Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "ep1"}}
		c.Request, _ = http.NewRequest(http.MethodDelete, "/api/endpoint/ep1", nil)

		handler.DeleteAPIEndpoint(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetAllAPIEndpointsByAPIVersionID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedEndpoints := []models.Endpoint{{ID: "ep1", HTTPMethod: "GET", Path: "/test"}}
		mockService.On("GetAllEndpointsByAPIVersionID", "v1").Return(expectedEndpoints, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "apiVersionID", Value: "v1"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/endpoints", nil)

		handler.GetAllAPIEndpointsByAPIVersionID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []models.Endpoint
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedEndpoints, response)
		mockService.AssertExpectations(t)
	})

	t.Run("API Version ID is required", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api//endpoints", nil)

		handler.GetAllAPIEndpointsByAPIVersionID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to fetch API endpoints", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetAllEndpointsByAPIVersionID", "v1").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "apiVersionID", Value: "v1"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/endpoints", nil)

		handler.GetAllAPIEndpointsByAPIVersionID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetAllCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedCategories := []models.Category{{ID: "1", Name: "Test Category"}}
		mockService.On("GetAllCategories").Return(expectedCategories, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/categories", nil)

		handler.GetAllCategories(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []models.Category
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedCategories, response)
		mockService.AssertExpectations(t)
	})

	t.Run("Failed to fetch categories", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetAllCategories").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/categories", nil)

		handler.GetAllCategories(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
