package tools

import (
	"net/http"

	"control-plane/internal/api/handlers"
)

type WorkflowsTool struct {
	Handler handlers.WorkflowHandler
}

func (t WorkflowsTool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Handler.ServeHTTP(w, r)
}
