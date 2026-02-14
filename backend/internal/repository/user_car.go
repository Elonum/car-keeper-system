package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserCarRepository struct {
	db *database.DB
}

func NewUserCarRepository(db *database.DB) *UserCarRepository {
	return &UserCarRepository{db: db}
}

func (r *UserCarRepository) Create(ctx context.Context, userID uuid.UUID, create model.UserCarCreate) (*model.UserCar, error) {
	var userCar model.UserCar
	query := `
		INSERT INTO user_cars (user_id, trim_id, color_id, vin, year, current_mileage, purchase_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING user_car_id, user_id, trim_id, color_id, vin, year, current_mileage, purchase_date, created_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		userID, create.TrimID, create.ColorID, create.VIN,
		create.Year, create.CurrentMileage, create.PurchaseDate,
	).Scan(
		&userCar.UserCarID, &userCar.UserID, &userCar.TrimID, &userCar.ColorID,
		&userCar.VIN, &userCar.Year, &userCar.CurrentMileage, &userCar.PurchaseDate,
		&userCar.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user car: %w", err)
	}

	return &userCar, nil
}

func (r *UserCarRepository) GetByID(ctx context.Context, userCarID uuid.UUID) (*model.UserCarWithDetails, error) {
	var userCar model.UserCarWithDetails
	query := `
		SELECT 
			uc.user_car_id, uc.user_id, uc.trim_id, uc.color_id, uc.vin, uc.year,
			uc.current_mileage, uc.purchase_date, uc.created_at,
			t.name as trim_name, b.name as brand_name, m.name as model_name,
			c.name as color_name, c.hex_code as color_hex
		FROM user_cars uc
		JOIN trims t ON uc.trim_id = t.trim_id
		JOIN generations g ON t.generation_id = g.generation_id
		JOIN models m ON g.model_id = m.model_id
		JOIN brands b ON m.brand_id = b.brand_id
		JOIN colors c ON uc.color_id = c.color_id
		WHERE uc.user_car_id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, userCarID).Scan(
		&userCar.UserCarID, &userCar.UserID, &userCar.TrimID, &userCar.ColorID,
		&userCar.VIN, &userCar.Year, &userCar.CurrentMileage, &userCar.PurchaseDate,
		&userCar.CreatedAt, &userCar.TrimName, &userCar.BrandName, &userCar.ModelName,
		&userCar.ColorName, &userCar.ColorHex,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user car not found")
		}
		return nil, fmt.Errorf("failed to get user car: %w", err)
	}

	return &userCar, nil
}

func (r *UserCarRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.UserCarWithDetails, error) {
	query := `
		SELECT 
			uc.user_car_id, uc.user_id, uc.trim_id, uc.color_id, uc.vin, uc.year,
			uc.current_mileage, uc.purchase_date, uc.created_at,
			t.name as trim_name, b.name as brand_name, m.name as model_name,
			c.name as color_name, c.hex_code as color_hex
		FROM user_cars uc
		JOIN trims t ON uc.trim_id = t.trim_id
		JOIN generations g ON t.generation_id = g.generation_id
		JOIN models m ON g.model_id = m.model_id
		JOIN brands b ON m.brand_id = b.brand_id
		JOIN colors c ON uc.color_id = c.color_id
		WHERE uc.user_id = $1
		ORDER BY uc.created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user cars: %w", err)
	}
	defer rows.Close()

	var userCars []model.UserCarWithDetails
	for rows.Next() {
		var userCar model.UserCarWithDetails
		if err := rows.Scan(
			&userCar.UserCarID, &userCar.UserID, &userCar.TrimID, &userCar.ColorID,
			&userCar.VIN, &userCar.Year, &userCar.CurrentMileage, &userCar.PurchaseDate,
			&userCar.CreatedAt, &userCar.TrimName, &userCar.BrandName, &userCar.ModelName,
			&userCar.ColorName, &userCar.ColorHex,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user car: %w", err)
		}
		userCars = append(userCars, userCar)
	}

	return userCars, nil
}

func (r *UserCarRepository) VINExists(ctx context.Context, vin string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user_cars WHERE vin = $1)`
	err := r.db.Pool.QueryRow(ctx, query, vin).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check VIN: %w", err)
	}
	return exists, nil
}

