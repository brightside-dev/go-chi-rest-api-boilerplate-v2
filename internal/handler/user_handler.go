package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/dto"
	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/repository"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/util"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserRepository repository.UserRepository
	dbLogger       *slog.Logger
}

func NewUserHandler(repo repository.UserRepository, dbLogger *slog.Logger) *UserHandler {
	return &UserHandler{
		UserRepository: repo,
		dbLogger:       dbLogger,
	}
}

func (h *UserHandler) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.UserRepository.GetAllUsers(r.Context())
		if err != nil {
			util.LogHTTPRequestError(h.dbLogger, "failed to get users", r)

			APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		var response []dto.UserResponse
		for _, user := range users {
			// Ensure FirstName is not empty
			if len(user.FirstName) == 0 {
				continue // Skip users with an empty first name
			}

			// Format the name with the first letter of FirstName capitalized
			formattedName := fmt.Sprintf("%s.%s", strings.ToUpper(string(user.FirstName[0])), user.LastName)

			// Append to response slice
			response = append(response, dto.UserResponse{
				ID:      user.ID,
				Name:    formattedName,
				Country: user.Country,
			})
		}

		APIResponse.SuccessResponse(w, r, response)
	}
}

func (h *UserHandler) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate "id" query parameter
		idParam := chi.URLParam(r, "id")
		if idParam == "" {
			util.LogHTTPRequestError(h.dbLogger, "missing id paramter", r)
			APIResponse.ErrorResponse(w, r, fmt.Errorf("missing 'id' parameter"), http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			util.LogHTTPRequestError(h.dbLogger, "id parameter is not an int", r)
			APIResponse.ErrorResponse(w, r, fmt.Errorf("id must be a valid integer"), http.StatusBadRequest)
			return
		}

		user, err := h.UserRepository.GetUserByID(r.Context(), id)
		if err != nil {
			util.LogHTTPRequestError(h.dbLogger, "failed to get user", r)
			APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		if user == nil {
			util.LogHTTPRequestError(h.dbLogger, "failed to get user", r)
			APIResponse.ErrorResponse(w, r, err, http.StatusNotFound)
			return
		}

		responseDTO := dto.UserResponse{
			ID:      user.ID,
			Name:    fmt.Sprintf("%s.%s", strings.ToUpper(string(user.FirstName[0])), user.LastName),
			Country: user.Country,
		}

		APIResponse.SuccessResponse(w, r, &responseDTO)
	}
}
