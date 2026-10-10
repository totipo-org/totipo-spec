# Totipo

Totipo is an encrypted, append-only TOTP vault format operating on a configured
durable store. Synchronization is optional and external.

The current normative specification is **v1/r19**, a design draft with moving
pre-release-candidate conformance evidence:
[Totipo Vault Format v1](spec/totipo-vault-format-v1.md).

The design uses complete-state TOKEN assertions, fixed 1024-byte encrypted objects,
keyed content addressing, a maximum of four explicit parents, and whole-state
conflicts. Cycles form causal-equivalence groups. Optional CLIENT_NAME and CLIENT_TIME
are informational and preserved exactly while each object is represented and
across all stages of one fold, without requiring persistent remembered history. Tombstones retain the full credential. Applications describe
observed state truthfully; known complete credentials remain available for TOTP
computation regardless of lifecycle or historical status.

A store contains canonical regular-file `vault` and directory `objects-v1/`. Only exact v1 TOKENs
contribute semantic state. Unknown sibling families are outside v1 interpretation;
future families define their own compatibility relationships. The vault root
provides authoring authority. VAULT_ID (SHA-256 of exact VAULT bytes) provides optional stable recognition.
VAULT is an immutable create-once bootstrap; object history is immutable. Credential,
KDF-policy, or root changes require migration to a fresh independent vault.
No per-client persistent graph database is required. Store loss or rollback can
lose history. Totipo validates and authenticates observed content, but does not
cryptographically guarantee a complete or freshest store view. Stronger freshness,
rollback detection, or history retention must be explicitly supplied by the
storage environment or application; synchronization remains optional.

| Path | Role |
| --- | --- |
| `spec/` | Current normative protocol |
| `vectors/` | Exact byte, negative, graph, fold, and storage workflow cases |
| `requirements/` | Moving pre-RC pins |
| `conformance/` | Go reference consumer and explicit fixture generator |
| `review/` | Historical decisions and review evidence |
| `tools/` | Structural and schema checks |

Normal repository qualification uses the pinned Nix environment:

```sh
nix flake check path:.
```

On Linux this runs formatting, structural/schema and integrity checks, Python
unit tests, all Go tests and manifest conformance cases, vet, race tests, and the
three bounded fuzz targets. Linux CI uses this authoritative gate. Supplemental
macOS/Windows jobs and Linux Go 1.23 checks provide portability evidence.
There is no application package or additional build command.

Use `nix develop` or the existing direnv setup when available. Go 1.23+ and Python
3.9+ are supported for direct local development checks:

```sh
make check
make race
make fuzz
go -C conformance vet ./...
```

`make conformance` executes every manifest case. `make verify` checks schemas,
case hashes, physical case coverage, and exact requirements pins. Normal checks
never regenerate fixtures. Explicit generator maintenance remains separate.
Make defaults to a writable Go cache under `.direnv/` and honors caller-provided
`GOCACHE`; Nix qualification uses temporary caches and fixed vendored Go inputs,
with Go proxy and checksum-service access disabled. Python tools use only the
standard library. Darwin development shells remain available; full Nix
qualification checks are Linux-only.
See the [case contract](vectors/FORMAT.md) and [Go consumer](conformance/README.md).

The [r19 immutable-VAULT report](review/V1_R19_IMMUTABLE_VAULT_REPORT.md) records
the create-only model, unchanged wire bytes, and adjusted conformance.
The [r18 hardening report](review/V1_R18_HARDENING_REPORT.md) records application
safety, conformance scopes, editorial cleanup, and the unchanged portable corpus.
The [r17 clarification report](review/V1_R17_HARDENING_REPORT.md) records the
threat-model clarification and unchanged 90-case corpus.
The [r16 rewrite report](review/V1_R16_REWRITE_REPORT.md) records the baseline,
five checkpoint hashes, per-case migration, and validation. Older reports remain
historical evidence. No release candidate is frozen by this rewrite.
