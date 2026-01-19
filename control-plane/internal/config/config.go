package config

import (
	"errors"
	"os"
)

type Config struct {
	Env         string
	DataPlaneURL string
	AuthzBypass bool
}

func Load() (Config, error) {
	cfg := Config{
		Env:          os.Getenv("ENV"),
		DataPlaneURL: os.Getenv("DATA_PLANE_URL"),
		AuthzBypass:  os.Getenv("AUTHZ_BYPASS") == "true",
	}
	return cfg, cfg.Validate()
}

func (c Config) Validate() error {
	if c.Env == "" {
		return errors.New("ENV is required")
	}
	if c.DataPlaneURL == "" {
		return errors.New("DATA_PLANE_URL is required")
	}
	if c.Env == "production" && c.AuthzBypass {
		return errors.New("AUTHZ_BYPASS is not allowed in production")
	}
	return nil
}
