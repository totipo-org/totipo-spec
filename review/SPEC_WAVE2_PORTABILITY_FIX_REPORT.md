# Spec Wave 2 portability containment fix

## 1. Starting state

Date: 2026-10-11. Work confined to `/home/niki/Sources/totipo-spec`.
Before editing, the required commands returned:

```text
git branch --show-current
main
git rev-parse HEAD
83a52243f4cd3993fc2e240820dd6c77f25310a5
git status --short
(empty)
```

Starting commit subject: `Consolidate Spec integrity verification` (Wave 2).
Protocol **v1/r19**; requirements `totipo-requirements-v1`, **moving-pre-rc**,
`spec_revision: r19`; **92 required / 92 manifest / 92 physical cases**.
Read AGENTS.md, the support-tooling audit, Wave 1 and Wave 2 reports, vectors
reader/tests, Makefile, qualification source filter and structure-test consumers.
Historical reports remain untouched. Inherited Wave 2 evidence records passing
Linux host checks and one human-reported Nix pass; neither qualifies this change.

## 2. Reported remote failure

The user supplied these failures on both macOS and Windows:

```text
--- FAIL: TestArtifactPins
    integrity_test.go:33: case symlink escapes vectors
--- FAIL: TestProfileRelationship
    integrity_test.go:86: case symlink escapes vectors
--- FAIL: TestManifestContract
    integrity_test.go:128: case symlink escapes vectors
```

Linux qualification passed. Remote runner path diagnostics were not supplied;
this report distinguishes a reproduced implementation defect from inferred
platform triggers. No fresh remote jobs were run or observed here.

## 3. Classification

**TEST_OR_QUALIFICATION_TOOLING**, **CONFORMANCE_GO**, **portability regression
fix**. No protocol semantic change. Go implementation/tests are effective Nix
qualification inputs and invalidate inherited Nix evidence. Build, dependency,
Nix and workflow configuration themselves are unchanged.

## 4. Exact error source and inputs

`vectors.Read` emits `case symlink escapes vectors` after the manifest entry's
lexical validation and before reading/checksumming/decoding its case bytes.
All three integrity tests call `copyIntegrityRoot`: `t.TempDir()` followed by
`os.CopyFS` of spec, vectors and requirements. Their first `Read(root)` fails
remotely before any mutation. The inputs are ordinary complete copied corpora,
not deliberately hostile case symlinks.

Old filesystem algorithm, for each entry:

1. Join lexical root, `vectors`, and validated entry path.
2. `EvalSymlinks` the candidate, resolving all its ancestors too.
3. `Abs` the lexical vectors base and the resolved candidate, ignoring Abs errors.
4. `Rel(base, candidate)`; reject errors or a nonlocal relative result.

Containment already used `filepath.Rel` and `filepath.IsLocal`, not a raw string
prefix. `Abs` cleans paths but does not resolve ancestor aliases. Only one side
passed through filesystem evaluation. The `cases/` textual prefix is a separate
manifest grammar rule and is not filesystem containment.

## 5. Proven root cause and reproduction

The base and candidate occupied different path namespaces. A candidate reached
through a legitimate root alias was resolved to its target, while its base
retained the alias spelling. `Rel` then correctly described those differently
spelled paths as outside one another; the caller incorrectly interpreted that
as evidence of a hostile case symlink.

Before editing production or test source, a disposable Python harness copied
spec/vectors/requirements into `/tmp/totipo-alias-vqt_j25c/real-root`, created
`alias-root -> real-root`, and invoked the existing Go verify-only CLI through
each root. Results:

```text
normal /tmp/totipo-alias-vqt_j25c/real-root exit 0
PASS: r19 manifest, profile and 92 case files verified
alias /tmp/totipo-alias-vqt_j25c/alias-root exit 1
FAIL: case symlink escapes vectors
exit status 1
lexical base: /tmp/totipo-alias-vqt_j25c/alias-root/vectors
resolved candidate: /tmp/totipo-alias-vqt_j25c/real-root/vectors/cases/graph/v1.graph.late-parent.001.json
```

No case symlink or repository-data modification was involved. The disposable
fixture was removed automatically; the transcript above preserves the failure
and its layout. The durable `TestReadRootAlias` now exercises the same class.

## 6. Platform explanation and evidence limits

Linux's ordinary temporary root in this session had no spelling-changing
ancestor alias; the two namespaces coincided. The three CI failures passed
locally as one focused group before editing (0.076s). A deliberately aliased
Linux root failed, demonstrating this is not intrinsically an OS-specific rule.

On macOS, a temporary ancestor such as `/var` can resolve to `/private/var`.
Resolving only the candidate makes a safe temporary corpus appear outside its
lexical root. Windows also evaluates reparse points and normalizes filesystem
spelling: local Go stdlib inspection of `path/filepath/symlink_windows.go` shows
`walkSymlinks` followed by `toNorm`, whose `normBase` uses `FindFirstFile` to
obtain actual names. Short/long path spellings can differ even without a case
file symlink. Case differences alone are not the claimed cause: Windows Rel
compares components with host case rules. Rel also checks volumes and rootedness.

These mechanisms explain the cross-platform vulnerability in the old algorithm.
The exact alias/short-name spelling in each remote job is **not independently
verified** from the supplied errors. The Linux reproduction proves the namespace
bug; fresh macOS/Windows jobs remain necessary to prove the remote fix.

## 7. New containment algorithm

`resolvedPath` applies checked `filepath.Abs` followed by checked
`filepath.EvalSymlinks`. `Read` resolves the vectors base once, then resolves each
lexically validated candidate identically. `containedPath` requires absolute
inputs, obtains `filepath.Rel`, and accepts only a successful local result via
`filepath.IsLocal`. Failed resolution or an unsafe relationship fails closed.

This uses host separator/volume semantics, including Windows drive/UNC rules;
there is no production OS branch, CI/temp-name special case, textual containment
prefix, or unconditional acceptance after evaluation. The existing `Rel == "."`
local policy is retained; subsequent file read/hash/shape checks remain in force.

## 8. Security properties preserved

Lexical validation remains separate and unchanged: paths must be local,
backslash-free, clean/canonical in slash form, and begin with `cases/`.
Traversal and noncanonical forms stay rejected before candidate evaluation.
Duplicate IDs/paths, exact physical enumeration, case hashes, identity, expected
outcomes and shapes retain their existing checks. Profile validation is unchanged.

Both a case-file symlink and a directory symlink resolving outside vectors are
rejected with the original containment error. Missing/unresolvable paths return
errors. No filesystem alias can make an invalid manifest path bypass its grammar.
The boundary remains the intended vectors root, as before; protocol storage
semantics and evaluators are unchanged. This fix does not introduce a new claim
of protection against concurrent filesystem replacement between check and read.

## 9. Regression tests

Added `conformance/internal/vectors/path_test.go`:

- `TestContainedPath`: child and same-base policy; parent, sibling, shared textual
  prefix and nonabsolute rejection. These deterministic tests need no symlinks.
  On Windows, additional native drive/case/UNC checks test cross-volume/share
  rejection and valid host spelling relationships without filesystem privileges.
- `TestReadRootAlias`: normal copied corpus loads 92 cases; aliased root loads
  deeply identical manifests/cases and passes the same profile relationship.
- `TestReadCaseSymlinkEscape`: file and ancestor-directory escapes fail with the
  intended error, using real temporary filesystem fixtures.
- `TestResolvedPathMissing`: resolution failure is rejected.

Existing hostile lexical/duplicate/inventory tests were neither weakened nor
removed. Windows fixture creation alone may skip for privilege error 1314 or
permission denial; all other creation errors fail. Deterministic boundary checks
and the ordinary copied-corpus tests still run without that privilege. No fixture
skipped on this Linux host. Production has no Windows exemption.

## 10. Focused results and baseline

Baseline `make check`: **1 PASS**, 15 Python tests (0.445s), seven Go test packages,
92-case conformance, 92-case schemas, and manifest/profile integrity PASS.
No baseline race or fuzz was needed: filesystem containment does not change
concurrent execution or the fuzz-sensitive parsers/models.

Focused Go invocations, all `-count=1`:

1. Before editing: the three reported CI failures grouped with exact anchored
   test-name regex, **PASS**, 0.076s.
2. After editing: the same three tests plus all four new path tests grouped by
   anchored regex, **PASS**, 0.121s; lexical and real escape subtests PASS.
3. Entire affected `./internal/vectors` package, **PASS**, 0.951s.

Focused Python invocations: 0. Focused fuzz runs: 0 per target. The pre-edit
verify-only reproduction made two CLI invocations (safe PASS, alias expected
failure); it was not a failing qualification suite. No generator was run.

## 11. Protected-byte and differential semantic evidence

Compared starting Git blobs to raw working-tree bytes and path inventories,
including untracked nonignored additions and deletions:

| Protected scope | Files | Result |
| --- | ---: | --- |
| spec | 1 | Identical bytes/paths |
| vectors | 97 | Identical bytes/paths |
| requirements | 2 | Identical bytes/paths |
| tools | 5 | Identical bytes/paths |
| .github | 2 | Identical bytes/paths |
| Existing review/history | 33 | Identical bytes/paths |

Makefile, flake.nix, flake.lock, go.work, every go.mod/go.sum and generator source
are identical. Root go.mod/go.sum remain absent. No dependency or external
package was added. Exactly 92 physical cases remain. Identical manifest, profile
and case bytes prove unchanged IDs, paths, hashes and authored expected outcomes;
normal, alias and final profile checks confirm the same required relationship.
Generation executions: **0**. No corpus/schema/profile adjustment was necessary.

## 12. Final host qualification

All included source inputs frozen before one final suite. Final `make check`:
structural PASS; **15 Python tests in 0.450s, OK**; full CLI **92 cases PASS**;
schemas for **92 cases PASS**; manifest/profile/**92 case files verified PASS**.
Vet: exit 0, no diagnostics. Formatting and `git diff --check`: PASS.

| Package under totipo/conformance | Normal test result | All-package race result |
| --- | --- | --- |
| cmd/generate-vectors | no test files | no test files |
| cmd/totipo-conformance | no test files | no test files |
| internal/cryptov1 | ok 0.503s | ok 1.823s |
| internal/graph | ok 0.008s | ok 1.048s |
| internal/object | ok 0.002s | ok 1.011s |
| internal/storage | ok 0.480s | ok 1.676s |
| internal/tlv | ok 0.002s | ok 1.013s |
| internal/totp | ok 0.006s | ok 1.014s |
| internal/vectors | ok 1.151s | ok 3.698s |

Normal and race retained all-package `-count=1`; race retained `-race`.
Final `make fuzz`: all targets PASS with unchanged `-run '^$'`,
`-fuzztime 10s`, `-parallel 2` policy and seeds/failure semantics.

| Target / package | Initial coverage | Executions | New interesting / total | Duration / result |
| --- | --- | ---: | --- | --- |
| FuzzDispatch / internal/object | 50/50 | 1,131,846 | 0 / 50 | 10.088s PASS |
| FuzzOpen / internal/cryptov1 | 16/16 | 578,741 | 0 / 16 | 11.013s PASS |
| FuzzArrivalAndDisappearance / internal/graph | 130/130 | 46,257 | 2 / 132 | 11.016s PASS |

Elapsed package times include worker shutdown beyond the 10s mutation budget.
Existing local fuzz-cache entries contributed initial coverage; they are not
cached Nix evidence or substitutes for this fresh run. No seed source changed.

## 13. Execution accounting

Baseline make check: 1 PASS. Focused Python/Go test invocations: 0 / 3 PASS.
Focused fuzz Dispatch/Open/graph: 0 / 0 / 0. Final complete suites: 1;
incomplete attempts: 0. Final check/vet/race/fuzz invocations: 1 / 1 / 1 / 1.
Each final fuzz target: 1. Generator: 0. Human Nix: 1,
human-reported PASS for the frozen fix; agent Nix: 0. Fresh remote CI observations: 0; cached Nix observations:
none. Inherited Wave 2 host/human-Nix evidence is separate from these counts.

Raw baseline/focused/final logs are local temporary evidence at
`/tmp/totipo-portability-qualification/`. The pre-edit three-test output was
observed directly in the session; the reproduction transcript is preserved above.
No final failure or fuzz reproduction input was produced.

## 14. Human Nix gate

**PASS — human-reported** for the frozen fix. The user reported “nix flake
check passed” after the source-freeze handoff. Human execution count: 1; agent
execution count: 0. This records the human result, not an agent-observed run.
AGENTS.md explicitly says
“The agent does not run Nix” and requires the final human command after source
freeze:

```sh
nix flake check path:.
```

Only this excluded report was updated after the human result; no qualification
source input was edited during this update. The actual filter includes these Go
implementation/test files. This new report is excluded; structure checks consume
only the specifically included archived history, which remains unchanged.
Report-only updates do not invalidate qualification. No `nix build` requested.

## 15. Remote portability status

**NOT VERIFIED**. Linux alias and security regressions passing do not establish
macOS or Windows success. After the user's commit/push, the unchanged portability
matrix must run and at least the previously failing macOS and Windows jobs must
pass. No commit/push or workflow dispatch was performed. The authoritative Linux Nix CI boundary and supplemental portability matrix
remain unchanged; no fresh remote result is claimed for either.

## 16. Changed files

- `conformance/internal/vectors/runner.go`: consistent resolved containment.
- `conformance/internal/vectors/path_test.go`: new portability/security regressions.
- `review/SPEC_WAVE2_PORTABILITY_FIX_REPORT.md`: this report.

## 17. Final Git state

Branch `main`; HEAD remains `83a52243f4cd3993fc2e240820dd6c77f25310a5`.
Index unchanged; all work unstaged/uncommitted:

```text
 M conformance/internal/vectors/runner.go
?? conformance/internal/vectors/path_test.go
?? review/SPEC_WAVE2_PORTABILITY_FIX_REPORT.md
```

No staging, commit, push, tag, release, generator, agent Nix or CI dispatch.
Host qualification and human Nix have passed; the milestone is locally complete.
Full portability completion additionally requires remote macOS/Windows PASS.
No stop condition was encountered.

## 18. Suggested commit subject

```text
Fix cross-platform vector path containment
```
