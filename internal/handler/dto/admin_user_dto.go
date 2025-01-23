package dto

type AdminUserCreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type AdminUserResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type AdminUserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminUserLoginResponse struct {
	AdminUser    AdminUserResponse `json:"admin_user"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
}
