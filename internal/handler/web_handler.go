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
	template.Render(w, r, "dashboard.html", nil)
}
