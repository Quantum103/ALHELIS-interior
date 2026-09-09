package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"user-service/internal/middleware"
	"user-service/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) GetProfile(ctx context.Context, userID int64) (*models.UserProfile, error) {
	args := m.Called(ctx, userID)

	var profile *models.UserProfile
	if args.Get(0) != nil {
		profile = args.Get(0).(*models.UserProfile)
	}

	return profile, args.Error(1)
}

func (m *mockUserService) CreateProfile(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestShowLK(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(mockUserService)

		profile := &models.UserProfile{
			UserID: 1,
			Name:   "Test",
			Phone:  "123456789",
		}

		service.On("GetProfile", mock.Anything, int64(1)).
			Return(profile, nil)

		handler := NewUserHandler(service)

		ctx := context.WithValue(
			context.Background(),
			middleware.UserIDKey,
			int64(1),
		)

		req := httptest.NewRequest(http.MethodGet, "/lk", nil)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		handler.ShowLK(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var result models.UserProfile

		err := json.NewDecoder(rec.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, profile.UserID, result.UserID)
		assert.Equal(t, profile.Name, result.Name)
		assert.Equal(t, profile.Phone, result.Phone)

		service.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		service := new(mockUserService)

		handler := NewUserHandler(service)

		req := httptest.NewRequest(http.MethodGet, "/lk", nil)
		rec := httptest.NewRecorder()

		handler.ShowLK(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		service.AssertNotCalled(t, "GetProfile")
	})

	t.Run("profile not found", func(t *testing.T) {
		service := new(mockUserService)

		service.On("GetProfile", mock.Anything, int64(1)).
			Return(nil, errors.New("profile not found"))

		handler := NewUserHandler(service)

		ctx := context.WithValue(
			context.Background(),
			middleware.UserIDKey,
			int64(1),
		)

		req := httptest.NewRequest(http.MethodGet, "/lk", nil)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		handler.ShowLK(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		service.AssertExpectations(t)
	})
}

func TestCreateProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(mockUserService)

		service.On("CreateProfile", mock.Anything, int64(1)).
			Return(nil)

		handler := NewUserHandler(service)

		ctx := context.WithValue(
			context.Background(),
			middleware.UserIDKey,
			int64(1),
		)

		req := httptest.NewRequest(http.MethodPost, "/profile", nil)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		handler.CreateProfile(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var result map[string]string

		err := json.NewDecoder(rec.Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, "Профиль создан", result["message"])

		service.AssertExpectations(t)
	})

	t.Run("method not allowed", func(t *testing.T) {
		service := new(mockUserService)

		handler := NewUserHandler(service)

		req := httptest.NewRequest(http.MethodGet, "/profile", nil)
		rec := httptest.NewRecorder()

		handler.CreateProfile(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

		service.AssertNotCalled(t, "CreateProfile")
	})

	t.Run("unauthorized", func(t *testing.T) {
		service := new(mockUserService)

		handler := NewUserHandler(service)

		req := httptest.NewRequest(http.MethodPost, "/profile", nil)
		rec := httptest.NewRecorder()

		handler.CreateProfile(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		service.AssertNotCalled(t, "CreateProfile")
	})

	t.Run("service error", func(t *testing.T) {
		service := new(mockUserService)

		service.On("CreateProfile", mock.Anything, int64(1)).
			Return(errors.New("database error"))

		handler := NewUserHandler(service)

		ctx := context.WithValue(
			context.Background(),
			middleware.UserIDKey,
			int64(1),
		)

		req := httptest.NewRequest(http.MethodPost, "/profile", nil)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		handler.CreateProfile(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		service.AssertExpectations(t)
	})
}
