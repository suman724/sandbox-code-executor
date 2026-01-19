package orchestration

import (
	"context"
	"testing"
	"time"
)

func TestWorkflowOrchestration(t *testing.T) {
	store := &mockWorkflowStore{statuses: map[string]WorkflowStatus{}}
	memory := NewMemoryStore()
	runner := &mockWorkflowRunner{t: t}
	now := time.Date(2026, 1, 19, 12, 0, 0, 0, time.UTC)
	svc := WorkflowService{
		Store:  store,
		Runner: runner,
		Memory: memory,
		Now:    func() time.Time { return now },
	}
	wf := Workflow{
		ID:       "wf-1",
		TenantID: "t-1",
		Steps: []WorkflowStep{
			{AgentID: "agent-1", Sequence: 1},
			{AgentID: "agent-2", Sequence: 2},
		},
	}
	if err := svc.Start(context.Background(), wf); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(store.created) != 1 {
		t.Fatalf("expected workflow to be created")
	}
	if store.statuses["wf-1"] != WorkflowFinished {
		t.Fatalf("expected workflow to finish, got %s", store.statuses["wf-1"])
	}
	if len(runner.calls) != 2 || runner.calls[0] == runner.calls[1] {
		t.Fatalf("expected sequential workflow steps, got %v", runner.calls)
	}
}

type mockWorkflowStore struct {
	created  []Workflow
	statuses map[string]WorkflowStatus
}

func (m *mockWorkflowStore) Create(ctx context.Context, workflow Workflow) error {
	_ = ctx
	m.created = append(m.created, workflow)
	return nil
}

func (m *mockWorkflowStore) UpdateStatus(ctx context.Context, id string, status WorkflowStatus) error {
	_ = ctx
	m.statuses[id] = status
	return nil
}

type mockWorkflowRunner struct {
	t     *testing.T
	calls []string
}

func (m *mockWorkflowRunner) RunStep(ctx context.Context, workflowID string, step WorkflowStep, memory SharedMemoryStore) (string, error) {
	if step.Sequence == 1 {
		_ = memory.Put(ctx, workflowID, "shared", "payload")
	}
	if step.Sequence == 2 {
		if _, ok := memory.Get(ctx, workflowID, "shared"); !ok {
			m.t.Fatalf("expected shared memory to contain key")
		}
	}
	m.calls = append(m.calls, step.ID)
	return "job-" + step.ID, nil
}
