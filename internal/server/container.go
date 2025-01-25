package server

import (
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/database"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type Container struct {
	// System
	Router         *chi.Mux
	SessionManager *scs.SessionManager
	JWTAuth        *jwtauth.JWTAuth
	// Handlers
	TestAPIHandler   *handler.TestAPIHandler
	WebHandler       *handler.WebHandler
	AuthAdminHandler *handler.AuthAdminHandler
	AuthHandler      *handler.AuthHandler
	AdminUserHandler *handler.AdminUserHandler
	UserHandler      *handler.UserHandler
}

func NewContainer(db database.Service) *Container {

	// Session Manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db.GetDB())
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// JWT Auth
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	adminUserRepo := repository.NewAdminUserRepository(db)
	refreshTokenRepo := repository.NewUserRefreshTokenRepository(db)

	// Handlers
	userHandler := handler.NewUserHandler(userRepo)
	authHandler := handler.NewAuthHandler(userRepo, refreshTokenRepo)
	testAPIHandler := handler.NewTestAPIHandler()
	authAdminHandler := handler.NewAuthAdminHandler(adminUserRepo, *sessionManager)
	adminUserHandler := handler.NewAdminUserHandler(adminUserRepo)
	webHandler := handler.NewWebHandler(*sessionManager)

	return &Container{
		// System
		Router:         chi.NewRouter(),
		SessionManager: sessionManager,
		JWTAuth:        tokenAuth,
		// Handlers
		TestAPIHandler:   testAPIHandler,
		WebHandler:       webHandler,
		AuthAdminHandler: authAdminHandler,
		AuthHandler:      authHandler,
		AdminUserHandler: adminUserHandler,
		UserHandler:      userHandler,
	}

}
