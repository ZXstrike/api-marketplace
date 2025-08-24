package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock type for the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateAPIKey(subscriptionID string) (string, error) {
	args := m.Called(subscriptionID)
	return args.String(0), args.Error(1)
}

func (m *MockService) DeleteAPIKey(apiKeyID string) error {
	args := m.Called(apiKeyID)
	return args.Error(0)
}

func TestCreateAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{
			"subscription_id": "sub123",
		}

		mockService.On("CreateAPIKey", "sub123").Return("new_api_key", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api-key", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateAPIKey(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "new_api_key", response["api_key"])
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid request data", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(gin.H{})
		c.Request, _ = http.NewRequest(http.MethodPost, "/api-key", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateAPIKey(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to create API key", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{
			"subscription_id": "sub123",
		}

		mockService.On("CreateAPIKey", "sub123").Return("", errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api-key", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateAPIKey(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("DeleteAPIKey", "key123").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, "/api-key?api_key_id=key123", nil)

		handler.DeleteAPIKey(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "API key deleted successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("Failed to delete API key", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("DeleteAPIKey", "key123").Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, "/api-key?api_key_id=key123", nil)

		handler.DeleteAPIKey(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
