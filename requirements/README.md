# Requirements profiles

[`v1-pre-rc.json`](v1-pre-rc.json) is the **moving pre-RC** profile for v1/r19.
It pins the exact specification, entire manifest, manifest schema, case schema,
and every manifest case ID/hash. There are no conditional capabilities.

`go run ./conformance/cmd/totipo-conformance -root . -verify-only` checks schemas,
exact physical case coverage, hashes, and profile membership. Omit `-verify-only`
for full conformance, which includes that verification and executes every case.
Neither regenerates fixtures; the committed corpus is authoritative input.
Explicit maintenance is described in [FORMAT.md](../vectors/FORMAT.md).

Do not freeze an RC profile until independent implementation interoperability and
required security/platform evidence are complete. The r16 rewrite does not tag,
release, or freeze anything.
