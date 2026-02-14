package model

import (
	"time"

	"github.com/google/uuid"
)

type ServiceAppointment struct {
	ServiceAppointmentID uuid.UUID  `db:"service_appointment_id" json:"service_appointment_id"`
	UserCarID           uuid.UUID  `db:"user_car_id" json:"user_car_id"`
	BranchID            uuid.UUID  `db:"branch_id" json:"branch_id"`
	ManagerID           *uuid.UUID `db:"manager_id" json:"manager_id,omitempty"`
	AppointmentDate     time.Time  `db:"appointment_date" json:"appointment_date"`
	Status              string     `db:"status" json:"status"`
	Description         *string    `db:"description" json:"description,omitempty"`
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at" json:"updated_at"`
}

type ServiceAppointmentCreate struct {
	UserCarID       uuid.UUID   `json:"user_car_id" validate:"required"`
	BranchID        uuid.UUID   `json:"branch_id" validate:"required"`
	ServiceTypeIDs  []uuid.UUID `json:"service_type_ids" validate:"required,min=1"`
	AppointmentDate time.Time   `json:"appointment_date" validate:"required"`
	Description     *string     `json:"description,omitempty"`
}

type ServiceAppointmentWithDetails struct {
	ServiceAppointment
	UserCarVIN     string    `db:"user_car_vin" json:"user_car_vin"`
	BranchName     string    `db:"branch_name" json:"branch_name"`
	BranchAddress  string    `db:"branch_address" json:"branch_address"`
	ManagerName    *string   `db:"manager_name" json:"manager_name,omitempty"`
	ServiceTypes   []ServiceType `json:"service_types"`
}

