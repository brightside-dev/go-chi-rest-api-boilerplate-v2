package handler

import (
	"net/http"
)

func LandingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Landing"))
}
