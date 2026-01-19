package orchestration

import "context"

type MemoryStore struct {
	items map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: map[string]string{}}
}

func (s *MemoryStore) Put(ctx context.Context, key string, value string) error {
	_ = ctx
	s.items[key] = value
	return nil
}

func (s *MemoryStore) Get(ctx context.Context, key string) (string, bool) {
	_ = ctx
	value, ok := s.items[key]
	return value, ok
}
