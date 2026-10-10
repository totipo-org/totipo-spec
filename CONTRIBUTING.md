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

Direct local development commands remain useful:

```sh
gofmt -l conformance
make check
make race
make fuzz
go -C conformance vet ./...
```

Use the preserved Nix environment when available. See [FORMAT.md](vectors/FORMAT.md)
for deliberate fixture maintenance. Every physical case must be manifest-listed
exactly once and required by the moving profile. Abstract graph/storage cases must
state their modeling limits. A synthetic cycle does not claim constructible
cryptographic bytes, and a modeled durability outcome does not prove host crash
behavior.

Report suspected vulnerabilities through [SECURITY.md](SECURITY.md). The project
uses the [Apache License 2.0](LICENSE).
