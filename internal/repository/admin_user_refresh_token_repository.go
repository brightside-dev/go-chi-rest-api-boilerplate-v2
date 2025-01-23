package repository

import (
	"context"
	"fmt"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/database"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/model"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/util"
)

type AdminUserRefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *model.AdminUserRefreshToken) error
	GetByToken(ctx context.Context, token string) (*model.AdminUserRefreshToken, error)
}

type adminUserRefreshTokenRepository struct {
	db database.Service
}

func NewAdminUserRefreshTokenRepository(db database.Service) AdminUserRefreshTokenRepository {
	return &adminUserRefreshTokenRepository{db: db}
}

func (r *adminUserRefreshTokenRepository) Create(ctx context.Context, refreshToken *model.AdminUserRefreshToken) error {
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

func (r *adminUserRefreshTokenRepository) GetByToken(ctx context.Context, token string) (*model.AdminUserRefreshToken, error) {
	var refreshToken model.AdminUserRefreshToken

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
