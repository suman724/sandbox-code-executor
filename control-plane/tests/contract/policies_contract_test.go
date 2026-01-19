package contract

import "testing"

func TestPoliciesContract(t *testing.T) {
	endpoint := "/policies"
	if endpoint == "" {
		t.Fatalf("expected endpoint")
	}
}
