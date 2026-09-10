package service

import (
	"context"
	"errors"
	"testing"

	"user-service/internal/models"
)

type mockUserRepository struct {
	getProfileFunc    func(ctx context.Context, userID int64) (*models.UserProfile, error)
	createProfileFunc func(ctx context.Context, userID int64, name string, phone string) error
}

func (m *mockUserRepository) GetProfile(ctx context.Context, userID int64) (*models.UserProfile, error) {
	return m.getProfileFunc(ctx, userID)
}

func (m *mockUserRepository) CreateProfile(ctx context.Context, userID int64, name string, phone string) error {
	return m.createProfileFunc(ctx, userID, name, phone)
}

func TestUserService_GetProfile_Success(t *testing.T) {
	expectedProfile := &models.UserProfile{UserID: 1, Name: "Test", Phone: "123"}
	mock := &mockUserRepository{
		getProfileFunc: func(ctx context.Context, userID int64) (*models.UserProfile, error) {
			return expectedProfile, nil
		},
	}

	svc := NewUserService(mock)
	profile, err := svc.GetProfile(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if profile.UserID != expectedProfile.UserID {
		t.Errorf("expected UserID %d, got %d", expectedProfile.UserID, profile.UserID)
	}
}

func TestUserService_GetProfile_Error(t *testing.T) {
	expectedErr := errors.New("profile not found")
	mock := &mockUserRepository{
		getProfileFunc: func(ctx context.Context, userID int64) (*models.UserProfile, error) {
			return nil, expectedErr
		},
	}

	svc := NewUserService(mock)
	_, err := svc.GetProfile(context.Background(), 1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestUserService_CreateProfile_Success(t *testing.T) {
	mock := &mockUserRepository{
		createProfileFunc: func(ctx context.Context, userID int64, name string, phone string) error {
			return nil
		},
	}

	svc := NewUserService(mock)
	err := svc.CreateProfile(context.Background(), 1, "Test", "123")

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestUserService_CreateProfile_Error(t *testing.T) {
	expectedErr := errors.New("db error")
	mock := &mockUserRepository{
		createProfileFunc: func(ctx context.Context, userID int64, name string, phone string) error {
			return expectedErr
		},
	}

	svc := NewUserService(mock)
	err := svc.CreateProfile(context.Background(), 1, "Test", "123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "failed to create profile in db: db error" {
		t.Errorf("expected 'failed to create profile in db: db error', got %v", err)
	}
}
