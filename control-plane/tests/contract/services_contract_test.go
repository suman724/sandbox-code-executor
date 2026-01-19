package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"control-plane/internal/api/handlers"
	"control-plane/internal/services"
)

func TestServicesContract(t *testing.T) {
	handler := handlers.ServiceHandler{
		Starter: func(service services.Service) (string, error) {
			return "http://proxy/" + service.ID, nil
		},
	}

	payload := map[string]any{
		"tenantId": "tenant-1",
		"policyId": "policy-1",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, rec.Code)
	}
}
