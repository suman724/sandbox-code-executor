package runtime

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(Auth)

	notImplemented := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})

	r.Post("/runs", RunHandler{}.ServeHTTP)
	r.Get("/runs/{runId}", notImplemented)
	r.Post("/runs/{runId}/terminate", notImplemented)

	return r
}
