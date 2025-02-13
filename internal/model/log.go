package model

import "time"

type Log struct {
	ID        int       `json:"id"`
	Domain    string    `json:"domain"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Context   string    `json:"context"`
	CreatedAt time.Time `json:"created_at"`
}
