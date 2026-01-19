package contract

import "testing"

func TestJobsContract(t *testing.T) {
	endpoint := "/jobs"
	if endpoint == "" {
		t.Fatalf("expected endpoint")
	}
}
