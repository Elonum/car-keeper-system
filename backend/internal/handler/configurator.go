package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	var createRequest struct {
		TrimID    string   `json:"trim_id"`
		ColorID   string   `json:"color_id"`
		OptionIDs []string `json:"option_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		BadRequest(w, "Invalid request body: "+err.Error())
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
			BadRequest(w, "Invalid option_id format: "+optIDStr)
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

	configIDStr := chi.URLParam(r, "id")
	configID, err := uuid.Parse(configIDStr)
	if err != nil {
		BadRequest(w, "Invalid configuration ID")
		return
	}

	var updateRequest struct {
		Status    *string   `json:"status,omitempty"`
		TrimID    *string   `json:"trim_id,omitempty"`
		ColorID   *string   `json:"color_id,omitempty"`
		OptionIDs []string  `json:"option_ids,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		BadRequest(w, "Invalid request body: "+err.Error())
		return
	}

	// Log the received request for debugging
	log.Printf("[UpdateConfiguration] Received request - ConfigID: %s, Status: %v, TrimID: %v, ColorID: %v, OptionIDs count: %d",
		configID.String(),
		updateRequest.Status,
		updateRequest.TrimID,
		updateRequest.ColorID,
		len(updateRequest.OptionIDs))

	// Check if this is a status-only update
	// Status-only: status is provided (not empty) AND no trim_id/color_id provided
	hasStatus := updateRequest.Status != nil && *updateRequest.Status != ""
	hasTrimOrColor := (updateRequest.TrimID != nil && *updateRequest.TrimID != "") || 
	                  (updateRequest.ColorID != nil && *updateRequest.ColorID != "")
	
	isStatusOnlyUpdate := hasStatus && !hasTrimOrColor
	
	log.Printf("[UpdateConfiguration] Analysis - hasStatus: %v, hasTrimOrColor: %v, isStatusOnlyUpdate: %v",
		hasStatus, hasTrimOrColor, isStatusOnlyUpdate)
	
	// If status is provided but empty string, treat it as not provided
	if updateRequest.Status != nil && *updateRequest.Status == "" {
		log.Printf("[UpdateConfiguration] Empty status detected, treating as not provided")
		hasStatus = false
		isStatusOnlyUpdate = false
		// Set status to nil to prevent any further processing
		updateRequest.Status = nil
	}
	
	if isStatusOnlyUpdate {
		log.Printf("[UpdateConfiguration] Processing status-only update with status: %s", *updateRequest.Status)
		if err := h.services.Configurator.UpdateConfigurationStatus(r.Context(), configID, *updateRequest.Status); err != nil {
			log.Printf("[UpdateConfiguration] Status update error: %v", err)
			BadRequest(w, err.Error())
			return
		}
		config, err := h.services.Configurator.GetConfiguration(r.Context(), configID)
		if err != nil {
			InternalServerError(w, err.Error())
			return
		}
		Success(w, config)
		return
	}

	// Full update - trim_id and color_id are required
	// Note: Status field is completely ignored in full updates (even if provided or empty)
	// Configuration keeps its current status
	log.Printf("[UpdateConfiguration] Processing full update - Status will be ignored")
	if updateRequest.TrimID == nil || *updateRequest.TrimID == "" {
		log.Printf("[UpdateConfiguration] Error: trim_id is required for full update")
		BadRequest(w, "trim_id is required for full update")
		return
	}
	if updateRequest.ColorID == nil || *updateRequest.ColorID == "" {
		log.Printf("[UpdateConfiguration] Error: color_id is required for full update")
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
				continue // Skip empty strings
			}
			optID, err := uuid.Parse(optIDStr)
			if err != nil {
				BadRequest(w, "Invalid option_id format: "+optIDStr)
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

	config, err := h.services.Configurator.UpdateConfiguration(r.Context(), configID, userID, create)
	if err != nil {
		BadRequest(w, err.Error())
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
		// Check if it's a business logic error (status check)
		if strings.Contains(err.Error(), "can only delete draft") {
			BadRequest(w, err.Error())
			return
		}
		InternalServerError(w, err.Error())
		return
	}

	Success(w, map[string]string{"message": "Configuration deleted"})
}

