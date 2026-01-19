package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Artifact struct {
	Name     string
	Path     string
	Size     int64
	Checksum string
}

func CaptureArtifact(path string) (Artifact, error) {
	if path == "" {
		return Artifact{}, errors.New("missing path")
	}
	info, err := os.Stat(path)
	if err != nil {
		return Artifact{}, err
	}
	if info.IsDir() {
		return Artifact{}, errors.New("artifact path is directory")
	}
	checksum, err := fileChecksum(path)
	if err != nil {
		return Artifact{}, err
	}
	return Artifact{
		Name:     filepath.Base(path),
		Path:     path,
		Size:     info.Size(),
		Checksum: checksum,
	}, nil
}

func fileChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
