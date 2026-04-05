package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/carkeeper/backend/internal/middleware"
	"github.com/carkeeper/backend/internal/model"
)

const (
	maxNameRunes     = 100 // matches varchar(100); count runes, not UTF-8 bytes
	maxEmailRunes    = 255
	maxPasswordRunes = 128
	maxPhoneRunes    = 30
)

func validateUserRegisterInput(in *model.UserRegisterInput) string {
	in.FirstName = strings.TrimSpace(in.FirstName)
	in.LastName = strings.TrimSpace(in.LastName)
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if in.FirstName == "" || utf8.RuneCountInString(in.FirstName) > maxNameRunes {
		return "first_name is required (max 100 characters)"
	}
	if in.LastName == "" || utf8.RuneCountInString(in.LastName) > maxNameRunes {
		return "last_name is required (max 100 characters)"
	}
	if in.Email == "" || utf8.RuneCountInString(in.Email) > maxEmailRunes {
		return "email is required (max 255 characters)"
	}
	if in.Phone != nil {
		p := strings.TrimSpace(*in.Phone)
		if p == "" {
			in.Phone = nil
		} else if utf8.RuneCountInString(p) > maxPhoneRunes {
			return "phone is too long (max 30 characters)"
		} else {
			in.Phone = &p
		}
	}
	pwLen := utf8.RuneCountInString(in.Password)
	if pwLen < 6 {
		return "password must be at least 6 characters"
	}
	if pwLen > maxPasswordRunes {
		return "password is too long (max 128 characters)"
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

