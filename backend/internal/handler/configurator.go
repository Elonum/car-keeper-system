package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carkeeper/backend/internal/middleware"
	"github.com/carkeeper/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetColors(w http.ResponseWriter, r *http.Request) {
	var isAvailable *bool
	if isAvailableStr := r.URL.Query().Get("is_available"); isAvailableStr != "" {
		if val, err := strconv.ParseBool(isAvailableStr); err == nil {
			isAvailable = &val
		}
	}

	colors, err := h.services.Configurator.GetColors(r.Context(), isAvailable)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, colors)
}

func (h *Handler) GetOptions(w http.ResponseWriter, r *http.Request) {
	trimIDStr := r.URL.Query().Get("trim_id")
	if trimIDStr == "" {
		BadRequest(w, "trim_id is required")
		return
	}

	trimID, err := uuid.Parse(trimIDStr)
	if err != nil {
		BadRequest(w, "Invalid trim_id")
		return
	}

	options, err := h.services.Configurator.GetOptions(r.Context(), trimID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, options)
}

func (h *Handler) CreateConfiguration(w http.ResponseWriter, r *http.Request) {
	userIDVal, ok := middleware.GetUserIDUUID(r.Context())
	if !ok {
		Unauthorized(w, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		userIDStr, ok := userIDVal.(string)
		if !ok {
			Unauthorized(w, "Invalid user ID")
			return
		}
		var err error
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			Unauthorized(w, "Invalid user ID format")
			return
		}
	}

	var create model.ConfigurationCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	config, err := h.services.Configurator.CreateConfiguration(r.Context(), userID, create)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	Success(w, config)
}

func (h *Handler) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	configIDStr := chi.URLParam(r, "id")
	configID, err := uuid.Parse(configIDStr)
	if err != nil {
		BadRequest(w, "Invalid configuration ID")
		return
	}

	config, err := h.services.Configurator.GetConfiguration(r.Context(), configID)
	if err != nil {
		NotFound(w, err.Error())
		return
	}
	Success(w, config)
}

func (h *Handler) UpdateConfiguration(w http.ResponseWriter, r *http.Request) {
	configIDStr := chi.URLParam(r, "id")
	configID, err := uuid.Parse(configIDStr)
	if err != nil {
		BadRequest(w, "Invalid configuration ID")
		return
	}

	var update struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if err := h.services.Configurator.UpdateConfigurationStatus(r.Context(), configID, update.Status); err != nil {
		BadRequest(w, err.Error())
		return
	}

	config, err := h.services.Configurator.GetConfiguration(r.Context(), configID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, config)
}

func (h *Handler) DeleteConfiguration(w http.ResponseWriter, r *http.Request) {
	configIDStr := chi.URLParam(r, "id")
	configID, err := uuid.Parse(configIDStr)
	if err != nil {
		BadRequest(w, "Invalid configuration ID")
		return
	}

	if err := h.services.Configurator.DeleteConfiguration(r.Context(), configID); err != nil {
		InternalServerError(w, err.Error())
		return
	}

	Success(w, map[string]string{"message": "Configuration deleted"})
}

