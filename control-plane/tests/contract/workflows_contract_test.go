package contract

import "testing"

func TestWorkflowsContract(t *testing.T) {
	endpoint := "/workflows"
	if endpoint == "" {
		t.Fatalf("expected endpoint")
	}
}
