package service

import (
	"context"
	"fmt"

	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/google/uuid"
)

type ServiceService struct {
	repo *repository.Repository
}

func NewServiceService(repos *repository.Repository) *ServiceService {
	return &ServiceService{repo: repos}
}

func (s *ServiceService) GetServiceTypes(ctx context.Context, category *string, isAvailable *bool) ([]model.ServiceType, error) {
	return s.repo.ServiceType.GetAll(ctx, category, isAvailable)
}

func (s *ServiceService) GetBranches(ctx context.Context, isActive *bool) ([]model.Branch, error) {
	return s.repo.Branch.GetAll(ctx, isActive)
}

func (s *ServiceService) CreateAppointment(ctx context.Context, userID uuid.UUID, create model.ServiceAppointmentCreate) (*model.ServiceAppointmentWithDetails, error) {
	// Verify user car belongs to user
	userCar, err := s.repo.UserCar.GetByID(ctx, create.UserCarID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user car: %w", err)
	}

	if userCar.UserID != userID {
		return nil, fmt.Errorf("user car does not belong to user")
	}

	// Verify branch exists and is active
	branch, err := s.repo.Branch.GetByID(ctx, create.BranchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	if !branch.IsActive {
		return nil, fmt.Errorf("branch is not active")
	}

	// Verify service types exist and are available
	if len(create.ServiceTypeIDs) == 0 {
		return nil, fmt.Errorf("at least one service type is required")
	}

	serviceTypes, err := s.repo.ServiceType.GetAll(ctx, nil, boolPtr(true))
	if err != nil {
		return nil, fmt.Errorf("failed to get service types: %w", err)
	}

	serviceTypeMap := make(map[uuid.UUID]bool)
	for _, st := range serviceTypes {
		serviceTypeMap[st.ServiceTypeID] = true
	}

	for _, stID := range create.ServiceTypeIDs {
		if !serviceTypeMap[stID] {
			return nil, fmt.Errorf("service type %s not found or not available", stID)
		}
	}

	// Create appointment
	appointment, err := s.repo.ServiceAppointment.Create(ctx, create)
	if err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	// Get full details
	appointmentWithDetails, err := s.repo.ServiceAppointment.GetByID(ctx, appointment.ServiceAppointmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment details: %w", err)
	}

	return appointmentWithDetails, nil
}

func (s *ServiceService) GetAppointment(ctx context.Context, appointmentID uuid.UUID) (*model.ServiceAppointmentWithDetails, error) {
	return s.repo.ServiceAppointment.GetByID(ctx, appointmentID)
}

func (s *ServiceService) GetUserAppointments(ctx context.Context, userID uuid.UUID) ([]model.ServiceAppointmentWithDetails, error) {
	return s.repo.ServiceAppointment.GetByUserID(ctx, userID)
}

func (s *ServiceService) CancelAppointment(ctx context.Context, appointmentID uuid.UUID) error {
	return s.repo.ServiceAppointment.UpdateStatus(ctx, appointmentID, "cancelled")
}

func boolPtr(b bool) *bool {
	return &b
}

