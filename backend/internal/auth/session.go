package auth

import (
	"net/http"
	"strings"
)

const SessionCookieName = "carkeeper_session"

// TokenFromRequest returns JWT from HttpOnly cookie or Authorization Bearer header.
func TokenFromRequest(r *http.Request) string {
	if c, err := r.Cookie(SessionCookieName); err == nil {
		if v := strings.TrimSpace(c.Value); v != "" {
			return v
		}
	}
	return bearerFromHeader(r.Header.Get("Authorization"))
}

func bearerFromHeader(authHeader string) string {
	parts := strings.SplitN(strings.TrimSpace(authHeader), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

// SetSessionCookie writes the JWT session cookie.
func SetSessionCookie(w http.ResponseWriter, token string, secure bool, maxAgeSeconds int) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAgeSeconds,
	})
}

// ClearSessionCookie removes the session cookie.
func ClearSessionCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
