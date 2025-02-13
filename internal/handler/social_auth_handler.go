package handler

import (
	"net/http"

	"github.com/brightside-dev/ronin-fitness-be/internal/handler/response"
	"github.com/brightside-dev/ronin-fitness-be/internal/service/oauth"
)

type SocialAuthHandler interface {
	HandleLoginURL() http.HandlerFunc
	HandleCallbackURL() http.HandlerFunc
}

type socialAuthHandler struct {
	APIResponse  response.APIResponseManager
	OAuthService oauth.OAuthService
}

func NewSocialAuthHandler(
	apiResponseManager response.APIResponseManager,
	oAuthService oauth.OAuthService,
) SocialAuthHandler {
	return &socialAuthHandler{
		APIResponse:  apiResponseManager,
		OAuthService: oAuthService,
	}
}

func (h *socialAuthHandler) HandleLoginURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oAuthLoginResponseDTO, err := h.OAuthService.GetLoginURL(w, r)
		if err != nil {
			h.APIResponse.ErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}

		h.APIResponse.SuccessResponse(w, r, oAuthLoginResponseDTO)
	}
}

func (h *socialAuthHandler) HandleCallbackURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
