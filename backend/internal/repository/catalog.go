package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type BrandRepository struct {
	db *database.DB
}

func NewBrandRepository(db *database.DB) *BrandRepository {
	return &BrandRepository{db: db}
}

func (r *BrandRepository) GetAll(ctx context.Context) ([]model.Brand, error) {
	query := `SELECT brand_id, name, country, created_at FROM brands ORDER BY name`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get brands: %w", err)
	}
	defer rows.Close()

	var brands []model.Brand
	for rows.Next() {
		var brand model.Brand
		if err := rows.Scan(&brand.BrandID, &brand.Name, &brand.Country, &brand.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan brand: %w", err)
		}
		brands = append(brands, brand)
	}

	return brands, nil
}

type ModelRepository struct {
	db *database.DB
}

func NewModelRepository(db *database.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

func (r *ModelRepository) GetByBrandID(ctx context.Context, brandID uuid.UUID) ([]model.Model, error) {
	query := `
		SELECT model_id, brand_id, name, segment, description, created_at
		FROM models
		WHERE brand_id = $1
		ORDER BY name
	`
	rows, err := r.db.Pool.Query(ctx, query, brandID)
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}
	defer rows.Close()

	var models []model.Model
	for rows.Next() {
		var m model.Model
		if err := rows.Scan(&m.ModelID, &m.BrandID, &m.Name, &m.Segment, &m.Description, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan model: %w", err)
		}
		models = append(models, m)
	}

	return models, nil
}

type GenerationRepository struct {
	db *database.DB
}

func NewGenerationRepository(db *database.DB) *GenerationRepository {
	return &GenerationRepository{db: db}
}

func (r *GenerationRepository) GetByModelID(ctx context.Context, modelID uuid.UUID) ([]model.Generation, error) {
	query := `
		SELECT generation_id, model_id, name, year_from, year_to, created_at
		FROM generations
		WHERE model_id = $1
		ORDER BY year_from DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, modelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get generations: %w", err)
	}
	defer rows.Close()

	var generations []model.Generation
	for rows.Next() {
		var g model.Generation
		if err := rows.Scan(&g.GenerationID, &g.ModelID, &g.Name, &g.YearFrom, &g.YearTo, &g.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan generation: %w", err)
		}
		generations = append(generations, g)
	}

	return generations, nil
}

type TrimRepository struct {
	db *database.DB
}

func NewTrimRepository(db *database.DB) *TrimRepository {
	return &TrimRepository{db: db}
}

func (r *TrimRepository) GetByID(ctx context.Context, trimID uuid.UUID) (*model.TrimWithDetails, error) {
	query := `
		SELECT 
			t.trim_id, t.generation_id, t.name, t.base_price,
			t.engine_type_id, t.transmission_id, t.drive_type_id,
			t.is_available, t.created_at, t.updated_at,
			b.name as brand_name, m.name as model_name, g.name as generation_name,
			et.name as engine_type, tr.name as transmission, dt.name as drive_type
		FROM trims t
		JOIN generations g ON t.generation_id = g.generation_id
		JOIN models m ON g.model_id = m.model_id
		JOIN brands b ON m.brand_id = b.brand_id
		JOIN engine_types et ON t.engine_type_id = et.engine_type_id
		JOIN transmissions tr ON t.transmission_id = tr.transmission_id
		JOIN drive_types dt ON t.drive_type_id = dt.drive_type_id
		WHERE t.trim_id = $1
	`

	var trim model.TrimWithDetails
	err := r.db.Pool.QueryRow(ctx, query, trimID).Scan(
		&trim.TrimID, &trim.GenerationID, &trim.Name, &trim.BasePrice,
		&trim.EngineTypeID, &trim.TransmissionID, &trim.DriveTypeID,
		&trim.IsAvailable, &trim.CreatedAt, &trim.UpdatedAt,
		&trim.BrandName, &trim.ModelName, &trim.GenerationName,
		&trim.EngineType, &trim.Transmission, &trim.DriveType,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("trim not found")
		}
		return nil, fmt.Errorf("failed to get trim: %w", err)
	}

	return &trim, nil
}

func (r *TrimRepository) GetWithFilters(ctx context.Context, filters model.TrimFilters) ([]model.TrimWithDetails, error) {
	var conditions []string
	var args []interface{}
	argPos := 1

	if len(filters.BrandID) > 0 {
		placeholders := make([]string, len(filters.BrandID))
		for i, id := range filters.BrandID {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, id)
			argPos++
		}
		conditions = append(conditions, fmt.Sprintf("b.brand_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(filters.EngineTypeID) > 0 {
		placeholders := make([]string, len(filters.EngineTypeID))
		for i, id := range filters.EngineTypeID {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, id)
			argPos++
		}
		conditions = append(conditions, fmt.Sprintf("t.engine_type_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(filters.TransmissionID) > 0 {
		placeholders := make([]string, len(filters.TransmissionID))
		for i, id := range filters.TransmissionID {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, id)
			argPos++
		}
		conditions = append(conditions, fmt.Sprintf("t.transmission_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(filters.DriveTypeID) > 0 {
		placeholders := make([]string, len(filters.DriveTypeID))
		for i, id := range filters.DriveTypeID {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, id)
			argPos++
		}
		conditions = append(conditions, fmt.Sprintf("t.drive_type_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if filters.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("t.base_price >= $%d", argPos))
		args = append(args, *filters.MinPrice)
		argPos++
	}

	if filters.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("t.base_price <= $%d", argPos))
		args = append(args, *filters.MaxPrice)
		argPos++
	}

	if filters.IsAvailable != nil {
		conditions = append(conditions, fmt.Sprintf("t.is_available = $%d", argPos))
		args = append(args, *filters.IsAvailable)
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT 
			t.trim_id, t.generation_id, t.name, t.base_price,
			t.engine_type_id, t.transmission_id, t.drive_type_id,
			t.is_available, t.created_at, t.updated_at,
			b.name as brand_name, m.name as model_name, g.name as generation_name,
			et.name as engine_type, tr.name as transmission, dt.name as drive_type
		FROM trims t
		JOIN generations g ON t.generation_id = g.generation_id
		JOIN models m ON g.model_id = m.model_id
		JOIN brands b ON m.brand_id = b.brand_id
		JOIN engine_types et ON t.engine_type_id = et.engine_type_id
		JOIN transmissions tr ON t.transmission_id = tr.transmission_id
		JOIN drive_types dt ON t.drive_type_id = dt.drive_type_id
		%s
		ORDER BY t.created_at DESC
	`, whereClause)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get trims: %w", err)
	}
	defer rows.Close()

	var trims []model.TrimWithDetails
	for rows.Next() {
		var trim model.TrimWithDetails
		if err := rows.Scan(
			&trim.TrimID, &trim.GenerationID, &trim.Name, &trim.BasePrice,
			&trim.EngineTypeID, &trim.TransmissionID, &trim.DriveTypeID,
			&trim.IsAvailable, &trim.CreatedAt, &trim.UpdatedAt,
			&trim.BrandName, &trim.ModelName, &trim.GenerationName,
			&trim.EngineType, &trim.Transmission, &trim.DriveType,
		); err != nil {
			return nil, fmt.Errorf("failed to scan trim: %w", err)
		}
		trims = append(trims, trim)
	}

	return trims, nil
}

