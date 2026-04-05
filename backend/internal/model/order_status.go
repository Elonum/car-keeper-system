package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatusDefinition struct {
	OrderStatusID   uuid.UUID `db:"order_status_id" json:"order_status_id"`
	Code            string    `db:"code" json:"code"`
	CustomerLabelRu string    `db:"customer_label_ru" json:"customer_label_ru"`
	AdminLabelRu    *string   `db:"admin_label_ru" json:"admin_label_ru,omitempty"`
	Description     *string   `db:"description" json:"description,omitempty"`
	SortOrder       int       `db:"sort_order" json:"sort_order"`
	IsActive        bool      `db:"is_active" json:"is_active"`
	IsTerminal      bool      `db:"is_terminal" json:"is_terminal"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// OrderStatusPublic is a minimal view for clients (active statuses list).
type OrderStatusPublic struct {
	Code            string `json:"code"`
	CustomerLabelRu string `json:"customer_label_ru"`
	IsTerminal      bool   `json:"is_terminal"`
	SortOrder       int    `json:"sort_order"`
}

type OrderStatusCreate struct {
	Code            string  `json:"code"`
	CustomerLabelRu string  `json:"customer_label_ru"`
	AdminLabelRu    *string `json:"admin_label_ru,omitempty"`
	Description     *string `json:"description,omitempty"`
	SortOrder       *int    `json:"sort_order,omitempty"`
	IsActive        *bool   `json:"is_active,omitempty"`
	IsTerminal      *bool   `json:"is_terminal,omitempty"`
}

type OrderStatusUpdate struct {
	CustomerLabelRu *string `json:"customer_label_ru,omitempty"`
	AdminLabelRu    *string `json:"admin_label_ru,omitempty"`
	Description     *string `json:"description,omitempty"`
	SortOrder       *int    `json:"sort_order,omitempty"`
	IsActive        *bool   `json:"is_active,omitempty"`
	IsTerminal      *bool   `json:"is_terminal,omitempty"`
	Code            *string `json:"code,omitempty"`
}
