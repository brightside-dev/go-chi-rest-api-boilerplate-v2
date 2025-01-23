package dto

type CreateUserRequest struct {
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  *string       `json:"token"`
	RefreshToken *string       `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  *string `json:"token"`
	RefreshToken *string `json:"refresh_token"`
}
