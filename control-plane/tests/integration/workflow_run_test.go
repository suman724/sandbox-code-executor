package integration

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
	created  []orchestration.Workflow
	statuses map[string]orchestration.WorkflowStatus
}

func (m *mockWorkflowStore) Create(ctx context.Context, workflow orchestration.Workflow) error {
	_ = ctx
	m.created = append(m.created, workflow)
	return nil
}

func (m *mockWorkflowStore) UpdateStatus(ctx context.Context, id string, status orchestration.WorkflowStatus) error {
	_ = ctx
	if m.statuses == nil {
		m.statuses = map[string]orchestration.WorkflowStatus{}
	}
	m.statuses[id] = status
	return nil
}

type mockWorkflowRunner struct {
	t     *testing.T
	calls []string
}

func (m *mockWorkflowRunner) RunStep(ctx context.Context, workflowID string, step orchestration.WorkflowStep, memory orchestration.SharedMemoryStore) (string, error) {
	if step.Sequence == 1 {
		_ = memory.Put(ctx, workflowID, "shared", "value")
	}
	if step.Sequence == 2 {
		if _, ok := memory.Get(ctx, workflowID, "shared"); !ok {
			m.t.Fatalf("expected shared memory to contain key")
		}
	}
	m.calls = append(m.calls, step.ID)
	return "job-" + step.ID, nil
}

func TestWorkflowRunIntegration(t *testing.T) {
	store := &mockWorkflowStore{statuses: map[string]orchestration.WorkflowStatus{}}
	memory := orchestration.NewMemoryStore()
	runner := &mockWorkflowRunner{t: t}
	service := orchestration.WorkflowService{
		Store:  store,
		Runner: runner,
		Memory: memory,
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
	if len(store.created) != 1 {
		t.Fatalf("expected workflow created")
	}
	if store.statuses[store.created[0].ID] != orchestration.WorkflowFinished {
		t.Fatalf("expected workflow finished")
	}
	if len(runner.calls) != 2 {
		t.Fatalf("expected two steps executed")
	}
}
