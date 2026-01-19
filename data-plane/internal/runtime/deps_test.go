package runtime

import "testing"

func TestDependencyAllowlist(t *testing.T) {
	policy := DependencyPolicy{Allowlist: []string{"a"}, Requested: []string{"a"}}
	if err := ValidateDependencies(policy); err != nil {
		t.Fatalf("expected allowlisted dependency")
	}
	policy = DependencyPolicy{Allowlist: []string{"a"}, Requested: []string{"b"}}
	if err := ValidateDependencies(policy); err == nil {
		t.Fatalf("expected non-allowlisted dependency to fail")
	}
	policy = DependencyPolicy{Allowlist: []string{}, Requested: []string{"a"}}
	if err := ValidateDependencies(policy); err == nil {
		t.Fatalf("expected empty allowlist to fail")
	}
}
