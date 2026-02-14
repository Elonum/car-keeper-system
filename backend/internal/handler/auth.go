package handler

import (
	"encoding/json"
	"net/http"

	"github.com/carkeeper/backend/internal/middleware"
	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var create model.UserCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	user, err := h.services.Auth.Register(r.Context(), create)
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
	userIDVal, ok := middleware.GetUserIDUUID(r.Context())
	if !ok {
		Unauthorized(w, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		// Try to parse as string
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

	user, err := h.services.Auth.GetUser(r.Context(), userID)
	if err != nil {
		NotFound(w, err.Error())
		return
	}

	Success(w, user)
}

