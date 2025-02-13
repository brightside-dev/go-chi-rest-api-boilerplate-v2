package repository

import (
	"context"
	"database/sql"

	"github.com/brightside-dev/ronin-fitness-be/database/client"
	"github.com/brightside-dev/ronin-fitness-be/internal/model"
)

type VerificationCodeRepository interface {
	Create(ctx context.Context, tx *sql.Tx, verificationCode *model.VerificationCode) (*model.VerificationCode, error)
	GetByUserID(ctx context.Context, userID int) (*model.VerificationCode, error)
	GetByCode(ctx context.Context, code string, userID int) (*model.VerificationCode, error)
}

type verificationCodeRepository struct {
	db client.DatabaseService
}

func NewVerificationCodeRepository(db client.DatabaseService) VerificationCodeRepository {
	return &verificationCodeRepository{db: db}
}

func (r *verificationCodeRepository) Create(ctx context.Context, tx *sql.Tx, verificationCode *model.VerificationCode) (*model.VerificationCode, error) {
	query := `
		INSERT INTO verification_codes (user_id, email, code, expires_at)
		VALUES (?, ?, ?, NOW() + INTERVAL 1 HOUR)
	`

	result, err := tx.ExecContext(ctx, query, verificationCode.UserID, verificationCode.Email, verificationCode.Code)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	verificationCode.ID = int(id)

	return verificationCode, nil
}

func (r *verificationCodeRepository) GetByUserID(ctx context.Context, userID int) (*model.VerificationCode, error) {
	query := `
		SELECT *
		FROM verification_codes
		WHERE user_id = ?
	`

	row := r.db.QueryRowContext(ctx, query, userID)

	var verificationCode model.VerificationCode
	err := row.Scan(
		&verificationCode.ID,
		&verificationCode.UserID,
		&verificationCode.Email,
		&verificationCode.Code,
		&verificationCode.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &verificationCode, nil
}

func (r *verificationCodeRepository) GetByCode(ctx context.Context, code string, userID int) (*model.VerificationCode, error) {
	query := `
		SELECT *
		FROM verification_codes
		WHERE code = ?
		AND user_id = ?
	`

	row := r.db.QueryRowContext(ctx, query, code, userID)

	var verificationCode model.VerificationCode
	err := row.Scan(
		&verificationCode.ID,
		&verificationCode.UserID,
		&verificationCode.Email,
		&verificationCode.Code,
		&verificationCode.CreatedAt,
		&verificationCode.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &verificationCode, nil
}
