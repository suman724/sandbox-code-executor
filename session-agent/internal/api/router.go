package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RouterDeps struct {
	HealthHandler           http.Handler
	StepsHandler            http.Handler
	AuthMiddleware          func(http.Handler) http.Handler
	SessionsHandler         http.Handler
	SessionTerminateHandler http.Handler
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
	sessionsHandler := deps.SessionsHandler
	if sessionsHandler == nil {
		sessionsHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "sessions handler not configured", http.StatusNotImplemented)
		})
	}
	terminateHandler := deps.SessionTerminateHandler
	if terminateHandler == nil {
		terminateHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "terminate handler not configured", http.StatusNotImplemented)
		})
	}

	router := chi.NewRouter()
	router.Get("/v1/health", healthHandler.ServeHTTP)
	router.Route("/v1", func(r chi.Router) {
		if deps.AuthMiddleware != nil {
			r.With(deps.AuthMiddleware).Post("/steps", stepsHandler.ServeHTTP)
			r.With(deps.AuthMiddleware).Post("/sessions", sessionsHandler.ServeHTTP)
			r.With(deps.AuthMiddleware).Post("/sessions/{sessionId}/terminate", terminateHandler.ServeHTTP)
			return
		}
		r.Post("/steps", stepsHandler.ServeHTTP)
		r.Post("/sessions", sessionsHandler.ServeHTTP)
		r.Post("/sessions/{sessionId}/terminate", terminateHandler.ServeHTTP)
	})

	return router
}
