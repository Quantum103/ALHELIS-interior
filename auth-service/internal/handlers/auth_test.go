package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-service/internal/middleware"
	"auth-service/internal/models"
	"auth-service/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(ctx context.Context, req models.RegisterRequest) (int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockAuthService) Login(ctx context.Context, req models.LoginRequest) (string, *models.UserResponse, error) {
	args := m.Called(ctx, req)

	var user *models.UserResponse
	if args.Get(1) != nil {
		user = args.Get(1).(*models.UserResponse)
	}

	return args.String(0), user, args.Error(2)
}

func (m *mockAuthService) GetMe(ctx context.Context, id int64) (*models.UserResponse, error) {
	args := m.Called(ctx, id)

	var user *models.UserResponse
	if args.Get(0) != nil {
		user = args.Get(0).(*models.UserResponse)
	}

	return user, args.Error(1)
}

func TestHandleRegister(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		mock       bool
		mockErr    error
		wantStatus int
	}{
		{
			name:       "success",
			method:     http.MethodPost,
			body:       `{"username":"test","email":"test@example.com","password":"password"}`,
			mock:       true,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "method not allowed",
			method:     http.MethodGet,
			body:       "",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "invalid json",
			method:     http.MethodPost,
			body:       `{invalid}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service error",
			method:     http.MethodPost,
			body:       `{"username":"test","email":"test@example.com","password":"password"}`,
			mock:       true,
			mockErr:    errors.New("database error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := new(mockAuthService)

			if tt.mock {
				auth.On(
					"Register",
					mock.Anything,
					mock.AnythingOfType("models.RegisterRequest"),
				).Return(int64(1), tt.mockErr)
			}

			server := NewServer(auth)

			req := httptest.NewRequest(
				tt.method,
				"/register",
				bytes.NewBufferString(tt.body),
			)
			rec := httptest.NewRecorder()

			server.HandleRegister(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.mock {
				auth.AssertExpectations(t)
			}
		})
	}
}

func TestHandleLogin(t *testing.T) {
	user := &models.UserResponse{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}

	tests := []struct {
		name       string
		method     string
		body       string
		token      string
		user       *models.UserResponse
		err        error
		wantStatus int
	}{
		{
			name:       "success",
			method:     http.MethodPost,
			body:       `{"username":"test","password":"password"}`,
			token:      "token",
			user:       user,
			wantStatus: http.StatusOK,
		},
		{
			name:       "method not allowed",
			method:     http.MethodGet,
			body:       "",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "invalid json",
			method:     http.MethodPost,
			body:       `{invalid}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid credentials",
			method:     http.MethodPost,
			body:       `{"username":"test","password":"wrong"}`,
			err:        service.ErrInvalidCredentials,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "service error",
			method:     http.MethodPost,
			body:       `{"username":"test","password":"password"}`,
			err:        errors.New("database error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := new(mockAuthService)

			if tt.method == http.MethodPost && tt.body != `{invalid}` {
				auth.On(
					"Login",
					mock.Anything,
					mock.AnythingOfType("models.LoginRequest"),
				).Return(tt.token, tt.user, tt.err)
			}

			server := NewServer(auth)

			req := httptest.NewRequest(
				tt.method,
				"/login",
				bytes.NewBufferString(tt.body),
			)
			rec := httptest.NewRecorder()

			server.HandleLogin(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.method == http.MethodPost && tt.body != `{invalid}` {
				auth.AssertExpectations(t)
			}
		})
	}
}

func TestHandleGetMe(t *testing.T) {
	user := &models.UserResponse{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}

	t.Run("success", func(t *testing.T) {
		auth := new(mockAuthService)

		auth.On("GetMe", mock.Anything, int64(1)).
			Return(user, nil)

		server := NewServer(auth)

		ctx := context.WithValue(
			context.Background(),
			middleware.UserIDKey,
			int64(1),
		)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		server.HandleGetMe(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.UserResponse
		err := json.NewDecoder(rec.Body).Decode(&response)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Username, response.Username)
		assert.Equal(t, user.Email, response.Email)

		auth.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		auth := new(mockAuthService)
		server := NewServer(auth)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()

		server.HandleGetMe(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		auth.AssertNotCalled(t, "GetMe")
	})

	t.Run("user not found", func(t *testing.T) {
		auth := new(mockAuthService)

		auth.On("GetMe", mock.Anything, int64(1)).
			Return(nil, errors.New("not found"))

		server := NewServer(auth)

		ctx := context.WithValue(
			context.Background(),
			middleware.UserIDKey,
			int64(1),
		)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		server.HandleGetMe(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		auth.AssertExpectations(t)
	})
}

func TestHandleAuthPage(t *testing.T) {
	t.Run("method not allowed", func(t *testing.T) {
		server := NewServer(new(mockAuthService))

		req := httptest.NewRequest(http.MethodPost, "/auth", nil)
		rec := httptest.NewRecorder()

		server.HandleAuthPage(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("get", func(t *testing.T) {
		server := NewServer(new(mockAuthService))

		req := httptest.NewRequest(http.MethodGet, "/auth", nil)
		rec := httptest.NewRecorder()

		server.HandleAuthPage(rec, req)

		assert.NotEqual(t, http.StatusMethodNotAllowed, rec.Code)
	})
}
