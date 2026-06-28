// Package config loads and validates application configuration.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// ─── Constants ────────────────────────────────────────────────────

const (
	defaultPort               = "8080"
	defaultJWTExpirationHours = 24
	minJWTSecretLength        = 32
)

// ─── Config ─────────────────────────────────────────────────────────

// Config holds all runtime configuration.
type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	JWTExpiration   time.Duration
	APIURL          string
	Env             string
	RedisURL        string
	OwnerEmail      string
	OwnerPassword   string
	PartnerEmail    string
	PartnerPassword string
}

// Load reads configuration from environment variables and validates it.
func Load() (*Config, error) {
	cfg := &Config{
		Port:            getEnv("PORT", defaultPort),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		APIURL:          getEnv("API_URL", "http://localhost:"+getEnv("PORT", defaultPort)),
		Env:             getEnv("ENV", "development"),
		RedisURL:        os.Getenv("REDIS_URL"),
		OwnerEmail:      getEnv("OWNER_EMAIL", "owner@home.local"),
		OwnerPassword:   os.Getenv("OWNER_PASSWORD"),
		PartnerEmail:    getEnv("PARTNER_EMAIL", "partner@home.local"),
		PartnerPassword: os.Getenv("PARTNER_PASSWORD"),
	}

	expirationHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", strconv.Itoa(defaultJWTExpirationHours)))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION_HOURS: %w", err)
	}
	cfg.JWTExpiration = time.Duration(expirationHours) * time.Hour

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that required fields are present and secure.
func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if len(c.JWTSecret) < minJWTSecretLength {
		return fmt.Errorf("JWT_SECRET must be at least %d characters", minJWTSecretLength)
	}

	if c.OwnerPassword == "" || c.PartnerPassword == "" {
		return fmt.Errorf("OWNER_PASSWORD and PARTNER_PASSWORD are required")
	}

	return nil
}

// IsDevelopment returns true if the environment is development.
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// getEnv returns the value of the environment variable or the default.
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
