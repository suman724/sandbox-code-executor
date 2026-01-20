package runtime

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Dependencies struct {
	RunHandler     RunHandler
	SessionHandler SessionHandler
}

func Router() http.Handler {
	return RouterWithDependencies(Dependencies{})
}

func RouterWithDependencies(deps Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Use(Auth)

	runHandler := deps.RunHandler
	r.Post("/runs", runHandler.ServeHTTP)
	r.Get("/runs/{runId}", runHandler.ServeHTTP)
	r.Post("/runs/{runId}/terminate", runHandler.ServeHTTP)

	sessionHandler := deps.SessionHandler
	r.Post("/sessions", sessionHandler.ServeHTTP)
	r.Post("/sessions/{sessionId}/steps", sessionHandler.ServeHTTP)
	r.Post("/sessions/{sessionId}/terminate", sessionHandler.ServeHTTP)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/openapi.yaml", OpenAPIHandler())
	r.Get("/docs", SwaggerUIHandler())

	return r
}
