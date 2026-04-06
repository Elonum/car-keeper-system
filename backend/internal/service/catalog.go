package service

import (
	"context"
	"strings"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/google/uuid"
)

type CatalogService struct {
	repo *repository.Repository
}

func NewCatalogService(repos *repository.Repository) *CatalogService {
	return &CatalogService{repo: repos}
}

func (s *CatalogService) GetBrands(ctx context.Context) ([]model.Brand, error) {
	return s.repo.Brand.GetAll(ctx)
}

func (s *CatalogService) GetModels(ctx context.Context, brandID uuid.UUID) ([]model.Model, error) {
	return s.repo.Model.GetByBrandID(ctx, brandID)
}

func (s *CatalogService) GetGenerations(ctx context.Context, modelID uuid.UUID) ([]model.Generation, error) {
	return s.repo.Generation.GetByModelID(ctx, modelID)
}

func (s *CatalogService) GetTrim(ctx context.Context, trimID uuid.UUID) (*model.TrimWithDetails, error) {
	return s.repo.Trim.GetByID(ctx, trimID)
}

func (s *CatalogService) GetTrims(ctx context.Context, filters model.TrimFilters) ([]model.TrimWithDetails, error) {
	return s.repo.Trim.GetWithFilters(ctx, filters)
}

func (s *CatalogService) GetEngineTypes(ctx context.Context) ([]model.EngineType, error) {
	return s.repo.Dictionary.GetEngineTypes(ctx)
}

func (s *CatalogService) GetTransmissions(ctx context.Context) ([]model.Transmission, error) {
	return s.repo.Dictionary.GetTransmissions(ctx)
}

func (s *CatalogService) GetDriveTypes(ctx context.Context) ([]model.DriveType, error) {
	return s.repo.Dictionary.GetDriveTypes(ctx)
}

func (s *CatalogService) AdminCreateBrand(ctx context.Context, name, country string) (*model.Brand, error) {
	name = strings.TrimSpace(name)
	country = strings.TrimSpace(country)
	if name == "" || country == "" {
		return nil, apperr.BadRequest("name and country are required")
	}
	return s.repo.Brand.Create(ctx, name, country)
}

func (s *CatalogService) AdminUpdateBrand(ctx context.Context, id uuid.UUID, name, country string) error {
	name = strings.TrimSpace(name)
	country = strings.TrimSpace(country)
	if name == "" || country == "" {
		return apperr.BadRequest("name and country are required")
	}
	return s.repo.Brand.Update(ctx, id, name, country)
}

func (s *CatalogService) AdminDeleteBrand(ctx context.Context, id uuid.UUID) error {
	return s.repo.Brand.Delete(ctx, id)
}

