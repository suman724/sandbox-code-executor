package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"control-plane/internal/audit"
)

type AuditHandler struct {
	Store audit.Store
}

func (h AuditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenantId")
	if tenantID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var since, until time.Time
	if raw := r.URL.Query().Get("since"); raw != "" {
		parsed, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		since = parsed
	}
	if raw := r.URL.Query().Get("until"); raw != "" {
		parsed, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		until = parsed
	}
	events, err := h.Store.List(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	filtered := make([]audit.Event, 0, len(events))
	for _, event := range events {
		if event.TenantID != tenantID {
			continue
		}
		if !since.IsZero() && event.Time.Before(since) {
			continue
		}
		if !until.IsZero() && event.Time.After(until) {
			continue
		}
		filtered = append(filtered, event)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(filtered)
}
