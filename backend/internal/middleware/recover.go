package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// JSONRecover catches panics and returns a generic JSON 500 without leaking stack traces.
func JSONRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rv := recover(); rv != nil {
				slog.Error("panic recovered",
					"path", r.URL.Path,
					"method", r.Method,
					"request_id", chimw.GetReqID(r.Context()),
					"panic", rv,
				)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"success": false,
					"error":   "An unexpected error occurred. Please try again later.",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
