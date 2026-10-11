package vectors

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func copyIntegrityRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{"spec", "vectors", "requirements"} {
		if err := os.CopyFS(filepath.Join(root, dir), os.DirFS("../../../"+dir)); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func writeIntegrityFile(t *testing.T, root, name string, b []byte) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), b, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestArtifactPins(t *testing.T) {
	root := copyIntegrityRoot(t)
	m, _, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyProfile(root, m); err != nil {
		t.Fatal(err)
	}
	profilePath := "requirements/v1-pre-rc.json"
	original, err := os.ReadFile(filepath.Join(root, profilePath))
	if err != nil {
		t.Fatal(err)
	}
	for field, artifact := range map[string]string{
		"spec_sha256":        "spec/totipo-vault-format-v1.md",
		"manifest_sha256":    "vectors/manifest.json",
		"schema_sha256":      "vectors/manifest.schema.json",
		"case_schema_sha256": "vectors/case.schema.json",
	} {
		t.Run(field, func(t *testing.T) {
			b, err := os.ReadFile(filepath.Join(root, artifact))
			if err != nil {
				t.Fatal(err)
			}
			writeIntegrityFile(t, root, artifact, append(b, '\n'))
			if err := VerifyProfile(root, m); err == nil || !strings.Contains(err.Error(), artifact) {
				t.Fatalf("byte drift: %v", err)
			}
			writeIntegrityFile(t, root, artifact, b)
			for _, missing := range []bool{false, true} {
				var p map[string]any
				if err := json.Unmarshal(original, &p); err != nil {
					t.Fatal(err)
				}
				p[field] = strings.Repeat("0", 64)
				if missing {
					delete(p, field)
				}
				changed, err := json.Marshal(p)
				if err != nil {
					t.Fatal(err)
				}
				writeIntegrityFile(t, root, profilePath, changed)
				if err := VerifyProfile(root, m); err == nil || !strings.Contains(err.Error(), artifact) {
					t.Fatalf("missing=%v: %v", missing, err)
				}
			}
			writeIntegrityFile(t, root, profilePath, original)
		})
	}
}

func TestProfileRelationship(t *testing.T) {
	root := copyIntegrityRoot(t)
	m, _, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, "requirements/v1-pre-rc.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, mutation := range []string{"revision", "missing", "extra", "duplicate", "unknown", "hash"} {
		t.Run(mutation, func(t *testing.T) {
			var p Profile
			if err := Decode(b, &p); err != nil {
				t.Fatal(err)
			}
			switch mutation {
			case "revision":
				p.Revision = "r99"
			case "missing":
				p.Required = p.Required[1:]
			case "extra":
				p.Required = append(p.Required, p.Required[0])
			case "duplicate":
				p.Required[1] = p.Required[0]
			case "unknown":
				p.Required[0].ID = "v1.unknown.001"
			case "hash":
				p.Required[0].SHA256 = strings.Repeat("0", 64)
			}
			changed, err := json.Marshal(p)
			if err != nil {
				t.Fatal(err)
			}
			writeIntegrityFile(t, root, "requirements/v1-pre-rc.json", changed)
			if err := VerifyProfile(root, m); err == nil {
				t.Fatal("accepted", mutation)
			}
		})
	}
}

func TestManifestContract(t *testing.T) {
	root := copyIntegrityRoot(t)
	m, _, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, mutation := range []string{"unsafe", "backslash", "duplicate-id", "duplicate-path", "identity", "shape", "case-bytes", "unlisted"} {
		t.Run(mutation, func(t *testing.T) {
			changed := m
			changed.Cases = append([]Entry(nil), m.Cases...)
			entry := &changed.Cases[0]
			casePath := entry.Path
			original, err := os.ReadFile(filepath.Join(root, "vectors", casePath))
			if err != nil {
				t.Fatal(err)
			}
			defer writeIntegrityFile(t, root, "vectors/"+casePath, original)
			switch mutation {
			case "unsafe":
				entry.Path = "../escape.json"
			case "backslash":
				entry.Path = `cases\escape.json`
			case "duplicate-id":
				changed.Cases[1].ID = entry.ID
			case "duplicate-path":
				changed.Cases[1].Path = entry.Path
			case "case-bytes":
				writeIntegrityFile(t, root, "vectors/"+casePath, append(original, '\n'))
			case "unlisted":
				writeIntegrityFile(t, root, "vectors/cases/unlisted.json", []byte("{}"))
				defer os.Remove(filepath.Join(root, "vectors/cases/unlisted.json"))
			default:
				var c Case
				if err := Decode(original, &c); err != nil {
					t.Fatal(err)
				}
				if mutation == "identity" {
					c.ID = "v1.other.001"
				} else {
					c.Operation = "unknown"
				}
				b, err := json.Marshal(c)
				if err != nil {
					t.Fatal(err)
				}
				entry.SHA256 = Hash(b) // Isolate identity/shape rejection past the checksum guard.
				writeIntegrityFile(t, root, "vectors/"+casePath, b)
			}
			b, err := json.Marshal(changed)
			if err != nil {
				t.Fatal(err)
			}
			writeIntegrityFile(t, root, "vectors/manifest.json", b)
			if _, _, err := Read(root); err == nil {
				t.Fatal("accepted", mutation)
			}
		})
	}
}
