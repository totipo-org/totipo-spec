// Package policy data below deliberately accepts only reviewed Make wiring.
// Make is executable configuration: includes, expansions, overrides and new
// targets require policy review rather than a partial Make parser.
package repopolicy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// reviewedMakefile enumerates the supported commands and dependency edges.
// Ordinary Make comments and blank lines are ignored; recipes retain shell text.
// GO/GOCACHE overrides supplied by the caller are trusted developer inputs.
const reviewedMakefile = `
GO ?= go
# Keep Go's build cache writable in the Nix jailed development environment.
export GOCACHE ?= $(CURDIR)/.direnv/go-build

.DEFAULT_GOAL := help
.PHONY: help spec-check test race fuzz conformance verify check architecture generate
.NOTPARALLEL: check

help:
	@printf '%s\n' 'make spec-check  Check the v1/r19 specification structure' 'make test        Run Go tests and independent wire oracle' 'make conformance Run all v1/r19 requirements' 'make verify      Verify manifest, files, and moving profile' 'make check       Run all required checks' 'make race        Run Go race tests' 'make fuzz        Run bounded parser fuzzing' 'make architecture Enforce the two-module implementation boundary' 'make generate    Deliberately regenerate fixtures and pins'

spec-check:
	$(GO) run ./conformance/cmd/check-spec -root .

test:
	$(GO) -C generator test -count=1 ./...
	$(GO) -C conformance test -count=1 ./...

race:
	$(GO) -C generator test -race -count=1 ./...
	$(GO) -C conformance test -race -count=1 ./...

fuzz:
	$(GO) -C conformance test ./internal/object -run '^$$' -fuzz FuzzDispatch -fuzztime 10s -parallel 2
	$(GO) -C conformance test ./internal/cryptov1 -run '^$$' -fuzz FuzzOpen -fuzztime 10s -parallel 2
	$(GO) -C conformance test ./internal/graph -run '^$$' -fuzz FuzzArrivalAndDisappearance -fuzztime 10s -parallel 2

conformance:
	$(GO) run ./conformance/cmd/totipo-conformance -root .

verify:
	$(GO) run ./conformance/cmd/totipo-conformance -root . -verify-only

architecture:
	$(GO) run ./conformance/cmd/check-architecture -root .

generate:
	$(GO) run ./generator/cmd/generate-vectors -root .

check: architecture spec-check test conformance verify
`

func wiringLines(body string) []string {
	var lines []string
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		// Recipe comments still undergo Make expansion; continued Make comments
		// can swallow a target and attach its recipe to the preceding target.
		// Neither is a harmless comment for executable-wiring purposes.
		if trimmed == "" || strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(line, "\t") && !strings.HasSuffix(trimmed, "\\") {
			continue
		}
		// Preserve recipe identity; leading non-recipe whitespace is insignificant.
		if strings.HasPrefix(line, "\t") {
			trimmed = "\t" + strings.TrimSpace(line[1:])
		}
		lines = append(lines, trimmed)
	}
	return lines
}

// Wiring fails closed on unreviewed executable Make configuration. This is
// intentionally an allowlist, not a general Make evaluator or a CI/Nix parser.
func Wiring(root string) error {
	b, err := os.ReadFile(filepath.Join(root, "Makefile"))
	if err != nil {
		return err
	}
	want, got := wiringLines(reviewedMakefile), wiringLines(string(b))
	if len(want) != len(got) {
		return fmt.Errorf("Make wiring requires review: executable line count")
	}
	for i := range want {
		if want[i] != got[i] {
			return fmt.Errorf("Make wiring requires review: executable line %d", i+1)
		}
	}
	return nil
}
