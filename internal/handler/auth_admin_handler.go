package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/dto"
	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/model"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthAdminHandler struct {
	AdminUserRepository repository.AdminUserRepository
	SessionManager      scs.SessionManager
}

func NewAuthAdminHandler(
	adminUserRepo repository.AdminUserRepository,
	sessionManager scs.SessionManager,
) *AuthAdminHandler {
	return &AuthAdminHandler{
		AdminUserRepository: adminUserRepo,
		SessionManager:      sessionManager,
	}
}

func (h *AuthAdminHandler) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := dto.AdminUserLoginRequest{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to decode request: %w", err), http.StatusBadRequest)
			return
		}

		// Login logic
		user, err := h.AdminUserRepository.GetByEmail(r.Context(), req.Email)
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("failed to get user: %w", err), http.StatusUnauthorized)
			return
		}

		if user == nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("user not found"), http.StatusUnauthorized)
			return
		}

		// Compare the password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("invalid password and/or email"), http.StatusUnauthorized)
			return
		}

		err = h.SessionManager.RenewToken(r.Context())
		if err != nil {
			//app.serverError(w, r, err)
			return
		}

		// Add the ID of the current user to the session, so that they are now
		// 'logged in'.
		h.SessionManager.Put(r.Context(), "adminUserID", &user.ID)

		// Redirect the user to the create snippet page.
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func (h *AuthAdminHandler) RegisterHandler() http.HandlerFunc {
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

		// Create a new adminUser
		adminUser := model.AdminUser{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  string(hashedPassword),
		}

		// Save adminUser to database
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

func (h *AuthAdminHandler) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logout logic
		APIResponse.SuccessResponse(w, r, nil)
	}
}
