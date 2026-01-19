package mcp

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"control-plane/internal/api/handlers"
	"control-plane/internal/api/middleware"
	"control-plane/internal/mcp/tools"
	"control-plane/internal/storage"
)

type Dependencies struct {
	JobsHandler      handlers.JobHandler
	SessionsHandler  handlers.SessionHandler
	WorkflowsHandler handlers.WorkflowHandler
	ArtifactStore    storage.ArtifactStore
}

func Router() http.Handler {
	return RouterWithDependencies(Dependencies{})
}

func RouterWithDependencies(deps Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Auth)

	jobsTool := tools.JobsTool{Handler: deps.JobsHandler}
	sessionsTool := tools.SessionsTool{Handler: deps.SessionsHandler}
	workflowsTool := tools.WorkflowsTool{Handler: deps.WorkflowsHandler}
	artifactsTool := tools.ArtifactsTool{Store: deps.ArtifactStore}

	r.Post("/tools/jobs", jobsTool.ServeHTTP)
	r.Get("/tools/jobs/{jobId}", jobsTool.ServeHTTP)

	r.Post("/tools/sessions", sessionsTool.ServeHTTP)
	r.Post("/tools/sessions/{sessionId}/steps", sessionsTool.ServeHTTP)

	r.Post("/tools/workflows", workflowsTool.ServeHTTP)

	r.Post("/tools/artifacts/upload", artifactsTool.Upload)
	r.Get("/tools/artifacts/{artifactId}/download", artifactsTool.Download)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}
