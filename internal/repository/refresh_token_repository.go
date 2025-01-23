package repository

import (
	"context"
	"fmt"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/database"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/model"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/util"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) error
	GetRefreshTokenByToken(ctx context.Context, token string) (*model.RefreshToken, error)
}

type refreshTokenRepository struct {
	db database.Service
}

func NewRefreshTokenRepository(db database.Service) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) error {
	// Create refresh token
	// Begin a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Ensure rollback is called if the function exits without committing
	defer func() {
		tx.Rollback()
	}()

	_, err = tx.ExecContext(
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

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *refreshTokenRepository) GetRefreshTokenByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken

	row := r.db.QueryRowContext(ctx, "SELECT * FROM refresh_tokens WHERE token = ?", token)

	var expiresAtRaw interface{}

	if err := row.Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.CreatedAt,
		&expiresAtRaw,
		&refreshToken.RevokedAt,
		&refreshToken.UserAgent,
		&refreshToken.IPAddress,
	); err != nil {
		return nil, err
	}

	if expiresAtRaw != nil {
		expiresAt, err := util.ParseDateTime(expiresAtRaw)
		if err != nil {
			return nil, err
		}
		refreshToken.ExpiresAt = expiresAt
	}

	fmt.Printf("Refresh Token: %+v\n", refreshToken)

	return &refreshToken, nil
}
