package service

import (
	"context"
	"fmt"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/authz"
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

func (s *ProfileService) GetUserCar(ctx context.Context, userCarID uuid.UUID, requester uuid.UUID, role string) (*model.UserCarWithDetails, error) {
	car, err := s.repo.UserCar.GetByID(ctx, userCarID)
	if err != nil {
		return nil, err
	}
	if !authz.IsOwnerOrStaff(car.UserID, requester, role) {
		return nil, fmt.Errorf("%w", apperr.ErrNotFound)
	}
	return car, nil
}

// UpdateProfile updates first name, last name, and phone for the authenticated user.
func (s *ProfileService) UpdateProfile(ctx context.Context, userID uuid.UUID, firstName, lastName string, phone *string) (*model.UserResponse, error) {
	user, err := s.repo.User.UpdateProfile(ctx, userID, firstName, lastName, phone)
	if err != nil {
		return nil, err
	}
	resp := user.ToResponse()
	return &resp, nil
}

