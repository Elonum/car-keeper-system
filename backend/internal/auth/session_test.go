package auth

import (
	"net/http"
	"testing"
)

func TestTokenFromRequest_CookiePreferred(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "from-cookie"})
	r.Header.Set("Authorization", "Bearer from-header")
	if got := TokenFromRequest(r); got != "from-cookie" {
		t.Fatalf("got %q", got)
	}
}

func TestTokenFromRequest_BearerFallback(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Authorization", "Bearer from-header")
	if got := TokenFromRequest(r); got != "from-header" {
		t.Fatalf("got %q", got)
	}
}
