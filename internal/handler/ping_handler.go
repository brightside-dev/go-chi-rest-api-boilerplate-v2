package handler

import (
	"encoding/json"
	"net/http"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/database"
)

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	}
}

func HealthHandler(db database.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonResp, _ := json.Marshal(db.Health())
		_, _ = w.Write(jsonResp)
	}
}
