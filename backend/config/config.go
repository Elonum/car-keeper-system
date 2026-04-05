package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database           DatabaseConfig
	Server             ServerConfig
	JWT                JWTConfig
	Storage            StorageConfig
	Env                string
	CORSAllowedOrigins []string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port string
	Host string
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

// StorageConfig controls on-disk document storage.
type StorageConfig struct {
	RootPath       string
	MaxUploadBytes int64
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "carkeeper"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "change-me-in-production"),
			ExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		},
		Storage: StorageConfig{
			RootPath:       getEnv("DOCUMENT_STORAGE_ROOT", "./data/documents"),
			MaxUploadBytes: getEnvAsInt64("DOCUMENT_MAX_UPLOAD_BYTES", 15<<20),
		},
		Env: getEnv("ENV", "development"),
		CORSAllowedOrigins: parseCSVOrigins(getEnv("CORS_ALLOWED_ORIGINS", "")),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.JWT.Secret == "change-me-in-production" && c.Env == "production" {
		return fmt.Errorf("JWT_SECRET must be changed in production")
	}
	if c.Storage.MaxUploadBytes < 1<<20 {
		c.Storage.MaxUploadBytes = 15 << 20
	}
	return nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c *JWTConfig) ExpiryDuration() time.Duration {
	return time.Duration(c.ExpiryHours) * time.Hour
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseCSVOrigins(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// CORSOrigins returns configured origins or localhost defaults for development.
func (c *Config) CORSOrigins() []string {
	if len(c.CORSAllowedOrigins) > 0 {
		return c.CORSAllowedOrigins
	}
	return []string{"http://localhost:3000", "http://localhost:5173"}
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

