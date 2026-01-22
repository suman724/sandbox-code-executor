package config

import (
	"errors"
	"os"
)

type Config struct {
	Env        string
	ListenAddr string
	AuthBypass bool
}

func Load() (Config, error) {
	cfg := Config{
		Env:        os.Getenv("ENV"),
		ListenAddr: getenv("SESSION_AGENT_ADDR", ":9000"),
		AuthBypass: os.Getenv("SESSION_AGENT_AUTH_BYPASS") == "true",
	}
	return cfg, cfg.Validate()
}

func (c Config) Validate() error {
	if c.Env == "" {
		return errors.New("ENV is required")
	}
	if c.ListenAddr == "" {
		return errors.New("SESSION_AGENT_ADDR is required")
	}
	if c.Env == "production" && c.AuthBypass {
		return errors.New("SESSION_AGENT_AUTH_BYPASS is not allowed in production")
	}
	return nil
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
