package sessionagent

type StepRequest struct {
	SessionID string `json:"sessionId"`
	StepID    string `json:"stepId"`
	Code      string `json:"code"`
	Runtime   string `json:"runtime,omitempty"`
}

type SessionRegisterRequest struct {
	SessionID    string `json:"sessionId"`
	Runtime      string `json:"runtime"`
	Token        string `json:"token,omitempty"`
	WorkspaceDir string `json:"workspaceDir,omitempty"`
}

type SessionRegisterResponse struct {
	SessionID string `json:"sessionId"`
	Status    string `json:"status"`
}

type SessionTerminateResponse struct {
	SessionID string `json:"sessionId"`
	Status    string `json:"status"`
}

type StepResult struct {
	StepID   string `json:"stepId"`
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode,omitempty"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

const (
	StepStatusCompleted = "completed"
	StepStatusFailed    = "failed"
)
