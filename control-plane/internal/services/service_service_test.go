package services

import (
	"context"
	"testing"
	"time"
)

type mockStore struct {
	last Service
}

func (m *mockStore) Create(ctx context.Context, service Service) error {
	_ = ctx
	m.last = service
	return nil
}

func (m *mockStore) UpdateStatus(ctx context.Context, id string, status Status, proxyURL string) error {
	_ = ctx
	m.last.ID = id
	m.last.Status = status
	m.last.ProxyURL = proxyURL
	return nil
}

type mockRunner struct {
	proxyURL string
}

func (m mockRunner) Start(ctx context.Context, service Service) (string, error) {
	_ = ctx
	_ = service
	return m.proxyURL, nil
}

func (m mockRunner) Stop(ctx context.Context, serviceID string) error {
	_ = ctx
	_ = serviceID
	return nil
}

func TestServiceLifecycle(t *testing.T) {
	store := &mockStore{}
	runner := mockRunner{proxyURL: "http://proxy/service-1"}
	service := ServiceService{
		Store:   store,
		Runner:  runner,
		NowFunc: func() time.Time { return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) },
	}

	svc, err := service.Start(context.Background(), Service{ID: "svc-1", TenantID: "tenant-1", PolicyID: "policy-1"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if svc.Status != StatusRunning || svc.ProxyURL == "" {
		t.Fatalf("expected service running with proxy url")
	}
	if store.last.Status != StatusRunning {
		t.Fatalf("expected store updated to running")
	}

	if err := service.Stop(context.Background(), "svc-1"); err != nil {
		t.Fatalf("expected stop to succeed, got %v", err)
	}
	if store.last.Status != StatusStopped {
		t.Fatalf("expected store updated to stopped")
	}
}
