package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"control-plane/internal/storage"
)

type JobStore struct {
	DB *sql.DB
}

func (s JobStore) Create(ctx context.Context, job storage.Job) error {
	if s.DB == nil {
		return errors.New("nil db")
	}
	_, err := s.DB.ExecContext(ctx, `insert into jobs (id, status) values (?, ?)`, job.ID, job.Status)
	return err
}

func (s JobStore) Get(ctx context.Context, id string) (storage.Job, error) {
	if s.DB == nil {
		return storage.Job{}, errors.New("nil db")
	}
	var job storage.Job
	err := s.DB.QueryRowContext(ctx, `select id, status from jobs where id = ?`, id).Scan(&job.ID, &job.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Job{}, errors.New("job not found")
		}
		return storage.Job{}, err
	}
	return job, nil
}

func (s JobStore) UpdateStatus(ctx context.Context, id string, status string) error {
	if s.DB == nil {
		return errors.New("nil db")
	}
	_, err := s.DB.ExecContext(ctx, `update jobs set status = ? where id = ?`, status, id)
	return err
}
