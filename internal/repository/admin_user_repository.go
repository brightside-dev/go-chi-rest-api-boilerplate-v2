package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/database"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/model"
)

type AdminUserRepository interface {
	GetAll(ctx context.Context) ([]model.AdminUser, error)
	GetByID(ctx context.Context, id int) (*model.AdminUser, error)
	GetByEmail(ctx context.Context, email string) (*model.AdminUser, error)
	Create(ctx context.Context, user *model.AdminUser) (*model.AdminUser, error)
}

type adminUserRepository struct {
	db database.Service
}

func NewAdminUserRepository(db database.Service) AdminUserRepository {
	return &adminUserRepository{db: db}
}

func (r *adminUserRepository) GetAll(ctx context.Context) ([]model.AdminUser, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM admin_users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.AdminUser
	for rows.Next() {
		var user model.AdminUser

		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (r *adminUserRepository) GetByID(ctx context.Context, id int) (*model.AdminUser, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM admin_users WHERE id = ?", id)

	var user model.AdminUser

	// Scan the row, using birthdayRaw to handle the birthday field temporarily
	if err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

func (r *adminUserRepository) GetByEmail(ctx context.Context, email string) (*model.AdminUser, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM admin_users WHERE email = ?", email)

	var user model.AdminUser

	// Scan the row, using birthdayRaw to handle the birthday field temporarily
	if err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

func (r *adminUserRepository) Create(ctx context.Context, user *model.AdminUser) (*model.AdminUser, error) {
	// Begin a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.AdminUser{}, err
	}

	// Ensure rollback is called if the function exits without committing
	defer func() {
		tx.Rollback()
	}()

	result, err := tx.ExecContext(
		ctx,
		"INSERT INTO admin_users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return &model.AdminUser{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return &model.AdminUser{}, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return &model.AdminUser{}, err
	}

	newAdminUser := &model.AdminUser{
		ID:        int(id),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}

	return newAdminUser, nil
}
