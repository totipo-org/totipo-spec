# Go v1/r19 reference/conformance consumer

This Go 1.23+ module is the repository's reference consumer, not a production client.
Run `make check`, `make race`, and `make fuzz` from the repository root. The CLI
accepts `-root /path/to/repository` and `-verify-only`. Every manifest case is
required; there is no optional capability selection. These cases provide core
protocol and modeled store-operation evidence under the specification’s Section 20
scopes. Passing them does not establish application conformance or live backend
crash durability.

Read-only integrity belongs to `vectors.Read` and `vectors.VerifyProfile`: exact
revision, raw artifact/case hashes, physical corpus coverage, and the exact
manifest/profile required set. `make spec-check` checks prose structure/history;
`make verify` runs schema-aware Go artifact and integrity verification.
The schemas remain the presence/type contract, including required false and zero
values; Go strict decoding and `ValidateShape` protect direct consumer inputs.
Verification never updates pins or derives expected outcomes from evaluation.

Packages separate TLV framing, exact TOKEN grammar, object/root cryptography,
TOTP computation, causal-equivalence graph evaluation, and storage outcomes.
Storage bytes pass physical-size, AEAD, semantic-length, zero-padding, and keyed-ID
checks before grammar validation. Required post-AEAD negative fixtures prove that
authentication alone cannot bypass length, padding, or keyed-ID checks. Invalid grammar contributes no TOKEN state.

`graph.Store` models a supplied valid observed set, not persistent security memory.
It uses cycle-safe reachability to compute all members of maximal SCCs. Missing
parents remain unresolved. Head records and retained historical nodes preserve exact
client metadata independently of the semantic-value projection for as long as
the object is represented; this creates no persistent history obligation. Fold construction
copies one operation’s complete value and exact metadata into every stage. Values
in graph vectors are symbolic complete tuples;
real TokenValue extraction is tested against canonical objects. Defensive same-ID
failure excludes that identity from the model without poisoning unrelated tokens.
The model is deliberately small, not a scalable graph-index implementation.

`internal/storage` models exact namespaces, candidate names, observation diagnostics,
immutable publication, initial VAULT creation, no-replace creation with post-publication revalidation, orphan-object safety,
and preservation of VAULT through object publication.
Observed bootstrap entries must be regular files and object namespaces directories;
wrong types are not followed or traversed. Backend booleans describe trusted API outcomes. They do not prove file/namespace
persistence, staging behavior, host syscall races, or remote completeness. A live
adapter must implement the specification's durability and truthful reporting rules.

Tests never regenerate fixtures. For an intentional reviewed protocol update:

```sh
make generate
make verify
make conformance
```

The sibling generator module independently constructs deterministic TOKEN bytes using public test roots and
preserves the independent RFC TOTP answers. Semantic graph/workflow expectations
are specified explicitly, not queried from the evaluator. Review the diff and
pins; independent interoperability remains necessary before RC freeze.
See the [case contract](../vectors/FORMAT.md) and
[r16 rewrite report](../review/V1_R16_REWRITE_REPORT.md).

The verifier-owned stdlib `internal/schema` interpreter validates the checked-in
manifest and case contracts before strict typed decoding in `vectors.Read`.
It detects duplicates at every depth, keeps integer tokens exact, counts Unicode
code points for schema string lengths, and checks required member presence
independently of Go zero values. Unsupported schema keywords or patterns fail
closed. Existing shape, hash, path/inventory, profile, and semantic checks remain.
Verify-only performs artifact verification and skips semantic execution; full
conformance executes semantics afterward.

The separate `cmd/check-spec` / `internal/specpolicy` command reads only the active
spec and archived history. The architecture guard keeps this policy independent
of semantic packages. The independent stdlib-only `internal/wireoracle` tests
read committed artifacts with their own duplicate-safe reader and independently
check selected HMAC/key composition, TOKEN fields, VAULT_ID, size, and post-AEAD
fixture predicates. Their import guard prohibits production and generator helper
reuse. They provide no independent AES-GCM or Argon2 evidence.
Generator/verifier isolation from W3A is unchanged. No Totipo support-tooling
command requires Python.
