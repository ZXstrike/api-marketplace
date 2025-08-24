package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZXstrike/shared/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock type for the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) Register(ctx context.Context, email, password, username string) error {
	args := m.Called(ctx, email, password, username)
	return args.Error(0)
}

func (m *MockService) Login(ctx context.Context, email, password string) (string, string, *models.User, error) {
	args := m.Called(ctx, email, password)
	if args.Get(2) == nil {
		return args.String(0), args.String(1), nil, args.Error(3)
	}
	return args.String(0), args.String(1), args.Get(2).(*models.User), args.Error(3)
}

func (m *MockService) VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	args := m.Called(tokenStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.RegisteredClaims), args.Error(1)
}

func (m *MockService) RefreshToken(tokenStr string) (string, error) {
	args := m.Called(tokenStr)
	return args.String(0), args.Error(1)
}

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{
			"email":    "test@example.com",
			"password": "password123",
			"username": "testuser",
		}

		mockService.On("Register", mock.Anything, "test@example.com", "password123", "testuser").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterHandler(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid request", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(gin.H{})
		c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterHandler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Registration failed", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{
			"email":    "test@example.com",
			"password": "password123",
			"username": "testuser",
		}

		mockService.On("Register", mock.Anything, "test@example.com", "password123", "testuser").Return(errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterHandler(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{
			"email":    "test@example.com",
			"password": "password123",
		}

		mockService.On("Login", mock.Anything, "test@example.com", "password123").Return("test_token", "user123", &models.User{}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.LoginHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "test_token", response["token"])
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid request", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(gin.H{})
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.LoginHandler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{
			"email":    "test@example.com",
			"password": "password123",
		}

		mockService.On("Login", mock.Anything, "test@example.com", "password123").Return("", "", nil, errors.New("invalid credentials"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.LoginHandler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestRefreshHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{
			"token": "old_token",
		}

		mockService.On("RefreshToken", "old_token").Return("new_token", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RefreshHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "new_token", response["token"])
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid request", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(gin.H{})
		c.Request, _ = http.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RefreshHandler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failed to refresh token", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{
			"token": "old_token",
		}

		mockService.On("RefreshToken", "old_token").Return("", errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RefreshHandler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertExpectations(t)
	})
}
