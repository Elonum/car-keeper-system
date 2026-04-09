package handler

import (
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/carkeeper/backend/internal/authz"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) AdminCreateBrand(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermCatalogManage); !ok {
		return
	}
	var body struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	}
	if !DecodeJSON(w, r, &body) {
		return
	}
	b, err := h.services.Catalog.AdminCreateBrand(r.Context(), body.Name, body.Country)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	JSON(w, http.StatusCreated, Response{Success: true, Data: b})
}

func (h *Handler) AdminUpdateBrand(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermCatalogManage); !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "Invalid brand ID")
		return
	}
	var body struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	}
	if !DecodeJSON(w, r, &body) {
		return
	}
	if err := h.services.Catalog.AdminUpdateBrand(r.Context(), id, body.Name, body.Country); err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, map[string]string{"message": "updated"})
}

func (h *Handler) AdminDeleteBrand(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermCatalogManage); !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "Invalid brand ID")
		return
	}
	if err := h.services.Catalog.AdminDeleteBrand(r.Context(), id); err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, map[string]string{"message": "deleted"})
}

func (h *Handler) AdminCreateServiceType(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermServiceManage); !ok {
		return
	}
	var body struct {
		Name            string  `json:"name"`
		Category        string  `json:"category"`
		Description     *string `json:"description"`
		Price           float64 `json:"price"`
		DurationMinutes *int    `json:"duration_minutes"`
		IsAvailable     bool    `json:"is_available"`
	}
	if !DecodeJSON(w, r, &body) {
		return
	}
	st, err := h.services.Service.AdminCreateServiceType(r.Context(), body.Name, body.Category, body.Description, body.Price, body.DurationMinutes, body.IsAvailable)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	JSON(w, http.StatusCreated, Response{Success: true, Data: st})
}

func (h *Handler) AdminUpdateServiceType(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermServiceManage); !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "Invalid service type ID")
		return
	}
	var body struct {
		Name            string  `json:"name"`
		Category        string  `json:"category"`
		Description     *string `json:"description"`
		Price           float64 `json:"price"`
		DurationMinutes *int    `json:"duration_minutes"`
		IsAvailable     bool    `json:"is_available"`
	}
	if !DecodeJSON(w, r, &body) {
		return
	}
	if err := h.services.Service.AdminUpdateServiceType(r.Context(), id, body.Name, body.Category, body.Description, body.Price, body.DurationMinutes, body.IsAvailable); err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, map[string]string{"message": "updated"})
}

func (h *Handler) AdminDeleteServiceType(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermServiceManage); !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "Invalid service type ID")
		return
	}
	if err := h.services.Service.AdminDeleteServiceType(r.Context(), id); err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, map[string]string{"message": "deleted"})
}

func (h *Handler) AdminUpdateBranch(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermServiceManage); !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "Invalid branch ID")
		return
	}
	var body struct {
		Name     *string `json:"name"`
		Address  *string `json:"address"`
		Phone    *string `json:"phone"`
		Email    *string `json:"email"`
		IsActive *bool   `json:"is_active"`
	}
	if !DecodeJSON(w, r, &body) {
		return
	}
	if err := h.services.Service.AdminUpdateBranch(r.Context(), id, body.Name, body.Address, body.Phone, body.Email, body.IsActive); err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, map[string]string{"message": "updated"})
}

func (h *Handler) AdminListModels(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermCatalogManage); !ok {
		return
	}
	items, err := h.services.Catalog.AdminListModels(r.Context())
	if err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, items)
}

func (h *Handler) AdminCreateModel(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermCatalogManage); !ok {
		return
	}
	var body struct {
		BrandID     string  `json:"brand_id"`
		Name        string  `json:"name"`
		Segment     *string `json:"segment"`
		Description *string `json:"description"`
	}
	if !DecodeJSON(w, r, &body) {
		return
	}
	brandID, err := uuid.Parse(strings.TrimSpace(body.BrandID))
	if err != nil {
		BadRequest(w, "Invalid brand ID")
		return
	}
	item, err := h.services.Catalog.AdminCreateModel(r.Context(), brandID, body.Name, body.Segment, body.Description)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	JSON(w, http.StatusCreated, Response{Success: true, Data: item})
}

func (h *Handler) AdminUpdateModel(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermCatalogManage); !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "Invalid model ID")
		return
	}
	var body struct {
		BrandID     string  `json:"brand_id"`
		Name        string  `json:"name"`
		Segment     *string `json:"segment"`
		Description *string `json:"description"`
	}
	if !DecodeJSON(w, r, &body) {
		return
	}
	brandID, err := uuid.Parse(strings.TrimSpace(body.BrandID))
	if err != nil {
		BadRequest(w, "Invalid brand ID")
		return
	}
	if err := h.services.Catalog.AdminUpdateModel(r.Context(), id, brandID, body.Name, body.Segment, body.Description); err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, map[string]string{"message": "updated"})
}

func (h *Handler) AdminDeleteModel(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermCatalogManage); !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "Invalid model ID")
		return
	}
	if err := h.services.Catalog.AdminDeleteModel(r.Context(), id); err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, map[string]string{"message": "deleted"})
}

func (h *Handler) AdminUploadModelImage(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermCatalogManage); !ok {
		return
	}
	modelID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "Invalid model ID")
		return
	}
	maxB := h.cfg.Storage.MaxUploadBytes
	if maxB <= 0 {
		maxB = 5 << 20
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxB+(8<<20))
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		BadRequest(w, "invalid multipart form")
		return
	}
	defer func() { _ = r.MultipartForm.RemoveAll() }()

	file, header, err := r.FormFile("file")
	if err != nil {
		BadRequest(w, "file field is required")
		return
	}
	defer file.Close()

	if header.Size <= 0 {
		BadRequest(w, "could not determine image size")
		return
	}
	if err := h.services.Catalog.AdminUploadModelImage(r.Context(), modelID, file, header.Size, header.Filename, header.Header.Get("Content-Type")); err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, map[string]string{"message": "uploaded"})
}

func (h *Handler) GetModelImage(w http.ResponseWriter, r *http.Request) {
	modelID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		BadRequest(w, "Invalid model ID")
		return
	}
	rc, mimeType, err := h.services.Catalog.OpenModelImage(r.Context(), modelID)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	defer rc.Close()
	if strings.TrimSpace(mimeType) == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("inline", map[string]string{"filename": "model-image"}))
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, rc)
}
