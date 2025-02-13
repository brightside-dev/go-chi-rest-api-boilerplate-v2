package handler

import (
	"log/slog"
	"net/http"

	"github.com/brightside-dev/ronin-fitness-be/internal/handler/response"
	"github.com/brightside-dev/ronin-fitness-be/internal/service"
	"github.com/go-chi/jwtauth/v5"
)

type ProfileHandler interface {
	GetMyProfile() http.HandlerFunc
	GetProfile() http.HandlerFunc
	UpdateProfile() http.HandlerFunc
	FollowProfile() http.HandlerFunc
	UnfollowProfile() http.HandlerFunc
	GetFollowers() http.HandlerFunc
	GetFollowing() http.HandlerFunc
	RemoveFollower() http.HandlerFunc
}

type profileHandler struct {
	APIResponse    response.APIResponseManager
	DBLogger       *slog.Logger
	JWTAuth        *jwtauth.JWTAuth
	ProfileService service.ProfileService
}

func NewProfileHandler(
	apiResponse response.APIResponseManager,
	dbLogger *slog.Logger,
	profileService service.ProfileService,
) ProfileHandler {
	return &profileHandler{
		DBLogger:       dbLogger,
		ProfileService: profileService,
	}
}

func (h *profileHandler) GetMyProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			h.APIResponse.ErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		}

		myProfileResponseDTO, err := h.ProfileService.GetMyProfileByUserID(w, r, claims["id"].(int))
		if err != nil {
			h.APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		h.APIResponse.SuccessResponse(w, r, myProfileResponseDTO)
	}
}

func (h *profileHandler) GetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileResponseDTO, err := h.ProfileService.GetProfileByUserID(w, r)
		if err != nil {
			h.APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		h.APIResponse.SuccessResponse(w, r, profileResponseDTO)
	}
}

func (h *profileHandler) UpdateProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileResponseDTO, err := h.ProfileService.UpdateProfile(w, r)
		if err != nil {
			h.APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		h.APIResponse.SuccessResponse(w, r, profileResponseDTO)
	}
}

func (h *profileHandler) FollowProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		followProfileResponseDTO, err := h.ProfileService.FollowProfile(w, r)
		if err != nil {
			h.APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		h.APIResponse.SuccessResponse(w, r, followProfileResponseDTO)
	}
}

func (h *profileHandler) UnfollowProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		unfollowProfileResponseDTO, err := h.ProfileService.UnfollowProfile(w, r)
		if err != nil {
			h.APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		h.APIResponse.SuccessResponse(w, r, unfollowProfileResponseDTO)
	}
}

func (h *profileHandler) GetFollowers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (h *profileHandler) GetFollowing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (h *profileHandler) RemoveFollower() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
