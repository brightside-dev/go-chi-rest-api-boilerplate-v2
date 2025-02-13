package handler

import (
	"net/http"

	"github.com/brightside-dev/ronin-fitness-be/internal/handler/response"
	"github.com/brightside-dev/ronin-fitness-be/internal/service"
	"github.com/brightside-dev/ronin-fitness-be/internal/service/email"
)

type AuthHandler interface {
	Login() http.HandlerFunc
	Register() http.HandlerFunc
	Logout() http.HandlerFunc
	RefreshToken() http.HandlerFunc
	VerifyAccount() http.HandlerFunc
}

type authHandler struct {
	APIResponse  response.APIResponseManager
	EmailService email.EmailService
	AuthService  service.AuthService
}

func NewAuthHandler(
	apiResponse response.APIResponseManager,
	emailService email.EmailService,
	authService service.AuthService,
) AuthHandler {
	return &authHandler{
		APIResponse:  apiResponse,
		EmailService: emailService,
		AuthService:  authService,
	}
}

func (h *authHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userLoginResponseDTO, err := h.AuthService.Login(w, r)
		if err != nil {
			h.APIResponse.ErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		}

		h.APIResponse.SuccessResponse(w, r, response.APIResponseDTO{
			Success: true,
			Data:    &userLoginResponseDTO,
		}, http.StatusOK)
	}
}

func (h *authHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userResponseDTO, err := h.AuthService.Register(w, r)
		if err != nil {
			h.APIResponse.ErrorResponse(w, r, err, http.StatusBadRequest)
			return
		}

		h.APIResponse.SuccessResponse(w, r, &userResponseDTO, http.StatusCreated)
	}
}

func (h *authHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.AuthService.Logout(w, r)
		if err != nil {
			h.APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}
		h.APIResponse.SuccessResponse(w, r, nil)
	}
}

func (h *authHandler) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRefreshTokenRepsonseDTO, err := h.AuthService.RefreshToken(w, r)
		if err != nil {
			h.APIResponse.ErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		}

		h.APIResponse.SuccessResponse(w, r, userRefreshTokenRepsonseDTO, http.StatusOK)
	}
}

func (h *authHandler) VerifyAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileWithUserDTO, err := h.AuthService.VerifyAccount(w, r)
		if err != nil {
			h.APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}
		h.APIResponse.SuccessResponse(w, r, profileWithUserDTO, http.StatusOK)
	}
}
