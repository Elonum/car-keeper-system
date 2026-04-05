package apperr

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// APIError is a safe-to-serialize HTTP error: Msg is shown to the client; Cause is for logs only.
type APIError struct {
	Status int
	Msg    string
	Code   string
	Cause  error
}

func (e *APIError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "api error"
}

func (e *APIError) Unwrap() error { return e.Cause }

func New(status int, msg string, cause error) *APIError {
	return &APIError{Status: status, Msg: msg, Cause: cause}
}

func BadRequest(msg string) *APIError {
	return New(http.StatusBadRequest, msg, nil)
}

func Unauthorized(msg string) *APIError {
	return New(http.StatusUnauthorized, msg, nil)
}

func Forbidden(msg string) *APIError {
	return New(http.StatusForbidden, msg, nil)
}

func NotFoundErr(msg string) *APIError {
	return New(http.StatusNotFound, msg, nil)
}

func Conflict(msg string) *APIError {
	return New(http.StatusConflict, msg, nil)
}

func PayloadTooLarge(msg string) *APIError {
	return New(http.StatusRequestEntityTooLarge, msg, nil)
}

// Internal wraps an unexpected failure; Msg is always generic for the client.
func Internal(cause error) *APIError {
	return &APIError{
		Status: http.StatusInternalServerError,
		Msg:    "An unexpected error occurred. Please try again later.",
		Cause:  cause,
	}
}

// Wrap returns a client-facing error unless err is already an *APIError.
func Wrap(err error, status int, public string) error {
	if err == nil {
		return nil
	}
	var ae *APIError
	if errors.As(err, &ae) {
		return err
	}
	return New(status, public, err)
}
