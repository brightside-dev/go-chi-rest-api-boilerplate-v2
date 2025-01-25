package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/dto"
	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/model"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/repository"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/template"
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

func (h *AuthAdminHandler) RedirectToLoginForm(form *LoginForm, w http.ResponseWriter, r *http.Request) {
	data := &template.TemplateData{}
	data.Form = form
	template.RenderLogin(w, r, "login", data)
}

type LoginForm struct {
	Email        string
	Password     string
	FormErrors   map[string]string
	SystemErrors map[string]string
}

func (h *AuthAdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{
		Email:        r.FormValue("email"),
		Password:     r.FormValue("password"),
		FormErrors:   map[string]string{},
		SystemErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Email) == "" {
		form.FormErrors["email"] = "Email cannot be blank"
		h.RedirectToLoginForm(&form, w, r)
		return
	}

	if strings.TrimSpace(form.Password) == "" {
		form.FormErrors["password"] = "Password cannot be blank"
		h.RedirectToLoginForm(&form, w, r)
		return
	}

	// Login logic
	user, err := h.AdminUserRepository.GetByEmail(r.Context(), form.Email)
	if err != nil {
		form.SystemErrors["message"] = err.Error()
		h.RedirectToLoginForm(&form, w, r)
		return
	}

	if user == nil {
		form.SystemErrors["message"] = "Invalid password and/or email"
		h.RedirectToLoginForm(&form, w, r)
		return
	}

	// Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		form.SystemErrors["message"] = "Invalid password and/or email"
		h.RedirectToLoginForm(&form, w, r)
		return
	}

	err = h.SessionManager.RenewToken(r.Context())
	if err != nil {
		form.SystemErrors["message"] = err.Error()
		h.RedirectToLoginForm(&form, w, r)
		return
	}

	// Add the ID of the current user to the session, so that they are now
	// 'logged in'.
	h.SessionManager.Put(r.Context(), "adminUserID", &user.ID)

	// flash success
	h.SessionManager.Put(r.Context(), "flash", "You've been logged in successfully")

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func (h *AuthAdminHandler) Register(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthAdminHandler) Logout(w http.ResponseWriter, r *http.Request) {

	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	h.SessionManager.Remove(r.Context(), "adminUserID")

	// flash success
	h.SessionManager.Put(r.Context(), "flash", "You've been logged out successfully")

	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}
