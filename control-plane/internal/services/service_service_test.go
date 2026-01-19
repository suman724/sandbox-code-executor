package services

import "testing"

func TestServiceLifecycle(t *testing.T) {
	svc := Service{ID: "svc-1", Status: StatusStarting}
	if svc.ID == "" || svc.Status == "" {
		t.Fatalf("expected service fields set")
	}
}
