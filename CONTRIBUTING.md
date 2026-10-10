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
nix flake check path:.
```

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

Use the preserved Nix environment when available. See [FORMAT.md](vectors/FORMAT.md)
for deliberate fixture maintenance. Every physical case must be manifest-listed
exactly once and required by the moving profile. Abstract graph/storage cases must
state their modeling limits. A synthetic cycle does not claim constructible
cryptographic bytes, and a modeled durability outcome does not prove host crash
behavior.

Report suspected vulnerabilities through [SECURITY.md](SECURITY.md). The project
uses the [Apache License 2.0](LICENSE).
