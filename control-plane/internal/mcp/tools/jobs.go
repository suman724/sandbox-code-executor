package tools

import (
	"net/http"

	"control-plane/internal/api/handlers"
)

type JobsTool struct {
	Handler handlers.JobHandler
}

func (t JobsTool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Handler.ServeHTTP(w, r)
}
