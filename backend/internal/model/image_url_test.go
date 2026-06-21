package model

import (
	"testing"

	"github.com/google/uuid"
)

func TestBuildModelImageURL(t *testing.T) {
	id := uuid.MustParse("80000000-0000-0000-0000-000000000001")
	key := "abc-123"
	got := BuildModelImageURL(id, key)
	want := "/api/catalog/models/80000000-0000-0000-0000-000000000001/image?v=abc-123"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if BuildModelImageURL(id, "") != "" {
		t.Fatal("expected empty url for empty key")
	}
}
