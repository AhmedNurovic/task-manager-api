package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ahmednurovic/task-manager-api/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&user.ID)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, email, password FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}