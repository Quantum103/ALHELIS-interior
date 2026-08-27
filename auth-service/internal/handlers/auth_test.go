package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-service/internal/middleware"
	"auth-service/internal/models"
)

type mockAuthService struct {
	mockRegister func(ctx context.Context, req models.RegisterRequest) error
	mockLogin    func(ctx context.Context, req models.LoginRequest) (string, *models.UserResponse, error)
	mockGetMe    func(ctx context.Context, id int64) (*models.UserResponse, error)
}

func (m *mockAuthService) Register(ctx context.Context, req models.RegisterRequest) error {
	return m.mockRegister(ctx, req)
}
func (m *mockAuthService) Login(ctx context.Context, req models.LoginRequest) (string, *models.UserResponse, error) {
	return m.mockLogin(ctx, req)
}
func (m *mockAuthService) GetMe(ctx context.Context, id int64) (*models.UserResponse, error) {
	return m.mockGetMe(ctx, id)
}

func TestServer_HandleRegister(t *testing.T) {
	mockSvc := &mockAuthService{
		mockRegister: func(ctx context.Context, req models.RegisterRequest) error {
			return nil
		},
	}
	server := NewServer(mockSvc)

	body := []byte(`{"username":"test","email":"test@test.com","password":"123"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	server.HandleRegister(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Ожидался статус 201, получен %d", rr.Code)
	}
}

func TestServer_HandleLogin(t *testing.T) {
	mockSvc := &mockAuthService{
		mockLogin: func(ctx context.Context, req models.LoginRequest) (string, *models.UserResponse, error) {
			return "mock_token", &models.UserResponse{ID: 1}, nil
		},
	}
	server := NewServer(mockSvc)

	body := []byte(`{"username":"test","password":"123"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	server.HandleLogin(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", rr.Code)
	}
}

func TestServer_HandleGetMe(t *testing.T) {
	mockSvc := &mockAuthService{
		mockGetMe: func(ctx context.Context, id int64) (*models.UserResponse, error) {
			return &models.UserResponse{ID: id, Username: "test"}, nil
		},
	}
	server := NewServer(mockSvc)

	req, _ := http.NewRequest(http.MethodGet, "/api/auth/me", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, int64(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	server.HandleGetMe(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", rr.Code)
	}

	reqUnauth, _ := http.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rrUnauth := httptest.NewRecorder()
	server.HandleGetMe(rrUnauth, reqUnauth)

	if rrUnauth.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус 401, получен %d", rrUnauth.Code)
	}
}
