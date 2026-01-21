package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"shared/sessionagent"
)

type AgentClient struct {
	HTTPClient *http.Client
}

type AgentRoute struct {
	Endpoint string
	Token    string
	AuthMode string
}

func NewAgentClient() *AgentClient {
	return &AgentClient{HTTPClient: &http.Client{Timeout: 10 * time.Second}}
}

func (c *AgentClient) RunStep(ctx context.Context, route AgentRoute, request sessionagent.StepRequest) (sessionagent.StepResult, error) {
	if route.Endpoint == "" {
		return sessionagent.StepResult{}, errors.New("agent endpoint not configured")
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return sessionagent.StepResult{}, fmt.Errorf("encode step request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, route.Endpoint+"/v1/steps", bytes.NewReader(payload))
	if err != nil {
		return sessionagent.StepResult{}, fmt.Errorf("create agent request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if route.AuthMode != "bypass" && route.Token != "" {
		req.Header.Set("X-Session-Token", route.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return sessionagent.StepResult{}, fmt.Errorf("agent request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return sessionagent.StepResult{}, fmt.Errorf("agent returned status %d", resp.StatusCode)
	}

	var result sessionagent.StepResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return sessionagent.StepResult{}, fmt.Errorf("decode agent response: %w", err)
	}

	return result, nil
}
