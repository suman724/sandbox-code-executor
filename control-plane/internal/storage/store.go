package storage

import "context"

type Job struct {
	ID     string
	Status string
}

type Session struct {
	ID     string
	Status string
}

type Policy struct {
	ID      string
	Version int
}

type AuditEvent struct {
	ID      string
	Action  string
	Outcome string
}

type Artifact struct {
	ID         string
	Name       string
	SizeBytes  int64
	Checksum   string
	StorageURI string
}

type JobStore interface {
	Create(ctx context.Context, job Job) error
	Get(ctx context.Context, id string) (Job, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

type SessionStore interface {
	Create(ctx context.Context, session Session) error
	Get(ctx context.Context, id string) (Session, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

type SessionStepStore interface {
	Append(ctx context.Context, step SessionStep) error
	List(ctx context.Context, sessionID string) ([]SessionStep, error)
}

type SessionStep struct {
	ID        string
	SessionID string
	Command   string
	Status    string
}

type PolicyStore interface {
	Upsert(ctx context.Context, policy Policy) error
}

type AuditStore interface {
	Append(ctx context.Context, event AuditEvent) error
	List(ctx context.Context) ([]AuditEvent, error)
}
