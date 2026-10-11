# Spec support-tooling pruning Wave 1

## 1. Starting revision/state

Date: 2026-10-11. Work was confined to `totipo-spec`. Before editing:

```text
git branch --show-current
main
git rev-parse HEAD
0591c59753d78c168e1214b32430ffd0d2f934ef
git status --short
(empty)
```

Clean `main` and committed `review/SPEC_SUPPORT_TOOLING_AUDIT.md` were verified.
Protocol is **v1/r19**; requirements are `totipo-requirements-v1`,
`moving-pre-rc`, `spec_revision: r19`. The current manifest contains **92 cases**.
No revision advancement occurred.

## 2. Classification and qualification boundary

**TEST_OR_QUALIFICATION_TOOLING / active-tree historical-tool cleanup**, as
requested. Under AGENTS.md's classifications, **BUILD_NIX_CI** applies because
the deleted Python files are Nix source-projection inputs. There is also a current
maintenance-document change. This is not an overall DOCS_ONLY milestone, and no
active Python verifier, Go consumer, generator, protocol or corpus logic changed.

The semantic-work baseline `make check` was not run: this is infrastructure/source
cleanup with explicit protected-byte evidence. The complete final host suite was
run once after the two deletions and maintenance policy stabilized. No narrower
docs-only gate was used.

Inherited evidence: the committed support-tooling audit cites passing host and
isolated-source evidence from the Nix/CI unification report. The Actions refresh
report leaves human Nix **PENDING** and remote CI **NOT RUN**. These are inherited
records, not executions or passes for this wave; an older r19 Nix pass does not
resolve the pending later gate.

## 3. Historical caller/reference analysis

Before deletion, repository-wide `rg` searches for `audit_r15` and `audit_r16`
found historical commands/links in the r15 simplification report and r16 rewrite
report, paths in the r15 audit JSON and r16 final-status snapshots, and descriptions
in the support-tooling audit and Nix/CI report. These are historical invocations
and inventories, not current qualification callers.

The scripts were read before deletion. The r15 audit compares r14 Git objects
against an r15 worktree, checking old applicability, frozen bytes and legacy
references. The r16 audit compares r15 Git objects against an r16 worktree with
fixed replacement mappings. Neither is a current r19 validator. Makefile,
flake.nix, workflow commands, Go source, and remaining Python checkers/tests have
no executable dependency on either program. No stop condition was encountered.

## 4. Historical retrieval policy

[CONTRIBUTING.md](../CONTRIBUTING.md) now states that historical one-off programs
are reproduced from their historical Git revision, and main need not retain
executables whose migration responsibility has ended. It gives the exact source
locations, commands, clean separate checkout guidance and a detached-worktree
example. Repository history retention is the dependency; no external archival
guarantee is claimed.

| Audit | Source revision/path | Fixed comparison baseline |
| --- | --- | --- |
| r14 → r15 | `5ff6977ae004a9c03a410b38f0f9d73ca2bc76d4:tools/audit_r15.py` | `0b63f77886f95cf9bab3ceb9cbb0916177fe3589` |
| r15 → r16 | `54f38d13d62dc811c98489eb1867d63a9bef235a:tools/audit_r16.py` | `0a01fc1f0493ec1fefc0c4b89f8b22490d586852` |

`git show` proved both historical source blobs equal their pre-deletion bytes.
Those revisions contain r15/105-case and r16/90-case manifests, respectively,
and their corresponding requirements revisions. `git cat-file -e` confirmed both
baseline commits are available. Both audit JSON files and both primary reports
at those historical source revisions equal their current preserved bytes.
Historical reproduction was established through committed source/input/output
availability and inspection; the old audit programs were not executed in this
wave. The r15 legacy-reference inventory depends on checkout contents, so it must
be interpreted in that historical worktree context.

## 5. Deleted files

Deleted exactly `tools/audit_r15.py` (125 lines) and `tools/audit_r16.py` (89 lines):
214 Python lines removed. No Go port, wrapper or replacement executable was added.

## 6. Preserved historical evidence

All existing `review/` files remain byte-identical to starting HEAD, including:

- [r15 audit JSON](V1_R15_CASE_AUDIT.json) and [r15 report](V1_R15_STATE_MODEL_SIMPLIFICATION_REPORT.md).
- [r16 audit JSON](V1_R16_CASE_AUDIT.json), [r16 rewrite report](V1_R16_REWRITE_REPORT.md),
  [r16 final-status snapshot](V1_R16_FINAL_STATUS.txt), and design checkpoints.
- [support-tooling audit](SPEC_SUPPORT_TOOLING_AUDIT.md) and archived revision history.

Original relative source links and commands in historical reports were left
verbatim. CONTRIBUTING explains their historical-tree context.

## 7. Protected-byte invariance

A raw-byte comparison used `git show 0591c59753d78c168e1214b32430ffd0d2f934ef:<path>`
against working-tree bytes, with starting/current path inventory checks including
untracked nonignored additions. Results:

| Scope | Starting files | Result |
| --- | ---: | --- |
| `spec/**` | 1 | Identical bytes; no additions/deletions |
| `vectors/**` | 97 | Identical bytes; no additions/deletions |
| `requirements/**` | 2 | Identical bytes; no additions/deletions |
| `conformance/**` | 23 | Identical bytes; no additions/deletions |
| `tools/**` | 7 | Exactly the two audit deletions; other five files identical; no additions |

`go.work` and all `go.mod`/`go.sum` files are identical; root `go.mod` and `go.sum`
are absent at both boundaries. Makefile, flake.nix, flake.lock and all `.github/`
inputs are identical. All existing reports/output evidence are identical.
Active Python checkers/tests, independent wire tests and the generator are intact.
No case/profile hashes were regenerated.

## 8. Live-call-graph result

Post-deletion repository-wide searches again found only the preserved historical
references plus the new CONTRIBUTING retrieval instructions and this report.
Separate inspection/search of Makefile, flake.nix, `.github/`, `tools/` and
`conformance/` found no live caller or import of either deleted audit.
Current qualification documentation does not require running either audit on
main; the new maintenance instructions explicitly require historical checkouts.
Unobserved external callers are outside this repository-local proof.

## 9. Host qualification commands/results

All final gates passed on the host, once each, with exit status 0. Raw logs are
retained locally at `/tmp/totipo-wave1-qualification.XeRuII/` (`check.log`,
`vet.log`, `race.log`, `fuzz.log`); these temporary files are not committed evidence.

`make check` emitted:

- PASS: Totipo v1/r19 structural checks and profile artifact hashes.
- Python discovery: **Ran 14 tests in 0.632s; OK**.
- Manifest conformance: **PASS: v1/r19 moving-pre-rc: 92 cases**.
- Schema verification: **PASS: JSON schemas for manifest and 92 cases**.
- Integrity verification: **PASS: r19 manifest, profile and 92 case files verified**.

Exact Go package results from normal and race runs:

| Package under `totipo/conformance/` | `make check` | `make race` |
| --- | --- | --- |
| `cmd/generate-vectors` | [no test files] | [no test files] |
| `cmd/totipo-conformance` | [no test files] | [no test files] |
| `internal/cryptov1` | ok, 0.589s | ok, 2.293s |
| `internal/graph` | ok, 0.009s | ok, 1.060s |
| `internal/object` | ok, 0.009s | ok, 1.041s |
| `internal/storage` | ok, 0.567s | ok, 2.172s |
| `internal/tlv` | ok, 0.004s | ok, 1.020s |
| `internal/totp` | ok, 0.007s | ok, 1.033s |
| `internal/vectors` | ok, 1.184s | ok, 3.537s |

`go -C conformance vet ./...`: PASS, exit 0, no diagnostics.
`make race`: PASS, all packages, `-race -count=1 ./...`.

`make fuzz` preserved the Makefile policy: each target used `-run '^$'`,
`-fuzztime 10s`, `-parallel 2`. Exact emitted final accounting:

| Target/package | Initial coverage | Final executions | New interesting / total | Package duration/result |
| --- | --- | ---: | --- | --- |
| `FuzzDispatch` / `internal/object` | 41/41 | 1,124,598 | 5 / 46 | 10.073s, PASS |
| `FuzzOpen` / `internal/cryptov1` | 15/15 | 686,527 | 1 / 16 | 11.014s, PASS |
| `FuzzArrivalAndDisappearance` / `internal/graph` | 121/121 | 164,773 | 5 / 126 | 10.112s, PASS |

The cryptov1 final progress line reported 11s elapsed under the unchanged 10s
bounded policy. No fuzz failure or reproduction input was produced; seeds and
failure semantics are unchanged. Formatting (`test -z "$(gofmt -l conformance)"`),
current-document local links, report whitespace and `git diff --check` passed.

| Execution accounting | This wave |
| --- | --- |
| Baseline `make check` | 0 (infrastructure/source cleanup) |
| Focused Python / Go test invocations | 0 / 0 |
| Focused fuzz: Dispatch / Open / graph | 0 / 0 / 0 |
| Final complete suites / incomplete attempts | 1 / 0 |
| Final `make check` / vet / race / fuzz | 1 / 1 / 1 / 1 |
| Each fuzz target within final `make fuzz` | 1 |
| Generator / historical audit execution | 0 / 0 |
| Human Nix / agent Nix | 1 / 0; human-reported PASS / agent not run |
| Fresh remote CI / cached Nix observations | 0 / none |

## 10. Generator non-run statement

No vector generator or regeneration command was run. Normal compilation/testing
of the generator command package is not execution of its generator entry point.
No expected output, schema, pin, vector, or protocol semantic change was made.

## 11. Nix source-filter effect

Static inspection of flake.nix confirms `qualificationSource` selects all
`tools/*.py` through `file.hasExt "py"`. Removing these two files therefore
changes the filtered source input for qualification and dependency derivations,
invalidating prior Nix qualification for the final tree. The filter was not
narrowed or refactored.

Nix executes vendoring, formatting, `make check`, vet, race and fuzz; Make invokes
the active structural/schema checkers and `test_*.py` discovery. Neither audit
is invoked or discovered as a test. Deletion removes source bytes, not an executed
Nix command. CONTRIBUTING and this report are excluded from qualificationSource
and are not consumed by current structure checks; archived revision history is
included and remains unchanged. All executable/filter inputs are frozen.

## 12. Human Nix status

**PASS — human-reported after handoff.** The user reported “nix flake check
passed” for the frozen Wave 1 inputs. Human execution count: 1; agent execution
count: 0. This records the human's result, not an independently observed agent run.
After the passing host boundary and frozen inputs, the handed-off human command was:

```sh
nix flake check path:.
```

AGENTS.md requires: “The agent does not run Nix.” Its “Nix last, human only” rule
requires this final routine pinned-environment gate after source inputs freeze.
No Nix command was run by the agent and no cached result was inferred. This
status/count update changes only this excluded report; no post-Nix executable or
qualification-source edits occurred. A later included-source edit would
invalidate the result.

## 13. Remote CI disposition

No commit/push or workflow dispatch occurred. Remote CI for this wave is **NOT
RUN / not claimed**. After the user's later commit/push, ordinary authoritative
Linux Nix CI and the unchanged supplemental portability matrix provide remote
evidence. The matrix was not reproduced locally.

## 14. Changed files

- Modified `CONTRIBUTING.md`: historical-tool retrieval policy and exact revisions.
- Deleted `tools/audit_r15.py` and `tools/audit_r16.py`.
- Added `review/SPEC_SUPPORT_TOOLING_WAVE1_REPORT.md`: this report.

## 15. Final Git state

Branch remains `main`; HEAD remains `0591c59753d78c168e1214b32430ffd0d2f934ef`.
Index is unchanged. All work is unstaged/uncommitted:

```text
 M CONTRIBUTING.md
 D tools/audit_r15.py
 D tools/audit_r16.py
?? review/SPEC_SUPPORT_TOOLING_WAVE1_REPORT.md
```

No stage, commit, push, tag, release, migration or semantic refactor occurred.

## 16. Suggested commit subject

```text
Retire obsolete Spec migration audit tools
```
