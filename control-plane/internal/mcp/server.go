package mcp

import (
	"errors"
	"net/http"
)

type Server struct {
	Addr   string
	Router http.Handler
}

func NewServer(addr string, router http.Handler) Server {
	return Server{Addr: addr, Router: router}
}

func (s Server) ListenAndServe() error {
	if s.Router == nil {
		return errors.New("missing mcp router")
	}
	if s.Addr == "" {
		return errors.New("missing mcp addr")
	}
	return http.ListenAndServe(s.Addr, s.Router)
}
