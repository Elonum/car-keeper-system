package middleware

import (
	"context"
	"net/http"

	"github.com/carkeeper/backend/internal/auth"
	"github.com/carkeeper/backend/internal/service"
	"github.com/google/uuid"
)

// UserIDFromContext returns the authenticated user id set by Auth or OptionalAuth middleware.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	v, ok := GetUserIDUUID(ctx)
	if !ok || v == nil {
		return uuid.Nil, false
	}
	switch id := v.(type) {
	case uuid.UUID:
		if id == uuid.Nil {
			return uuid.Nil, false
		}
		return id, true
	case string:
		parsed, err := uuid.Parse(id)
		if err != nil {
			return uuid.Nil, false
		}
		return parsed, true
	default:
		return uuid.Nil, false
	}
}

// OptionalAuthMiddleware attaches user_id and role when a valid Bearer token is present.
// Invalid or missing tokens are ignored (anonymous request).
func OptionalAuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := auth.TokenFromRequest(r)
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}
			claims, err := authService.ValidateToken(token)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
