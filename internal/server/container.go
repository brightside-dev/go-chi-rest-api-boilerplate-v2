package server

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/database"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/email"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/logger"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	_ "github.com/joho/godotenv/autoload"
)

type Container struct {
	// System
	Env            string
	Router         *chi.Mux
	SessionManager *scs.SessionManager
	JWTAuth        *jwtauth.JWTAuth
	DBLogger       *slog.Logger

	// Handlers
	TestAPIHandler   *handler.TestAPIHandler
	WebHandler       *handler.WebHandler
	AuthAdminHandler *handler.AuthAdminHandler
	AuthHandler      *handler.AuthHandler
	AdminUserHandler *handler.AdminUserHandler
	UserHandler      *handler.UserHandler

	// Services
	EmailService *email.EmailService
}

func NewContainer(db database.Service) *Container {
	// Get the environment
	env := os.Getenv("APP_ENV")

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

	// DB Logger
	logger := NewLogger(db.GetDB())

	// Services
	emailService := email.NewEmailService(logger)
	// pushService, err := push.NewPushService(logger)
	// if err != nil {
	// 	logger.Error("Failed to create push service")
	// }

	// Handlers
	userHandler := handler.NewUserHandler(userRepo, logger)
	authHandler := handler.NewAuthHandler(emailService, userRepo, refreshTokenRepo, logger)
	testAPIHandler := handler.NewTestAPIHandler()
	authAdminHandler := handler.NewAuthAdminHandler(adminUserRepo, *sessionManager)
	adminUserHandler := handler.NewAdminUserHandler(adminUserRepo)
	webHandler := handler.NewWebHandler(*sessionManager)

	return &Container{
		// System
		Env:            env,
		Router:         chi.NewRouter(),
		SessionManager: sessionManager,
		JWTAuth:        tokenAuth,
		DBLogger:       logger,

		// Handlers
		TestAPIHandler:   testAPIHandler,
		WebHandler:       webHandler,
		AuthAdminHandler: authAdminHandler,
		AuthHandler:      authHandler,
		AdminUserHandler: adminUserHandler,
		UserHandler:      userHandler,

		// Services
		EmailService: emailService,
	}

}

func NewLogger(db *sql.DB) *slog.Logger {
	// Create a database log handler.
	dbLogHandler := logger.NewDBLogHandler(db, slog.LevelInfo)

	// Create a TextHandler with AddSource enabled
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})

	// Create a MultiHandler with both handlers (console and DB).
	multiHandler := logger.NewMultiHandler(textHandler, dbLogHandler)

	// Create the logger with the MultiHandler.
	logger := slog.New(multiHandler)

	return logger
}
