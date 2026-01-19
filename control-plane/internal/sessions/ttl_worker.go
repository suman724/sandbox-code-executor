package sessions

import (
	"context"
	"errors"
	"time"
)

type ExpiredSessionStore interface {
	ListExpired(ctx context.Context, before time.Time) ([]Session, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

type CleanupFunc func(ctx context.Context, sessionID string) error

type TTLWorker struct {
	Store   ExpiredSessionStore
	Cleanup CleanupFunc
	Now     func() time.Time
}

func (w TTLWorker) Run(ctx context.Context) error {
	if w.Store == nil {
		return errors.New("missing session store")
	}
	now := time.Now
	if w.Now != nil {
		now = w.Now
	}
	expired, err := w.Store.ListExpired(ctx, now())
	if err != nil {
		return err
	}
	for _, session := range expired {
		if w.Cleanup != nil {
			_ = w.Cleanup(ctx, session.ID)
		}
		if err := w.Store.UpdateStatus(ctx, session.ID, string(StatusExpired)); err != nil {
			return err
		}
	}
	return nil
}
