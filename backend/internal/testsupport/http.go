package testsupport

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type APIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error"`
}

func DoJSON(t *testing.T, handler http.Handler, method, path string, body any, token string) (*httptest.ResponseRecorder, APIResponse) {
	t.Helper()
	return DoJSONWithCookies(t, handler, method, path, body, token, nil)
}

func DoJSONWithCookies(t *testing.T, handler http.Handler, method, path string, body any, token string, cookies []*http.Cookie) (*httptest.ResponseRecorder, APIResponse) {
	t.Helper()
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		r = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, path, r)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var resp APIResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	return rr, resp
}

func SessionCookieFromResponse(t *testing.T, rr *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	for _, c := range rr.Result().Cookies() {
		if c.Name == "carkeeper_session" && c.Value != "" {
			return c
		}
	}
	return nil
}

func ParseDataMap(t *testing.T, raw json.RawMessage) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("parse data: %v body=%s", err, string(raw))
	}
	return m
}

func ParseDataArray(t *testing.T, raw json.RawMessage) []map[string]any {
	t.Helper()
	var arr []map[string]any
	if err := json.Unmarshal(raw, &arr); err != nil {
		t.Fatalf("parse data array: %v body=%s", err, string(raw))
	}
	return arr
}

func TokenFromLogin(t *testing.T, rr *httptest.ResponseRecorder, raw json.RawMessage) string {
	t.Helper()
	m := ParseDataMap(t, raw)
	if tok, _ := m["token"].(string); tok != "" {
		return tok
	}
	if c := SessionCookieFromResponse(t, rr); c != nil && c.Value != "" {
		return c.Value
	}
	t.Fatal("missing session token in login response")
	return ""
}
