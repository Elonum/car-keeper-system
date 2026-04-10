package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func hasModelImageColumns(ctx context.Context, db *database.DB) (bool, error) {
	var count int
	err := db.Pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_schema = 'public'
			AND table_name = 'models'
			AND column_name IN ('image_key', 'image_mime')
	`).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check models image columns: %w", err)
	}
	return count == 2, nil
}

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

// Create inserts a brand (admin catalog).
func (r *BrandRepository) Create(ctx context.Context, name, country string) (*model.Brand, error) {
	var b model.Brand
	err := r.db.Pool.QueryRow(ctx, `
		INSERT INTO brands (name, country) VALUES ($1, $2)
		RETURNING brand_id, name, country, created_at
	`, name, country).Scan(&b.BrandID, &b.Name, &b.Country, &b.CreatedAt)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return &b, nil
}

// Update updates brand name and country.
func (r *BrandRepository) Update(ctx context.Context, brandID uuid.UUID, name, country string) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE brands SET name = $1, country = $2 WHERE brand_id = $3`, name, country, brandID)
	if err != nil {
		return apperr.Internal(err)
	}
	return nil
}

// Delete removes a brand (fails if models reference it).
func (r *BrandRepository) Delete(ctx context.Context, brandID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM brands WHERE brand_id = $1`, brandID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return apperr.BadRequest("Cannot delete brand: referenced by models")
		}
		return apperr.Internal(err)
	}
	return nil
}

type ModelRepository struct {
	db *database.DB
}

func NewModelRepository(db *database.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

func (r *ModelRepository) GetByBrandID(ctx context.Context, brandID uuid.UUID) ([]model.Model, error) {
	hasImageCols, err := hasModelImageColumns(ctx, r.db)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	imageCols := "NULL::text as image_key, NULL::text as image_mime, NULL::text as image_url"
	if hasImageCols {
		imageCols = "image_key, image_mime, CASE WHEN image_key IS NOT NULL THEN '/api/catalog/models/' || model_id::text || '/image' ELSE NULL END as image_url"
	}
	query := fmt.Sprintf(`
		SELECT 
			model_id, brand_id, name, segment, description, %s,
			created_at
		FROM models
		WHERE brand_id = $1
		ORDER BY name
	`, imageCols)
	rows, err := r.db.Pool.Query(ctx, query, brandID)
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}
	defer rows.Close()

	var models []model.Model
	for rows.Next() {
		var m model.Model
		if err := rows.Scan(&m.ModelID, &m.BrandID, &m.Name, &m.Segment, &m.Description, &m.ImageKey, &m.ImageMime, &m.ImageURL, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan model: %w", err)
		}
		models = append(models, m)
	}

	return models, nil
}

func (r *ModelRepository) GetAll(ctx context.Context) ([]model.Model, error) {
	hasImageCols, err := hasModelImageColumns(ctx, r.db)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	imageCols := "NULL::text as image_key, NULL::text as image_mime, NULL::text as image_url"
	if hasImageCols {
		imageCols = "image_key, image_mime, CASE WHEN image_key IS NOT NULL THEN '/api/catalog/models/' || model_id::text || '/image' ELSE NULL END as image_url"
	}
	query := fmt.Sprintf(`
		SELECT 
			model_id, brand_id, name, segment, description, %s,
			created_at
		FROM models
		ORDER BY created_at DESC
	`, imageCols)
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}
	defer rows.Close()

	var models []model.Model
	for rows.Next() {
		var m model.Model
		if err := rows.Scan(&m.ModelID, &m.BrandID, &m.Name, &m.Segment, &m.Description, &m.ImageKey, &m.ImageMime, &m.ImageURL, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan model: %w", err)
		}
		models = append(models, m)
	}
	return models, nil
}

func (r *ModelRepository) Create(ctx context.Context, brandID uuid.UUID, name string, segment, description *string) (*model.Model, error) {
	var m model.Model
	hasImageCols, err := hasModelImageColumns(ctx, r.db)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	query := `
		INSERT INTO models (brand_id, name, segment, description)
		VALUES ($1, $2, $3, $4)
		RETURNING model_id, brand_id, name, segment, description, NULL::text as image_key, NULL::text as image_mime, NULL::text as image_url, created_at
	`
	if hasImageCols {
		query = `
			INSERT INTO models (brand_id, name, segment, description)
			VALUES ($1, $2, $3, $4)
			RETURNING model_id, brand_id, name, segment, description, image_key, image_mime,
				CASE WHEN image_key IS NOT NULL THEN '/api/catalog/models/' || model_id::text || '/image' ELSE NULL END as image_url,
				created_at
		`
	}
	err = r.db.Pool.QueryRow(ctx, query, brandID, name, segment, description).Scan(&m.ModelID, &m.BrandID, &m.Name, &m.Segment, &m.Description, &m.ImageKey, &m.ImageMime, &m.ImageURL, &m.CreatedAt)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return &m, nil
}

func (r *ModelRepository) Update(ctx context.Context, modelID, brandID uuid.UUID, name string, segment, description *string) error {
	cmd, err := r.db.Pool.Exec(ctx, `
		UPDATE models
		SET brand_id = $1, name = $2, segment = $3, description = $4
		WHERE model_id = $5
	`, brandID, name, segment, description, modelID)
	if err != nil {
		return apperr.Internal(err)
	}
	if cmd.RowsAffected() == 0 {
		return apperr.NotFoundErr("Model not found")
	}
	return nil
}

func (r *ModelRepository) Delete(ctx context.Context, modelID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM models WHERE model_id = $1`, modelID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return apperr.BadRequest("Cannot delete model: referenced by generations")
		}
		return apperr.Internal(err)
	}
	return nil
}

func (r *ModelRepository) GetImageMeta(ctx context.Context, modelID uuid.UUID) (key string, mimeType string, err error) {
	hasImageCols, colErr := hasModelImageColumns(ctx, r.db)
	if colErr != nil {
		return "", "", apperr.Internal(colErr)
	}
	if !hasImageCols {
		return "", "", apperr.BadRequest("Database schema is outdated. Add models.image_key and models.image_mime")
	}
	err = r.db.Pool.QueryRow(ctx, `
		SELECT image_key, COALESCE(image_mime, 'application/octet-stream')
		FROM models
		WHERE model_id = $1
	`, modelID).Scan(&key, &mimeType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", apperr.NotFoundErr("Model not found")
		}
		return "", "", apperr.Internal(err)
	}
	if strings.TrimSpace(key) == "" {
		return "", "", apperr.NotFoundErr("Model image not found")
	}
	return key, mimeType, nil
}

func (r *ModelRepository) SetImageMeta(ctx context.Context, modelID uuid.UUID, key, mimeType string) error {
	hasImageCols, colErr := hasModelImageColumns(ctx, r.db)
	if colErr != nil {
		return apperr.Internal(colErr)
	}
	if !hasImageCols {
		return apperr.BadRequest("Database schema is outdated. Add models.image_key and models.image_mime")
	}
	cmd, err := r.db.Pool.Exec(ctx, `
		UPDATE models
		SET image_key = $1, image_mime = $2
		WHERE model_id = $3
	`, key, mimeType, modelID)
	if err != nil {
		return apperr.Internal(err)
	}
	if cmd.RowsAffected() == 0 {
		return apperr.NotFoundErr("Model not found")
	}
	return nil
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
	hasImageCols, err := hasModelImageColumns(ctx, r.db)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	imageURLSelect := "NULL::text as image_url"
	if hasImageCols {
		imageURLSelect = "CASE WHEN m.image_key IS NOT NULL THEN '/api/catalog/models/' || m.model_id::text || '/image' ELSE NULL END as image_url"
	}
	query := fmt.Sprintf(`
		SELECT 
			t.trim_id, t.generation_id, t.name, t.base_price,
			t.engine_type_id, t.transmission_id, t.drive_type_id,
			t.is_available, t.created_at, t.updated_at,
			b.name as brand_name, m.name as model_name, g.name as generation_name,
			m.segment as segment,
			%s,
			et.name as engine_type, tr.name as transmission, dt.name as drive_type
		FROM trims t
		JOIN generations g ON t.generation_id = g.generation_id
		JOIN models m ON g.model_id = m.model_id
		JOIN brands b ON m.brand_id = b.brand_id
		JOIN engine_types et ON t.engine_type_id = et.engine_type_id
		JOIN transmissions tr ON t.transmission_id = tr.transmission_id
		JOIN drive_types dt ON t.drive_type_id = dt.drive_type_id
		WHERE t.trim_id = $1
	`, imageURLSelect)

	var trim model.TrimWithDetails
	err = r.db.Pool.QueryRow(ctx, query, trimID).Scan(
		&trim.TrimID, &trim.GenerationID, &trim.Name, &trim.BasePrice,
		&trim.EngineTypeID, &trim.TransmissionID, &trim.DriveTypeID,
		&trim.IsAvailable, &trim.CreatedAt, &trim.UpdatedAt,
		&trim.BrandName, &trim.ModelName, &trim.GenerationName,
		&trim.Segment, &trim.ImageURL,
		&trim.EngineType, &trim.Transmission, &trim.DriveType,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w", apperr.ErrNotFound)
		}
		return nil, apperr.Internal(err)
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

	if filters.GenerationID != nil {
		conditions = append(conditions, fmt.Sprintf("t.generation_id = $%d", argPos))
		args = append(args, *filters.GenerationID)
		argPos++
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

	hasImageCols, err := hasModelImageColumns(ctx, r.db)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	imageURLSelect := "NULL::text as image_url"
	if hasImageCols {
		imageURLSelect = "CASE WHEN m.image_key IS NOT NULL THEN '/api/catalog/models/' || m.model_id::text || '/image' ELSE NULL END as image_url"
	}
	query := fmt.Sprintf(`
		SELECT 
			t.trim_id, t.generation_id, t.name, t.base_price,
			t.engine_type_id, t.transmission_id, t.drive_type_id,
			t.is_available, t.created_at, t.updated_at,
			b.name as brand_name, m.name as model_name, g.name as generation_name,
			m.segment as segment,
			%s,
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
	`, imageURLSelect, whereClause)

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
			&trim.Segment, &trim.ImageURL,
			&trim.EngineType, &trim.Transmission, &trim.DriveType,
		); err != nil {
			return nil, fmt.Errorf("failed to scan trim: %w", err)
		}
		trims = append(trims, trim)
	}

	return trims, nil
}
