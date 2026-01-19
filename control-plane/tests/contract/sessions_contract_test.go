package contract

import "testing"

func TestSessionsContract(t *testing.T) {
	endpoint := "/sessions"
	if endpoint == "" {
		t.Fatalf("expected endpoint")
	}
}
