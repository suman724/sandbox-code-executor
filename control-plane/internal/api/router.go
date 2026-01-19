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
	"control-plane/internal/storage"
)

type Dependencies struct {
	JobService      *orchestration.JobService
	JobStore        storage.JobStore
	SessionService  *sessions.Service
	Stepper         *sessions.StepService
	PolicyStore     policy.Store
	AuditStore      audit.Store
	WorkflowService *orchestration.WorkflowService
	ServiceStarter  func(service services.Service) (string, error)
}

func Router() http.Handler {
	return RouterWithDependencies(Dependencies{})
}

func RouterWithDependencies(deps Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Auth)

	notImplemented := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})

	jobService := orchestration.JobService{}
	if deps.JobService != nil {
		jobService = *deps.JobService
	}
	jobStore := deps.JobStore
	if jobStore == nil && deps.JobService != nil {
		jobStore = deps.JobService.Store
	}
	r.Post("/jobs", handlers.JobHandler{Service: jobService, Store: jobStore}.ServeHTTP)
	r.Get("/jobs/{jobId}", notImplemented)

	sessionService := sessions.Service{}
	if deps.SessionService != nil {
		sessionService = *deps.SessionService
	}
	stepper := sessions.StepService{}
	if deps.Stepper != nil {
		stepper = *deps.Stepper
	}
	r.Post("/sessions", handlers.SessionHandler{Service: sessionService, Stepper: stepper}.ServeHTTP)
	r.Post("/sessions/{sessionId}/steps", notImplemented)

	r.Post("/artifacts/upload", notImplemented)
	r.Get("/artifacts/{artifactId}/download", notImplemented)

	policyStore := deps.PolicyStore
	if policyStore == nil {
		policyStore = policy.NewInMemoryStore()
	}
	auditStore := deps.AuditStore
	if auditStore == nil {
		auditStore = &audit.InMemoryStore{}
	}
	r.Post("/policies", handlers.PolicyHandler{Store: policyStore}.ServeHTTP)
	r.Get("/audit/events", handlers.AuditHandler{Store: auditStore}.ServeHTTP)

	workflowService := orchestration.WorkflowService{}
	if deps.WorkflowService != nil {
		workflowService = *deps.WorkflowService
	}
	r.Post("/workflows", handlers.WorkflowHandler{Service: workflowService}.ServeHTTP)
	r.Post("/services", handlers.ServiceHandler{Starter: deps.ServiceStarter}.ServeHTTP)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}
