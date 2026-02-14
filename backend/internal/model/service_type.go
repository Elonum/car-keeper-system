package model

import (
	"time"

	"github.com/google/uuid"
)

type ServiceType struct {
	ServiceTypeID  uuid.UUID `db:"service_type_id" json:"service_type_id"`
	Name           string    `db:"name" json:"name"`
	Category       string    `db:"category" json:"category"`
	Description    *string   `db:"description" json:"description,omitempty"`
	Price          float64   `db:"price" json:"price"`
	DurationMinutes *int     `db:"duration_minutes" json:"duration_minutes,omitempty"`
	IsAvailable    bool      `db:"is_available" json:"is_available"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

