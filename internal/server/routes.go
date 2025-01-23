package server

import (
	"net/http"

	"github.com/brightside-dev/boxing-be/cmd/web"
	"github.com/brightside-dev/boxing-be/internal/handler"
	"github.com/brightside-dev/boxing-be/internal/repository"

	"github.com/a-h/templ"
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

	// START API
	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository(s.db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(s.db)
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
		r.Get("/", handler.PingHandler())
		r.Get("/health", handler.HealthHandler(s.db))
	})
	// END API

	// START ADMIN
	// Serve static files
	fileServer := http.FileServer(http.FS(web.Files))
	r.Group(func(r chi.Router) {
		r.Handle("/assets/*", fileServer)
		r.Get("/web", templ.Handler(web.HelloForm()).ServeHTTP)
		r.Post("/hello", web.HelloWebHandler)
	})
	// END ADMIN

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
