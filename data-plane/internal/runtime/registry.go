package runtime

import (
	"errors"
	"os"
	"os/exec"
)

type Adapter interface {
	Run(code string) error
}

type Registry struct {
	Adapters map[string]Adapter
}

func NewRegistry() Registry {
	return Registry{Adapters: map[string]Adapter{}}
}

func DefaultRegistry() Registry {
	registry := NewRegistry()
	registry.Register("python", ExecAdapter{Command: "python3"})
	registry.Register("node", ExecAdapter{Command: "node"})
	return registry
}

func (r Registry) Register(language string, adapter Adapter) {
	if r.Adapters == nil {
		r.Adapters = map[string]Adapter{}
	}
	r.Adapters[language] = adapter
}

func (r Registry) Adapter(language string) (Adapter, bool) {
	adapter, ok := r.Adapters[language]
	return adapter, ok
}

func (r Registry) Supports(language string) bool {
	_, ok := r.Adapters[language]
	return ok
}

type ExecAdapter struct {
	Command string
	Args    []string
}

func (a ExecAdapter) Run(code string) error {
	if a.Command == "" {
		return errors.New("missing command")
	}
	file, err := os.CreateTemp("", "run-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err := file.WriteString(code); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	args := append([]string{}, a.Args...)
	args = append(args, file.Name())
	cmd := exec.Command(a.Command, args...)
	return cmd.Run()
}
