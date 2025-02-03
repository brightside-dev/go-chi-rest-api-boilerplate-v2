package service

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/database"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/dto"
	customError "github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/handler/error"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/repository"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/util"

	"github.com/go-chi/chi/v5"
)

type UserService struct {
	db             *database.Service
	dbLogger       *slog.Logger
	UserRepository repository.UserRepositoryInterface
}

func NewUserService(
	db *database.Service,
	dbLogger *slog.Logger,
	userRepository repository.UserRepositoryInterface,
) *UserService {
	return &UserService{
		db:             db,
		dbLogger:       dbLogger,
		UserRepository: userRepository,
	}
}

func (s *UserService) GetUsers(w http.ResponseWriter, r *http.Request) ([]dto.UserResponse, error) {
	var usersRepsonseDTO []dto.UserResponse

	users, err := s.UserRepository.GetAllUsers(r.Context())
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
			ID:      user.ID,
			Name:    formattedName,
			Country: user.Country,
		})
	}

	return usersRepsonseDTO, nil
}

func (s *UserService) GetUser(w http.ResponseWriter, r *http.Request) (dto.UserResponse, error) {
	var userResponseDTO dto.UserResponse

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return userResponseDTO, fmt.Errorf("missing id parameter")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		util.LogWithContext(s.dbLogger, slog.LevelError, "id parameter is not an int", nil, r)
		return userResponseDTO, fmt.Errorf("id must be a valid integer")
	}

	user, err := s.UserRepository.GetUserByID(r.Context(), id)
	if err != nil {
		util.LogWithContext(s.dbLogger, slog.LevelError, "failed to get user", nil, r)
		return userResponseDTO, customError.ErrInternalServerError
	}

	if user == nil {
		util.LogWithContext(s.dbLogger, slog.LevelError, "failed to get user", nil, r)
		return userResponseDTO, customError.ErrInternalServerError
	}

	return dto.UserResponse{
		ID:      user.ID,
		Name:    fmt.Sprintf("%s.%s", strings.ToUpper(string(user.FirstName[0])), user.LastName),
		Country: user.Country,
	}, nil
}
