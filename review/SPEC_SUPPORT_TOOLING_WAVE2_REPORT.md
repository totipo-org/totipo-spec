# Spec support-tooling consolidation Wave 2

## 1. Starting state

Date: 2026-10-11. Repository: `/home/niki/Sources/totipo-spec` only.
Starting branch: `main`; HEAD: `09601efe1a410598fb124beb0d30b3be231a8e79`;
`git status --short`: empty. The audit is committed in `0591c59`; Wave 1 is
committed in starting HEAD `09601ef`. Protocol/profile/manifest: **v1/r19**, moving
pre-RC, **92 required cases / 92 physical case JSON files**. No revision advance.

Read AGENTS, CONTRIBUTING, README, the support-tooling audit and Wave 1 report,
all five active Python files, all vectors-package Go files, the consumer CLI,
Makefile, flake.nix, schemas and portability workflow. No generator execution.
Inherited Wave 1 evidence records passing host qualification and one
human-reported Nix pass. Those are prior-tree observations, not this wave's gate.
Remote CI for this wave has not run.

## 2. Classification and scope

**TEST_OR_QUALIFICATION_TOOLING**, with AGENTS classifications
**PYTHON_VERIFIER** and **CONFORMANCE_GO**. **BUILD_NIX_CI** applies to effective
qualification-source invalidation: included Go/Python bytes change, although
Makefile, Nix files, lock, dependency declarations and CI configuration do not.
There is no protocol, corpus, schema, requirement, expected-result or generator
change. Python remains required. All work is unstaged/uncommitted.

## 3. Verification responsibilities before and after

| Responsibility | Before | After |
| --- | --- | --- |
| Four raw artifact pins | Python check_spec and Go VerifyProfile | Go VerifyProfile only |
| Manifest/profile exact revision and required set | Go Read/VerifyProfile, repeated revision literals | Same Go authority; explicit ExpectedRevision shared within reader and CLI |
| Spec sections/history/policy anchors | Python check_spec | Same Python rules, unchanged |
| Language-neutral shape contract and required presence | JSON schemas interpreted by Python | Same schemas and interpreter |
| Strict consumer decoding and execution prerequisites | Go uniqueJSON, Decode, ValidateShape | Same logic, expanded negative tests |
| Physical completeness, paths, case identity and hashes | Go Read | Same logic, expanded mutation tests |
| Independent wire/key formulas | Python test_r16_vectors | Unchanged |
| Expected-result execution | Go Run, TestCorpus, full CLI | Unchanged, all 92 cases |

The current responsibility model is documented in
[conformance/README.md](../conformance/README.md). No shared writer/verifier
entry point, generated expectation, or schema port was introduced.

## 4. Artifact-pin authority and r18 regression

`vectors.Read` validates the manifest's expected revision, unique identities and
paths, physical coverage, raw case hashes, case identity/outcome and payload
shape. `vectors.VerifyProfile`, supplied that validated manifest, validates the
profile format/status/protocol/revision, hashes the exact spec, manifest and both
schema files, and requires all-and-only manifest case IDs with matching hashes.
Together they are the single read-only integrity authority.

Removed only the Python profile JSON read, four raw hash comparisons, unused
hashlib/json imports and hashing claim in its success message. `make spec-check`
now explicitly checks structure/history; `make verify` remains the standalone
schema plus integrity gate. Editorial byte drift still fails `make check` and
`make verify`, including the actual r18 stale-spec-pin regression. A valid
structure check alone intentionally makes no integrity claim.

`TestArtifactPins` accepts current bytes and rejects a trailing newline in each
of the four artifacts, and separately missing/stale values for each pin.
`TestProfileRelationship` rejects wrong revision, missing/extra/duplicate/unknown
required IDs and wrong case hashes. Expectations come from the authored profile;
verification never rewrites it or computes replacement expected pins. Test-only
case hash changes isolate identity/shape rejection in temporary copies; they do
not touch committed bytes or run generation.

## 5. Complete structural-rule disposition

Classes: **S** semantic/live structural invariant; **P** useful process guard;
**B** brittle exact prose fingerprint; **D** duplicate integrity check. B describes
the current implementation of an invariant, not permission to drop it. All rows
except the final D row are retained. A prefix comparison against starting HEAD
proved the complete pre-pin checker source identical after removing the two
unused imports. No structural assertion, regex, section split or anchor changed.

| Existing rule(s) | Class | Disposition |
| --- | --- | --- |
| Split active text/history at section 21; exactly consecutive sections 1–21 | S | Retain |
| Explicit protocol version 1 and revision r19 header | S/B | Retain independent spec assertion |
| Archived-history link; r1–r15 archived and r16–r19 current entries | P/B | Retain continuity and exact link checks |
| Active retired terms: DEVICE, DEVICE_ID, P-256, ECDSA, DER, SIGNATURE, AUTHOR_DEVICE_ID, AUTHOR_TIME, PROVENANCE, OPAQUE_ROUTABLE, OPAQUE_UNSCOPED, OBJECT_TYPE, OBJECT_VERSION, VAULT_BINDING, READY, PROCESSING_INCOMPLETE, HISTORY_REGRESSION, HISTORY_MEMORY_LOST, candidate-use, K_signature_context | P | Retain active-text token-boundary exclusions |
| Case-insensitive provenance and opaque exclusions | P | Retain |
| Sole TOKEN grammar; MAX_PARENTS=4 and parent-count bounds; 1005/1006/1024 anchors | S/B | Retain |
| Exact VAULT_ID formula and causal-equivalence formula | S/B | Retain |
| Create-only state, optional external synchronization (including current newline), non-exhaustive enumeration, no parent-publication prerequisite | S/B | Retain |
| Ordered registry tags 0x0001–0x000c and all 12 field names | S | Retain |
| Every active numeric Section reference resolves to 1–20 | S | Retain |
| Storage-layout code block begins with vault; no uppercase VAULT entry or pathname claim | S/B | Retain |
| Fold table has exactly TokenValue, CLIENT_NAME, CLIENT_TIME; exact-presence rule for both metadata fields | S/B | Retain |
| TokenValue tuple excludes metadata; head section contains value and both metadata fields; no head metadata erasure | S/B | Retain |
| Store table: vault regular file only, objects-v1 directory only | S | Retain |
| Symlinks must not deliberately be followed/interpreted/traversed | S/B | Retain |
| No openat2, O_NOFOLLOW, O_PATH, statx or SecureDirectoryStream prescription; optional hardening anchor | P/B | Retain portability guard |
| Metadata preserved while represented/returned/exposed/retained; no persistent retention or advisory-history/cache/cross-run requirement | S/B | Retain |
| History unavailability alone does not prohibit informed authorship; new CSPRNG TOKEN_ID and empty parent set | S/B | Retain |
| Empty-password interactive confirmation before creation | S/B | Retain |
| Creation cannot silently initialize; explicit recovery/reconfiguration/new-location workflow; orphan data not garbage or deleted | S/B | Retain |
| Reopen canonical vault regular file; NO-REPLACE; success only after publication/persistence/revalidation | S/B | Retain |
| Tombstoning/deletion cannot be described as secure erasure | S/B | Retain |
| No in-place password change/rewrap; no root reuse for migration; INFORMATIVE and possibly incomplete synchronization anchors | S/B | Retain |
| Retired VAULT_FINGERPRINT, old fingerprint label, CURRENT != BASE, stageReplacement, compare-before-replace, latest password wrapper | P | Retain exclusions |
| CSPRNG 32-byte root, lifetime immutability, no in-place rotation, bootstrap 0x01, never newer version, same-root alternate-bootstrap anchor | S/B | Retain |
| No replacement even with exact-identical vault bytes | S/B | Retain |
| Core protocol, Store/writer, Application scope rows; core metadata coverage | S | Retain |
| Section 12 per-object versus Section 17 fold metadata, applicable core requirements, GUI boundary, corpus/application distinction, no stronger atomic primitive | S/B | Retain |
| Four exact raw artifact pins including editorial history | D | Remove from Python; surviving Go authority and regressions above |

Exact prose fingerprints remain candidates for a separately reviewed structural
rewrite in Wave 3. This wave deliberately makes no claim that substring checks
prove protocol behavior. New Python mutations cover section ordering, revision,
reference bounds, current/archive history entries, archive link, registry tag and
storage entry type. Passing structure tests use only spec/archive files.

## 6. Schema/Go shape-protection matrix

The required contract is the **union**, not equivalence of schema and typed Go
structs. Neither schema file nor either JSON parser was removed or changed.

| Protection | Schema / Python | Go | Wave 2 disposition |
| --- | --- | --- | --- |
| Duplicate members at any depth | pairs hook, including schemas | uniqueJSON before typed Decode | Keep both reader boundaries |
| Unknown nested members / mixed operation payloads | additionalProperties:false and operation conditionals | DisallowUnknownFields; ValidateShape counts payloads | Keep; nested Go negative expanded |
| Required false booleans | required + boolean, including graph expect, storage diagnostics and workflow flags | Plain bool cannot distinguish omitted from false | Python remains authority for presence; add all five workflow flag omission tests |
| Required zero values | required + integer, including TOTP t0 and row time/counter | uint64 accepts zero; missing value may also become zero | Python retains presence; add t0 omission plus rows/period omissions |
| Exact integer/u64 behavior | Python integers, explicit type excluding bool, bounds through u64 max | Typed uint64 rejects overflow, negative, fraction and bool | Keep; Go explicitly tests 0/max acceptance and four invalid types/ranges |
| String lengths, patterns, enums and array limits | Schema-subset validate | Some typed and execution checks, not complete schema interpretation | Keep Python; no interpreter migration |
| Operation-specific fields | All nine conditional schema branches and nested required lists | ValidateShape enforces payload, root/input/semantic prerequisites, kind and selected lengths/results | Keep both complementary protections |
| Safe paths | Resolved listed paths must remain within vectors | Canonical local cases/ path, no backslash, symlink containment | Keep; new Go traversal/backslash negatives |
| Manifest/case identity and expectation | Schema constrains each separately | Read cross-checks exact identity/expected result | Go authority; new identity mutation bypasses hash guard only in temp test |
| Unique manifest IDs/paths, exact physical corpus | Not supplied by schemas alone | Read maps plus WalkDir | Keep; duplicate ID/path and unlisted-case negatives |
| Profile equality and raw hashes | Not schema evaluator responsibility | Read/VerifyProfile | Single Go authority |

Operation contract inventory: crypto requires input/semantic/root/crypto;
dispatch semantic/root/crypto; bootstrap root/bootstrap; totp totp; graph graph;
fold fold; storage root/storage; workflow workflow; post-aead semantic/root/
post_aead, in addition to format/id/operation/expected. Nested schemas retain all
crypto intermediates, bootstrap fields, TOTP source/notes/algorithm/digits/period/
t0/secret/rows, graph action-specific add/remove/evaluate fields and expected
snapshots, fold frontier/stage IDs/parents/token, storage namespace/entries/classes/
diagnostics, workflow action/kind/bytes/five booleans/result, and authenticated
negative object ID/plaintext/object/defect. Their requirements are unchanged.

No accidental shape parsing was removed: there is no demonstrated whole-parser
equivalence, especially for false/zero presence. The removed duplicate profile
parse was demonstrably redundant with Go. Keeping parser boundaries also protects
the independently usable Go consumer and Python wire oracle.

## 7. CLI/orchestration repetition

**Deferred.** `make check` still runs spec-check, tests, full conformance, then
verify (Python schemas and Go -verify-only). Portability CI still invokes both
public modes. Full CLI and TestCorpus still execute all 92 cases. Normal, race
and fuzz remain distinct qualification evidence.

Full and verify-only share Read/VerifyProfile, so their integrity computation
repeats. However the command package currently has no dedicated CLI behavior
tests; removing a mode invocation would remove its existing smoke evidence.
This wave does not substitute library tests for argument/exit/skip-execution
behavior. A future orchestration change should first add focused CLI tests for
valid verify-only, loader/profile failure, positional/flag failure, and full-mode
evaluator failure versus verify-only skipping evaluation. Make targets and the
five-job supplemental portability matrix remain byte-identical.

## 8. Revision-literal inventory

| Location | Responsibility / decision |
| --- | --- |
| Spec header/history and current prose | Independent authored r19 assertion; unchanged |
| requirements/v1-pre-rc.json and vectors/manifest.json | Authored machine-readable r19 metadata, validated; unchanged |
| vectors/manifest.schema.json | Independent r19 const; unchanged |
| vectors/case.schema.json title | Historical r16 descriptive label, not current revision validation; unchanged |
| Go Read, VerifyProfile, CLI summaries | Four current literals consolidated to explicit vectors.ExpectedRevision = r19 |
| Generator manifest/profile/output literals | Retain independent r19 writer declarations; generator does not set verifier expectation |
| Python check_spec header, output, history upper bound | Retain independent spec/version/history assertions |
| Python test/Go test r19 inputs | Independent mutations and expectations; retained |
| Make help, Nix derivation versions/output, README/CONTRIBUTING/conformance/requirements docs | Descriptive labels; no auto-derivation or new synchronization machinery |
| test_r16_vectors name/docstring, object package comment, old history | Historical grammar labels; no revision authority |

Go expected revision is a literal verifier assertion, never read from generated
artifacts. Schema and spec assertions remain independent. CLI output text is
unchanged. This removes three redundant Go literals without letting a generator
or incoming manifest choose what the verifier accepts.

## 9. Independence retained

Generator versus verifier; authored expected outcomes versus Run; requirements
expected set versus writer case list; Python HMAC/HKDF/TLV/wire formulas versus
production helpers; graph closure oracle versus traversal; genuinely authenticated
negative-fixture proof versus generic failure; known RFC answers; normal tests
versus all-package race versus three bounded fuzz targets; fixed Linux Nix versus
supplemental portability. No independent wire-oracle source changed. No Python
elimination, evaluator rewrite, expected-result update or generator command.

## 10. Test-first differential and moved/added tests

Baseline `make check` passed against unchanged production inputs: **14 Python
tests**, seven Go test packages, full **92-case** execution, schema validation,
and verify-only success. Existing negatives captured duplicate/unknown JSON,
required false, u64 limits, revision, operation mixing and authenticated defects.

Before responsibility edits, added and passed Go artifact/profile/corpus mutations
against the original verifier. The old Python pin regressions also passed before
delegation. Added structure mutations passed against the original checker before
removing its pin block. The new zero-presence schema test was first run against
the unchanged schema interpreter. Its initial filename was incorrectly written
as v1.totp.sha1.001.json; that focused run failed with FileNotFoundError. Corrected
to the existing v1.totp.rfc6238-sha1.001.json, reran the focused group, and passed.
Failure log is preserved; no fixture or expectation was changed to resolve it.

Moved Python pin behaviors to `integrity_test.go`; removed Python pin-harness
coupling. The current-pass assertion is supplied by Go integrity tests and the
separate two-file Python structure test. Added profile set/revision negatives,
unsafe/backslash paths, duplicate IDs/paths, identity mismatch, incorrect shape,
raw case drift and unlisted physical cases. Expanded Go strict decoding and
added explicit u64 zero/max acceptance. Python adds structure mutations, all five
workflow boolean-presence checks and required TOTP zero/operation-field checks.

Focused Python invocations: **4** (one failed, three passed), logging 12 tests
with one error, then 12 pass, 7 pass during temporary pre-consolidation test
composition, and 12 pass after consolidation. Focused Go invocations: **3**, all
passed (0.096s, 0.078s, 0.082s). No focused fuzz runs. Raw logs are local temporary
evidence in `/tmp/totipo-wave2-qualification/`, including both differential logs,
retry, structure-before and focused logs. Tests and source, not temporary logs,
are the durable regression evidence.

## 11. Protected-byte result

Compared starting Git blobs to working-tree raw bytes and path inventories,
including additions/deletions: **spec 1 file, vectors 97 files (including all
JSON schemas/cases/manifest), requirements 2 files: identical, no additions or
deletions**. All 32 existing review files/history reports also remain identical.
Generator, independent Python oracle and schema interpreter are identical.

Makefile, flake.nix, flake.lock, all .github inputs, go.work and every go.mod/go.sum
remain identical. Root go.mod/go.sum are absent at both boundaries. Changes under
conformance/tools are only those listed below. No vector regeneration, case
rewrite, authored expected-result change or pin update occurred.

## 12. Final host qualification

One stable final suite, all exit 0. `make check`: structural PASS; **15 Python
tests in 0.474s, OK**; **92 cases** executed by full CLI; schemas for **92 cases**
PASS; manifest/profile/**92 case files** verified by verify-only PASS. Vet passed
without diagnostics. Formatting and `git diff --check` passed.

| Package under totipo/conformance | Final normal tests | Final all-package race |
| --- | --- | --- |
| cmd/generate-vectors | no test files | no test files |
| cmd/totipo-conformance | no test files | no test files |
| internal/cryptov1 | ok 0.933s | ok 1.842s |
| internal/graph | ok 0.012s | ok 1.044s |
| internal/object | ok 0.003s | ok 1.011s |
| internal/storage | ok 0.745s | ok 1.737s |
| internal/tlv | ok 0.002s | ok 1.010s |
| internal/totp | ok 0.003s | ok 1.011s |
| internal/vectors | ok 1.576s | ok 3.602s |

Normal and race retain `-count=1 ./...`; race retains `-race`. Logs:
`check.log`, `vet.log`, `race.log`, `fuzz.log` in the temporary directory above.
No failing final attempt. Afterward only this excluded report was added.

## 13. Fuzz and execution accounting

Each final fuzz target used exactly `-run '^$' -fuzztime 10s -parallel 2` with its
existing target name/package. No seed source or failure semantics changed.

| Target / package | Initial coverage | Executions | New interesting / total | Duration / result |
| --- | --- | ---: | --- | --- |
| FuzzDispatch / internal/object | 46/46 | 924,554 | 4 / 50 | 10.046s PASS |
| FuzzOpen / internal/cryptov1 | 16/16 | 772,371 | 0 / 16 | 11.013s PASS |
| FuzzArrivalAndDisappearance / internal/graph | 126/126 | 53,657 | 4 / 130 | 10.858s PASS |

Reported elapsed time can exceed 10s during final worker shutdown; bounded policy
is unchanged. No fuzz failure or reproduction input was produced.

| Execution | This wave |
| --- | --- |
| Baseline make check | 1 PASS |
| Focused Python / Go invocations | 4 (1 failed, 3 passed) / 3 PASS |
| Focused fuzz Dispatch / Open / graph | 0 / 0 / 0 |
| Final complete suites / incomplete attempts | 1 / 0 |
| Final make check / vet / race / fuzz | 1 / 1 / 1 / 1 PASS |
| Each fuzz target within final make fuzz | 1 |
| Generator executions | 0 |
| Human Nix / agent Nix | 1 / 0; human-reported PASS / not run |
| Fresh remote CI / cached Nix observations | 0 / none |

Static byte/source/format/link checks are not unit-test invocations. Inherited
Wave 1 host/human Nix evidence is kept separate from these execution counts.

## 14. Nix source filter and human gate

The actual filter includes spec/vectors/requirements, conformance Go/module files,
tools Python, Makefile, go.work and archived pre-r16 history. Changed Go/Python
files, including the new tests, invalidate the inherited Nix result. The filter,
fixed vendorHash/lock and Linux/supplemental portability boundary are unchanged.
README and this report are excluded; conformance README is excluded by the Go
file filter. Neither is consumed by structural checks.

All qualification-source inputs are frozen. **Human Nix: PASS — human-reported
after handoff.** The user reported “nix flake check passed” for the frozen Wave 2
inputs. Human execution count: 1; agent execution count: 0. This records the
human’s result, not an independently observed agent run. AGENTS.md says “The
agent does not run Nix” and requires the final human gate after source freeze.
The handed-off command was:

```sh
nix flake check path:.
```

No agent Nix command, artificial package materialization or cached-pass claim.
This status/count update changes only this excluded report; no qualification-source
input was edited during the update. A later included-source edit invalidates the
human result. Remote Linux and
supplemental portability CI remain post-commit/push evidence, not locally claimed.

## 15. LOC and tool-call reduction

Physical executable source/test lines, including comments/blanks:

| Scope | Before | After | Delta |
| --- | ---: | ---: | ---: |
| Python, five files | 493 | 486 | -7 |
| Go, including generator and tests | 3,044 | 3,242 | +198 |
| Make / Nix / workflow | unchanged | unchanged | 0 |
| Total executable/orchestration delta | | | +191 |

Production checker shrinks 16 lines; explicit Go expected revision adds 3 lines:
net **13 fewer production verification lines**. New/moved tests account for the
net overall increase; this is responsibility consolidation and better negative
evidence, not a claimed net LOC pruning. No runtime, language or public tool is
removed. Documentation/report LOC is excluded from this implementation count.

Each check_spec invocation removes one profile parse and four artifact hashes.
Python structure-test subprocesses fall from 13 (current + 4 drift + 8 pins) to
9 (current + 7 structure mutations + archived-entry mutation): **4 fewer
subprocesses** per discovery. Including the standalone spec check, 14 becomes
10 such invocations per make check. The moved Go tests run in-process. No Make/CI
CLI invocation reduction is claimed; both Go run commands remain.

## 16. Precise remaining Python responsibilities for Wave 3

- `check_spec.py` (115 lines): current revision/version assertions; section,
  reference, registry, history/archive, active retired-term and storage/policy/
  metadata/scope anchors listed above. No hashes or profile parsing remain.
- `test_check_spec.py` (56 lines): standalone structure success with only two
  inputs, and eight mutation failures. Preserve these behaviors when migrating.
- `check_vectors.py` (101 lines): duplicate-safe JSON reading; checked-in schema
  subset interpretation (types/const/enum/properties/required/additionalProperties,
  items/array limits, integer bounds, patterns/string limits, allOf/if/then/else),
  visited unknown-keyword rejection; manifest-listed case validation and resolved
  path containment. It is not a general schema implementation or full meta-schema
  validator. Go struct decoding alone cannot replace it.
- `test_check_vectors.py` (111 lines): ten test methods protecting existing shape
  negatives plus required workflow false flags and TOTP zero/field presence.
- `test_r16_vectors.py` (103 lines): unchanged three-test independent stdlib
  HMAC/HKDF/key/ID/nonce/AAD/padding/TLV oracle, VAULT hash and maximum size, and
  post-AEAD plaintext-defect construction evidence. It still does not claim
  independent AES/Argon2 validation. Port it last without production-helper reuse.

Wave 3 must preserve exact presence/type behavior and strict duplicate handling
before removing either parser. Review Go JSON field matching/default/null behavior
against schemas rather than assuming typed decoding is schema equivalence.

## 17. Final Git state

Branch `main`; HEAD remains `09601efe1a410598fb124beb0d30b3be231a8e79`.
Index unchanged; no staging, commit, push, tag, release or CI dispatch.

```text
 M conformance/README.md
 M conformance/cmd/totipo-conformance/main.go
 M conformance/internal/vectors/runner.go
 M conformance/internal/vectors/runner_test.go
 M tools/check_spec.py
 M tools/test_check_spec.py
 M tools/test_check_vectors.py
?? conformance/internal/vectors/integrity_test.go
?? review/SPEC_SUPPORT_TOOLING_WAVE2_REPORT.md
```

No stop condition was encountered. Host work is complete and the human reported
passing pinned-environment qualification. Historical reports are preserved verbatim.
