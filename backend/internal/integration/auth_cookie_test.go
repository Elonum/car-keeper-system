package integration_test

import (
	"net/http"
	"testing"

	"github.com/carkeeper/backend/internal/auth"
	"github.com/carkeeper/backend/internal/testsupport"
)

func TestAuth_LoginSetsSessionCookie_AndMeWorksWithoutBearer(t *testing.T) {
	email, pass := registerFreshCustomerCredentials(t)
	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    email,
		"password": pass,
	}, "")
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("login status=%d err=%s", rr.Code, resp.Error)
	}
	cookie := testsupport.SessionCookieFromResponse(t, rr)
	if cookie == nil {
		t.Fatal("expected carkeeper_session cookie")
	}
	if cookie.Name != auth.SessionCookieName {
		t.Fatalf("cookie name=%s", cookie.Name)
	}
	if cookie.HttpOnly == false {
		t.Fatal("session cookie must be HttpOnly")
	}

	rr, resp = testsupport.DoJSONWithCookies(t, testHandler, http.MethodGet, "/api/auth/me", nil, "", []*http.Cookie{cookie})
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("me with cookie status=%d err=%s", rr.Code, resp.Error)
	}
	me := testsupport.ParseDataMap(t, resp.Data)
	if me["email"] != email {
		t.Fatalf("email=%v", me["email"])
	}
}

func TestAuth_LogoutClearsSession(t *testing.T) {
	token := loginSeedUser(t, "customer@carkeeper.ru")
	rr, _ := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/auth/logout", nil, token)
	if rr.Code != http.StatusOK {
		t.Fatalf("logout status=%d", rr.Code)
	}
	cleared := false
	for _, c := range rr.Result().Cookies() {
		if c.Name == auth.SessionCookieName && (c.MaxAge < 0 || c.Value == "") {
			cleared = true
		}
	}
	if !cleared {
		t.Fatal("expected cleared session cookie")
	}
}
