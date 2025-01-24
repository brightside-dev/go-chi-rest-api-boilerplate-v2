package server

import (
	"net/http"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", handler.HealthHandler(s.db))
	r.Get("/api/ping", handler.PingHandler())

	// START API
	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository(s.db)
	refreshTokenRepo := repository.NewUserRefreshTokenRepository(s.db)
	userHandler := handler.NewUserHandler(userRepo)
	authHandler := handler.NewAuthHandler(userRepo, refreshTokenRepo)
	r.Group(func(r chi.Router) {
		// Auth Middleware
		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(APIAuthMiddleware)
		r.Get("/api/users", userHandler.GetUsersHandler())
		r.Get("/api/users/{id}", userHandler.GetUserHandler())
		r.Post("/api/create-user", userHandler.CreateUserHandler())
	})

	r.Group(func(r chi.Router) {
		r.Post("/api/auth/login", authHandler.LoginHandler())
		r.Post("/api/auth/register", authHandler.RegisterHandler())
		r.Post("/api/auth/refresh-token", authHandler.RefreshTokenHandler())
	})
	// END API

	adminUserRepo := repository.NewAdminUserRepository(s.db)
	adminRefreshTokenRepo := repository.NewAdminUserRefreshTokenRepository(s.db)
	authAdminHandler := handler.NewAuthAdminHandler(adminUserRepo, adminRefreshTokenRepo)
	adminUserHandler := handler.NewAdminUserHandler(adminUserRepo)
	r.Group(func(r chi.Router) {
		// Auth Middleware
		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(APIAdminAuthMiddleware)
		r.Get("/api/admin/users", adminUserHandler.GetUsersHandler())
		r.Get("/api/admin/users/{id}", adminUserHandler.GetUserHandler())
		r.Post("/api/admin/create-admin-user", adminUserHandler.CreateUserHandler())
	})

	r.Group(func(r chi.Router) {
		r.Post("/api/admin/auth/login", authAdminHandler.LoginHandler())
		r.Post("/api/admin/auth/register", authAdminHandler.RegisterHandler())
		r.Post("/api/admin/auth/refresh-token", authAdminHandler.RefreshTokenHandler())
	})

	// Web
	// Serve the static files
	// fileServer := http.FileServer(http.Dir(".ui/static/"))
	// r.Get("/static", http.StripPrefix("/static/", fileServer).ServeHTTP)

	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))

	webHandler := handler.NewWebHandler()
	r.Group(func(r chi.Router) {
		r.Get("/dashboard", webHandler.Dashboard)
		r.Get("/users", webHandler.Users)
		r.Get("/login", webHandler.Login)
	})

	return r
}

// func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
// 	resp := make(map[string]string)
// 	resp["message"] = "Hello World"

// 	jsonResp, err := json.Marshal(resp)
// 	if err != nil {
// 		log.Fatalf("error handling JSON marshal. Err: %v", err)
// 	}

// 	_, _ = w.Write(jsonResp)
// }
