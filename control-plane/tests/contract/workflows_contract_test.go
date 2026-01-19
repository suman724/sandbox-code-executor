package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"control-plane/internal/api/handlers"
	"control-plane/internal/orchestration"
)

type mockWorkflowStore struct {
	created []orchestration.Workflow
}

func (m *mockWorkflowStore) Create(ctx context.Context, workflow orchestration.Workflow) error {
	_ = ctx
	m.created = append(m.created, workflow)
	return nil
}

func (m *mockWorkflowStore) UpdateStatus(ctx context.Context, id string, status orchestration.WorkflowStatus) error {
	_ = ctx
	_ = id
	_ = status
	return nil
}

type mockWorkflowRunner struct{}

func (mockWorkflowRunner) RunStep(ctx context.Context, workflowID string, step orchestration.WorkflowStep, memory orchestration.SharedMemoryStore) (string, error) {
	_ = ctx
	_ = workflowID
	_ = step
	_ = memory
	return "job-1", nil
}

func TestWorkflowsContractCreate(t *testing.T) {
	store := &mockWorkflowStore{}
	service := orchestration.WorkflowService{
		Store:  store,
		Runner: mockWorkflowRunner{},
		Memory: orchestration.NewMemoryStore(),
	}
	handler := handlers.WorkflowHandler{Service: service}

	payload := map[string]any{
		"tenantId": "tenant-1",
		"steps":    []string{"agent-1", "agent-2"},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/workflows", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, rec.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["id"] == "" || resp["status"] == "" {
		t.Fatalf("expected id and status")
	}
	if len(store.created) != 1 {
		t.Fatalf("expected workflow created")
	}
}
