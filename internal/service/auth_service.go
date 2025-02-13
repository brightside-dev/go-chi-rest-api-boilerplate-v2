package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brightside-dev/ronin-fitness-be/database/client"
	"github.com/brightside-dev/ronin-fitness-be/internal/handler/dto"
	customError "github.com/brightside-dev/ronin-fitness-be/internal/handler/error"
	"github.com/brightside-dev/ronin-fitness-be/internal/model"
	"github.com/brightside-dev/ronin-fitness-be/internal/repository"
	"github.com/brightside-dev/ronin-fitness-be/internal/service/email"
	"github.com/brightside-dev/ronin-fitness-be/internal/util"
	"github.com/go-playground/validator/v10"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/joho/godotenv/autoload"
)

type AuthService interface {
	Login(w http.ResponseWriter, r *http.Request) (*dto.AuthLoginResponse, error)
	Register(w http.ResponseWriter, r *http.Request) (*dto.AuthRegisterResponse, error)
	Logout(w http.ResponseWriter, r *http.Request) error
	RefreshToken(w http.ResponseWriter, r *http.Request) (*dto.RefreshTokenResponse, error)
	VerifyAccount(w http.ResponseWriter, r *http.Request) (*dto.UserResponse, error)
}

type authService struct {
	DB                         client.DatabaseService
	Logger                     *slog.Logger
	Validate                   *validator.Validate
	TokenAuth                  *jwtauth.JWTAuth
	EmailService               email.EmailService
	UserRepository             repository.UserRepository
	RefreshTokenRepository     repository.RefreshTokenRepository
	ProfileReposistory         repository.ProfileRepository
	VerificationCodeRepository repository.VerificationCodeRepository
}

func NewAuthService(
	db client.DatabaseService,
	logger *slog.Logger,
	validate *validator.Validate,
	tokenAuth *jwtauth.JWTAuth,
	emailService email.EmailService,
	userRepository repository.UserRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	profileRepository repository.ProfileRepository,
	verificationCodeRepository repository.VerificationCodeRepository,
) AuthService {
	return &authService{
		DB:                         db,
		Logger:                     logger,
		TokenAuth:                  tokenAuth,
		Validate:                   validate,
		EmailService:               emailService,
		UserRepository:             userRepository,
		RefreshTokenRepository:     refreshTokenRepository,
		ProfileReposistory:         profileRepository,
		VerificationCodeRepository: verificationCodeRepository,
	}
}

func (s *authService) Login(w http.ResponseWriter, r *http.Request) (*dto.AuthLoginResponse, error) {
	req := dto.AuthLoginRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInvalidRequestBody
	}

	// Login logic
	user, err := s.UserRepository.GetByEmail(r.Context(), req.Email)
	if err != nil {
		util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInvalidEmailOrPassword
	}

	if user == nil {
		return nil, customError.ErrInvalidEmailOrPassword
	}

	// Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, customError.ErrInvalidEmailOrPassword
	}

	result, err := util.WithTransaction(r.Context(), s.DB.GetDB(), func(tx *sql.Tx) (interface{}, error) {
		// Create access and refresh tokens
		refreshTokenResponseDTO, err := s.createTokens(r, tx, user.ID)
		if err != nil {
			return nil, err
		}

		return refreshTokenResponseDTO, nil
	})
	if err != nil {
		util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInternalServerError
	}

	refreshTokenResponseDTO, ok := result.(*dto.RefreshTokenResponse)
	if !ok {
		util.LogWithContext(s.Logger, slog.LevelError, "unexpected result type", nil, r)
		return nil, customError.ErrInternalServerError
	}

	authLoginResponseDTO := dto.AuthLoginResponse{
		UserResponse: dto.UserResponse{
			UserID:     user.ID,
			Name:       fmt.Sprintf("%s.%s", strings.ToUpper(string(user.FirstName[0])), user.LastName),
			Country:    user.Country,
			IsVerified: user.IsVerified,
		},
		RefreshTokenResponse: *refreshTokenResponseDTO,
	}

	return &authLoginResponseDTO, nil
}

func (s *authService) Register(w http.ResponseWriter, r *http.Request) (*dto.AuthRegisterResponse, error) {
	req := dto.AuthRegisterRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInvalidRequestBody
	}

	// validate the request
	err = s.Validate.Struct(req)
	if err != nil {
		return nil, util.FormatValidationError(err.(validator.ValidationErrors))
	}

	// Parse the Birthday string into time.Time
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
		return nil, err
	}

	// Create a new user
	user := model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Country:   req.Country,
		Birthday:  birthday,
	}

	result, err := util.WithTransaction(r.Context(), s.DB.GetDB(), func(tx *sql.Tx) (interface{}, error) {
		newUser, err := s.UserRepository.Create(r.Context(), tx, &user)
		if err != nil {
			return nil, err
		}

		// Create a verification code
		verificationCode := model.VerificationCode{
			UserID: newUser.ID,
			Email:  newUser.Email,
			Code:   util.GenerateVerificationCode(),
		}

		newVerificationCode, err := s.VerificationCodeRepository.Create(r.Context(), tx, &verificationCode)
		if err != nil {
			return nil, err
		}

		// Create access and refresh tokens
		refreshTokenResponseDTO, err := s.createTokens(r, tx, newUser.ID)
		if err != nil {
			return nil, err
		}

		//Send Email
		data := map[string]string{
			"name":   newUser.FirstName,
			"code":   newVerificationCode.Code,
			"expiry": newVerificationCode.ExpiresAt.UTC().Format(time.RFC3339),
		}

		err = s.EmailService.Send("verify_email", "Account Verification - Ronin Fitness", []string{newUser.Email}, data)
		if err != nil {
			util.LogWithContext(
				s.Logger,
				slog.LevelError,
				"failed to send email to user",
				map[string]interface{}{
					"userId": newUser.ID,
					"email":  newUser.Email,
				},
				nil)
		}

		authRegisterResponseDTO := dto.AuthRegisterResponse{
			UserID:               newUser.ID,
			RefreshTokenResponse: *refreshTokenResponseDTO,
		}

		return &authRegisterResponseDTO, nil
	})
	if err != nil {
		return nil, err
	}

	authRegisterResponseDTO, ok := result.(*dto.AuthRegisterResponse)
	if !ok {
		util.LogWithContext(s.Logger, slog.LevelError, "unexpected result type", nil, r)
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return authRegisterResponseDTO, nil
}

func (s *authService) Logout(w http.ResponseWriter, r *http.Request) error {
	req := dto.RefreshTokenRequest{}

	// Decode request body to get the refresh token
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	// Validate refresh token
	refreshToken, err := s.RefreshTokenRepository.GetByToken(r.Context(), req.RefreshToken)
	if err != nil || refreshToken == nil {
		return customError.ErrUnAuthorized
	}

	// Delete the refresh token from the database
	err = s.RefreshTokenRepository.DeleteByToken(r.Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) RefreshToken(w http.ResponseWriter, r *http.Request) (*dto.RefreshTokenResponse, error) {
	req := dto.RefreshTokenRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	// Get the refresh token from the database
	refreshToken, err := s.RefreshTokenRepository.GetByToken(r.Context(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if refreshToken == nil {
		return nil, err
	}

	if refreshToken.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, err
	}

	result, err := util.WithTransaction(r.Context(), s.DB.GetDB(), func(tx *sql.Tx) (interface{}, error) {
		// Create access and refresh tokens
		refreshTokenResponseDTO, err := s.createTokens(r, tx, refreshToken.UserID)
		if err != nil {
			return nil, err
		}

		return refreshTokenResponseDTO, nil
	})
	if err != nil {
		return nil, err
	}

	refreshTokenResponseDTO, ok := result.(*dto.RefreshTokenResponse)
	if !ok {
		util.LogWithContext(s.Logger, slog.LevelError, "unexpected result type", nil, r)
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return refreshTokenResponseDTO, nil
}

func (s *authService) VerifyAccount(w http.ResponseWriter, r *http.Request) (*dto.UserResponse, error) {
	req := dto.AuthVerifyRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	// Get the verification code from the database
	verificationCode, err := s.VerificationCodeRepository.GetByCode(r.Context(), req.Code, req.UserID)
	if err != nil {
		util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
		return nil, err
	}

	if verificationCode == nil {
		return nil, fmt.Errorf("verification code not found")
	}

	if verificationCode.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, fmt.Errorf("verification code expired")
	}

	user, err := s.UserRepository.GetByID(r.Context(), req.UserID)
	if err != nil {
		util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	user.IsVerified = 1

	result, err := util.WithTransaction(r.Context(), s.DB.GetDB(), func(tx *sql.Tx) (interface{}, error) {
		// Update the user's account to be verified
		err := s.UserRepository.Update(r.Context(), tx, user)
		if err != nil {
			util.LogWithContext(s.Logger, slog.LevelError, "test", nil, r)
			return nil, err
		}

		data := map[string]string{
			"name": user.FirstName,
		}

		err = s.EmailService.Send("welcome_email", "Welcome to Ronin Fitness!", []string{user.Email}, data)
		if err != nil {
			util.LogWithContext(
				s.Logger,
				slog.LevelError,
				"failed to send email to user",
				map[string]interface{}{
					"userId": user.ID,
					"email":  user.Email,
				},
				nil)
		}

		return &dto.UserResponse{
			UserID:     user.ID,
			Name:       fmt.Sprintf("%s.%s", strings.ToUpper(string(user.FirstName[0])), user.LastName),
			Country:    user.Country,
			IsVerified: user.IsVerified,
		}, nil

	})
	if err != nil {
		return nil, err
	}

	userResponseDTO, ok := result.(*dto.UserResponse)
	if !ok {
		util.LogWithContext(s.Logger, slog.LevelError, "unexpected result type", nil, r)
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return userResponseDTO, nil
}

func (s *authService) createTokens(r *http.Request, tx *sql.Tx, userID int) (*dto.RefreshTokenResponse, error) {
	// generate token access and refresh tokens
	// Generate the access token (short-lived)
	_, accessTokenString, err := s.TokenAuth.Encode(map[string]interface{}{
		"sub": strconv.Itoa(userID),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(15 * time.Minute).Unix(), // Access token valid for 15 minutes
	})
	if err != nil {

		return nil, customError.ErrInternalServerError
	}

	// Generate the refresh token (long-lived)
	_, refreshTokenString, err := s.TokenAuth.Encode(map[string]interface{}{
		"sub": strconv.Itoa(userID),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(), // Refresh token valid for 30 days
	})
	if err != nil {
		util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrUnAuthorized
	}

	refreshToken := model.RefreshToken{
		Token:     refreshTokenString,
		UserID:    userID,
		UserAgent: r.UserAgent(),
		IPAddress: r.RemoteAddr,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	err = s.RefreshTokenRepository.Create(r.Context(), tx, &refreshToken)
	if err != nil {
		util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInternalServerError
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
