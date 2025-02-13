package model

import "time"

type RefreshToken struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	UserID    int       `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	ExpiresAt time.Time `json:"updated_at"`
	Revoked   bool      `json:"revoked_at"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
}
