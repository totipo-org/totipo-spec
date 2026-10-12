# Totipo support tooling Wave 3B structural and schema validation

## 1. Starting branch and HEAD

Work is restricted to `/home/niki/Sources/totipo-spec`; no sibling repository was
inspected. Required checks ran before inspection: branch `main`, HEAD
`ac6a1489cc4b557155b9abab2f9996ad9d4a5edb`, clean status, no untracked files,
and `git diff --cached --quiet` exit 0. W3A and its boundary report are committed
at this starting revision. W3A's report records host PASS and human-reported Nix
PASS; these are inherited evidence, not W3B qualification. Historical reports
are preserved unchanged.

**W3B complete: implementation, differential parity, final host qualification,
human-reported Nix, and post-Nix frozen-input verification PASS.**
No W3C, protocol, generator, or zero-Python completion is claimed.

## 2. Classification

Requested: SUPPORT_TOOLING, LANGUAGE_SIMPLIFICATION, STRUCTURAL_VALIDATION,
SCHEMA_VALIDATION, WAVE_3B; conceptual PYTHON_TOOLING_REMOVAL.
Actual AGENTS ladder: CONFORMANCE_GO, PYTHON_VERIFIER (replacement/removal of
verification implementation/tests), BUILD_NIX_CI, DOCS_ONLY (current guidance).
No SPEC_TEXT, VECTOR_CORPUS, REQUIREMENTS, GENERATOR, or protocol change.

## 3. Baseline protocol and artifact inventory

Protocol **v1/r19**, moving pre-RC; **92 manifest cases, 92 physical case JSON
files, 92 required profile entries**. Protected inventory: spec 1, vectors 97,
requirements 2: **100 files**. Generator: 8 files, all unchanged.
Starting workspace has exactly `./generator` and `./conformance`; exactly two
modules, no root module, local replace bridge, or cross-role import. The W3A
architecture guard confirms the unchanged generator stdlib-only closure.

Before edits, `/tmp/totipo-w3b-ac6a1489/` captured protected/generator paths and
SHA-256 values (`baseline-hashes.json`), exact bytes of all five Python files
(`baseline-tools/`), relevant Make expansions (`make-*.txt`), and CI commands
(`baseline-ci.yml`). `start-tree/` is an outside archive of starting HEAD for
vendor comparison. No generator was run.

Protected artifact SHA-256 values remain:

| Artifact | SHA-256 |
| --- | --- |
| Active spec | `8bb76b890086eb4bf89271edd5861e02f833f3ffa11a83b5865080ed4e4228cb` |
| Manifest | `3953dbc315b4dcb3d0d31dd399bf82a15cb38c49fde2e2b93ea81c5cfe0ec714` |
| Manifest schema | `f265771f904be57d19602dd6da9a572c16b3f80fdd166c5177721aaa3f9a91ec` |
| Case schema | `99805d442f13872cd4febe9ac8ae4f36fd25607580bc1e457ff5ae2efe199e23` |

The profile and all other protected hashes are in the captured inventory.
Exact byte/path equality establishes unchanged case IDs, expectations, hashes,
profile required set, pins, schemas, and revision without regeneration.

## 4. Baseline qualification

One baseline `make check`, conformance vet, `make race` (both modules), and
`make fuzz`: all PASS. The explicit W3B request requires broader baseline
coverage than the ladder's default. No baseline Nix or generator execution.
Python discovery measured **15 tests, OK**; a separate verbose discovery also
passed all 15 and captured individual names in `baseline-python-names.log`.
Both original checker CLIs passed; full/verify-only passed 92 cases.

Individual baseline test names (all PASS):

| Module/class | Methods |
| --- | --- |
| test_check_spec.StructureTests | test_current_structure_passes; test_structure_mutations_fail |
| test_check_vectors.SchemaContractTests | test_duplicate_members; test_head_records_required_and_exact_metadata_types; test_manifest_revision_and_retired_extension; test_nested_unknown_and_mixed_payloads; test_no_vault_replacement_or_rewrap_contract; test_post_aead_shape; test_required_false_result_cannot_be_omitted; test_required_zero_and_operation_fields; test_unsigned_timestamp_exactness; test_workflow_boolean_presence |
| test_r16_vectors.IndependentWireTests | test_fields_and_key_intermediates; test_post_aead_fixture_inputs; test_vault_id_and_maximum |


Baseline Go tests/race: generator cmd/generate-vectors, internal/fixture,
internal/wire; conformance cryptov1, graph, object, repopolicy, storage, tlv,
totp, vectors: all PASS. Existing check-architecture/totipo-conformance command
packages compiled without command tests. Vet PASS. Fuzz results are in §33.

## 5. Python responsibility inventory

| Starting file | Responsibility | W3B disposition |
| --- | --- | --- |
| tools/check_spec.py | Two-file prose/history/registry/policy assertions | Replaced with Go command and policy package; deleted after parity |
| tools/test_check_spec.py | Current two-input CLI success and eight mutations | Ported policy/CLI tests; deleted after parity |
| tools/check_vectors.py | Duplicate-safe JSON, schema subset, listed-path containment; reader used by oracle | Normal acceptance moved to Go; file retained byte-identical for oracle |
| tools/test_check_vectors.py | Ten schema contract/presence/type/operation tests | Ported/expanded in Go; deleted after parity |
| tools/test_r16_vectors.py | Three independent selected wire/key/VAULT/post-AEAD tests | Retained byte-identical and active |

Starting Python: 486 lines in five files. Removed: 282 lines in three files.
Remaining: 204 lines in two files. Python remains an independent evidence path.

## 6. Go check-spec architecture

`conformance/cmd/check-spec` calls `internal/specpolicy.Check`; policy imports
only stdlib. Default root is `.`; optional `-root` selects a temporary tree.
The only repository inputs are the active spec and archived pre-r16 history.
Vectors, schemas, requirements, and semantic packages are unnecessary. The
command remains separate from vector semantic execution and prints the existing
success message. Violations produce contextual diagnostics and exit 1;
argument refusal returns 2. `main` maps the tested run function to process exit.

## 7. Complete migrated spec-rule inventory

All original assertions were inventoried before implementation and retained:

- Split at Section 21; exact ordered headings 1–21; active numeric Section
  references restricted to 1–20; protocol 1/revision r19 header.
- Exact archive link; r1–r15 headings in archive, r16–r19 in current history.
- Token-boundary retired list: DEVICE, DEVICE_ID, P-256, ECDSA, DER, SIGNATURE,
  AUTHOR_DEVICE_ID, AUTHOR_TIME, PROVENANCE, OPAQUE_ROUTABLE, OPAQUE_UNSCOPED,
  OBJECT_TYPE, OBJECT_VERSION, VAULT_BINDING, READY, PROCESSING_INCOMPLETE,
  HISTORY_REGRESSION, HISTORY_MEMORY_LOST, candidate-use, K_signature_context.
  Case-insensitive provenance/opaque exclusions remain separate.
- Sole TOKEN grammar; MAX_PARENTS=4 and count bound; 1005/1006/1024 anchors;
  exact VAULT_ID and causal-equivalence formula wording.
- Create-only state; optional/external synchronization; no exhaustive
  enumeration requirement; parent availability not a publication prerequisite.
- Exact ordered 12-tag registry, 0x0001–0x000c: TOKEN_ID, PARENT_COUNT, PARENT_ID,
  STATUS, ISSUER, ACCOUNT, ALGORITHM, DIGITS, PERIOD, SECRET_BYTES, CLIENT_NAME,
  CLIENT_TIME.
- Layout block starts with `vault`; no uppercase VAULT layout entry or VAULT
  pathname claim. Exact entry table: vault regular file, objects-v1 directory.
- Fold table has exactly TokenValue/CLIENT_NAME/CLIENT_TIME; metadata retains
  same exact presence; TokenValue tuple excludes metadata; head descriptions
  include metadata and value and prohibit erasing head metadata.
- Symlink follow/interpret/traversal bans; no openat2, O_NOFOLLOW, O_PATH, statx,
  SecureDirectoryStream prescriptions; optional implementation hardening anchor.
- Represented/returned/exposed/retained metadata preserved exactly; no persistent
  retention requirement; no advisory-history/remembered-head/cache/cross-run rule.
- Informed authorship despite unavailable history; new CSPRNG TOKEN_ID and empty
  parent set anchors.
- Explicit confirmation before empty-password creation.
- No silent initialization; explicit recovery/reconfiguration/new-location
  workflow; orphan objects not garbage or deleted; reopen regular vault;
  NO-REPLACE; persistence/publication/revalidation before reporting success.
- No secure-erasure claim from token tombstoning/deletion.
- No in-place password change/rewrap or root reuse; INFORMATIVE migration note
  and not-yet-synchronized-data limitation.
- Retired VAULT_FINGERPRINT, totipo/v1/vault-fingerprint, CURRENT != BASE,
  stageReplacement, compare-before-replace, latest password wrapper exclusions.
- CSPRNG 32-byte root; lifetime immutability; no in-place rotation;
  BOOTSTRAP_VERSION=0x01; never newer; alternate same-root representation anchor.
- No VAULT replacement even when exact-identical.
- Core protocol/Store-writer/Application scope rows; core metadata preservation;
  Sections 12/17 metadata scope; applicable core requirements; GUI boundary;
  portable corpus insufficient for application conformance; no stronger atomic
  primitive imposed by scope.

History stays outside active retired-term checks. Artifact pins remain the
verifier's responsibility. No brittle prose guard was removed or redesigned.
Text handling preserves universal newlines, Python Unicode whitespace, and
Unicode decimal section references; explicit parity boundaries exercise these.

## 8. Spec-policy test parity

Policy tests cover current files, the exact original first-occurrence mutations,
expanded literal/structure mutations, archive mutation, every retired term,
history exemptions, uppercase/underscore token boundaries, CRLF, Unicode
references and whitespace. Command tests run against a temporary tree containing
only the two inputs, prove success/mutation failure, missing-input refusal, and
argument errors. **138 old/new CLI comparisons, zero disposition drift** (§23).
No repository dependency was introduced beyond the original two-file boundary.

## 9. Schema-engine architecture

Verifier-owned stdlib `internal/schema` parses generic JSON, compiles the checked-in
schema source into a validated tree, and validates instances. It imports no
semantic verifier or generator package. Schemas remain independently authored,
byte-identical JSON contracts; neither schema nor Go structs generate the other.
No external JSON-schema/CLI dependency was added.

## 10. Duplicate-safe JSON model

A `json.Decoder` token walk with `UseNumber` checks each object member before
insertion into `map[string]any`. Arrays are `[]any`; scalars retain JSON types.
Duplicate names, including escaped spellings resolving to the same name, fail
at every depth. Trailing JSON and invalid UTF-8 fail. Exact case remains distinct.
`vectors.Decode` uses this parser too, replacing its old duplicate-only walk,
then retains DisallowUnknownFields and the trailing-input check.

## 11. Exact presence

`required` tests object membership with the Go comma-ok lookup. Missing, null,
false, zero, empty string, and empty array remain distinct. No typed zero value
substitutes for presence. Tests remove **5,625 required members** across the
manifest and 92 cases, including selected conditional branches, and assert that
false/zero/empty-array/empty-string representatives were encountered. The CLI
removes required TOTP t0 in a repinned temporary tree to demonstrate that a
plausible typed zero cannot bypass schema acceptance.

## 12. Keyword support

Supported current subset: $schema, $id, title, description annotations; type,
const, enum; properties, required, additionalProperties; items, minItems,
maxItems; minimum, maximum; pattern, minLength, maxLength; allOf, if, then, else.
Current source inventory: manifest 11 type declarations, case 366; case has
15 if/then declarations and 3 allOf declarations. No current else declaration,
but the supported else path has focused tests. Allowed type names are object,
array, string, integer, boolean, null; current numeric bounds use integer tokens.
This deliberately is not a general JSON Schema implementation.

## 13. Schema-source validation

Compilation parses duplicates safely, requires object schema nodes, checks all
keywords and value types, type/required/enum lists, child schemas, nonnegative
item/length bounds, exact integer range bounds, valid supported regexes, and
conditional structure. All properties/items/conditional/allOf descendants are
preflighted even if no current instance visits them. Unknown keywords or future
unsupported patterns fail closed. Tests include duplicate source keys, malformed
nodes/lists/bounds/annotations and unknown keywords in inactive branches.

Measured intentional strengthening: old Python accepts `{}` against a schema
whose unused property contains an unsupported keyword; Go rejects that schema
source eagerly. This is the explicitly required W3B source-preflight boundary,
not drift in acceptance under either current contract. The separate measurement
is in `source-preflight.log`; ordinary parity rows use valid current/subset
contracts, plus visited unsupported-keyword/duplicate-source rejections.

## 14. Numeric semantics

Numbers retain json.Number spelling. Integer typing requires an integer token,
then compares through math/big.Int: no float64 rounding, including u64max.
Fractions and exponent spellings fail integer type; bool is a separate type.
Tests cover zero, negatives, lower/upper bounds, u32max/u64max, overflow, 1.0,
1e0, 0e0, and bool/string/null substitutions.

Python's actual const/enum equality quirk is preserved with exact math/big.Rat
comparison: const true accepts 1/1.0/1e0; TOTP const 0 accepts false/0.0/0e0.
Type=integer still rejects them. The existing strict Go typed decoder subsequently
rejects unsuitable numeric spellings/bools for typed fields, just as the old
combined Python-plus-Go qualification did. This was measured, not silently fixed.

## 15. Unicode length semantics

Schema lengths use Unicode code points (utf8.RuneCountInString), matching Python
len rather than UTF-8 bytes. Tests include accented characters, supplementary
characters, decomposed sequences, and an escaped lone surrogate (one replacement
rune in Go, one code point in Python for the length check). Protocol wire-byte
limits remain in unchanged semantic packages, separate from schema string length.

## 16. Regex compatibility

The union of current patterns has **seven** distinct expressions:
`^[0-9a-f]{64}$`, `^r[0-9]+$`,
`^v1\.[a-z0-9-]+(?:\.[a-z0-9-]+)+$`, `^(?:[0-9a-f]{2})*$`,
`^[0-9]{6,8}$`, `^[0-9a-f]+$`, `^[0-9a-f]{174}$`.
(The early progress update's eight-pattern count was corrected by the inventory.)
Go handles their classes, repetition, and noncapturing groups directly. Validation
also tries stripping exactly one terminal newline to reproduce Python `$`;
length validation still counts the original string. Focused/differential tests
cover accepted boundaries, final newline, double newline, CRLF, non-ASCII and
invalid prefix/suffix. An explicit allowlist rejects unreviewed future syntax.

## 17. Conditional validation

allOf checks every subschema. if selects then on success and else on validation
failure; missing branches impose no additional condition. All current operation
contracts and graph add/remove/evaluate contracts are interpreted independently
of ValidateShape. Presence-removal tests walk selected branches; differential
mutations remove payloads, add incompatible payloads, alter operations, and
exercise enums/consts/branch restrictions. Dedicated tests exercise else.

## 18. Manifest integration

Read compiles both schema sources and validates raw manifest JSON against
manifest.schema.json before strict typed decode. Existing independent expected
format/protocol/revision, nonempty inventory, categories/IDs, normative status,
kind, paths, sections, hashes, duplicate IDs/paths and physical inventory checks
remain. Profile checks remain in VerifyProfile after Read as before.

## 19. Case integration

Every listed case passes case.schema.json before typed decoding. All nine
operation forms (crypto, dispatch, bootstrap, totp, graph, fold, storage,
workflow, post-aead), nested structures, and optional metadata are covered.
Current cases all pass. The original ten Python contract tests have explicit
Go equivalents plus broader required-member, scalar and schema-source tests.
Optional client_name/time remain optional but are exact-type checked when present.

## 20. Path and inventory interaction

Existing lexical safety, resolved containment, root aliases, symlink handling,
hash checks and all-and-only physical inventory code are retained unchanged.
No second path validator was introduced. Old Python/new Go path differential
rejects all four fixtures: parent escape, absolute outside path, file symlink
escape, directory symlink escape. Existing Go path/integrity tests still run,
including host-independent lexical/volume tests and root-alias acceptance.
Pre-existing stricter Go lexical/inventory policy remains authoritative.

## 21. vectors.Read acceptance order

Compile both duplicate-safe schema sources → duplicate-safe manifest schema
validation → strict manifest decode/header checks → existing entry/path/resolved
containment checks → case bytes/hash → duplicate-safe case schema validation →
strict case decode → identity/expectation and ValidateShape → complete physical
inventory. A successful return therefore includes schema, duplicate, path,
inventory and typed acceptance. CLI then checks profile pins, and only full mode
executes semantic Run. Callers need no separate schema command.

## 22. Verify-only and full CLI behavior

Both modes now reject schema-invalid artifacts before semantic execution.
Command integration tests use temporary copies and repin deliberately mutated
case bytes so checksum/profile failures cannot mask the schema result. They
cover valid corpus, invalid manifest, absent case t0, unsupported schema source,
argument/load/profile refusal, and a schema-valid wrong TOTP answer: verify-only
passes it while full mode fails semantic execution. The CLI's tested run seam
avoids weakening W3A's subprocess prohibition; main still provides process exits.

## 23. Differential method and results

Before deleting any replaced Python file, an outside Go adapter was built against
the real schema package and removed from the checkout. Original Python bytes
were saved outside first. The temporary Python harness sends JSON requests to
that Go binary and compares accept/reject with the original validator. Every
original schema test runs while its root validation calls are recorded. Current
manifest/all 92 cases, deduplicated structure/operation mutations, and every
constraint-bearing fragment supply further rows. No permanent test depends on
Python or this adapter.

Final result: **1,726 schema comparisons: 159 accepted, 1,567 rejected, zero
drift**. Spec comparison invokes original Python and new Go CLIs side by side
on temporary two-file trees: **138 comparisons, zero drift**. Path: **4/4**
rejections agree. Evidence: `differential-3.log`, `schema-requests.jsonl`,
`schema-results.txt`, `schema-labels.json`, `schema-drift.json` (empty),
`spec-results.json`, `path-results.json`, all under the outside evidence path.
An earlier successful pass covered 134 spec rows; four text boundaries were
added before the final 138-row pass. No corpus/schema edit was needed.

## 24. Expanded mutation coverage

Rows cover representative required top-level and nested removals for each
populated structure and selected operation; false/zero/empty values and omissions;
scalar/object/array nulls; wrong-case keys; unknown nested/top-level properties;
escaped/top/nested duplicates; bool-as-integer, number-as-string, fraction/exponent
tokens; lower/upper range violations; string limits and all seven patterns;
enum/const mismatches; item type and array bounds; unknown/mixed operations;
missing operation payloads; incompatible payloads and conditional violations;
duplicate schema source and unsupported visited keyword; and path escapes.
Permanent Go tests add eager source-preflight failures and 5,625 independent
required-member removals without committing mutated corpus fixtures.

## 25. Removed Python files

After successful parity: tools/check_spec.py, tools/test_check_spec.py,
tools/test_check_vectors.py. Their Make/CI invocations are replaced. The Nix
Python fileset naturally excludes deleted files; no per-file list needed editing.

## 26. Retained independent wire oracle

Both tools/test_r16_vectors.py and tools/check_vectors.py match starting bytes.
The oracle still independently computes its own HMAC/HKDF block composition,
TLV/field encoding, VAULT hash, capacity and post-AEAD plaintext predicates in
Python, with its original limitations. It imports the local duplicate-safe read
helper; the rest of the retained Python schema checker is dormant in normal
qualification. It was deliberately not shrunk or rewritten to enable deletion.

## 27. Make changes

All public targets preserved. spec-check runs Go check-spec; verify runs the
schema-aware Go verifier -verify-only; test still runs both modules and Python
discovery. check remains architecture/spec-check/test/conformance/verify, serial.
Race and all three fuzz recipes/flags are unchanged. generate remains explicit;
normal qualification never executes it. Help now mentions the wire oracle.

## 28. CI changes

Only the old checker step label/commands change: Go check-spec and schema-aware
verify-only replace Python checker invocations; Python discovery/setup remains.
The existing later full/verify-only steps remain. Actions, matrix, Go/Python
versions, module tests/vet/format, architecture and standalone-module checks,
Linux authoritative Nix and race/fuzz coverage are unchanged. No push/dispatch.

## 29. Nix, vendor and dependencies

Only the flake comment naming the spec checker changes. Existing Go/native/module
source family already includes new policy/schema files; retained Python stays
in the fileset and runtime inputs. Make supplies the new commands. README,
CONTRIBUTING, conformance README and this report are excluded and unconsumed.
Archived pre-r16 history remains included and consumed.

No dependency/version, lock, vendorHash, module metadata, or closure change.
Generator remains stdlib-only; conformance external packages remain
x/crypto/{argon2,blake2b,hkdf} v0.36.0 and x/sys/cpu v0.31.0.
Two non-Nix `go work vendor -o <outside-path>` executions, on the starting archive
and final workspace, produced **75 identical path/byte/mode files**, including
modules.txt. Existing vendorHash retained; no unexpected vendor difference.
Root go.mod/go.sum/go.work.sum and generator/go.sum remain absent.

## 30. Documentation

README, CONTRIBUTING and conformance README explain Go prose policy, independent
checked-in schema interpretation before typed acceptance, exact required presence,
the retained Python oracle/reader and W3C handoff, and unchanged W3A ownership.
No historical report or protected vectors documentation was changed.

## 31. Generator non-involvement

All 8 generator files remain byte/path-identical to starting HEAD. Zero generator
CLI executions, zero regeneration, zero semantic or build-plumbing edits.
Ordinary generator unit tests/vet/race still run as required; compilation and
module vendoring do not execute generator main.

## 32. W3A boundary regression

Architecture guard PASS before final and within final make check. Exactly two
modules/workspace members; generator→verifier imports 0, verifier→generator 0,
shared repository modules/packages 0, local replace bridges 0, cross-role
subprocesses 0. New static policy restricts check-spec to specpolicy and stdlib,
and specpolicy to stdlib. Negative tests cover command→vectors,
policy→object, and policy→command violations. No subprocess test exception added.
Actual dependency closure still confirms the reviewed primitive versions only.
After the human Nix PASS, two architecture command invocations passed (initial
verification and logged final confirmation); neither repeated the host suite.

## 33. Final host qualification

Host Go 1.26.8, Python 3.14.7, Linux amd64. One stabilized final sequence:
`make check`, generator vet, conformance vet, `make race`, `make fuzz`: all PASS,
no incomplete final attempt. Formatting (`gofmt -l generator conformance` empty)
and `git diff --check` PASS. Architecture PASS; full/verify-only each pass 92.

Final Go test/race package results: generator cmd/generate-vectors,
internal/fixture, internal/wire all PASS. Conformance cmd/check-spec,
cmd/totipo-conformance, internal/cryptov1, graph, object, repopolicy, schema,
specpolicy, storage, tlv, totp, vectors all PASS. cmd/check-architecture compiles
with no command tests. Both vets exit 0. Logs are `final-*.log` and
`final-host-results.json` outside Git.

| Fuzz target/package | Unchanged settings | Baseline | Final |
| --- | --- | --- | --- |
| FuzzDispatch / object | -run '^$' -fuzztime 10s -parallel 2 | PASS 10.092s | PASS 10.088s |
| FuzzOpen / cryptov1 | -run '^$' -fuzztime 10s -parallel 2 | PASS 10.079s | PASS 11.011s |
| FuzzArrivalAndDisappearance / graph | -run '^$' -fuzztime 10s -parallel 2 | PASS 10.056s | PASS 10.471s |

Runtime may exceed the requested budget while workers finish. Seeds and failure
semantics unchanged; only Go's ignored fuzz cache accumulated discoveries.

## 34. Final Python tests

**3 tests, OK** in final make check: test_fields_and_key_intermediates,
test_post_aead_fixture_inputs, test_vault_id_and_maximum. All belong to unchanged
IndependentWireTests in tools/test_r16_vectors.py. No Python spec/schema checker
or tests remain in normal qualification.

## 35. Protected-byte comparison

Final exact paths/SHA-256 compare to pre-edit capture: all 100 protected files and
8 generator files identical, zero additions/deletions/differences. Retained
Python, go.work, module files/sums and flake.lock also match starting Git bytes.
Manifest IDs/count/expectations/raw hashes/profile set/pins/schema bytes and
protocol remain exactly unchanged. No expectation repair or artifact production.

## 36. Execution accounting and development failures

Counts are top-level agent invocations; guard-internal Go inspection subprocesses
are not separate qualification runs.

| Execution | Count/status |
| --- | --- |
| Baseline make check / conformance vet / make race / make fuzz | 1 / 1 / 1 / 1; all PASS |
| Baseline bounded targets Dispatch / Open / graph | 1 / 1 / 1; PASS |
| Additional baseline verbose Python discovery | 1; 15/15 PASS, names captured |
| Focused spec-policy/command Go invocations | 5; 3 PASS, 1 missing-package development failure, 1 boundary-test assertion failure |
| Focused schema(/vectors) Go invocations | 3; all PASS, first before schema tests existed |
| Focused CLI/guard Go invocation | 1; PASS |
| Differential harness attempts | 3; first harness setup failed, two parity passes (134 then 138 spec rows) |
| Original Python schema suite inside differential attempts | 3; first 10 setup errors, subsequent 10/10 PASS twice |
| Separate path differential | 1; four agreeing rejections |
| Eager schema-preflight comparison | 1; intentional strengthened source rejection measured |
| Focused fuzz outside baseline/final | 0 / 0 / 0 |
| Final complete suites / incomplete attempts | 1 / 0 |
| Final make check / generator vet / conformance vet / race / fuzz | 1 / 1 / 1 / 1 / 1; PASS |
| Final bounded targets Dispatch / Open / graph | 1 / 1 / 1; PASS |
| Remaining oracle discovery | 1 within final make check; 3/3 PASS |
| Non-Nix workspace vendor generations | 2; identical |
| Generator executions | 0 |
| Human W3B Nix | 1; PASS, human-reported |
| Agent Nix | 0 |
| Cached W3B Nix evidence | none |
| Remote CI / push / dispatch | 0 / 0 / 0; NOT YET RUN |

Development exceptions preserved in transcript/outside logs:

- Initial static keyword inventory incorrectly passed a dict to Counter.update;
  corrected to keys before deriving the actual inventory. No qualification failed.
- Temporary mechanical spec-port authoring script initially missed tuple assignment,
  then a regex flag argument. The latter left the policy package unwritten, causing
  one focused compile attempt to fail. The completed Go policy was reviewed and
  passed focused/current/differential checks; the authoring tool is outside Git.
- First differential attempt loaded original test bytes from the outside snapshot,
  whose relative BASE then pointed outside the corpus: ten setup errors. BASE was
  explicitly pointed at the unchanged checkout before the two successful passes.
  No implementation acceptance drift occurred.
- First extra Unicode reference test replaced every literal Section 12, also
  breaking an independently required scope anchor. The test was corrected to add
  a reference without destroying that anchor; the corresponding original-Python
  boundary comparison then agreed. No policy rule was weakened.

No baseline/final check, vet, race, fuzz, committed-corpus, or successful parity
failure occurred; no fuzz reproduction input was generated. Full logs remain
local temporary evidence; this report preserves durable results and counts.

## 37. Frozen-input inventory

All intended new source/test/doc/report files are made Git-known with
intent-to-add only. `git diff --cached --quiet` must remain PASS. Frozen path,
state (including explicit deletions), mode and SHA-256 inventory and Git membership
are saved outside the repository in `frozen-inputs.json`, `frozen-membership.json`
and `frozen-summary.json`. Freeze covers every Git-known file except this excluded
report: generator/conformance/remaining tools, Makefile, workspace, flake/lock,
workflow, current docs, spec/vectors/requirements and archived history. Temporary
adapter/harness, logs, vendor trees, caches and original Python copies stay outside
Git source membership. The freeze records **286 paths** (including the three explicit deletions),
**287 Git-known paths**, inventory SHA-256
`7d59025727f7c8899c9813f054de1ae90879ecc071aa87341aef2e6c1c8d6092`.
The nine new intended files (eight Go source/test files and this report) are
intent-added; cached content remains empty.

This report is excluded by the actual flake fileset and is not consumed by the
policy, schema, vector, architecture, or generator implementations; it alone may
be finalized after the human gate. All effective source/supporting inputs are
frozen before requesting Nix. No included input may change after a PASS without
requalification.

## 38. Human gate

**PASS — human-reported on 2026-10-12; human execution count 1; agent Nix
count 0.** The user reported “nix flake check passed” for the requested command
from the repository root after the recorded freeze:

```sh
nix flake check
```

No path: source, rebuild request, package materialization or agent Nix execution.
Cache reuse versus fresh derivation execution was not reported and is not inferred.

Post-PASS verification checked all **286 frozen path/state/mode/hash records**,
**287 Git-known membership paths**, all **100 protected artifact files** and
**8 generator files**, unchanged main/starting HEAD, no unrelated untracked files,
and no staged content: PASS. Architecture guard: PASS, retaining two independent
modules and generator stdlib-only closure. Only this excluded report was updated
afterward; no included or supporting frozen input changed. The freeze and index
checks passed again after report finalization. Post-Nix evidence is recorded
outside Git in `post-nix-freeze-check.log`, `post-nix-architecture.log`, and
`post-nix-final-check.log`.

## 39. Remote CI

**NOT YET RUN** for this uncommitted milestone. No push or workflow dispatch.
Human commit/push will obtain authoritative Linux and supplemental macOS/Windows/
Go portability evidence. Local Linux host results do not establish remote passes.

## 40. Limitations and W3C handoff

Human Nix PASS is user-reported; cache/fresh-build details are unknown.
Post-gate frozen-input and architecture verification passed. Remote CI remains
unrun. The schema interpreter
supports the current closed subset/patterns, with eager source validation; future
extensions require explicit review/tests. Temporary parity machinery is not a
permanent qualification dependency. CLI tests exercise the run/exit-code seam;
main remains a direct os.Exit wrapper. The independent Python oracle's original
AEAD/Argon2/coverage limitations are unchanged. W3C alone owns its eventual
migration and removal of the remaining reader/runtime island. No new oracle,
independent crypto implementation, external interoperability or zero-Python claim.

## 41. Final Git and index state

Branch/HEAD stay main / ac6a1489cc4b557155b9abab2f9996ad9d4a5edb. Changes remain
uncommitted. New files are intent-added only; existing edits/deletions unstaged.
No staged file contents, unrelated untracked files, generated corpus, temporary
harness, vendor output, or evidence is included. Final membership/status and
empty cached diff are recorded with the freeze outside Git. No reset, commit,
push, tag or dispatch. Post-Nix verification confirms the same state; only this
excluded report was finalized after the gate.

## 42. Suggested commit subject

`Move spec and schema validation into Go`
