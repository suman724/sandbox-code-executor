package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"control-plane/internal/storage"
)

type JobStore struct {
	Pool *pgxpool.Pool
}

type SessionStore struct {
	Pool *pgxpool.Pool
}

func (s JobStore) Create(ctx context.Context, job storage.Job) error {
	if s.Pool == nil {
		return errors.New("nil pool")
	}
	_, err := s.Pool.Exec(ctx, `insert into jobs (id, status) values ($1, $2)`, job.ID, job.Status)
	return err
}

func (s JobStore) Get(ctx context.Context, id string) (storage.Job, error) {
	if s.Pool == nil {
		return storage.Job{}, errors.New("nil pool")
	}
	var job storage.Job
	err := s.Pool.QueryRow(ctx, `select id, status from jobs where id = $1`, id).Scan(&job.ID, &job.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.Job{}, errors.New("job not found")
		}
		return storage.Job{}, err
	}
	return job, nil
}

func (s JobStore) UpdateStatus(ctx context.Context, id string, status string) error {
	if s.Pool == nil {
		return errors.New("nil pool")
	}
	_, err := s.Pool.Exec(ctx, `update jobs set status = $1 where id = $2`, status, id)
	return err
}

func (s SessionStore) Create(ctx context.Context, session storage.Session) error {
	if s.Pool == nil {
		return errors.New("nil pool")
	}
	_, err := s.Pool.Exec(ctx, `insert into sessions (id, status) values ($1, $2)`, session.ID, session.Status)
	return err
}

func (s SessionStore) Get(ctx context.Context, id string) (storage.Session, error) {
	if s.Pool == nil {
		return storage.Session{}, errors.New("nil pool")
	}
	var session storage.Session
	err := s.Pool.QueryRow(ctx, `select id, status from sessions where id = $1`, id).Scan(&session.ID, &session.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.Session{}, errors.New("session not found")
		}
		return storage.Session{}, err
	}
	return session, nil
}

func (s SessionStore) UpdateStatus(ctx context.Context, id string, status string) error {
	if s.Pool == nil {
		return errors.New("nil pool")
	}
	_, err := s.Pool.Exec(ctx, `update sessions set status = $1 where id = $2`, status, id)
	return err
}
