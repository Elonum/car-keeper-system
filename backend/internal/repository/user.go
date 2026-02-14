package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, userCreate model.UserCreate) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	role := "customer"
	if userCreate.Role != "" {
		role = userCreate.Role
	}

	var user model.User
	query := `
		INSERT INTO users (first_name, last_name, email, phone, role, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING user_id, first_name, last_name, email, phone, role, created_at, updated_at
	`

	err = r.db.Pool.QueryRow(ctx, query,
		userCreate.FirstName,
		userCreate.LastName,
		userCreate.Email,
		userCreate.Phone,
		role,
		hashedPassword,
	).Scan(
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	var user model.User
	query := `
		SELECT user_id, first_name, last_name, email, phone, role, created_at, updated_at
		FROM users
		WHERE user_id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	var passwordHash string
	query := `
		SELECT user_id, first_name, last_name, email, phone, role, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.Role,
		&passwordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Store password hash for verification (not in model)
	_ = passwordHash

	return &user, nil
}

func (r *UserRepository) VerifyPassword(ctx context.Context, email, password string) (*model.User, error) {
	var user model.User
	var passwordHash string
	query := `
		SELECT user_id, first_name, last_name, email, phone, role, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.Role,
		&passwordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &user, nil
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email: %w", err)
	}
	return exists, nil
}

