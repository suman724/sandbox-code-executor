package main

import (
	"log"
	"net/http"

	"session-agent/internal/api"
	"session-agent/internal/api/handlers"
	"session-agent/internal/api/middleware"
	"session-agent/internal/config"
	"session-agent/internal/runtime"
	"session-agent/internal/telemetry"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	logger := telemetry.NewLogger("session-agent")
	runner := runtime.NewRunner()
	stepHandler := handlers.StepHandler{Runner: runner}

	var authMiddleware func(http.Handler) http.Handler
	if cfg.AuthBypass {
		authMiddleware = middleware.AuthBypassMiddleware
	} else {
		authMiddleware = middleware.AuthTokenMiddleware(cfg.AuthToken)
	}

	router := api.NewRouter(api.RouterDeps{
		StepsHandler:  stepHandler,
		AuthMiddleware: authMiddleware,
	})

	logger.Info("starting server")
	if err := http.ListenAndServe(cfg.ListenAddr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
