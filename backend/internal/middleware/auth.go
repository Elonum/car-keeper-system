package middleware

import (
	"context"
	"net/http"

	"github.com/carkeeper/backend/internal/auth"
	"github.com/carkeeper/backend/internal/service"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserRoleKey contextKey = "role"

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := auth.TokenFromRequest(r)
			if token == "" {
				writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			claims, err := authService.AuthenticateRequest(r.Context(), token)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDUUID(ctx context.Context) (interface{}, bool) {
	userID := ctx.Value(UserIDKey)
	if userID == nil {
		return nil, false
	}
	return userID, true
}

func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}

