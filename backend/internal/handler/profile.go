package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/carkeeper/backend/internal/middleware"
	"github.com/carkeeper/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) CreateUserCar(w http.ResponseWriter, r *http.Request) {
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
		TrimID        string  `json:"trim_id"`
		ColorID       string  `json:"color_id"`
		VIN           string  `json:"vin"`
		Year          int     `json:"year"`
		CurrentMileage int    `json:"current_mileage"`
		PurchaseDate  *string `json:"purchase_date,omitempty"`
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

	create := model.UserCarCreate{
		TrimID:        trimID,
		ColorID:       colorID,
		VIN:           createRequest.VIN,
		Year:          createRequest.Year,
		CurrentMileage: createRequest.CurrentMileage,
	}

	if createRequest.PurchaseDate != nil && *createRequest.PurchaseDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *createRequest.PurchaseDate)
		if err != nil {
			BadRequest(w, "Invalid purchase_date format. Use YYYY-MM-DD")
			return
		}
		create.PurchaseDate = &parsedDate
	}

	userCar, err := h.services.Profile.CreateUserCar(r.Context(), userID, create)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	Success(w, userCar)
}

func (h *Handler) GetUserCar(w http.ResponseWriter, r *http.Request) {
	userCarIDStr := chi.URLParam(r, "id")
	userCarID, err := uuid.Parse(userCarIDStr)
	if err != nil {
		BadRequest(w, "Invalid user car ID")
		return
	}

	userCar, err := h.services.Profile.GetUserCar(r.Context(), userCarID)
	if err != nil {
		NotFound(w, err.Error())
		return
	}
	Success(w, userCar)
}

func (h *Handler) GetUserConfigurations(w http.ResponseWriter, r *http.Request) {
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

	configurations, err := h.services.Configurator.GetUserConfigurations(r.Context(), userID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, configurations)
}

