package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"control-plane/internal/orchestration"
)

type WorkflowHandler struct {
	Service orchestration.WorkflowService
}

type workflowRequest struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
}

func (h WorkflowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req workflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	wf := orchestration.Workflow{ID: req.ID, TenantID: req.TenantID, Status: orchestration.WorkflowQueued, CreatedAt: time.Now()}
	if err := h.Service.Start(r.Context(), wf); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
