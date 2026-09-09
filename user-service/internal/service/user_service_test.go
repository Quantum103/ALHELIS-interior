package service

import (
	"context"
	"errors"
	"testing"

	"user-service/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetProfile(ctx context.Context, userID int64) (*models.UserProfile, error) {
	args := m.Called(ctx, userID)

	var profile *models.UserProfile
	if args.Get(0) != nil {
		profile = args.Get(0).(*models.UserProfile)
	}

	return profile, args.Error(1)
}

func (m *mockUserRepository) CreateProfile(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestGetProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)

		profile := &models.UserProfile{
			UserID: 1,
			Name:   "Test",
			Phone:  "123456789",
		}

		repo.On("GetProfile", mock.Anything, int64(1)).
			Return(profile, nil)

		service := NewUserService(repo)

		result, err := service.GetProfile(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, profile, result)

		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(mockUserRepository)

		repo.On("GetProfile", mock.Anything, int64(1)).
			Return(nil, errors.New("profile not found"))

		service := NewUserService(repo)

		result, err := service.GetProfile(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, result)

		repo.AssertExpectations(t)
	})
}

func TestCreateProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)

		repo.On("CreateProfile", mock.Anything, int64(1)).
			Return(nil)

		service := NewUserService(repo)

		err := service.CreateProfile(context.Background(), 1)

		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(mockUserRepository)

		repo.On("CreateProfile", mock.Anything, int64(1)).
			Return(errors.New("database error"))

		service := NewUserService(repo)

		err := service.CreateProfile(context.Background(), 1)

		assert.Error(t, err)

		repo.AssertExpectations(t)
	})
}
