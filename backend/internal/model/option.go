package model

import (
	"time"

	"github.com/google/uuid"
)

type Option struct {
	OptionID    uuid.UUID `db:"option_id" json:"option_id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description,omitempty"`
	Price       float64   `db:"price" json:"price"`
	IsAvailable bool      `db:"is_available" json:"is_available"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

