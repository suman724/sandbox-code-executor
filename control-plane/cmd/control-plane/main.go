package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"control-plane/internal/api"
	"control-plane/internal/audit"
	"control-plane/internal/config"
	"control-plane/internal/orchestration"
	"control-plane/internal/policy"
	"control-plane/internal/sessions"
	"control-plane/internal/storage"
	"control-plane/pkg/client"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	addr := ":" + getenv("PORT", "8080")
	log.Printf("control-plane starting env=%s addr=%s", cfg.Env, addr)

	stores, err := storage.NewStoreSet(context.Background(), cfg.DatabaseDriver, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("storage error: %v", err)
	}
	if stores.Close != nil {
		defer func() {
			if err := stores.Close(); err != nil {
				log.Printf("storage close error: %v", err)
			}
		}()
	}

	allowAll := policy.StaticRulesetResolver{RulesetText: "package policy\n default allow = true"}
	evaluator := &policy.OPAEvaluator{Resolver: allowAll}
	enforcer := orchestration.PolicyEnforcer{Evaluator: evaluator}
	dataPlaneClient := client.DataPlaneClient{BaseURL: cfg.DataPlaneURL}

	jobService := orchestration.JobService{
		Store:    stores.JobStore,
		Client:   dataPlaneClient,
		Enforcer: enforcer,
	}
	sessionService := sessions.Service{
		Store:    stores.SessionStore,
		Client:   dataPlaneClient,
		Enforcer: enforcer,
		Logger:   audit.StdoutLogger{},
	}

	deps := api.Dependencies{
		JobService:     &jobService,
		JobStore:       stores.JobStore,
		SessionService: &sessionService,
		PolicyStore:    policy.NewInMemoryStore(),
		AuditStore:     &audit.InMemoryStore{},
	}

	if err := http.ListenAndServe(addr, api.RouterWithDependencies(deps)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
