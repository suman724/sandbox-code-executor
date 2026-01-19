package audit

import (
	"context"

	"control-plane/internal/storage"
)

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

type StorageAdapter struct {
	Store storage.AuditStore
}

func (s StorageAdapter) Append(ctx context.Context, event Event) error {
	if s.Store == nil {
		return ErrStoreUnavailable
	}
	return s.Store.Append(ctx, storage.AuditEvent{
		ID:      event.Detail,
		Action:  event.Action,
		Outcome: event.Outcome,
	})
}

func (s StorageAdapter) List(ctx context.Context) ([]Event, error) {
	if s.Store == nil {
		return nil, ErrStoreUnavailable
	}
	items, err := s.Store.List(ctx)
	if err != nil {
		return nil, err
	}
	events := make([]Event, 0, len(items))
	for _, item := range items {
		events = append(events, Event{
			Action:  item.Action,
			Outcome: item.Outcome,
			Detail:  item.ID,
		})
	}
	return events, nil
}

var ErrStoreUnavailable = storageError("audit store unavailable")

type storageError string

func (e storageError) Error() string { return string(e) }
