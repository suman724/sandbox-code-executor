package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"control-plane/internal/orchestration"
)

type WorkflowHandler struct {
	Service orchestration.WorkflowService
}

type workflowRequest struct {
	TenantID string   `json:"tenantId"`
	Steps    []string `json:"steps"`
}

type workflowResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (h WorkflowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req workflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.TenantID == "" || len(req.Steps) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	workflowID := "workflow-" + time.Now().UTC().Format("20060102150405")
	steps := make([]orchestration.WorkflowStep, 0, len(req.Steps))
	for i, agentID := range req.Steps {
		steps = append(steps, orchestration.WorkflowStep{
			ID:       workflowID + "-step-" + strconv.Itoa(i+1),
			Sequence: i + 1,
			AgentID:  agentID,
			Status:   orchestration.WorkflowStepQueued,
		})
	}
	wf := orchestration.Workflow{
		ID:        workflowID,
		TenantID:  req.TenantID,
		Status:    orchestration.WorkflowQueued,
		CreatedAt: time.Now().UTC(),
		Steps:     steps,
	}
	if err := h.Service.Start(r.Context(), wf); err != nil {
		log.Printf("workflows: create error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(workflowResponse{ID: wf.ID, Status: string(wf.Status)})
}
