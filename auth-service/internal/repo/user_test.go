package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestUserRepository_Integration(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("Пропуск интеграционного теста: не задана переменная TEST_DB_DSN")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Не удалось подключиться к тестовой БД: %v", err)
	}
	defer pool.Close()

	_, err = pool.Exec(ctx, "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Ошибка очистки таблицы: %v", err)
	}

	repo := NewUserRepository(pool)
	testUsername := "testuser"
	testEmail := "test@example.com"
	testHash := "hashed_password_123"

	t.Run("Create User", func(t *testing.T) {
		err := repo.Create(ctx, testUsername, testEmail, testHash)
		if err != nil {
			t.Errorf("Ожидался nil, получено %v", err)
		}
	})

	t.Run("Get User By Login", func(t *testing.T) {
		user, hash, err := repo.GetByLogin(ctx, testEmail)
		if err != nil {
			t.Errorf("Ожидался nil, получено %v", err)
		}
		if user.Username != testUsername || hash != testHash {
			t.Errorf("Получены неверные данные пользователя")
		}
	})

	t.Run("Get User By ID", func(t *testing.T) {
		user, err := repo.GetByID(ctx, 1)
		if err != nil {
			t.Errorf("Ожидался nil, получено %v", err)
		}
		if user.Email != testEmail {
			t.Errorf("Получены неверные данные пользователя")
		}
	})
}
