package repository

import (
	"context"
	"database/sql"

	"github.com/brightside-dev/ronin-fitness-be/database/client"
	"github.com/brightside-dev/ronin-fitness-be/internal/model"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, tx *sql.Tx, refreshToken *model.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	DeleteByToken(ctx context.Context, token string) error
}

type refreshTokenRepository struct {
	db client.DatabaseService
}

func NewUserRefreshTokenRepository(db client.DatabaseService) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, tx *sql.Tx, refreshToken *model.RefreshToken) error {
	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO refresh_tokens (user_id, token, expires_at, ip_address, user_agent) VALUES (?, ?, ?, ?, ?)",
		refreshToken.UserID,
		refreshToken.Token,
		refreshToken.ExpiresAt,
		refreshToken.IPAddress,
		refreshToken.UserAgent,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *refreshTokenRepository) GetByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken

	row := r.db.QueryRowContext(ctx, "SELECT * FROM refresh_tokens WHERE token = ?", token)

	if err := row.Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.CreatedAt,
		&refreshToken.ExpiresAt,
		&refreshToken.Revoked,
		&refreshToken.UserAgent,
		&refreshToken.IPAddress,
	); err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *refreshTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE token = ?", token)
	if err != nil {
		return err
	}

	return nil
}
