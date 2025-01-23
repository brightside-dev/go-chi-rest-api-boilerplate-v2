package dto

type UserCreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Country   string `json:"country"`
	Birthday  string `json:"birthday"`
}

type UserResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
}

type UserRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserRefreshTokenResponse struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
