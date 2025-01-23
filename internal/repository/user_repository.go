package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/brightside-dev/boxing-be/internal/database"
	"github.com/brightside-dev/boxing-be/internal/model"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
}

type userRepository struct {
	db database.Service
}

func NewUserRepository(db database.Service) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var birthdayRaw interface{}
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&birthdayRaw,
			&user.Country,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		// Convert birthdayRaw to time.Time
		if birthdayRaw != nil {
			switch v := birthdayRaw.(type) {
			case time.Time:
				user.Birthday = v
			case []uint8:
				// Convert from string (in []uint8) to time.Time
				parsedTime, err := time.Parse("2006-01-02", string(v)) // Adjust the format as per your DB field format
				if err != nil {
					return nil, fmt.Errorf("failed to parse birthday: %w", err)
				}
				user.Birthday = parsedTime
			default:
				return nil, fmt.Errorf("unexpected type for birthday: %T", v)
			}
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)

	var user model.User
	var birthdayRaw interface{}

	// Scan the row, using birthdayRaw to handle the birthday field temporarily
	if err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&birthdayRaw,
		&user.Country,
		&user.CreatedAt,
		&user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("user not found: %w", err)
	}

	fmt.Printf("user: %+v\n", user)
	// Convert birthdayRaw to time.Time
	if birthdayRaw != nil {
		switch v := birthdayRaw.(type) {
		case time.Time:
			user.Birthday = v
		case []uint8:
			// Convert from string (in []uint8) to time.Time
			parsedTime, err := time.Parse("2006-01-02", string(v)) // Adjust the format as per your DB field format
			if err != nil {
				return nil, fmt.Errorf("failed to parse birthday: %w", err)
			}
			user.Birthday = parsedTime
		default:
			return nil, fmt.Errorf("unexpected type for birthday: %T", v)
		}
	}

	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	// Begin a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.User{}, err
	}

	// Ensure rollback is called if the function exits without committing
	defer func() {
		tx.Rollback()
	}()

	result, err := tx.ExecContext(
		ctx,
		"INSERT INTO users (first_name, last_name, email, password, country, birthday) VALUES (?, ?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Password, user.Country, user.Birthday)
	if err != nil {
		return &model.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return &model.User{}, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return &model.User{}, err
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
