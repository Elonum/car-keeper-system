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
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var resp APIResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	return rr, resp
}

func ParseDataMap(t *testing.T, raw json.RawMessage) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("parse data: %v body=%s", err, string(raw))
	}
	return m
}

func TokenFromLogin(t *testing.T, raw json.RawMessage) string {
	t.Helper()
	m := ParseDataMap(t, raw)
	tok, _ := m["token"].(string)
	if tok == "" {
		t.Fatal("missing token in login response")
	}
	return tok
}
