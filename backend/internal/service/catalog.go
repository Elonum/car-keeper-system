package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/carkeeper/backend/internal/storage"
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

func (s *CatalogService) AdminListModels(ctx context.Context) ([]model.Model, error) {
	return s.repo.Model.GetAll(ctx)
}

func (s *CatalogService) AdminCreateModel(ctx context.Context, brandID uuid.UUID, name string, segment, description *string) (*model.Model, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apperr.BadRequest("name is required")
	}
	segment = normalizeOptional(segment, 100)
	description = normalizeOptional(description, 2000)
	return s.repo.Model.Create(ctx, brandID, name, segment, description)
}

func (s *CatalogService) AdminUpdateModel(ctx context.Context, id, brandID uuid.UUID, name string, segment, description *string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return apperr.BadRequest("name is required")
	}
	segment = normalizeOptional(segment, 100)
	description = normalizeOptional(description, 2000)
	return s.repo.Model.Update(ctx, id, brandID, name, segment, description)
}

func (s *CatalogService) AdminDeleteModel(ctx context.Context, id uuid.UUID) error {
	return s.repo.Model.Delete(ctx, id)
}

func (s *CatalogService) AdminUploadModelImage(ctx context.Context, modelID uuid.UUID, r io.Reader, size int64, fileName, mimeHint string) error {
	if s.store == nil {
		return apperr.Internal(fmt.Errorf("file storage is not configured"))
	}
	if size <= 0 {
		return apperr.BadRequest("image file is required")
	}
	maxUpload := s.maxUploadSize
	if maxUpload <= 0 {
		maxUpload = 5 << 20
	}
	if size > maxUpload {
		return apperr.BadRequest("image too large")
	}

	head := make([]byte, 512)
	n, err := io.ReadFull(r, head)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return apperr.BadRequest("failed to read image")
	}
	mimeType := normalizeImageMime(http.DetectContentType(head[:n]), mimeHint, fileName)
	if mimeType == "" {
		return apperr.BadRequest("only jpeg/png/webp images are allowed")
	}

	key := uuid.NewString()
	payload := io.MultiReader(bytes.NewReader(head[:n]), r)
	if err := s.store.Store(ctx, key, payload, size); err != nil {
		return apperr.Internal(err)
	}
	if err := s.repo.Model.SetImageMeta(ctx, modelID, key, mimeType); err != nil {
		_ = s.store.Remove(ctx, key)
		return err
	}
	return nil
}

func (s *CatalogService) OpenModelImage(ctx context.Context, modelID uuid.UUID) (io.ReadCloser, string, error) {
	if s.store == nil {
		return nil, "", apperr.Internal(fmt.Errorf("file storage is not configured"))
	}
	key, mimeType, err := s.repo.Model.GetImageMeta(ctx, modelID)
	if err != nil {
		return nil, "", err
	}
	rc, err := s.store.Open(ctx, key)
	if err != nil {
		return nil, "", apperr.NotFoundErr("Model image not found")
	}
	return rc, mimeType, nil
}

func normalizeOptional(v *string, maxLen int) *string {
	if v == nil {
		return nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}
	if maxLen > 0 {
		rs := []rune(s)
		if len(rs) > maxLen {
			s = string(rs[:maxLen])
		}
	}
	return &s
}

func normalizeImageMime(detected, hint, fileName string) string {
	candidate := strings.TrimSpace(strings.ToLower(detected))
	if isAllowedImageMime(candidate) {
		return candidate
	}
	candidate = strings.TrimSpace(strings.ToLower(strings.Split(hint, ";")[0]))
	if isAllowedImageMime(candidate) {
		return candidate
	}
	if ext := strings.ToLower(filepath.Ext(fileName)); ext != "" {
		byExt := mime.TypeByExtension(ext)
		byExt = strings.TrimSpace(strings.ToLower(strings.Split(byExt, ";")[0]))
		if isAllowedImageMime(byExt) {
			return byExt
		}
	}
	return ""
}

func isAllowedImageMime(m string) bool {
	return m == "image/jpeg" || m == "image/png" || m == "image/webp"
}
