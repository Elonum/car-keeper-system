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
	now := time.Now()
	if err := validate.AppointmentDate(create.AppointmentDate, now); err != nil {
		return nil, err
	}
	if err := validate.AppointmentDescription(create.Description); err != nil {
		return nil, err
	}

	// Deduplicate service type IDs (client may send duplicates; DB PK would reject second insert)
	if len(create.ServiceTypeIDs) > 0 {
		seen := make(map[uuid.UUID]struct{}, len(create.ServiceTypeIDs))
		uniq := make([]uuid.UUID, 0, len(create.ServiceTypeIDs))
		for _, id := range create.ServiceTypeIDs {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			uniq = append(uniq, id)
		}
		create.ServiceTypeIDs = uniq
	}

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

	if len(create.ServiceTypeIDs) == 0 {
		return nil, fmt.Errorf("at least one service type is required")
	}

	selectedTypes, err := s.repo.ServiceType.GetByIDs(ctx, create.ServiceTypeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get service types: %w", err)
	}
	if len(selectedTypes) != len(create.ServiceTypeIDs) {
		return nil, fmt.Errorf("one or more service types are not available")
	}

	create.DurationMinutes = totalDurationMinutes(selectedTypes)
	if err := validateAppointmentSlot(branch, create.AppointmentDate, create.DurationMinutes); err != nil {
		return nil, err
	}

	appointment, err := s.repo.ServiceAppointment.Create(ctx, create, branch.ConcurrentBays)
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

func (s *ServiceService) GetAppointment(ctx context.Context, appointmentID uuid.UUID, requester uuid.UUID, role string) (*model.ServiceAppointmentWithDetails, error) {
	a, err := s.repo.ServiceAppointment.GetByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	if !authz.IsOwnerOrStaff(a.OwnerUserID, requester, role) {
		return nil, fmt.Errorf("%w", apperr.ErrNotFound)
	}
	return a, nil
}

func (s *ServiceService) GetUserAppointments(ctx context.Context, userID uuid.UUID) ([]model.ServiceAppointmentWithDetails, error) {
	return s.repo.ServiceAppointment.GetByUserID(ctx, userID)
}

func (s *ServiceService) CancelAppointment(ctx context.Context, appointmentID uuid.UUID, requester uuid.UUID, role string) error {
	a, err := s.repo.ServiceAppointment.GetByID(ctx, appointmentID)
	if err != nil {
		return err
	}
	if !authz.IsOwnerOrStaff(a.OwnerUserID, requester, role) {
		return fmt.Errorf("%w", apperr.ErrForbidden)
	}
	switch a.Status {
	case "cancelled":
		return fmt.Errorf("appointment already cancelled")
	case "completed":
		return fmt.Errorf("cannot cancel completed appointment")
	}
	return s.repo.ServiceAppointment.UpdateStatus(ctx, appointmentID, "cancelled")
}

func boolPtr(b bool) *bool {
	return &b
}

