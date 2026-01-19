package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"control-plane/internal/sessions"
)

type SessionHandler struct {
	Service sessions.Service
}

type sessionRequest struct {
	ID       string        `json:"id"`
	TenantID string        `json:"tenant_id"`
	PolicyID string        `json:"policy_id"`
	TTL      time.Duration `json:"ttl"`
}

type sessionResponse struct {
	SessionID string `json:"session_id"`
	RunID     string `json:"run_id"`
}

func (h SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req sessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("sessions: decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	session := sessions.Session{
		ID:       req.ID,
		TenantID: req.TenantID,
		PolicyID: req.PolicyID,
		TTL:      req.TTL,
		Status:   sessions.StatusActive,
	}
	runID, err := h.Service.CreateSession(r.Context(), session)
	if err != nil {
		log.Printf("sessions: create error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sessionResponse{SessionID: session.ID, RunID: runID})
}
