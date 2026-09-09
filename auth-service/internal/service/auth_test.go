package service

import (
	"auth-service/internal/models"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Create(ctx context.Context, username, email, passwordHash string) (int64, error) {
	args := m.Called(ctx, username, email, passwordHash)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserRepository) GetByLogin(ctx context.Context, login string) (*models.UserResponse, string, error) {
	args := m.Called(ctx, login)

	var user *models.UserResponse
	if args.Get(0) != nil {
		user = args.Get(0).(*models.UserResponse)
	}

	return user, args.String(1), args.Error(2)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int64) (*models.UserResponse, error) {
	args := m.Called(ctx, id)

	var user *models.UserResponse
	if args.Get(0) != nil {
		user = args.Get(0).(*models.UserResponse)
	}

	return user, args.Error(1)
}

type mockProfileCreator struct {
	mock.Mock
}

func (m *mockProfileCreator) CreateProfile(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestRegister(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)
		profile := new(mockProfileCreator)

		repo.On(
			"Create",
			mock.Anything,
			"test",
			"test@example.com",
			mock.AnythingOfType("string"),
		).Return(int64(1), nil)

		profile.On("CreateProfile", mock.Anything, int64(1)).
			Return(nil)

		s := NewAuthService(repo, "secret", profile)

		req := models.RegisterRequest{
			Username: "test",
			Email:    "test@example.com",
			Password: "password",
		}

		id, err := s.Register(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)

		repo.AssertExpectations(t)
		profile.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(mockUserRepository)
		profile := new(mockProfileCreator)

		repo.On(
			"Create",
			mock.Anything,
			"test",
			"test@example.com",
			mock.AnythingOfType("string"),
		).Return(int64(0), errors.New("database error"))

		s := NewAuthService(repo, "secret", profile)

		req := models.RegisterRequest{
			Username: "test",
			Email:    "test@example.com",
			Password: "password",
		}

		id, err := s.Register(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, int64(0), id)

		profile.AssertNotCalled(t, "CreateProfile")
	})

	t.Run("profile error", func(t *testing.T) {
		repo := new(mockUserRepository)
		profile := new(mockProfileCreator)

		repo.On(
			"Create",
			mock.Anything,
			"test",
			"test@example.com",
			mock.AnythingOfType("string"),
		).Return(int64(1), nil)

		profile.On("CreateProfile", mock.Anything, int64(1)).
			Return(errors.New("profile error"))

		s := NewAuthService(repo, "secret", profile)

		req := models.RegisterRequest{
			Username: "test",
			Email:    "test@example.com",
			Password: "password",
		}

		id, err := s.Register(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
	})
}

func TestLogin(t *testing.T) {
	user := &models.UserResponse{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("password"),
		bcrypt.DefaultCost,
	)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)

		repo.On(
			"GetByLogin",
			mock.Anything,
			"test",
		).Return(user, string(passwordHash), nil)

		s := NewAuthService(repo, "secret", nil)

		tokenString, resultUser, err := s.Login(
			context.Background(),
			models.LoginRequest{
				Username: "test",
				Password: "password",
			},
		)

		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)
		assert.Equal(t, user, resultUser)

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {
				return []byte("secret"), nil
			},
		)

		assert.NoError(t, err)
		assert.True(t, token.Valid)

		claims := token.Claims.(jwt.MapClaims)

		assert.Equal(t, float64(1), claims["sub"])
		assert.Equal(t, "test@example.com", claims["email"])

		exp := int64(claims["exp"].(float64))
		assert.True(t, exp > time.Now().Unix())

		repo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := new(mockUserRepository)

		repo.On(
			"GetByLogin",
			mock.Anything,
			"test",
		).Return(nil, "", pgx.ErrNoRows)

		s := NewAuthService(repo, "secret", nil)

		token, resultUser, err := s.Login(
			context.Background(),
			models.LoginRequest{
				Username: "test",
				Password: "password",
			},
		)

		assert.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Empty(t, token)
		assert.Nil(t, resultUser)
	})

	t.Run("wrong password", func(t *testing.T) {
		repo := new(mockUserRepository)

		repo.On(
			"GetByLogin",
			mock.Anything,
			"test",
		).Return(user, string(passwordHash), nil)

		s := NewAuthService(repo, "secret", nil)

		token, resultUser, err := s.Login(
			context.Background(),
			models.LoginRequest{
				Username: "test",
				Password: "wrong",
			},
		)

		assert.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Empty(t, token)
		assert.Nil(t, resultUser)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(mockUserRepository)

		repo.On(
			"GetByLogin",
			mock.Anything,
			"test",
		).Return(nil, "", errors.New("database error"))

		s := NewAuthService(repo, "secret", nil)

		token, resultUser, err := s.Login(
			context.Background(),
			models.LoginRequest{
				Username: "test",
				Password: "password",
			},
		)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, resultUser)
	})
}

func TestGetMe(t *testing.T) {
	user := &models.UserResponse{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}

	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)

		repo.On("GetByID", mock.Anything, int64(1)).
			Return(user, nil)

		s := NewAuthService(repo, "secret", nil)

		result, err := s.GetMe(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, user, result)

		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(mockUserRepository)

		repo.On("GetByID", mock.Anything, int64(1)).
			Return(nil, errors.New("not found"))

		s := NewAuthService(repo, "secret", nil)

		result, err := s.GetMe(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
