package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/carkeeper/backend/internal/apperr"
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
	userID, _, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	userCars, err := h.services.Profile.GetUserCars(r.Context(), userID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, userCars)
}

func (h *Handler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := RequesterAndRole(w, r)
	if !ok {
		return
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
	userID, _, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	appointments, err := h.services.Service.GetUserAppointments(r.Context(), userID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, appointments)
}

func (h *Handler) GetAppointment(w http.ResponseWriter, r *http.Request) {
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	appointmentIDStr := chi.URLParam(r, "id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		BadRequest(w, "Invalid appointment ID")
		return
	}

	appointment, err := h.services.Service.GetAppointment(r.Context(), appointmentID, requester, role)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			NotFound(w, "Appointment not found")
			return
		}
		NotFound(w, err.Error())
		return
	}
	Success(w, appointment)
}

func (h *Handler) CancelAppointment(w http.ResponseWriter, r *http.Request) {
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	appointmentIDStr := chi.URLParam(r, "id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		BadRequest(w, "Invalid appointment ID")
		return
	}

	if err := h.services.Service.CancelAppointment(r.Context(), appointmentID, requester, role); err != nil {
		switch {
		case errors.Is(err, apperr.ErrForbidden):
			Forbidden(w, "Not allowed to cancel this appointment")
		default:
			BadRequest(w, err.Error())
		}
		return
	}

	appointment, err := h.services.Service.GetAppointment(r.Context(), appointmentID, requester, role)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			NotFound(w, "Appointment not found")
			return
		}
		InternalServerError(w, err.Error())
		return
	}
	Success(w, appointment)
}

