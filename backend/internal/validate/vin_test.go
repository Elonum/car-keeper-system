package validate

import "testing"

func TestVIN_Valid(t *testing.T) {
	if msg := VIN("1HGBH41JXMN109186"); msg != "" {
		t.Fatalf("expected valid VIN, got %q", msg)
	}
}

func TestVIN_InvalidLength(t *testing.T) {
	if msg := VIN("SHORT"); msg == "" {
		t.Fatal("expected error for short VIN")
	}
}

func TestVIN_ForbiddenLetters(t *testing.T) {
	if msg := VIN("1HGBH41JXMN10918I"); msg == "" {
		t.Fatal("expected error for letter I in VIN")
	}
}

func TestNormalizeVIN(t *testing.T) {
	got := NormalizeVIN("  abcdefgh1jklmn234  ")
	if got != "ABCDEFGH1JKLMN234" {
		t.Fatalf("got %q", got)
	}
}
