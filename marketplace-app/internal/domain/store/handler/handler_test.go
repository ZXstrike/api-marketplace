package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZXstrike/marketplace-app/internal/domain/store/service"
	"github.com/ZXstrike/shared/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock type for the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) GetStoreByUserID(ctx context.Context, userID string) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockService) GetStoreByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockService) GetAllStores(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockService) CreateStore(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockService) UpdateStore(ctx context.Context, userID string, description string) error {
	args := m.Called(ctx, userID, description)
	return args.Error(0)
}

func (m *MockService) GetStoreApis(ctx context.Context, userID string) ([]service.ApisData, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]service.ApisData), args.Error(1)
}

func TestGetStoreByUserIDHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedStore := &models.User{ID: "user123", Username: "testuser"}
		mockService.On("GetStoreByUserID", mock.Anything, "user123").Return(expectedStore, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "userID", Value: "user123"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/store/user/user123", nil)

		handler.GetStoreByUserIDHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, *expectedStore, response)
		mockService.AssertExpectations(t)
	})

	t.Run("Failed to get store", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetStoreByUserID", mock.Anything, "user123").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "userID", Value: "user123"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/store/user/user123", nil)

		handler.GetStoreByUserIDHandler(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetStoreByUsernameHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedStore := &models.User{ID: "user123", Username: "testuser"}
		mockService.On("GetStoreByUsername", mock.Anything, "testuser").Return(expectedStore, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "username", Value: "testuser"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/store/username/testuser", nil)

		handler.GetStoreByUsernameHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, *expectedStore, response)
		mockService.AssertExpectations(t)
	})

	t.Run("Failed to get store", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetStoreByUsername", mock.Anything, "testuser").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "username", Value: "testuser"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/store/username/testuser", nil)

		handler.GetStoreByUsernameHandler(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetAllStoresHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedStores := []models.User{{ID: "user123", Username: "testuser"}}
		mockService.On("GetAllStores", mock.Anything).Return(expectedStores, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/stores", nil)

		handler.GetAllStoresHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedStores, response)
		mockService.AssertExpectations(t)
	})

	t.Run("Failed to get stores", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetAllStores", mock.Anything).Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/stores", nil)

		handler.GetAllStoresHandler(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCreateStoreHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("CreateStore", mock.Anything, "user123").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodPost, "/store", nil)

		handler.CreateStoreHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/store", nil)

		handler.CreateStoreHandler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failed to create store", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("CreateStore", mock.Anything, "user123").Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodPost, "/store", nil)

		handler.CreateStoreHandler(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateStoreHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{"description": "new description"}
		mockService.On("UpdateStore", mock.Anything, "user123", "new description").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPut, "/store", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateStoreHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(gin.H{"description": "new description"})
		c.Request, _ = http.NewRequest(http.MethodPut, "/store", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateStoreHandler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid request", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(gin.H{})
		c.Request, _ = http.NewRequest(http.MethodPut, "/store", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateStoreHandler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to update store", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{"description": "new description"}
		mockService.On("UpdateStore", mock.Anything, "user123", "new description").Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPut, "/store", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateStoreHandler(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetStoreApisHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedApis := []service.ApisData{{ID: "api1", Name: "Test API", SubsCount: 10}}
		mockService.On("GetStoreApis", mock.Anything, "user123").Return(expectedApis, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/store/apis", nil)

		handler.GetStoreApisHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []service.ApisData
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedApis, response)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/store/apis", nil)

		handler.GetStoreApisHandler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failed to get store APIs", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetStoreApis", mock.Anything, "user123").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/store/apis", nil)

		handler.GetStoreApisHandler(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
