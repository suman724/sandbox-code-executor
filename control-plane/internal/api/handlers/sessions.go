package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"control-plane/internal/sessions"
)

type SessionHandler struct {
	Service sessions.Service
	Stepper sessions.StepService
}

type sessionRequest struct {
	TenantID   string `json:"tenantId"`
	AgentID    string `json:"agentId"`
	PolicyID   string `json:"policyId"`
	TTLSeconds int    `json:"ttlSeconds"`
}

type sessionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type stepRequest struct {
	Command string `json:"command"`
}

type stepResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func (h SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && chi.URLParam(r, "sessionId") != "" {
		h.handleStep(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req sessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("sessions: decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	session := sessions.Session{
		ID:       "session-" + time.Now().UTC().Format("20060102150405"),
		TenantID: req.TenantID,
		AgentID:  req.AgentID,
		PolicyID: req.PolicyID,
		TTL:      time.Duration(req.TTLSeconds) * time.Second,
		Status:   sessions.StatusActive,
	}
	_, err := h.Service.CreateSession(r.Context(), session)
	if err != nil {
		log.Printf("sessions: create error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sessionResponse{ID: session.ID, Status: string(session.Status)})
}

func (h SessionHandler) handleStep(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req stepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Command == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if h.Stepper.Runner == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	result, err := h.Stepper.Run(r.Context(), chi.URLParam(r, "sessionId"), req.Command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(stepResponse{
		ID:     result.ID,
		Status: "accepted",
		Stdout: result.Stdout,
		Stderr: result.Stderr,
	})
}
