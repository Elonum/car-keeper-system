package model

import (
	"time"

	"github.com/google/uuid"
)

type UserCar struct {
	UserCarID     uuid.UUID  `db:"user_car_id" json:"user_car_id"`
	UserID        uuid.UUID  `db:"user_id" json:"user_id"`
	TrimID        uuid.UUID  `db:"trim_id" json:"trim_id"`
	ColorID       uuid.UUID  `db:"color_id" json:"color_id"`
	VIN           string     `db:"vin" json:"vin"`
	Year          int        `db:"year" json:"year"`
	CurrentMileage int       `db:"current_mileage" json:"current_mileage"`
	PurchaseDate  *time.Time `db:"purchase_date" json:"purchase_date,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}

type UserCarCreate struct {
	TrimID        uuid.UUID  `json:"trim_id" validate:"required"`
	ColorID        uuid.UUID  `json:"color_id" validate:"required"`
	VIN            string     `json:"vin" validate:"required,len=17"`
	Year           int        `json:"year" validate:"required,min=1900"`
	CurrentMileage int        `json:"current_mileage" validate:"min=0"`
	PurchaseDate   *time.Time `json:"purchase_date,omitempty"`
}

type UserCarWithDetails struct {
	UserCar
	TrimName   string `db:"trim_name" json:"trim_name"`
	BrandName  string `db:"brand_name" json:"brand_name"`
	ModelName  string `db:"model_name" json:"model_name"`
	ColorName  string `db:"color_name" json:"color_name"`
	ColorHex   *string `db:"color_hex" json:"color_hex,omitempty"`
}

