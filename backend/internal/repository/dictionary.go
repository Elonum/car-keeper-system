package repository

import (
	"context"
	"fmt"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/model"
)

type DictionaryRepository struct {
	db *database.DB
}

func NewDictionaryRepository(db *database.DB) *DictionaryRepository {
	return &DictionaryRepository{db: db}
}

func (r *DictionaryRepository) GetEngineTypes(ctx context.Context) ([]model.EngineType, error) {
	query := `SELECT engine_type_id, name, created_at FROM engine_types ORDER BY name`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get engine types: %w", err)
	}
	defer rows.Close()

	var engineTypes []model.EngineType
	for rows.Next() {
		var et model.EngineType
		if err := rows.Scan(&et.EngineTypeID, &et.Name, &et.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan engine type: %w", err)
		}
		engineTypes = append(engineTypes, et)
	}

	return engineTypes, nil
}

func (r *DictionaryRepository) GetTransmissions(ctx context.Context) ([]model.Transmission, error) {
	query := `SELECT transmission_id, name, created_at FROM transmissions ORDER BY name`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get transmissions: %w", err)
	}
	defer rows.Close()

	var transmissions []model.Transmission
	for rows.Next() {
		var t model.Transmission
		if err := rows.Scan(&t.TransmissionID, &t.Name, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan transmission: %w", err)
		}
		transmissions = append(transmissions, t)
	}

	return transmissions, nil
}

func (r *DictionaryRepository) GetDriveTypes(ctx context.Context) ([]model.DriveType, error) {
	query := `SELECT drive_type_id, name, created_at FROM drive_types ORDER BY name`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get drive types: %w", err)
	}
	defer rows.Close()

	var driveTypes []model.DriveType
	for rows.Next() {
		var dt model.DriveType
		if err := rows.Scan(&dt.DriveTypeID, &dt.Name, &dt.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan drive type: %w", err)
		}
		driveTypes = append(driveTypes, dt)
	}

	return driveTypes, nil
}

