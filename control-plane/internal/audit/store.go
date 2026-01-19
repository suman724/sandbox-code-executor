package audit

import "context"

type Store interface {
	Append(ctx context.Context, event Event) error
	List(ctx context.Context) ([]Event, error)
}

type InMemoryStore struct {
	events []Event
}

func (s *InMemoryStore) Append(ctx context.Context, event Event) error {
	_ = ctx
	s.events = append(s.events, event)
	return nil
}

func (s *InMemoryStore) List(ctx context.Context) ([]Event, error) {
	_ = ctx
	return append([]Event(nil), s.events...), nil
}
