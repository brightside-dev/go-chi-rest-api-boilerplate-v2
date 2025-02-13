package handler

import (
	"log/slog"
	"net/http"

	"github.com/brightside-dev/ronin-fitness-be/internal/handler/response"
	"github.com/brightside-dev/ronin-fitness-be/internal/service"
	"github.com/brightside-dev/ronin-fitness-be/internal/util"
)

type UserHandler interface {
	GetUsers() http.HandlerFunc
	GetUser() http.HandlerFunc
}

type userHandler struct {
	APIResponse response.APIResponseManager
	DBLogger    *slog.Logger
	UserService service.UserService
}

func NewUserHandler(
	userService service.UserService,
	dbLogger *slog.Logger,
) UserHandler {
	return &userHandler{
		DBLogger:    dbLogger,
		UserService: userService,
	}
}

func (h *userHandler) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usersResponseDTO, err := h.UserService.GetUsers(w, r)
		if err != nil {
			util.LogWithContext(h.DBLogger, slog.LevelError, err.Error(), nil, r)
			h.APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		h.APIResponse.SuccessResponse(w, r, usersResponseDTO)
	}
}

func (h *userHandler) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userResponseDTO, err := h.UserService.GetUser(w, r)
		if err != nil {
			util.LogWithContext(h.DBLogger, slog.LevelError, err.Error(), nil, r)
			h.APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		h.APIResponse.SuccessResponse(w, r, &userResponseDTO)
	}
}
