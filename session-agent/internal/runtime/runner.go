package runtime

import (
	"context"
	"errors"
	"strings"
	"sync"

	"shared/sessionagent"
)

type Session struct {
	ID           string
	Runtime      string
	Token        string
	WorkspaceDir string
	Process      *sessionProcess
}

type Runner struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewRunner() *Runner {
	return &Runner{sessions: make(map[string]*Session)}
}

func (r *Runner) RegisterSession(req sessionagent.SessionRegisterRequest) (*Session, error) {
	if req.SessionID == "" {
		return nil, errors.New("missing session id")
	}
	if strings.TrimSpace(req.Runtime) == "" {
		return nil, errors.New("missing runtime")
	}

	r.mu.Lock()
	if session, ok := r.sessions[req.SessionID]; ok {
		if req.Token != "" {
			session.Token = req.Token
		}
		if req.WorkspaceDir != "" {
			session.WorkspaceDir = req.WorkspaceDir
		}
		r.mu.Unlock()
		return session, nil
	}
	r.mu.Unlock()

	process, err := startSessionProcess(req.Runtime, req.WorkspaceDir)
	if err != nil {
		return nil, err
	}
	session := &Session{
		ID:           req.SessionID,
		Runtime:      req.Runtime,
		Token:        req.Token,
		WorkspaceDir: req.WorkspaceDir,
		Process:      process,
	}
	r.mu.Lock()
	r.sessions[req.SessionID] = session
	r.mu.Unlock()
	return session, nil
}

func (r *Runner) GetSession(sessionID string) (*Session, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, ok := r.sessions[sessionID]
	return session, ok
}

func (r *Runner) Authorize(sessionID string, token string) error {
	if token == "" {
		return errors.New("missing session token")
	}
	r.mu.RLock()
	session, ok := r.sessions[sessionID]
	r.mu.RUnlock()
	if !ok {
		return errors.New("session not registered")
	}
	if session.Token == "" || session.Token != token {
		return errors.New("invalid session token")
	}
	return nil
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

	session, ok := r.GetSession(req.SessionID)
	if !ok {
		return sessionagent.StepResult{}, errors.New("session not registered")
	}
	if session.Process == nil {
		return sessionagent.StepResult{}, errors.New("session process not available")
	}
	stdout, stderr, err := session.Process.RunStep(req.Code)
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

func (r *Runner) RemoveSession(sessionID string) {
	r.mu.Lock()
	session, ok := r.sessions[sessionID]
	if ok {
		delete(r.sessions, sessionID)
	}
	r.mu.Unlock()
	if ok && session.Process != nil {
		_ = session.Process.Close()
	}
}
