package handler

import (
	"log/slog"
	"net/http"

	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/service"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/service/email"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/util"
)

type AuthHandler struct {
	EmailService *email.EmailService
	AuthService  service.AuthServiceInterface
}

func NewAuthHandler(
	emailService *email.EmailService,
	authService service.AuthServiceInterface,
) *AuthHandler {
	return &AuthHandler{
		EmailService: emailService,
		AuthService:  authService,
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userLoginResponseDTO, err := h.AuthService.Login(w, r)
		if err != nil {
			util.LogWithContext(h.EmailService.Logger, slog.LevelError, err.Error(), nil, r)
			APIResponse.ErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		}

		APIResponse.SuccessResponse(w, r, APIResponse.APIResponse{
			Success: true,
			Data:    &userLoginResponseDTO,
		}, http.StatusOK)
	}
}

func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userResponseDTO, err := h.AuthService.Register(w, r)
		if err != nil {
			util.LogWithContext(h.EmailService.Logger, slog.LevelError, err.Error(), nil, r)
			APIResponse.ErrorResponse(w, r, err, http.StatusBadRequest)
			return
		}

		APIResponse.SuccessResponse(w, r, &userResponseDTO, http.StatusCreated)
	}
}

func (h *AuthHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.AuthService.Logout(w, r)
		if err != nil {
			util.LogWithContext(h.EmailService.Logger, slog.LevelError, err.Error(), nil, r)
			APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}
		APIResponse.SuccessResponse(w, r, nil)
	}
}

func (h *AuthHandler) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRefreshTokenRepsonseDTO, err := h.AuthService.RefreshToken(w, r)
		if err != nil {
			util.LogWithContext(h.EmailService.Logger, slog.LevelError, err.Error(), nil, r)
			APIResponse.ErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		}

		APIResponse.SuccessResponse(w, r, userRefreshTokenRepsonseDTO, http.StatusOK)
	}
}
