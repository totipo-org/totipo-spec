package specpolicy

import (
	"os"
	"strings"
	"testing"
)

func inputs(t *testing.T) (string, string) {
	t.Helper()
	s, e := os.ReadFile("../../../spec/totipo-vault-format-v1.md")
	if e != nil {
		t.Fatal(e)
	}
	a, e := os.ReadFile("../../../review/V1_PRE_R16_REVISION_HISTORY.md")
	if e != nil {
		t.Fatal(e)
	}
	return string(s), string(a)
}
func TestCurrent(t *testing.T) {
	s, a := inputs(t)
	if e := Validate(s, a); e != nil {
		t.Fatal(e)
	}
}
func TestMutations(t *testing.T) {
	s, a := inputs(t)
	for _, m := range [][2]string{
		{"../review/V1_PRE_R16_REVISION_HISTORY.md", "W3B_REMOVED"},
		{"optional implementation hardening", "W3B_REMOVED"},
		{"does not require persistent retention", "W3B_REMOVED"},
		{"interactive application MUST require explicit confirmation", "W3B_REMOVED"},
		{"v1 has no in-place password change and defines no password-rewrap operation", "W3B_REMOVED"},
		{"INFORMATIVE", "W3B_REMOVED"},
		{"MUST NOT replace it, even if it is exact-identical", "W3B_REMOVED"},
		{"**Revision:** r19", "W3B_REMOVED"},
		{"**Protocol version:** 1", "W3B_REMOVED"},
		{"same exact presence", "W3B_REMOVED"},
		{"generate a new CSPRNG TOKEN_ID", "W3B_REMOVED"},
		{"empty parent set", "W3B_REMOVED"},
		{"TOKEN is the sole semantic object grammar", "W3B_REMOVED"},
		{"MAX_PARENTS = 4", "W3B_REMOVED"},
		{"0 <= PARENT_COUNT <= 4", "W3B_REMOVED"},
		{"1005", "W3B_REMOVED"},
		{"1006", "W3B_REMOVED"},
		{"1024", "W3B_REMOVED"},
		{"VAULT_ID = SHA-256(exact canonical VAULT representation)", "W3B_REMOVED"},
		{"A ≡c B iff A reaches B and B reaches A", "W3B_REMOVED"},
		{"Totipo v1 protocol state is create-only after vault initialization", "W3B_REMOVED"},
		{"Synchronization is\noptional and external", "W3B_REMOVED"},
		{"MUST NOT require exhaustive enumeration", "W3B_REMOVED"},
		{"Parent availability is not a publication prerequisite", "W3B_REMOVED"},
		{"MUST NOT be classified as garbage", "W3B_REMOVED"},
		{"Creation MUST NOT delete them", "W3B_REMOVED"},
		{"Reopen canonical `vault` as a regular file", "W3B_REMOVED"},
		{"NO-REPLACE", "W3B_REMOVED"},
		{"K_root = CSPRNG(32 bytes)", "W3B_REMOVED"},
		{"immutable for the complete lifetime", "W3B_REMOVED"},
		{"MUST NOT be rotated in-place", "W3B_REMOVED"},
		{"BOOTSTRAP_VERSION = 0x01", "W3B_REMOVED"},
		{"NEVER a newer version", "W3B_REMOVED"},
		{"same-root alternate legitimate bootstrap representations", "W3B_REMOVED"},
		{"per-object metadata rule in Section 12", "W3B_REMOVED"},
		{"Section 17 additionally governs operation-wide fold metadata", "W3B_REMOVED"},
		{"core library need not guarantee the behavior of a GUI", "W3B_REMOVED"},
		{"portable corpus success alone does not establish application conformance", "W3B_REMOVED"},
		{"No stronger atomic primitive is imposed by this scope", "W3B_REMOVED"},
		{"## 20.", "## 22."},
		{"**Revision:** r19", "**Revision:** r99"},
		{"Section 12", "Section 99"},
		{"### v1/r19\n", "### v1/r99\n"},
		{"../review/V1_PRE_R16_REVISION_HISTORY.md", "../review/missing.md"},
		{"0x0001", "0x00ff"},
		{"regular file only", "any entry type"},
		{"## 21. Revision history", "## 21. Missing history"},
		{"configured Totipo store", "unconfigured store"},
		{"    vault", "    VAULT"},
		{"same exact presence", "different presence"},
		{"TokenValue is the complete tuple", "TokenValue is some tuple"},
		{"MUST NOT erase metadata from head objects", "may erase metadata from head objects"},
		{"symlink", "link"},
		{"For as long as", "While"},
		{"Unavailability of history alone MUST NOT prohibit authorship", "History absence may prohibit authorship"},
		{"| Core protocol |", "| Wrong scope |"},
	} {
		t.Run(m[0], func(t *testing.T) {
			if !strings.Contains(s, m[0]) {
				t.Fatal("missing mutation anchor")
			}
			if e := Validate(strings.ReplaceAll(s, m[0], m[1]), a); e == nil {
				t.Fatal("mutation accepted")
			}
		})
	}
	if e := Validate(s, strings.ReplaceAll(a, "### v1/r1\n", "### v1/r99\n")); e == nil {
		t.Fatal("archive mutation accepted")
	}
}
func TestActiveHistorySeparation(t *testing.T) {
	s, a := inputs(t)
	for _, term := range []string{"DEVICE", "DEVICE_ID", "P-256", "ECDSA", "DER", "SIGNATURE", "AUTHOR_DEVICE_ID", "AUTHOR_TIME", "PROVENANCE", "OPAQUE_ROUTABLE", "OPAQUE_UNSCOPED", "OBJECT_TYPE", "OBJECT_VERSION", "VAULT_BINDING", "READY", "PROCESSING_INCOMPLETE", "HISTORY_REGRESSION", "HISTORY_MEMORY_LOST", "candidate-use", "K_signature_context", "openat2", "O_NOFOLLOW", "O_PATH", "statx", "SecureDirectoryStream", "VAULT_FINGERPRINT", "totipo/v1/vault-fingerprint", "CURRENT != BASE", "stageReplacement", "compare-before-replace", "latest password wrapper", "provenance", "opaque"} {
		t.Run(term, func(t *testing.T) {
			changed := strings.Replace(s, "## 21. Revision history", term+"\n## 21. Revision history", 1)
			if e := Validate(changed, a); e == nil {
				t.Fatal("active retired term accepted")
			}
			if e := Validate(s+"\n"+term, a+"\n"+term); e != nil {
				t.Fatal("historical term rejected", e)
			}
		})
	}
	// The old negative-lookaround rule permits an uppercase/underscore neighbor.
	for _, term := range []string{"XDEVICEX", "_DEVICE_", "XREADYX"} {
		if e := Validate(strings.Replace(s, "## 21. Revision history", term+"\n## 21. Revision history", 1), a); e != nil {
			t.Fatal(e)
		}
	}
}

func TestTextBoundaries(t *testing.T) {
	s, a := inputs(t)
	for _, changed := range []string{strings.ReplaceAll(s, "\n", "\r\n"), strings.Replace(s, "## 21. Revision history", "Section ١٢\n## 21. Revision history", 1), strings.ReplaceAll(s, "optional implementation hardening", "optional\u001cimplementation hardening")} {
		// The literal hardening anchor is deliberately not whitespace-flattened.
		want := !strings.Contains(changed, "optional\u001cimplementation")
		if e := Validate(changed, a); (e == nil) != want {
			t.Fatal(e)
		}
	}
	changed := strings.Replace(s, "configured Totipo store\n", "configured Totipo store\nvault\n\u00a0VAULT\u00a0\n", 1)
	if e := Validate(changed, a); e == nil {
		t.Fatal("Unicode whitespace uppercase layout accepted")
	}
}

func TestOriginalMutationSet(t *testing.T) {
	s, a := inputs(t)
	for _, m := range [][2]string{
		{"## 20.", "## 22."}, {"**Revision:** r19", "**Revision:** r99"}, {"Section 12", "Section 99"}, {"### v1/r19\n", "### v1/r99\n"}, {"../review/V1_PRE_R16_REVISION_HISTORY.md", "../review/missing.md"}, {"0x0001", "0x00ff"}, {"regular file only", "any entry type"},
	} {
		if e := Validate(strings.Replace(s, m[0], m[1], 1), a); e == nil {
			t.Fatal("original mutation accepted", m[0])
		}
	}
}
