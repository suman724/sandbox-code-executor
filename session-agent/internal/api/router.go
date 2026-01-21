package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RouterDeps struct {
	HealthHandler http.Handler
	StepsHandler  http.Handler
	AuthMiddleware func(http.Handler) http.Handler
}

func NewRouter(deps RouterDeps) http.Handler {
	healthHandler := deps.HealthHandler
	if healthHandler == nil {
		healthHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	}

	stepsHandler := deps.StepsHandler
	if stepsHandler == nil {
		stepsHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "steps handler not configured", http.StatusNotImplemented)
		})
	}

	router := chi.NewRouter()
	router.Get("/v1/health", healthHandler.ServeHTTP)
	router.Route("/v1", func(r chi.Router) {
		if deps.AuthMiddleware != nil {
			r.With(deps.AuthMiddleware).Post("/steps", stepsHandler.ServeHTTP)
			return
		}
		r.Post("/steps", stepsHandler.ServeHTTP)
	})

	return router
}
