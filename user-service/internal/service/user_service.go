package service

import (
	"context"
	"fmt"
	"log"

	"user-service/internal/models"
)

type UserRepository interface {
	GetProfile(ctx context.Context, userID int64) (*models.UserProfile, error)
	CreateProfile(ctx context.Context, userID int64, name string, phone string) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (*models.UserProfile, error) {
	return s.repo.GetProfile(ctx, userID)
}

func (s *UserService) CreateProfile(
	ctx context.Context,
	userID int64,
	name string,
	phone string,
) error {

	err := s.repo.CreateProfile(
		ctx,
		userID,
		name,
		phone,
	)

	if err != nil {
		log.Printf(
			"[user-service] ОШИБКА БД userID=%d: %v",
			userID,
			err,
		)

		return fmt.Errorf(
			"failed to create profile in db: %w",
			err,
		)
	}

	return nil
}
