package oauth_client

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/brightside-dev/ronin-fitness-be/internal/handler/dto"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	_ "github.com/joho/godotenv/autoload"
)

type GoogleAuth interface {
	GetLoginURL() (string, error)
	HandleCallback(authCode string) (*dto.OAuthUserResponse, error)
}

type googleAuth struct{}

var (
	oauth2GoogleConfig = oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
)

func NewGoogleAuth() GoogleAuth {
	return &googleAuth{}
}

func (g *googleAuth) GetLoginURL() (string, error) {
	byte := make([]byte, 16)
	rand.Read(byte)
	oauthStateString := base64.URLEncoding.EncodeToString(byte)

	return oauth2GoogleConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline), nil
}

func (g *googleAuth) HandleCallback(authCode string) (*dto.OAuthUserResponse, error) {
	// Exchange the code for a token
	token, err := oauth2GoogleConfig.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, err
	}

	// Use the token to get user information
	client := oauth2GoogleConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse user info and use it
	// (you can use this data to create or update the user in your database)
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &dto.OAuthUserResponse{
		FirstName: userInfo["given_name"].(string),
		LastName:  userInfo["family_name"].(string),
		Email:     userInfo["email"].(string),
	}, nil
}
