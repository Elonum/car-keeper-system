package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/validate"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) CreateUserCar(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := RequesterAndRole(w, r)
	if !ok {
		return
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
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	userCarIDStr := chi.URLParam(r, "id")
	userCarID, err := uuid.Parse(userCarIDStr)
	if err != nil {
		BadRequest(w, "Invalid user car ID")
		return
	}

	userCar, err := h.services.Profile.GetUserCar(r.Context(), userCarID, requester, role)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			NotFound(w, "User car not found")
			return
		}
		NotFound(w, err.Error())
		return
	}
	Success(w, userCar)
}

func (h *Handler) GetUserConfigurations(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	configurations, err := h.services.Configurator.GetUserConfigurations(r.Context(), userID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, configurations)
}

// UpdateProfile patches first name, last name, and phone. Email is read-only (change requires a verified flow).
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	var body struct {
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
		Phone     *string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	fn, ln, msg := validate.Names(body.FirstName, body.LastName)
	if msg != "" {
		BadRequest(w, msg)
		return
	}
	ph, msg := validate.PhonePtr(body.Phone)
	if msg != "" {
		BadRequest(w, msg)
		return
	}

	user, err := h.services.Profile.UpdateProfile(r.Context(), userID, fn, ln, ph)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}
	Success(w, user)
}

// ChangePassword requires the current password (proves session holder knows credentials).
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}
	if body.CurrentPassword == "" {
		BadRequest(w, "current_password is required")
		return
	}

	if err := h.services.Auth.ChangePassword(r.Context(), userID, body.CurrentPassword, body.NewPassword); err != nil {
		BadRequest(w, err.Error())
		return
	}
	Success(w, map[string]string{"status": "ok"})
}

