package orchestration

import "time"

type WorkflowStatus string

const (
	WorkflowQueued   WorkflowStatus = "queued"
	WorkflowRunning  WorkflowStatus = "running"
	WorkflowFinished WorkflowStatus = "finished"
)

type Workflow struct {
	ID        string
	TenantID  string
	Status    WorkflowStatus
	CreatedAt time.Time
}
