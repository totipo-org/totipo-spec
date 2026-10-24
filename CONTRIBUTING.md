# Contributing

The [v1/r19 specification](spec/totipo-vault-format-v1.md) is normative. Design
checkpoints in `review/` are historical inputs, and the Go consumer is conformance
evidence rather than production code.

Protocol changes need rationale and corresponding cases. Review exact wire changes,
case IDs, schema changes, and moving-profile pins together. Do not regenerate
expectations merely to make a failing consumer pass. Preserve historical reports
and checkpoints. Pre-RC case IDs may be removed or replaced when their concepts
are removed; document every delta.

## Qualification policy

Start milestones on clean `main`, unless the user explicitly authorizes continuing
identified uncommitted work. Record the starting HEAD, Git state, protocol revision,
and inherited qualification status. Preserve historical reports. Work only within
the authorized repository and scope.

The repository's Nix flake is the supported development environment. The owner
uses direnv; agents use ordinary Go commands in the activated environment. The
human-owned authoritative routine qualification command is `nix flake check`.
Agents do not execute Nix. Remote CI supplies supplemental portability evidence
after an owner-authorized commit/push.

Run the human gate from the repository root:

```sh
nix flake check
```

### Classify before qualifying

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

### Baseline and focused development

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

### Final host boundary

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

### Generation and independent verification

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

### Documentation, source invariance and invalidation

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

### Human gate and Git-aware source membership

Request the human `nix flake check` only after all intended milestone files are
Git-known and included source inputs are frozen. Existing tracked changes remain
unstaged. New intended files may use `git add --intent-to-add <file>` only;
require `git diff --cached --quiet` to pass. Do not stage content without explicit
authorization. Leave intent-to-add membership for owner review/commit.
Include the milestone review report among intended Git-known files before the
human gate. The Git-aware working tree includes unstaged tracked modifications.

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

Use the Nix development environment through the owner's direnv setup or
`nix develop`. On Linux, the human gate runs the complete pinned suite:
formatting, structural/schema and integrity checks, independent Go wire-oracle
tests, all module tests and manifest cases, vet, race, and bounded fuzzing.
Linux CI uses the same gate; supplemental macOS/Windows and Linux Go 1.23 jobs
provide portability evidence. Full Nix qualification is Linux-only; Darwin
development shells remain available. There is no application package or
additional build command. Nix uses temporary caches and fixed vendored Go inputs,
with Go proxy and checksum-service access disabled.

## Historical audit reproduction

Historical one-off audit programs are reproduced from the historical Git revision
that contains them. Current main need not retain executable migration tools after
their migration responsibility has ended. Use a separate clean historical checkout
with repository history available, rather than running an old script against the
current corpus:

| Audit | Historical checkout revision | Command from that checkout |
| --- | --- | --- |
| r14 → r15 | `5ff6977ae004a9c03a410b38f0f9d73ca2bc76d4` | `python3 tools/audit_r15.py` |
| r15 → r16 | `54f38d13d62dc811c98489eb1867d63a9bef235a` | `python3 tools/audit_r16.py` |

For example, `git worktree add --detach /tmp/totipo-r15-audit
5ff6977ae004a9c03a410b38f0f9d73ca2bc76d4` creates the r15 checkout; run its
command from that directory. Use the r16 revision and a different directory for
r16. The scripts read their historical worktree and fixed baseline Git objects
(`0b63f77886f95cf9bab3ceb9cbb0916177fe3589` and
`0a01fc1f0493ec1fefc0c4b89f8b22490d586852`, respectively). The r15 legacy-reference
inventory also reflects files present in that checkout. This retrieval convention
relies on normal repository history retention, not an external archival guarantee.
The original [r15 report](review/V1_R15_STATE_MODEL_SIMPLIFICATION_REPORT.md),
[r15 audit JSON](review/V1_R15_CASE_AUDIT.json),
[r16 report](review/V1_R16_REWRITE_REPORT.md), and
[r16 audit JSON](review/V1_R16_CASE_AUDIT.json) remain historical evidence;
their original relative source links and commands refer to their historical tree.

## Fixture maintenance and implementation independence

See [FORMAT.md](vectors/FORMAT.md)
for deliberate fixture maintenance. Every physical case must be manifest-listed
exactly once and required by the moving profile. Abstract graph/storage cases must
state their modeling limits. A synthetic cycle does not claim constructible
cryptographic bytes, and a modeled durability outcome does not prove host crash
behavior.

Report suspected vulnerabilities through [SECURITY.md](SECURITY.md). The project
uses the [Apache License 2.0](LICENSE).

The sibling `totipo/generator` and `totipo/conformance` modules intentionally own
separate Totipo data models, protocol constants, and composition logic. Small
duplicated formulas are independence evidence. No shared repository Totipo Go
package or cross-role subprocess delegation is allowed. `go.work` is developer
convenience, not isolation. `go run ./conformance/cmd/check-architecture -root .`
checks source, metadata, and resolved test dependencies. Serialized artifacts on
disk are the interface.

Nix and GitHub Actions directly invoke Go commands. Review those entry points
to preserve required qualification and keep generation separate. Architecture
checks enforce semantic isolation; they do not validate orchestration text.
Generator unit tests check writer policy without executing fixture generation.
Full conformance and verify-only build only the consumer:

```sh
go run ./conformance/cmd/totipo-conformance -root .
go run ./conformance/cmd/totipo-conformance -root . -verify-only
```

Full conformance includes artifact verification. CLI regressions hash canonical
artifact paths and bytes around successful and rejected verification in temporary
artifact trees.

For standalone module qualification, use `GOWORK=off go -C generator test
-mod=readonly -count=1 ./...` and the corresponding conformance command; vet
each module with `GOWORK=off go -C <module> vet -mod=readonly ./...`.

`go run ./generator/cmd/generate-vectors -root .` is an explicitly authorized
maintenance operation that writes fixtures and pins. Routine Nix and CI
qualification never regenerates vectors. Qualify a
writer change by copying spec/vectors/requirements to a temporary tree, retaining
the RFC TOTP and bootstrap source known answers, running
`go run ./generator/cmd/generate-vectors -root /path/to/temp-tree`, and comparing
all output bytes and inventories before normal verification. Do not derive
authored graph/storage expectations from verifier results.

Spec/prose policy runs via `go run ./conformance/cmd/check-spec -root .`.
It requires only the active spec
and archived pre-r16 history. `vectors.Read` validates duplicate-safe JSON against
both independently authored checked-in schemas before strict typed acceptance;
required fields are checked by exact member presence rather than Go zero values.
Schema keyword or pattern extensions require explicit support and parity review.

The independent `conformance/internal/wireoracle` Go tests preserve selected
wire/key, VAULT_ID, capacity, and post-AEAD fixture checks. They use only stdlib
and their own duplicate-safe serialized-artifact reader; the architecture guard
rejects imports of generator, verifier, or external implementation packages.
Do not replace their formulas or constants with production helpers. These tests
do not provide independent AES-GCM or Argon2 evidence. No Totipo support-tooling
command requires Python; the historical audit commands above apply only to their
specified historical checkouts. Normal checks do not execute the generator.
