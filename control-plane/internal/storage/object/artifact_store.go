package object

import (
	"context"

	"control-plane/internal/storage"
)

type ArtifactStore struct{}

func (ArtifactStore) Put(ctx context.Context, artifact storage.Artifact) error {
	_ = ctx
	_ = artifact
	return nil
}

func (ArtifactStore) Get(ctx context.Context, id string) (storage.Artifact, error) {
	_ = ctx
	_ = id
	return storage.Artifact{}, nil
}
