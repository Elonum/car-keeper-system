package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/carkeeper/backend/config"
	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/carkeeper/backend/internal/validate"
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

func (s *AuthService) Register(ctx context.Context, in model.UserRegisterInput) (*model.UserResponse, error) {
	create := model.UserCreate{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		Phone:     in.Phone,
		Password:  in.Password,
		Role:      "customer",
	}
	// Check if email already exists
	exists, err := s.repo.User.EmailExists(ctx, create.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperr.Conflict("This email is already registered")
	}

	// Create user
	user, err := s.repo.User.Create(ctx, create)
	if err != nil {
		return nil, apperr.Internal(err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *AuthService) Login(ctx context.Context, login model.UserLogin) (string, *model.UserResponse, error) {
	// Verify credentials
	user, err := s.repo.User.VerifyPassword(ctx, login.Email, login.Password)
	if err != nil {
		if errors.Is(err, apperr.ErrInvalidCredentials) {
			return "", nil, apperr.Unauthorized("Invalid email or password")
		}
		return "", nil, err
	}

	// Generate JWT token
	token, err := s.generateToken(user.UserID, user.Role)
	if err != nil {
		return "", nil, apperr.Internal(err)
	}

	response := user.ToResponse()
	return token, &response, nil
}

func (s *AuthService) GetUser(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error) {
	user, err := s.repo.User.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, apperr.NotFoundErr("User not found")
		}
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// ChangePassword requires the current password and validates the new password strength.
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	if msg := validate.NewPassword(newPassword); msg != "" {
		return apperr.BadRequest(msg)
	}
	if err := s.repo.User.VerifyPasswordForUserID(ctx, userID, currentPassword); err != nil {
		if errors.Is(err, apperr.ErrInvalidCredentials) {
			return apperr.Unauthorized("Current password is incorrect")
		}
		return err
	}
	if msg := validate.NewPasswordMustDifferFromCurrent(currentPassword, newPassword); msg != "" {
		return apperr.BadRequest(msg)
	}
	if err := s.repo.User.UpdatePassword(ctx, userID, newPassword); err != nil {
		return apperr.Internal(err)
	}
	return nil
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

