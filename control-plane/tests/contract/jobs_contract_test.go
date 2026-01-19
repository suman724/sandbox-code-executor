package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"control-plane/internal/api/handlers"
	"control-plane/internal/orchestration"
	"control-plane/internal/policy"
	"control-plane/internal/storage"
	"control-plane/pkg/client"
)

func TestJobsContractCreate(t *testing.T) {
	dataPlane := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"run_id":"run-1"}`))
	}))
	t.Cleanup(dataPlane.Close)

	store := &mockStore{job: storage.Job{ID: "job-1", Status: "queued"}}
	service := orchestration.JobService{
		Store:  store,
		Client: client.DataPlaneClient{BaseURL: dataPlane.URL, Client: dataPlane.Client()},
		Enforcer: orchestration.PolicyEnforcer{
			Evaluator: allowAllEvaluator{},
		},
	}
	handler := handlers.JobHandler{Service: service, Store: store}

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
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["id"] == "" || resp["status"] == "" {
		t.Fatalf("expected id and status")
	}
}

func TestJobsContractGet(t *testing.T) {
	store := &mockStore{job: storage.Job{ID: "job-1", Status: "running"}}
	handler := handlers.JobHandler{Store: store}

	req := httptest.NewRequest(http.MethodGet, "/jobs/job-1", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("jobId", "job-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}
}

type mockStore struct {
	job storage.Job
}

func (m *mockStore) Create(ctx context.Context, job storage.Job) error {
	m.job = job
	return nil
}

func (m *mockStore) Get(ctx context.Context, id string) (storage.Job, error) {
	_ = ctx
	return m.job, nil
}

func (m *mockStore) UpdateStatus(ctx context.Context, id string, status string) error {
	_ = id
	m.job.Status = status
	return nil
}

type allowAllEvaluator struct{}

func (allowAllEvaluator) Evaluate(ctx context.Context, input any) (policy.Decision, error) {
	_ = ctx
	_ = input
	return policy.Decision{Allowed: true}, nil
}
