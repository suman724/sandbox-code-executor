package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"control-plane/internal/services"
)

type ServiceHandler struct {
	Starter func(service services.Service) error
}

type serviceRequest struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	PolicyID string `json:"policy_id"`
}

func (h ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req serviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("services: decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	svc := services.Service{ID: req.ID, TenantID: req.TenantID, PolicyID: req.PolicyID, Status: services.StatusStarting}
	if h.Starter != nil {
		if err := h.Starter(svc); err != nil {
			log.Printf("services: start error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusAccepted)
}
