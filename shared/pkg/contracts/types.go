package contracts

type JobCreate struct {
	TenantID  TenantID `json:"tenant_id"`
	AgentID   AgentID  `json:"agent_id"`
	PolicyID  PolicyID `json:"policy_id"`
	Language  string   `json:"language"`
	Code      string   `json:"code"`
	Artifacts []string `json:"artifacts,omitempty"`
}

type Job struct {
	ID           JobID    `json:"id"`
	Status       string   `json:"status"`
	ExitStatus   int      `json:"exit_status,omitempty"`
	OutputRef    string   `json:"output_ref,omitempty"`
	ErrorRef     string   `json:"error_ref,omitempty"`
	ArtifactRefs []string `json:"artifact_refs,omitempty"`
}

type SessionCreate struct {
	TenantID   TenantID `json:"tenant_id"`
	AgentID    AgentID  `json:"agent_id"`
	PolicyID   PolicyID `json:"policy_id"`
	TTLSeconds int      `json:"ttl_seconds"`
}

type Session struct {
	ID        SessionID `json:"id"`
	Status    string    `json:"status"`
	ExpiresAt string    `json:"expires_at"`
}

type SessionStepCreate struct {
	Command string `json:"command"`
}

type SessionStep struct {
	ID     SessionStepID `json:"id"`
	Status string        `json:"status"`
}

type ArtifactUpload struct {
	TenantID  TenantID `json:"tenant_id"`
	Name      string   `json:"name"`
	SizeBytes int64    `json:"size_bytes"`
}

type Artifact struct {
	ID          ArtifactID `json:"id"`
	Name        string     `json:"name"`
	SizeBytes   int64      `json:"size_bytes"`
	DownloadURL string     `json:"download_url,omitempty"`
}

type Policy struct {
	TenantID TenantID `json:"tenant_id"`
	Name     string   `json:"name"`
	Version  int      `json:"version"`
	Ruleset  string   `json:"ruleset"`
}

type AuditEvent struct {
	ID           AuditEventID `json:"id"`
	Timestamp    string       `json:"timestamp"`
	ActorID      AgentID      `json:"actor_id"`
	Action       string       `json:"action"`
	ResourceType string       `json:"resource_type"`
	ResourceID   string       `json:"resource_id"`
	Outcome      string       `json:"outcome"`
}

type WorkflowCreate struct {
	TenantID TenantID `json:"tenant_id"`
	Steps    []string `json:"steps"`
}

type Workflow struct {
	ID     WorkflowID `json:"id"`
	Status string     `json:"status"`
}

type ServiceCreate struct {
	TenantID TenantID `json:"tenant_id"`
	PolicyID PolicyID `json:"policy_id"`
}

type Service struct {
	ID       ServiceID `json:"id"`
	Status   string    `json:"status"`
	ProxyURL string    `json:"proxy_url,omitempty"`
}
