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
