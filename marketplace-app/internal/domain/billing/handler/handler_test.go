package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ZXstrike/shared/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock type for the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) GetBillingInfo(userID, walletType string) (map[string]interface{}, error) {
	args := m.Called(userID, walletType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockService) TopUp(userID, walletType string, amount float64) error {
	args := m.Called(userID, walletType, amount)
	return args.Error(0)
}

func (m *MockService) ProcessPayment(userID, walletType string, amount float64) error {
	args := m.Called(userID, walletType, amount)
	return args.Error(0)
}

func (m *MockService) GetPaymentHistory(userID string) ([]models.PaymentTransaction, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PaymentTransaction), args.Error(1)
}

func (m *MockService) ProcessAndCreatePayout(ctx context.Context, payoutData *models.ProviderPayout) (*models.ProviderPayout, error) {
	args := m.Called(ctx, payoutData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProviderPayout), args.Error(1)
}

func TestGetBillingInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedInfo := map[string]interface{}{"balance": 100.0}
		mockService.On("GetBillingInfo", "user123", "consumer").Return(expectedInfo, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/billing?wallet_type=consumer", nil)

		handler.GetBillingInfo(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedInfo, response)
		mockService.AssertExpectations(t)
	})

	t.Run("User not authenticated", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/billing?wallet_type=consumer", nil)

		handler.GetBillingInfo(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Missing wallet_type", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/billing", nil)

		handler.GetBillingInfo(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to retrieve billing info", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetBillingInfo", "user123", "consumer").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/billing?wallet_type=consumer", nil)

		handler.GetBillingInfo(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTopUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{"amount": 50.0}
		mockService.On("TopUp", "user123", "consumer", 50.0).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing/topup", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.TopUp(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("User not authenticated", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(gin.H{"amount": 50.0})
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing/topup", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.TopUp(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid request", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(gin.H{})
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing/topup", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.TopUp(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to top up", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{"amount": 50.0}
		mockService.On("TopUp", "user123", "consumer", 50.0).Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing/topup", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.TopUp(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProcessPayment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{"amount": 20.0, "wallet_type": "consumer"}
		mockService.On("ProcessPayment", "user123", "consumer", 20.0).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing/payment", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessPayment(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("User not authenticated", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(gin.H{"amount": 20.0, "wallet_type": "consumer"})
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing/payment", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessPayment(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid request", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(gin.H{})
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing/payment", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessPayment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to process payment", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{"amount": 20.0, "wallet_type": "consumer"}
		mockService.On("ProcessPayment", "user123", "consumer", 20.0).Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing/payment", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessPayment(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetPaymentHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedHistory := []models.PaymentTransaction{{ID: "1", Amount: 10.0}}
		mockService.On("GetPaymentHistory", "user123").Return(expectedHistory, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/billing/history", nil)

		handler.GetPaymentHistory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string][]models.PaymentTransaction
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedHistory, response["payment_history"])
		mockService.AssertExpectations(t)
	})

	t.Run("User not authenticated", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/billing/history", nil)

		handler.GetPaymentHistory(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failed to retrieve history", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetPaymentHistory", "user123").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/billing/history", nil)

		handler.GetPaymentHistory(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestHandleCreatePayout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		startTime := time.Now()
		endTime := time.Now().Add(24 * time.Hour)

		reqBody := gin.H{
			"provider_id":  "prov123",
			"period_start": startTime,
			"period_end":   endTime,
		}

		payoutData := &models.ProviderPayout{
			ProviderID:  "prov123",
			PeriodStart: startTime,
			PeriodEnd:   endTime,
			GrossAmount: 1000.00,
			PlatformFee: 50.00,
			NetAmount:   950.00,
		}

		mockService.On("ProcessAndCreatePayout", mock.Anything, mock.AnythingOfType("*models.ProviderPayout")).Return(payoutData, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing/payout", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.HandleCreatePayout(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.ProviderPayout
		json.Unmarshal(w.Body.Bytes(), &response)
		// Can't directly compare time objects due to potential monotonic clock differences
		assert.Equal(t, payoutData.ProviderID, response.ProviderID)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid request", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(gin.H{})
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing/payout", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.HandleCreatePayout(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to create payout", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		startTime := time.Now()
		endTime := time.Now().Add(24 * time.Hour)

		reqBody := gin.H{
			"provider_id":  "prov123",
			"period_start": startTime,
			"period_end":   endTime,
		}

		mockService.On("ProcessAndCreatePayout", mock.Anything, mock.AnythingOfType("*models.ProviderPayout")).Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing/payout", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.HandleCreatePayout(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
