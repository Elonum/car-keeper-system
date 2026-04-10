package service

import (
	"context"
	"fmt"
	"strings"
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

	create.ServiceTypeIDs = dedupeUUIDs(create.ServiceTypeIDs)

	// Verify user car belongs to user
	userCar, err := s.repo.UserCar.GetByID(ctx, create.UserCarID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user car: %w", err)
	}

	if userCar.UserID != userID {
		return nil, apperr.Forbidden("user car does not belong to user")
	}

	// Verify branch exists and is active
	branch, err := s.repo.Branch.GetByID(ctx, create.BranchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	if !branch.IsActive {
		return nil, apperr.BadRequest("branch is not active")
	}

	if len(create.ServiceTypeIDs) == 0 {
		return nil, apperr.BadRequest("at least one service type is required")
	}

	selectedTypes, err := s.repo.ServiceType.GetByIDs(ctx, create.ServiceTypeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get service types: %w", err)
	}
	if len(selectedTypes) != len(create.ServiceTypeIDs) {
		return nil, apperr.BadRequest("one or more service types are not available")
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
	if !authz.IsOwnerOrHasPermission(a.OwnerUserID, requester, role, authz.PermAppointmentsViewAny) {
		return nil, fmt.Errorf("%w", apperr.ErrNotFound)
	}
	return a, nil
}

func (s *ServiceService) GetUserAppointments(ctx context.Context, userID uuid.UUID) ([]model.ServiceAppointmentWithDetails, error) {
	return s.repo.ServiceAppointment.GetByUserID(ctx, userID)
}

// ListAllAppointmentsForStaff returns all appointments (caller must enforce permission).
func (s *ServiceService) ListAllAppointmentsForStaff(ctx context.Context) ([]model.ServiceAppointmentWithDetails, error) {
	return s.repo.ServiceAppointment.ListAllWithDetails(ctx)
}

func (s *ServiceService) CancelAppointment(ctx context.Context, appointmentID uuid.UUID, requester uuid.UUID, role string) error {
	a, err := s.repo.ServiceAppointment.GetByID(ctx, appointmentID)
	if err != nil {
		return err
	}
	if !authz.IsOwnerOrHasPermission(a.OwnerUserID, requester, role, authz.PermAppointmentsViewAny) {
		return fmt.Errorf("%w", apperr.ErrForbidden)
	}
	switch a.Status {
	case "cancelled":
		return apperr.BadRequest("appointment already cancelled")
	case "completed":
		return apperr.BadRequest("cannot cancel completed appointment")
	}
	ok, err := s.repo.ServiceAppointment.UpdateStatusIfCurrent(ctx, appointmentID, "scheduled", "cancelled")
	if err != nil {
		return err
	}
	if !ok {
		return apperr.Conflict("appointment status changed, refresh and retry")
	}
	return nil
}

// RescheduleAppointment updates appointment_date for the car owner; services and duration stay unchanged.
func (s *ServiceService) RescheduleAppointment(ctx context.Context, userID uuid.UUID, appointmentID uuid.UUID, newDate time.Time) (*model.ServiceAppointmentWithDetails, error) {
	now := time.Now()
	if err := validate.AppointmentDate(newDate, now); err != nil {
		return nil, err
	}
	a, err := s.repo.ServiceAppointment.GetByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	if a.OwnerUserID != userID {
		return nil, fmt.Errorf("%w", apperr.ErrForbidden)
	}
	if a.Status != "scheduled" {
		return nil, apperr.BadRequest("only scheduled appointments can be rescheduled")
	}

	branch, err := s.repo.Branch.GetByID(ctx, a.BranchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}
	if !branch.IsActive {
		return nil, apperr.BadRequest("branch is not active")
	}
	if err := validateAppointmentSlot(branch, newDate, a.DurationMinutes); err != nil {
		return nil, err
	}

	if err := s.repo.ServiceAppointment.RescheduleOwned(ctx, appointmentID, userID, a.BranchID, newDate, a.DurationMinutes, branch.ConcurrentBays); err != nil {
		return nil, err
	}

	return s.repo.ServiceAppointment.GetByID(ctx, appointmentID)
}

func boolPtr(b bool) *bool {
	return &b
}

// AdminCreateServiceType creates a catalog service offering.
func (s *ServiceService) AdminCreateServiceType(ctx context.Context, name, category string, description *string, price float64, durationMinutes *int, isAvailable bool) (*model.ServiceType, error) {
	name = strings.TrimSpace(name)
	category = strings.TrimSpace(category)
	if name == "" || category == "" {
		return nil, apperr.BadRequest("name and category are required")
	}
	if price < 0 {
		return nil, apperr.BadRequest("invalid price")
	}
	return s.repo.ServiceType.Create(ctx, name, category, description, price, durationMinutes, isAvailable)
}

// AdminUpdateServiceType updates a service type.
func (s *ServiceService) AdminUpdateServiceType(ctx context.Context, id uuid.UUID, name, category string, description *string, price float64, durationMinutes *int, isAvailable bool) error {
	name = strings.TrimSpace(name)
	category = strings.TrimSpace(category)
	if name == "" || category == "" {
		return apperr.BadRequest("name and category are required")
	}
	return s.repo.ServiceType.Update(ctx, id, name, category, description, price, durationMinutes, isAvailable)
}

// AdminDeleteServiceType removes a service type.
func (s *ServiceService) AdminDeleteServiceType(ctx context.Context, id uuid.UUID) error {
	return s.repo.ServiceType.Delete(ctx, id)
}

// AdminUpdateBranch updates branch operational fields.
func (s *ServiceService) AdminUpdateBranch(ctx context.Context, id uuid.UUID, name, address *string, phone, email *string, isActive *bool) error {
	return s.repo.Branch.UpdateBranch(ctx, id, name, address, phone, email, isActive)
}
