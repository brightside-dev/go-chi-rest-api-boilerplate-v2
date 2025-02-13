package dto

type AuthRegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Country   string `json:"country" validate:"required"`
	Birthday  string `json:"birthday" validate:"required"`
}

type AuthRegisterResponse struct {
	UserID               int `json:"user_id"`
	RefreshTokenResponse `json:"refresh_token"`
}

type AuthLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthLoginResponse struct {
	UserResponse         `json:"profile"`
	RefreshTokenResponse `json:"refresh_token"`
}

type AuthVerifyRequest struct {
	UserID int    `json:"user_id" validate:"required"`
	Code   string `json:"code" validate:"required"`
}

type OAuthLoginURLRequest struct {
	Client string `json:"client" validate:"required"`
}

type OAuthLoginURLResponse struct {
	LoginURL string `json:"login_url"`
}

type OAuthCallbackURLRequest struct {
	AuthCode string `json:"auth_code" validate:"required"`
	Client   string `json:"client" validate:"required"`
}

type OAuthCallbackURLResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type OAuthUserResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
