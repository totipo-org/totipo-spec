# Tooling simplification implementation report

Date: 2026-10-14. Classifications: **BUILD_NIX_CI**, **CONFORMANCE_GO**,
plus documentation corrections. This is not a docs-only qualification.

## Starting state and inherited evidence

Repository: `/home/niki/Sources/totipo-spec`; branch `main`.
Starting full HEAD: `e2bb221cad1e1cb3432df3d7278e6eb7d3854c30`
(`tooling simplification audit`). Starting `git status --short` was empty;
the audit was tracked and committed at that HEAD. AGENTS.md was read completely.
No sibling repository, isolated worktree, or delegated agent was used.

Active specification: protocol 1, v1/r19, design draft; moving-pre-RC profile.
Manifest and physical case inventory: **92** each; no revision/profile change.
`totipo/generator`, `totipo/conformance`, and go.work retain Go **1.23.0**.
Generator has no external module requirements. Conformance retains
`golang.org/x/crypto v0.36.0` and indirect `golang.org/x/sys v0.31.0`.
Host: `go version go1.26.8 linux/amd64`.

The committed audit cites W3D passing host qualification and one human-reported
Nix PASS on 2026-10-12, with cached/fresh provenance unknown and exact-commit
portability CI pending. Those are inherited historical observations, not passes
for this milestone. Remote status was not queried or upgraded.

Baseline `make check` ran once before edits and passed: architecture, spec policy,
three generator packages, thirteen tested conformance packages (plus the
check-architecture command with no tests), full 92-case conformance, then
verify-only acceptance of 92 case files. No baseline race or fuzz was run.

## Implemented scope and invocation reduction

Option A only. Nine existing files changed, plus this new report. No optional
Hash test deletion, Nix tool-list factoring, or Option B change was implemented.

Make before:

```text
check: architecture spec-check test conformance verify
```

Make after:

```text
check: architecture spec-check test conformance
```

Both sequences run architecture and specification policy commands, then
`go -C generator test -count=1 ./...`,
`go -C conformance test -count=1 ./...`, then full conformance.
The old sequence additionally ran the CLI with `-verify-only`.
`make -n check` confirms the final five commands and no generator main invocation.
`generate` remains a separate deliberate writer target, with no dependency edge
from qualification. Make race, fuzz, public verify, and all other recipes remain
unchanged. Formatting and both module vets remain separate final host/Nix gates.

The exact Make allowlist remains in `repopolicy/wiring.go`; only its expected
check prerequisites changed. Only the two affected wiring mutation fixture rows
changed. All thirteen wiring negatives remain, including hidden generation,
architecture removal, and removal of full conformance/artifact verification.
No policy rule was broadened or new Make parser introduced.

Portability CI now runs `check-spec` alone in the renamed early `Spec policy`
step, and omits the trailing verify-only integrity step. Full semantic CLI
execution remains mandatory. Linux Nix, all five supplemental portability
combinations, formatting, both module vets/tests, standalone module tests,
architecture checks, and immutable Action references are unchanged.

Executable/configuration savings: one fewer CLI process per `make check`, two
fewer per portability job (ten across the five-job matrix). No runtime speedup
was measured. Workflow shrank by four physical lines. Make and mirrored policy
line counts are unchanged. Across the four code/config files: 6 inserted,
10 deleted, net **4 fewer physical lines**. Across all nine existing files:
30 inserted, 26 deleted, net **4 more lines**, due to clearer documentation.
This report is additional review evidence, excluded from those reduction counts.

## Independent implementations and verify-only behavior

Generator source, generator tests, semantic consumer packages, oracle tests,
CLI source/tests, schemas, cases, expected outcomes, and profile are unchanged.
The two Go modules retain separate models, constants and Totipo composition.
No shared helper, cross-role import, delegation, or semantic library was added.
No generator execution or vector regeneration occurred.

Full CLI still calls `vectors.Read` and `vectors.VerifyProfile` before either
mode branches. Schema, manifest, case inventory/hash, profile and revision checks
therefore remain mandatory before semantic execution. `make verify` remains
operational and passed independently with 92 case files.

Focused command (one invocation, two packages):

```sh
go -C conformance test -count=1 -v ./internal/repopolicy ./cmd/totipo-conformance
```

Both packages passed. Existing tests demonstrated valid acceptance in both modes;
manifest, case-member and schema-source refusal in both modes; wrong profile
revision refusal; malformed/missing input and argument refusal; and read-only
path/byte snapshots for successful and failed checks. The deliberately repinned
wrong TOTP answer passed verify-only and failed full semantic execution, as
expected. Wrong profile identity refusal is explicitly tested in verify-only;
the unchanged common pre-branch VerifyProfile call protects full mode too.

All 32 source/metadata policy mutation subtests and thirteen wiring mutation
subtests passed by rejecting the prohibited changes. Existing closure, symlink,
vendor containment, oracle import and generator/verifier coupling negatives also
passed. Final all-package tests additionally cover artifact/profile relationships,
manifest tampering, protocol/crypto, graph/fold/storage, schemas, and the independent
wire oracle. Expected rejections are successful tests, not qualification failures.

## Documentation corrections

- `vectors/FORMAT.md`: fixes the generator command to
  `go run ./generator/cmd/generate-vectors -root .`; clarifies independent
  composition, shared standard crypto primitives, authoritative committed input,
  and artifact verification in both CLI modes. Only deliberate-maintenance prose
  changed; the artifact format description and serialized schemas did not.
- `CONTRIBUTING.md`: removes the superseded stale-path workaround.
- Root README, conformance README, and requirements README: clarify separate
  modules, independent semantics, full conformance's included verification,
  deliberate generation, and verify-only use without semantic execution.
  The maintenance example uses generation followed by full conformance.
- AGENTS.md was inspected and remains unchanged; its qualification ladder and
  human-owned Nix policy remain applicable.

## Protected identities and dependency accounting

Compared Git blob identities from starting HEAD with `git hash-object` of working
files, including inventory checks for additions/deletions. Of the 100 tracked
spec/vectors/requirements paths, exactly two prose paths changed:
`vectors/FORMAT.md` and `requirements/README.md`. The other **98 blobs** are
identical: all 92 case files, spec, manifest, both schemas, moving profile, and
vectors README. No protected path was added or deleted.

Reproducible identity inventory command:

```sh
git ls-tree -r e2bb221cad1e1cb3432df3d7278e6eb7d3854c30 -- spec vectors requirements |
  awk '$4 != "vectors/FORMAT.md" && $4 != "requirements/README.md"'
```

SHA-256 of that complete 98-row Git blob/path inventory (including final newline):
`340a875aeff05e32f2ac0c51fa0b09a43706741bf19ce3012f4bd79c3ef9467b`.
Each listed starting blob was individually matched against its final file.
Key starting = final Git blob identities:

| Artifact | Git blob |
| --- | --- |
| Active specification | `66478873fdbb753d5c2248e3ae1e87d6e19cd10f` |
| Manifest | `e6ed939ef1f1894161b40b63820dd18d623df4c9` |
| Manifest schema | `37b662aace94e51452d83f22222764c37675deb3` |
| Case schema | `f960103f8898022e14ee5ba3514a49ac4a7b1196` |
| Requirements profile | `e88eb3869dd8fdaa903caf0c94066d3715bfee8c` |
| generator/go.mod | `c3d1f1479eb6bcf65eb6d38f87d2d067f209d3cb` |
| conformance/go.mod | `f380167a10486852baa71a3b4c7622b2bed9e2c9` |
| conformance/go.sum | `927b25585c3f8fe1f3c1fc68078350081e9fcba1` |
| go.work | `1c01d150bd3662c457e080ea96898cf2fb09a979` |
| flake.nix | `1100eb7fd467df43c1a2c03ae7e620b94d6bc628` |
| flake.lock | `be8c1fae2544af340b328d4dc047df5619af251c` |

All generator files are identical. Within conformance, only README and the two
listed wiring files differ. No tools directory exists. Root go.mod, go.sum,
go.work.sum and generator/go.sum remain absent. All module/dependency checksums,
Nix lock inputs, vendorHash, release/version/profile pins and external Action
references remain unchanged. No cache regeneration or package materialization
was requested. Existing tag inventory remains `v0-rc1`; no tag action occurred.
Historical review reports remain unchanged.

## Final host qualification

Qualification inputs were stabilized before the final suite. Results are recorded
below after execution. Logs live only under `/tmp/totipo-simplification-*.log`;
no build log, generated output, cache, or temporary fixture is Git-known.

The single final complete suite passed, with no incomplete or superseded attempt:

```sh
make check
go -C generator vet ./...
go -C conformance vet ./...
test -z "$(gofmt -l generator conformance)"
make race
make fuzz
go run ./conformance/cmd/totipo-conformance -root .
git diff --check
```

Both module test suites are executed by `make check`; no duplicate all-package
host invocation was needed. Architecture and v1/r19 spec policy passed. Full
conformance passed **92 cases** inside make and again in the independent CLI.
Both vets exited 0; formatting returned no paths; whitespace check passed.
Python test invocations/executions: **0/0** (no current Python verifier).

| Package (module-relative) | Final tests | Final race |
| --- | --- | --- |
| generator: cmd/generate-vectors | PASS 0.002s | PASS 1.012s |
| generator: internal/fixture | PASS 0.002s | PASS 1.010s |
| generator: internal/wire | PASS 0.014s | PASS 1.102s |
| conformance: cmd/check-architecture | No test files | No test files |
| conformance: cmd/check-spec | PASS 0.010s | PASS 1.093s |
| conformance: cmd/totipo-conformance | PASS 3.615s | PASS 10.440s |
| conformance: internal/cryptov1 | PASS 0.942s | PASS 2.073s |
| conformance: internal/graph | PASS 0.008s | PASS 1.054s |
| conformance: internal/object | PASS 0.008s | PASS 1.024s |
| conformance: internal/repopolicy | PASS 1.168s | PASS 2.408s |
| conformance: internal/schema | PASS 0.928s | PASS 4.377s |
| conformance: internal/specpolicy | PASS 0.856s | PASS 5.488s |
| conformance: internal/storage | PASS 0.825s | PASS 2.083s |
| conformance: internal/tlv | PASS 0.012s | PASS 1.017s |
| conformance: internal/totp | PASS 0.005s | PASS 1.022s |
| conformance: internal/vectors | PASS 2.098s | PASS 6.588s |
| conformance: internal/wireoracle | PASS 0.076s | PASS 1.270s |

Race retained both all-package commands with `-race -count=1 ./...`.
Fuzz retained the exact existing Make commands, each with
`-run '^$' -fuzztime 10s -parallel 2`:

| Target / package | Result | Reported package duration | Executions |
| --- | --- | --- | --- |
| FuzzDispatch / internal/object | PASS | 10.084s | 1,184,709 |
| FuzzOpen / internal/cryptov1 | PASS | 11.014s | 634,323 |
| FuzzArrivalAndDisappearance / internal/graph | PASS | 10.535s | 174,270 |

Go's bounded-run completion accounts for elapsed time beyond the nominal 10s.
Starting fuzz coverage was 53/18/155 inputs respectively; graph discovered three
additional interesting inputs in the local Go fuzz cache. No failure or
reproduction artifact was produced; committed seeds and failure semantics remain
unchanged. Local caches are outside Git membership.

Execution accounting:

- Baseline make check: **1**.
- Focused Go test invocations: **1**, covering two packages; focused Python: **0**.
- Focused fuzz runs Dispatch/Open/ArrivalAndDisappearance: **0/0/0**.
- Final complete host suite: **1**; incomplete attempts: **0**.
- Final make check/vet/race/fuzz: **1/2 module vet commands/1/1**.
- Final module test commands within make: **2**, plus **2** race module commands.
- Independent public make verify: **1**; independent full CLI: **1**.
- Make dry-run inspection: **1**. Generator executions: **0**.
- Agent Nix: **0**; human Nix for this milestone: **1, owner-reported PASS**.
- Nix cached/fresh provenance: **not supplied by owner**. Remote CI: not run or queried.

There were no test failures. One initial inventory command used exclude pathspec
syntax unsupported by `git ls-tree`; it was corrected to the explicit awk filter
shown above and all 98 blobs were compared successfully. An interim log read
occurred before fuzz had started; final completed logs supply the results above.
Neither read-only inspection issue was a qualification attempt or changed inputs.

## Human Nix gate and source freeze

**Owner-reported PASS: `nix flake check`, supplied on 2026-10-14.**
The owner reported one successful human execution for this milestone. Detailed
build output and cached/fresh provenance were not supplied. The frozen source
SHA-256 inventory was rechecked after that report and passed; qualification inputs
are unchanged. Only this excluded report was updated to record the result.
AGENTS.md and the user explicitly reserve this command for the human. No agent
Nix command, separate build, dependency update, or cache regeneration occurred.
The workflow's Linux Nix job remains authoritative after owner commit/push;
supplemental portability evidence for these uncommitted changes is pending.

The actual flake filter was inspected: spec, vectors and requirements directories,
all repository Go/native/module inputs, Makefile, go.work, and archived pre-r16
history are included. This milestone changes included Make/policy inputs and
vectors/requirements prose; its required human Nix gate is now recorded as passed.
The flake and lock themselves are unchanged. The new report, root README,
CONTRIBUTING, conformance README and workflow are outside qualificationSource.
Specpolicy reads only active spec and archived history; repository policy inspects
source/module/Make inputs, not the new report. No changed documentation is consumed
by current structure tests.

Before final qualification, a SHA-256 inventory was captured for all tracked
spec/vectors/requirements/generator/conformance inputs, Makefile, go.work,
flake.nix, flake.lock and archived pre-r16 history. Rechecking that inventory after
qualification passed. No qualification input was edited after the final suite;
only this excluded report was completed. The source freeze covers every changed
qualification input. Later included-source edits invalidate the affected evidence.

The new report is Git-known with `git add --intent-to-add` only. All existing
modifications remain unstaged. Final `git diff --cached --quiet` exits 0;
`git diff --check` exits 0. Nonignored untracked inventory is empty. Exactly the
nine intended existing paths and this report differ; there are no deletions or
unintended additions. Local build/cache/log state has not been added.

Final branch: `main`; full HEAD remains
`e2bb221cad1e1cb3432df3d7278e6eb7d3854c30`.
Full final `git status --short`:

```text
 M .github/workflows/conformance.yml
 M CONTRIBUTING.md
 M Makefile
 M README.md
 M conformance/README.md
 M conformance/internal/repopolicy/wiring.go
 M conformance/internal/repopolicy/wiring_test.go
 M requirements/README.md
 A review/SPEC_TOOLING_SIMPLIFICATION_REPORT.md
 M vectors/FORMAT.md
```

The leading blank in ` A` denotes unstaged new content with intent-to-add;
no file content is staged. All changes remain uncommitted. No commit, tag, push,
release, GitHub settings change, or work outside this repository occurred.

The human Nix gate is complete. Supplemental remote CI remains pending a later
owner-authorized commit/push; no such action is part of this milestone.
No host failure or implementation issue remains. Stop at this bounded Option A
implementation; no Desktop or Android work is part of this milestone.

## Correction: Nix-first qualification and removal of Make

Date: 2026-10-14. This section supersedes the initial Option A implementation's
current-state conclusions above while preserving its implementation and execution
history. The owner's new instruction explicitly approves Nix-first/no-Make
qualification, deletion of the exact Make validator, and updated agent guidance.
This is a correction to the same uncommitted simplification, not a new milestone.
Classifications remain **BUILD_NIX_CI** and **CONFORMANCE_GO**, plus documentation.
No protocol, generator or serialized corpus change is authorized or implemented.

### Actual correction baseline

Case **B** applied: branch `main`, full HEAD
`e2bb221cad1e1cb3432df3d7278e6eb7d3854c30`, with only the previously qualified
uncommitted simplification. Starting status:

```text
 M .github/workflows/conformance.yml
 M CONTRIBUTING.md
 M Makefile
 M README.md
 M conformance/README.md
 M conformance/internal/repopolicy/wiring.go
 M conformance/internal/repopolicy/wiring_test.go
 M requirements/README.md
 A review/SPEC_TOOLING_SIMPLIFICATION_REPORT.md
 M vectors/FORMAT.md
```

No content was staged. The changed path list matched the prior intended inventory,
and all prior frozen qualification inputs matched their recorded SHA-256 values.
There were no unrelated changes. The original report was saved locally before
appending this section; its original text is preserved above. Historical committed
review reports are unchanged. No sibling repository, worktree or sub-agent was used.

Baseline evidence is the earlier host PASS and owner-reported Nix PASS recorded
above. No additional baseline suite was run; this deliberately reuses the qualified
uncommitted baseline authorized by the owner. The earlier Nix PASS is **invalidated
for the corrected implementation**, retained only as historical evidence.
Protocol v1/r19, 92 cases, Go module versions, host Go version and all dependency
identities remain as recorded above.

### Deleted orchestration and retained isolation

Deleted completely:

| File | Physical lines at correction baseline | Purpose |
| --- | --- | --- |
| Makefile | 40 | Public task recipes and dependencies |
| conformance/internal/repopolicy/wiring.go | 95 | Exact Make text allowlist and comparison |
| conformance/internal/repopolicy/wiring_test.go | 45 | Thirteen wiring mutation subtests and harmless-text test |

Removed the three-line Wiring call from `policy.go` and the one-line Makefile
construction fixture from `policy_test.go`. Total removed across these five paths:
**184 physical lines**. Two wiring test functions and thirteen wiring mutation
rows are retired. No semantic isolation test was removed. Nix adds twelve lines
and removes nine (net +3), so code/config reduction for this correction is
**181 physical lines**. The existing workflow's earlier net four-line reduction
is inherited; combined code/config reduction against committed HEAD is 185 lines.
Documentation/report line changes are separate from those implementation totals.

The retained Check sequence is Source → Metadata → Dependencies. All 32 existing
source/metadata mutation subtests remain. Retained restrictions cover public and
internal cross-role imports, shared/third modules, replacements and workspace
bridges, external semantic dependencies, oracle implementation dependencies,
process delegation, native/plugin/linkname bridges, and dependency containment.
Exact module/dependency metadata policy remains unchanged. No orchestration
allowlist, parser, checksum guard, new runner, script, Nix app, or workflow was
introduced. The retained source scan checks implementation isolation only.

The removed wiring guard's accidental-orchestration regression protection is now
owned by review of ordinary Nix/Actions configuration. Local Go tests do not prove
provider behavior or independently enforce future orchestration edits.

### Nix and CI commands

Before, the qualification build phase ran formatting, `make check`, both module
vets, `make race`, and `make fuzz`. After, the same existing derivation directly
runs the following sequence (after its unchanged vendoring/cache/environment setup):

```sh
go version
test -z "$(gofmt -l generator conformance)"
go run ./conformance/cmd/check-architecture -root .
go run ./conformance/cmd/check-spec -root .
go -C generator test -count=1 ./...
go -C conformance test -count=1 ./...
go run ./conformance/cmd/totipo-conformance -root .
go -C generator vet ./...
go -C conformance vet ./...
go -C generator test -race -count=1 ./...
go -C conformance test -race -count=1 ./...
go -C conformance test ./internal/object -run '^$' -fuzz FuzzDispatch -fuzztime 10s -parallel 2
go -C conformance test ./internal/cryptov1 -run '^$' -fuzz FuzzOpen -fuzztime 10s -parallel 2
go -C conformance test ./internal/graph -run '^$' -fuzz FuzzArrivalAndDisappearance -fuzztime 10s -parallel 2
```

Full conformance performs artifact/schema/profile verification before semantics;
there is one top-level full acceptance invocation in routine Nix qualification.
CLI and corpus unit tests retain their behavior-focused verification coverage.
There is no separate redundant verify-only qualification process.

Removed the deleted Makefile from qualificationSource and removed explicit
gnumake entries from qualification nativeBuildInputs, the development shell and
both agent jail tool lists (four entries). All other source selection, pinned Go,
Nix inputs, fixed workspace vendoring, vendorHash, offline proxy/sumdb settings,
GOTOOLCHAIN=local, sandbox setup and writable temporary cache paths are unchanged.
No Nix command or dependency/cache update was executed by the agent. Nix evaluation
and sandbox execution were reserved for the new human gate; the owner's subsequent
PASS is recorded below, separately from agent source inspection.

The portability workflow already used direct Go after the initial simplification;
**zero further workflow edits** were necessary in this correction. Linux Nix and
all five portability combinations, standalone-module tests, both vets/tests,
formatting, spec/architecture checks, full CLI and immutable Action references
remain unchanged from the correction baseline. No Windows Nix requirement or
new platform race/fuzz job was added. The existing Linux Nix race/fuzz coverage
is preserved by the direct sequence above. CI/provider execution is unverified
until an owner-authorized commit/push.

### Generation boundary and active documentation

Confirmed generator main defines `-root` with default `.`. The deliberate writer
remains `go run ./generator/cmd/generate-vectors -root .`. No generation occurred.
Nix and workflow source have no generator main invocation; module tests compile
and test the writer without executing its main. No ordinary qualification path
invokes fixture generation. This is source review plus retained source restrictions,
not a new orchestration validator.

Updated AGENTS, root README, CONTRIBUTING, conformance README, requirements README,
vectors FORMAT and vectors README to use direct Go commands. The latter was an
additional active operational reference found during the repository-wide search.
AGENTS distinguishes focused agent Go checks, the complete direct host boundary,
final human Nix qualification and explicitly authorized generation. The supported
development environment is the Nix flake through the owner's direnv setup.
No operational Make reference remains outside historical review evidence.

### Protected identities

The active specification, all 92 cases, manifest, two schemas and moving profile
are byte-identical to both committed HEAD and the qualified correction baseline.
All generator source/tests and consumer semantic/CLI/oracle source/tests remain
unchanged. Only repopolicy's four-line caller/fixture removal and deleted wiring
files change Go source. No revision, outcome, schema contract, artifact pin,
module dependency, vendorHash, flake.lock, Action pin or release identity changed.
Root go.mod/go.sum/go.work.sum, generator/go.sum and tools remain absent.

The exact artifact and dependency blob identities in the original table above
still match. For current protected artifact inventory, exclude the three deliberately
edited prose files (vectors FORMAT, vectors README and requirements README) from
the original 100-path tree. All **97 remaining blobs** individually match their
starting Git identities, with no protected path addition/deletion:

```sh
git ls-tree -r e2bb221cad1e1cb3432df3d7278e6eb7d3854c30 -- spec vectors requirements |
  awk '$4 != "vectors/FORMAT.md" && $4 != "requirements/README.md" && $4 != "vectors/README.md"'
```

SHA-256 of this complete 97-row blob/path inventory, including final newline:
`367bd582246841223fc3158d59f07243a8778065c6523326a09df360fff9b7be`.
Each row was compared to `git hash-object` of its final working file. Manifest
entry count and physical case-file count are both **92**. No protected input was
repinned or regenerated to pass a check.

The earlier flake.nix source blob is intentionally superseded by this correction;
flake.lock and dependency/version identities remain identical. The three edited
artifact-directory Markdown files are prose exceptions, not serialized artifacts.

### Correction host qualification and execution counts

Focused execution (one test invocation plus four real CLI invocations):

```sh
go -C conformance test -count=1 -v ./internal/repopolicy ./cmd/totipo-conformance
go run ./conformance/cmd/check-architecture -root .
go run ./conformance/cmd/check-spec -root .
go run ./conformance/cmd/totipo-conformance -root .
go run ./conformance/cmd/totipo-conformance -root . -verify-only
```

All passed. All 32 retained source/metadata mutation negatives, three actual closure
refusals and symlink/vendor containment checks passed. The real architecture CLI
passed with no Makefile present. CLI tests retained valid/refusal behavior in both
modes, malformed/schema/profile identity refusal, semantic execution differences,
and successful/failed read-only artifact snapshots. Expected negative refusals
are successful tests. Focused package durations: repopolicy 0.259s; CLI 2.760s.
Full CLI and verify-only each accepted all 92 committed cases.

At the stable source boundary, the complete direct host sequence shown in the
Nix section above ran once, except `go version` was recorded during baseline
inspection, and with `git diff --check` added. No duplicate full conformance run
was added merely to replace a removed recipe. The final full CLI accepted 92
cases; architecture and spec policy passed. Both module vets exited 0; formatting
returned no Go files; whitespace check passed.

| Package (module-relative) | Final tests | Final race |
| --- | --- | --- |
| generator: cmd/generate-vectors | PASS 0.002s | PASS 1.010s |
| generator: internal/fixture | PASS 0.002s | PASS 1.013s |
| generator: internal/wire | PASS 0.014s | PASS 1.100s |
| conformance: cmd/check-architecture | No test files | No test files |
| conformance: cmd/check-spec | PASS 0.016s | PASS 1.091s |
| conformance: cmd/totipo-conformance | PASS 3.828s | PASS 10.575s |
| conformance: internal/cryptov1 | PASS 0.971s | PASS 2.359s |
| conformance: internal/graph | PASS 0.015s | PASS 1.041s |
| conformance: internal/object | PASS 0.003s | PASS 1.015s |
| conformance: internal/repopolicy | PASS 0.887s | PASS 1.940s |
| conformance: internal/schema | PASS 0.956s | PASS 4.979s |
| conformance: internal/specpolicy | PASS 0.780s | PASS 5.970s |
| conformance: internal/storage | PASS 0.753s | PASS 2.152s |
| conformance: internal/tlv | PASS 0.017s | PASS 1.041s |
| conformance: internal/totp | PASS 0.021s | PASS 1.025s |
| conformance: internal/vectors | PASS 1.636s | PASS 7.016s |
| conformance: internal/wireoracle | PASS 0.051s | PASS 1.327s |

Three generator and thirteen conformance packages passed tests and race, plus one
conformance command with no test files. All-package settings remained
`-count=1 ./...`, with `-race` for the two race commands.

| Fuzz target | Configured budget/workers | Result and package duration | Executions |
| --- | --- | --- | --- |
| FuzzDispatch | 10s / 2 | PASS 11.011s | 1,201,777 |
| FuzzOpen | 10s / 2 | PASS 11.012s | 1,085,326 |
| FuzzArrivalAndDisappearance | 10s / 2 | PASS 11.021s | 40,047 |

Every fuzz command used `-run '^$'`; no test selection, seed or failure semantics
changed. Starting local coverage was 2/2/1 inputs; new interesting inputs were
8/3/39 respectively. These remain local cache state. The current host's ordinary
writable Go cache differs from the earlier recipe-provided cache; this explains
different initial cached corpus counts, not a committed seed change. Bounded-run
completion exceeded nominal 10s slightly. There were no failures or reproducers.

Correction execution counts, separate from the original pass above:

- New baseline suites: **0** (qualified case B evidence reused).
- Focused Go test invocations: **1**, two packages; focused CLI runs: **4**.
- Focused fuzz: **0/0/0**. Python invocations/executions: **0/0**.
- Complete final host suites: **1**; incomplete/superseded attempts: **0**.
- Final module test commands: **2**; module vet commands: **2**;
  all-package race commands: **2**; bounded fuzz commands: **3**.
- Final architecture/spec/full CLI commands: **1 each**; formatting: **1**.
- Final CLI package tests include verify-only regressions; no extra final
  verify-only process was added after the successful focused CLI run.
- Generator executions: **0**. Agent Nix: **0**. Human Nix for corrected inputs:
  **1, owner-reported PASS**. The earlier one human PASS applies only to
  pre-correction inputs.
- Current Nix cached/fresh provenance: not supplied by owner. Remote CI: not run
  or queried.

No required check failed. All logs are local `/tmp/totipo-nix-first-*.log` files;
none entered Git source membership. No new repository testing infrastructure exists.

### Corrected source freeze, Git state and human gate

All included inputs were frozen before final host qualification. A SHA-256
inventory covers existing spec/vectors/requirements/generator/conformance inputs,
go.work, flake.nix, flake.lock and included archived history. Rechecking after
qualification passed. Deleted Make/wiring paths remain absent. No included-source
edit followed the final suite; only this excluded report was completed.
The workflow is identical to the qualified correction baseline. No dependency,
lock, cache, toolchain or serialized fixture update occurred.

The report retains its existing intent-to-add membership. No new file was added
for the correction; all three deletions and remaining edits are unstaged.
`git diff --cached --quiet` and `git diff --check` pass; nonignored untracked
inventory is empty. Exactly the intended fifteen paths differ from committed HEAD
(twelve changed/added paths and three deletions). Earlier simplification changes
are preserved except where this explicit correction supersedes them.

Final branch `main`; full HEAD unchanged:
`e2bb221cad1e1cb3432df3d7278e6eb7d3854c30`.
Final status:

```text
 M .github/workflows/conformance.yml
 M AGENTS.md
 M CONTRIBUTING.md
 D Makefile
 M README.md
 M conformance/README.md
 M conformance/internal/repopolicy/policy.go
 M conformance/internal/repopolicy/policy_test.go
 D conformance/internal/repopolicy/wiring.go
 D conformance/internal/repopolicy/wiring_test.go
 M flake.nix
 M requirements/README.md
 A review/SPEC_TOOLING_SIMPLIFICATION_REPORT.md
 M vectors/FORMAT.md
 M vectors/README.md
```

**Owner-reported PASS: `nix flake check`, supplied on 2026-10-14 for the corrected
Nix-first implementation.** This is one successful human execution; detailed
build output and cached/fresh provenance were not supplied. The frozen source
SHA-256 inventory was rechecked after the owner's report and passed. Only this
excluded report was updated to record the result. The agent has not evaluated or
executed Nix. All intended files are Git-known and no content is staged.
Any later included-source edit invalidates the affected result;
an excluded report-only status update does not. No separate Nix build or dependency
cache update is expected.

Supplemental macOS/Windows/Linux CI remains unverified pending a later
owner-authorized commit/push. Local Go/source inspection does not establish Nix
or GitHub provider execution. No host failure or known implementation issue
remains. No commit, push, tag, release, vector generation or GitHub settings change
occurred. All work remains uncommitted. Stop after this correction.
