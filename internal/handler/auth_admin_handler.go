package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/dto"
	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/model"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/repository"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthAdminHandler struct {
	AdminUserRepository             repository.AdminUserRepository
	AdminUserRefreshTokenRepository repository.AdminUserRefreshTokenRepository
}

func NewAuthAdminHandler(
	adminUserRepo repository.AdminUserRepository,
	adminUserRefreshTokenRepo repository.AdminUserRefreshTokenRepository,
) *AuthAdminHandler {
	return &AuthAdminHandler{
		AdminUserRepository:             adminUserRepo,
		AdminUserRefreshTokenRepository: adminUserRefreshTokenRepo,
	}
}

func (h *AuthAdminHandler) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := dto.AdminUserLoginRequest{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to decode request: %w", err), http.StatusBadRequest)
			return
		}

		// Login logic
		adminUser, err := h.AdminUserRepository.GetByEmail(r.Context(), req.Email)
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to get adminUser: %w", err), http.StatusUnauthorized)
			return
		}

		if adminUser == nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("adminUser not found"), http.StatusUnauthorized)
			return
		}

		// Compare the password
		err = bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(req.Password))
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("invalid password and/or email"), http.StatusUnauthorized)
			return
		}

		// generate token access and refresh tokens
		// Generate the access token (short-lived)
		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		_, accessTokenString, err := tokenAuth.Encode(map[string]interface{}{
			"sub": strconv.Itoa(adminUser.ID),
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(15 * time.Minute).Unix(), // Access token valid for 15 minutes
		})
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to generate access token: %w", err), http.StatusInternalServerError)
			return
		}

		// Generate the refresh token (long-lived)
		_, refreshTokenString, err := tokenAuth.Encode(map[string]interface{}{
			"sub": strconv.Itoa(adminUser.ID),
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(30 * 24 * time.Hour).Unix(), // Refresh token valid for 30 days
		})
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to generate refresh token: %w", err), http.StatusInternalServerError)
			return
		}

		refreshToken := model.AdminUserRefreshToken{
			Token:     refreshTokenString,
			UserID:    adminUser.ID,
			UserAgent: r.UserAgent(),
			IPAddress: r.RemoteAddr,
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		}

		h.AdminUserRefreshTokenRepository.Create(r.Context(), &refreshToken)

		adminUserResponseDTO := dto.AdminUserResponse{
			ID:        adminUser.ID,
			FirstName: adminUser.FirstName,
			LastName:  adminUser.LastName,
			Email:     adminUser.Email,
		}

		repsonseDTO := dto.AdminUserLoginResponse{
			AdminUser:    adminUserResponseDTO,
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
		}

		APIResponse.SuccessResponse(w, r, APIResponse.APIResponse{
			Success: true,
			Data:    repsonseDTO,
		}, http.StatusOK)
	}
}

func (h *AuthAdminHandler) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := dto.AdminUserCreateRequest{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to decode request: %w", err), http.StatusBadRequest)
			return
		}

		// Check for missing or empty fields
		if req.FirstName == "" || req.LastName == "" || req.Email == "" ||
			req.Password == "" {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("missing required fields"), http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to hash password: %w", err), http.StatusInternalServerError)
			return
		}

		// Create a new adminUser
		adminUser := model.AdminUser{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  string(hashedPassword),
		}

		// Save adminUser to database
		newUser, err := h.AdminUserRepository.Create(r.Context(), &adminUser)
		if err != nil {
			APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		responseDTO := dto.AdminUserResponse{
			ID:        newUser.ID,
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Email:     newUser.Email,
		}

		APIResponse.SuccessResponse(w, r, &responseDTO, http.StatusCreated)
	}
}

func (h *AuthAdminHandler) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logout logic
		APIResponse.SuccessResponse(w, r, nil)
	}
}

func (h *AuthAdminHandler) RefreshTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Refresh token logic
		req := dto.UserRefreshTokenRequest{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to decode request: %w", err), http.StatusBadRequest)
			return
		}

		// Get the refresh token from the database
		refreshToken, err := h.AdminUserRefreshTokenRepository.GetByToken(r.Context(), req.RefreshToken)
		if err != nil {
			APIResponse.ErrorResponse(w, r, err, http.StatusNotFound)
			return
		}

		if refreshToken == nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("refresh token not found"), http.StatusNotFound)
			return
		}

		if refreshToken.ExpiresAt.Unix() < time.Now().Unix() {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("refresh token expired"), http.StatusUnauthorized)
			return
		}

		// Generate the access token (short-lived)
		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		_, accessTokenString, err := tokenAuth.Encode(map[string]interface{}{
			"sub":  strconv.Itoa(refreshToken.UserID),
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(15 * time.Minute).Unix(), // Access token valid for 15 minutes
			"role": "admin",
		})
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to generate access token: %w", err), http.StatusInternalServerError)
			return
		}

		responseDTO := dto.UserRefreshTokenResponse{
			AccessToken:  accessTokenString,
			RefreshToken: req.RefreshToken,
		}

		APIResponse.SuccessResponse(w, r, responseDTO, http.StatusOK)
	}
}
