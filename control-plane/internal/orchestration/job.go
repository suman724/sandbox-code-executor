package orchestration

type JobStatus string

const (
	JobQueued   JobStatus = "queued"
	JobRunning  JobStatus = "running"
	JobFailed   JobStatus = "failed"
	JobFinished JobStatus = "finished"
)

type Job struct {
	ID           string
	TenantID     string
	AgentID      string
	PolicyID     string
	Language     string
	Code         string
	Workspace    string
	Status       JobStatus
	OutputRef    string
	ErrorRef     string
	ArtifactRefs []string
}
