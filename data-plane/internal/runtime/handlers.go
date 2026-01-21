package runtime

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"shared/sessionagent"
)

type RunHandler struct {
	Runner Runner
	Store  RunStore
}

type runRequest struct {
	JobID        string `json:"jobId"`
	PolicyID     string `json:"policyId"`
	Language     string `json:"language"`
	Code         string `json:"code"`
	WorkspaceRef string `json:"workspaceRef"`
}

type runResponse struct {
	RunID string `json:"run_id"`
}

type Run struct {
	ID           string
	JobID        string
	Status       string
	ExitStatus   int
	OutputRef    string
	ErrorRef     string
	ArtifactRefs []string
}

type RunStore interface {
	Create(ctx context.Context, run Run) error
	Get(ctx context.Context, id string) (Run, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

type SessionHandler struct {
	Runtime  SessionRuntime
	Registry SessionRegistry
	Agent    *AgentClient
}

type sessionRequest struct {
	SessionID    string `json:"sessionId"`
	PolicyID     string `json:"policyId"`
	WorkspaceRef string `json:"workspaceRef"`
	Runtime      string `json:"runtime"`
}

type sessionResponse struct {
	ID        string `json:"id"`
	RuntimeID string `json:"runtimeId"`
	Status    string `json:"status"`
}

type sessionStepRequest struct {
	Command string `json:"command"`
	Runtime string `json:"runtime"`
}

type sessionStepResponse struct {
	Status string `json:"status"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func (h SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/terminate") {
		h.handleTerminate(w, r)
		return
	}
	if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/steps") {
		h.handleStep(w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.handleCreate(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h SessionHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	if h.Runtime == nil || h.Registry == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	var req sessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("sessions: decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.SessionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	route, err := h.Runtime.StartSession(r.Context(), req.SessionID, req.PolicyID, req.WorkspaceRef, req.Runtime)
	if err != nil {
		log.Printf("sessions: start error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := h.Registry.Put(req.SessionID, route); err != nil {
		log.Printf("sessions: registry error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(sessionResponse{ID: req.SessionID, RuntimeID: route.RuntimeID, Status: "running"})
}

func (h SessionHandler) handleStep(w http.ResponseWriter, r *http.Request) {
	if h.Runtime == nil || h.Registry == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var req sessionStepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Command == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Runtime == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	stepID := "step-" + time.Now().UTC().Format("20060102150405.000")
	route, ok := h.Registry.Get(sessionID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	output, err := h.Runtime.RunStep(r.Context(), route.RuntimeID, req.Command)
	if err == nil {
		log.Printf("sessions: route session_id=%s runtime_id=%s endpoint=%s auth_mode=%s step_id=%s status=completed", sessionID, route.RuntimeID, route.Endpoint, route.AuthMode, stepID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(sessionStepResponse{
			Status: "accepted",
			Stdout: output.Stdout,
			Stderr: output.Stderr,
		})
		return
	}
	if h.Agent == nil || route.Endpoint == "" {
		log.Printf("sessions: route session_id=%s runtime_id=%s endpoint=%s auth_mode=%s step_id=%s status=failed", sessionID, route.RuntimeID, route.Endpoint, route.AuthMode, stepID)
		log.Printf("sessions: step error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	agentResult, agentErr := h.Agent.RunStep(r.Context(), AgentRoute{
		Endpoint: route.Endpoint,
		Token:    route.Token,
		AuthMode: route.AuthMode,
	}, sessionagent.StepRequest{
		SessionID: sessionID,
		StepID:    stepID,
		Code:      req.Command,
		Runtime:   req.Runtime,
	})
	if agentErr != nil {
		log.Printf("sessions: route session_id=%s runtime_id=%s endpoint=%s auth_mode=%s step_id=%s status=failed", sessionID, route.RuntimeID, route.Endpoint, route.AuthMode, stepID)
		log.Printf("sessions: agent step error: %v", agentErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("sessions: route session_id=%s runtime_id=%s endpoint=%s auth_mode=%s step_id=%s status=%s", sessionID, route.RuntimeID, route.Endpoint, route.AuthMode, stepID, agentResult.Status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(sessionStepResponse{
		Status: agentResult.Status,
		Stdout: agentResult.Stdout,
		Stderr: agentResult.Stderr,
	})
}

func (h SessionHandler) handleTerminate(w http.ResponseWriter, r *http.Request) {
	if h.Runtime == nil || h.Registry == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	route, ok := h.Registry.Get(sessionID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := h.Runtime.TerminateSession(r.Context(), route.RuntimeID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Registry.Delete(sessionID)
	w.WriteHeader(http.StatusAccepted)
}

func (h RunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/terminate") {
		h.handleTerminate(w, r)
		return
	}
	switch r.Method {
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodGet:
		h.handleGet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h RunHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req runRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("runs: decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	runID, err := h.Runner.Run(r.Context(), req.JobID, req.Language, req.Code)
	if err != nil {
		log.Printf("runs: run error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if h.Store != nil {
		_ = h.Store.Create(r.Context(), Run{ID: runID, JobID: req.JobID, Status: "running"})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(runResponse{RunID: runID})
	log.Printf("runs: accepted job_id=%s run_id=%s ts=%s", req.JobID, runID, time.Now().UTC().Format(time.RFC3339))
}

func (h RunHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	if h.Store == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	runID := chi.URLParam(r, "runId")
	if runID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	run, err := h.Store.Get(r.Context(), runID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(run)
}

func (h RunHandler) handleTerminate(w http.ResponseWriter, r *http.Request) {
	if h.Store == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	runID := chi.URLParam(r, "runId")
	if runID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Store.UpdateStatus(r.Context(), runID, "terminated"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
