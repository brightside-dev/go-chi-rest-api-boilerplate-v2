package service

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/brightside-dev/ronin-fitness-be/database/client"
	"github.com/brightside-dev/ronin-fitness-be/internal/handler/dto"
	customError "github.com/brightside-dev/ronin-fitness-be/internal/handler/error"
	"github.com/brightside-dev/ronin-fitness-be/internal/repository"
	"github.com/brightside-dev/ronin-fitness-be/internal/util"

	"github.com/go-chi/chi/v5"
)

type UserService interface {
	GetUsers(w http.ResponseWriter, r *http.Request) ([]dto.UserResponse, error)
	GetUser(w http.ResponseWriter, r *http.Request) (dto.UserResponse, error)
}

type userService struct {
	DB             client.DatabaseService
	DBLogger       *slog.Logger
	UserRepository repository.UserRepository
}

func NewUserService(
	DB client.DatabaseService,
	DBLogger *slog.Logger,
	userRepository repository.UserRepository,
) UserService {
	return &userService{
		DB:             DB,
		DBLogger:       DBLogger,
		UserRepository: userRepository,
	}
}

func (s *userService) GetUsers(w http.ResponseWriter, r *http.Request) ([]dto.UserResponse, error) {
	var usersRepsonseDTO []dto.UserResponse

	users, err := s.UserRepository.GetAll(r.Context())
	if err != nil {
		return usersRepsonseDTO, err
	}

	for _, user := range users {
		// Ensure FirstName is not empty
		if len(user.FirstName) == 0 {
			continue // Skip users with an empty first name
		}

		// Format the name with the first letter of FirstName capitalized
		formattedName := fmt.Sprintf("%s.%s", strings.ToUpper(string(user.FirstName[0])), user.LastName)

		// Append to response slice
		usersRepsonseDTO = append(usersRepsonseDTO, dto.UserResponse{
			UserID:  user.ID,
			Name:    formattedName,
			Country: user.Country,
		})
	}

	return usersRepsonseDTO, nil
}

func (s *userService) GetUser(w http.ResponseWriter, r *http.Request) (dto.UserResponse, error) {
	var userResponseDTO dto.UserResponse

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return userResponseDTO, fmt.Errorf("missing id parameter")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, "id parameter is not an int", nil, r)
		return userResponseDTO, fmt.Errorf("id must be a valid integer")
	}

	user, err := s.UserRepository.GetByID(r.Context(), id)
	if err != nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, "failed to get user", nil, r)
		return userResponseDTO, customError.ErrInternalServerError
	}

	if user == nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, "failed to get user", nil, r)
		return userResponseDTO, customError.ErrInternalServerError
	}

	return dto.UserResponse{
		UserID:  user.ID,
		Name:    fmt.Sprintf("%s.%s", strings.ToUpper(string(user.FirstName[0])), user.LastName),
		Country: user.Country,
	}, nil
}
