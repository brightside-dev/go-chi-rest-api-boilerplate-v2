package repository

import (
	"context"

	"github.com/brightside-dev/ronin-fitness-be/database/client"
	"github.com/brightside-dev/ronin-fitness-be/internal/model"
)

type ProfileFollowRepository interface {
	Create(ctx context.Context, profileFollow *model.ProfileFollow) (*model.ProfileFollow, error)
	Delete(ctx context.Context, profileFollow *model.ProfileFollow) error
}

type profileFollowRepository struct {
	db client.DatabaseService
}

func NewProfileFollowRepository(db client.DatabaseService) ProfileFollowRepository {
	return &profileFollowRepository{db: db}
}

func (r *profileFollowRepository) Create(ctx context.Context, profileFollow *model.ProfileFollow) (*model.ProfileFollow, error) {
	query := `
		INSERT INTO profile_follows (profile_id, follower_profile_id)
		VALUES (?, ?)
	`

	res, err := r.db.ExecContext(ctx, query, profileFollow.ProfileID, profileFollow.FollowerProfileID)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	profileFollow.ID = int(id)

	return profileFollow, nil
}

func (r *profileFollowRepository) Delete(ctx context.Context, profileFollow *model.ProfileFollow) error {
	query := `DELETE FROM profile_follows WHERE profile_id = ? AND follower_profile_id = ?`

	_, err := r.db.ExecContext(ctx, query, profileFollow.ProfileID, profileFollow.FollowerProfileID)
	if err != nil {
		return err
	}

	return nil
}
