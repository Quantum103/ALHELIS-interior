package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	pb "amelli/proto"
	"auth-service/internal/middleware"
	"auth-service/internal/models"
	"auth-service/internal/service"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, req models.RegisterRequest) error
	Login(ctx context.Context, req models.LoginRequest) (string, *models.UserResponse, error)
	GetMe(ctx context.Context, id int64) (*models.UserResponse, error)
}

type Server struct {
	pb.UnimplementedAlhelisServiceServer
	authService AuthServiceInterface
}

func NewServer(authService AuthServiceInterface) *Server {
	return &Server{authService: authService}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := s.authService.Register(ctx, req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка при создании пользователя")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Аккаунт успешно создан"})
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tokenString, user, err := s.authService.Login(ctx, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": tokenString,
		"user":         user,
	})
}

func (s *Server) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Неавторизован")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := s.authService.GetMe(ctx, userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (s *Server) HandleAuthPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "/app/frontend/auth.html")
}
