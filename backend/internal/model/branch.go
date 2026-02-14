package model

import (
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	BranchID  uuid.UUID `db:"branch_id" json:"branch_id"`
	Name      string    `db:"name" json:"name"`
	Address   string    `db:"address" json:"address"`
	Phone     *string   `db:"phone" json:"phone,omitempty"`
	Email     *string   `db:"email" json:"email,omitempty"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

