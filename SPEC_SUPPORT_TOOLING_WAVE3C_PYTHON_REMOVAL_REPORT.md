# Wave 3C — Python removal and verifier cutover

## Starting state and scope

Date: 2026-10-12. Repository only: `/home/niki/Sources/totipo-spec`.
Starting branch **main**, HEAD **df820576c5076b9c9434a4f4d64a6b7c316b24e4**.
Before edits, `git status --short` was empty; the required HEAD and branch
commands confirmed the above. Protocol remains **v1/r19**, moving pre-RC.
No sibling repository or outside agent state was inspected or used.

Classifications: **CONFORMANCE_GO**, **PYTHON_VERIFIER** (oracle/checker removal),
**BUILD_NIX_CI**, **DOCS_ONLY** (current guidance). No SPEC_TEXT, REQUIREMENTS,
VECTOR_CORPUS, GENERATOR, or semantic change. No verifier acceptance fix was
needed. Historical reports are untouched. No commit, tag, release, push, reset,
or workflow dispatch was performed.

**W3C complete: cutover, host qualification, human-reported Nix/development
qualification and post-gate frozen-input verification PASS.**

## W3A/W3B assumptions carried forward

Read the architecture audit, W3A boundary report and W3B structural/schema report.
W3A and W3B record completed host qualification and human-reported Nix PASS;
these are inherited evidence for their inputs, not new W3C passes. Cache versus
fresh derivation execution is not known. Remote portability passes are not inferred.
W3B reports 1,726 schema and 138 spec-policy differential comparisons with zero
disposition drift, four agreed escaping-path rejections, and no cutover blocker.
Its eager schema-source preflight is an explicitly documented strengthening.
Existing Python const/enum, terminal-newline pattern and Unicode policy behavior
remains deliberately implemented in Go; no compatibility rules were changed here.

W3B already made structural/schema verification authoritative in Go. Its exact
remaining runtime island was `tools/test_r16_vectors.py` (three independent wire
methods) importing the duplicate-safe `read` helper in `tools/check_vectors.py`.
The latter's schema CLI was dormant. `tools/check_spec.py` and both structural/
schema Python test modules were already deleted by W3B, not by this wave.

The architecture audit §§17, 23–24 explicitly assigns W3C the independent
stdlib-only `conformance/internal/wireoracle` tests, guarded against production
imports. This preserves existing oracle evidence; it is not another schema or
semantic verifier implementation wave. Two sibling Go modules, independently
owned models/constants/composition, serialized artifacts as the interface,
stdlib-only generator, and the existing primitive dependency allowlist remain.

## Exact changes and entry points

| File | Change |
| --- | --- |
| `Makefile` | Remove PYTHON variable/discovery; describe Go oracle in help |
| `flake.nix` | Remove tools Python fileset, Python native input/version/environment setting, dev-shell and both jail package entries |
| `.github/workflows/conformance.yml` | Remove setup-python and discovery; rename policy/schema step |
| `README.md` | Remove tools directory/prerequisite/handoff claims; describe Go-only support tooling |
| `CONTRIBUTING.md` | Go focused loop and guarded independent oracle guidance |
| `conformance/README.md` | Describe independent Go wire tests and limits |
| `AGENTS.md` | Update live example, source-filter description and post-cutover accounting guidance |
| `conformance/internal/repopolicy/policy.go` | Enforce stdlib-only oracle imports, including tests/inactive source |
| `conformance/internal/repopolicy/policy_test.go` | Four negative oracle import fixtures: verifier, schema, generator, external primitive |
| `conformance/internal/wireoracle/oracle_test.go` | New independent reader, three migrated oracle families, reader and mutation regressions |
| `tools/check_vectors.py` | Delete obsolete checker/reader (101 lines) |
| `tools/test_r16_vectors.py` | Delete migrated oracle (103 lines) |
| `SPEC_SUPPORT_TOOLING_WAVE3C_PYTHON_REMOVAL_REPORT.md` | New report at the requested repository-root path |

The empty tools directory and its generated local bytecode cache were removed.
No third-party dependencies, module declarations, lock, vendorHash, generator
implementation or output changed. The CI matrix, pinned Actions other than the
removed Python setup, Linux authority and supplemental portability roles remain.

| Previous path | Sole live replacement after W3C |
| --- | --- |
| Python spec checker, already retired in W3B | `make spec-check` → `go run ./conformance/cmd/check-spec -root .` |
| Python schema CLI, dormant since W3B | `make verify` → `go run ./conformance/cmd/totipo-conformance -root . -verify-only` |
| Python oracle discovery in make test/CI | `go -C conformance test -count=1 ./...`, including wireoracle |
| Python package requirements in qualification/development/jails | Existing Go and Make tooling |

`make check` still runs architecture, spec policy, both modules' tests, full
conformance and verify-only. `make conformance` still executes every case;
`make verify` still checks serialized schemas/artifacts/pins without semantic
execution. No normal check invokes `make generate` or generator main.

## Migration parity and deletion gate

Structural/schema deletion evidence is inherited from W3B, whose report records
the full old/new differential method and results. W3C did not rerun deleted
Python source or reimplement those checks. Coverage of the requested boundaries:

| Boundary | Evidence |
| --- | --- |
| Pristine repository | W3B parity; W3C baseline check: Python 3/3 OK, Go policy/full/verify PASS, 92 cases |
| Malformed JSON | W3B duplicate-safe parser/schema comparisons; W3C oracle comparison rejects malformed and duplicate-member cases |
| Schema-invalid manifest/case, missing required material | W3B differential mutations, 5,625 required-member removals, repinned CLI negatives; unchanged Go regressions pass at W3C baseline/final |
| Manifest/physical-case mismatch, missing file | Existing `TestManifestContract/unlisted` and `TestResolvedPathMissing` pass at baseline/final |
| Incorrect case digest/hash | Existing `TestManifestContract/case-bytes` passes at baseline/final |
| Incorrect requirements pins/checksums | Existing `TestArtifactPins`, `TestProfileRelationship/hash` pass at baseline/final |
| Structural/revision inconsistency | W3B 138 spec CLI comparisons; current specpolicy tests, `TestManifestTampering` and profile revision regressions pass |

Hashes, physical inventory and requirements pins were already Go responsibilities
in the old combined qualification path. The old Python schema CLI alone did not
check them. No false claim of Python-only coverage or new differential rerun is
made for these unchanged responsibilities.

Before Python deletion, W3C built the real oracle test binary and compared it
with all three original Python methods on temporary copies of vectors. Both
read serialized files; neither generated fixtures. Original Python bytes were
saved before edits; temporary copies were automatically removed after each row.
The original checkout remained untouched by mutations.

**22 comparisons: one accepted pristine corpus, 21 rejected mutations, zero
disposition drift.** Rejections: eight crypto fields (id key, object-root key,
object ID, object key, nonce, AAD, padding, object size); five TOKEN inputs
(identity, status, period, client name, client time); wrong VAULT_ID; maximum
length and padding; repaired nonzero-padding, object-ID-mismatch and invalid-
semantic-length defects; malformed JSON; duplicate members. Python executed
three tests per row, 66 method executions total. Go selected exactly the three
migrated test families per row. Diagnostics were not compared.

Go retains direct HMAC extract/single-block expand with independent domain
literals, raw TLV walk and independently assembled expected fields/widths,
VAULT record SHA-256, 1005-byte maximum and zero padding, 1008/1024-byte sizes,
and all three original post-AEAD input predicates. An independently summed
maximum-field-size assertion supplements the migrated maximum check, as the
audit requests. Reader regressions reject nested/escaped duplicates, malformed
JSON, trailing values and non-object roots. Sixteen permanent mutation subtests
exercise crypto, TOKEN inputs and repaired post-AEAD defects. Four new guard
negatives protect the oracle imports. No new independent AES-GCM, Argon2,
full interoperability, graph oracle or protocol rule is claimed.

Local evidence: `/tmp/totipo-w3c-df820576/` contains baseline/final logs,
original-tools, parity harness/log/results, dependency listings, hashes and
vendor comparison. These are local supplementary evidence; the report records
the durable results and counts. Permanent qualification has no Python dependency.

## Isolation proof

```text
go.work: exactly ./generator and ./conformance; Go 1.23.0
  totipo/generator
    cmd/generate-vectors, internal/fixture, internal/wire
    repository dependencies: own module only; external packages: none
  totipo/conformance
    cmd/check-architecture, cmd/check-spec, cmd/totipo-conformance
    internal/repopolicy, specpolicy, schema, vectors
    internal/tlv, object, cryptov1, graph, storage, totp
    internal/wireoracle (test-only; independent stdlib reader/formulas)
```

Mechanical inspection used `go -C <module> list -deps -test` for both modules,
plus the oracle alone; `make architecture` checked complete source, module/
workspace metadata and resolved test dependencies. Generator's nonstdlib closure
contains only its own three packages and test variants. Conformance contains
only its own packages and the four authorized primitives:
x/crypto/{argon2,blake2b,hkdf} v0.36.0 and x/sys/cpu v0.31.0.
Oracle's nonstdlib closure contains only itself and its synthesized Go test
variants. There are no generator → verifier or verifier → generator imports,
shared Totipo implementation packages, root/third module, replace bridge, or
cross-role process delegation. Standalone GOWORK=off tests also pass for both.

`go list -m all` shows the two workspace main modules and the existing external
module graph; graph-only x/net/x/term/x/text entries are not imported packages or
new dependencies. That inspection generated an untracked go.work.sum containing
three go.mod sums; it was saved in local evidence and removed. It was never
Git-known or an intended input. An inventory check caught its presence; no
qualification suite failed and no source or dependency declaration changed.
Root go.mod/go.sum/go.work.sum and generator/go.sum remain absent.

Intentionally common material is serialized spec/schemas/corpus/profile and
standard cryptographic primitives already authorized by W3A. Models, protocol
constants and Totipo composition remain separately owned. The oracle's own
reader uses raw maps/json.Number, not either module's models or schema helpers.
The guard checks all oracle imports and also rejects process/native/linkname
bridges under the existing rules. This establishes code/dependency separation,
not a mathematical proof of independent correctness or protection from a
malicious rewrite/copy of implementation or guards.

> Can a bug in the generator's implementation automatically make the verifier
> accept the same bad generated artifact merely because they execute the same
> implementation code?

**No.** They do not execute shared Totipo implementation code. Correlated human
mistakes and shared primitive defects remain possible; serialized-artifact
verification, independently authored expectations and selected independent
oracle formulas provide separate evidence.

## Byte invariance and dependencies

Before editing, sorted SHA-256/path inventories captured spec, vectors,
requirements and generator. Final enumeration and hashes match exactly,
including additions/deletions: **100 protected files** (1 spec, 97 vectors,
2 requirements) and **8 generator files**. `git diff --quiet HEAD --` those trees
also passes. Manifest/full/verify outputs confirm **92 cases**, physical cases
and required entries unchanged. No regeneration or generator CLI execution.

Protected inventory SHA-256 (sorted `sha256sum` lines with repository paths):
`36ed54470fda50d87a6c8c563671af2962c4f34f19af566227eb330ac8b9c770`.
Combined protected+generator inventory:
`7381a8ab915ba8a8d812ba3e98ac034b693bd23081ee48e74822da183e27748c`.

| Artifact | Unchanged SHA-256 |
| --- | --- |
| Spec | `8bb76b890086eb4bf89271edd5861e02f833f3ffa11a83b5865080ed4e4228cb` |
| Manifest | `3953dbc315b4dcb3d0d31dd399bf82a15cb38c49fde2e2b93ea81c5cfe0ec714` |
| Manifest schema | `f265771f904be57d19602dd6da9a572c16b3f80fdd166c5177721aaa3f9a91ec` |
| Case schema | `99805d442f13872cd4febe9ac8ae4f36fd25607580bc1e457ff5ae2efe199e23` |
| Moving profile | `7242fc557a7ca56452c934bedc5e7dc2835248f5cd36fa220f05588dc967ec80` |

All pre-existing conformance files except README and the two guard files match
starting HEAD; all history, module files, workspace, flake.lock and generator
files are unchanged. Tools has exactly the two intended deletions. Two non-Nix
`go work vendor -o <outside-path>` runs (starting HEAD archive and final tree)
produce **75 identical files**, paths/bytes/modes, including modules.txt.
Existing vendorHash is retained. Changed flake source/runtime inputs invalidate
W3B's Nix gate for this wave; no existing Nix result is reused as W3C evidence.

## Validation and execution accounting

Host: **Go 1.26.8 linux/amd64**. Final suite runs with failing python/python3
sentinels first on PATH; the sentinel invocation log remains absent. Final
validation invoked no Python, including standalone-module tests.

| Check | Agent execution/result |
| --- | --- |
| Baseline make check | 1, PASS; Python 3 tests OK, all Go packages PASS, 92 cases |
| Baseline vet/race/fuzz | 0/0/0 |
| Focused Go test invocations | 2, PASS: oracle; then oracle + repopolicy |
| Focused Python unit-test CLI invocations | 0 outside baseline; parity harness 1 invokes 22 three-method suites in-process |
| Parity Go build / test-binary runs | 1 compile-only build / 22 runs; dispositions agree |
| Focused fuzz Dispatch/Open/graph | 0/0/0 |
| Final complete host suites / incomplete attempts | 1 / 0, PASS |
| Final make check | 1, PASS; Python 0 tests, Go policy/architecture/tests/full/verify PASS |
| Explicit final make conformance / make verify | 1/1 additional invocations, both PASS, 92 cases |
| Final vet generator / conformance | 1/1, PASS |
| Final make race | 1, PASS, all packages in both modules |
| Final make fuzz | 1, all three targets PASS |
| Standalone GOWORK=off Go tests | 1 per module, PASS, -mod=readonly, proxy/sumdb off |
| Formatting / git diff --check | PASS |
| Generator CLI / vector regeneration | 0/0 |
| W3C human nix flake check / nix develop check | 1/1, both PASS, human-reported |
| Agent Nix | 0, not run |
| Cached W3C Nix / remote CI | Cache use unknown / NOT RUN |

Both module test/race/standalone runs cover every package in the isolation tree.
All three generator packages and 13 conformance packages with tests PASS;
cmd/check-architecture compiles with no test files. The new oracle has five
Test functions; the first three traverse committed manifest-listed artifacts.
Conformance vet and generator vet emitted no diagnostics.

| Fuzz target | Policy | Final result / package runtime |
| --- | --- | --- |
| FuzzDispatch, internal/object | -run '^$' -fuzztime 10s -parallel 2 | PASS / 10.080s |
| FuzzOpen, internal/cryptov1 | same | PASS / 10.800s |
| FuzzArrivalAndDisappearance, internal/graph | same | PASS / 10.104s |

Seeds, failure semantics and full race coverage are unchanged. No test failure,
fuzz reproduction input, or incomplete final suite occurred. A shell cleanup
command was rejected before execution by the command safety filter; cleanup was
performed through the scoped authoring script instead. No permission request
or destructive fallback was needed. The generated workspace-sum inventory issue
is described above rather than hidden as a qualification pass.

## Remaining Python references

Repository source inventory excluding .git/.direnv has **zero physical .py files**.
Makefile, flake.nix, .github and .envrc have zero Python/checker-name matches.
The live command expansion contains no Python or generation invocation.
Remaining matches are intentional:

- Historical review reports retain their original Python commands and evidence.
- CONTRIBUTING's r14→r15 and r15→r16 commands explicitly apply to the named
  historical checkouts, not current tooling.
- README, CONTRIBUTING and conformance README state Python is no longer required.
- AGENTS retains PYTHON_VERIFIER classification and historical execution/
  invalidation accounting vocabulary; its live example now uses Go.
- Existing specpolicy/schema code comments, helper/test names describe preserved
  Python whitespace, numeric-equality and regex compatibility semantics.
- Existing generator test comment describes historically Python-checked fields.
- New oracle comment records the migrated evidence's provenance; this report
  records removal/parity. None is an interpreter invocation.

Ignored dependency/tool caches are outside repository source membership and
outside this claim. No Nix transitive closure inventory or zero-Python closure
claim is made.

## Human gate and frozen handoff

Nix is **unavailable on the agent PATH** (`command -v nix` returned no path).
AGENTS also requires human-only Nix execution. On 2026-10-12 the owner reported
"both passed" for the two requested commands from this checkout:

```sh
nix flake check
nix develop --command make check
```

The first is the final pinned qualification gate; the second validates the normal
development environment after Python removal. Both are **PASS, human-reported**,
one execution each. Agent Nix execution remains zero. Cache reuse versus fresh
derivation execution and detailed human logs were not supplied and are not inferred.
All intended new files are made Git-known using intent-to-add only; require
`git diff --cached --quiet`. Existing edits/deletions remain unstaged. Local
frozen path/state/mode/SHA-256 inventory excludes only this report, which the
actual flake fileset and structure consumers exclude. Any later flake-source
edit requires renewed qualification. Report-only edits do not invalidate it.
No artificial package materialization is required. Remote Linux/macOS/Windows
matrix evidence awaits the owner's later commit/push; no dispatch was made.

The repository is now operationally Python-free: **No Totipo support-tooling
command requires Python.** W3C is complete with no unresolved local qualification
gate. Remote CI remains deferred until the owner's later commit/push. All changes
remain uncommitted for review.

Frozen handoff: **286 Git-known path/state/mode/hash records**, excluding only
this report, including explicit absent/deleted paths. Inventory SHA-256:
`78bf91941a9fa96ea0b602f6513925cb6f913e335429d8ed350f088ecddba290`.
Source membership is captured separately in `frozen-membership.bin`. All local
evidence, vendor trees, test binaries and sentinel scripts remain outside Git.

After the human PASS, all 286 frozen input records and exact Git membership were
compared again: PASS. All 100 protected artifact files and eight generator files
still match baseline; branch/HEAD remain main/df820576c5076b9c9434a4f4d64a6b7c316b24e4.
The cached diff remains empty and `git diff --check` passes. Only this excluded
report was updated after the gate; no qualification input changed. The final host
suite was not repeated for this report-only update.

## Git evidence

The following is the final uncommitted diff/status, including this report.
New files have intent-to-add entries only; `git diff --cached --quiet` passes.

<!-- git-evidence -->
```text
$ git diff --stat
 .github/workflows/conformance.yml                    |   6 +-
 AGENTS.md                                            |  13 +-
 CONTRIBUTING.md                                      |  15 +-
 Makefile                                             |   4 +-
 README.md                                            |  20 +--
 SPEC_SUPPORT_TOOLING_WAVE3C_PYTHON_REMOVAL_REPORT.md | 360 ++++++++++++++++++++++++++++++++++++++
 conformance/README.md                                |  11 +-
 conformance/internal/repopolicy/policy.go            |   3 +
 conformance/internal/repopolicy/policy_test.go       |   4 +
 conformance/internal/wireoracle/oracle_test.go       | 364 +++++++++++++++++++++++++++++++++++++++
 flake.nix                                            |   8 +-
 tools/check_vectors.py                               | 101 -----------
 tools/test_r16_vectors.py                            | 103 -----------
 13 files changed, 766 insertions(+), 246 deletions(-)

$ git status --short
 M .github/workflows/conformance.yml
 M AGENTS.md
 M CONTRIBUTING.md
 M Makefile
 M README.md
 A SPEC_SUPPORT_TOOLING_WAVE3C_PYTHON_REMOVAL_REPORT.md
 M conformance/README.md
 M conformance/internal/repopolicy/policy.go
 M conformance/internal/repopolicy/policy_test.go
 A conformance/internal/wireoracle/oracle_test.go
 M flake.nix
 D tools/check_vectors.py
 D tools/test_r16_vectors.py
```
