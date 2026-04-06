package service

import (
	"context"
	"fmt"
	"time"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/authz"
	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/carkeeper/backend/internal/validate"
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
	create.VIN = validate.NormalizeVIN(create.VIN)
	if msg := validate.VIN(create.VIN); msg != "" {
		return nil, apperr.BadRequest(msg)
	}
	y := time.Now().UTC().Year()
	if create.Year < 1900 || create.Year > y+1 {
		return nil, apperr.BadRequest("Invalid vehicle year")
	}
	if create.CurrentMileage < 0 {
		return nil, apperr.BadRequest("Mileage must be non-negative")
	}

	exists, err := s.repo.UserCar.VINExists(ctx, create.VIN)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if exists {
		return nil, apperr.Conflict("This VIN is already registered")
	}

	userCar, err := s.repo.UserCar.Create(ctx, userID, create)
	if err != nil {
		return nil, apperr.Internal(err)
	}

	userCarWithDetails, err := s.repo.UserCar.GetByID(ctx, userCar.UserCarID)
	if err != nil {
		return nil, err
	}

	return userCarWithDetails, nil
}

func (s *ProfileService) GetUserCar(ctx context.Context, userCarID uuid.UUID, requester uuid.UUID, role string) (*model.UserCarWithDetails, error) {
	car, err := s.repo.UserCar.GetByID(ctx, userCarID)
	if err != nil {
		return nil, err
	}
	if !authz.IsOwnerOrHasPermission(car.UserID, requester, role, authz.PermGarageViewAny) {
		return nil, fmt.Errorf("%w", apperr.ErrNotFound)
	}
	return car, nil
}

// DeleteUserCar removes the car if owned by userID.
func (s *ProfileService) DeleteUserCar(ctx context.Context, userID, userCarID uuid.UUID) error {
	return s.repo.UserCar.DeleteOwned(ctx, userID, userCarID)
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
