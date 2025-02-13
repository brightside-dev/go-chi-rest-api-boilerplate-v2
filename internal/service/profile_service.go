package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/brightside-dev/ronin-fitness-be/database/client"
	"github.com/brightside-dev/ronin-fitness-be/internal/handler/dto"
	customError "github.com/brightside-dev/ronin-fitness-be/internal/handler/error"
	"github.com/brightside-dev/ronin-fitness-be/internal/model"
	"github.com/brightside-dev/ronin-fitness-be/internal/repository"
	"github.com/brightside-dev/ronin-fitness-be/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ProfileService interface {
	GetMyProfileByUserID(w http.ResponseWriter, r *http.Request, userID int) (*dto.MyProfileResponse, error)
	GetProfileByUserID(w http.ResponseWriter, r *http.Request) (*dto.ProfileResponse, error)
	FollowProfile(w http.ResponseWriter, r *http.Request) (*dto.FollowProfilesResponse, error)
	UnfollowProfile(w http.ResponseWriter, r *http.Request) (*dto.FollowProfilesResponse, error)
	UpdateProfile(w http.ResponseWriter, r *http.Request) (*dto.MyProfileResponse, error)
}

type profileService struct {
	DB                       client.DatabaseService
	DBLogger                 *slog.Logger
	Validate                 *validator.Validate
	ProfileRepository        repository.ProfileRepository
	ProfilesFollowRepository repository.ProfileFollowRepository
	UserRepository           repository.UserRepository
}

func NewProfileService(
	db client.DatabaseService,
	dbLogger *slog.Logger,
	validator *validator.Validate,
	profileRepository repository.ProfileRepository,
	profilesFollowRepository repository.ProfileFollowRepository,
	userRepository repository.UserRepository,
) ProfileService {
	return &profileService{
		DB:                db,
		DBLogger:          dbLogger,
		Validate:          validator,
		ProfileRepository: profileRepository,
		UserRepository:    userRepository,
	}
}

func (s *profileService) GetMyProfileByUserID(w http.ResponseWriter, r *http.Request, userID int) (*dto.MyProfileResponse, error) {
	profileWithUser, err := s.ProfileRepository.GetByUserID(r.Context(), userID)
	if err != nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInternalServerError
	}

	return &dto.MyProfileResponse{
		ProfileID:         profileWithUser.ID,
		DisplayName:       profileWithUser.DisplayName,
		AvatarVersion:     profileWithUser.AvatarVersion,
		Privacy:           profileWithUser.Privacy,
		FitnessExperience: profileWithUser.FitnessExperience,
		ExperiencePoints:  profileWithUser.ExperiencePoints,
	}, nil
}

func (s *profileService) GetProfileByUserID(w http.ResponseWriter, r *http.Request) (*dto.ProfileResponse, error) {
	idParam := chi.URLParam(r, "userId")
	if idParam == "" {
		util.LogWithContext(s.DBLogger, slog.LevelError, "missing user_id paramter", nil, r)
		return nil, fmt.Errorf("missing id parameter")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, "id parameter is not an int", nil, r)
		return nil, fmt.Errorf("id must be a valid integer")
	}

	profileWithUser, err := s.ProfileRepository.GetByUserID(r.Context(), id)
	if err != nil {
		return nil, err
	}

	return &dto.ProfileResponse{
		ProfileID:         profileWithUser.ID,
		DisplayName:       profileWithUser.DisplayName,
		AvatarVersion:     profileWithUser.AvatarVersion,
		Privacy:           profileWithUser.Privacy,
		FitnessExperience: profileWithUser.FitnessExperience,
		ExperiencePoints:  profileWithUser.ExperiencePoints,
		User: dto.UserResponse{
			UserID:  profileWithUser.UserID,
			Name:    fmt.Sprintf("%s.%s", strings.ToUpper(string(profileWithUser.FirstName[0])), profileWithUser.LastName),
			Country: profileWithUser.Country,
		},
	}, nil
}

func (s *profileService) FollowProfile(w http.ResponseWriter, r *http.Request) (*dto.FollowProfilesResponse, error) {
	req := dto.FollowProfilesRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInvalidRequestBody
	}

	if req.FollowingProfileID == 0 || req.FollowerProfileID == 0 {
		util.LogWithContext(s.DBLogger, slog.LevelError, "following_profile_id or follower_profile_id parameter missing", nil, r)
		return nil, fmt.Errorf("following_profile_id or follower_profile_id parameter missing")
	}

	profileFollow := &model.ProfileFollow{
		ProfileID:         req.FollowingProfileID,
		FollowerProfileID: req.FollowerProfileID,
	}

	newProfileFollow, err := s.ProfilesFollowRepository.Create(r.Context(), profileFollow)
	if err != nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInternalServerError
	}

	return &dto.FollowProfilesResponse{
		FollowingProfileID: newProfileFollow.ProfileID,
	}, nil
}

func (s *profileService) UpdateProfile(w http.ResponseWriter, r *http.Request) (*dto.MyProfileResponse, error) {
	req := dto.ProfileUpdateRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInvalidRequestBody
	}

	err = s.Validate.Struct(req)
	if err != nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, err.Error(), nil, r)
		return nil, err
	}

	profileWithUser, err := s.ProfileRepository.GetByUserID(r.Context(), req.UserID)
	if err != nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInternalServerError
	}

	profile := model.Profile{
		ID:                     req.ProfileID,
		UserID:                 req.UserID,
		DisplayName:            profileWithUser.DisplayName,
		Privacy:                profileWithUser.Privacy,
		AvatarVersion:          profileWithUser.AvatarVersion,
		IsNotificationsEnabled: profileWithUser.IsNotificationsEnabled,
		FitnessExperience:      profileWithUser.FitnessExperience,
		ExperiencePoints:       profileWithUser.ExperiencePoints,
	}

	isUpdating := false

	if (profile.DisplayName != req.DisplayName) && req.DisplayName != "" {
		profile.DisplayName = req.DisplayName
		isUpdating = true
	}

	if (profile.AvatarVersion != req.AvatarVersion) && req.AvatarVersion != 0 {
		profile.AvatarVersion = req.AvatarVersion
		isUpdating = true
	}

	if (profile.IsNotificationsEnabled != req.IsNotificationsEnabled) && req.IsNotificationsEnabled {
		profile.IsNotificationsEnabled = req.IsNotificationsEnabled
		isUpdating = true
	}

	if (profile.Privacy != req.Privacy) && req.Privacy != "" {
		profile.Privacy = req.Privacy
		isUpdating = true
	}

	if (profile.FitnessExperience != req.FitnessExperience) && req.FitnessExperience != "" {
		profile.FitnessExperience = req.FitnessExperience
		isUpdating = true
	}

	if !isUpdating {
		return nil, customError.ErrNothingToUpdate
	}

	tx, err := s.DB.BeginTx(r.Context(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	updatedProfile, err := s.ProfileRepository.Update(r.Context(), tx, &profile)
	if err != nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInternalServerError
	}

	return &dto.MyProfileResponse{
		ProfileID:         updatedProfile.ID,
		DisplayName:       updatedProfile.DisplayName,
		AvatarVersion:     updatedProfile.AvatarVersion,
		Privacy:           updatedProfile.Privacy,
		FitnessExperience: updatedProfile.FitnessExperience,
		ExperiencePoints:  updatedProfile.ExperiencePoints,
	}, nil
}

func (s *profileService) UnfollowProfile(w http.ResponseWriter, r *http.Request) (*dto.FollowProfilesResponse, error) {
	req := dto.FollowProfilesRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInvalidRequestBody
	}

	if req.FollowingProfileID == 0 || req.FollowerProfileID == 0 {
		util.LogWithContext(s.DBLogger, slog.LevelError, "following_profile_id or follower_profile_id parameter missing", nil, r)
		return nil, fmt.Errorf("following_profile_id or follower_profile_id parameter missing")
	}

	profileFollow := &model.ProfileFollow{
		ProfileID:         req.FollowingProfileID,
		FollowerProfileID: req.FollowerProfileID,
	}

	err = s.ProfilesFollowRepository.Delete(r.Context(), profileFollow)
	if err != nil {
		util.LogWithContext(s.DBLogger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInternalServerError
	}

	return &dto.FollowProfilesResponse{
		FollowingProfileID: profileFollow.ProfileID,
	}, nil
}
