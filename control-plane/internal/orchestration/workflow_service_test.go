package orchestration

import (
	"context"
	"testing"
	"time"
)

func TestWorkflowOrchestration(t *testing.T) {
	svc := WorkflowService{}
	wf := Workflow{ID: "wf-1", TenantID: "t-1", Status: WorkflowQueued, CreatedAt: time.Now()}
	if err := svc.Start(context.Background(), wf); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
