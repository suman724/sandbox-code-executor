package sessions

import "time"

type Status string

const (
	StatusActive    Status = "active"
	StatusExpired   Status = "expired"
	StatusTerminated Status = "terminated"
)

type Session struct {
	ID        string
	TenantID  string
	PolicyID  string
	TTL       time.Duration
	ExpiresAt time.Time
	Status    Status
}
