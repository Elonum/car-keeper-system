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

type ColorRepository struct {
	db *database.DB
}

func NewColorRepository(db *database.DB) *ColorRepository {
	return &ColorRepository{db: db}
}

func (r *ColorRepository) GetAll(ctx context.Context, isAvailable *bool) ([]model.Color, error) {
	query := `SELECT color_id, name, hex_code, price_delta, is_available, created_at FROM colors`
	var args []interface{}
	
	if isAvailable != nil {
		query += ` WHERE is_available = $1`
		args = append(args, *isAvailable)
	}
	query += ` ORDER BY name`

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get colors: %w", err)
	}
	defer rows.Close()

	var colors []model.Color
	for rows.Next() {
		var color model.Color
		if err := rows.Scan(&color.ColorID, &color.Name, &color.HexCode, &color.PriceDelta, &color.IsAvailable, &color.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan color: %w", err)
		}
		colors = append(colors, color)
	}

	return colors, nil
}

func (r *ColorRepository) GetByID(ctx context.Context, colorID uuid.UUID) (*model.Color, error) {
	var color model.Color
	query := `SELECT color_id, name, hex_code, price_delta, is_available, created_at FROM colors WHERE color_id = $1`
	
	err := r.db.Pool.QueryRow(ctx, query, colorID).Scan(
		&color.ColorID, &color.Name, &color.HexCode, &color.PriceDelta, &color.IsAvailable, &color.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("color not found")
		}
		return nil, fmt.Errorf("failed to get color: %w", err)
	}

	return &color, nil
}

type OptionRepository struct {
	db *database.DB
}

func NewOptionRepository(db *database.DB) *OptionRepository {
	return &OptionRepository{db: db}
}

func (r *OptionRepository) GetByTrimID(ctx context.Context, trimID uuid.UUID) ([]model.Option, error) {
	query := `
		SELECT o.option_id, o.name, o.description, o.price, o.is_available, o.created_at
		FROM options o
		JOIN trim_options tro ON o.option_id = tro.option_id
		WHERE tro.trim_id = $1 AND o.is_available = true
		ORDER BY o.name
	`
	rows, err := r.db.Pool.Query(ctx, query, trimID)
	if err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
	}
	defer rows.Close()

	var options []model.Option
	for rows.Next() {
		var opt model.Option
		if err := rows.Scan(&opt.OptionID, &opt.Name, &opt.Description, &opt.Price, &opt.IsAvailable, &opt.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan option: %w", err)
		}
		options = append(options, opt)
	}

	return options, nil
}

func (r *OptionRepository) GetByIDs(ctx context.Context, optionIDs []uuid.UUID) ([]model.Option, error) {
	if len(optionIDs) == 0 {
		return []model.Option{}, nil
	}

	query := `
		SELECT option_id, name, description, price, is_available, created_at
		FROM options
		WHERE option_id = ANY($1)
	`
	rows, err := r.db.Pool.Query(ctx, query, optionIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
	}
	defer rows.Close()

	var options []model.Option
	for rows.Next() {
		var opt model.Option
		if err := rows.Scan(&opt.OptionID, &opt.Name, &opt.Description, &opt.Price, &opt.IsAvailable, &opt.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan option: %w", err)
		}
		options = append(options, opt)
	}

	return options, nil
}

type ConfigurationRepository struct {
	db *database.DB
}

func NewConfigurationRepository(db *database.DB) *ConfigurationRepository {
	return &ConfigurationRepository{db: db}
}

func (r *ConfigurationRepository) Create(ctx context.Context, userID uuid.UUID, create model.ConfigurationCreate, totalPrice float64) (*model.Configuration, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var config model.Configuration
	query := `
		INSERT INTO configurations (user_id, trim_id, color_id, status, total_price)
		VALUES ($1, $2, $3, 'draft', $4)
		RETURNING configuration_id, user_id, trim_id, color_id, status, total_price, created_at, updated_at
	`

	err = tx.QueryRow(ctx, query, userID, create.TrimID, create.ColorID, totalPrice).Scan(
		&config.ConfigurationID, &config.UserID, &config.TrimID, &config.ColorID,
		&config.Status, &config.TotalPrice, &config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create configuration: %w", err)
	}

	// Insert configuration options
	if len(create.OptionIDs) > 0 {
		optQuery := `
			INSERT INTO configuration_options (configuration_id, option_id)
			SELECT $1, unnest($2::uuid[])
		`
		_, err = tx.Exec(ctx, optQuery, config.ConfigurationID, create.OptionIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to insert configuration options: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &config, nil
}

func (r *ConfigurationRepository) GetByID(ctx context.Context, configID uuid.UUID) (*model.ConfigurationWithDetails, error) {
	var config model.ConfigurationWithDetails
	query := `
		SELECT 
			c.configuration_id, c.user_id, c.trim_id, c.color_id, c.status, c.total_price,
			c.created_at, c.updated_at,
			t.name as trim_name, col.name as color_name, col.hex_code as color_hex
		FROM configurations c
		JOIN trims t ON c.trim_id = t.trim_id
		JOIN colors col ON c.color_id = col.color_id
		WHERE c.configuration_id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, configID).Scan(
		&config.ConfigurationID, &config.UserID, &config.TrimID, &config.ColorID,
		&config.Status, &config.TotalPrice, &config.CreatedAt, &config.UpdatedAt,
		&config.TrimName, &config.ColorName, &config.ColorHex,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("configuration not found")
		}
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	// Get options
	optQuery := `
		SELECT o.option_id, o.name, o.description, o.price, o.is_available, o.created_at
		FROM options o
		JOIN configuration_options co ON o.option_id = co.option_id
		WHERE co.configuration_id = $1
	`
	optRows, err := r.db.Pool.Query(ctx, optQuery, configID)
	if err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
	}
	defer optRows.Close()

	for optRows.Next() {
		var opt model.Option
		if err := optRows.Scan(&opt.OptionID, &opt.Name, &opt.Description, &opt.Price, &opt.IsAvailable, &opt.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan option: %w", err)
		}
		config.Options = append(config.Options, opt)
	}

	return &config, nil
}

func (r *ConfigurationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.ConfigurationWithDetails, error) {
	query := `
		SELECT 
			c.configuration_id, c.user_id, c.trim_id, c.color_id, c.status, c.total_price,
			c.created_at, c.updated_at,
			t.name as trim_name, col.name as color_name, col.hex_code as color_hex
		FROM configurations c
		JOIN trims t ON c.trim_id = t.trim_id
		JOIN colors col ON c.color_id = col.color_id
		WHERE c.user_id = $1
		ORDER BY c.created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get configurations: %w", err)
	}
	defer rows.Close()

	var configs []model.ConfigurationWithDetails
	for rows.Next() {
		var config model.ConfigurationWithDetails
		if err := rows.Scan(
			&config.ConfigurationID, &config.UserID, &config.TrimID, &config.ColorID,
			&config.Status, &config.TotalPrice, &config.CreatedAt, &config.UpdatedAt,
			&config.TrimName, &config.ColorName, &config.ColorHex,
		); err != nil {
			return nil, fmt.Errorf("failed to scan configuration: %w", err)
		}

		// Get options for each configuration
		optQuery := `
			SELECT o.option_id, o.name, o.description, o.price, o.is_available, o.created_at
			FROM options o
			JOIN configuration_options co ON o.option_id = co.option_id
			WHERE co.configuration_id = $1
		`
		optRows, err := r.db.Pool.Query(ctx, optQuery, config.ConfigurationID)
		if err == nil {
			for optRows.Next() {
				var opt model.Option
				if err := optRows.Scan(&opt.OptionID, &opt.Name, &opt.Description, &opt.Price, &opt.IsAvailable, &opt.CreatedAt); err == nil {
					config.Options = append(config.Options, opt)
				}
			}
			optRows.Close()
		}

		configs = append(configs, config)
	}

	return configs, nil
}

func (r *ConfigurationRepository) UpdateStatus(ctx context.Context, configID uuid.UUID, status string) error {
	query := `UPDATE configurations SET status = $1 WHERE configuration_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, configID)
	if err != nil {
		return fmt.Errorf("failed to update configuration status: %w", err)
	}
	return nil
}

func (r *ConfigurationRepository) Update(ctx context.Context, configID uuid.UUID, update model.ConfigurationCreate, totalPrice float64) (*model.Configuration, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update configuration
	query := `
		UPDATE configurations 
		SET trim_id = $1, color_id = $2, total_price = $3, updated_at = NOW()
		WHERE configuration_id = $4
		RETURNING configuration_id, user_id, trim_id, color_id, status, total_price, created_at, updated_at
	`
	var config model.Configuration
	err = tx.QueryRow(ctx, query, update.TrimID, update.ColorID, totalPrice, configID).Scan(
		&config.ConfigurationID, &config.UserID, &config.TrimID, &config.ColorID,
		&config.Status, &config.TotalPrice, &config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update configuration: %w", err)
	}

	// Delete existing options
	_, err = tx.Exec(ctx, `DELETE FROM configuration_options WHERE configuration_id = $1`, configID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete existing options: %w", err)
	}

	// Insert new options
	if len(update.OptionIDs) > 0 {
		optQuery := `
			INSERT INTO configuration_options (configuration_id, option_id)
			SELECT $1, unnest($2::uuid[])
		`
		_, err = tx.Exec(ctx, optQuery, configID, update.OptionIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to insert configuration options: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &config, nil
}

func (r *ConfigurationRepository) Delete(ctx context.Context, configID uuid.UUID) error {
	query := `DELETE FROM configurations WHERE configuration_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, configID)
	if err != nil {
		return fmt.Errorf("failed to delete configuration: %w", err)
	}
	return nil
}

