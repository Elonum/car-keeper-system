package service

import (
	"context"
	"fmt"

	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/google/uuid"
)

type ProfileService struct {
	repo *repository.Repository
}

func NewProfileService(repos *repository.Repository) *ProfileService {
	return &ProfileService{repo: repos}
}

func (s *ProfileService) GetUserCars(ctx context.Context, userID uuid.UUID) ([]model.UserCarWithDetails, error) {
	return s.repo.UserCar.GetByUserID(ctx, userID)
}

func (s *ProfileService) CreateUserCar(ctx context.Context, userID uuid.UUID, create model.UserCarCreate) (*model.UserCarWithDetails, error) {
	exists, err := s.repo.UserCar.VINExists(ctx, create.VIN)
	if err != nil {
		return nil, fmt.Errorf("failed to check VIN: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("VIN already exists")
	}

	userCar, err := s.repo.UserCar.Create(ctx, userID, create)
	if err != nil {
		return nil, fmt.Errorf("failed to create user car: %w", err)
	}

	userCarWithDetails, err := s.repo.UserCar.GetByID(ctx, userCar.UserCarID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user car details: %w", err)
	}

	return userCarWithDetails, nil
}

func (s *ProfileService) GetUserCar(ctx context.Context, userCarID uuid.UUID) (*model.UserCarWithDetails, error) {
	return s.repo.UserCar.GetByID(ctx, userCarID)
}

