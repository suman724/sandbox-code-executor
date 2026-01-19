package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"control-plane/internal/policy"
)

type PolicyHandler struct {
	Store policy.Store
}

type policyRequest struct {
	TenantID string `json:"tenantId"`
	Name     string `json:"name"`
	Version  int    `json:"version"`
	Ruleset  string `json:"ruleset"`
}

func (h PolicyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req policyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("policies: decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.TenantID == "" || req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	policyID := req.TenantID + ":" + req.Name
	if err := h.Store.Upsert(r.Context(), policy.Policy{ID: policyID, Version: req.Version, Ruleset: req.Ruleset}); err != nil {
		log.Printf("policies: upsert error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
