package dto

type ProfileRequest struct {
	UserID int `json:"user_id" validate:"required"`
}

type MyProfileResponse struct {
	ProfileID         int    `json:"id"`
	DisplayName       string `json:"display_name"`
	AvatarVersion     int    `json:"avatar_version"`
	Privacy           string `json:"privacy"`
	FitnessExperience string `json:"fitness_experience"`
	ExperiencePoints  int    `json:"experience_points"`
}

type ProfileResponse struct {
	ProfileID         int          `json:"id"`
	DisplayName       string       `json:"display_name"`
	AvatarVersion     int          `json:"avatar_version"`
	Privacy           string       `json:"privacy"`
	FitnessExperience string       `json:"fitness_experience"`
	ExperiencePoints  int          `json:"experience_points"`
	User              UserResponse `json:"user"`
}

type ProfileUpdateRequest struct {
	UserID                 int    `json:"user_id" validate:"required"`
	ProfileID              int    `json:"profile_id" validate:"required"`
	DisplayName            string `json:"display_name" `
	AvatarVersion          int    `json:"avatar_version"`
	IsNotificationsEnabled bool   `json:"is_notifications_enabled"`
	Privacy                string `json:"privacy"`
	FitnessExperience      string `json:"fitness_experience"`
}

type FollowProfilesRequest struct {
	FollowingProfileID int `json:"following_profile_id" validate:"required"`
	FollowerProfileID  int `json:"follower_profile_id" validate:"required"`
}

type FollowProfilesResponse struct {
	FollowingProfileID int `json:"following_profile_id"`
}
