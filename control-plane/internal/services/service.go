package services

import "time"

type Status string

const (
	StatusStarting Status = "starting"
	StatusRunning  Status = "running"
	StatusStopped  Status = "stopped"
)

type Service struct {
	ID        string
	TenantID  string
	PolicyID  string
	Status    Status
	StartedAt time.Time
	StoppedAt time.Time
	ProxyURL  string
}
