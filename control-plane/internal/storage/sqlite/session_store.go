package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"control-plane/internal/storage"
)

type SessionStore struct {
	DB *sql.DB
}

func (s SessionStore) Create(ctx context.Context, session storage.Session) error {
	if s.DB == nil {
		return errors.New("nil db")
	}
	_, err := s.DB.ExecContext(ctx, `insert into sessions (id, status, runtime_id) values (?, ?, ?)`, session.ID, session.Status, session.RuntimeID)
	return err
}

func (s SessionStore) Get(ctx context.Context, id string) (storage.Session, error) {
	if s.DB == nil {
		return storage.Session{}, errors.New("nil db")
	}
	var session storage.Session
	err := s.DB.QueryRowContext(ctx, `select id, status, runtime_id from sessions where id = ?`, id).Scan(&session.ID, &session.Status, &session.RuntimeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Session{}, errors.New("session not found")
		}
		return storage.Session{}, err
	}
	return session, nil
}

func (s SessionStore) UpdateStatus(ctx context.Context, id string, status string) error {
	if s.DB == nil {
		return errors.New("nil db")
	}
	_, err := s.DB.ExecContext(ctx, `update sessions set status = ? where id = ?`, status, id)
	return err
}
