package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"control-plane/internal/api/handlers"
	"control-plane/internal/api/middleware"
	"control-plane/internal/audit"
	"control-plane/internal/orchestration"
	"control-plane/internal/policy"
	"control-plane/internal/services"
	"control-plane/internal/sessions"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Auth)

	notImplemented := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})

	r.Post("/jobs", handlers.JobHandler{Service: orchestration.JobService{}}.ServeHTTP)
	r.Get("/jobs/{jobId}", notImplemented)

	r.Post("/sessions", handlers.SessionHandler{Service: sessions.Service{}}.ServeHTTP)
	r.Post("/sessions/{sessionId}/steps", notImplemented)

	r.Post("/artifacts/upload", notImplemented)
	r.Get("/artifacts/{artifactId}/download", notImplemented)

	r.Post("/policies", handlers.PolicyHandler{Store: policy.Store(nil)}.ServeHTTP)
	r.Get("/audit/events", handlers.AuditHandler{Store: &audit.InMemoryStore{}}.ServeHTTP)

	r.Post("/workflows", handlers.WorkflowHandler{Service: orchestration.WorkflowService{}}.ServeHTTP)
	r.Post("/services", handlers.ServiceHandler{Starter: nil}.ServeHTTP)

	return r
}
