package repo

import (
	"context"
	"errors"
	"fmt"

	"user-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
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
