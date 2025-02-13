package dto

type UserResponse struct {
	UserID     int    `json:"id"`
	Name       string `json:"name"`
	Country    string `json:"country"`
	IsVerified int    `json:"is_verified"`
}
