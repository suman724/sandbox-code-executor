package runtime

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"shared/sessionagent"
)

type Session struct {
	ID      string
	Runtime string
}

type Runner struct {
	mu       sync.Mutex
	sessions map[string]*Session
}

func NewRunner() *Runner {
	return &Runner{sessions: make(map[string]*Session)}
}

func (r *Runner) EnsureSession(sessionID string, runtime string) *Session {
	r.mu.Lock()
	defer r.mu.Unlock()

	if session, ok := r.sessions[sessionID]; ok {
		return session
	}

	session := &Session{ID: sessionID, Runtime: runtime}
	r.sessions[sessionID] = session
	return session
}

func (r *Runner) GetSession(sessionID string) (*Session, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, ok := r.sessions[sessionID]
	return session, ok
}

func (r *Runner) RemoveSession(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, sessionID)
}

func (r *Runner) RunStep(ctx context.Context, req sessionagent.StepRequest) (sessionagent.StepResult, error) {
	if req.SessionID == "" {
		return sessionagent.StepResult{}, errors.New("missing session id")
	}
	if req.StepID == "" {
		return sessionagent.StepResult{}, errors.New("missing step id")
	}
	if strings.TrimSpace(req.Code) == "" {
		return sessionagent.StepResult{}, errors.New("missing step code")
	}

	session := r.EnsureSession(req.SessionID, req.Runtime)
	stdout, stderr, err := runCommand(ctx, session.Runtime, req.Code)
	status := sessionagent.StepStatusCompleted
	if err != nil {
		status = sessionagent.StepStatusFailed
		stderr = strings.TrimSpace(strings.Join([]string{stderr, err.Error()}, "\n"))
	}

	return sessionagent.StepResult{
		StepID: req.StepID,
		Status: status,
		Stdout: stdout,
		Stderr: stderr,
	}, nil
}

func runCommand(ctx context.Context, runtime string, code string) (string, string, error) {
	cmd := buildCommand(runtime, code)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String(), stderr.String(), err
	}
	return stdout.String(), stderr.String(), nil
}

func buildCommand(runtime string, code string) *exec.Cmd {
	normalized := strings.ToLower(strings.TrimSpace(runtime))
	if normalized == "python" || normalized == "python3" {
		return exec.Command("python3", "-c", code)
	}
	if normalized == "node" || normalized == "javascript" {
		return exec.Command("node", "-e", code)
	}
	return exec.Command("sh", "-c", fmt.Sprintf("%s", code))
}
