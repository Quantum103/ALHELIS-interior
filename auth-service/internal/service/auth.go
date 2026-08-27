package service

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("неверный логин или пароль")

type UserRepository interface {
	Create(ctx context.Context, username, email, passwordHash string) error
	GetByLogin(ctx context.Context, login string) (*models.UserResponse, string, error)
	GetByID(ctx context.Context, id int64) (*models.UserResponse, error)
}

type AuthService struct {
	repo      UserRepository
	jwtSecret []byte
}

func NewAuthService(repo UserRepository, secret string) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: []byte(secret),
	}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.Create(ctx, req.Username, req.Email, string(hashedPassword))
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (string, *models.UserResponse, error) {
	user, hash, err := s.repo.GetByLogin(ctx, req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

func (s *AuthService) GetMe(ctx context.Context, id int64) (*models.UserResponse, error) {
	return s.repo.GetByID(ctx, id)
}
