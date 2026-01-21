package unit

import (
	"testing"

	"data-plane/internal/runtime"
)

func TestSessionRegistryStoresRoute(t *testing.T) {
	registry := runtime.NewInMemorySessionRegistry()
	route := runtime.SessionRoute{
		RuntimeID: "runtime-1",
		Endpoint:  "http://localhost:9000",
		Token:     "token",
		AuthMode:  "bypass",
	}
	if err := registry.Put("session-1", route); err != nil {
		t.Fatalf("put route: %v", err)
	}
	stored, ok := registry.Get("session-1")
	if !ok {
		t.Fatalf("expected route to exist")
	}
	if stored.RuntimeID != route.RuntimeID {
		t.Fatalf("expected runtime id %q, got %q", route.RuntimeID, stored.RuntimeID)
	}
	if stored.Endpoint != route.Endpoint {
		t.Fatalf("expected endpoint %q, got %q", route.Endpoint, stored.Endpoint)
	}
	if stored.AuthMode != route.AuthMode {
		t.Fatalf("expected auth mode %q, got %q", route.AuthMode, stored.AuthMode)
	}
	registry.Delete("session-1")
	if _, ok := registry.Get("session-1"); ok {
		t.Fatalf("expected route to be deleted")
	}
}
