package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/carkeeper/backend/config"
	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewAuthService(repos *repository.Repository, cfg *config.Config) *AuthService {
	return &AuthService{repo: repos, cfg: cfg}
}

func (s *AuthService) Register(ctx context.Context, create model.UserCreate) (*model.UserResponse, error) {
	// Check if email already exists
	exists, err := s.repo.User.EmailExists(ctx, create.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Create user
	user, err := s.repo.User.Create(ctx, create)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *AuthService) Login(ctx context.Context, login model.UserLogin) (string, *model.UserResponse, error) {
	// Verify credentials
	user, err := s.repo.User.VerifyPassword(ctx, login.Email, login.Password)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user.UserID, user.Role)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	response := user.ToResponse()
	return token, &response, nil
}

func (s *AuthService) GetUser(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error) {
	user, err := s.repo.User.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *AuthService) generateToken(userID uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(s.cfg.JWT.ExpiryDuration()).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}

func (s *AuthService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	role, _ := claims["role"].(string)

	return &TokenClaims{
		UserID: userID,
		Role:   role,
	}, nil
}

type TokenClaims struct {
	UserID uuid.UUID
	Role   string
}

