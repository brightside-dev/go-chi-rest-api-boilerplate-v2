package model

import "time"

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Birthday  time.Time `json:"birthday"`
	Country   string    `json:"country"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type AdminUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Log struct {
	ID        int    `json:"id"`
	Domain    string `json:"domain"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Context   string `json:"context"`
	CreatedAt string `json:"created_at"`
}

type UserRefreshToken struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	UserID    int       `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	ExpiresAt time.Time `json:"updated_at"`
	RevokedAt bool      `json:"revoked_at"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
}

type AdminUserRefreshToken struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	UserID    int       `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	ExpiresAt time.Time `json:"updated_at"`
	RevokedAt bool      `json:"revoked_at"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
}
