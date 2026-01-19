package runtime

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
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
