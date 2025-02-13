package model

import "time"

type ProfileFollow struct {
	ID                int       `json:"id"`
	ProfileID         int       `json:"following_profile_id"`
	FollowerProfileID int       `json:"follower_profile_id"`
	CreatedAt         time.Time `json:"created_at"`
}
