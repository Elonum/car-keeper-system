package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/carkeeper/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetBrands(w http.ResponseWriter, r *http.Request) {
	brands, err := h.services.Catalog.GetBrands(r.Context())
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, brands)
}

func (h *Handler) GetModels(w http.ResponseWriter, r *http.Request) {
	brandIDStr := r.URL.Query().Get("brand_id")
	if brandIDStr == "" {
		BadRequest(w, "brand_id is required")
		return
	}

	brandID, err := uuid.Parse(brandIDStr)
	if err != nil {
		BadRequest(w, "Invalid brand_id")
		return
	}

	models, err := h.services.Catalog.GetModels(r.Context(), brandID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, models)
}

func (h *Handler) GetGenerations(w http.ResponseWriter, r *http.Request) {
	modelIDStr := r.URL.Query().Get("model_id")
	if modelIDStr == "" {
		BadRequest(w, "model_id is required")
		return
	}

	modelID, err := uuid.Parse(modelIDStr)
	if err != nil {
		BadRequest(w, "Invalid model_id")
		return
	}

	generations, err := h.services.Catalog.GetGenerations(r.Context(), modelID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, generations)
}

func (h *Handler) GetTrims(w http.ResponseWriter, r *http.Request) {
	filters := model.TrimFilters{}

	// Parse brand_id
	if brandIDStr := r.URL.Query().Get("brand_id"); brandIDStr != "" {
		for _, idStr := range strings.Split(brandIDStr, ",") {
			if id, err := uuid.Parse(strings.TrimSpace(idStr)); err == nil {
				filters.BrandID = append(filters.BrandID, id)
			}
		}
	}

	// Parse engine_type_id
	if engineTypeIDStr := r.URL.Query().Get("engine_type_id"); engineTypeIDStr != "" {
		for _, idStr := range strings.Split(engineTypeIDStr, ",") {
			if id, err := uuid.Parse(strings.TrimSpace(idStr)); err == nil {
				filters.EngineTypeID = append(filters.EngineTypeID, id)
			}
		}
	}

	// Parse transmission_id
	if transmissionIDStr := r.URL.Query().Get("transmission_id"); transmissionIDStr != "" {
		for _, idStr := range strings.Split(transmissionIDStr, ",") {
			if id, err := uuid.Parse(strings.TrimSpace(idStr)); err == nil {
				filters.TransmissionID = append(filters.TransmissionID, id)
			}
		}
	}

	// Parse drive_type_id
	if driveTypeIDStr := r.URL.Query().Get("drive_type_id"); driveTypeIDStr != "" {
		for _, idStr := range strings.Split(driveTypeIDStr, ",") {
			if id, err := uuid.Parse(strings.TrimSpace(idStr)); err == nil {
				filters.DriveTypeID = append(filters.DriveTypeID, id)
			}
		}
	}

	// Parse price range
	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &minPrice
		}
	}

	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &maxPrice
		}
	}

	// Parse is_available
	if isAvailableStr := r.URL.Query().Get("is_available"); isAvailableStr != "" {
		if isAvailable, err := strconv.ParseBool(isAvailableStr); err == nil {
			filters.IsAvailable = &isAvailable
		}
	}

	trims, err := h.services.Catalog.GetTrims(r.Context(), filters)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, trims)
}

func (h *Handler) GetTrim(w http.ResponseWriter, r *http.Request) {
	trimIDStr := chi.URLParam(r, "id")
	trimID, err := uuid.Parse(trimIDStr)
	if err != nil {
		BadRequest(w, "Invalid trim ID")
		return
	}

	trim, err := h.services.Catalog.GetTrim(r.Context(), trimID)
	if err != nil {
		NotFound(w, err.Error())
		return
	}
	Success(w, trim)
}

func (h *Handler) GetEngineTypes(w http.ResponseWriter, r *http.Request) {
	engineTypes, err := h.services.Catalog.GetEngineTypes(r.Context())
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, engineTypes)
}

func (h *Handler) GetTransmissions(w http.ResponseWriter, r *http.Request) {
	transmissions, err := h.services.Catalog.GetTransmissions(r.Context())
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, transmissions)
}

func (h *Handler) GetDriveTypes(w http.ResponseWriter, r *http.Request) {
	driveTypes, err := h.services.Catalog.GetDriveTypes(r.Context())
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, driveTypes)
}

