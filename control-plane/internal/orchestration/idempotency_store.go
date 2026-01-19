package orchestration

import "context"

type IdempotencyStore interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Put(ctx context.Context, key string, value string) error
}

type InMemoryIdempotencyStore struct {
	items map[string]string
}

func NewInMemoryIdempotencyStore() *InMemoryIdempotencyStore {
	return &InMemoryIdempotencyStore{items: map[string]string{}}
}

func (s *InMemoryIdempotencyStore) Get(ctx context.Context, key string) (string, bool, error) {
	_ = ctx
	value, ok := s.items[key]
	return value, ok, nil
}

func (s *InMemoryIdempotencyStore) Put(ctx context.Context, key string, value string) error {
	_ = ctx
	s.items[key] = value
	return nil
}
