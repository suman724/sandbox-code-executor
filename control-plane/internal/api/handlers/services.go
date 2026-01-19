package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"control-plane/internal/services"
)

type ServiceHandler struct {
	Starter  func(service services.Service) (string, error)
}

type serviceRequest struct {
	TenantID string `json:"tenantId"`
	PolicyID string `json:"policyId"`
}

type serviceResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	ProxyURL string `json:"proxyUrl,omitempty"`
}

func (h ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req serviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("services: decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.TenantID == "" || req.PolicyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	svc := services.Service{
		ID:       "service-" + time.Now().UTC().Format("20060102150405"),
		TenantID: req.TenantID,
		PolicyID: req.PolicyID,
		Status:   services.StatusStarting,
		StartedAt: time.Now().UTC(),
	}
	if h.Starter != nil {
		proxyURL, err := h.Starter(svc)
		if err != nil {
			log.Printf("services: start error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		svc.ProxyURL = proxyURL
		svc.Status = services.StatusRunning
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(serviceResponse{ID: svc.ID, Status: string(svc.Status), ProxyURL: svc.ProxyURL})
}
