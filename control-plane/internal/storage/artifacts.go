package storage

import "context"

type ArtifactStore interface {
	Put(ctx context.Context, artifact Artifact) error
	Get(ctx context.Context, id string) (Artifact, error)
	SignedDownloadURL(id string) (string, error)
}
