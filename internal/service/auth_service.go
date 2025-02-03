package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/database"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/dto"
	customError "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/error"
	APIResponse "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/response"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/model"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/repository"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/service/email"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/util"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db                     *database.Service
	EmailService           *email.EmailService
	UserRepository         repository.UserRepository
	RefreshTokenRepository repository.UserRefreshTokenRepository
}

func NewAuthService(
	db *database.Service,
	emailService *email.EmailService,
	userRepository repository.UserRepository,
	refreshTokenRepository repository.UserRefreshTokenRepository,
) *AuthService {
	return &AuthService{
		db:                     db,
		EmailService:           emailService,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) (dto.UserLoginResponse, error) {

	userLoginResponse := dto.UserLoginResponse{}

	req := dto.UserLoginRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		return userLoginResponse, customError.ErrInvalidRequestBody
	}

	// Login logic
	user, err := s.UserRepository.GetUserByEmail(r.Context(), req.Email)
	if err != nil || user == nil {
		return userLoginResponse, customError.ErrInvalidEmailOrPassword
	}

	// Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return userLoginResponse, customError.ErrInvalidEmailOrPassword
	}

	// generate token access and refresh tokens
	// Generate the access token (short-lived)
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, accessTokenString, err := tokenAuth.Encode(map[string]interface{}{
		"sub": strconv.Itoa(user.ID),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(15 * time.Minute).Unix(), // Access token valid for 15 minutes
	})
	if err != nil {
		// TODO log to DB
		//dbError := customError.NewSystemError(err)
		return userLoginResponse, customError.ErrInternalServerError
	}

	// Generate the refresh token (long-lived)
	_, refreshTokenString, err := tokenAuth.Encode(map[string]interface{}{
		"sub": strconv.Itoa(user.ID),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(), // Refresh token valid for 30 days
	})
	if err != nil {
		return userLoginResponse, customError.ErrUnAuthorized
	}

	refreshToken := model.UserRefreshToken{
		Token:     refreshTokenString,
		UserID:    user.ID,
		UserAgent: r.UserAgent(),
		IPAddress: r.RemoteAddr,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	s.RefreshTokenRepository.Create(r.Context(), &refreshToken)

	userResponseDTO := dto.UserResponse{
		ID:      user.ID,
		Name:    fmt.Sprintf("%s.%s", strings.ToUpper(string(user.FirstName[0])), user.LastName),
		Country: user.Country,
	}

	userLoginResponse.User = userResponseDTO
	userLoginResponse.AccessToken = accessTokenString
	userLoginResponse.RefreshToken = refreshTokenString

	return userLoginResponse, nil
}

func (s *AuthService) Register(w http.ResponseWriter, r *http.Request) (dto.UserResponse, error) {
	userResponseDTO := dto.UserResponse{}

	req := dto.UserCreateRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return userResponseDTO, customError.ErrInvalidRequestBody
	}

	// Check for missing or empty fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" ||
		req.Password == "" || req.Country == "" || req.Birthday == "" {
		return userResponseDTO, fmt.Errorf("missing required fields")
	}

	// Parse the Birthday string into time.Time
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return userResponseDTO, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return userResponseDTO, err
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

	// Save user to database
	newUser, err := s.UserRepository.CreateUser(r.Context(), &user)
	if err != nil {
		return userResponseDTO, err
	}

	// Send Email
	err = s.EmailService.Send("welcome_email", "Welcome", []string{newUser.Email}, nil)
	if err != nil {
		util.LogWithContext(
			s.EmailService.Logger,
			slog.LevelError,
			"failed to send email to user",
			map[string]interface{}{
				"userId": newUser.ID,
				"email":  newUser.Email,
			},
			nil)
	}

	userResponseDTO.ID = newUser.ID
	userResponseDTO.Name = fmt.Sprintf("%s.%s", strings.ToUpper(string(newUser.FirstName[0])), newUser.LastName)
	userResponseDTO.Country = newUser.Country

	return userResponseDTO, nil
}

func (s *AuthService) Logout(w http.ResponseWriter, r *http.Request) error {
	req := dto.UserRefreshTokenRequest{}

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

	// Return a success response
	APIResponse.SuccessResponse(w, r, map[string]string{"message": "Successfully logged out"})

	return nil
}

func (s *AuthService) RefreshToken(w http.ResponseWriter, r *http.Request) (dto.UserRefreshTokenResponse, error) {
	userRefreshTokenResponseDTO := dto.UserRefreshTokenResponse{}

	req := dto.UserRefreshTokenRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return userRefreshTokenResponseDTO, err
	}

	// Get the refresh token from the database
	refreshToken, err := s.RefreshTokenRepository.GetByToken(r.Context(), req.RefreshToken)
	if err != nil {
		return userRefreshTokenResponseDTO, err
	}

	if refreshToken == nil {
		return userRefreshTokenResponseDTO, err
	}

	if refreshToken.ExpiresAt.Unix() < time.Now().Unix() {
		return userRefreshTokenResponseDTO, err
	}

	// Generate the access token (short-lived)
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, accessTokenString, err := tokenAuth.Encode(map[string]interface{}{
		"sub": strconv.Itoa(refreshToken.UserID),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(15 * time.Minute).Unix(), // Access token valid for 15 minutes
	})
	if err != nil {
		return userRefreshTokenResponseDTO, err
	}

	userRefreshTokenResponseDTO.AccessToken = accessTokenString
	userRefreshTokenResponseDTO.RefreshToken = req.RefreshToken

	return userRefreshTokenResponseDTO, nil
}
