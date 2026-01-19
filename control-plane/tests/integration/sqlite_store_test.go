package integration

import (
	"context"
	"testing"

	"control-plane/internal/orchestration"
	"control-plane/internal/storage"
)

func TestSQLiteStoresIntegration(t *testing.T) {
	ctx := context.Background()
	stores, err := storage.NewStoreSet(ctx, "sqlite", "file:memdb1?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("open sqlite stores: %v", err)
	}
	if stores.Close != nil {
		defer func() {
			if err := stores.Close(); err != nil {
				t.Fatalf("close stores: %v", err)
			}
		}()
	}
	if stores.DB == nil {
		t.Fatalf("expected sqlite db handle")
	}

	job := storage.Job{ID: "job-1", Status: "queued"}
	if err := stores.JobStore.Create(ctx, job); err != nil {
		t.Fatalf("create job: %v", err)
	}
	if err := stores.JobStore.UpdateStatus(ctx, job.ID, "running"); err != nil {
		t.Fatalf("update job: %v", err)
	}
	gotJob, err := stores.JobStore.Get(ctx, job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}
	if gotJob.Status != "running" {
		t.Fatalf("expected job status running, got %s", gotJob.Status)
	}

	session := storage.Session{ID: "session-1", Status: "active"}
	if err := stores.SessionStore.Create(ctx, session); err != nil {
		t.Fatalf("create session: %v", err)
	}
	if err := stores.SessionStore.UpdateStatus(ctx, session.ID, "expired"); err != nil {
		t.Fatalf("update session: %v", err)
	}
	gotSession, err := stores.SessionStore.Get(ctx, session.ID)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if gotSession.Status != "expired" {
		t.Fatalf("expected session status expired, got %s", gotSession.Status)
	}

	if err := stores.PolicyStore.Upsert(ctx, storage.Policy{ID: "tenant-1:default", Version: 1}); err != nil {
		t.Fatalf("upsert policy: %v", err)
	}

	if err := stores.AuditStore.Append(ctx, storage.AuditEvent{ID: "event-1", Action: "job_accepted", Outcome: "ok"}); err != nil {
		t.Fatalf("append audit: %v", err)
	}
	events, err := stores.AuditStore.List(ctx)
	if err != nil {
		t.Fatalf("list audit: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}

	idempotency := orchestration.SQLiteIdempotencyStore{DB: stores.DB}
	if _, ok, err := idempotency.Get(ctx, "key-1"); err != nil || ok {
		t.Fatalf("expected empty idempotency entry")
	}
	if err := idempotency.Put(ctx, "key-1", "value-1"); err != nil {
		t.Fatalf("put idempotency: %v", err)
	}
	value, ok, err := idempotency.Get(ctx, "key-1")
	if err != nil {
		t.Fatalf("get idempotency: %v", err)
	}
	if !ok || value != "value-1" {
		t.Fatalf("expected idempotency value")
	}
}
