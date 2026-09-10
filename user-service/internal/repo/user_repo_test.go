package repo

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockDB struct {
	queryRowFunc func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	execFunc     func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

func (m *mockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return m.queryRowFunc(ctx, sql, args...)
}

func (m *mockDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return m.execFunc(ctx, sql, args...)
}

type mockRow struct {
	scanFunc func(dest ...interface{}) error
}

func (m *mockRow) Scan(dest ...interface{}) error {
	return m.scanFunc(dest...)
}

func TestGetProfile_Success(t *testing.T) {
	mock := &mockDB{
		queryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
			return &mockRow{
				scanFunc: func(dest ...interface{}) error {
					if len(dest) >= 3 {
						if ptr, ok := dest[0].(*int64); ok {
							*ptr = 1
						}
						if ptr, ok := dest[1].(*string); ok {
							*ptr = "TestName"
						}
						if ptr, ok := dest[2].(*string); ok {
							*ptr = "123456"
						}
					}
					return nil
				},
			}
		},
	}

	repo := NewUserRepository(mock)
	profile, err := repo.GetProfile(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if profile.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", profile.UserID)
	}
	if profile.Name != "TestName" {
		t.Errorf("expected Name 'TestName', got %s", profile.Name)
	}
	if profile.Phone != "123456" {
		t.Errorf("expected Phone '123456', got %s", profile.Phone)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	mock := &mockDB{
		queryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
			return &mockRow{
				scanFunc: func(dest ...interface{}) error {
					return pgx.ErrNoRows
				},
			}
		},
	}

	repo := NewUserRepository(mock)
	_, err := repo.GetProfile(context.Background(), 1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "profile not found" {
		t.Errorf("expected 'profile not found', got %v", err)
	}
}

func TestGetProfile_Error(t *testing.T) {
	expectedErr := errors.New("db error")
	mock := &mockDB{
		queryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
			return &mockRow{
				scanFunc: func(dest ...interface{}) error {
					return expectedErr
				},
			}
		},
	}

	repo := NewUserRepository(mock)
	_, err := repo.GetProfile(context.Background(), 1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "get profile: db error" {
		t.Errorf("expected 'get profile: db error', got %v", err)
	}
}
