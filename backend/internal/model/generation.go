package model

import (
	"time"

	"github.com/google/uuid"
)

type Generation struct {
	GenerationID uuid.UUID `db:"generation_id" json:"generation_id"`
	ModelID      uuid.UUID `db:"model_id" json:"model_id"`
	Name         string    `db:"name" json:"name"`
	YearFrom     int       `db:"year_from" json:"year_from"`
	YearTo       *int      `db:"year_to" json:"year_to,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

