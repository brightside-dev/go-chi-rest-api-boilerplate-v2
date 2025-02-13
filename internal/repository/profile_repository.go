package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/brightside-dev/ronin-fitness-be/database/client"
	"github.com/brightside-dev/ronin-fitness-be/internal/model"
)

type ProfileRepository interface {
	GetByUserID(ctx context.Context, id int) (*model.ProfileWithUser, error)
	Create(ctx context.Context, tx *sql.Tx, profile *model.Profile) (*model.Profile, error)
	Update(ctx context.Context, tx *sql.Tx, profile *model.Profile) (*model.Profile, error)
}

type profileRepository struct {
	db client.DatabaseService
}

func NewProfileRepository(db client.DatabaseService) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) GetByUserID(ctx context.Context, userID int) (*model.ProfileWithUser, error) {
	fmt.Println(userID)
	query := `
        SELECT 
            p.id, p.user_id, p.display_name, p.privacy, p.avatar_version, 
            p.is_notifications_enabled, p.fitness_experience, p.experience_points, 
            p.created_at, p.updated_at, 
            u.first_name, u.last_name, u.country 
        FROM profiles p 
        INNER JOIN users u ON p.user_id = u.id 
        WHERE p.user_id = ?
    `
	row := r.db.QueryRowContext(ctx, query, userID)

	var profile model.ProfileWithUser
	err := row.Scan(
		&profile.ID,
		&profile.UserID,
		&profile.DisplayName,
		&profile.Privacy,
		&profile.AvatarVersion,
		&profile.IsNotificationsEnabled,
		&profile.FitnessExperience,
		&profile.ExperiencePoints,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&profile.FirstName,
		&profile.LastName,
		&profile.Country,
	)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *profileRepository) Create(ctx context.Context, tx *sql.Tx, profile *model.Profile) (*model.Profile, error) {
	query := `
		INSERT INTO profiles 
			(user_id, display_name, privacy, avatar_version, is_notifications_enabled, fitness_experience, experience_points) 
		VALUES 
			(?, ?, ?, ?, ?, ?, ?)
	`
	res, err := tx.ExecContext(
		ctx, query,
		profile.UserID,
		profile.DisplayName,
		profile.Privacy,
		profile.AvatarVersion,
		profile.IsNotificationsEnabled,
		profile.FitnessExperience,
		profile.ExperiencePoints,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	newProfile := &model.Profile{
		ID:                     int(id),
		UserID:                 profile.UserID,
		DisplayName:            profile.DisplayName,
		Privacy:                profile.Privacy,
		AvatarVersion:          profile.AvatarVersion,
		IsNotificationsEnabled: profile.IsNotificationsEnabled,
		FitnessExperience:      profile.FitnessExperience,
		ExperiencePoints:       profile.ExperiencePoints,
	}

	return newProfile, nil
}

func (r *profileRepository) Update(ctx context.Context, tx *sql.Tx, profile *model.Profile) (*model.Profile, error) {
	query := `
		UPDATE profiles p
		JOIN users u ON p.user_id = u.id
		SET 
			p.display_name = COALESCE(NULLIF(?, ''), p.display_name),
			p.avatar_version = COALESCE(NULLIF(?, 0), p.avatar_version),
			p.privacy = COALESCE(NULLIF(?, ''), p.privacy),
			p.fitness_experience = COALESCE(NULLIF(?, ''), p.fitness_experience)
		WHERE p.user_id = ?
		AND (
			p.display_name != COALESCE(NULLIF(?, ''), p.display_name) OR
			p.avatar_version != COALESCE(NULLIF(?, 0), p.avatar_version) OR
			p.privacy != COALESCE(NULLIF(?, ''), p.privacy) OR
			p.fitness_experience != COALESCE(NULLIF(?, ''), p.fitness_experience)
		);

		SELECT 
			p.id AS profile_id,
			p.display_name,
			p.avatar_version,
			p.privacy,
			p.fitness_experience,
			p.experience_points,
			u.id AS user_id,
			u.first_name,
			u.last_name,
			u.country
		FROM profiles p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?;
	`
	_, err := r.db.ExecContext(
		ctx, query,
		profile.DisplayName,
		profile.AvatarVersion,
		profile.Privacy,
		profile.FitnessExperience,
		profile.IsNotificationsEnabled,
		profile.ExperiencePoints,
		profile.ID,
	)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
