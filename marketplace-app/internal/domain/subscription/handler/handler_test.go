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

func (m *MockService) SubscribeToAPI(userID string, apiVersionID string) error {
	args := m.Called(userID, apiVersionID)
	return args.Error(0)
}

func (m *MockService) UnsubscribeFromAPI(userID string, subscriptionID string) error {
	args := m.Called(userID, subscriptionID)
	return args.Error(0)
}

func (m *MockService) GetSubscription(userID string, subscriptionID string) (*models.Subscription, error) {
	args := m.Called(userID, subscriptionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockService) GetSubscriptionsByUserID(userID string) ([]models.Subscription, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Subscription), args.Error(1)
}

func TestSubscribeToAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{"api_version_id": "v1"}
		mockService.On("SubscribeToAPI", "user123", "v1").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.SubscribeToAPI(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid request", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(gin.H{})
		c.Request, _ = http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.SubscribeToAPI(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("User not found", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(gin.H{"api_version_id": "v1"})
		c.Request, _ = http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.SubscribeToAPI(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to subscribe", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{"api_version_id": "v1"}
		mockService.On("SubscribeToAPI", "user123", "v1").Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.SubscribeToAPI(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUnsubscribeFromAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{"subscription_id": "sub123"}
		mockService.On("UnsubscribeFromAPI", "user123", "sub123").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/unsubscribe", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UnsubscribeFromAPI(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid request", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(gin.H{})
		c.Request, _ = http.NewRequest(http.MethodPost, "/unsubscribe", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UnsubscribeFromAPI(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("User not found", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(gin.H{"subscription_id": "sub123"})
		c.Request, _ = http.NewRequest(http.MethodPost, "/unsubscribe", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UnsubscribeFromAPI(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to unsubscribe", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{"subscription_id": "sub123"}
		mockService.On("UnsubscribeFromAPI", "user123", "sub123").Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/unsubscribe", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UnsubscribeFromAPI(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedSub := &models.Subscription{ID: "sub123"}
		mockService.On("GetSubscription", "user123", "sub123").Return(expectedSub, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/subscription?subscriptionID=sub123", nil)

		handler.GetSubscription(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.Subscription
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, *expectedSub, response)
		mockService.AssertExpectations(t)
	})

	t.Run("Missing subscriptionID", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/subscription", nil)

		handler.GetSubscription(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("User not found", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/subscription?subscriptionID=sub123", nil)

		handler.GetSubscription(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to get subscription", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetSubscription", "user123", "sub123").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/subscription?subscriptionID=sub123", nil)

		handler.GetSubscription(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetSubscriptions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedSubs := []models.Subscription{{ID: "sub123"}}
		mockService.On("GetSubscriptionsByUserID", "user123").Return(expectedSubs, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/subscriptions", nil)

		handler.GetSubscriptions(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []models.Subscription
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedSubs, response)
		mockService.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/subscriptions", nil)

		handler.GetSubscriptions(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to get subscriptions", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetSubscriptionsByUserID", "user123").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/subscriptions", nil)

		handler.GetSubscriptions(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
