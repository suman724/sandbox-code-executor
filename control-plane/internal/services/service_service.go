package services

import (
	"context"
	"errors"
	"time"
)

type Store interface {
	Create(ctx context.Context, service Service) error
	UpdateStatus(ctx context.Context, id string, status Status, proxyURL string) error
}

type Starter interface {
	Start(ctx context.Context, service Service) (string, error)
	Stop(ctx context.Context, serviceID string) error
}

type ServiceService struct {
	Store   Store
	Runner  Starter
	NowFunc func() time.Time
}

func (s ServiceService) Start(ctx context.Context, service Service) (Service, error) {
	if service.ID == "" {
		return Service{}, errors.New("missing service id")
	}
	if service.TenantID == "" || service.PolicyID == "" {
		return Service{}, errors.New("missing tenant or policy id")
	}
	if s.Store == nil || s.Runner == nil {
		return Service{}, errors.New("missing store or runner")
	}
	now := time.Now
	if s.NowFunc != nil {
		now = s.NowFunc
	}
	service.Status = StatusStarting
	service.StartedAt = now()
	if err := s.Store.Create(ctx, service); err != nil {
		return Service{}, err
	}
	proxyURL, err := s.Runner.Start(ctx, service)
	if err != nil {
		_ = s.Store.UpdateStatus(ctx, service.ID, StatusStopped, "")
		return Service{}, err
	}
	service.Status = StatusRunning
	service.ProxyURL = proxyURL
	if err := s.Store.UpdateStatus(ctx, service.ID, StatusRunning, proxyURL); err != nil {
		return Service{}, err
	}
	return service, nil
}

func (s ServiceService) Stop(ctx context.Context, serviceID string) error {
	if serviceID == "" {
		return errors.New("missing service id")
	}
	if s.Runner == nil || s.Store == nil {
		return errors.New("missing store or runner")
	}
	if err := s.Runner.Stop(ctx, serviceID); err != nil {
		return err
	}
	return s.Store.UpdateStatus(ctx, serviceID, StatusStopped, "")
}
