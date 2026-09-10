package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"user-service/internal/middleware"
	"user-service/internal/models"
)

type mockUserService struct {
	getProfileFunc    func(ctx context.Context, userID int64) (*models.UserProfile, error)
	createProfileFunc func(ctx context.Context, userID int64, name string, phone string) error
}

func (m *mockUserService) GetProfile(ctx context.Context, userID int64) (*models.UserProfile, error) {
	return m.getProfileFunc(ctx, userID)
}

func (m *mockUserService) CreateProfile(ctx context.Context, userID int64, name string, phone string) error {
	return m.createProfileFunc(ctx, userID, name, phone)
}

func contextWithUserID(userID int64) context.Context {
	return context.WithValue(context.Background(), middleware.UserIDKey, userID)
}

func TestShowLK_Unauthorized(t *testing.T) {
	handler := NewUserHandler(&mockUserService{})
	req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
	w := httptest.NewRecorder()

	handler.ShowLK(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestShowLK_NotFound(t *testing.T) {
	mock := &mockUserService{
		getProfileFunc: func(ctx context.Context, userID int64) (*models.UserProfile, error) {
			return nil, errors.New("not found")
		},
	}
	handler := NewUserHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil).WithContext(contextWithUserID(1))
	w := httptest.NewRecorder()

	handler.ShowLK(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestShowLK_Success(t *testing.T) {
	expectedProfile := &models.UserProfile{UserID: 1, Name: "Test"}
	mock := &mockUserService{
		getProfileFunc: func(ctx context.Context, userID int64) (*models.UserProfile, error) {
			return expectedProfile, nil
		},
	}
	handler := NewUserHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil).WithContext(contextWithUserID(1))
	w := httptest.NewRecorder()

	handler.ShowLK(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCreateProfile_Unauthorized(t *testing.T) {
	handler := NewUserHandler(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/api/user/profile", nil)
	w := httptest.NewRecorder()

	handler.CreateProfile(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCreateProfile_InternalError(t *testing.T) {
	mock := &mockUserService{
		createProfileFunc: func(ctx context.Context, userID int64, name string, phone string) error {
			return errors.New("db error")
		},
	}
	handler := NewUserHandler(mock)
	req := httptest.NewRequest(http.MethodPost, "/api/user/profile", nil).WithContext(contextWithUserID(1))
	w := httptest.NewRecorder()

	handler.CreateProfile(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestCreateProfile_Success(t *testing.T) {
	mock := &mockUserService{
		createProfileFunc: func(ctx context.Context, userID int64, name string, phone string) error {
			return nil
		},
	}
	handler := NewUserHandler(mock)
	req := httptest.NewRequest(http.MethodPost, "/api/user/profile", nil).WithContext(contextWithUserID(1))
	w := httptest.NewRecorder()

	handler.CreateProfile(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}
