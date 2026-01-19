package handlers

import (
	"encoding/json"
	"net/http"

	"control-plane/internal/audit"
)

type AuditHandler struct {
	Store audit.Store
}

func (h AuditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	events, err := h.Store.List(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(events)
}
