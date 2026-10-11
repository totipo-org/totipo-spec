package vectors

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"syscall"
	"testing"
)

// These boundary tests run even when the host cannot create symlinks.
func TestContainedPath(t *testing.T) {
	base := filepath.Join(t.TempDir(), "vectors")
	for _, tt := range []struct {
		name, base, candidate string
		want                  bool
	}{
		{"child", base, filepath.Join(base, "cases", "case.json"), true},
		{"same-root", base, base, true}, // Preserve the existing Rel/IsLocal policy.
		{"parent", base, filepath.Dir(base), false},
		{"sibling", base, filepath.Join(filepath.Dir(base), "outside", "case.json"), false},
		{"shared-prefix", base, base + "-outside" + string(filepath.Separator) + "case.json", false},
		{"relative-base", "vectors", filepath.Join(base, "case.json"), false},
		{"relative-candidate", base, "cases/case.json", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := containedPath(tt.base, tt.candidate); got != tt.want {
				t.Fatalf("containedPath(%q, %q) = %v; want %v", tt.base, tt.candidate, got, tt.want)
			}
		})
	}
	if runtime.GOOS == "windows" {
		for _, tt := range []struct {
			base, candidate string
			want            bool
		}{
			{`C:\vectors`, `D:\vectors\case.json`, false},
			{`C:\vectors`, `c:\VECTORS\cases\case.json`, true},
			{`\\server\share\vectors`, `\\server\other\vectors\case.json`, false},
			{`\\server\share\vectors`, `\\server\share\vectors\cases\case.json`, true},
		} {
			if got := containedPath(tt.base, tt.candidate); got != tt.want {
				t.Errorf("containedPath(%q, %q) = %v; want %v", tt.base, tt.candidate, got, tt.want)
			}
		}
	}
}

func pathTestSymlink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		// Windows may require Developer Mode or the symlink privilege. Only fixture
		// creation is skipped; TestContainedPath still checks the rejection boundary.
		if runtime.GOOS == "windows" && (errors.Is(err, syscall.Errno(1314)) || errors.Is(err, os.ErrPermission)) {
			t.Skipf("symlink fixture unavailable without Windows symlink privilege: %v", err)
		}
		t.Fatal(err)
	}
}

func TestReadRootAlias(t *testing.T) {
	root := copyIntegrityRoot(t)
	wantManifest, wantCases, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(wantCases) != 92 {
		t.Fatalf("normal copied corpus: %d cases", len(wantCases))
	}
	alias := filepath.Join(t.TempDir(), "alias-root")
	pathTestSymlink(t, root, alias)
	gotManifest, gotCases, err := Read(alias)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotManifest, wantManifest) || !reflect.DeepEqual(gotCases, wantCases) {
		t.Fatal("alias changed corpus identity, hashes, paths or outcomes")
	}
	if err := VerifyProfile(alias, gotManifest); err != nil {
		t.Fatal(err)
	}
}

func TestReadCaseSymlinkEscape(t *testing.T) {
	for _, directory := range []bool{false, true} {
		name := "file"
		if directory {
			name = "directory"
		}
		t.Run(name, func(t *testing.T) {
			root := copyIntegrityRoot(t)
			m, _, err := Read(root)
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(root, "vectors", m.Cases[0].Path)
			outside := filepath.Join(t.TempDir(), "outside")
			if directory {
				path = filepath.Dir(path)
			}
			if err := os.Rename(path, outside); err != nil {
				t.Fatal(err)
			}
			pathTestSymlink(t, outside, path)
			if _, _, err := Read(root); err == nil || err.Error() != "case symlink escapes vectors" {
				t.Fatalf("escape: %v", err)
			}
		})
	}
}

func TestResolvedPathMissing(t *testing.T) {
	if _, err := resolvedPath(filepath.Join(t.TempDir(), "missing", "vectors")); err == nil {
		t.Fatal("unresolvable root accepted")
	}
}
