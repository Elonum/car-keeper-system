package validate

import (
	"strings"
	"testing"
)

func TestNewsFields_Valid(t *testing.T) {
	ti, co, msg := NewsFields(" Title ", " Body ")
	if msg != "" || ti != "Title" || co != "Body" {
		t.Fatalf("got %q %q msg %q", ti, co, msg)
	}
}

func TestNewsFields_EmptyContent(t *testing.T) {
	_, _, msg := NewsFields("Title", "   ")
	if msg != "content is required" {
		t.Fatalf("got %q", msg)
	}
}

func TestNewsFields_ContentTooLong(t *testing.T) {
	_, _, msg := NewsFields("Title", strings.Repeat("x", NewsContentMax+1))
	if msg != "content is too long" {
		t.Fatalf("got %q", msg)
	}
}
