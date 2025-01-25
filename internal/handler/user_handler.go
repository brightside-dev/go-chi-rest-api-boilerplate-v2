package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/dto"
	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/model"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/repository"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserRepository repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepository: repo}
}

func (h *UserHandler) GetUsers() http.HandlerFunc {
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

func (h *UserHandler) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate "id" query parameter
		idParam := chi.URLParam(r, "id")
		if idParam == "" {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("missing 'id' parameter"), http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("id must be a valid integer"), http.StatusBadRequest)
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

		APIResponse.SuccessResponse(w, r, &responseDTO)
	}
}

func (h *UserHandler) CreateUser() http.HandlerFunc {
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

		APIResponse.SuccessResponse(w, r, &responseDTO, http.StatusCreated)
	}
}
