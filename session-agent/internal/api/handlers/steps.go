package handlers

import (
	"encoding/json"
	"net/http"

	"session-agent/internal/api/middleware"
	"session-agent/internal/runtime"
	"shared/sessionagent"
)

type StepHandler struct {
	Runner       *runtime.Runner
	RequireToken bool
}

func (h StepHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.Runner == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	var req sessionagent.StepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if h.RequireToken {
		token := middleware.TokenFromRequest(r)
		if err := h.Runner.Authorize(req.SessionID, token); err != nil {
			http.Error(w, "invalid session token", http.StatusUnauthorized)
			return
		}
	}
	result, err := h.Runner.RunStep(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}
