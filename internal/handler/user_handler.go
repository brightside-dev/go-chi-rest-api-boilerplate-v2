package handler

import (
	"log/slog"
	"net/http"

	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/service"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/util"
)

type UserHandler struct {
	dbLogger    *slog.Logger
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService, dbLogger *slog.Logger) *UserHandler {
	return &UserHandler{
		dbLogger:    dbLogger,
		UserService: userService,
	}
}

func (h *UserHandler) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usersResponseDTO, err := h.UserService.GetUsers(w, r)
		if err != nil {
			util.LogWithContext(h.dbLogger, slog.LevelError, err.Error(), nil, r)
			APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		APIResponse.SuccessResponse(w, r, usersResponseDTO)
	}
}

func (h *UserHandler) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userResponseDTO, err := h.UserService.GetUser(w, r)
		if err != nil {
			util.LogWithContext(h.dbLogger, slog.LevelError, err.Error(), nil, r)
			APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		APIResponse.SuccessResponse(w, r, &userResponseDTO)
	}
}
