package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brightside-dev/boxing-be/internal/handler/dto"
	APIResponse "github.com/brightside-dev/boxing-be/internal/handler/response"
	"github.com/brightside-dev/boxing-be/internal/model"
	"github.com/brightside-dev/boxing-be/internal/repository"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserRepository repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepository: repo}
}

func (h *UserHandler) GetUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.UserRepository.GetAllUsers(r.Context())
		if err != nil {
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

func (h *UserHandler) GetUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			APIResponse.ErrorResponse(w, r, err, http.StatusBadRequest)
			return
		}

		user, err := h.UserRepository.GetUserByID(r.Context(), id)
		if err != nil {
			APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		if user == nil {
			err = fmt.Errorf("user with ID %d not found", id)
			APIResponse.ErrorResponse(w, r, err, http.StatusNotFound)
			return
		}

		responseDTO := dto.UserResponse{
			ID:      user.ID,
			Name:    fmt.Sprintf("%s.%s", strings.ToUpper(string(user.FirstName[0])), user.LastName),
			Country: user.Country,
		}

		APIResponse.SuccessResponse(w, r, responseDTO)
	}
}

func (h *UserHandler) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := dto.CreateUserRequest{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Parse the Birthday string into time.Time
		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Create a new user
		user := model.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  req.Password,
			Country:   req.Country,
			Birthday:  birthday,
		}

		// Save user to database
		_, err = h.UserRepository.CreateUser(r.Context(), &user)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
