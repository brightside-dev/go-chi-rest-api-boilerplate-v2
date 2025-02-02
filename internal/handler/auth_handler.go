package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/email"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/dto"
	customError "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/error"
	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/model"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/repository"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/util"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	EmailService           *email.EmailService
	UserRepository         repository.UserRepository
	RefreshTokenRepository repository.UserRefreshTokenRepository
	dbLogger               *slog.Logger
}

func NewAuthHandler(
	emailService *email.EmailService,
	userRepo repository.UserRepository,
	refreshTokenRepo repository.UserRefreshTokenRepository,
	dbLogger *slog.Logger,
) *AuthHandler {
	return &AuthHandler{
		EmailService:           emailService,
		UserRepository:         userRepo,
		RefreshTokenRepository: refreshTokenRepo,
		dbLogger:               dbLogger,
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := dto.UserLoginRequest{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			APIResponse.ErrorResponse(w, r, customError.ErrInvalidRequestBody, http.StatusInternalServerError)
			return
		}

		// Login logic
		user, err := h.UserRepository.GetUserByEmail(r.Context(), req.Email)
		if err != nil || user == nil {
			APIResponse.ErrorResponse(w, r, customError.ErrInvalidEmailOrPassword, http.StatusUnauthorized)
			return
		}

		// Compare the password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			APIResponse.ErrorResponse(w, r, customError.ErrInvalidEmailOrPassword, http.StatusUnauthorized)
			return
		}

		// generate token access and refresh tokens
		// Generate the access token (short-lived)
		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		_, accessTokenString, err := tokenAuth.Encode(map[string]interface{}{
			"sub": strconv.Itoa(user.ID),
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(15 * time.Minute).Unix(), // Access token valid for 15 minutes
		})
		if err != nil {
			// TODO log to DB
			//dbError := customError.NewSystemError(err)
			APIResponse.ErrorResponse(w, r, customError.ErrInternalServerError, http.StatusInternalServerError)
			return
		}

		// Generate the refresh token (long-lived)
		_, refreshTokenString, err := tokenAuth.Encode(map[string]interface{}{
			"sub": strconv.Itoa(user.ID),
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(30 * 24 * time.Hour).Unix(), // Refresh token valid for 30 days
		})
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to generate refresh token: %w", err), http.StatusInternalServerError)
			return
		}

		refreshToken := model.UserRefreshToken{
			Token:     refreshTokenString,
			UserID:    user.ID,
			UserAgent: r.UserAgent(),
			IPAddress: r.RemoteAddr,
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		}

		h.RefreshTokenRepository.Create(r.Context(), &refreshToken)

		userResponseDTO := dto.UserResponse{
			ID:      user.ID,
			Name:    fmt.Sprintf("%s.%s", strings.ToUpper(string(user.FirstName[0])), user.LastName),
			Country: user.Country,
		}

		repsonseDTO := dto.UserLoginResponse{
			User:         userResponseDTO,
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
		}

		APIResponse.SuccessResponse(w, r, APIResponse.APIResponse{
			Success: true,
			Data:    repsonseDTO,
		}, http.StatusOK)
	}
}

func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := dto.UserCreateRequest{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to decode request: %w", err), http.StatusBadRequest)
			return
		}

		// Check for missing or empty fields
		if req.FirstName == "" || req.LastName == "" || req.Email == "" ||
			req.Password == "" || req.Country == "" || req.Birthday == "" {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("missing required fields"), http.StatusBadRequest)
			return
		}

		// Parse the Birthday string into time.Time
		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to parse birthday: %w", err), http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to hash password: %w", err), http.StatusInternalServerError)
			return
		}

		// Create a new user
		user := model.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  string(hashedPassword),
			Country:   req.Country,
			Birthday:  birthday,
		}

		// Save user to database
		newUser, err := h.UserRepository.CreateUser(r.Context(), &user)
		if err != nil {
			APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		responseDTO := dto.UserResponse{
			ID:      newUser.ID,
			Name:    fmt.Sprintf("%s.%s", strings.ToUpper(string(newUser.FirstName[0])), newUser.LastName),
			Country: newUser.Country,
		}

		// Send Email
		err = h.EmailService.Send("welcome_email", "Welcome", []string{newUser.Email}, nil)
		if err != nil {
			util.LogWithContext(
				h.dbLogger,
				slog.LevelError,
				"failed to send email to user",
				map[string]interface{}{
					"userId": newUser.ID,
					"email":  newUser.Email,
				},
				nil)
		}

		APIResponse.SuccessResponse(w, r, &responseDTO, http.StatusCreated)
	}
}

func (h *AuthHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logout logic
		APIResponse.SuccessResponse(w, r, nil)
	}
}

func (h *AuthHandler) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Refresh token logic
		req := dto.UserRefreshTokenRequest{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to decode request: %w", err), http.StatusBadRequest)
			return
		}

		// Get the refresh token from the database
		refreshToken, err := h.RefreshTokenRepository.GetByToken(r.Context(), req.RefreshToken)
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
			"sub": strconv.Itoa(refreshToken.UserID),
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(15 * time.Minute).Unix(), // Access token valid for 15 minutes
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
