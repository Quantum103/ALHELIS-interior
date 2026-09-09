package service

import (
	"context"

	"user-service/internal/models"
	"user-service/internal/repo"
)

type UserService struct {
	repo *repo.UserRepository
}

func NewUserService(repo *repo.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetProfile(
	ctx context.Context,
	userID int64,
) (*models.UserProfile, error) {
	return s.repo.GetProfile(ctx, userID)
}

func (s *UserService) CreateProfile(
	ctx context.Context,
	userID int64,
) error {
	return s.repo.CreateProfile(ctx, userID)
}
