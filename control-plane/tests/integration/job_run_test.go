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
	"control-plane/internal/policy"
	"control-plane/internal/storage"
	"control-plane/pkg/client"
)

type mockStore struct {
	created []storage.Job
	updated []string
}

func (m *mockStore) Create(ctx context.Context, job storage.Job) error {
	_ = ctx
	m.created = append(m.created, job)
	return nil
}

func (m *mockStore) Get(ctx context.Context, id string) (storage.Job, error) {
	_ = ctx
	return storage.Job{ID: id, Status: string(orchestration.JobQueued)}, nil
}

func (m *mockStore) UpdateStatus(ctx context.Context, id string, status string) error {
	_ = ctx
	m.updated = append(m.updated, id+":"+status)
	return nil
}

type allowAllEvaluator struct{}

func (allowAllEvaluator) Evaluate(ctx context.Context, input any) (policy.Decision, error) {
	_ = ctx
	_ = input
	return policy.Decision{Allowed: true}, nil
}

func TestJobRunIntegration(t *testing.T) {
	dataPlane := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/runs" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"run_id":"run-123"}`))
	}))
	t.Cleanup(dataPlane.Close)

	store := &mockStore{}
	service := orchestration.JobService{
		Store:  store,
		Client: client.DataPlaneClient{BaseURL: dataPlane.URL, Client: dataPlane.Client()},
		Enforcer: orchestration.PolicyEnforcer{
			Evaluator: allowAllEvaluator{},
		},
	}
	handler := handlers.JobHandler{Service: service}

	payload := map[string]any{
		"tenantId": "tenant-1",
		"agentId":  "agent-1",
		"policyId": "policy-1",
		"language": "python",
		"code":     "print('ok')",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, rec.Code)
	}
	if len(store.created) != 1 {
		t.Fatalf("expected job to be created")
	}
}
