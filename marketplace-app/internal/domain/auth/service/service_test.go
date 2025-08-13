package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"testing"

	"github.com/ZXstrike/shared/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// MockRepository is a mock implementation of the repositories.Repository interface.
type MockRepository struct {
	mock.Mock
}

// GetByID implements repositories.Repository.
func (m *MockRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	panic("unimplemented")
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRepository) CreateRole(ctx context.Context, role *models.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func TestRegister(t *testing.T) {
	// Generate dummy keys for service initialization
	privateKey, _ := ecdsa.GenerateKey(nil, rand.Reader)

	testCases := []struct {
		name          string
		email         string
		password      string
		username      string
		setupMock     func(mockRepo *MockRepository)
		expectErr     bool
		expectedError string
	}{
		{
			name:     "Successful Registration",
			email:    "test@example.com",
			password: "password123",
			username: "testuser",
			setupMock: func(mockRepo *MockRepository) {
				consumerRole := &models.Role{Name: "consumer"}
				mockRepo.On("GetRoleByName", mock.Anything, "consumer").Return(consumerRole, nil).Once()
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name:     "Successful Registration - Role Not Found, Then Created",
			email:    "newrole@example.com",
			password: "password123",
			username: "newroleuser",
			setupMock: func(mockRepo *MockRepository) {
				consumerRole := &models.Role{Name: "consumer"}
				// First call to GetRoleByName fails
				mockRepo.On("GetRoleByName", mock.Anything, "consumer").Return(nil, gorm.ErrRecordNotFound).Once()
				// Then we expect CreateRole to be called
				mockRepo.On("CreateRole", mock.Anything, mock.AnythingOfType("*models.Role")).Return(nil).Once()
				// Second call to GetRoleByName succeeds
				mockRepo.On("GetRoleByName", mock.Anything, "consumer").Return(consumerRole, nil).Once()
				// Finally, Create user is called
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name:     "Database Error on GetRoleByName",
			email:    "dberror@example.com",
			password: "password123",
			username: "dberroruser",
			setupMock: func(mockRepo *MockRepository) {
				mockRepo.On("GetRoleByName", mock.Anything, "consumer").Return(nil, errors.New("db connection failed")).Once()
			},
			expectErr:     true,
			expectedError: "db connection failed",
		},
		{
			name:     "Database Error on Create",
			email:    "createerror@example.com",
			password: "password123",
			username: "createerroruser",
			setupMock: func(mockRepo *MockRepository) {
				consumerRole := &models.Role{Name: "consumer"}
				mockRepo.On("GetRoleByName", mock.Anything, "consumer").Return(consumerRole, nil).Once()
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("failed to insert user")).Once()
			},
			expectErr:     true,
			expectedError: "failed to insert user",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			// Pass nil for keys as they aren't used in Register
			authService := New(mockRepo, privateKey, &privateKey.PublicKey)
			tc.setupMock(mockRepo)

			// Execute
			err := authService.Register(context.Background(), tc.email, tc.password, tc.username)

			// Assert
			if tc.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// Verify that all expected mock calls were made
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(nil, rand.Reader)
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	testUser := &models.User{
		ID:           "user-123",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}

	testCases := []struct {
		name          string
		email         string
		password      string
		setupMock     func(mockRepo *MockRepository)
		expectErr     bool
		expectedError string
	}{
		{
			name:     "Successful Login",
			email:    "test@example.com",
			password: "password123",
			setupMock: func(mockRepo *MockRepository) {
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(testUser, nil).Once()
			},
			expectErr: false,
		},
		{
			name:     "User Not Found",
			email:    "notfound@example.com",
			password: "password123",
			setupMock: func(mockRepo *MockRepository) {
				mockRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectErr:     true,
			expectedError: "invalid credentials",
		},
		{
			name:     "Incorrect Password",
			email:    "test@example.com",
			password: "wrongpassword",
			setupMock: func(mockRepo *MockRepository) {
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(testUser, nil).Once()
			},
			expectErr:     true,
			expectedError: "invalid credentials",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			authService := New(mockRepo, privateKey, &privateKey.PublicKey)
			tc.setupMock(mockRepo)

			token, userID, user, err := authService.Login(context.Background(), tc.email, tc.password)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Empty(t, token)
				assert.Empty(t, userID)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Equal(t, testUser.ID, userID)
				assert.NotNil(t, user)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
