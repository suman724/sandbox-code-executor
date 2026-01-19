package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"control-plane/internal/api/handlers"
	"control-plane/internal/audit"
)

func TestAuditQueryIntegration(t *testing.T) {
	store := &audit.InMemoryStore{}
	now := time.Now().UTC()

	_ = store.Append(context.Background(), audit.Event{
		TenantID: "tenant-1",
		Action:   "job_created",
		Outcome:  "ok",
		Time:     now,
		Detail:   "job-1",
	})
	_ = store.Append(context.Background(), audit.Event{
		TenantID: "tenant-2",
		Action:   "job_created",
		Outcome:  "ok",
		Time:     now,
		Detail:   "job-2",
	})

	handler := handlers.AuditHandler{Store: store}
	req := httptest.NewRequest(http.MethodGet, "/audit/events?tenantId=tenant-1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}
	var events []audit.Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(events) != 1 || events[0].TenantID != "tenant-1" {
		t.Fatalf("expected tenant-1 events only")
	}
}
