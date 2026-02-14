package model

import (
	"time"

	"github.com/google/uuid"
)

type Brand struct {
	BrandID   uuid.UUID `db:"brand_id" json:"brand_id"`
	Name      string    `db:"name" json:"name"`
	Country   string    `db:"country" json:"country"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

