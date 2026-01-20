package sessions

import (
	"context"
	"errors"

	"control-plane/internal/storage"
)

type StorageStepStore struct {
	Store storage.SessionStepStore
}

func (s StorageStepStore) AppendStep(ctx context.Context, step SessionStep) error {
	if s.Store == nil {
		return errors.New("missing step store")
	}
	return s.Store.Append(ctx, storage.SessionStep{
		ID:        step.ID,
		SessionID: step.SessionID,
		Command:   step.Command,
		Status:    step.Status,
	})
}

func (s StorageStepStore) ListSteps(ctx context.Context, sessionID string) ([]SessionStep, error) {
	if s.Store == nil {
		return nil, errors.New("missing step store")
	}
	steps, err := s.Store.List(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	result := make([]SessionStep, 0, len(steps))
	for _, step := range steps {
		result = append(result, SessionStep{
			ID:        step.ID,
			SessionID: step.SessionID,
			Command:   step.Command,
			Status:    step.Status,
		})
	}
	return result, nil
}
