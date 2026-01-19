package runtime

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type RunHandler struct {
	Runner Runner
}

type runRequest struct {
	JobID    string `json:"job_id"`
	Language string `json:"language"`
	Code     string `json:"code"`
}

type runResponse struct {
	RunID string `json:"run_id"`
}

func (h RunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(runResponse{RunID: runID})
	log.Printf("runs: accepted job_id=%s run_id=%s ts=%s", req.JobID, runID, time.Now().UTC().Format(time.RFC3339))
}
