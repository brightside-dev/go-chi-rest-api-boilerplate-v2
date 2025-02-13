package http

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/brightside-dev/ronin-fitness-be/database/client"
	"github.com/brightside-dev/ronin-fitness-be/internal/handler"
	"github.com/brightside-dev/ronin-fitness-be/internal/handler/response"
	"github.com/brightside-dev/ronin-fitness-be/internal/repository"
	"github.com/brightside-dev/ronin-fitness-be/internal/service"
	"github.com/brightside-dev/ronin-fitness-be/internal/service/email"
	"github.com/brightside-dev/ronin-fitness-be/internal/service/logger"
	"github.com/brightside-dev/ronin-fitness-be/internal/service/oauth"
	"github.com/go-playground/validator/v10"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
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

	APIResponseManager response.APIResponseManager

	// Handlers
	// TestAPIHandler   *handler.TestAPIHandler
	CMSHandler handler.CMSHandler
	// AuthAdminHandler *handler.AuthAdminHandler
	AuthHandler       handler.AuthHandler
	SocialAuthHandler handler.SocialAuthHandler
	// AdminUserHandler *handler.AdminUserHandler
	UserHandler    handler.UserHandler
	ProfileHandler handler.ProfileHandler

	// Services
	EmailService     email.EmailService
	AdminUserService service.AdminUserService
}

func NewContainer(db client.DatabaseService) *Container {
	// Get the environment
	env := os.Getenv("APP_ENV")
	// Chi Router
	router := chi.NewRouter()
	// Session Manager
	sessionManager := NewSessionManager(db.GetDB())
	// HTTP Response Manager
	apiResponseManager := response.NewAPIResponseManager()
	// JWT Auth
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	// Repositories
	userRepository := repository.NewUserRepository(db)
	adminUserRepository := repository.NewAdminUserRepository(db)
	refreshTokenRepository := repository.NewUserRefreshTokenRepository(db)
	profileRepository := repository.NewProfileRepository(db)
	profileFollowsRepository := repository.NewProfileFollowRepository(db)
	verificationCodeRepository := repository.NewVerificationCodeRepository(db)

	// DB Logger
	logger := NewLogger(db.GetDB())

	// Request Validator
	validator := validator.New()

	// Services
	emailService := email.NewEmailService(logger)
	oAuthService := oauth.NewOAuthService(logger)
	authService := service.NewAuthService(db, logger, validator, tokenAuth, emailService, userRepository, refreshTokenRepository, profileRepository, verificationCodeRepository)
	userService := service.NewUserService(db, logger, userRepository)
	adminUserService := service.NewAdminUserService(adminUserRepository)
	profileService := service.NewProfileService(db, logger, validator, profileRepository, profileFollowsRepository, userRepository)
	// // pushService, err := push.NewPushService(logger)
	// // if err != nil {
	// // 	logger.Error("Failed to create push service")
	// // }

	// // Handlers
	userHandler := handler.NewUserHandler(userService, logger)
	authHandler := handler.NewAuthHandler(apiResponseManager, emailService, authService)
	socialAuthHandler := handler.NewSocialAuthHandler(apiResponseManager, oAuthService)
	profileHandler := handler.NewProfileHandler(apiResponseManager, logger, profileService)
	// testAPIHandler := handler.NewTestAPIHandler()
	// authAdminHandler := handler.NewAuthAdminHandler(adminUserRepo, *sessionManager)
	// adminUserHandler := handler.NewAdminUserHandler(adminUserRepo)
	CMSHandler := handler.NewCMSHandler(*sessionManager, adminUserService)

	return &Container{
		// System
		Env:                env,
		Router:             router,
		SessionManager:     sessionManager,
		JWTAuth:            tokenAuth,
		DBLogger:           logger,
		APIResponseManager: apiResponseManager,

		// // Handlers
		// TestAPIHandler:   testAPIHandler,
		CMSHandler:        CMSHandler,
		AuthHandler:       authHandler,
		SocialAuthHandler: socialAuthHandler,
		// AdminUserHandler: adminUserHandler,
		UserHandler:    userHandler,
		ProfileHandler: profileHandler,

		// Services
		EmailService:     emailService,
		AdminUserService: adminUserService,
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

func NewSessionManager(db *sql.DB) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return sessionManager
}
