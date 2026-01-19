package object

import (
	"context"
	"errors"
	"strings"

	"control-plane/internal/storage"
)

type ArtifactStore struct {
	BaseURL string
}

func (s ArtifactStore) Put(ctx context.Context, artifact storage.Artifact) error {
	_ = ctx
	if artifact.ID == "" {
		return errors.New("missing artifact id")
	}
	return nil
}

func (s ArtifactStore) Get(ctx context.Context, id string) (storage.Artifact, error) {
	_ = ctx
	if id == "" {
		return storage.Artifact{}, errors.New("missing artifact id")
	}
	return storage.Artifact{ID: id}, nil
}

func (s ArtifactStore) SignedDownloadURL(id string) (string, error) {
	if id == "" {
		return "", errors.New("missing artifact id")
	}
	if s.BaseURL == "" {
		return "", errors.New("missing base url")
	}
	return strings.TrimRight(s.BaseURL, "/") + "/artifacts/" + id, nil
}
