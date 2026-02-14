package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/carkeeper/backend/internal/middleware"
	"github.com/carkeeper/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetServiceTypes(w http.ResponseWriter, r *http.Request) {
	var category *string
	var isAvailable *bool

	if cat := r.URL.Query().Get("category"); cat != "" {
		category = &cat
	}

	if isAvailableStr := r.URL.Query().Get("is_available"); isAvailableStr != "" {
		if val, err := strconv.ParseBool(isAvailableStr); err == nil {
			isAvailable = &val
		}
	}

	serviceTypes, err := h.services.Service.GetServiceTypes(r.Context(), category, isAvailable)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, serviceTypes)
}

func (h *Handler) GetBranches(w http.ResponseWriter, r *http.Request) {
	var isActive *bool
	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		if val, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &val
		}
	}

	branches, err := h.services.Service.GetBranches(r.Context(), isActive)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, branches)
}

func (h *Handler) GetUserCars(w http.ResponseWriter, r *http.Request) {
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

	userCars, err := h.services.Profile.GetUserCars(r.Context(), userID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, userCars)
}

func (h *Handler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
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

	var create model.ServiceAppointmentCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	// Parse appointment_date if it's a string
	if create.AppointmentDate.IsZero() {
		if dateStr := r.URL.Query().Get("appointment_date"); dateStr != "" {
			if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
				create.AppointmentDate = t
			}
		}
	}

	appointment, err := h.services.Service.CreateAppointment(r.Context(), userID, create)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	Success(w, appointment)
}

func (h *Handler) GetUserAppointments(w http.ResponseWriter, r *http.Request) {
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

	appointments, err := h.services.Service.GetUserAppointments(r.Context(), userID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, appointments)
}

func (h *Handler) GetAppointment(w http.ResponseWriter, r *http.Request) {
	appointmentIDStr := chi.URLParam(r, "id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		BadRequest(w, "Invalid appointment ID")
		return
	}

	appointment, err := h.services.Service.GetAppointment(r.Context(), appointmentID)
	if err != nil {
		NotFound(w, err.Error())
		return
	}
	Success(w, appointment)
}

func (h *Handler) CancelAppointment(w http.ResponseWriter, r *http.Request) {
	appointmentIDStr := chi.URLParam(r, "id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		BadRequest(w, "Invalid appointment ID")
		return
	}

	if err := h.services.Service.CancelAppointment(r.Context(), appointmentID); err != nil {
		BadRequest(w, err.Error())
		return
	}

	appointment, err := h.services.Service.GetAppointment(r.Context(), appointmentID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, appointment)
}

