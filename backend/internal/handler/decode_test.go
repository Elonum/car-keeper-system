package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSON_InvalidBody(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid`))
	var dst map[string]any
	if DecodeJSON(rr, req, &dst) {
		t.Fatal("expected decode failure")
	}
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status %d", rr.Code)
	}
}

func TestDecodeJSON_EmptyBody(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	var dst map[string]any
	if DecodeJSON(rr, req, &dst) {
		t.Fatal("expected failure on empty body")
	}
}

func TestSuccessResponseShape(t *testing.T) {
	rr := httptest.NewRecorder()
	Success(rr, map[string]string{"ok": "true"})
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"success":true`) {
		t.Fatalf("body %s", rr.Body.String())
	}
}
