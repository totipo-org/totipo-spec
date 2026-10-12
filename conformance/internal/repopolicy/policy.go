// Package repopolicy enforces implementation ownership, not protocol semantics.
package repopolicy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var roles = []string{"generator", "conformance"}
var primitives = map[string]bool{
	"golang.org/x/crypto/argon2":  true,
	"golang.org/x/crypto/blake2b": true,
	"golang.org/x/crypto/hkdf":    true,
	"golang.org/x/sys/cpu":        true,
}

func within(base, path string) bool {
	r, e := filepath.Rel(base, path)
	return e == nil && filepath.IsLocal(r)
}
func roleOf(root, path string) string {
	for _, role := range roles {
		if within(filepath.Join(root, role), path) {
			return role
		}
	}
	return ""
}
func ignoredDir(name string) bool {
	return name == ".git" || name == ".direnv" || name == "vendor" || name == "__pycache__"
}

// Static checks include tests, generated files and inactive build-tagged files.
// All source symlinks are rejected, including directory links; no source traversal
// can silently skip an implementation. Vendor code is checked in the loaded graph.
func Source(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && filepath.Dir(path) == root && ignoredDir(d.Name()) {
			return filepath.SkipDir
		}
		role := roleOf(root, path)
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("source symlink requires review: %s", path)
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "go.mod" && path != filepath.Join(root, "generator/go.mod") && path != filepath.Join(root, "conformance/go.mod") {
			return fmt.Errorf("third/root module: %s", path)
		}
		ext := filepath.Ext(path)
		if role != "" && (ext == ".s" || ext == ".S" || ext == ".c" || ext == ".h" || ext == ".cc" || ext == ".syso") {
			return fmt.Errorf("native source bridge requires review: %s", path)
		}
		if ext != ".go" {
			return nil
		}
		if role == "" {
			return fmt.Errorf("Go source outside approved roots: %s", path)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		f, err := parser.ParseFile(token.NewFileSet(), path, b, parser.ParseComments)
		if err != nil {
			return err
		}
		for _, group := range f.Comments {
			for _, c := range group.List {
				if strings.Contains(c.Text, "go:linkname") || strings.Contains(c.Text, "go:generate") {
					return fmt.Errorf("source directive requires review: %s", path)
				}
			}
		}
		for _, imp := range f.Imports {
			name, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				return err
			}
			specCommand := within(filepath.Join(root, "conformance/cmd/check-spec"), path)
			specPolicy := within(filepath.Join(root, "conformance/internal/specpolicy"), path)
			if (specCommand || specPolicy) && (strings.HasPrefix(name, "totipo/") || strings.Contains(strings.Split(name, "/")[0], ".")) && !(specCommand && name == "totipo/conformance/internal/specpolicy") {
				return fmt.Errorf("spec checker dependency requires review: %s in %s", name, path)
			}
			if within(filepath.Join(root, "conformance/internal/wireoracle"), path) && (strings.HasPrefix(name, "totipo/") || strings.Contains(strings.Split(name, "/")[0], ".")) {
				return fmt.Errorf("wire oracle must use only stdlib: %s in %s", name, path)
			}
			own := "totipo/" + role
			if strings.HasPrefix(name, "totipo/") && name != own && !strings.HasPrefix(name, own+"/") {
				return fmt.Errorf("cross-role/repository import %s in %s", name, path)
			}
			if name == "C" || name == "plugin" || name == "unsafe" {
				return fmt.Errorf("bridge import %s in %s", name, path)
			}
			if name == "os/exec" && path != filepath.Join(root, "conformance/internal/repopolicy/policy.go") {
				return fmt.Errorf("process delegation requires review: %s", path)
			}
			if strings.Contains(strings.Split(name, "/")[0], ".") && !(role == "conformance" && primitives[name]) {
				return fmt.Errorf("unreviewed external source import %s", name)
			}
		}
		return nil
	})
}

func goCommand(root, dir string, args ...string) ([]byte, error) {
	c := exec.Command("go", args...)
	c.Dir = dir
	// Preserve GOFLAGS so the pinned workspace vendor mode remains effective.
	c.Env = os.Environ()
	for _, kv := range []string{"GOWORK=" + filepath.Join(root, "go.work"), "GOPROXY=off", "GOSUMDB=off", "GOTOOLCHAIN=local", "GOENV=off"} {
		key := strings.SplitN(kv, "=", 2)[0] + "="
		filtered := c.Env[:0]
		for _, old := range c.Env {
			if !strings.HasPrefix(old, key) {
				filtered = append(filtered, old)
			}
		}
		c.Env = append(filtered, kv)
	}
	var stderr bytes.Buffer
	c.Stderr = &stderr
	b, err := c.Output()
	if err != nil {
		return nil, fmt.Errorf("go %v: %w: %s", args, err, stderr.String())
	}
	return b, nil
}

type module struct {
	Path, Version, Dir, GoVersion string
	Main                          bool
	Replace                       *module
}
type modFile struct {
	Ignore  []json.RawMessage
	Module  *module
	Go      string
	Require []struct {
		Path, Version string
		Indirect      bool
	}
	Replace   []json.RawMessage
	Exclude   []json.RawMessage
	Retract   []json.RawMessage
	Toolchain string
	Tool      []json.RawMessage
	Godebug   []json.RawMessage
}

func decodeMetadata(b []byte, dst any) error {
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	return d.Decode(dst)
}
func Metadata(root string) error {
	for _, role := range roles {
		path := filepath.Join(root, role, "go.mod")
		b, err := goCommand(root, root, "mod", "edit", "-json", path)
		if err != nil {
			return err
		}
		var m modFile
		if err = decodeMetadata(b, &m); err != nil {
			return err
		}
		if m.Module == nil || m.Module.Path != "totipo/"+role || m.Go != "1.23.0" || len(m.Replace)+len(m.Exclude)+len(m.Retract)+len(m.Tool)+len(m.Godebug)+len(m.Ignore) != 0 || m.Toolchain != "" {
			return fmt.Errorf("module metadata policy: %s", role)
		}
		want := map[string]string{}
		if role == "conformance" {
			want = map[string]string{"golang.org/x/crypto": "v0.36.0", "golang.org/x/sys": "v0.31.0"}
		}
		if len(m.Require) != len(want) {
			return fmt.Errorf("unexpected requirements: %s", role)
		}
		for _, r := range m.Require {
			if want[r.Path] != r.Version || r.Indirect != (r.Path == "golang.org/x/sys") {
				return fmt.Errorf("unreviewed dependency: %s", r.Path)
			}
			delete(want, r.Path)
		}
	}
	b, err := goCommand(root, root, "work", "edit", "-json", filepath.Join(root, "go.work"))
	if err != nil {
		return err
	}
	var w struct {
		Go, Toolchain    string
		Use              []struct{ DiskPath, ModulePath string }
		Replace, Godebug []json.RawMessage
	}
	if err = decodeMetadata(b, &w); err != nil {
		return err
	}
	if w.Go != "1.23.0" || w.Toolchain != "" || len(w.Use) != 2 || len(w.Replace)+len(w.Godebug) != 0 {
		return fmt.Errorf("workspace metadata policy")
	}
	seen := map[string]bool{}
	for _, u := range w.Use {
		if u.DiskPath != "./generator" && u.DiskPath != "./conformance" {
			return fmt.Errorf("extra workspace member")
		}
		if seen[u.DiskPath] || u.ModulePath != "" {
			return fmt.Errorf("workspace alias/duplicate")
		}
		seen[u.DiskPath] = true
	}
	return nil
}

type pkg struct {
	ImportPath, Dir      string
	Standard, Incomplete bool
	Module               *module
	Error                *struct{ Err string }
	DepsErrors           []struct{ Err string }
}

func checkDependencyDir(root string, p pkg) error {
	// The Go tool owns standard-library paths. Windows toolchain installations
	// may use junctions that EvalSymlinks cannot traverse. Only non-standard
	// sources need repository/cache/vendor directory policy checks.
	if p.Standard || p.Dir == "" {
		return nil
	}
	resolved, err := filepath.EvalSymlinks(p.Dir)
	if err != nil {
		return fmt.Errorf("resolve dependency directory %s (%q): %w", p.ImportPath, p.Dir, err)
	}
	// A vendor package may not resolve into another source tree.
	if within(filepath.Join(root, "vendor"), p.Dir) {
		// Compare canonical paths on both sides: temporary/repository ancestors
		// may be aliases (for example /var -> /private/var on macOS). Resolve
		// only the root so links out of vendor still fail the containment check.
		resolvedRoot, err := filepath.EvalSymlinks(root)
		if err != nil {
			return fmt.Errorf("resolve repository root %q: %w", root, err)
		}
		if !within(filepath.Join(resolvedRoot, "vendor"), resolved) {
			return fmt.Errorf("vendor source symlink escape: %s", p.ImportPath)
		}
	}
	return nil
}

func Dependencies(root string) error {
	flags, err := goCommand(root, root, "env", "GOFLAGS")
	if err != nil {
		return err
	}
	for _, flag := range strings.Fields(string(flags)) {
		if flag != "-mod=vendor" && flag != "-mod=readonly" {
			return fmt.Errorf("dependency inspection flag requires review: %s", flag)
		}
	}

	cacheBytes, err := goCommand(root, root, "env", "GOMODCACHE")
	if err != nil {
		return err
	}
	cache := strings.TrimSpace(string(cacheBytes))
	for _, role := range roles {
		b, err := goCommand(root, filepath.Join(root, role), "list", "-deps", "-test", "-json", "./...")
		if err != nil {
			return err
		}
		dec := json.NewDecoder(bytes.NewReader(b))
		count := 0
		for {
			var p pkg
			err := dec.Decode(&p)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			count++
			if p.Incomplete || p.Error != nil || len(p.DepsErrors) != 0 {
				return fmt.Errorf("dependency load error: %s", p.ImportPath)
			}
			if err := checkDependencyDir(root, p); err != nil {
				return err
			}
			if p.Standard {
				continue
			}
			if p.Module == nil || p.Module.Replace != nil {
				return fmt.Errorf("unowned/replaced dependency: %s", p.ImportPath)
			}
			if p.Module.Path == "totipo/"+role {
				if !p.Module.Main || !within(filepath.Join(root, role), p.Dir) {
					return fmt.Errorf("repository dependency escape: %s", p.ImportPath)
				}
				continue
			}
			if role != "conformance" || !primitives[p.ImportPath] {
				return fmt.Errorf("unreviewed dependency closure: %s", p.ImportPath)
			}
			versions := map[string]string{"golang.org/x/crypto": "v0.36.0", "golang.org/x/sys": "v0.31.0"}
			if p.Module.Main || versions[p.Module.Path] != p.Module.Version || within(root, p.Dir) && !within(cache, p.Dir) && !within(filepath.Join(root, "vendor"), p.Dir) {
				return fmt.Errorf("dependency resolution escape: %s", p.ImportPath)
			}
			// Resolved primitive sources must belong to their cache module or workspace vendor.
			if !(p.Module.Dir == filepath.Join(cache, filepath.FromSlash(p.Module.Path)+"@"+p.Module.Version) && within(p.Module.Dir, p.Dir)) && !within(filepath.Join(root, "vendor"), p.Dir) {
				return fmt.Errorf("external dependency directory escape: %s", p.ImportPath)
			}
		}
		if count == 0 {
			return fmt.Errorf("empty package closure: %s", role)
		}
	}
	return nil
}
func Check(root string) error {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	root, err = filepath.EvalSymlinks(absolute)
	if err != nil {
		return fmt.Errorf("resolve repository root %q: %w", absolute, err)
	}
	if err = Source(root); err != nil {
		return err
	}
	if err = Metadata(root); err != nil {
		return err
	}
	if err = Wiring(root); err != nil {
		return err
	}
	return Dependencies(root)
}
