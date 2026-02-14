package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Email     string    `db:"email" json:"email"`
	Phone     *string   `db:"phone" json:"phone,omitempty"`
	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UserCreate struct {
	FirstName string  `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string  `json:"last_name" validate:"required,min=1,max=100"`
	Email     string  `json:"email" validate:"required,email"`
	Phone     *string `json:"phone,omitempty"`
	Password  string  `json:"password" validate:"required,min=6"`
	Role      string  `json:"role,omitempty"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone,omitempty"`
	Role      string    `json:"role"`
	FullName  string    `json:"full_name"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		UserID:    u.UserID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		FullName:  u.FirstName + " " + u.LastName,
	}
}

