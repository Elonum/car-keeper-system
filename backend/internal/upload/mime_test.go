package upload

import "testing"

func TestResolveDocumentMIME_fromMagicBytes(t *testing.T) {
	pdf := []byte("%PDF-1.4\n")
	got := ResolveDocumentMIME(pdf, "text/plain", "evil.exe")
	if got != "application/pdf" {
		t.Fatalf("got %q, want application/pdf", got)
	}
}

func TestResolveDocumentMIME_rejectsSpoofedHint(t *testing.T) {
	got := ResolveDocumentMIME([]byte("not a real file"), "application/x-msdownload", "file.bin")
	if got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}
