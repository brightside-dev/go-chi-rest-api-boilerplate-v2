package model

import "time"

type WorkoutReaction struct {
	ID        int       `json:"id"`
	WorkoutID int       `json:"workout_id"`
	ProfileID int       `json:"profile_id"`
	Reaction  string    `json:"reaction"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
