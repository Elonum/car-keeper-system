package model

import (
	"time"

	"github.com/google/uuid"
)

type Color struct {
	ColorID    uuid.UUID `db:"color_id" json:"color_id"`
	Name       string    `db:"name" json:"name"`
	HexCode    *string   `db:"hex_code" json:"hex_code,omitempty"`
	PriceDelta float64   `db:"price_delta" json:"price_delta"`
	IsAvailable bool     `db:"is_available" json:"is_available"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

