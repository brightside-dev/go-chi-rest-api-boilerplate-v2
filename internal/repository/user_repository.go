package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/brightside-dev/ronin-fitness-be/database/client"
	"github.com/brightside-dev/ronin-fitness-be/internal/model"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error)
	Update(ctx context.Context, tx *sql.Tx, user *model.User) error
}

type userRepository struct {
	db client.DatabaseService
}

func NewUserRepository(db client.DatabaseService) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.Birthday,
			&user.Country,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)

	var user model.User
	if err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Birthday,
		&user.Country,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE email = ?", email)

	var user model.User
	// Scan the row, using birthdayRaw to handle the birthday field temporarily
	if err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Birthday,
		&user.Country,
		&user.CreatedAt,
		&user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error) {
	result, err := tx.ExecContext(
		ctx,
		"INSERT INTO users (first_name, last_name, email, password, country, birthday) VALUES (?, ?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Password, user.Country, user.Birthday)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		ID:        int(id),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Country:   user.Country,
		Birthday:  user.Birthday,
	}

	return newUser, nil
}

func (r *userRepository) Update(ctx context.Context, tx *sql.Tx, user *model.User) error {

	fmt.Println("Updating user: ", user)
	queryUpdate := `
		UPDATE users u
		SET 
			u.first_name = COALESCE(NULLIF(?, ''), u.first_name),
			u.last_name = COALESCE(NULLIF(?, ''), u.last_name),
			u.email = COALESCE(NULLIF(?, ''), u.email),
			u.password = COALESCE(NULLIF(?, ''), u.password),
			u.country = COALESCE(NULLIF(?, ''), u.country),
			u.birthday = COALESCE(NULLIF(?, '0000-00-00'), u.birthday),
			u.is_verified = COALESCE(?, u.is_verified)
		WHERE u.id = ?
	`

	_, err := tx.ExecContext(
		ctx,
		queryUpdate,
		user.FirstName, user.LastName, user.Email, user.Password, user.Country, user.Birthday, user.IsVerified, user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
