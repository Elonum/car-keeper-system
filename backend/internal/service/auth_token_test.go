package service

import (
	"testing"

	"github.com/carkeeper/backend/config"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/google/uuid"
)

func TestAuthService_ValidateToken_RoundTrip(t *testing.T) {
	cfg := config.TestConfig()
	svc := NewAuthService(&repository.Repository{}, cfg)
	userID := uuid.New()

	token, err := svc.generateToken(userID, "customer")
	if err != nil {
		t.Fatal(err)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != userID || claims.Role != "customer" {
		t.Fatalf("claims mismatch: %+v", claims)
	}
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	cfg := config.TestConfig()
	svc := NewAuthService(&repository.Repository{}, cfg)
	if _, err := svc.ValidateToken("not-a-jwt"); err == nil {
		t.Fatal("expected error for garbage token")
	}
}
