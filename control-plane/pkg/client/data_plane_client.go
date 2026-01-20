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
	RunID string `json:"run_id"`
}

type SessionCreateRequest struct {
	SessionID    string `json:"sessionId"`
	PolicyID     string `json:"policyId,omitempty"`
	WorkspaceRef string `json:"workspaceRef"`
}

type SessionResponse struct {
	ID        string `json:"id"`
	RuntimeID string `json:"runtimeId"`
	Status    string `json:"status"`
}

type SessionStepRequest struct {
	Command string `json:"command"`
}

type SessionStepResponse struct {
	Status string `json:"status"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type DataPlaneClient struct {
	BaseURL   string
	AuthToken string
	Client    *http.Client
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

func (c DataPlaneClient) StartSession(ctx context.Context, req SessionCreateRequest) (SessionResponse, error) {
	if c.BaseURL == "" {
		return SessionResponse{}, errors.New("missing base url")
	}
	client := c.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return SessionResponse{}, err
	}
	url := strings.TrimRight(c.BaseURL, "/") + "/sessions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return SessionResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.AuthToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return SessionResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return SessionResponse{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var decoded SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return SessionResponse{}, err
	}
	return decoded, nil
}

func (c DataPlaneClient) RunSessionStep(ctx context.Context, sessionID string, command string) (SessionStepResponse, error) {
	if c.BaseURL == "" {
		return SessionStepResponse{}, errors.New("missing base url")
	}
	if sessionID == "" {
		return SessionStepResponse{}, errors.New("missing session id")
	}
	if command == "" {
		return SessionStepResponse{}, errors.New("missing command")
	}
	client := c.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	body, err := json.Marshal(SessionStepRequest{Command: command})
	if err != nil {
		return SessionStepResponse{}, err
	}
	url := strings.TrimRight(c.BaseURL, "/") + "/sessions/" + sessionID + "/steps"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return SessionStepResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.AuthToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return SessionStepResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return SessionStepResponse{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var decoded SessionStepResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return SessionStepResponse{}, err
	}
	return decoded, nil
}
