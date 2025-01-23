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

	r.Get("/", handler.PingHandler())
	r.Get("/health", handler.HealthHandler(s.db))

	r.Group(func(r chi.Router) {
		// Initialize repositories and handlers
		userRepo := repository.NewUserRepository(s.db)
		userHandler := handler.NewUserHandler(userRepo)

		r.Get("/api/users", userHandler.GetUsersHandler())
		r.Get("/api/users/{id}", userHandler.GetUserHandler())
		r.Post("/api/create-user", userHandler.CreateUserHandler())
	})

	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)
	r.Get("/web", templ.Handler(web.HelloForm()).ServeHTTP)
	r.Post("/hello", web.HelloWebHandler)

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
