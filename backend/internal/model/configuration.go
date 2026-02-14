package model

import (
	"time"

	"github.com/google/uuid"
)

type Configuration struct {
	ConfigurationID uuid.UUID `db:"configuration_id" json:"configuration_id"`
	UserID          uuid.UUID `db:"user_id" json:"user_id"`
	TrimID          uuid.UUID `db:"trim_id" json:"trim_id"`
	ColorID         uuid.UUID `db:"color_id" json:"color_id"`
	Status          string    `db:"status" json:"status"`
	TotalPrice      float64   `db:"total_price" json:"total_price"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type ConfigurationCreate struct {
	TrimID     uuid.UUID   `json:"trim_id" validate:"required"`
	ColorID    uuid.UUID   `json:"color_id" validate:"required"`
	OptionIDs  []uuid.UUID `json:"option_ids"`
}

type ConfigurationWithDetails struct {
	Configuration
	TrimName   string `db:"trim_name" json:"trim_name"`
	ColorName  string `db:"color_name" json:"color_name"`
	ColorHex   *string `db:"color_hex" json:"color_hex,omitempty"`
	Options    []Option `json:"options"`
}

