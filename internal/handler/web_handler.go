package handler

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/template"
)

type WebHandler struct {
	SessionManager scs.SessionManager
}

func NewWebHandler(sessionManager scs.SessionManager) *WebHandler {
	return &WebHandler{
		SessionManager: sessionManager,
	}
}

func (wc *WebHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	data := template.NewTemplateData(r, &wc.SessionManager)

	template.RenderDashboard(w, r, "home", data)
}

func (wc *WebHandler) Users(w http.ResponseWriter, r *http.Request) {
	template.RenderDashboard(w, r, "users", nil)
}
