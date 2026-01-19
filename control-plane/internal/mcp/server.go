package mcp

import "net/http"

type Server struct {
	Router http.Handler
}

func NewServer(router http.Handler) Server {
	return Server{Router: router}
}
