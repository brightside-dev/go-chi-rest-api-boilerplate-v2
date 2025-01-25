package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/dto"
	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/model"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/repository"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type AdminUserHandler struct {
	AdminUserRepository repository.AdminUserRepository
}

func NewAdminUserHandler(repo repository.AdminUserRepository) *AdminUserHandler {
	return &AdminUserHandler{AdminUserRepository: repo}
}

func (h *AdminUserHandler) GetUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.AdminUserRepository.GetAll(r.Context())
		if err != nil {
			APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		var response []dto.AdminUserResponse
		for _, user := range users {
			// Append to response slice
			response = append(response, dto.AdminUserResponse{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			})
		}

		APIResponse.SuccessResponse(w, r, response)
	}
}

func (h *AdminUserHandler) GetUserHandler() http.HandlerFunc {
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

		user, err := h.AdminUserRepository.GetByID(r.Context(), id)
		if err != nil {
			APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		if user == nil {
			err = fmt.Errorf("user with ID %d not found", id)
			APIResponse.ErrorResponse(w, r, err, http.StatusNotFound)
			return
		}

		responseDTO := dto.AdminUserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}

		APIResponse.SuccessResponse(w, r, &responseDTO)
	}
}

func (h *AdminUserHandler) CreateUserHandler() http.HandlerFunc {
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

		// Create a new user
		adminUser := model.AdminUser{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  string(hashedPassword),
		}

		// Save user to database
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
