package sessions

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"control-plane/internal/orchestration"
	"control-plane/internal/policy"
	"control-plane/internal/storage"
	"control-plane/pkg/client"
)

type mockSessionStore struct{}

func (mockSessionStore) Create(ctx context.Context, session storage.Session) error {
	_ = ctx
	if session.ID == "" {
		return storageError("missing id")
	}
	return nil
}

func (mockSessionStore) Get(ctx context.Context, id string) (storage.Session, error) {
	_ = ctx
	return storage.Session{ID: id, Status: string(StatusActive)}, nil
}

func (mockSessionStore) UpdateStatus(ctx context.Context, id string, status string) error {
	_ = ctx
	_ = id
	_ = status
	return nil
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

type mockEvaluator struct {
	allowed bool
}

func (m mockEvaluator) Evaluate(ctx context.Context, input any) (policy.Decision, error) {
	_ = ctx
	_ = input
	return policy.Decision{Allowed: m.allowed}, nil
}

func TestSessionLifecycle(t *testing.T) {
	svc := Service{
		Store:  mockSessionStore{},
		Client: client.DataPlaneClient{BaseURL: "http://data-plane", Client: stubClient("run-1")},
		Enforcer: orchestration.PolicyEnforcer{
			Evaluator: mockEvaluator{allowed: true},
		},
	}
	_, err := svc.CreateSession(context.Background(), Session{ID: "s-1", TTL: time.Minute, Status: StatusActive})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSessionPolicyDenied(t *testing.T) {
	svc := Service{
		Store:  mockSessionStore{},
		Client: client.DataPlaneClient{BaseURL: "http://data-plane", Client: stubClient("run-1")},
		Enforcer: orchestration.PolicyEnforcer{
			Evaluator: mockEvaluator{allowed: false},
		},
	}
	_, err := svc.CreateSession(context.Background(), Session{ID: "s-1", TTL: time.Minute, Status: StatusActive})
	if err == nil {
		t.Fatalf("expected policy denial error")
	}
}
