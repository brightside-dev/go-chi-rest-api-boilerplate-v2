package model

import "time"

type User struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	Birthday   time.Time `json:"birthday"`
	Country    string    `json:"country"`
	IsVerified int       `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
