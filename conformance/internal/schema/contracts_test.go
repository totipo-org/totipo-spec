package schema

import (
	"encoding/json"
	"os"
	"testing"
)

func caseValue(t *testing.T, path string) map[string]any {
	t.Helper()
	b, e := os.ReadFile("../../../vectors/cases/" + path + ".json")
	if e != nil {
		t.Fatal(e)
	}
	return parsed(t, b).(map[string]any)
}
func TestMigratedPythonContracts(t *testing.T) {
	s := contract(t, "case")
	for _, trial := range []struct {
		name, path string
		change     func(map[string]any)
	}{
		{"required false", "graph/v1.graph.sequential.001", func(c map[string]any) {
			steps := c["graph"].(map[string]any)["steps"].([]any)
			delete(steps[len(steps)-1].(map[string]any)["expect"].(map[string]any), "conflicting")
		}},
		{"head metadata type", "graph/v1.graph.equal-concurrent.001", func(c map[string]any) {
			steps := c["graph"].(map[string]any)["steps"].([]any)
			steps[len(steps)-1].(map[string]any)["expect"].(map[string]any)["head_objects"].([]any)[0].(map[string]any)["client_time"] = json.Number("18446744073709551616")
		}},
		{"required heads", "graph/v1.graph.equal-concurrent.001", func(c map[string]any) {
			steps := c["graph"].(map[string]any)["steps"].([]any)
			delete(steps[len(steps)-1].(map[string]any)["expect"].(map[string]any), "head_objects")
		}},
		{"nested unknown", "graph/v1.graph.sequential.001", func(c map[string]any) {
			steps := c["graph"].(map[string]any)["steps"].([]any)
			steps[len(steps)-1].(map[string]any)["expect"].(map[string]any)["typo"] = false
		}},
		{"mixed payload", "graph/v1.graph.sequential.001", func(c map[string]any) {
			c["root_hex"] = "0000000000000000000000000000000000000000000000000000000000000000"
		}},
		{"post aead defect", "crypto/v1.crypto.nonzero-padding.001", func(c map[string]any) { c["post_aead"].(map[string]any)["defect"] = "failed-aead" }},
		{"post aead width", "crypto/v1.crypto.nonzero-padding.001", func(c map[string]any) {
			p := c["post_aead"].(map[string]any)
			p["object_hex"] = p["object_hex"].(string)[2:]
		}},
		{"replacement", "vault/v1.vault.create-existing.001", func(c map[string]any) { c["workflow"].(map[string]any)["action"] = "replace" }},
		{"rewrap", "vault/v1.vault.create-existing.001", func(c map[string]any) { c["workflow"].(map[string]any)["action"] = "rewrap" }},
		{"base extension", "vault/v1.vault.create-existing.001", func(c map[string]any) { c["workflow"].(map[string]any)["base_hex"] = "01" }},
	} {
		t.Run(trial.name, func(t *testing.T) {
			c := caseValue(t, trial.path)
			if e := s.Validate(c); e != nil {
				t.Fatal(e)
			}
			trial.change(c)
			if e := s.Validate(c); e == nil {
				t.Fatal("mutation accepted")
			}
		})
	}
	for _, bad := range []any{json.Number("-1"), json.Number("18446744073709551616"), json.Number("1.5"), true} {
		c := caseValue(t, "metadata/v1.metadata.client-time-u64max.001")
		c["input"].(map[string]any)["client_time"] = bad
		if e := s.Validate(c); e == nil {
			t.Fatal("unsigned timestamp accepted", bad)
		}
	}
	for _, field := range []string{"readable", "complete", "durable", "orphan_objects", "parents_available"} {
		c := caseValue(t, "vault/v1.vault.create-existing.001")
		w := c["workflow"].(map[string]any)
		w[field] = false
		if e := s.Validate(c); e != nil {
			t.Fatal(e)
		}
		delete(w, field)
		if e := s.Validate(c); e == nil {
			t.Fatal("absent boolean accepted", field)
		}
	}
	for _, field := range []string{"t0", "rows", "period"} {
		c := caseValue(t, "totp/v1.totp.rfc6238-sha1.001")
		delete(c["totp"].(map[string]any), field)
		if e := s.Validate(c); e == nil {
			t.Fatal("absent field accepted", field)
		}
	}
	ms := contract(t, "manifest")
	b, e := os.ReadFile("../../../vectors/manifest.json")
	if e != nil {
		t.Fatal(e)
	}
	m := parsed(t, b).(map[string]any)
	entry := m["cases"].([]any)[0].(map[string]any)
	entry["applicability"] = map[string]any{"kind": "baseline"}
	if e := ms.Validate(m); e == nil {
		t.Fatal("retired extension accepted")
	}
	delete(entry, "applicability")
	m["spec_revision"] = "r99"
	if e := ms.Validate(m); e == nil {
		t.Fatal("revision accepted")
	}
}
