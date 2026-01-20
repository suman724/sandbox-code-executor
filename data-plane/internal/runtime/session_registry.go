package runtime

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type SessionRegistry interface {
	Put(sessionID string, runtimeID string) error
	Get(sessionID string) (string, bool)
	Delete(sessionID string)
}

type InMemorySessionRegistry struct {
	mu    sync.RWMutex
	items map[string]string
}

func NewInMemorySessionRegistry() *InMemorySessionRegistry {
	return &InMemorySessionRegistry{items: map[string]string{}}
}

func (r *InMemorySessionRegistry) Put(sessionID string, runtimeID string) error {
	if sessionID == "" || runtimeID == "" {
		return errors.New("missing session or runtime id")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[sessionID] = runtimeID
	return nil
}

func (r *InMemorySessionRegistry) Get(sessionID string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	runtimeID, ok := r.items[sessionID]
	return runtimeID, ok
}

func (r *InMemorySessionRegistry) Delete(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, sessionID)
}

type FileSessionRegistry struct {
	mu    sync.RWMutex
	path  string
	items map[string]string
}

func NewFileSessionRegistry(path string) (*FileSessionRegistry, error) {
	if path == "" {
		return nil, errors.New("missing registry path")
	}
	registry := &FileSessionRegistry{
		path:  path,
		items: map[string]string{},
	}
	if err := registry.load(); err != nil {
		return nil, err
	}
	return registry, nil
}

func (r *FileSessionRegistry) Put(sessionID string, runtimeID string) error {
	if sessionID == "" || runtimeID == "" {
		return errors.New("missing session or runtime id")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[sessionID] = runtimeID
	return r.persistLocked()
}

func (r *FileSessionRegistry) Get(sessionID string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	runtimeID, ok := r.items[sessionID]
	return runtimeID, ok
}

func (r *FileSessionRegistry) Delete(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, sessionID)
	_ = r.persistLocked()
}

func (r *FileSessionRegistry) load() error {
	content, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(content) == 0 {
		return nil
	}
	return json.Unmarshal(content, &r.items)
}

func (r *FileSessionRegistry) persistLocked() error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}
	temp := r.path + ".tmp"
	payload, err := json.Marshal(r.items)
	if err != nil {
		return err
	}
	if err := os.WriteFile(temp, payload, 0o644); err != nil {
		return err
	}
	return os.Rename(temp, r.path)
}
