package tools

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"control-plane/internal/storage"
)

type ArtifactsTool struct {
	Store storage.ArtifactStore
}

type uploadRequest struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	SizeBytes  int64  `json:"sizeBytes"`
	Checksum   string `json:"checksum"`
	StorageURI string `json:"storageUri"`
}

func (t ArtifactsTool) Upload(w http.ResponseWriter, r *http.Request) {
	if t.Store == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	var req uploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	artifact := storage.Artifact{
		ID:         req.ID,
		Name:       req.Name,
		SizeBytes:  req.SizeBytes,
		Checksum:   req.Checksum,
		StorageURI: req.StorageURI,
	}
	if err := t.Store.Put(r.Context(), artifact); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": req.ID})
}

func (t ArtifactsTool) Download(w http.ResponseWriter, r *http.Request) {
	if t.Store == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	artifactID := chi.URLParam(r, "artifactId")
	if artifactID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url, err := t.Store.SignedDownloadURL(artifactID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"url": url})
}
