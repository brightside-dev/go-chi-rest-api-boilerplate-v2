package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/database"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*model.User), args.Error(1)
}

// MockRefreshTokenRepository is a mock implementation of UserRefreshTokenRepository
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *model.UserRefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetByToken(ctx context.Context, token string) (*model.UserRefreshToken, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*model.UserRefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

// MockEmailService is a mock implementation of EmailService
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) Send(templateName, subject string, to []string, data map[string]string) error {
	args := m.Called(templateName, subject, to, data)
	return args.Error(0)
}

func TestAuthService_Login(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)
	mockEmailService := new(MockEmailService)
	db := new(database.Service)
	logger := new(slog.Logger)

	authService := NewAuthService(db, logger, mockEmailService, mockUserRepo, mockRefreshTokenRepo)

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password",
	})
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	mockUserRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(&model.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}, nil)

	mockRefreshTokenRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	userLoginResponse, err := authService.Login(rr, req)
	assert.NoError(t, err)
	assert.NotEmpty(t, userLoginResponse.AccessToken)
	assert.NotEmpty(t, userLoginResponse.RefreshToken)
	mockUserRepo.AssertExpectations(t)
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthService_Register(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)
	mockEmailService := new(MockEmailService)
	db := new(database.Service)
	logger := new(slog.Logger)

	authService := NewAuthService(db, logger, mockEmailService, mockUserRepo, mockRefreshTokenRepo)

	reqBody, _ := json.Marshal(map[string]string{
		"firstName": "John",
		"lastName":  "Doe",
		"email":     "test@example.com",
		"password":  "password",
		"country":   "Australia",
		"birthday":  "1990-01-01",
	})
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockUserRepo.On("CreateUser", mock.Anything, mock.Anything).Return(&model.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "hashedpassword",
	}, nil)

	mockEmailService.On("Send", "welcome_email", "Welcome", []string{"test@example.com"}, nil).Return(nil)

	userResponse, err := authService.Register(rr, req)
	assert.NoError(t, err)
	assert.Equal(t, "John.Doe", userResponse.Name)
	mockUserRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

func TestAuthService_Logout(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)
	mockEmailService := new(MockEmailService)
	db := new(database.Service)
	logger := new(slog.Logger)

	authService := NewAuthService(db, logger, mockEmailService, mockUserRepo, mockRefreshTokenRepo)

	reqBody, _ := json.Marshal(map[string]string{
		"refreshToken": "test-refresh-token",
	})
	req, err := http.NewRequest("POST", "/logout", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockRefreshTokenRepo.On("GetByToken", mock.Anything, "test-refresh-token").Return(&model.UserRefreshToken{
		Token: "test-refresh-token",
	}, nil)

	mockRefreshTokenRepo.On("DeleteByToken", mock.Anything, "test-refresh-token").Return(nil)

	err = authService.Logout(rr, req)
	assert.NoError(t, err)
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthService_RefreshToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)
	mockEmailService := new(MockEmailService)
	db := new(database.Service)
	logger := new(slog.Logger)

	authService := NewAuthService(db, logger, mockEmailService, mockUserRepo, mockRefreshTokenRepo)

	reqBody, _ := json.Marshal(map[string]string{
		"refreshToken": "test-refresh-token",
	})
	req, err := http.NewRequest("POST", "/refresh-token", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockRefreshTokenRepo.On("GetByToken", mock.Anything, "test-refresh-token").Return(&model.UserRefreshToken{
		Token:     "test-refresh-token",
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}, nil)

	userRefreshTokenResponse, err := authService.RefreshToken(rr, req)
	assert.NoError(t, err)
	assert.NotEmpty(t, userRefreshTokenResponse.AccessToken)
	mockRefreshTokenRepo.AssertExpectations(t)
}
