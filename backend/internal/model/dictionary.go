package model

import (
	"time"

	"github.com/google/uuid"
)

type EngineType struct {
	EngineTypeID uuid.UUID `db:"engine_type_id" json:"engine_type_id"`
	Name         string    `db:"name" json:"name"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Transmission struct {
	TransmissionID uuid.UUID `db:"transmission_id" json:"transmission_id"`
	Name          string    `db:"name" json:"name"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type DriveType struct {
	DriveTypeID uuid.UUID `db:"drive_type_id" json:"drive_type_id"`
	Name        string    `db:"name" json:"name"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

