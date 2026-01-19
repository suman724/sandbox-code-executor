package config

import (
	"errors"
	"os"
)

type Config struct {
	Env            string
	DataPlaneURL   string
	DatabaseDriver string
	DatabaseURL    string
	ArtifactBucket string
	OtelEndpoint   string
	OtelService    string
	AuthIssuer     string
	AuthAudience   string
	MCPAddr        string
	AuthzBypass    bool
}

func Load() (Config, error) {
	cfg := Config{
		Env:            os.Getenv("ENV"),
		DataPlaneURL:   os.Getenv("DATA_PLANE_URL"),
		DatabaseDriver: getenv("DATABASE_DRIVER", "postgres"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		ArtifactBucket: os.Getenv("ARTIFACT_BUCKET"),
		OtelEndpoint:   os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OtelService:    os.Getenv("OTEL_SERVICE_NAME"),
		AuthIssuer:     os.Getenv("AUTH_ISSUER"),
		AuthAudience:   os.Getenv("AUTH_AUDIENCE"),
		MCPAddr:        os.Getenv("MCP_ADDR"),
		AuthzBypass:    os.Getenv("AUTHZ_BYPASS") == "true",
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
	if c.DatabaseDriver != "postgres" && c.DatabaseDriver != "sqlite" {
		return errors.New("DATABASE_DRIVER must be postgres or sqlite")
	}
	if c.DatabaseDriver == "sqlite" && c.Env == "production" {
		return errors.New("DATABASE_DRIVER=sqlite is not allowed in production")
	}
	if c.DatabaseDriver == "postgres" && c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required for postgres")
	}
	if c.DatabaseDriver == "sqlite" && c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required for sqlite")
	}
	if c.Env == "production" && c.AuthzBypass {
		return errors.New("AUTHZ_BYPASS is not allowed in production")
	}
	return nil
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
