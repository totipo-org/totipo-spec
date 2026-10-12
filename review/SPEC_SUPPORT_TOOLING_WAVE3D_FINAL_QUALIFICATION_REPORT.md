# Wave 3D — isolation enforcement and final qualification

## 1. Starting state

Date: 2026-10-12. Repository: /home/niki/Sources/totipo-spec only.
Starting branch **main**, full HEAD
**e9c9f5fbeaeb211a2d01a57cfcc92de821912ac6**.
Initial git status was empty. W3C and its completed report are committed at this
HEAD (subject: Remove python verifier). No sibling repository or agent worktree
was used. Protocol and moving requirements revision: **v1/r19**.

Classifications: **CONFORMANCE_GO**, **BUILD_NIX_CI** (qualification policy),
**DOCS_ONLY** (supporting guidance/report). No protocol, corpus, requirements,
generator behavior, dependencies, Make targets, CI matrix, or Nix filter changes.

Starting commands: make architecture, spec-check, test, verify, conformance,
check, race, fuzz; make generate is deliberate maintenance only. check executes
architecture/spec-check/both-module tests/full conformance/verify-only.
Baseline make check ran once and passed: **92 manifest cases**, all Go tests,
structural/schema verification; **zero Python tests**.

## 2. Inherited Wave 3A–3C evidence

Read the completed W3A boundary, W3B structural/schema and W3C removal reports.
Their historical contents remain unchanged.

- W3A: separate modules and independently owned models, literals and writer
  composition; temporary regeneration comparisons ultimately matched all
  100 protected files. Host and human-reported Nix PASS.
- W3B: independent Go schema/spec policy; 1,726 schema and 138 structure
  differential comparisons with zero disposition drift, four agreed escaping
  path rejections, permanent required-member/negative coverage. Eager schema
  preflight is the documented strengthening. Host and human-reported Nix PASS.
- W3C: Python cutover committed; 22 independent-wire comparisons with zero drift,
  stdlib-only Go oracle retained; all operational Python removed. Host PASS.
  Latest recorded human evidence on 2026-10-12 is “both passed” for
  nix flake check and nix develop --command make check, one execution each.
  Its post-gate frozen-input check passed.

These are inherited evidence for prior inputs. Human cache versus fresh build
provenance is unknown. They do not qualify W3D changes. No inherited remote
platform result is upgraded to a pass.

## 3. Final architecture and module graph

```text
go.work (Go 1.23.0; exactly ./generator and ./conformance)
  totipo/generator (Go 1.23.0)
    cmd/generate-vectors -> own internal/{fixture,wire} -> stdlib
        |
        | deliberate maintenance writes serialized files
        v
  committed spec/ + vectors/ + requirements/
        |
        | read serialized schemas, manifest, cases and pins
        v
  totipo/conformance (Go 1.23.0)
    cmd/totipo-conformance -> internal/vectors -> schema + semantic packages
    cmd/check-spec -> internal/specpolicy -> stdlib
    cmd/check-architecture -> internal/repopolicy -> stdlib + Go inspection
    internal/wireoracle tests -> own serialized reader/formulas -> stdlib
    reviewed primitives: x/crypto/{argon2,blake2b,hkdf}, x/sys/cpu
```

There is no generator/verifier implementation edge in either direction, no
shared semantic module, no root/third module, and no replace directive.
The workspace is convenience; metadata, source and resolved closure checks
establish isolation. Direct verification compiles the consumer only.

## 4. Automatically enforced invariants

The established check-architecture CLI remains the sole architecture gate.
Source scanning covers tests, generated-marked and inactive-tag source; module
and workspace metadata reject third modules, unauthorized dependencies and
replacement bridges. Resolved package/test closures reject indirect coupling,
load errors, unreviewed primitives, ownership escapes and vendor symlink escapes.
Existing source bridge, process delegation, spec-policy and oracle restrictions
remain. Generator dependency closure is stdlib-only.

W3D adds **Wiring** to this same gate: an explicit allowlist of the current
Makefile's executable declarations, dependency edges and recipes. Generation
remains confined to make generate; normal checking cannot gain generation,
include another Makefile, hide execution in variable expansion, or remove
architecture/schema verification without policy review. Unknown executable
wiring fails closed. Ordinary Make comments and blank lines are ignored;
recipe comments and backslash-continued comments require review because they
can affect execution. No documentation/history text is scanned for forbidden
command strings. This deliberately avoids a partial general Make interpreter.

Schema/structure execution remains independent: check-spec reads spec/history;
vectors.Read interprets committed schema sources before typed acceptance.
Successful generation does not confer verification status.

## 5. Negative-test evidence

Disposable trees are used; the live repository was never intentionally broken.
Focused tests and final all-package tests passed:

| Violation | Regression evidence |
| --- | --- |
| Conformance imports generator | verifier-import; cross-role source rejection |
| Generator imports conformance | generator-import; tests/tagged/generated source also covered |
| Shared semantic implementation | third-package, outside-source, third/root-module; new shared-semantic-module using example.org/semantics |
| Metadata/workspace bypass | require-bridge, module replace, workspace-replace, extra-workspace; new primitive-local-replacement |
| Indirect or unavailable dependency | actual dependency-closure sibling/load/external fixtures |
| Generation in normal checking | generation-prerequisite, generation-recipe, recursive-generation, alias |
| Hidden executable Make configuration | shell-expansion, override, include, recipe-prefix |
| Weakened verification | remove-architecture, remove-schema-verification, weaken-verifier |
| Comment-related Make bypass | continued-comment, recipe-comment-expansion |

There are **32 source/metadata mutation subtests** and **13 wiring mutation
subtests**, all rejecting as expected. Harmless Make comments, documentation
command strings and an ordinary explanatory Go import comment are accepted.
Existing vendor containment, root-alias, missing-directory and source-symlink
regressions remain active.

## 6. Dependency and standalone results

GOWORK=off, GOPROXY=off, GOSUMDB=off standalone list -mod=readonly -deps -test
-json ./... passed separately for both modules. Standalone tests and vet also
passed in both modules. Module/workspace manifests and checksums are unchanged.
Root go.mod/go.sum/go.work.sum and generator/go.sum remain absent.

Generator has three own packages plus test variants, no external packages.
Conformance's external imported closure is exactly:
golang.org/x/crypto/{argon2,blake2b,hkdf} **v0.36.0** and
golang.org/x/sys/cpu **v0.31.0**. Nonstdlib external intersection: zero.
Standalone mod graph confirms only the original direct pins; x/net v0.21.0,
x/term v0.30.0 and x/text v0.23.0 are transitive module-metadata edges from
x/crypto, not imported semantic packages. Both standalone go mod verify
commands reported “all modules verified”. No dependency repair or update.

One non-Nix go work vendor in an outside archived copy produced **75 files**,
including modules.txt with exactly the reviewed four primitive packages.
Updated policy source/tests and the actual gate passed in that isolated tree
with GOFLAGS=-mod=vendor, an empty GOMODCACHE and proxy/sumdb disabled.
No live vendor tree or manifest was changed.

## 7. Read-only verification proof

New TestCLIVerificationLeavesArtifactsUnchanged copies only the 100 canonical
artifact files to disposable roots, without generator or module source. It
compares path-to-SHA-256 maps before and after both real CLI run modes:
verify-only and full conformance. It checks successful input and malformed
manifest refusal. All four invocations preserve artifact file paths and bytes.
This calls the production run seam; no subprocess exception was introduced.

Host qualification separately compared every protected path/hash before and
after normal commands, and compared tracked bytes against starting Git HEAD.
Generator execution and Python sentinels guarded the successful full suites;
neither forbidden invocation log exists. Generator unit tests compile/test writer
code as explicit repository policy; they do not execute generator main or write
canonical artifacts. Go build/test caches and synthetic test outputs are outside
canonical source. No timestamp-only evidence is used.

## 8. Complete host qualification and execution accounting

Host: **Go 1.26.8, Linux amd64**. Final inputs pass make check, make verify,
make conformance, both module vets, make race, make fuzz, formatting and
git diff --check. All 92 conformance cases pass; schema/spec-policy negatives,
architecture negatives and independent wire tests pass.

Generator: cmd/generate-vectors, internal/fixture, internal/wire all PASS.
Conformance: cmd/check-spec, cmd/totipo-conformance and internal/{cryptov1,
graph,object,repopolicy,schema,specpolicy,storage,tlv,totp,vectors,wireoracle}
all PASS. cmd/check-architecture compiles with no command tests. All-package
race coverage includes both modules and all those tests. Both vets emit no
diagnostics. No final scope, seeds, flags or failure semantics were reduced.

Counts are top-level invocations, excluding inspection subprocesses within
architecture tests and the CLI run seam within unit tests:

| Execution | Count / result |
| --- | --- |
| Baseline make check | 1 PASS |
| Baseline vet/race/fuzz | 0/0/0 |
| Focused ordinary Go test invocations | 2 PASS: policy+CLI, then policy comment-bypass tests |
| Focused Python / Python tests in any W3D suite | 0 / 0 |
| Focused fuzz outside complete suites | 0 per target |
| Complete host suites | 2 PASS; second qualifies final comment-hardening inputs |
| Incomplete final attempts | 1 harness failure before architecture Go execution |
| Final-phase make check attempts | 3: 1 harness failure, 2 PASS |
| Final-phase explicit make verify / make conformance | 2 / 2 PASS, additional to make check |
| Final-phase workspace vet per module | 2 each PASS |
| Final-phase make race / make fuzz | 2 / 2 PASS |
| Standalone dependency listing / mod graph / mod verify | 1 each per module, PASS |
| Standalone test/vet generator | 1 / 1 PASS |
| Standalone test/vet conformance | 2 / 2 PASS, repeated after comment fix |
| Isolated vendored policy tests / architecture CLI | 2 / 2 PASS, repeated after comment fix |
| Non-Nix workspace vendor generation | 1 PASS, outside checkout |
| Generator CLI / canonical regeneration | 0 / 0 |
| Agent Nix / human W3D Nix | 0 / 1; PASS, human-reported |
| Remote execution/dispatch/commit/push | 0; PENDING |

The incomplete attempt used a sentinel with a /bin/bash shebang; that path does
not exist on this host. make stopped with error 127 before running the checker.
Failure logs are preserved. Corrected sentinels use resolved shell paths.
After the first complete PASS, code review found Make comment interpretation
bypasses; policy and negative tests were hardened and the full suite repeated.
There was no protocol/test assertion failure or fuzz failure/reproducer.

All fuzz runs use exactly -run '^$' -fuzztime 10s -parallel 2:

| Target / package | Earlier complete PASS runtime | Final PASS runtime |
| --- | --- | --- |
| FuzzDispatch / object | 10.043s | 10.094s |
| FuzzOpen / cryptov1 | 10.109s | 10.047s |
| FuzzArrivalAndDisappearance / graph | 11.030s | 10.061s |

Runtime includes worker shutdown and may exceed the requested fuzz budget.
Evidence logs, dependency JSON, sentinels, vendor copy and inventories are local
supplementary evidence at /tmp/totipo-w3d-evidence; this report preserves results.

## 9. Python and active documentation audit

No tracked Python source or active Make/Nix/Actions Python requirement remains.
Successful complete suites used failing python/python3 sentinels without any
invocation. Go schema/structure and independent wire-oracle replacements run.
Historical reports/audit-checkout instructions and compatibility comments are
preserved. No zero-Python transitive Nix closure claim is made.

README, CONTRIBUTING and conformance README already explain separate Go roles,
committed artifacts, explicit generation, Go-only support and current commands.
CONTRIBUTING now documents wiring enforcement, its trust boundary, command
coverage and standalone test/vet commands. The obsolete generator path in
protected vectors/FORMAT.md remains byte-identical and is explicitly superseded
by CONTRIBUTING, as before. No protocol prose cleanup or unrelated rewrite.

## 10. Final byte-invariance audit

Against e9c9f5fbeaeb211a2d01a57cfcc92de821912ac6:
**100/100 protected files identical**, including path inventory; no additions,
deletions or byte differences. spec 1, vectors 97, requirements 2. Manifest,
physical case files and required profile entries remain **92/92/92**.
All eight generator files are unchanged. Production semantic consumer, schema,
spec-policy and oracle behavior/source remain unchanged; only architecture policy
and consumer tests change. No authored expectations or pins were repaired.

Sorted protected SHA-256 inventory digest:
**36ed54470fda50d87a6c8c563671af2962c4f34f19af566227eb330ac8b9c770**.

| Artifact | Unchanged SHA-256 |
| --- | --- |
| Spec | 8bb76b890086eb4bf89271edd5861e02f833f3ffa11a83b5865080ed4e4228cb |
| Manifest | 3953dbc315b4dcb3d0d31dd399bf82a15cb38c49fde2e2b93ea81c5cfe0ec714 |
| Manifest schema | f265771f904be57d19602dd6da9a572c16b3f80fdd166c5177721aaa3f9a91ec |
| Case schema | 99805d442f13872cd4febe9ac8ae4f36fd25607580bc1e457ff5ae2efe199e23 |
| Requirements profile | 7242fc557a7ca56452c934bedc5e7dc2835248f5cd36fa220f05588dc967ec80 |

go.work, both go.mod files, conformance/go.sum, Makefile, flake.nix/lock,
workflows, historical reports and every other unchanged tracked path retain
starting bytes. Source additions are limited to the two wiring Go files/report.

## 11. Vendor portability regression

Existing checkDependencyDir still trusts Go-owned stdlib paths without trying
to resolve Windows toolchain junctions. Nonstdlib directories remain resolved,
with contextual failures. Both vendor and repository ancestors are canonicalized
for macOS aliases; escaping vendor symlinks remain rejected.
TestDependencyDirStandardLibrary, TestDependencyDirVendor and
TestDependencyDirVendorAliasedRoot pass on this Linux host, including escape
rejections. No symlink containment or vendor boundary was weakened.
Windows may skip symlink creation only when privileges are unavailable;
Windows/macOS execution is not inferred from host PASS.

## 12. Nix qualification and frozen handoff

**PASS — human-reported; W3D human execution count 1; agent count 0.**
On 2026-10-12 the owner reported “nix flake check passed” for the frozen
checkout. Cache reuse versus fresh derivation execution was not supplied and
is not inferred. Post-gate verification passed for all 287 frozen file hashes
and modes, exact Git membership, all 100 protected artifact hashes, unchanged
branch/HEAD, empty cached diff and diff whitespace checks. Only this excluded
report was updated afterward; no qualification input changed.
AGENTS.md states “The agent does not run Nix.” The owner gate, after frozen host
qualification and Git-aware membership, is:

```sh
nix flake check
```

qualificationSource includes spec/vectors/requirements, all Go/native/module
inputs (including both new wiring files), Makefile, go.work and archived pre-r16
history. Linux derivation executes make check, both vets, race and fuzz offline
with the existing vendorHash. No fileset change was necessary. Workflow/flake
definitions are evaluated by their established outer qualification paths; no
workflow change is made. README/CONTRIBUTING/AGENTS and this report are excluded
and not consumed by structure checks.

All new intended files, including this report, receive intent-to-add only.
No staged file content is permitted; git diff --cached --quiet must pass.
Frozen actual path/mode/hash inventory and exact Git membership are retained
outside Git, excluding this report only. Later qualification-input changes
invalidate a human PASS; report-only updates do not. Prior W3C human PASS and
unknown cache provenance remain prior evidence only.

## 13. Remote qualification

**Linux Nix: PENDING. macOS: PENDING. Windows: PENDING for W3D.**
No exact final commit exists and no commit/push/dispatch is authorized.
The host has no gh executable; authoritative remote results were not obtained.

Inspected existing matrix, unchanged: authoritative ubuntu-24.04 Nix plus
supplemental macos-latest and windows-latest on Go 1.23.x/stable, and
ubuntu-24.04 on Go 1.23.x (five portability combinations).
Every portability job runs structure/schema verification, both module tests/vet,
the actual architecture command, standalone tests and full/verify-only consumer
commands. Thus policy negatives, CLI read-only tests and vendor portability
regressions are exercised on all matrix platforms. Linux Nix additionally
qualifies offline vendoring, all-package race and the three fuzz targets.
No redundant job or matrix expansion is needed.

## 14. Changed files and final Git state

| Path | Purpose |
| --- | --- |
| CONTRIBUTING.md | Enforcement/trust boundaries and standalone commands |
| conformance/cmd/totipo-conformance/main_test.go | Canonical path/hash read-only CLI regression |
| conformance/internal/repopolicy/policy.go | Invoke existing gate's wiring check |
| conformance/internal/repopolicy/policy_test.go | Shared-module/replacement negatives and Make fixture |
| conformance/internal/repopolicy/wiring.go | Reviewed executable Make wiring allowlist |
| conformance/internal/repopolicy/wiring_test.go | Thirteen negative wiring cases and harmless-text acceptance |
| SPEC_SUPPORT_TOOLING_WAVE3D_FINAL_QUALIFICATION_REPORT.md | This report |

Branch and HEAD remain main/e9c9f5fbeaeb211a2d01a57cfcc92de821912ac6.
No commit, tag, release, push, dependency update, actual-content staging, or
unrelated source membership. Final status is appended below after intent-to-add.

## 15. Limits and residual risks

This gate is repository policy, not a hostile-code sandbox or proof of independent
algorithmic ancestry. Arbitrary policy rewrites, malicious replacement toolchains,
caller GO/GOFLAGS/environment overrides and direct manual commands are trusted
review boundaries. Nix/Actions wiring is reviewed and executed through established
paths, not parsed by the architecture gate. The Make allowlist intentionally
requires explicit review for even benign executable-wiring changes.

Read-only snapshots prove final paths/content for exercised commands, not absence
of a transient write restored before return, nor every possible future input.
Isolation/source checks and reviewed code supplement that observation. Oracle
scope remains selected formulas, not independent AES-GCM/Argon2 or external
interoperability. Local /tmp evidence is supplementary, not a durable CI archive.

## 16. Completion recommendation

Host implementation/qualification is ready for owner review; no architectural
migration blocker remains. Host and human W3D Nix qualification PASS.
**Wave 3 final closure remains PENDING** until exact owner-committed inputs
pass authoritative Linux and supplemental macOS/Windows CI. Do not claim Wave 3
complete while those gates are pending. No later migration wave is started.

## Frozen handoff and Git evidence

Final inputs are frozen: **287 Git-known files excluding this report**, captured
with actual SHA-256 and filesystem mode inventories. Membership includes all
three new intended files using intent-to-add only. No staged content and no
unintended nonignored untracked file remain.

| Local inventory | SHA-256 |
| --- | --- |
| frozen-sha256.txt | 2f63f8c66a2d608921be00854e1a9ce4efbde03eef82779ec754b6fb0c1cf30a |
| frozen-modes.txt | e7e02b94218bcc87605fd8d86d7ac5d1ae82156ab0fbdd40bc7fb93629b42430 |
| frozen-membership.bin | e5ba763f2e29713cc8ffc5ac256a97b0027dc77019a5cf12aa79f710d775aa89 |

Protected inventory was compared again at handoff: PASS, all 100 paths/hashes.
Unchanged generator, workspace/module files, Make/Nix/CI, and historical reports
were compared to starting HEAD again: PASS. Empty cached diff, formatting and
diff whitespace checks PASS. This excluded report is the only file finalized
after the freeze; the successful final host suite was not repeated for prose.

```text
$ git status --short
 M CONTRIBUTING.md
 A SPEC_SUPPORT_TOOLING_WAVE3D_FINAL_QUALIFICATION_REPORT.md
 M conformance/cmd/totipo-conformance/main_test.go
 M conformance/internal/repopolicy/policy.go
 M conformance/internal/repopolicy/policy_test.go
 A conformance/internal/repopolicy/wiring.go
 A conformance/internal/repopolicy/wiring_test.go

$ git diff --cached --quiet
(exit 0)
```

Human W3D Nix: PASS, one human-reported execution, with post-gate frozen-input
verification PASS. Exact-commit remote platform statuses remain PENDING.
