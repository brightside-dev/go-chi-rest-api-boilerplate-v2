package server

import (
	"fmt"
	"net/http"
	"time"

	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/go-chi/jwtauth/v5"
)

func APIAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if token exists
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil || claims == nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("invalid JWT Token"), http.StatusUnauthorized)
			return
		}

		// Verify the token's "exp" exists and is a string
		expClaim, ok := claims["exp"].(time.Time)
		if !ok {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("invalid JWT Token"), http.StatusUnauthorized)
			return
		}

		// Check if the token has expired
		if time.Now().After(expClaim) {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("JWT Token has expired"), http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func APIAdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if token exists
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil || claims == nil {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("invalid JWT Token"), http.StatusUnauthorized)
			return
		}

		// Verify the token's "exp" exists and is a string
		expClaim, ok := claims["exp"].(time.Time)
		if !ok {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("invalid JWT Token"), http.StatusUnauthorized)
			return
		}

		// Check if the token has expired
		if time.Now().After(expClaim) {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("JWT Token has expired"), http.StatusUnauthorized)
			return
		}

		// Check if the user is an admin
		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			APIResponse.ErrorResponse(w, r, fmt.Errorf("unauthorized access"), http.StatusForbidden)
			return
		}

		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
