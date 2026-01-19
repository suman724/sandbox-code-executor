package orchestration

import (
	"context"
	"sync"
)

type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: map[string]map[string]string{}}
}

func (s *MemoryStore) Put(ctx context.Context, workflowID string, key string, value string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.items[workflowID] == nil {
		s.items[workflowID] = map[string]string{}
	}
	s.items[workflowID][key] = value
	return nil
}

func (s *MemoryStore) Get(ctx context.Context, workflowID string, key string) (string, bool) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := s.items[workflowID]
	if items == nil {
		return "", false
	}
	value, ok := items[key]
	return value, ok
}
