package runtime

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type sessionProcess struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  *bufio.Reader
	runtime string
	mu      sync.Mutex
}

type replResponse struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error"`
}

func startSessionProcess(runtime string, workspaceDir string) (*sessionProcess, error) {
	normalized := strings.ToLower(strings.TrimSpace(runtime))
	var cmd *exec.Cmd
	if normalized == "python" || normalized == "python3" {
		cmd = exec.Command("python3", "-u", "-c", pythonReplScript)
	} else if normalized == "node" || normalized == "javascript" {
		cmd = exec.Command("node", "-e", nodeReplScript)
	} else {
		return nil, errors.New("unsupported runtime")
	}
	if workspaceDir != "" {
		if err := os.MkdirAll(workspaceDir, 0o750); err != nil {
			return nil, err
		}
		cmd.Dir = workspaceDir
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		return nil, err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		return nil, err
	}
	return &sessionProcess{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  bufio.NewReader(stdoutPipe),
		runtime: runtime,
	}, nil
}

func (p *sessionProcess) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stdin != nil {
		_ = p.stdin.Close()
	}
	if p.cmd != nil && p.cmd.Process != nil {
		return p.cmd.Process.Kill()
	}
	return nil
}

func (p *sessionProcess) RunStep(code string) (string, string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	payload, err := json.Marshal(map[string]string{"code": code})
	if err != nil {
		return "", "", err
	}
	if _, err := p.stdin.Write(append(payload, '\n')); err != nil {
		return "", "", err
	}
	line, err := p.stdout.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	var resp replResponse
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &resp); err != nil {
		return "", "", err
	}
	stderr := resp.Stderr
	if resp.Error != "" {
		if stderr != "" {
			stderr = stderr + "\n" + resp.Error
		} else {
			stderr = resp.Error
		}
	}
	return resp.Stdout, stderr, nil
}
