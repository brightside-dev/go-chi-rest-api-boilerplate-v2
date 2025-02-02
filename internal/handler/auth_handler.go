package handler

import (
	"net/http"

	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/service"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/service/email"
)

type AuthHandler struct {
	EmailService *email.EmailService
	AuthService  *service.AuthService
}

func NewAuthHandler(
	emailService *email.EmailService,
	authService *service.AuthService,
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
			APIResponse.ErrorResponse(w, r, err, http.StatusBadRequest)
			return
		}

		APIResponse.SuccessResponse(w, r, &userResponseDTO, http.StatusCreated)
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
		userRefreshTokenRepsonseDTO, err := h.AuthService.RefreshToken(w, r)
		if err != nil {
			APIResponse.ErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		}

		APIResponse.SuccessResponse(w, r, userRefreshTokenRepsonseDTO, http.StatusOK)
	}
}
