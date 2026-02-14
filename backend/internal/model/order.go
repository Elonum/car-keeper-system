package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID         uuid.UUID  `db:"order_id" json:"order_id"`
	UserID          uuid.UUID  `db:"user_id" json:"user_id"`
	ConfigurationID uuid.UUID  `db:"configuration_id" json:"configuration_id"`
	ManagerID       *uuid.UUID `db:"manager_id" json:"manager_id,omitempty"`
	Status          string     `db:"status" json:"status"`
	FinalPrice      float64    `db:"final_price" json:"final_price"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

type OrderCreate struct {
	ConfigurationID uuid.UUID `json:"configuration_id" validate:"required"`
}

type OrderWithDetails struct {
	Order
	Configuration ConfigurationWithDetails `json:"configuration"`
	ManagerName   *string                  `db:"manager_name" json:"manager_name,omitempty"`
}

