package runtime

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type LocalSessionRuntime struct {
	mu        sync.RWMutex
	processes map[string]*sessionProcess
}

type sessionProcess struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mu     sync.Mutex
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
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		return "", err
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		return "", err
	}
	r.mu.Lock()
	r.processes[sessionID] = &sessionProcess{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdoutPipe),
	}
	r.mu.Unlock()
	return sessionID, nil
}

func (r *LocalSessionRuntime) RunStep(ctx context.Context, runtimeID string, command string) (StepOutput, error) {
	_ = ctx
	if runtimeID == "" {
		return StepOutput{}, errors.New("missing runtime id")
	}
	if command == "" {
		return StepOutput{}, errors.New("missing command")
	}
	r.mu.RLock()
	process, ok := r.processes[runtimeID]
	r.mu.RUnlock()
	if !ok {
		return StepOutput{}, errors.New("runtime not found")
	}
	process.mu.Lock()
	defer process.mu.Unlock()
	token := fmt.Sprintf("step-%d", time.Now().UTC().UnixNano())
	payload := fmt.Sprintf("out=$(mktemp) err=$(mktemp); (%s) 1>\"$out\" 2>\"$err\"; printf '__STDOUT__%s\\n'; cat \"$out\"; printf '\\n__STDERR__%s\\n'; cat \"$err\"; printf '\\n__END__%s\\n'; rm -f \"$out\" \"$err\"\n", command, token, token, token)
	if _, err := io.WriteString(process.stdin, payload); err != nil {
		return StepOutput{}, err
	}
	stdout, stderr, err := readStepOutput(process.stdout, token)
	if err != nil {
		return StepOutput{}, err
	}
	return StepOutput{Stdout: stdout, Stderr: stderr}, nil
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

func readStepOutput(reader *bufio.Reader, token string) (string, string, error) {
	stdoutMarker := "__STDOUT__" + token
	stderrMarker := "__STDERR__" + token
	endMarker := "__END__" + token

	if _, err := readUntilMarker(reader, stdoutMarker); err != nil {
		return "", "", err
	}
	stdout, err := readUntilMarker(reader, stderrMarker)
	if err != nil {
		return "", "", err
	}
	stderr, err := readUntilMarker(reader, endMarker)
	if err != nil {
		return "", "", err
	}
	return stdout, stderr, nil
}

func readUntilMarker(reader *bufio.Reader, marker string) (string, error) {
	var buf bytes.Buffer
	for {
		line, err := reader.ReadString('\n')
		if err != nil && len(line) == 0 {
			return "", err
		}
		text := line
		text = strings.TrimRight(text, "\r\n")
		if text == marker {
			return buf.String(), nil
		}
		if len(line) > 0 {
			buf.WriteString(line)
		}
		if err != nil {
			return buf.String(), err
		}
	}
}
