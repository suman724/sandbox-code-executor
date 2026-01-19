package runtime

import "testing"

type noopAdapter struct{}

func (noopAdapter) Run(code string) error {
	_ = code
	return nil
}

func TestRegistryLookup(t *testing.T) {
	reg := NewRegistry()
	reg.Register("python", noopAdapter{})
	if !reg.Supports("python") {
		t.Fatalf("expected registry to support python")
	}
	if _, ok := reg.Adapter("python"); !ok {
		t.Fatalf("expected adapter")
	}
}
