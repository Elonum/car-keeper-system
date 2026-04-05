package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/apperr"
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

	// Public registration is always customer; privileged roles use separate admin flows.
	role := "customer"

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
			return nil, fmt.Errorf("%w", apperr.ErrNotFound)
		}
		return nil, apperr.Internal(err)
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
			return nil, fmt.Errorf("%w", apperr.ErrInvalidCredentials)
		}
		return nil, apperr.Internal(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("%w", apperr.ErrInvalidCredentials)
	}

	return &user, nil
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, apperr.Internal(err)
	}
	return exists, nil
}

// VerifyPasswordForUserID checks the plaintext password against the stored hash for userID.
func (r *UserRepository) VerifyPasswordForUserID(ctx context.Context, userID uuid.UUID, plain string) error {
	var passwordHash string
	err := r.db.Pool.QueryRow(ctx,
		`SELECT password_hash FROM users WHERE user_id = $1`,
		userID,
	).Scan(&passwordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w", apperr.ErrInvalidCredentials)
		}
		return apperr.Internal(err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plain)); err != nil {
		return fmt.Errorf("%w", apperr.ErrInvalidCredentials)
	}
	return nil
}

// UpdatePassword replaces the password hash for the user.
func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, plain string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	ct, err := r.db.Pool.Exec(ctx,
		`UPDATE users SET password_hash = $1, updated_at = now() WHERE user_id = $2`,
		string(hashedPassword), userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// UpdateProfile updates editable profile fields (email is not changed here).
func (r *UserRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, firstName, lastName string, phone *string) (*model.User, error) {
	var user model.User
	err := r.db.Pool.QueryRow(ctx, `
		UPDATE users
		SET first_name = $1, last_name = $2, phone = $3, updated_at = now()
		WHERE user_id = $4
		RETURNING user_id, first_name, last_name, email, phone, role, created_at, updated_at
	`, firstName, lastName, phone, userID).Scan(
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
			return nil, fmt.Errorf("%w", apperr.ErrNotFound)
		}
		return nil, apperr.Internal(err)
	}
	return &user, nil
}

