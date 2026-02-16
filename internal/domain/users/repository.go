package users

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByUserID(ctx context.Context, user_id string) (*User, error)
	CountByDate(ctx context.Context, date string) (int, error)
	Create(ctx context.Context, user *User) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.GetContext(ctx, &user,
		"SELECT id, user_id, username, email FROM users WHERE email=$1",
		email,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.db.GetContext(ctx, &user,
		"SELECT id, user_id, username, email, password FROM users WHERE username=$1",
		username,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUserID(ctx context.Context, user_id string) (*User, error) {
	var user User
	err := r.db.GetContext(ctx, &user,
		"SELECT id, user_id, username, email FROM users WHERE user_id=$1",
		user_id,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CountByDate(ctx context.Context, date string) (int, error) {
	var count int

	pattern := fmt.Sprintf("USER-%s-%%", date)

	query := `
		SELECT COUNT(*) 
		FROM users 
		WHERE user_id LIKE $1
	`

	err := r.db.GetContext(ctx, &count, query, pattern)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (user_id, fullname, password, email, username, updated_at) VALUES ($1, $2, $3, $4, $5, NOW())`
	_, err := r.db.ExecContext(ctx, query, user.USER_ID, user.Fullname, user.Password, user.Email, user.Username)
	return err
}
