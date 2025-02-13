package http

import (
	"net/http"

	// New import

	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterRoutes(c *Container) http.Handler {
	r := chi.NewRouter()
	r.Use(LoggerMiddleware)
	r.Use(CorsMiddleware)
	r.Use(SessionManagerMiddleware(c.SessionManager))
	r.Use(JWTVerifierMiddleware(c.JWTAuth))

	r.Group(func(r chi.Router) {
		r.Post("/api/auth/login", c.AuthHandler.Login())
		r.Post("/api/auth/register", c.AuthHandler.Register())
		r.Post("/api/auth/refresh-token", c.AuthHandler.RefreshToken())
		r.Post("/api/auth/verify-account", c.AuthHandler.VerifyAccount())
		r.Post("/api/auth/login-url", c.SocialAuthHandler.HandleLoginURL())
		r.Get("/api/auth/callback-url", c.SocialAuthHandler.HandleCallbackURL())
	})

	r.Group(func(r chi.Router) {
		r.Use(APIAuthMiddleware)
		r.Get("/api/users", c.UserHandler.GetUsers())
		r.Get("/api/users/{id}", c.UserHandler.GetUser())

		// Profile
		r.Get("/api/profile/me", c.ProfileHandler.GetMyProfile())
		r.Get("/api/profile/{userId}", c.ProfileHandler.GetProfile())
		r.Post("/api/profile/follow", c.ProfileHandler.FollowProfile())
		r.Post("/api/profile/unfollow", c.ProfileHandler.UnfollowProfile())
	})

	// Web
	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	r.Handle("/admin/assets/*", http.StripPrefix("/admin/assets/", fileServer))

	// Public Admin Routes
	r.Group(func(r chi.Router) {
		r.Post("/admin/login", c.CMSHandler.Login)
	})

	// Protected Admin Routes
	r.Group(func(r chi.Router) {
		r.Use(AdminSessionAuthMiddleware(c.SessionManager))
		r.Get("/admin/login", c.CMSHandler.LoginForm)
		r.Get("/admin/dashboard", c.CMSHandler.Dashboard)
		r.Get("/admin/users", c.CMSHandler.Users)
		r.Get("/admin/logout", c.CMSHandler.Logout)
	})

	return r
}
