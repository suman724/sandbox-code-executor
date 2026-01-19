package mcp

import "net/http"

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{mux: http.NewServeMux()}
}

func (r *Router) Handle(path string, handler http.Handler) {
	r.mux.Handle(path, handler)
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
