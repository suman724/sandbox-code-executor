package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"control-plane/internal/orchestration"
)

type JobHandler struct {
	Service orchestration.JobService
}

type jobRequest struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	PolicyID  string `json:"policy_id"`
	Language  string `json:"language"`
	Code      string `json:"code"`
	Workspace string `json:"workspace"`
}

type jobResponse struct {
	JobID string `json:"job_id"`
	RunID string `json:"run_id"`
}

func (h JobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req jobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("jobs: decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	job := orchestration.Job{
		ID:        req.ID,
		TenantID:  req.TenantID,
		PolicyID:  req.PolicyID,
		Language:  req.Language,
		Code:      req.Code,
		Workspace: req.Workspace,
		Status:    orchestration.JobQueued,
	}
	runID, err := h.Service.CreateJob(r.Context(), job)
	if err != nil {
		log.Printf("jobs: create error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(jobResponse{JobID: job.ID, RunID: runID})
	log.Printf("jobs: accepted job_id=%s run_id=%s ts=%s", job.ID, runID, time.Now().UTC().Format(time.RFC3339))
}
