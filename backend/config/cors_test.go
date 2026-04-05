package config

import "testing"

func TestParseCSVOrigins(t *testing.T) {
	if len(parseCSVOrigins("")) != 0 {
		t.Fatal("empty string")
	}
	got := parseCSVOrigins(" https://a.test ,https://b.test ")
	if len(got) != 2 || got[0] != "https://a.test" || got[1] != "https://b.test" {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestConfig_CORSOrigins_Default(t *testing.T) {
	c := &Config{}
	def := c.CORSOrigins()
	if len(def) != 2 {
		t.Fatalf("expected 2 default origins, got %d", len(def))
	}
}

func TestConfig_CORSOrigins_Custom(t *testing.T) {
	c := &Config{CORSAllowedOrigins: []string{"https://app.example.com"}}
	if got := c.CORSOrigins(); len(got) != 1 || got[0] != "https://app.example.com" {
		t.Fatalf("unexpected: %#v", got)
	}
}
