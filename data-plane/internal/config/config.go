package config

import (
	"errors"
	"os"
)

type Config struct {
	Env              string
	RuntimeNamespace string
	RuntimeClass     string
	ArtifactRoot     string
	OtelEndpoint     string
	OtelService      string
	AuthIssuer       string
	AuthAudience     string
	AuthzBypass      bool
}

func Load() (Config, error) {
	cfg := Config{
		Env:              os.Getenv("ENV"),
		RuntimeNamespace: os.Getenv("RUNTIME_NAMESPACE"),
		RuntimeClass:     os.Getenv("RUNTIME_CLASS"),
		ArtifactRoot:     os.Getenv("ARTIFACT_ROOT"),
		OtelEndpoint:     os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OtelService:      os.Getenv("OTEL_SERVICE_NAME"),
		AuthIssuer:       os.Getenv("AUTH_ISSUER"),
		AuthAudience:     os.Getenv("AUTH_AUDIENCE"),
		AuthzBypass:      os.Getenv("AUTHZ_BYPASS") == "true",
	}
	return cfg, cfg.Validate()
}

func (c Config) Validate() error {
	if c.Env == "" {
		return errors.New("ENV is required")
	}
	if c.Env == "production" && c.AuthzBypass {
		return errors.New("AUTHZ_BYPASS is not allowed in production")
	}
	return nil
}
