package config

import (
	"errors"
	"os"
)

type Config struct {
	Env                 string
	RuntimeNamespace    string
	RuntimeClass        string
	ArtifactRoot        string
	SessionRuntime      string
	SessionRegistry     string
	SessionRegistryPath string
	SessionImage        string
	SessionImagePython  string
	SessionImageNode    string
	AgentEndpoint       string
	AgentAuthMode       string
	AgentAuthToken      string
	OtelEndpoint        string
	OtelService         string
	AuthIssuer          string
	AuthAudience        string
	AuthzBypass         bool
}

func Load() (Config, error) {
	cfg := Config{
		Env:                 os.Getenv("ENV"),
		RuntimeNamespace:    os.Getenv("RUNTIME_NAMESPACE"),
		RuntimeClass:        os.Getenv("RUNTIME_CLASS"),
		ArtifactRoot:        os.Getenv("ARTIFACT_ROOT"),
		SessionRuntime:      getenv("SESSION_RUNTIME_BACKEND", "local"),
		SessionRegistry:     getenv("SESSION_REGISTRY_BACKEND", "memory"),
		SessionRegistryPath: os.Getenv("SESSION_REGISTRY_PATH"),
		SessionImage:        os.Getenv("SESSION_RUNTIME_IMAGE"),
		SessionImagePython:  os.Getenv("SESSION_RUNTIME_IMAGE_PYTHON"),
		SessionImageNode:    os.Getenv("SESSION_RUNTIME_IMAGE_NODE"),
		AgentEndpoint:       os.Getenv("SESSION_AGENT_ENDPOINT"),
		AgentAuthMode:       getenv("SESSION_AGENT_AUTH_MODE", "enforced"),
		AgentAuthToken:      os.Getenv("SESSION_AGENT_AUTH_TOKEN"),
		OtelEndpoint:        os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OtelService:         os.Getenv("OTEL_SERVICE_NAME"),
		AuthIssuer:          os.Getenv("AUTH_ISSUER"),
		AuthAudience:        os.Getenv("AUTH_AUDIENCE"),
		AuthzBypass:         os.Getenv("AUTHZ_BYPASS") == "true",
	}
	return cfg, cfg.Validate()
}

func (c Config) Validate() error {
	if c.Env == "" {
		return errors.New("ENV is required")
	}
	if c.SessionRuntime != "local" && c.SessionRuntime != "k8s" {
		return errors.New("SESSION_RUNTIME_BACKEND must be local or k8s")
	}
	if c.SessionRegistry != "memory" && c.SessionRegistry != "file" {
		return errors.New("SESSION_REGISTRY_BACKEND must be memory or file")
	}
	if c.SessionRegistry == "file" && c.SessionRegistryPath == "" {
		return errors.New("SESSION_REGISTRY_PATH is required for file registry backend")
	}
	if c.Env == "production" && c.SessionRegistry == "memory" {
		return errors.New("SESSION_REGISTRY_BACKEND must be persistent in production")
	}
	if c.Env == "production" && c.AuthzBypass {
		return errors.New("AUTHZ_BYPASS is not allowed in production")
	}
	if c.AgentAuthMode != "enforced" && c.AgentAuthMode != "bypass" {
		return errors.New("SESSION_AGENT_AUTH_MODE must be enforced or bypass")
	}
	if c.Env == "production" && c.AgentAuthMode == "bypass" {
		return errors.New("SESSION_AGENT_AUTH_MODE=bypass is not allowed in production")
	}
	if c.AgentAuthMode == "enforced" && c.AgentAuthToken == "" {
		return errors.New("SESSION_AGENT_AUTH_TOKEN is required when auth mode is enforced")
	}
	return nil
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
