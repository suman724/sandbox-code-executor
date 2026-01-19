package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"control-plane/internal/api/handlers"
	"control-plane/internal/policy"
)

func TestPoliciesContract(t *testing.T) {
	store := policy.NewInMemoryStore()
	handler := handlers.PolicyHandler{Store: store}

	payload := map[string]any{
		"tenantId": "tenant-1",
		"name":     "default",
		"version":  1,
		"ruleset":  "package policy\n default allow = true",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/policies", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}
}
