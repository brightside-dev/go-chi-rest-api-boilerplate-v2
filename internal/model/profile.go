package model

import "time"

type Profile struct {
	ID                     int       `json:"id"`
	UserID                 int       `json:"user_id"`
	DisplayName            string    `json:"display_name"`
	Privacy                string    `json:"privacy"`
	AvatarVersion          int       `json:"avatar_version"`
	IsNotificationsEnabled bool      `json:"is_notifications_enabled"`
	FitnessExperience      string    `json:"fitness_experience"`
	ExperiencePoints       int       `json:"experience_points"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type ProfileWithUser struct {
	Profile
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Country   string `json:"country"`
}
