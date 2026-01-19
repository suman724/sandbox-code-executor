package orchestration

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"control-plane/internal/policy"
	"control-plane/internal/storage"
	"control-plane/pkg/client"
)

type mockJobStore struct{}

func (mockJobStore) Create(ctx context.Context, job storage.Job) error {
	_ = ctx
	if job.ID == "" {
		return storageError("missing id")
	}
	return nil
}

func (mockJobStore) Get(ctx context.Context, id string) (storage.Job, error) {
	_ = ctx
	return storage.Job{ID: id, Status: string(JobQueued)}, nil
}

func (mockJobStore) UpdateStatus(ctx context.Context, id string, status string) error {
	_ = ctx
	_ = id
	_ = status
	return nil
}

type mockEvaluator struct {
	allowed bool
}

func (m mockEvaluator) Evaluate(ctx context.Context, input any) (policy.Decision, error) {
	_ = ctx
	_ = input
	return policy.Decision{Allowed: m.allowed}, nil
}

type storageError string

func (e storageError) Error() string { return string(e) }

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func stubClient(runID string) *http.Client {
	return &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			body := io.NopCloser(strings.NewReader(`{"run_id":"` + runID + `"}`))
			return &http.Response{
				StatusCode: http.StatusAccepted,
				Body:       body,
				Header:     make(http.Header),
			}, nil
		}),
	}
}

func TestJobServiceLifecycle(t *testing.T) {
	svc := JobService{
		Store:  mockJobStore{},
		Client: client.DataPlaneClient{BaseURL: "http://data-plane", Client: stubClient("run-1")},
		Enforcer: PolicyEnforcer{
			Evaluator: mockEvaluator{allowed: true},
		},
	}
	_, err := svc.CreateJob(context.Background(), Job{ID: "job-1", Language: "go", Status: JobQueued})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestJobServicePolicyDenied(t *testing.T) {
	svc := JobService{
		Store:  mockJobStore{},
		Client: client.DataPlaneClient{BaseURL: "http://data-plane", Client: stubClient("run-1")},
		Enforcer: PolicyEnforcer{
			Evaluator: mockEvaluator{allowed: false},
		},
	}
	_, err := svc.CreateJob(context.Background(), Job{ID: "job-1", Language: "go", Status: JobQueued})
	if err == nil {
		t.Fatalf("expected policy denial error")
	}
}
