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

Follow the [interactive qualification ladder](AGENTS.md). Classify each milestone
and record its execution counts. For ordinary semantic/conformance work, run
`make check` once at baseline, use focused Python/Go tests while editing, then
run the final host suite once after inputs stabilize:

```sh
make check
go -C conformance vet ./...
make race
make fuzz
git diff --check
```

Preserve applicable formatting checks (`gofmt -l conformance`). Final race and
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
