package model

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ModelID     uuid.UUID `db:"model_id" json:"model_id"`
	BrandID     uuid.UUID `db:"brand_id" json:"brand_id"`
	Name        string    `db:"name" json:"name"`
	Segment     *string   `db:"segment" json:"segment,omitempty"`
	Description *string   `db:"description" json:"description,omitempty"`
	ImageKey    *string   `db:"image_key" json:"-"`
	ImageMime   *string   `db:"image_mime" json:"-"`
	ImageURL    *string   `db:"image_url" json:"image_url,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
