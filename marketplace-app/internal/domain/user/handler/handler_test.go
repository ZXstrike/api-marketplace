package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func (m *MockService) GetUserProfile(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockService) UpdateUserProfile(ctx context.Context, id string, description string) error {
	args := m.Called(ctx, id, description)
	return args.Error(0)
}

func (m *MockService) ChangeUserPassword(ctx context.Context, id string, oldPass string, newPass string) error {
	args := m.Called(ctx, id, oldPass, newPass)
	return args.Error(0)
}

func (m *MockService) UpdateUserProfilePicture(ctx context.Context, id string, file *multipart.FileHeader) (string, error) {
	args := m.Called(ctx, id, file)
	return args.String(0), args.Error(1)
}

func TestUserProfileHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedUser := &models.User{ID: "user123", Username: "testuser"}
		mockService.On("GetUserProfile", mock.Anything, "user123").Return(expectedUser, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "user123"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/user/profile/user123", nil)

		handler.UserProfileHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, *expectedUser, response)
		mockService.AssertExpectations(t)
	})

	t.Run("Failed to get profile", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		mockService.On("GetUserProfile", mock.Anything, "user123").Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "user123"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/user/profile/user123", nil)

		handler.UserProfileHandler(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetMyProfileHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedUser := &models.User{ID: "user123", Username: "testuser"}
		mockService.On("GetUserProfile", mock.Anything, "user123").Return(expectedUser, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")
		c.Request, _ = http.NewRequest(http.MethodGet, "/user/me", nil)

		handler.GetMyProfileHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, *expectedUser, response)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/user/me", nil)

		handler.GetMyProfileHandler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestUpdateUserProfileHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{"description": "new description"}
		mockService.On("UpdateUserProfile", mock.Anything, "user123", "new description").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPut, "/user/profile", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateUserProfileHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(gin.H{"description": "new description"})
		c.Request, _ = http.NewRequest(http.MethodPut, "/user/profile", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateUserProfileHandler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid request", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(gin.H{})
		c.Request, _ = http.NewRequest(http.MethodPut, "/user/profile", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateUserProfileHandler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestChangePasswordHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		reqBody := gin.H{"old_password": "old_password", "new_password": "new_password"}
		mockService.On("ChangeUserPassword", mock.Anything, "user123", "old_password", "new_password").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/user/password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ChangePasswordHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(gin.H{"old_password": "old_password", "new_password": "new_password"})
		c.Request, _ = http.NewRequest(http.MethodPost, "/user/password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ChangePasswordHandler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid request", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body, _ := json.Marshal(gin.H{})
		c.Request, _ = http.NewRequest(http.MethodPost, "/user/password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ChangePasswordHandler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdateProfilePictureHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		// Create a dummy file for upload
		tmpfile, err := os.CreateTemp("", "test.png")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())
		tmpfile.WriteString("dummy content")
		tmpfile.Close()

		mockService.On("UpdateUserProfilePicture", mock.Anything, "user123", mock.AnythingOfType("*multipart.FileHeader")).Return("http://example.com/new_pic.png", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("profile_picture", filepath.Base(tmpfile.Name()))
		assert.NoError(t, err)
		file, err := os.Open(tmpfile.Name())
		assert.NoError(t, err)
		defer file.Close()
		io.Copy(part, file)
		writer.Close()

		c.Request, _ = http.NewRequest(http.MethodPost, "/user/picture", body)
		c.Request.Header.Set("Content-Type", writer.FormDataContentType())

		handler.UpdateProfilePictureHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetUserProfilePictureHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedUser := &models.User{ID: "user123", ProfilePictureURL: "http://example.com/pic.png"}
		mockService.On("GetUserProfile", mock.Anything, "user123").Return(expectedUser, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "user123"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/user/picture/user123", nil)

		handler.GetUserProfilePictureHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "http://example.com/pic.png", response["profile_picture"])
		mockService.AssertExpectations(t)
	})

	t.Run("User ID is required", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/user/picture/", nil)

		handler.GetUserProfilePictureHandler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Profile picture not found", func(t *testing.T) {
		mockService := new(MockService)
		handler := New(mockService)

		expectedUser := &models.User{ID: "user123", ProfilePictureURL: ""}
		mockService.On("GetUserProfile", mock.Anything, "user123").Return(expectedUser, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "user123"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/user/picture/user123", nil)

		handler.GetUserProfilePictureHandler(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}
