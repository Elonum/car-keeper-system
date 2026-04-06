package repository

import (
	"context"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
)

type RoleRepository struct {
	db *database.DB
}

func NewRoleRepository(db *database.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) ListAll(ctx context.Context) ([]model.RoleDefinition, error) {
	query := `
		SELECT role_id, code, name_ru, description, sort_order, is_staff, created_at, updated_at
		FROM role_definitions
		ORDER BY sort_order, code
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	defer rows.Close()

	var out []model.RoleDefinition
	for rows.Next() {
		var d model.RoleDefinition
		if err := rows.Scan(
			&d.RoleID, &d.Code, &d.NameRu, &d.Description, &d.SortOrder, &d.IsStaff,
			&d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, apperr.Internal(err)
		}
		out = append(out, d)
	}
	return out, nil
}

// ListStaffCodes returns role codes with is_staff = true (for authz bootstrap).
func (r *RoleRepository) ListStaffCodes(ctx context.Context) ([]string, error) {
	query := `SELECT code FROM role_definitions WHERE is_staff = true ORDER BY sort_order, code`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, apperr.Internal(err)
		}
		codes = append(codes, c)
	}
	return codes, nil
}

// LoadRolePermissions returns role_code -> permission codes from role_permissions.
func (r *RoleRepository) LoadRolePermissions(ctx context.Context) (map[string][]string, error) {
	query := `
		SELECT role_code, permission_code
		FROM role_permissions
		ORDER BY role_code, permission_code
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	defer rows.Close()

	out := make(map[string][]string)
	for rows.Next() {
		var roleCode, perm string
		if err := rows.Scan(&roleCode, &perm); err != nil {
			return nil, apperr.Internal(err)
		}
		out[roleCode] = append(out[roleCode], perm)
	}
	return out, nil
}
