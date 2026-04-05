package middleware

import (
	"net/http"
	"strings"
)

// LimitRequestBody wraps r.Body with http.MaxBytesReader for JSON and similar small bodies.
// Multipart uploads should use their own limits (e.g. documents).
func LimitRequestBody(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if maxBytes <= 0 || r.Body == nil {
				next.ServeHTTP(w, r)
				return
			}
			switch r.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
				next.ServeHTTP(w, r)
				return
			}
			ct := r.Header.Get("Content-Type")
			if strings.HasPrefix(ct, "multipart/form-data") {
				next.ServeHTTP(w, r)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
