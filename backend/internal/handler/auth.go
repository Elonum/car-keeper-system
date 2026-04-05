package handler

import (
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
	if !DecodeJSON(w, r, &in) {
		return
	}
	if msg := validateUserRegisterInput(&in); msg != "" {
		BadRequest(w, msg)
		return
	}

	user, err := h.services.Auth.Register(r.Context(), in)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	Success(w, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var login model.UserLogin
	if !DecodeJSON(w, r, &login) {
		return
	}

	token, user, err := h.services.Auth.Login(r.Context(), login)
	if err != nil {
		HandleError(w, r, err)
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
		HandleError(w, r, err)
		return
	}

	Success(w, user)
}

