package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

type BranchRepository struct {
	db *database.DB
}

func NewBranchRepository(db *database.DB) *BranchRepository {
	return &BranchRepository{db: db}
}

func (r *BranchRepository) GetAll(ctx context.Context, isActive *bool) ([]model.Branch, error) {
	query := `SELECT branch_id, name, address, phone, email, is_active, created_at, updated_at FROM branches`
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
		if err := rows.Scan(&branch.BranchID, &branch.Name, &branch.Address, &branch.Phone, &branch.Email, &branch.IsActive, &branch.CreatedAt, &branch.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan branch: %w", err)
		}
		branches = append(branches, branch)
	}

	return branches, nil
}

func (r *BranchRepository) GetByID(ctx context.Context, branchID uuid.UUID) (*model.Branch, error) {
	var branch model.Branch
	query := `SELECT branch_id, name, address, phone, email, is_active, created_at, updated_at FROM branches WHERE branch_id = $1`

	err := r.db.Pool.QueryRow(ctx, query, branchID).Scan(
		&branch.BranchID, &branch.Name, &branch.Address, &branch.Phone, &branch.Email,
		&branch.IsActive, &branch.CreatedAt, &branch.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("branch not found")
		}
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	return &branch, nil
}

type ServiceAppointmentRepository struct {
	db *database.DB
}

func NewServiceAppointmentRepository(db *database.DB) *ServiceAppointmentRepository {
	return &ServiceAppointmentRepository{db: db}
}

func (r *ServiceAppointmentRepository) Create(ctx context.Context, create model.ServiceAppointmentCreate) (*model.ServiceAppointment, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var appointment model.ServiceAppointment
	query := `
		INSERT INTO service_appointments (user_car_id, branch_id, appointment_date, status, description)
		VALUES ($1, $2, $3, 'scheduled', $4)
		RETURNING service_appointment_id, user_car_id, branch_id, manager_id, appointment_date, status, description, created_at, updated_at
	`

	err = tx.QueryRow(ctx, query, create.UserCarID, create.BranchID, create.AppointmentDate, create.Description).Scan(
		&appointment.ServiceAppointmentID, &appointment.UserCarID, &appointment.BranchID,
		&appointment.ManagerID, &appointment.AppointmentDate, &appointment.Status,
		&appointment.Description, &appointment.CreatedAt, &appointment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	// Insert service appointment types
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
			sa.appointment_date, sa.status, sa.description, sa.created_at, sa.updated_at,
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
		&appointment.ManagerID, &appointment.AppointmentDate, &appointment.Status,
		&appointment.Description, &appointment.CreatedAt, &appointment.UpdatedAt,
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
			sa.appointment_date, sa.status, sa.description, sa.created_at, sa.updated_at,
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
			&appointment.ManagerID, &appointment.AppointmentDate, &appointment.Status,
			&appointment.Description, &appointment.CreatedAt, &appointment.UpdatedAt,
			&appointment.UserCarVIN, &appointment.BranchName, &appointment.BranchAddress,
			&appointment.ManagerName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}

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

func (r *ServiceAppointmentRepository) UpdateStatus(ctx context.Context, appointmentID uuid.UUID, status string) error {
	query := `UPDATE service_appointments SET status = $1 WHERE service_appointment_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, appointmentID)
	if err != nil {
		return fmt.Errorf("failed to update appointment status: %w", err)
	}
	return nil
}

