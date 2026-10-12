package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"totipo/conformance/internal/vectors"
)

func tree(t *testing.T) string {
	t.Helper()
	r := t.TempDir()
	for _, dir := range []string{"spec", "vectors", "requirements"} {
		if e := os.CopyFS(filepath.Join(r, dir), os.DirFS("../../../"+dir)); e != nil {
			t.Fatal(e)
		}
	}
	return r
}
func read(t *testing.T, r, name string) map[string]any {
	t.Helper()
	b, e := os.ReadFile(filepath.Join(r, name))
	if e != nil {
		t.Fatal(e)
	}
	var v map[string]any
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	if e := d.Decode(&v); e != nil {
		t.Fatal(e)
	}
	return v
}
func write(t *testing.T, r, name string, v any) []byte {
	t.Helper()
	b, e := json.Marshal(v)
	if e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(filepath.Join(r, name), b, 0600); e != nil {
		t.Fatal(e)
	}
	return b
}
func call(r string, verify bool) (int, string) {
	args := []string{"-root", r}
	if verify {
		args = append(args, "-verify-only")
	}
	var out, errout bytes.Buffer
	code := run(args, &out, &errout)
	return code, out.String() + errout.String()
}

// Repin deliberate temp-tree mutations so a schema rejection cannot be masked
// by an unrelated checksum/profile failure.
func changeCase(t *testing.T, r, path string, change func(map[string]any)) {
	c := read(t, r, "vectors/"+path)
	change(c)
	b := write(t, r, "vectors/"+path, c)
	m := read(t, r, "vectors/manifest.json")
	id := c["id"].(string)
	for _, x := range m["cases"].([]any) {
		e := x.(map[string]any)
		if e["id"] == id {
			e["sha256"] = vectors.Hash(b)
		}
	}
	mb := write(t, r, "vectors/manifest.json", m)
	p := read(t, r, "requirements/v1-pre-rc.json")
	p["manifest_sha256"] = vectors.Hash(mb)
	for _, x := range p["required_cases"].([]any) {
		e := x.(map[string]any)
		if e["id"] == id {
			e["sha256"] = vectors.Hash(b)
		}
	}
	write(t, r, "requirements/v1-pre-rc.json", p)
}
func TestCLIValidAndSchemaRefusal(t *testing.T) {
	r := tree(t)
	for _, verify := range []bool{true, false} {
		if code, out := call(r, verify); code != 0 || !strings.Contains(out, "92") {
			t.Fatal(code, out)
		}
	}
	for _, kind := range []string{"manifest", "case", "schema-source"} {
		t.Run(kind, func(t *testing.T) {
			r := tree(t)
			switch kind {
			case "manifest":
				m := read(t, r, "vectors/manifest.json")
				delete(m, "cases")
				write(t, r, "vectors/manifest.json", m)
			case "case":
				changeCase(t, r, "cases/totp/v1.totp.rfc6238-sha1.001.json", func(c map[string]any) { delete(c["totp"].(map[string]any), "t0") })
			case "schema-source":
				s := read(t, r, "vectors/case.schema.json")
				s["unsupported"] = true
				write(t, r, "vectors/case.schema.json", s)
			}
			for _, verify := range []bool{true, false} {
				if code, out := call(r, verify); code != 1 || !strings.Contains(out, "schema") {
					t.Fatal(code, out)
				}
			}
		})
	}
}
func TestCLIVerifySkipsSemanticExecution(t *testing.T) {
	r := tree(t)
	changeCase(t, r, "cases/totp/v1.totp.rfc6238-sha1.001.json", func(c map[string]any) {
		c["totp"].(map[string]any)["rows"].([]any)[0].(map[string]any)["code"] = "00000000"
	})
	if code, out := call(r, true); code != 0 {
		t.Fatal(code, out)
	}
	if code, out := call(r, false); code != 1 || !strings.Contains(out, "TOTP mismatch") {
		t.Fatal(code, out)
	}
}
func TestCLIArgumentsAndLoadRefusal(t *testing.T) {
	for _, args := range [][]string{{"extra"}, {"-unknown"}, {"-root"}} {
		var out bytes.Buffer
		if code := run(args, &out, &out); code != 2 {
			t.Fatal(args, code, out.String())
		}
	}
	if code, out := call(t.TempDir(), true); code != 1 {
		t.Fatal(code, out)
	}
	r := tree(t)
	p := read(t, r, "requirements/v1-pre-rc.json")
	p["spec_revision"] = "r99"
	write(t, r, "requirements/v1-pre-rc.json", p)
	if code, out := call(r, true); code != 1 || !strings.Contains(out, "profile") {
		t.Fatal(code, out)
	}
}

// Snapshot paths as well as bytes: additions, removals and rewrites must fail.
// The disposable tree contains only committed artifacts, with no generator,
// module source, or executable available beneath the supplied root.
func artifactSnapshot(t *testing.T, root string) map[string][32]byte {
	t.Helper()
	result := map[string][32]byte{}
	for _, dir := range []string{"spec", "vectors", "requirements"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			result[rel] = sha256.Sum256(b)
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return result
}

func TestCLIVerificationLeavesArtifactsUnchanged(t *testing.T) {
	for _, invalid := range []bool{false, true} {
		t.Run(map[bool]string{false: "valid", true: "invalid"}[invalid], func(t *testing.T) {
			r := tree(t)
			wantCode := 0
			if invalid {
				// Failure must leave the bad serialized input untouched as well.
				if err := os.WriteFile(filepath.Join(r, "vectors/manifest.json"), []byte("{"), 0600); err != nil {
					t.Fatal(err)
				}
				wantCode = 1
			}
			before := artifactSnapshot(t, r)
			if len(before) != 100 {
				t.Fatalf("artifact inventory: %d", len(before))
			}
			for _, verify := range []bool{true, false} {
				if code, out := call(r, verify); code != wantCode {
					t.Fatal(code, out)
				}
				if !reflect.DeepEqual(before, artifactSnapshot(t, r)) {
					t.Fatal("verification changed canonical artifact paths or bytes")
				}
			}
		})
	}
}
