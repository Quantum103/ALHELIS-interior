package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"auth-service/internal/middleware"
	"auth-service/internal/models"
	"auth-service/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(
	ctx context.Context,
	req models.RegisterRequest,
) (int64, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuthService) Login(
	ctx context.Context,
	req models.LoginRequest,
) (string, *models.UserResponse, error) {
	args := m.Called(ctx, req)

	var user *models.UserResponse
	if args.Get(1) != nil {
		user = args.Get(1).(*models.UserResponse)
	}

	return args.String(0), user, args.Error(2)
}

func (m *MockAuthService) GetMe(
	ctx context.Context,
	id int64,
) (*models.UserResponse, error) {
	args := m.Called(ctx, id)

	var user *models.UserResponse
	if args.Get(0) != nil {
		user = args.Get(0).(*models.UserResponse)
	}

	return user, args.Error(1)
}

func TestHandleRegister(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		setupMock      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{
				"username": "pavel",
				"email": "pavel@test.com",
				"password": "123456"
			}`,
			setupMock: func(m *MockAuthService) {
				m.On(
					"Register",
					mock.Anything,
					models.RegisterRequest{
						Username: "pavel",
						Email:    "pavel@test.com",
						Password: "123456",
					},
				).Return(int64(1), nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			body:           `{invalid json}`,
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: `{
				"username": "pavel",
				"email": "pavel@test.com",
				"password": "123456"
			}`,
			setupMock: func(m *MockAuthService) {
				m.On(
					"Register",
					mock.Anything,
					models.RegisterRequest{
						Username: "pavel",
						Email:    "pavel@test.com",
						Password: "123456",
					},
				).Return(int64(0), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "method not allowed",
			body:           `{}`,
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			server := NewServer(mockService)

			method := http.MethodPost
			if tt.name == "method not allowed" {
				method = http.MethodGet
			}

			req := httptest.NewRequest(
				method,
				"/auth/register",
				strings.NewReader(tt.body),
			)

			rec := httptest.NewRecorder()

			server.HandleRegister(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandleLogin(t *testing.T) {
	user := &models.UserResponse{
		ID:       1,
		Username: "pavel",
		Email:    "pavel@test.com",
	}

	tests := []struct {
		name           string
		body           string
		setupMock      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{
				"username": "pavel",
				"password": "123456"
			}`,
			setupMock: func(m *MockAuthService) {
				m.On(
					"Login",
					mock.Anything,
					models.LoginRequest{
						Username: "pavel",
						Password: "123456",
					},
				).Return("test-token", user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid json",
			body:           `{invalid json}`,
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid credentials",
			body: `{
				"username": "pavel",
				"password": "wrong"
			}`,
			setupMock: func(m *MockAuthService) {
				m.On(
					"Login",
					mock.Anything,
					models.LoginRequest{
						Username: "pavel",
						Password: "wrong",
					},
				).Return(
					"",
					nil,
					service.ErrInvalidCredentials,
				)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "service error",
			body: `{
				"username": "pavel",
				"password": "123456"
			}`,
			setupMock: func(m *MockAuthService) {
				m.On(
					"Login",
					mock.Anything,
					models.LoginRequest{
						Username: "pavel",
						Password: "123456",
					},
				).Return(
					"",
					nil,
					errors.New("database error"),
				)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "method not allowed",
			body:           `{}`,
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			server := NewServer(mockService)

			method := http.MethodPost
			if tt.name == "method not allowed" {
				method = http.MethodGet
			}

			req := httptest.NewRequest(
				method,
				"/auth/login",
				strings.NewReader(tt.body),
			)

			rec := httptest.NewRecorder()

			server.HandleLogin(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.name == "success" {
				var response map[string]interface{}

				err := json.NewDecoder(rec.Body).Decode(&response)

				assert.NoError(t, err)
				assert.Equal(t, "test-token", response["access_token"])
				assert.NotNil(t, response["user"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandleGetMe(t *testing.T) {
	user := &models.UserResponse{
		ID:       1,
		Username: "pavel",
		Email:    "pavel@test.com",
	}

	t.Run("success", func(t *testing.T) {
		mockService := new(MockAuthService)

		mockService.On(
			"GetMe",
			mock.Anything,
			int64(1),
		).Return(user, nil)

		server := NewServer(mockService)

		req := httptest.NewRequest(
			http.MethodGet,
			"/auth/me",
			nil,
		)

		ctx := context.WithValue(
			req.Context(),
			middleware.UserIDKey,
			int64(1),
		)

		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		server.HandleGetMe(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.UserResponse

		err := json.NewDecoder(rec.Body).Decode(&response)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), response.ID)
		assert.Equal(t, "pavel", response.Username)
		assert.Equal(t, "pavel@test.com", response.Email)

		mockService.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockService := new(MockAuthService)

		server := NewServer(mockService)

		req := httptest.NewRequest(
			http.MethodGet,
			"/auth/me",
			nil,
		)

		rec := httptest.NewRecorder()

		server.HandleGetMe(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		mockService.AssertNotCalled(t, "GetMe")
	})

	t.Run("user not found", func(t *testing.T) {
		mockService := new(MockAuthService)

		mockService.On(
			"GetMe",
			mock.Anything,
			int64(1),
		).Return(nil, errors.New("user not found"))

		server := NewServer(mockService)

		req := httptest.NewRequest(
			http.MethodGet,
			"/auth/me",
			nil,
		)

		ctx := context.WithValue(
			req.Context(),
			middleware.UserIDKey,
			int64(1),
		)

		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		server.HandleGetMe(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		mockService.AssertExpectations(t)
	})
}
