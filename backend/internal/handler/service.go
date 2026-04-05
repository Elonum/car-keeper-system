package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

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

func (h *Handler) GetBranchAvailability(w http.ResponseWriter, r *http.Request) {
	_, _, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}
	branchIDStr := chi.URLParam(r, "branchID")
	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		BadRequest(w, "invalid branch id")
		return
	}
	dateStr := strings.TrimSpace(r.URL.Query().Get("date"))
	if dateStr == "" {
		BadRequest(w, "date is required")
		return
	}
	raw := strings.TrimSpace(r.URL.Query().Get("service_type_ids"))
	if raw == "" {
		BadRequest(w, "service_type_ids is required")
		return
	}
	parts := strings.Split(raw, ",")
	var ids []uuid.UUID
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := uuid.Parse(p)
		if err != nil {
			BadRequest(w, "invalid service_type_ids")
			return
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		BadRequest(w, "service_type_ids is required")
		return
	}
	avail, err := h.services.Service.BranchAvailability(r.Context(), branchID, dateStr, ids)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}
	Success(w, avail)
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

	if create.AppointmentDate.IsZero() {
		BadRequest(w, "appointment_date is required")
		return
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

