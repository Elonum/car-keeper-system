package upload

import "testing"

func TestResolveImageMIME_FromMagic(t *testing.T) {
	pngHead := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	got := ResolveImageMIME(pngHead, "", "photo.bin")
	if got != "image/png" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveImageMIME_WebP(t *testing.T) {
	head := append([]byte("RIFFxxxx"), []byte("WEBP")...)
	head = append(head, make([]byte, 4)...)
	got := ResolveImageMIME(head, "", "photo.webp")
	if got != "image/webp" {
		t.Fatalf("got %q, want image/webp", got)
	}
}

func TestResolveImageMIME_JPEG(t *testing.T) {
	head := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}
	got := ResolveImageMIME(head, "", "photo.jpg")
	if got != "image/jpeg" {
		t.Fatalf("got %q, want image/jpeg", got)
	}
}

func TestResolveImageMIME_RejectsUnknown(t *testing.T) {
	got := ResolveImageMIME([]byte("not-image"), "application/pdf", "x.pdf")
	if got != "" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveImageMIME_RejectsExtensionOnly(t *testing.T) {
	got := ResolveImageMIME(nil, "image/png", "photo.png")
	if got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}
