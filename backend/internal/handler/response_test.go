package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccess_nilSliceEncodesAsEmptyArray(t *testing.T) {
	rr := httptest.NewRecorder()
	var items []string
	Success(rr, items)

	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}

	var body Response
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !body.Success {
		t.Fatal("expected success")
	}
	raw, err := json.Marshal(body.Data)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "[]" {
		t.Fatalf("got %s, want []", string(raw))
	}
}
