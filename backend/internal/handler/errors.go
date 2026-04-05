package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/carkeeper/backend/internal/apperr"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// HandleError writes a JSON error, logs internal details, and never exposes err.Error() for unknown errors.
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	var ae *apperr.APIError
	if errors.As(err, &ae) {
		if ae.Status >= http.StatusInternalServerError && ae.Cause != nil {
			slog.Error("server error",
				"status", ae.Status,
				"public_msg", ae.Msg,
				"path", r.URL.Path,
				"method", r.Method,
				"request_id", chimw.GetReqID(r.Context()),
				"err", ae.Cause,
			)
		} else if ae.Cause != nil {
			slog.Debug("client error",
				"status", ae.Status,
				"public_msg", ae.Msg,
				"path", r.URL.Path,
				"method", r.Method,
				"request_id", chimw.GetReqID(r.Context()),
				"err", ae.Cause,
			)
		}
		Error(w, ae.Status, ae.Msg)
		return
	}

	if errors.Is(err, apperr.ErrNotFound) {
		NotFound(w, "Resource not found")
		return
	}
	if errors.Is(err, apperr.ErrForbidden) {
		Forbidden(w, "Access denied")
		return
	}
	if errors.Is(err, apperr.ErrInvalidCredentials) {
		Unauthorized(w, "Invalid email or password")
		return
	}

	slog.Error("unhandled error",
		"path", r.URL.Path,
		"method", r.Method,
		"request_id", chimw.GetReqID(r.Context()),
		"err", err,
	)
	InternalServerError(w, "An unexpected error occurred. Please try again later.")
}

// DecodeJSON decodes the body without leaking parser internals to the client.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if r.Body == nil {
		BadRequest(w, "Request body is required")
		return false
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		var mb *http.MaxBytesError
		if errors.As(err, &mb) {
			Error(w, http.StatusRequestEntityTooLarge, "Request body is too large")
			return false
		}
		if errors.Is(err, io.EOF) {
			BadRequest(w, "Request body is required")
			return false
		}
		BadRequest(w, "Invalid JSON")
		return false
	}
	return true
}
