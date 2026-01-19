package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"control-plane/internal/orchestration"
	"control-plane/internal/storage"
)

type JobHandler struct {
	Service orchestration.JobService
	Store   storage.JobStore
}

type jobRequest struct {
	TenantID string `json:"tenantId"`
	AgentID  string `json:"agentId"`
	PolicyID string `json:"policyId"`
	Language string `json:"language"`
	Code     string `json:"code"`
}

type jobResponse struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	ExitStatus int    `json:"exit_status,omitempty"`
}

func (h JobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.handleGet(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req jobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("jobs: decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	job := orchestration.Job{
		ID:        "job-" + time.Now().UTC().Format("20060102150405"),
		TenantID:  req.TenantID,
		AgentID:   req.AgentID,
		PolicyID:  req.PolicyID,
		Language:  req.Language,
		Code:      req.Code,
		Status:    orchestration.JobQueued,
	}
	job.Workspace = job.ID
	_, err := h.Service.CreateJob(r.Context(), job)
	if err != nil {
		log.Printf("jobs: create error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(jobResponse{ID: job.ID, Status: string(job.Status)})
	log.Printf("jobs: accepted job_id=%s ts=%s", job.ID, time.Now().UTC().Format(time.RFC3339))
}

func (h JobHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	if h.Store == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jobID := chi.URLParam(r, "jobId")
	if jobID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	job, err := h.Store.Get(r.Context(), jobID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(jobResponse{ID: job.ID, Status: job.Status})
}
