# Interactive qualification ladder

Start each milestone on clean `main`; record the starting HEAD and protocol
revision. Preserve historical reports. This ladder governs interactive agent work;
it does not change the protocol, requirements, corpus, expected outcomes,
generator semantics, conformance implementation, Make targets, or CI matrix.
Remote CI remains authoritative supplemental portability evidence after commit/push.

Optimization must come from staging expensive checks at boundaries,
never from deleting final conformance/race/fuzz coverage.

## Classify before qualifying

Record all applicable classifications in each substantial report:

| Classification | Changed inputs |
| --- | --- |
| `DOCS_ONLY` | Documentation/process without qualification source changes |
| `SPEC_TEXT` | Specification text, including editorial changes |
| `REQUIREMENTS` | Requirement profiles or pins |
| `VECTOR_CORPUS` | Cases, schemas, manifest, or expected outcomes |
| `CONFORMANCE_GO` | Go consumer or its tests |
| `PYTHON_VERIFIER` | Python verification/structure logic or its tests |
| `GENERATOR` | Fixture generator or its tests |
| `BUILD_NIX_CI` | Build, dependency, Nix, source-filter, or CI inputs |

Use multiple classifications for mixed milestones. An infrastructure-only
milestone must establish source/byte invariance; changes to qualification logic
still require the appropriate final qualification.

## Baseline once

For ordinary semantic/conformance work, run `make check` once at baseline.
The committed starting revision carries previous qualification evidence; cite
its actual status without upgrading pending evidence to a pass. Do not run race
or fuzz at baseline by default. If the milestone directly changes race/fuzz-
sensitive code, run the relevant focused target.

## Focused inner loop

Use the narrowest relevant command while editing, for example:

```sh
python3 -m unittest <specific module/test>
go -C conformance test ./internal/graph
go -C conformance test ./internal/cryptov1
```

For a specific conformance case, use the exact relevant manifest/test command
supported by the consumer. For fuzz-target development, run only the affected
bounded target using the existing Makefile flags. Do not repeatedly run
`make check`, vet, race, all three fuzz targets, or Nix after minor edits.

Generation and verification are separate. Never regenerate committed vectors
and accept them merely because generation succeeds. Generator changes need
applicable temporary generation/regenerate-and-compare evidence, review of
intended differences, and final normal corpus verification. Successful generation
cannot substitute for verification or justify changing expectations to pass.

If `make check`, race, or fuzz fails, preserve the failure output and any
reproduction input. Reproduce the specific Python test, Go package/test, or fuzz
target; fix narrowly and rerun that focused target. Return to the complete final
suite after stabilization. Do not restart all three fuzz targets while diagnosing
one deterministic test failure.

## Final host boundary

When spec/conformance inputs are stable, run the complete non-Nix suite once:

```sh
make check
go -C conformance vet ./...
make race
make fuzz
git diff --check
```

Keep any applicable formatting/structure checks required by changed inputs.
Record the exact manifest case count, Python test count/results, Go package
results, vet result, race result, and each fuzz target's duration/result.
Final race and fuzz remain mandatory for semantic/conformance code changes.
Do not repeat this suite after report-only edits outside qualification inputs.

The bounded fuzz policy remains exactly the Makefile policy: `FuzzDispatch`
(`internal/object`), `FuzzOpen` (`internal/cryptov1`), and
`FuzzArrivalAndDisappearance` (`internal/graph`), each with `-run '^$'`,
`-fuzztime 10s`, and `-parallel 2`. Preserve seeds and failure semantics.
`make race` retains all-package race coverage. The ladder changes when these
checks run, never their coverage or meaning.

## Docs-only profile and byte invariance

When `spec/`, `vectors/`, `requirements/`, `conformance/`, `tools/`, and `Makefile`
are unchanged, use documentation/static checks required by the actual changed
inputs and `git diff --check`. Do not automatically run race or 30+ seconds of
fuzz because prose changed. A docs-only milestone does not need a baseline
`make check` or the full host suite unless actual qualification dependencies
consume the changed documents.

For infrastructure/docs-only work, explicitly compare bytes against the starting
revision for `spec/`, `vectors/`, `requirements/`, `conformance/`, relevant `tools/`,
`go.work`, and every `go.mod`/`go.sum` (record absent root files). Include path
additions/deletions in the comparison. Verify unchanged build/filter inputs when
relying on existing Nix evidence. This is more useful than rerunning fuzz for
pure prose edits.

Inspect the actual source filter and structure-test dependencies before deciding
whether docs invalidate qualification. Currently `flake.nix` includes spec,
vectors, requirements, conformance `.go`/`go.mod`/`go.sum`, tools `.py`, Makefile,
go.work, and `review/V1_PRE_R16_REVISION_HISTORY.md` in `qualificationSource`.
README, CONTRIBUTING, AGENTS, and other review reports are excluded and are not
consumed by current structure checks. Reinspect this when inputs/filter change.

## Invalidation matrix

| Later change | Evidence invalidated / action |
| --- | --- |
| Spec, vectors, requirements, or conformance | Final full check/vet/race/fuzz suite; Nix when included in its source |
| Go/Python tests or verification logic | Appropriate final suite, since qualification logic changed; Nix for included source |
| Generator | Applicable generation comparison and final normal corpus verification; full host suite and Nix for included source |
| Nix, vendorHash, lock, or flake-source inputs | Final Nix; affected host checks for changed qualification inputs |
| Historical report excluded from qualification | No test/fuzz/Nix invalidation |
| Review/history consumed by structure tests or source filter | Rerun affected qualification, including Nix if a flake-source input |
| README/AGENTS/other process docs | Inspect actual filter and consumers; no guessed invalidation |

## Nix last, human only

The agent does not run Nix. The final routine pinned-environment gate is the
human command:

```sh
nix flake check path:.
```

Request it only after all flake-source inputs are frozen: Go/Python source,
vectors/spec/requirements, Nix files/lock, deliberately included review/history,
and qualification tooling. Any later edit to a flake-source input invalidates
that result. Docs-only edits outside those inputs do not require a fresh request.
An older pending gate remains pending, not passed. No `nix build` or artificial
package materialization is required. Preserve Linux qualification and the
existing supplemental portability matrix.

## Report execution counts

State baseline `make check` count, focused Python/Go test invocation counts,
focused fuzz runs by target, final complete suite count (and any incomplete
attempts), and human Nix count/status. Distinguish inherited evidence, agent
execution, human execution, cached Nix results, and remote CI evidence.

Clean semantic milestone target: baseline `make check`: 1; final
`make check`/vet/race/fuzz: 1; human Nix: 1. Docs-only milestones may correctly
record zero for all of these. Record exceptions and failures rather than hiding
repeated attempts. Leave staging/committing to the user's requested workflow.
