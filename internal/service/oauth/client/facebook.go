package oauth_client

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"

	"github.com/brightside-dev/ronin-fitness-be/internal/handler/dto"
	_ "github.com/joho/godotenv/autoload"
)

type FacebookAuth interface {
	GetLoginURL() (string, error)
	HandleCallback(authCode string) (*dto.OAuthUserResponse, error)
}

type facebookAuth struct{}

func NewFacebookAuth() FacebookAuth {
	return &facebookAuth{}
}

var oauth2FacebookConfig = oauth2.Config{
	ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),     // Your Facebook App ID
	ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"), // Your Facebook App Secret
	RedirectURL:  os.Getenv("OAUTH_CALLBACK_URL"),     // Custom redirect URI for mobile apps
	Scopes:       []string{"public_profile", "email"}, // Scopes to request
	Endpoint:     facebook.Endpoint,
}

func (f *facebookAuth) GetLoginURL() (string, error) {
	byte := make([]byte, 16)
	rand.Read(byte)
	oauthStateString := base64.URLEncoding.EncodeToString(byte)

	return oauth2FacebookConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline), nil
}

func (f *facebookAuth) HandleCallback(authCode string) (*dto.OAuthUserResponse, error) {
	// Exchange the code for an access token
	token, err := oauth2FacebookConfig.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, err
	}

	// Use the token to request user information from Facebook Graph API
	client := oauth2FacebookConfig.Client(context.Background(), token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,first_name,last_name,email")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse user information
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	// Handle missing fields gracefully
	firstName, _ := userInfo["first_name"].(string)
	lastName, _ := userInfo["last_name"].(string)
	email, _ := userInfo["email"].(string)

	return &dto.OAuthUserResponse{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}, nil
}
