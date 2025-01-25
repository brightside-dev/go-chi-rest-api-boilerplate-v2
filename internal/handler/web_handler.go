package handler

import (
	"net/http"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/template"
)

type WebHandler struct {
}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

func (wc *WebHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// data := &template.TemplateData{
	// 	AdminUser: 2021,
	// }
	template.RenderDashboard(w, r, "home", nil)
}

func (wc *WebHandler) Users(w http.ResponseWriter, r *http.Request) {
	template.RenderDashboard(w, r, "users", nil)
}

func (wc *WebHandler) Login(w http.ResponseWriter, r *http.Request) {
	data := &template.TemplateData{
		Form: LoginForm{},
	}

	template.RenderLogin(w, r, "login", data)
}
