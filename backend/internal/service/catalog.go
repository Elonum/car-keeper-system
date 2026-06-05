package service

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/carkeeper/backend/internal/storage"
	"github.com/carkeeper/backend/internal/upload"
	"github.com/carkeeper/backend/internal/validate"
	"github.com/google/uuid"
)

type CatalogService struct {
	repo          *repository.Repository
	store         storage.FileStorage
	maxUploadSize int64
}

func NewCatalogService(repos *repository.Repository, store storage.FileStorage, maxUploadSize int64) *CatalogService {
	return &CatalogService{repo: repos, store: store, maxUploadSize: maxUploadSize}
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
	var msg string
	name, msg = validate.BrandName(name)
	if msg != "" {
		return nil, apperr.BadRequest(msg)
	}
	country, msg = validate.BrandCountry(country)
	if msg != "" {
		return nil, apperr.BadRequest(msg)
	}
	return s.repo.Brand.Create(ctx, name, country)
}

func (s *CatalogService) AdminUpdateBrand(ctx context.Context, id uuid.UUID, name, country string) error {
	var msg string
	name, msg = validate.BrandName(name)
	if msg != "" {
		return apperr.BadRequest(msg)
	}
	country, msg = validate.BrandCountry(country)
	if msg != "" {
		return apperr.BadRequest(msg)
	}
	return s.repo.Brand.Update(ctx, id, name, country)
}

func (s *CatalogService) AdminDeleteBrand(ctx context.Context, id uuid.UUID) error {
	return s.repo.Brand.Delete(ctx, id)
}

func (s *CatalogService) AdminListModels(ctx context.Context) ([]model.Model, error) {
	return s.repo.Model.GetAll(ctx)
}

func (s *CatalogService) AdminCreateModel(ctx context.Context, brandID uuid.UUID, name string, segment, description *string) (*model.Model, error) {
	var msg string
	name, msg = validate.ModelName(name)
	if msg != "" {
		return nil, apperr.BadRequest(msg)
	}
	segment, msg = validate.ModelSegment(segment)
	if msg != "" {
		return nil, apperr.BadRequest(msg)
	}
	description, msg = validate.ModelDescription(description)
	if msg != "" {
		return nil, apperr.BadRequest(msg)
	}
	return s.repo.Model.Create(ctx, brandID, name, segment, description)
}

func (s *CatalogService) AdminUpdateModel(ctx context.Context, id, brandID uuid.UUID, name string, segment, description *string) error {
	var msg string
	name, msg = validate.ModelName(name)
	if msg != "" {
		return apperr.BadRequest(msg)
	}
	segment, msg = validate.ModelSegment(segment)
	if msg != "" {
		return apperr.BadRequest(msg)
	}
	description, msg = validate.ModelDescription(description)
	if msg != "" {
		return apperr.BadRequest(msg)
	}
	return s.repo.Model.Update(ctx, id, brandID, name, segment, description)
}

func (s *CatalogService) AdminDeleteModel(ctx context.Context, id uuid.UUID) error {
	var oldKey string
	if s.store != nil {
		if key, _, err := s.repo.Model.GetImageMeta(ctx, id); err == nil {
			oldKey = key
		}
	}
	if err := s.repo.Model.Delete(ctx, id); err != nil {
		return err
	}
	if oldKey != "" && s.store != nil {
		_ = s.store.Remove(ctx, oldKey)
	}
	return nil
}

func (s *CatalogService) maxModelImageBytes() int64 {
	if s.maxUploadSize > 0 {
		return s.maxUploadSize
	}
	return 5 << 20
}

func (s *CatalogService) AdminUploadModelImage(ctx context.Context, modelID uuid.UUID, r io.Reader, fileName, mimeHint string) error {
	if s.store == nil {
		return apperr.Internal(fmt.Errorf("file storage is not configured"))
	}
	maxUpload := s.maxModelImageBytes()

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	limited := io.LimitReader(r, maxUpload+1)
	written, err := io.Copy(buf, limited)
	if err != nil {
		return apperr.BadRequest("failed to read image")
	}
	if written > maxUpload {
		return apperr.BadRequest("image too large")
	}
	if written == 0 {
		return apperr.BadRequest("image file is required")
	}

	data := buf.Bytes()
	mimeType := upload.ResolveImageMIME(data, mimeHint, fileName)
	if mimeType == "" {
		return apperr.BadRequest("only jpeg/png/webp images are allowed")
	}

	var oldKey string
	if key, _, err := s.repo.Model.GetImageMeta(ctx, modelID); err == nil {
		oldKey = key
	}

	key := uuid.NewString()
	if err := s.store.Store(ctx, key, bytes.NewReader(data), written); err != nil {
		return apperr.Internal(err)
	}
	if err := s.repo.Model.SetImageMeta(ctx, modelID, key, mimeType); err != nil {
		_ = s.store.Remove(ctx, key)
		return err
	}
	if oldKey != "" && oldKey != key {
		_ = s.store.Remove(ctx, oldKey)
	}
	return nil
}

// OpenModelImage streams a model image. etagKey is the storage key (stable cache validator).
func (s *CatalogService) OpenModelImage(ctx context.Context, modelID uuid.UUID) (io.ReadCloser, string, string, error) {
	if s.store == nil {
		return nil, "", "", apperr.Internal(fmt.Errorf("file storage is not configured"))
	}
	key, mimeType, err := s.repo.Model.GetImageMeta(ctx, modelID)
	if err != nil {
		return nil, "", "", err
	}
	rc, err := s.store.Open(ctx, key)
	if err != nil {
		return nil, "", "", apperr.NotFoundErr("Model image not found")
	}
	return rc, mimeType, key, nil
}
