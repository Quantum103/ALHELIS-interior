package repository

import (
	"context"

	"auth-service/internal/models"

	"github.com/jackc/pgx/v5"
)

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type UserRepository struct {
	db DB
}

func NewUserRepository(db DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, username, email, passwordHash string) (int64, error) {
	query := `
			INSERT INTO users (username, email, password_hash)
			VALUES ($1, $2, $3)
			RETURNING id
		`

	var userID int64

	err := r.db.QueryRow(
		ctx,
		query,
		username,
		email,
		passwordHash,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*models.UserResponse, string, error) {
	query := `SELECT id, username, email, password_hash FROM users WHERE email = $1 OR username = $1`
	var user models.UserResponse
	var hash string
	err := r.db.QueryRow(ctx, query, login).Scan(&user.ID, &user.Username, &user.Email, &hash)
	return &user, hash, err
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.UserResponse, error) {
	query := `SELECT id, username, email FROM users WHERE id = $1`
	var user models.UserResponse
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email)
	return &user, err
}
