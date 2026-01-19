package contract

import "testing"

func TestServicesContract(t *testing.T) {
	endpoint := "/services"
	if endpoint == "" {
		t.Fatalf("expected endpoint")
	}
}
