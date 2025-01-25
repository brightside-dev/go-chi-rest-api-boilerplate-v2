package server

import (
	"net/http"

	// New import

	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterRoutes(container *Container) http.Handler {
	r := chi.NewRouter()
	r.Use(LoggerMiddleware)
	r.Use(CorsMiddleware)
	r.Use(SessionManagerMiddleware(container.SessionManager))
	r.Use(JWTVerifierMiddleware(container.JWTAuth))

	r.Get("/api/health", container.TestAPIHandler.HealthHandler(s.db))
	r.Get("/api/ping", container.TestAPIHandler.PingHandler())

	// START API
	r.Group(func(r chi.Router) {
		// Auth Middleware
		r.Use(APIAuthMiddleware)
		r.Get("/api/users", container.UserHandler.GetUsersHandler())
		r.Get("/api/users/{id}", container.UserHandler.GetUserHandler())
		r.Post("/api/create-user", container.UserHandler.CreateUserHandler())
	})

	r.Group(func(r chi.Router) {
		r.Post("/api/auth/login", container.AuthHandler.LoginHandler())
		r.Post("/api/auth/register", container.AuthHandler.RegisterHandler())
		r.Post("/api/auth/refresh-token", container.AuthHandler.RefreshTokenHandler())
	})
	// END API

	// Web
	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	r.Handle("/admin/assets/*", http.StripPrefix("/admin/assets/", fileServer))

	// TODO for testing - remove later
	r.Group(func(r chi.Router) {
		// Auth Middleware
		r.Get("/api/admin/users", container.AdminUserHandler.GetUsersHandler())
		r.Get("/api/admin/users/{id}", container.AdminUserHandler.GetUserHandler())
		r.Post("/api/admin/create-admin-user", container.AdminUserHandler.CreateUserHandler())
	})

	// Public Admin Routes
	r.Group(func(r chi.Router) {
		r.Post("/admin/login", container.AuthAdminHandler.LoginPostHandler)
		r.Post("/admin/register", container.AuthAdminHandler.RegisterPostHandler)
	})

	// Protected Admin Routes
	r.Group(func(r chi.Router) {
		r.Use(AdminSessionAuthMiddleware(container.SessionManager))
		r.Get("/admin/login", container.WebHandler.LoginFormHandler)
		r.Get("/admin/dashboard", container.WebHandler.Dashboard)
		r.Get("/admin/users", container.WebHandler.Users)
		r.Get("/admin/logout", container.AuthAdminHandler.Logout)
	})

	return r
}
