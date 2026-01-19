package main

import (
	"log"
	"net/http"
	"os"

	"data-plane/internal/config"
	"data-plane/internal/runtime"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	addr := ":" + getenv("PORT", "8081")
	log.Printf("data-plane starting env=%s addr=%s", cfg.Env, addr)

	if err := http.ListenAndServe(addr, runtime.Router()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
