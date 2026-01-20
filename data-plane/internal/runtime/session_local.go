package runtime

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
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
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  *bufio.Reader
	runtime string
	repl    bool
	mu     sync.Mutex
}

const pythonReplScript = `import contextlib
import io
import json
import sys
import traceback

globals_ns = {"__name__": "__main__"}

for line in sys.stdin:
    line = line.rstrip("\n")
    if not line:
        continue
    try:
        req = json.loads(line)
    except Exception:
        continue
    code = req.get("code", "")
    out = io.StringIO()
    err = io.StringIO()
    failure = ""
    try:
        with contextlib.redirect_stdout(out), contextlib.redirect_stderr(err):
            exec(code, globals_ns)
    except Exception:
        failure = traceback.format_exc()
    resp = {"stdout": out.getvalue(), "stderr": err.getvalue(), "error": failure}
    sys.stdout.write(json.dumps(resp) + "\n")
    sys.stdout.flush()
`

const nodeReplScript = `const readline = require("readline");
const vm = require("vm");

const context = vm.createContext({
  console,
  require,
  process,
  Buffer,
  setTimeout,
  setInterval,
  clearTimeout,
  clearInterval,
});

const rl = readline.createInterface({
  input: process.stdin,
  crlfDelay: Infinity,
});

rl.on("line", (line) => {
  if (!line) {
    return;
  }
  let req;
  try {
    req = JSON.parse(line);
  } catch (err) {
    return;
  }
  let stdout = "";
  let stderr = "";
  let error = "";

  const originalStdoutWrite = process.stdout.write.bind(process.stdout);
  const originalStderrWrite = process.stderr.write.bind(process.stderr);
  const originalConsoleLog = console.log;
  const originalConsoleError = console.error;

  process.stdout.write = (chunk, encoding, cb) => {
    stdout += chunk instanceof Buffer ? chunk.toString() : chunk;
    if (typeof cb === "function") {
      cb();
    }
    return true;
  };
  process.stderr.write = (chunk, encoding, cb) => {
    stderr += chunk instanceof Buffer ? chunk.toString() : chunk;
    if (typeof cb === "function") {
      cb();
    }
    return true;
  };
  console.log = (...args) => {
    stdout += args.join(" ") + "\n";
  };
  console.error = (...args) => {
    stderr += args.join(" ") + "\n";
  };

  try {
    vm.runInContext(req.code || "", context);
  } catch (err) {
    error = err && err.stack ? err.stack : String(err);
  }

  process.stdout.write = originalStdoutWrite;
  process.stderr.write = originalStderrWrite;
  console.log = originalConsoleLog;
  console.error = originalConsoleError;

  const resp = JSON.stringify({ stdout, stderr, error });
  originalStdoutWrite(resp + "\n");
});
`

func NewLocalSessionRuntime() *LocalSessionRuntime {
	return &LocalSessionRuntime{processes: map[string]*sessionProcess{}}
}

func (r *LocalSessionRuntime) StartSession(ctx context.Context, sessionID string, policyID string, workspaceRef string, runtime string) (string, error) {
	_ = ctx
	_ = policyID
	_ = workspaceRef
	if sessionID == "" {
		return "", errors.New("missing session id")
	}
	normalized := strings.ToLower(strings.TrimSpace(runtime))
	cmd := exec.Command("sh")
	repl := false
	if normalized == "python" || normalized == "python3" {
		cmd = exec.Command("python3", "-u", "-c", pythonReplScript)
		repl = true
	} else if normalized == "node" || normalized == "javascript" {
		cmd = exec.Command("node", "-e", nodeReplScript)
		repl = true
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		return "", err
	}
	if repl {
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = cmd.Stdout
	}
	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		return "", err
	}
	r.mu.Lock()
	r.processes[sessionID] = &sessionProcess{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  bufio.NewReader(stdoutPipe),
		runtime: runtime,
		repl:    repl,
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
	if process.repl {
		return runStepRepl(process, command)
	}
	token := fmt.Sprintf("step-%d", time.Now().UTC().UnixNano())
	wrapped := command
	if strings.EqualFold(process.runtime, "python") {
		heredoc := "PY" + strings.ReplaceAll(token, "-", "")
		wrapped = fmt.Sprintf("python3 - <<'%s'\n%s\n%s", heredoc, command, heredoc)
	} else if strings.EqualFold(process.runtime, "node") || strings.EqualFold(process.runtime, "javascript") {
		heredoc := "JS" + strings.ReplaceAll(token, "-", "")
		wrapped = fmt.Sprintf("node - <<'%s'\n%s\n%s", heredoc, command, heredoc)
	}
	payload := fmt.Sprintf("out=$(mktemp) err=$(mktemp); (%s) 1>\"$out\" 2>\"$err\"; printf '__STDOUT__%s\\n'; cat \"$out\"; printf '\\n__STDERR__%s\\n'; cat \"$err\"; printf '\\n__END__%s\\n'; rm -f \"$out\" \"$err\"\n", wrapped, token, token, token)
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

type replResponse struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error"`
}

func runStepRepl(process *sessionProcess, command string) (StepOutput, error) {
	payload, err := json.Marshal(map[string]string{"code": command})
	if err != nil {
		return StepOutput{}, err
	}
	if _, err := process.stdin.Write(append(payload, '\n')); err != nil {
		return StepOutput{}, err
	}
	line, err := process.stdout.ReadString('\n')
	if err != nil {
		return StepOutput{}, err
	}
	var resp replResponse
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &resp); err != nil {
		return StepOutput{}, err
	}
	stderr := resp.Stderr
	if resp.Error != "" {
		if stderr != "" {
			stderr = stderr + "\n" + resp.Error
		} else {
			stderr = resp.Error
		}
	}
	return StepOutput{Stdout: resp.Stdout, Stderr: stderr}, nil
}
