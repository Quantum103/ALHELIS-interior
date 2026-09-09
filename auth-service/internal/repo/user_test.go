package repository

import (
	"context"
	"errors"
	"testing"

	"auth-service/internal/models"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer db.Close()

		repo := NewUserRepository(db)

		db.ExpectQuery("INSERT INTO users").
			WithArgs("test", "test@example.com", "hash").
			WillReturnRows(
				pgxmock.NewRows([]string{"id"}).
					AddRow(int64(1)),
			)

		id, err := repo.Create(
			context.Background(),
			"test",
			"test@example.com",
			"hash",
		)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
		assert.NoError(t, db.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		db, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer db.Close()

		repo := NewUserRepository(db)

		db.ExpectQuery("INSERT INTO users").
			WithArgs("test", "test@example.com", "hash").
			WillReturnError(errors.New("database error"))

		id, err := repo.Create(
			context.Background(),
			"test",
			"test@example.com",
			"hash",
		)

		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.NoError(t, db.ExpectationsWereMet())
	})
}

func TestUserRepository_GetByLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer db.Close()

		repo := NewUserRepository(db)

		db.ExpectQuery("SELECT id, username, email, password_hash FROM users").
			WithArgs("test").
			WillReturnRows(
				pgxmock.NewRows(
					[]string{"id", "username", "email", "password_hash"},
				).AddRow(
					int64(1),
					"test",
					"test@example.com",
					"hash",
				),
			)

		user, hash, err := repo.GetByLogin(
			context.Background(),
			"test",
		)

		assert.NoError(t, err)
		assert.Equal(t, &models.UserResponse{
			ID:       1,
			Username: "test",
			Email:    "test@example.com",
		}, user)
		assert.Equal(t, "hash", hash)
		assert.NoError(t, db.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		db, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer db.Close()

		repo := NewUserRepository(db)

		db.ExpectQuery("SELECT id, username, email, password_hash FROM users").
			WithArgs("test").
			WillReturnError(errors.New("database error"))

		user, hash, err := repo.GetByLogin(
			context.Background(),
			"test",
		)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, hash)
		assert.NoError(t, db.ExpectationsWereMet())
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer db.Close()

		repo := NewUserRepository(db)

		db.ExpectQuery("SELECT id, username, email FROM users").
			WithArgs(int64(1)).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{"id", "username", "email"},
				).AddRow(
					int64(1),
					"test",
					"test@example.com",
				),
			)

		user, err := repo.GetByID(
			context.Background(),
			1,
		)

		assert.NoError(t, err)
		assert.Equal(t, &models.UserResponse{
			ID:       1,
			Username: "test",
			Email:    "test@example.com",
		}, user)
		assert.NoError(t, db.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		db, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer db.Close()

		repo := NewUserRepository(db)

		db.ExpectQuery("SELECT id, username, email FROM users").
			WithArgs(int64(1)).
			WillReturnError(errors.New("database error"))

		user, err := repo.GetByID(
			context.Background(),
			1,
		)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.NoError(t, db.ExpectationsWereMet())
	})
}
