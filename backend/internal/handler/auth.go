package handler

import (
	"encoding/json"
	"net/http"

	"github.com/carkeeper/backend/internal/middleware"
	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/validate"
)

func validateUserRegisterInput(in *model.UserRegisterInput) string {
	fn, ln, msg := validate.Names(in.FirstName, in.LastName)
	if msg != "" {
		return msg
	}
	in.FirstName, in.LastName = fn, ln
	em, msg := validate.Email(in.Email)
	if msg != "" {
		return msg
	}
	in.Email = em
	ph, msg := validate.PhonePtr(in.Phone)
	if msg != "" {
		return msg
	}
	in.Phone = ph
	if msg := validate.NewPassword(in.Password); msg != "" {
		return msg
	}
	return ""
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var in model.UserRegisterInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}
	if msg := validateUserRegisterInput(&in); msg != "" {
		BadRequest(w, msg)
		return
	}

	user, err := h.services.Auth.Register(r.Context(), in)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	Success(w, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var login model.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	token, user, err := h.services.Auth.Login(r.Context(), login)
	if err != nil {
		Unauthorized(w, err.Error())
		return
	}

	Success(w, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		Unauthorized(w, "User not authenticated")
		return
	}

	user, err := h.services.Auth.GetUser(r.Context(), userID)
	if err != nil {
		NotFound(w, err.Error())
		return
	}

	Success(w, user)
}

