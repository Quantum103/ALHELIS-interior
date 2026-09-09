package repo

import (
	"context"
	"errors"
	"testing"

	"user-service/internal/models"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_GetProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer db.Close()

		repository := NewUserRepository(db)

		db.ExpectQuery(`SELECT user_id, name, phone FROM user_profiles WHERE user_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(
				pgxmock.NewRows([]string{"user_id", "name", "phone"}).
					AddRow(int64(1), "Pavel", "123456789"),
			)

		profile, err := repository.GetProfile(
			context.Background(),
			1,
		)

		assert.NoError(t, err)
		assert.Equal(t, &models.UserProfile{
			UserID: 1,
			Name:   "Pavel",
			Phone:  "123456789",
		}, profile)

		assert.NoError(t, db.ExpectationsWereMet())
	})

	t.Run("profile not found", func(t *testing.T) {
		db, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer db.Close()

		repository := NewUserRepository(db)

		db.ExpectQuery(`SELECT user_id, name, phone FROM user_profiles WHERE user_id = \$1`).
			WithArgs(int64(1)).
			WillReturnError(errors.New("profile not found"))

		profile, err := repository.GetProfile(
			context.Background(),
			1,
		)

		assert.Error(t, err)
		assert.Nil(t, profile)

		assert.NoError(t, db.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer db.Close()

		repository := NewUserRepository(db)

		db.ExpectQuery(`SELECT user_id, name, phone FROM user_profiles WHERE user_id = \$1`).
			WithArgs(int64(1)).
			WillReturnError(errors.New("database error"))

		profile, err := repository.GetProfile(
			context.Background(),
			1,
		)

		assert.Error(t, err)
		assert.Nil(t, profile)
		assert.Contains(t, err.Error(), "get profile")

		assert.NoError(t, db.ExpectationsWereMet())
	})
}

func TestUserRepository_CreateProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer db.Close()

		repository := NewUserRepository(db)

		db.ExpectExec(`INSERT INTO user_profiles`).
			WithArgs(int64(1)).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = repository.CreateProfile(
			context.Background(),
			1,
		)

		assert.NoError(t, err)
		assert.NoError(t, db.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer db.Close()

		repository := NewUserRepository(db)

		db.ExpectExec(`INSERT INTO user_profiles`).
			WithArgs(int64(1)).
			WillReturnError(errors.New("database error"))

		err = repository.CreateProfile(
			context.Background(),
			1,
		)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")

		assert.NoError(t, db.ExpectationsWereMet())
	})
}
