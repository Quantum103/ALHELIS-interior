package service

import (
	"context"
	"testing"

	"auth-service/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	mockCreate     func(ctx context.Context, username, email, passwordHash string) error
	mockGetByLogin func(ctx context.Context, login string) (*models.UserResponse, string, error)
	mockGetByID    func(ctx context.Context, id int64) (*models.UserResponse, error)
}

func (m *mockUserRepository) Create(ctx context.Context, username, email, passwordHash string) error {
	return m.mockCreate(ctx, username, email, passwordHash)
}
func (m *mockUserRepository) GetByLogin(ctx context.Context, login string) (*models.UserResponse, string, error) {
	return m.mockGetByLogin(ctx, login)
}
func (m *mockUserRepository) GetByID(ctx context.Context, id int64) (*models.UserResponse, error) {
	return m.mockGetByID(ctx, id)
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := &mockUserRepository{
		mockCreate: func(ctx context.Context, username, email, passwordHash string) error {
			return nil
		},
	}
	svc := NewAuthService(mockRepo, "test-secret-key")
	req := models.RegisterRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
	}

	err := svc.Register(context.Background(), req)
	if err != nil {
		t.Errorf("Ожидался nil, получена ошибка: %v", err)
	}
}

func TestAuthService_Login(t *testing.T) {
	password := "mypassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockRepo := &mockUserRepository{
		mockGetByLogin: func(ctx context.Context, login string) (*models.UserResponse, string, error) {
			return &models.UserResponse{ID: 1, Username: "test", Email: "test@test.com"}, string(hashedPassword), nil
		},
	}
	svc := NewAuthService(mockRepo, "test-secret-key")
	t.Run("Valid Credentials", func(t *testing.T) {
		req := models.LoginRequest{Username: "test", Password: password}
		token, user, err := svc.Login(context.Background(), req)

		if err != nil {
			t.Errorf("Ожидался nil, получена ошибка: %v", err)
		}
		if token == "" {
			t.Errorf("Ожидался JWT токен, получена пустая строка")
		}
		if user.ID != 1 {
			t.Errorf("Ожидался ID=1, получено %d", user.ID)
		}
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		req := models.LoginRequest{Username: "test", Password: "wrongpassword"}
		_, _, err := svc.Login(context.Background(), req)

		if err != ErrInvalidCredentials {
			t.Errorf("Ожидалась ошибка ErrInvalidCredentials, получена %v", err)
		}
	})
}
