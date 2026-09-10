package service

import (
	"context"
	"errors"
	"testing"

	"auth-service/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	createFunc     func(ctx context.Context, username, email, passwordHash string) (int64, error)
	getByLoginFunc func(ctx context.Context, login string) (*models.UserResponse, string, error)
	getByIDFunc    func(ctx context.Context, id int64) (*models.UserResponse, error)
}

func (m *mockUserRepo) Create(ctx context.Context, username, email, passwordHash string) (int64, error) {
	return m.createFunc(ctx, username, email, passwordHash)
}

func (m *mockUserRepo) GetByLogin(ctx context.Context, login string) (*models.UserResponse, string, error) {
	return m.getByLoginFunc(ctx, login)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int64) (*models.UserResponse, error) {
	return m.getByIDFunc(ctx, id)
}

type mockProfileCreator struct {
	createProfileFunc func(ctx context.Context, userID int64, name, phone string) error
}

func (m *mockProfileCreator) CreateProfile(ctx context.Context, userID int64, name, phone string) error {
	return m.createProfileFunc(ctx, userID, name, phone)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := &mockUserRepo{
		createFunc: func(ctx context.Context, username, email, passwordHash string) (int64, error) {
			if username != "testuser" || email != "test@test.com" {
				t.Errorf("Неверные аргументы Create: %s, %s", username, email)
			}
			return 42, nil // Возвращаем тестовый ID
		},
	}

	mockProfile := &mockProfileCreator{
		createProfileFunc: func(ctx context.Context, userID int64, name, phone string) error {
			if userID != 42 {
				t.Errorf("Ожидался userID 42, получен %d", userID)
			}
			if name != "testuser" {
				t.Errorf("Ожидалось имя 'testuser', получено '%s'", name)
			}
			if phone != "" {
				t.Errorf("Ожидался пустой телефон, получен '%s'", phone)
			}
			return nil
		},
	}

	authService := NewAuthService(mockRepo, "test-secret-key", mockProfile)

	req := models.RegisterRequest{
		Username: "testuser",
		Email:    "test@test.com",
		Password: "password123",
	}

	userID, err := authService.Register(context.Background(), req)

	if err != nil {
		t.Fatalf("Ожидалась ошибка nil, получена: %v", err)
	}
	if userID != 42 {
		t.Errorf("Ожидался userID 42, получен %d", userID)
	}
}

func TestAuthService_Register_ProfileCreatorError(t *testing.T) {
	mockRepo := &mockUserRepo{
		createFunc: func(ctx context.Context, username, email, passwordHash string) (int64, error) {
			return 42, nil
		},
	}

	expectedErr := errors.New("gRPC ошибка создания профиля")
	mockProfile := &mockProfileCreator{
		createProfileFunc: func(ctx context.Context, userID int64, name, phone string) error {
			return expectedErr
		},
	}

	authService := NewAuthService(mockRepo, "test-secret-key", mockProfile)
	req := models.RegisterRequest{Username: "test", Email: "t@t.com", Password: "123"}

	userID, err := authService.Register(context.Background(), req)

	// Проверка
	if err == nil {
		t.Fatal("Ожидалась ошибка, но получена nil")
	}
	if userID != 0 {
		t.Errorf("При ошибке userID должен быть 0, получен %d", userID)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	plainPassword := "supersecret"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	expectedUser := &models.UserResponse{
		ID:    10,
		Email: "user@example.com",
	}

	mockRepo := &mockUserRepo{
		getByLoginFunc: func(ctx context.Context, login string) (*models.UserResponse, string, error) {
			return expectedUser, string(hashedPassword), nil
		},
	}

	authService := NewAuthService(mockRepo, "my-jwt-secret", nil)
	req := models.LoginRequest{Username: "user@example.com", Password: plainPassword}

	tokenString, user, err := authService.Login(context.Background(), req)

	// Проверка
	if err != nil {
		t.Fatalf("Ожидалась ошибка nil, получена: %v", err)
	}
	if tokenString == "" {
		t.Error("Токен не должен быть пустым")
	}
	if user.ID != expectedUser.ID {
		t.Errorf("Ожидался пользователь с ID %d, получен %d", expectedUser.ID, user.ID)
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("my-jwt-secret"), nil
	})
	if !token.Valid {
		t.Error("Сгенерированный токен недействителен")
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	// Подготовка: пользователь не найден
	mockRepo := &mockUserRepo{
		getByLoginFunc: func(ctx context.Context, login string) (*models.UserResponse, string, error) {
			return nil, "", pgx.ErrNoRows
		},
	}

	authService := NewAuthService(mockRepo, "secret", nil)
	req := models.LoginRequest{Username: "nobody", Password: "123"}

	_, _, err := authService.Login(context.Background(), req)

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Ожидалась ошибка ErrInvalidCredentials, получена: %v", err)
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	wrongHash, _ := bcrypt.GenerateFromPassword([]byte("другой_пароль"), bcrypt.DefaultCost)

	mockRepo := &mockUserRepo{
		getByLoginFunc: func(ctx context.Context, login string) (*models.UserResponse, string, error) {
			return &models.UserResponse{ID: 1}, string(wrongHash), nil
		},
	}

	authService := NewAuthService(mockRepo, "secret", nil)
	req := models.LoginRequest{Username: "user", Password: "неверный_пароль"}

	_, _, err := authService.Login(context.Background(), req)

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Ожидалась ошибка ErrInvalidCredentials при неверном пароле, получена: %v", err)
	}
}

func TestAuthService_GetMe_Success(t *testing.T) {
	// Подготовка
	expectedUser := &models.UserResponse{ID: 99, Email: "me@test.com"}

	mockRepo := &mockUserRepo{
		getByIDFunc: func(ctx context.Context, id int64) (*models.UserResponse, error) {
			if id != 99 {
				t.Errorf("Ожидался ID 99, получен %d", id)
			}
			return expectedUser, nil
		},
	}

	authService := NewAuthService(mockRepo, "secret", nil)

	user, err := authService.GetMe(context.Background(), 99)

	if err != nil {
		t.Fatalf("Ожидалась ошибка nil, получена: %v", err)
	}
	if user.ID != expectedUser.ID {
		t.Errorf("Ожидался пользователь с ID %d, получен %d", expectedUser.ID, user.ID)
	}
}
