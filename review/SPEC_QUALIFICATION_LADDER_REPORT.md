# Specification qualification ladder report

## 1. Starting HEAD/revision

Branch: `main`. Starting HEAD: `31b79c9d568fe402026035d27c13832ffc7ad243`
(`Unify Spec Nix and CI`). Initial `git status --short` was empty. Current
specification and moving requirements are v1/r19, with 92 manifest cases.
No protocol revision or normative content changed.

## 2. Existing expensive qualification

Inspected README.md, CONTRIBUTING.md, Makefile, flake.nix,
.github/workflows/conformance.yml, tools/check_spec.py and its tests, and
review/SPEC_NIX_CI_UNIFICATION_REPORT.md. No AGENTS.md existed.

The Makefile owns `make check`, all-package race tests, and three bounded fuzz
targets. The Linux flake gate performs formatting, make check, vet, race, and
fuzz. Supplemental CI retains macOS/Windows Go 1.23.x/stable and Linux Go 1.23.x.
Neither the targets nor the CI matrix changed.

The committed Nix/CI report records passing host and isolated-source suites:
92 manifest cases, 14 Python tests, seven Go test packages (two command packages
without tests), vet, race, and all three fuzz targets. That report explicitly
leaves human Nix pending and remote CI unrun. These are inherited observations,
not executions or fresh passes from this milestone; no Nix success is inferred.

## 3. Milestone classifications

This milestone: **DOCS_ONLY**. Only AGENTS.md, README.md, CONTRIBUTING.md, and this
new report change. The ladder requires all applicable classifications in every
substantial report:

| Classification | Scope |
| --- | --- |
| DOCS_ONLY | Documentation/process |
| SPEC_TEXT | Specification text |
| REQUIREMENTS | Profiles and pins |
| VECTOR_CORPUS | Cases, schemas, manifest, expected outcomes |
| CONFORMANCE_GO | Go consumer and tests |
| PYTHON_VERIFIER | Python verification and tests |
| GENERATOR | Fixture generation and tests |
| BUILD_NIX_CI | Build, dependencies, Nix, filtering, CI |

## 4. Baseline policy

Ordinary semantic/conformance work runs `make check` once at baseline. Prior
committed qualification evidence is cited with its actual status. No default
baseline race/fuzz; directly affected race/fuzz-sensitive code gets focused
coverage. This docs-only milestone does not require baseline make check.

## 5. Focused inner-loop policy

Use specific Python unittest modules/tests, affected Go packages/tests, exact
relevant supported manifest/test commands, or the affected bounded fuzz target.
Avoid repeating full check/vet/race/fuzz/Nix after each edit. Preserve failures
and reproduction inputs, reproduce narrowly, fix and rerun the focused target,
then return to the final suite after stabilization. A deterministic test failure
does not justify restarting every fuzz target during diagnosis.

## 6. Generator/verification separation

Generation is maintenance evidence, never an inner-loop replacement for
verification. Generator changes use applicable temporary generation/comparison
and reviewed differences, followed by final normal corpus verification. No
committed vectors or expectations may be accepted solely because generation
succeeds. This milestone ran no generator and changed no generator semantics.

## 7. Final full suite

For stable semantic/conformance inputs, run once:

```sh
make check
go -C conformance vet ./...
make race
make fuzz
git diff --check
```

Record exact manifest count, Python tests, Go package results, vet, race, and all
fuzz targets/durations. Preserve applicable formatting checks. Report-only edits
outside qualification inputs do not require repetition. No complete host suite
was needed or run for this milestone.

## 8. Docs-only policy

With spec/, vectors/, requirements/, conformance/, tools/, and Makefile unchanged,
use protected-byte comparison, actual applicable documentation/static checks,
and git diff --check. Do not automatically run race or 30+ seconds of fuzz for
prose. Inspect actual structure consumers and the Nix source filter first.
README and CONTRIBUTING now link to AGENTS.md and explain the boundaries.

## 9. Invalidation matrix

| Change | Invalidation |
| --- | --- |
| Spec/vector/requirements/conformance | Final full host suite; Nix for included source |
| Go/Python test or verification logic | Appropriate final suite; included source invalidates Nix |
| Generator | Applicable generation comparison, final corpus verification, full host suite and Nix for included source |
| Nix/vendorHash/lock/flake-source | Final Nix; affected host checks as applicable |
| Historical report excluded from qualification | No test/fuzz/Nix invalidation |
| Review/history consumed by structure tests/source filter | Affected qualification, including Nix if included |
| README/AGENTS/process docs | Inspect actual consumers/filter; do not guess |

## 10. Nix-last rule

The final routine pinned gate remains human `nix flake check path:.`, requested
only after every qualification source is frozen, including deliberately consumed
history and tooling. Later flake-source edits invalidate it. The agent runs no
Nix. There is no `nix build` or artificial package materialization requirement.

## 11. Race/fuzz preservation

Optimization must come from staging expensive checks at boundaries,
never from deleting final conformance/race/fuzz coverage.

Final semantic/conformance qualification still requires `make race` and
`make fuzz`. Unchanged Makefile policy:

| Fuzz target | Package | Duration / parallelism |
| --- | --- | --- |
| FuzzDispatch | internal/object | 10s / 2 |
| FuzzOpen | internal/cryptov1 | 10s / 2 |
| FuzzArrivalAndDisappearance | internal/graph | 10s / 2 |

Each retains `-run '^$'`, existing seeds, and failure semantics. Race retains
`-race -count=1 ./...`. No coverage was removed and no CI portability change made.

## 12. Execution-count convention

Reports distinguish baseline make-check runs, focused Python/Go test invocations,
focused fuzz by target, complete final suites and incomplete attempts, and human
Nix runs/status. Clean semantic target: baseline make check 1, final complete
check/vet/race/fuzz 1, human Nix 1. Failures/repeats are recorded, not hidden.

This milestone:

| Execution | Count |
| --- | --- |
| Baseline make check | 0 |
| Focused Python unit-test invocations | 0 |
| Focused Go test invocations | 0 |
| Focused fuzz runs (each target) | 0 |
| Final complete host suites / incomplete attempts | 0 / 0 |
| Vet / race / full fuzz | 0 / 0 / 0 |
| Human Nix during this milestone | 0; not required for these inputs |
| Agent Nix / generator | 0 / 0 |

Static documentation and byte checks are not Python unit-test executions.

## 13. Protected-byte invariance

A pre-edit SHA-256 inventory covers 163 tracked protected/build/history files.
Final comparison verifies identical bytes and protected path membership:
131 files under spec/, vectors/, requirements/, conformance/, tools/, plus
go.work; 28 existing historical review files; Makefile, flake.nix, flake.lock,
and .github/workflows/conformance.yml. Zero changed, added, or deleted protected
files. All original history remains untouched; this report is the only addition
under review/.

Root go.mod and go.sum were absent at baseline and remain absent. The actual
conformance/go.mod and conformance/go.sum are included in the byte comparison.
Protocol, requirements, vectors, expected outcomes, Go/Python implementation,
generator, and qualification tooling remain byte-identical.

## 14. Validation

PASS: protected SHA-256 and path-membership comparison; local Markdown link
existence and balanced code-fence checks for all four changed/new documents;
`git diff --check`, plus whitespace checks for new untracked Markdown files.
Manual review checked the ladder against the unchanged Make targets, source
filter, structural consumers, CI, and existing report. No existing structure
check consumes these four documents, so none requires rerunning for this edit.
No race/fuzz/full-suite execution was substituted for byte evidence.

## 15. Nix disposition

No new human Nix request is required. The actual qualificationSource includes
spec/, vectors/, requirements/, conformance Go/module files, tools Python,
Makefile, go.work, and review/V1_PRE_R16_REVISION_HISTORY.md. The four changed/new
documents are excluded and are not read by current structure checks. Nix files,
lock, vendorHash, and all included source bytes remain unchanged. Therefore this
milestone does not invalidate existing Nix evidence. The previous report's human
gate remains pending; this milestone neither completes nor restarts it.
Remote CI was not run; it remains authoritative supplemental portability evidence
after commit/push.

## 16. Final Git state

Branch and HEAD unchanged: main at 31b79c9d568fe402026035d27c13832ffc7ad243.
Modified tracked files: README.md and CONTRIBUTING.md. New untracked files:
AGENTS.md and review/SPEC_QUALIFICATION_LADDER_REPORT.md. The index is unchanged;
everything is unstaged/uncommitted. No commit, push, or workflow dispatch occurred.
