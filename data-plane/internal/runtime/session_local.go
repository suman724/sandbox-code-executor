package runtime

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"sync"
)

type LocalSessionRuntime struct {
	mu        sync.RWMutex
	processes map[string]*sessionProcess
}

type sessionProcess struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
}

func NewLocalSessionRuntime() *LocalSessionRuntime {
	return &LocalSessionRuntime{processes: map[string]*sessionProcess{}}
}

func (r *LocalSessionRuntime) StartSession(ctx context.Context, sessionID string, policyID string, workspaceRef string) (string, error) {
	_ = ctx
	_ = policyID
	_ = workspaceRef
	if sessionID == "" {
		return "", errors.New("missing session id")
	}
	cmd := exec.Command("sh")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		return "", err
	}
	r.mu.Lock()
	r.processes[sessionID] = &sessionProcess{cmd: cmd, stdin: stdin}
	r.mu.Unlock()
	return sessionID, nil
}

func (r *LocalSessionRuntime) RunStep(ctx context.Context, runtimeID string, command string) error {
	_ = ctx
	if runtimeID == "" {
		return errors.New("missing runtime id")
	}
	if command == "" {
		return errors.New("missing command")
	}
	r.mu.RLock()
	process, ok := r.processes[runtimeID]
	r.mu.RUnlock()
	if !ok {
		return errors.New("runtime not found")
	}
	_, err := io.WriteString(process.stdin, command+"\n")
	return err
}

func (r *LocalSessionRuntime) TerminateSession(ctx context.Context, runtimeID string) error {
	_ = ctx
	r.mu.Lock()
	process, ok := r.processes[runtimeID]
	if ok {
		delete(r.processes, runtimeID)
	}
	r.mu.Unlock()
	if !ok {
		return errors.New("runtime not found")
	}
	if process.stdin != nil {
		_ = process.stdin.Close()
	}
	if process.cmd != nil && process.cmd.Process != nil {
		_ = process.cmd.Process.Kill()
	}
	return nil
}
