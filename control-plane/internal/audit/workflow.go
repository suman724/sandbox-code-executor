package audit

import (
	"context"
	"time"
)

func WorkflowEvent(ctx context.Context, logger Logger, tenantID string, workflowID string, action string) error {
	return logger.Log(ctx, Event{
		TenantID: tenantID,
		Action:   action,
		Outcome:  "ok",
		Time:     time.Now(),
		Detail:   workflowID,
	})
}
