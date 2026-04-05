package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/carkeeper/backend/internal/apperr"
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
		HandleError(w, r, err)
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
		HandleError(w, r, err)
		return
	}
	Success(w, options)
}

func (h *Handler) CreateConfiguration(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	var createRequest struct {
		TrimID    string   `json:"trim_id"`
		ColorID   string   `json:"color_id"`
		OptionIDs []string `json:"option_ids"`
	}
	if !DecodeJSON(w, r, &createRequest) {
		return
	}

	trimID, err := uuid.Parse(createRequest.TrimID)
	if err != nil {
		BadRequest(w, "Invalid trim_id format")
		return
	}

	colorID, err := uuid.Parse(createRequest.ColorID)
	if err != nil {
		BadRequest(w, "Invalid color_id format")
		return
	}

	var optionIDs []uuid.UUID
	for _, optIDStr := range createRequest.OptionIDs {
		optID, err := uuid.Parse(optIDStr)
		if err != nil {
			BadRequest(w, "Invalid option_id format")
			return
		}
		optionIDs = append(optionIDs, optID)
	}

	create := model.ConfigurationCreate{
		TrimID:    trimID,
		ColorID:   colorID,
		OptionIDs: optionIDs,
	}

	config, err := h.services.Configurator.CreateConfiguration(r.Context(), userID, create)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	Success(w, config)
}

func (h *Handler) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	configIDStr := chi.URLParam(r, "id")
	configID, err := uuid.Parse(configIDStr)
	if err != nil {
		BadRequest(w, "Invalid configuration ID")
		return
	}

	config, err := h.services.Configurator.GetConfiguration(r.Context(), configID, requester, role)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, config)
}

func (h *Handler) UpdateConfiguration(w http.ResponseWriter, r *http.Request) {
	userID, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	configIDStr := chi.URLParam(r, "id")
	configID, err := uuid.Parse(configIDStr)
	if err != nil {
		BadRequest(w, "Invalid configuration ID")
		return
	}

	var updateRequest struct {
		Status    *string  `json:"status,omitempty"`
		TrimID    *string  `json:"trim_id,omitempty"`
		ColorID   *string  `json:"color_id,omitempty"`
		OptionIDs []string `json:"option_ids,omitempty"`
	}
	if !DecodeJSON(w, r, &updateRequest) {
		return
	}

	hasStatus := updateRequest.Status != nil && *updateRequest.Status != ""
	hasTrimOrColor := (updateRequest.TrimID != nil && *updateRequest.TrimID != "") ||
		(updateRequest.ColorID != nil && *updateRequest.ColorID != "")

	isStatusOnlyUpdate := hasStatus && !hasTrimOrColor

	if updateRequest.Status != nil && *updateRequest.Status == "" {
		hasStatus = false
		isStatusOnlyUpdate = false
		updateRequest.Status = nil
	}

	if isStatusOnlyUpdate {
		if err := h.services.Configurator.UpdateConfigurationStatus(r.Context(), configID, userID, role, *updateRequest.Status); err != nil {
			if errors.Is(err, apperr.ErrForbidden) {
				Forbidden(w, "Not allowed to update this configuration status")
				return
			}
			HandleError(w, r, err)
			return
		}
		config, err := h.services.Configurator.GetConfiguration(r.Context(), configID, userID, role)
		if err != nil {
			HandleError(w, r, err)
			return
		}
		Success(w, config)
		return
	}

	if updateRequest.TrimID == nil || *updateRequest.TrimID == "" {
		BadRequest(w, "trim_id is required for full update")
		return
	}
	if updateRequest.ColorID == nil || *updateRequest.ColorID == "" {
		BadRequest(w, "color_id is required for full update")
		return
	}

	trimID, err := uuid.Parse(*updateRequest.TrimID)
	if err != nil {
		BadRequest(w, "Invalid trim_id format")
		return
	}

	colorID, err := uuid.Parse(*updateRequest.ColorID)
	if err != nil {
		BadRequest(w, "Invalid color_id format")
		return
	}

	var optionIDs []uuid.UUID
	if updateRequest.OptionIDs != nil {
		for _, optIDStr := range updateRequest.OptionIDs {
			if optIDStr == "" {
				continue
			}
			optID, err := uuid.Parse(optIDStr)
			if err != nil {
				BadRequest(w, "Invalid option_id format")
				return
			}
			optionIDs = append(optionIDs, optID)
		}
	}

	create := model.ConfigurationCreate{
		TrimID:    trimID,
		ColorID:   colorID,
		OptionIDs: optionIDs,
	}

	config, err := h.services.Configurator.UpdateConfiguration(r.Context(), configID, userID, role, create)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	Success(w, config)
}

func (h *Handler) DeleteConfiguration(w http.ResponseWriter, r *http.Request) {
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	configIDStr := chi.URLParam(r, "id")
	configID, err := uuid.Parse(configIDStr)
	if err != nil {
		BadRequest(w, "Invalid configuration ID")
		return
	}

	if err := h.services.Configurator.DeleteConfiguration(r.Context(), configID, requester, role); err != nil {
		HandleError(w, r, err)
		return
	}

	Success(w, map[string]string{"message": "Configuration deleted"})
}
