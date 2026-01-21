package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"session-agent/internal/api/middleware"
	"session-agent/internal/runtime"
	"shared/sessionagent"
)

type SessionHandler struct {
	Runner       *runtime.Runner
	RequireToken bool
}

func (h SessionHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.Runner == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	var req sessionagent.SessionRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.SessionID == "" || req.Runtime == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	headerToken := middleware.TokenFromRequest(r)
	if req.Token == "" {
		req.Token = headerToken
	} else if headerToken != "" && req.Token != headerToken {
		http.Error(w, "session token mismatch", http.StatusBadRequest)
		return
	}
	if h.RequireToken && req.Token == "" {
		http.Error(w, "missing session token", http.StatusUnauthorized)
		return
	}
	if _, err := h.Runner.RegisterSession(req); err != nil {
		if errors.Is(err, runtime.ErrSessionRuntimeMismatch) {
			http.Error(w, "session runtime mismatch", http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(sessionagent.SessionRegisterResponse{
		SessionID: req.SessionID,
		Status:    "registered",
	})
}

func (h SessionHandler) Terminate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.Runner == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if h.RequireToken {
		token := middleware.TokenFromRequest(r)
		if err := h.Runner.Authorize(sessionID, token); err != nil {
			http.Error(w, "invalid session token", http.StatusUnauthorized)
			return
		}
	}
	h.Runner.RemoveSession(sessionID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(sessionagent.SessionTerminateResponse{
		SessionID: sessionID,
		Status:    "terminated",
	})
}
