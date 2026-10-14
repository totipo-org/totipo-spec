# Nix-first interactive qualification

Start milestones on clean `main`, unless the user explicitly authorizes continuing
identified uncommitted work. Record the starting HEAD, Git state, protocol revision,
and inherited qualification status. Preserve historical reports. Work only within
the authorized repository and scope.

The repository's Nix flake is the supported development environment. The owner
uses direnv; agents use ordinary Go commands in the activated environment. The
human-owned authoritative routine qualification command is `nix flake check`.
Agents do not execute Nix. Remote CI supplies supplemental portability evidence
after an owner-authorized commit/push.

## Classify before qualifying

Record every applicable classification in substantial reports:

| Classification | Changed inputs |
| --- | --- |
| `DOCS_ONLY` | Documentation/process without qualification source changes |
| `SPEC_TEXT` | Specification text, including editorial changes |
| `REQUIREMENTS` | Requirement profiles or pins |
| `VECTOR_CORPUS` | Cases, schemas, manifest, or expected outcomes |
| `CONFORMANCE_GO` | Go consumer or its tests |
| `PYTHON_VERIFIER` | Historical Python verification/structure logic or tests |
| `GENERATOR` | Fixture generator or its tests |
| `BUILD_NIX_CI` | Build, dependency, Nix, source-filter, or CI inputs |

Infrastructure-only work must establish protected source/byte invariance.
Qualification-logic changes still require the applicable final checks.

## Baseline and focused development

Record baseline evidence before editing; do not upgrade pending evidence to a
pass. For ordinary semantic/conformance work without a qualified baseline, run
both module tests and full conformance once:

```sh
go -C generator test -count=1 ./...
go -C conformance test -count=1 ./...
go run ./conformance/cmd/totipo-conformance -root .
```

Do not run baseline race/fuzz by default. If changing race/fuzz-sensitive code,
use the relevant focused command. While editing, use the narrowest Go package or
test, for example:

```sh
go -C conformance test -count=1 ./internal/schema
go -C conformance test -count=1 ./internal/graph
go -C conformance test -count=1 ./internal/cryptov1
```

Do not repeatedly run all-package tests, vet, race, all fuzz targets, or request
Nix after minor edits. Preserve failure output and reproduction inputs; reproduce
the specific package/test or bounded fuzz target, fix narrowly, then return to
complete final qualification after stabilization.

## Final host boundary

For semantic, conformance, or qualification-logic changes, run the complete
non-Nix suite once after included inputs stabilize:

```sh
test -z "$(gofmt -l generator conformance)"
go run ./conformance/cmd/check-architecture -root .
go run ./conformance/cmd/check-spec -root .
go -C generator test -count=1 ./...
go -C conformance test -count=1 ./...
go run ./conformance/cmd/totipo-conformance -root .
go -C generator vet ./...
go -C conformance vet ./...
go -C generator test -race -count=1 ./...
go -C conformance test -race -count=1 ./...
go -C conformance test ./internal/object -run '^$' -fuzz FuzzDispatch -fuzztime 10s -parallel 2
go -C conformance test ./internal/cryptov1 -run '^$' -fuzz FuzzOpen -fuzztime 10s -parallel 2
go -C conformance test ./internal/graph -run '^$' -fuzz FuzzArrivalAndDisappearance -fuzztime 10s -parallel 2
git diff --check
```

Preserve all-package race coverage, fuzz seeds, failure semantics and these exact
bounded fuzz settings. Optimization comes from staging expensive checks at
boundaries, never deleting final coverage. Record case counts, package results,
vet/race results, each fuzz duration/result, and Python test executions (currently
zero). Report-only edits excluded from qualification do not require a repeat.
Use writable Go caches; Nix qualification supplies temporary build/fuzz caches.

## Generation and independent verification

The generator and consumer remain separate Go modules with independent Totipo
semantic implementations. Serialized artifacts are their interface. Routine
qualification never invokes the generator. Explicitly authorized maintenance uses:

```sh
go run ./generator/cmd/generate-vectors -root .
```

Do not generate without deliberate authorization. Generator changes require
applicable temporary generation/comparison, review of intended differences, and
normal verification of the committed corpus. Successful generation is not proof
of correctness and must never justify changing expectations to pass.
Full conformance includes artifact verification. Verification without semantic
execution remains available:

```sh
go run ./conformance/cmd/totipo-conformance -root . -verify-only
```

## Documentation, source invariance and invalidation

For docs/process-only work, use applicable static checks and `git diff --check`;
do not automatically rerun race/fuzz or the full suite. Inspect actual consumers
and the flake filter before deciding whether documentation invalidates Nix.
Compare bytes and path additions/deletions against the starting revision for
spec, vectors, requirements, generator/conformance, relevant tools, go.work and
all go.mod/go.sum files. Record absent root module/sum files.

The current qualificationSource includes spec, vectors, requirements, repository
Go/native sources and go.mod/go.sum, go.work, and archived
`review/V1_PRE_R16_REVISION_HISTORY.md`. Root README, CONTRIBUTING, AGENTS and
other review reports are excluded and not consumed by structure checks.
Reinspect dependencies whenever filter or consumers change.

Spec/artifact/semantic or verification-logic changes invalidate applicable final
host qualification and Nix. Generator changes also need generation-comparison
evidence. Nix/lock/vendorHash/filter changes require renewed human Nix and any
affected host checks. Excluded report-only edits do not invalidate results.

## Human gate and Git-aware source membership

Request the human `nix flake check` only after all intended milestone files are
Git-known and included source inputs are frozen. Existing tracked changes remain
unstaged. New intended files may use `git add --intent-to-add <file>` only;
require `git diff --cached --quiet` to pass. Do not stage content without explicit
authorization. Leave intent-to-add membership for owner review/commit.

Keep local caches, secrets, generated output and build logs outside Git source
membership. Use Git-aware current-checkout forms for human Nix commands; explicit
`path:` sources are reserved for intentional raw/non-Git source semantics.
No separate Nix build or package materialization is routinely needed.

Any later included-source edit invalidates the human result; excluded report
updates do not. Record pending gates accurately. Distinguish agent execution,
human-reported execution, cached/fresh provenance and remote CI evidence.
Record baseline/focused/final command counts, incomplete attempts, failures and
human Nix count/status. Leave commits and release operations to the explicitly
authorized user workflow.
