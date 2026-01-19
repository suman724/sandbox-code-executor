package policy

import (
	"context"
	"errors"
	"sync"
)

type Policy struct {
	ID      string
	Version int
	Ruleset string
}

type Store interface {
	Upsert(ctx context.Context, policy Policy) error
}

var ErrStalePolicyVersion = errors.New("stale policy version")

type InMemoryStore struct {
	mu      sync.RWMutex
	entries map[string]Policy
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{entries: map[string]Policy{}}
}

func (s *InMemoryStore) Upsert(ctx context.Context, policy Policy) error {
	_ = ctx
	if policy.ID == "" {
		return errors.New("missing policy id")
	}
	if policy.Version <= 0 {
		return errors.New("invalid policy version")
	}
	if policy.Ruleset == "" {
		return errors.New("missing ruleset")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.entries[policy.ID]; ok {
		if policy.Version <= existing.Version {
			return ErrStalePolicyVersion
		}
	}
	s.entries[policy.ID] = policy
	return nil
}
