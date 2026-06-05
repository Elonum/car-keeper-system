package upload

import "testing"

func TestResolveImageMIME_FromMagic(t *testing.T) {
	pngHead := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	got := ResolveImageMIME(pngHead, "", "photo.bin")
	if got != "image/png" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveImageMIME_RejectsUnknown(t *testing.T) {
	got := ResolveImageMIME([]byte("not-image"), "application/pdf", "x.pdf")
	if got != "" {
		t.Fatalf("got %q", got)
	}
}
