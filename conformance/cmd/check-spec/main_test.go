package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"spec/totipo-vault-format-v1.md", "review/V1_PRE_R16_REVISION_HISTORY.md"} {
		b, e := os.ReadFile(filepath.Join("../../..", name))
		if e != nil {
			t.Fatal(e)
		}
		path := filepath.Join(root, name)
		if e := os.MkdirAll(filepath.Dir(path), 0755); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(path, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
	var out, errout bytes.Buffer
	if code := run([]string{"-root", root}, &out, &errout); code != 0 {
		t.Fatal(code, errout.String())
	}
	// This isolated tree has only the two policy inputs, no vectors/profile/schema.
	path := filepath.Join(root, "spec/totipo-vault-format-v1.md")
	b, _ := os.ReadFile(path)
	if e := os.WriteFile(path, []byte(strings.Replace(string(b), "**Revision:** r19", "**Revision:** r99", 1)), 0600); e != nil {
		t.Fatal(e)
	}
	if code := run([]string{"-root", root}, &out, &errout); code != 1 || !strings.Contains(errout.String(), "spec policy") {
		t.Fatal(code, errout.String())
	}
	for _, args := range [][]string{{"extra"}, {"-unknown"}, {"-root"}} {
		if code := run(args, &out, &errout); code != 2 {
			t.Fatal(args, code)
		}
	}
	if code := run([]string{"-root", t.TempDir()}, &out, &errout); code != 1 {
		t.Fatal(code)
	}
}
