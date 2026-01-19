package runtime

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(Auth)

	handler := RunHandler{}
	r.Post("/runs", handler.ServeHTTP)
	r.Get("/runs/{runId}", handler.ServeHTTP)
	r.Post("/runs/{runId}/terminate", handler.ServeHTTP)

	return r
}
