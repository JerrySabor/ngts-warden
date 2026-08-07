package catalog

import (
	"net/http"
	"testing"
)

func TestOfficialOperationInventory(t *testing.T) {
	spec, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	ops := spec.Operations()
	if len(spec.Paths) != 95 {
		t.Fatalf("paths = %d, want 95", len(spec.Paths))
	}
	if len(ops) != 147 {
		t.Fatalf("operations = %d, want 147", len(ops))
	}
	groups := map[string]bool{}
	methods := map[string]int{}
	seen := map[string]bool{}
	for _, op := range ops {
		groups[op.Group] = true
		methods[op.Method]++
		key := op.Group + "/" + op.Action
		if seen[key] {
			t.Fatalf("duplicate generated command %s", key)
		}
		seen[key] = true
	}
	if len(groups) != 24 {
		t.Fatalf("groups = %d, want 24", len(groups))
	}
	for method, want := range map[string]int{http.MethodGet: 55, http.MethodPost: 52, http.MethodDelete: 18, http.MethodPut: 12, http.MethodPatch: 10} {
		if methods[method] != want {
			t.Errorf("%s = %d, want %d", method, methods[method], want)
		}
	}
}
