package oauth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brightside-dev/ronin-fitness-be/internal/handler/dto"
	customError "github.com/brightside-dev/ronin-fitness-be/internal/handler/error"
	oauth_client "github.com/brightside-dev/ronin-fitness-be/internal/service/oauth/client"
	"github.com/brightside-dev/ronin-fitness-be/internal/util"
)

type OAuthService interface {
	GetLoginURL(w http.ResponseWriter, r *http.Request) (*dto.OAuthLoginURLResponse, error)
	HandleCallback(w http.ResponseWriter, r *http.Request) (*dto.OAuthCallbackURLResponse, error)
}

type oAuthService struct {
	Logger       *slog.Logger
	FacebookAuth oauth_client.FacebookAuth
	GoogleAuth   oauth_client.GoogleAuth
}

func NewOAuthService(logger *slog.Logger) OAuthService {
	return &oAuthService{
		Logger:       logger,
		FacebookAuth: oauth_client.NewFacebookAuth(),
		GoogleAuth:   oauth_client.NewGoogleAuth(),
	}
}

func (s *oAuthService) GetLoginURL(w http.ResponseWriter, r *http.Request) (*dto.OAuthLoginURLResponse, error) {
	req := dto.OAuthLoginURLRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInvalidRequestBody
	}

	switch req.Client {
	case "facebook":
		url, err := s.FacebookAuth.GetLoginURL()
		if err != nil {
			util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
			return nil, err
		}

		return &dto.OAuthLoginURLResponse{LoginURL: url}, nil
	case "google":
		url, err := s.GoogleAuth.GetLoginURL()
		if err != nil {
			util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
			return nil, err
		}

		return &dto.OAuthLoginURLResponse{LoginURL: url}, nil
	default:
		return nil, fmt.Errorf("unsupported social login: %s", req.Client)
	}
}

func (s *oAuthService) HandleCallback(w http.ResponseWriter, r *http.Request) (*dto.OAuthCallbackURLResponse, error) {
	req := dto.OAuthCallbackURLRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
		return nil, customError.ErrInvalidRequestBody
	}

	switch req.Client {
	case "facebook":
		oAuthUserResponse, err := s.FacebookAuth.HandleCallback(req.AuthCode)
		if err != nil {
			util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
			return nil, err
		}

		return &dto.OAuthCallbackURLResponse{
			FirstName: oAuthUserResponse.FirstName,
			LastName:  oAuthUserResponse.LastName,
			Email:     oAuthUserResponse.Email,
		}, nil
	case "google":
		oAuthUserResponse, err := s.GoogleAuth.HandleCallback(req.AuthCode)
		if err != nil {
			util.LogWithContext(s.Logger, slog.LevelError, err.Error(), nil, r)
			return nil, err
		}

		return &dto.OAuthCallbackURLResponse{
			FirstName: oAuthUserResponse.FirstName,
			LastName:  oAuthUserResponse.LastName,
			Email:     oAuthUserResponse.Email,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported social login: %s", req.Client)
	}
}
