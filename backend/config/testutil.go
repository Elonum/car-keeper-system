package config

import "os"

// TestConfig returns settings suitable for unit and integration tests.
func TestConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     envOr("DB_HOST", "localhost"),
			Port:     5432,
			User:     envOr("DB_USER", "postgres"),
			Password: envOr("DB_PASSWORD", "postgres"),
			Name:     envOr("DB_NAME", "carkeeper"),
			SSLMode:  envOr("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Host:             "127.0.0.1",
			Port:             "0",
			MaxJSONBodyBytes: 1 << 20,
		},
		JWT: JWTConfig{
			Secret:       "carkeeper-test-jwt-secret",
			ExpiryHours:  24,
			SecureCookie: false,
		},
		Storage: StorageConfig{
			RootPath:       envOr("DOCUMENT_STORAGE_ROOT", "./testdata/documents"),
			MaxUploadBytes: 15 << 20,
		},
		Env: "test",
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
