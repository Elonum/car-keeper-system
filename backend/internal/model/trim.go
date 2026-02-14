package model

import (
	"time"

	"github.com/google/uuid"
)

type Trim struct {
	TrimID        uuid.UUID `db:"trim_id" json:"trim_id"`
	GenerationID  uuid.UUID `db:"generation_id" json:"generation_id"`
	Name          string    `db:"name" json:"name"`
	BasePrice     float64   `db:"base_price" json:"base_price"`
	EngineTypeID  uuid.UUID `db:"engine_type_id" json:"engine_type_id"`
	TransmissionID uuid.UUID `db:"transmission_id" json:"transmission_id"`
	DriveTypeID   uuid.UUID `db:"drive_type_id" json:"drive_type_id"`
	IsAvailable   bool      `db:"is_available" json:"is_available"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type TrimWithDetails struct {
	Trim
	BrandName      string `db:"brand_name" json:"brand_name"`
	ModelName      string `db:"model_name" json:"model_name"`
	GenerationName string `db:"generation_name" json:"generation_name"`
	EngineType     string `db:"engine_type" json:"engine_type"`
	Transmission   string `db:"transmission" json:"transmission"`
	DriveType      string `db:"drive_type" json:"drive_type"`
}

type TrimFilters struct {
	BrandID       []uuid.UUID
	EngineTypeID  []uuid.UUID
	TransmissionID []uuid.UUID
	DriveTypeID   []uuid.UUID
	MinPrice      *float64
	MaxPrice      *float64
	IsAvailable   *bool
}

