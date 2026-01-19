package tools

import (
	"net/http"

	"control-plane/internal/api/handlers"
)

type SessionsTool struct {
	Handler handlers.SessionHandler
}

func (t SessionsTool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Handler.ServeHTTP(w, r)
}
