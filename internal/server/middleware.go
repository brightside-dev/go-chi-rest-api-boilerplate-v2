package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return middleware.Logger(next)
}

func CorsMiddleware(next http.Handler) http.Handler {
	handler := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return handler(next)
}

func SessionManagerMiddleware(sessionManager *scs.SessionManager) func(http.Handler) http.Handler {
	return sessionManager.LoadAndSave
}

func JWTVerifierMiddleware(jwtAuth *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return jwtauth.Verifier(jwtAuth)
}

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
