package handler

import (
	"net/http"
	"strings"

	"github.com/brightside-dev/ronin-fitness-be/internal/service"
	"github.com/brightside-dev/ronin-fitness-be/internal/template"

	"github.com/alexedwards/scs/v2"
)

type CMSHandler interface {
	Dashboard(w http.ResponseWriter, r *http.Request)
	Users(w http.ResponseWriter, r *http.Request)
	LoginForm(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type cmsHandler struct {
	SessionManager   scs.SessionManager
	AdminUserService service.AdminUserService
}

func NewCMSHandler(
	sessionManager scs.SessionManager,
	adminUserService service.AdminUserService,
) CMSHandler {
	return &cmsHandler{
		SessionManager:   sessionManager,
		AdminUserService: adminUserService,
	}
}

func (h *cmsHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	data := template.NewTemplateData(r, &h.SessionManager)

	template.RenderDashboard(w, r, "home", data)
}

func (h *cmsHandler) Users(w http.ResponseWriter, r *http.Request) {
	template.RenderDashboard(w, r, "users", nil)
}

type LoginForm struct {
	Email        string
	Password     string
	FormErrors   map[string]string
	SystemErrors map[string]string
}

func (h *cmsHandler) LoginForm(w http.ResponseWriter, r *http.Request) {
	data := template.NewTemplateData(r, &h.SessionManager)
	data.Form = LoginForm{}

	template.RenderLogin(w, r, "login", data)
}

func (h *cmsHandler) RedirectToLoginForm(form *LoginForm, w http.ResponseWriter, r *http.Request) {
	data := &template.TemplateData{}
	data.Form = form
	template.RenderLogin(w, r, "login", data)
}

func (h *cmsHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	adminUser, err := h.AdminUserService.Login(form.Email, form.Password)
	if err != nil {
		form.SystemErrors["message"] = err.Error()
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
	h.SessionManager.Put(r.Context(), "adminUserID", &adminUser.ID)

	// flash success
	h.SessionManager.Put(r.Context(), "flash", "You've been logged in successfully")

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func (h *cmsHandler) Logout(w http.ResponseWriter, r *http.Request) {

	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	h.SessionManager.Remove(r.Context(), "adminUserID")

	// flash success
	h.SessionManager.Put(r.Context(), "flash", "You've been logged out successfully")

	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}
