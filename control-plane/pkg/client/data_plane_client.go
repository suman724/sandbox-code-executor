package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type RunRequest struct {
	JobID        string `json:"jobId"`
	PolicyID     string `json:"policyId,omitempty"`
	Language     string `json:"language"`
	Code         string `json:"code"`
	WorkspaceRef string `json:"workspaceRef"`
}

type RunResponse struct {
	RunID string
}

type DataPlaneClient struct {
	BaseURL string
	AuthToken string
	Client   *http.Client
}

func (c DataPlaneClient) StartRun(ctx context.Context, req RunRequest) (RunResponse, error) {
	if c.BaseURL == "" {
		return RunResponse{}, errors.New("missing base url")
	}
	client := c.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return RunResponse{}, err
	}
	url := strings.TrimRight(c.BaseURL, "/") + "/runs"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return RunResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.AuthToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return RunResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return RunResponse{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var decoded RunResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return RunResponse{}, err
	}
	return decoded, nil
}
