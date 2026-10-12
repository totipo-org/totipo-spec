package repopolicy

import (
	"strings"
	"testing"
)

func TestWiringRejectsGenerationAndBypasses(t *testing.T) {
	for _, tt := range []struct{ name, old, replacement string }{
		{"generation-prerequisite", "check: architecture spec-check test conformance verify", "check: generate architecture spec-check test conformance verify"},
		{"generation-recipe", "$(GO) run ./conformance/cmd/totipo-conformance -root . -verify-only", "$(GO) run ./generator/cmd/generate-vectors -root ."},
		{"recursive-generation", "verify:\n", "verify:\n\t$(MAKE) generate\n"},
		{"alias", "verify:\n", "verify: refresh\n\nrefresh: generate\n"},
		{"shell-expansion", "GO ?= go", "GO ?= $(shell go run ./generator/cmd/generate-vectors -root .)"},
		{"override", "GO ?= go", "override GO = generator-wrapper"},
		{"include", "GO ?= go", "include hidden.mk\nGO ?= go"},
		{"remove-architecture", "check: architecture spec-check", "check: spec-check"},
		{"remove-schema-verification", "check: architecture spec-check test conformance verify", "check: architecture spec-check test"},
		{"weaken-verifier", "$(GO) run ./conformance/cmd/totipo-conformance -root . -verify-only", "@true"},
		{"recipe-prefix", ".DEFAULT_GOAL := help", ".RECIPEPREFIX := >\n.DEFAULT_GOAL := help"},
		{"continued-comment", "generate:\n", "# hidden continuation \\\ngenerate:\n"},
		{"recipe-comment-expansion", "verify:\n", "verify:\n\t# $(shell go run ./generator/cmd/generate-vectors -root .)\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := tree(t)
			if !strings.Contains(reviewedMakefile, tt.old) {
				t.Fatal("fixture does not mutate wiring")
			}
			put(t, r, "Makefile", strings.Replace(reviewedMakefile, tt.old, tt.replacement, 1))
			if err := Check(r); err == nil || !strings.Contains(err.Error(), "Make wiring") {
				t.Fatalf("want wiring rejection: %v", err)
			}
		})
	}
}

func TestWiringAllowsHarmlessText(t *testing.T) {
	r := tree(t)
	put(t, r, "Makefile", "# make generate; totipo/shared; python\n\n"+reviewedMakefile+"\n# include hidden.mk\n")
	put(t, r, "historical.md", "go run ./generator/cmd/generate-vectors -root .\n")
	put(t, r, "generator/a.go", "package generator\n// Importing totipo/conformance is forbidden.\n")
	if err := Check(r); err != nil {
		t.Fatal(err)
	}
}
