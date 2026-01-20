package sessions

import "time"

type Status string

const (
	StatusActive     Status = "active"
	StatusExpired    Status = "expired"
	StatusTerminated Status = "terminated"
)

type Session struct {
	ID           string
	TenantID     string
	AgentID      string
	PolicyID     string
	Runtime      string
	TTL          time.Duration
	ExpiresAt    time.Time
	Status       Status
	RuntimeID    string
	LastActivity time.Time
	Steps        []SessionStep
}

type SessionStep struct {
	ID         string
	SessionID  string
	Sequence   int
	Command    string
	Status     string
	StartedAt  time.Time
	FinishedAt time.Time
}
