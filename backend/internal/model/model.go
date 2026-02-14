package model

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ModelID     uuid.UUID  `db:"model_id" json:"model_id"`
	BrandID     uuid.UUID  `db:"brand_id" json:"brand_id"`
	Name        string     `db:"name" json:"name"`
	Segment     *string    `db:"segment" json:"segment,omitempty"`
	Description *string    `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

