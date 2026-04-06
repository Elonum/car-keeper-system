package handler

import (
	"net/http"

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
		Name            string   `json:"name"`
		Category        string   `json:"category"`
		Description     *string  `json:"description"`
		Price           float64  `json:"price"`
		DurationMinutes *int     `json:"duration_minutes"`
		IsAvailable     bool     `json:"is_available"`
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
		Name            string   `json:"name"`
		Category        string   `json:"category"`
		Description     *string  `json:"description"`
		Price           float64  `json:"price"`
		DurationMinutes *int     `json:"duration_minutes"`
		IsAvailable     bool     `json:"is_available"`
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
