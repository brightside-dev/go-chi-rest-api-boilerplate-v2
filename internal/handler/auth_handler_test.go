package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/dto"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/service/email"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(w http.ResponseWriter, r *http.Request) (dto.UserLoginResponse, error) {
	args := m.Called(w, r)
	return args.Get(0).(dto.UserLoginResponse), args.Error(1)
}

func (m *MockAuthService) Register(w http.ResponseWriter, r *http.Request) (dto.UserResponse, error) {
	args := m.Called(w, r)
	return args.Get(0).(dto.UserResponse), args.Error(1)
}

func (m *MockAuthService) Logout(w http.ResponseWriter, r *http.Request) error {
	args := m.Called(w, r)
	return args.Error(0)
}

func (m *MockAuthService) RefreshToken(w http.ResponseWriter, r *http.Request) (dto.UserRefreshTokenResponse, error) {
	args := m.Called(w, r)
	return args.Get(0).(dto.UserRefreshTokenResponse), args.Error(1)
}

func TestAuthHandler_Login(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockEmailService := new(email.EmailService)

	authHandler := NewAuthHandler(mockEmailService, mockAuthService)

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password",
	})
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockAuthService.On("Login", mock.Anything, mock.Anything).Return(dto.UserLoginResponse{
		User: dto.UserResponse{
			ID:      1,
			Name:    "John Doe",
			Country: "Australia",
		},
		AccessToken:  "token",
		RefreshToken: "refresh-token",
	}, nil)

	handler := authHandler.Login()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Register(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockEmailService := new(email.EmailService)

	authHandler := NewAuthHandler(mockEmailService, mockAuthService)

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password",
	})
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockAuthService.On("Register", mock.Anything, mock.Anything).Return(dto.UserResponse{
		ID:      1,
		Name:    "John Doe",
		Country: "Australia",
	}, nil)

	handler := authHandler.Register()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Logout(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockEmailService := new(email.EmailService)

	authHandler := NewAuthHandler(mockEmailService, mockAuthService)

	req, err := http.NewRequest("POST", "/logout", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockAuthService.On("Logout", mock.Anything, mock.Anything).Return(nil)

	handler := authHandler.Logout()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockEmailService := new(email.EmailService)

	authHandler := NewAuthHandler(mockEmailService, mockAuthService)

	req, err := http.NewRequest("POST", "/refresh-token", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockAuthService.On("RefreshToken", mock.Anything, mock.Anything).Return(dto.UserRefreshTokenResponse{
		AccessToken:  "token",
		RefreshToken: "refresh-token",
	}, nil)

	handler := authHandler.RefreshToken()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAuthService.AssertExpectations(t)
}
