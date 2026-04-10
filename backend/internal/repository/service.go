package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ServiceTypeRepository struct {
	db *database.DB
}

func NewServiceTypeRepository(db *database.DB) *ServiceTypeRepository {
	return &ServiceTypeRepository{db: db}
}

func (r *ServiceTypeRepository) GetAll(ctx context.Context, category *string, isAvailable *bool) ([]model.ServiceType, error) {
	query := `SELECT service_type_id, name, category, description, price, duration_minutes, is_available, created_at FROM service_types WHERE 1=1`
	var args []interface{}
	argPos := 1

	if category != nil {
		query += fmt.Sprintf(` AND category = $%d`, argPos)
		args = append(args, *category)
		argPos++
	}

	if isAvailable != nil {
		query += fmt.Sprintf(` AND is_available = $%d`, argPos)
		args = append(args, *isAvailable)
		argPos++
	}

	query += ` ORDER BY category, name`

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get service types: %w", err)
	}
	defer rows.Close()

	var serviceTypes []model.ServiceType
	for rows.Next() {
		var st model.ServiceType
		if err := rows.Scan(&st.ServiceTypeID, &st.Name, &st.Category, &st.Description, &st.Price, &st.DurationMinutes, &st.IsAvailable, &st.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan service type: %w", err)
		}
		serviceTypes = append(serviceTypes, st)
	}

	return serviceTypes, nil
}

func (r *ServiceTypeRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.ServiceType, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := `
		SELECT service_type_id, name, category, description, price, duration_minutes, is_available, created_at
		FROM service_types WHERE service_type_id = ANY($1) AND is_available = true
	`
	rows, err := r.db.Pool.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get service types by ids: %w", err)
	}
	defer rows.Close()

	var out []model.ServiceType
	for rows.Next() {
		var st model.ServiceType
		if err := rows.Scan(&st.ServiceTypeID, &st.Name, &st.Category, &st.Description, &st.Price, &st.DurationMinutes, &st.IsAvailable, &st.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan service type: %w", err)
		}
		out = append(out, st)
	}
	return out, nil
}

// Create inserts a service type row.
func (r *ServiceTypeRepository) Create(ctx context.Context, name, category string, description *string, price float64, durationMinutes *int, isAvailable bool) (*model.ServiceType, error) {
	var st model.ServiceType
	err := r.db.Pool.QueryRow(ctx, `
		INSERT INTO service_types (name, category, description, price, duration_minutes, is_available)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING service_type_id, name, category, description, price, duration_minutes, is_available, created_at
	`, name, category, description, price, durationMinutes, isAvailable).Scan(
		&st.ServiceTypeID, &st.Name, &st.Category, &st.Description, &st.Price, &st.DurationMinutes, &st.IsAvailable, &st.CreatedAt,
	)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return &st, nil
}

// Update patches a service type.
func (r *ServiceTypeRepository) Update(ctx context.Context, id uuid.UUID, name, category string, description *string, price float64, durationMinutes *int, isAvailable bool) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE service_types
		SET name = $1, category = $2, description = $3, price = $4, duration_minutes = $5, is_available = $6
		WHERE service_type_id = $7
	`, name, category, description, price, durationMinutes, isAvailable, id)
	if err != nil {
		return apperr.Internal(err)
	}
	return nil
}

// Delete removes a service type (fails if referenced by appointments).
func (r *ServiceTypeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM service_types WHERE service_type_id = $1`, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return apperr.BadRequest("Cannot delete: service type is in use")
		}
		return apperr.Internal(err)
	}
	return nil
}

type BranchRepository struct {
	db *database.DB
}

func NewBranchRepository(db *database.DB) *BranchRepository {
	return &BranchRepository{db: db}
}

func (r *BranchRepository) GetAll(ctx context.Context, isActive *bool) ([]model.Branch, error) {
	query := `SELECT branch_id, name, address, phone, email, is_active, timezone, workday_start_minutes, workday_end_minutes, slot_step_minutes, concurrent_bays, created_at, updated_at FROM branches`
	var args []interface{}

	if isActive != nil {
		query += ` WHERE is_active = $1`
		args = append(args, *isActive)
	}
	query += ` ORDER BY name`

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}
	defer rows.Close()

	var branches []model.Branch
	for rows.Next() {
		var branch model.Branch
		if err := rows.Scan(&branch.BranchID, &branch.Name, &branch.Address, &branch.Phone, &branch.Email, &branch.IsActive,
			&branch.Timezone, &branch.WorkdayStartMinutes, &branch.WorkdayEndMinutes, &branch.SlotStepMinutes, &branch.ConcurrentBays,
			&branch.CreatedAt, &branch.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan branch: %w", err)
		}
		branches = append(branches, branch)
	}

	return branches, nil
}

func (r *BranchRepository) GetByID(ctx context.Context, branchID uuid.UUID) (*model.Branch, error) {
	var branch model.Branch
	query := `SELECT branch_id, name, address, phone, email, is_active, timezone, workday_start_minutes, workday_end_minutes, slot_step_minutes, concurrent_bays, created_at, updated_at FROM branches WHERE branch_id = $1`

	err := r.db.Pool.QueryRow(ctx, query, branchID).Scan(
		&branch.BranchID, &branch.Name, &branch.Address, &branch.Phone, &branch.Email,
		&branch.IsActive, &branch.Timezone, &branch.WorkdayStartMinutes, &branch.WorkdayEndMinutes, &branch.SlotStepMinutes, &branch.ConcurrentBays,
		&branch.CreatedAt, &branch.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("branch not found")
		}
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	return &branch, nil
}

// UpdateBranch updates branch fields used in operations UI.
func (r *BranchRepository) UpdateBranch(ctx context.Context, id uuid.UUID, name, address *string, phone, email *string, isActive *bool) error {
	q := `UPDATE branches SET `
	var sets []string
	var args []interface{}
	n := 1
	if name != nil {
		sets = append(sets, fmt.Sprintf("name = $%d", n))
		args = append(args, *name)
		n++
	}
	if address != nil {
		sets = append(sets, fmt.Sprintf("address = $%d", n))
		args = append(args, *address)
		n++
	}
	if phone != nil {
		sets = append(sets, fmt.Sprintf("phone = $%d", n))
		args = append(args, *phone)
		n++
	}
	if email != nil {
		sets = append(sets, fmt.Sprintf("email = $%d", n))
		args = append(args, *email)
		n++
	}
	if isActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", n))
		args = append(args, *isActive)
		n++
	}
	if len(sets) == 0 {
		return nil
	}
	q += strings.Join(sets, ", ")
	q += fmt.Sprintf(", updated_at = now() WHERE branch_id = $%d", n)
	args = append(args, id)
	_, err := r.db.Pool.Exec(ctx, q, args...)
	if err != nil {
		return apperr.Internal(err)
	}
	return nil
}

type ServiceAppointmentRepository struct {
	db *database.DB
}

func NewServiceAppointmentRepository(db *database.DB) *ServiceAppointmentRepository {
	return &ServiceAppointmentRepository{db: db}
}

func (r *ServiceAppointmentRepository) CountOverlappingScheduled(ctx context.Context, branchID uuid.UUID, windowStart, windowEnd time.Time) (int, error) {
	var n int
	err := r.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*)::int FROM service_appointments
		WHERE branch_id = $1
		  AND status = 'scheduled'
		  AND appointment_date < $3
		  AND (appointment_date + (duration_minutes || ' minutes')::interval) > $2
	`, branchID, windowStart, windowEnd).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("failed to count overlapping appointments: %w", err)
	}
	return n, nil
}

func (r *ServiceAppointmentRepository) Create(ctx context.Context, create model.ServiceAppointmentCreate, concurrentBays int) (*model.ServiceAppointment, error) {
	if create.DurationMinutes <= 0 {
		return nil, fmt.Errorf("invalid duration")
	}
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1::text))`, create.BranchID.String()); err != nil {
		return nil, fmt.Errorf("lock branch: %w", err)
	}

	winEnd := create.AppointmentDate.Add(time.Duration(create.DurationMinutes) * time.Minute)
	var overlap int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*)::int FROM service_appointments
		WHERE branch_id = $1
		  AND status = 'scheduled'
		  AND appointment_date < $3
		  AND (appointment_date + (duration_minutes || ' minutes')::interval) > $2
	`, create.BranchID, create.AppointmentDate, winEnd).Scan(&overlap)
	if err != nil {
		return nil, fmt.Errorf("overlap check: %w", err)
	}
	if overlap >= concurrentBays {
		return nil, fmt.Errorf("this time slot is no longer available")
	}

	var appointment model.ServiceAppointment
	query := `
		INSERT INTO service_appointments (user_car_id, branch_id, appointment_date, duration_minutes, status, description)
		VALUES ($1, $2, $3, $4, 'scheduled', $5)
		RETURNING service_appointment_id, user_car_id, branch_id, manager_id, appointment_date, duration_minutes, status, description, created_at, updated_at
	`

	err = tx.QueryRow(ctx, query, create.UserCarID, create.BranchID, create.AppointmentDate, create.DurationMinutes, create.Description).Scan(
		&appointment.ServiceAppointmentID, &appointment.UserCarID, &appointment.BranchID,
		&appointment.ManagerID, &appointment.AppointmentDate, &appointment.DurationMinutes, &appointment.Status,
		&appointment.Description, &appointment.CreatedAt, &appointment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	if len(create.ServiceTypeIDs) > 0 {
		query = `
			INSERT INTO service_appointment_types (service_appointment_id, service_type_id)
			SELECT $1, unnest($2::uuid[])
		`
		_, err = tx.Exec(ctx, query, appointment.ServiceAppointmentID, create.ServiceTypeIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to insert service appointment types: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &appointment, nil
}

func (r *ServiceAppointmentRepository) GetByID(ctx context.Context, appointmentID uuid.UUID) (*model.ServiceAppointmentWithDetails, error) {
	var appointment model.ServiceAppointmentWithDetails
	query := `
		SELECT 
			sa.service_appointment_id, sa.user_car_id, sa.branch_id, sa.manager_id,
			sa.appointment_date, sa.duration_minutes, sa.status, sa.description, sa.created_at, sa.updated_at,
			uc.user_id,
			uc.vin as user_car_vin, b.name as branch_name, b.address as branch_address,
			u.first_name || ' ' || u.last_name as manager_name
		FROM service_appointments sa
		JOIN user_cars uc ON sa.user_car_id = uc.user_car_id
		JOIN branches b ON sa.branch_id = b.branch_id
		LEFT JOIN users u ON sa.manager_id = u.user_id
		WHERE sa.service_appointment_id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, appointmentID).Scan(
		&appointment.ServiceAppointmentID, &appointment.UserCarID, &appointment.BranchID,
		&appointment.ManagerID, &appointment.AppointmentDate, &appointment.DurationMinutes, &appointment.Status,
		&appointment.Description, &appointment.CreatedAt, &appointment.UpdatedAt,
		&appointment.OwnerUserID,
		&appointment.UserCarVIN, &appointment.BranchName, &appointment.BranchAddress,
		&appointment.ManagerName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("appointment not found")
		}
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	// Get service types
	query = `
		SELECT st.service_type_id, st.name, st.category, st.description, st.price, st.duration_minutes, st.is_available, st.created_at
		FROM service_types st
		JOIN service_appointment_types sat ON st.service_type_id = sat.service_type_id
		WHERE sat.service_appointment_id = $1
	`
	rows, err := r.db.Pool.Query(ctx, query, appointmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service types: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var st model.ServiceType
		if err := rows.Scan(&st.ServiceTypeID, &st.Name, &st.Category, &st.Description, &st.Price, &st.DurationMinutes, &st.IsAvailable, &st.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan service type: %w", err)
		}
		appointment.ServiceTypes = append(appointment.ServiceTypes, st)
	}

	return &appointment, nil
}

func (r *ServiceAppointmentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.ServiceAppointmentWithDetails, error) {
	query := `
		SELECT 
			sa.service_appointment_id, sa.user_car_id, sa.branch_id, sa.manager_id,
			sa.appointment_date, sa.duration_minutes, sa.status, sa.description, sa.created_at, sa.updated_at,
			uc.vin as user_car_vin, b.name as branch_name, b.address as branch_address,
			u.first_name || ' ' || u.last_name as manager_name
		FROM service_appointments sa
		JOIN user_cars uc ON sa.user_car_id = uc.user_car_id
		JOIN branches b ON sa.branch_id = b.branch_id
		LEFT JOIN users u ON sa.manager_id = u.user_id
		WHERE uc.user_id = $1
		ORDER BY sa.appointment_date DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %w", err)
	}
	defer rows.Close()

	var appointments []model.ServiceAppointmentWithDetails
	for rows.Next() {
		var appointment model.ServiceAppointmentWithDetails
		if err := rows.Scan(
			&appointment.ServiceAppointmentID, &appointment.UserCarID, &appointment.BranchID,
			&appointment.ManagerID, &appointment.AppointmentDate, &appointment.DurationMinutes, &appointment.Status,
			&appointment.Description, &appointment.CreatedAt, &appointment.UpdatedAt,
			&appointment.UserCarVIN, &appointment.BranchName, &appointment.BranchAddress,
			&appointment.ManagerName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}
		appointment.OwnerUserID = userID

		// Get service types for each appointment
		query = `
			SELECT st.service_type_id, st.name, st.category, st.description, st.price, st.duration_minutes, st.is_available, st.created_at
			FROM service_types st
			JOIN service_appointment_types sat ON st.service_type_id = sat.service_type_id
			WHERE sat.service_appointment_id = $1
		`
		rows2, err := r.db.Pool.Query(ctx, query, appointment.ServiceAppointmentID)
		if err == nil {
			for rows2.Next() {
				var st model.ServiceType
				if err := rows2.Scan(&st.ServiceTypeID, &st.Name, &st.Category, &st.Description, &st.Price, &st.DurationMinutes, &st.IsAvailable, &st.CreatedAt); err == nil {
					appointment.ServiceTypes = append(appointment.ServiceTypes, st)
				}
			}
			rows2.Close()
		}

		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

// ListAllWithDetails returns all appointments with details (staff / admin).
func (r *ServiceAppointmentRepository) ListAllWithDetails(ctx context.Context) ([]model.ServiceAppointmentWithDetails, error) {
	query := `
		SELECT 
			sa.service_appointment_id, sa.user_car_id, sa.branch_id, sa.manager_id,
			sa.appointment_date, sa.duration_minutes, sa.status, sa.description, sa.created_at, sa.updated_at,
			uc.vin as user_car_vin, b.name as branch_name, b.address as branch_address,
			u.first_name || ' ' || u.last_name as manager_name,
			uc.user_id,
			owner.email::text,
			TRIM(BOTH FROM owner.first_name || ' ' || owner.last_name)::text
		FROM service_appointments sa
		JOIN user_cars uc ON sa.user_car_id = uc.user_car_id
		JOIN users owner ON uc.user_id = owner.user_id
		JOIN branches b ON sa.branch_id = b.branch_id
		LEFT JOIN users u ON sa.manager_id = u.user_id
		ORDER BY sa.appointment_date DESC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %w", err)
	}
	defer rows.Close()

	var appointments []model.ServiceAppointmentWithDetails
	for rows.Next() {
		var appointment model.ServiceAppointmentWithDetails
		if err := rows.Scan(
			&appointment.ServiceAppointmentID, &appointment.UserCarID, &appointment.BranchID,
			&appointment.ManagerID, &appointment.AppointmentDate, &appointment.DurationMinutes, &appointment.Status,
			&appointment.Description, &appointment.CreatedAt, &appointment.UpdatedAt,
			&appointment.UserCarVIN, &appointment.BranchName, &appointment.BranchAddress,
			&appointment.ManagerName,
			&appointment.OwnerUserID,
			&appointment.OwnerEmail, &appointment.OwnerName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}

		q2 := `
			SELECT st.service_type_id, st.name, st.category, st.description, st.price, st.duration_minutes, st.is_available, st.created_at
			FROM service_types st
			JOIN service_appointment_types sat ON st.service_type_id = sat.service_type_id
			WHERE sat.service_appointment_id = $1
		`
		rows2, err := r.db.Pool.Query(ctx, q2, appointment.ServiceAppointmentID)
		if err == nil {
			for rows2.Next() {
				var st model.ServiceType
				if err := rows2.Scan(&st.ServiceTypeID, &st.Name, &st.Category, &st.Description, &st.Price, &st.DurationMinutes, &st.IsAvailable, &st.CreatedAt); err == nil {
					appointment.ServiceTypes = append(appointment.ServiceTypes, st)
				}
			}
			rows2.Close()
		}

		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *ServiceAppointmentRepository) UpdateStatus(ctx context.Context, appointmentID uuid.UUID, status string) error {
	query := `UPDATE service_appointments SET status = $1 WHERE service_appointment_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, appointmentID)
	if err != nil {
		return fmt.Errorf("failed to update appointment status: %w", err)
	}
	return nil
}

func (r *ServiceAppointmentRepository) UpdateStatusIfCurrent(ctx context.Context, appointmentID uuid.UUID, fromStatus, toStatus string) (bool, error) {
	tag, err := r.db.Pool.Exec(ctx, `
		UPDATE service_appointments
		SET status = $1, updated_at = now()
		WHERE service_appointment_id = $2 AND status = $3
	`, toStatus, appointmentID, fromStatus)
	if err != nil {
		return false, fmt.Errorf("failed to update appointment status conditionally: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

// RescheduleOwned moves a scheduled appointment to newStart after advisory lock and overlap check (excluding this row).
func (r *ServiceAppointmentRepository) RescheduleOwned(ctx context.Context, appointmentID, ownerUserID, branchID uuid.UUID, newStart time.Time, durationMin, concurrentBays int) error {
	if durationMin <= 0 {
		return fmt.Errorf("invalid duration")
	}
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1::text))`, branchID.String()); err != nil {
		return fmt.Errorf("lock branch: %w", err)
	}

	var stub int
	err = tx.QueryRow(ctx, `
		SELECT 1
		FROM service_appointments sa
		JOIN user_cars uc ON sa.user_car_id = uc.user_car_id
		WHERE sa.service_appointment_id = $1
		  AND uc.user_id = $2
		  AND sa.status = 'scheduled'
	`, appointmentID, ownerUserID).Scan(&stub)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("appointment not found or cannot be rescheduled")
		}
		return fmt.Errorf("verify appointment: %w", err)
	}

	newEnd := newStart.Add(time.Duration(durationMin) * time.Minute)
	var overlap int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*)::int FROM service_appointments
		WHERE branch_id = $1
		  AND status = 'scheduled'
		  AND service_appointment_id <> $2
		  AND appointment_date < $4
		  AND (appointment_date + (duration_minutes || ' minutes')::interval) > $3
	`, branchID, appointmentID, newStart, newEnd).Scan(&overlap)
	if err != nil {
		return fmt.Errorf("overlap check: %w", err)
	}
	if overlap >= concurrentBays {
		return fmt.Errorf("this time slot is no longer available")
	}

	tag, err := tx.Exec(ctx, `
		UPDATE service_appointments
		SET appointment_date = $1, updated_at = now()
		WHERE service_appointment_id = $2
		  AND status = 'scheduled'
	`, newStart, appointmentID)
	if err != nil {
		return fmt.Errorf("failed to reschedule: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("appointment could not be updated")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return nil
}
