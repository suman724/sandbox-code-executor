package main

import (
	"log"
	"net/http"
	"os"

	"control-plane/internal/api"
	"control-plane/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	addr := ":" + getenv("PORT", "8080")
	log.Printf("control-plane starting env=%s addr=%s", cfg.Env, addr)

	if err := http.ListenAndServe(addr, api.Router()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
