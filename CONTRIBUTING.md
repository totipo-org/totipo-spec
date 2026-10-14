# Contributing

The [v1/r19 specification](spec/totipo-vault-format-v1.md) is normative. Design
checkpoints in `review/` are historical inputs, and the Go consumer is conformance
evidence rather than production code.

Protocol changes need rationale and corresponding cases. Review exact wire changes,
case IDs, schema changes, and moving-profile pins together. Do not regenerate
expectations merely to make a failing consumer pass. Preserve historical reports
and checkpoints. Pre-RC case IDs may be removed or replaced when their concepts
are removed; document every delta.

The agent does not run Nix. The human normal qualification gate is exactly:

```sh
nix flake check
```

Local pre-commit qualification uses the Git-aware working tree, including
unstaged modifications to tracked files. For each new intended milestone file,
the agent may use `git add --intent-to-add <file>` to establish source membership
only. Require `git diff --cached --quiet` to pass: staging actual file contents
remains prohibited without explicit user authorization. Keep existing modified
tracked files unstaged and leave intent-to-add entries for human review/commit.
Make every intended new milestone file, including its review report, Git-known
before the human gate. Ignored/local build, cache, and secret state must remain
outside Git source membership. Explicit `path:` sources remain allowed only for
intentional raw/non-Git paths, such as generated or exported source.

On Linux it runs the complete pinned spec/conformance suite, including formatting,
vet, race, and bounded fuzzing. Linux CI uses the same gate; supplemental
macOS/Windows and Linux Go 1.23 jobs provide portability evidence. No additional
`nix build` is required. Full Nix qualification checks are Linux-only; Darwin
shells remain available for development.

Use the Nix development environment through the owner's direnv setup or
`nix develop`. Follow [AGENTS.md](AGENTS.md) for baseline evidence, focused
iteration, final host checks and execution accounting. Focused commands include:

```sh
go -C generator test -count=1 ./...
go -C conformance test -count=1 ./...
```

Preserve applicable formatting checks (`gofmt -l generator conformance`). Final race and
all three bounded fuzz targets remain mandatory for semantic/conformance code
changes. Preserve failures and diagnose with focused commands before repeating
the final suite. Keep generation separate from verification.

For docs/process-only changes, verify protected-byte invariance and use only
applicable documentation/static checks plus `git diff --check`. Inspect actual
structure dependencies and the flake source filter before requesting human Nix;
request the gate only after all its inputs are frozen. Current README,
CONTRIBUTING, AGENTS, and new review reports are excluded; the archived
`review/V1_PRE_R16_REVISION_HISTORY.md` is included and consumed by structure
checks. Report-only edits outside qualification inputs do not invalidate results.

Optimization must come from staging expensive checks at boundaries,
never from deleting final conformance/race/fuzz coverage.

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

Use the preserved Nix environment when available. See [FORMAT.md](vectors/FORMAT.md)
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
