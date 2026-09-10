package repo

import (
	"context"
	"errors"
	"fmt"

	"user-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type UserRepository struct {
	db DB
}

func NewUserRepository(db DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
func (r *UserRepository) GetProfile(ctx context.Context, userID int64) (*models.UserProfile, error) {
	var profile models.UserProfile

	err := r.db.QueryRow(
		ctx,
		`
		SELECT user_id, name, phone
		FROM user_profiles
		WHERE user_id = $1
		`,
		userID,
	).Scan(
		&profile.UserID,
		&profile.Name,
		&profile.Phone,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("profile not found")
		}

		return nil, fmt.Errorf("get profile: %w", err)
	}

	return &profile, nil
}
func (r *UserRepository) CreateProfile(ctx context.Context, userID int64, name string, phone string) error {
	query := `
		INSERT INTO user_profiles (
			user_id,
			name,
			phone,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (user_id) DO NOTHING
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		name,
		phone,
	)

	if err != nil {
		return fmt.Errorf("repository create profile: %w", err)
	}

	return nil
}
