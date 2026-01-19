package orchestration

import "time"

type WorkflowStatus string

const (
	WorkflowQueued   WorkflowStatus = "queued"
	WorkflowRunning  WorkflowStatus = "running"
	WorkflowFinished WorkflowStatus = "finished"
	WorkflowFailed   WorkflowStatus = "failed"
)

type WorkflowStepStatus string

const (
	WorkflowStepQueued    WorkflowStepStatus = "queued"
	WorkflowStepRunning   WorkflowStepStatus = "running"
	WorkflowStepSucceeded WorkflowStepStatus = "succeeded"
	WorkflowStepFailed    WorkflowStepStatus = "failed"
)

type Workflow struct {
	ID          string
	TenantID    string
	Status      WorkflowStatus
	CreatedAt   time.Time
	StartedAt   time.Time
	CompletedAt time.Time
	Steps       []WorkflowStep
}

type WorkflowStep struct {
	ID         string
	WorkflowID string
	Sequence   int
	AgentID    string
	JobID      string
	Status     WorkflowStepStatus
	StartedAt  time.Time
	EndedAt    time.Time
}
