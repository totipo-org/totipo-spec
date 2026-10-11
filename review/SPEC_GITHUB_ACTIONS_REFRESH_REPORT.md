# Spec GitHub Actions refresh report

## Starting state and classification

Reviewed 2026-10-10. Started on clean `main` at
`b77f4f2a5f73c6edc781268e51bbfcfcbf1f40db`. Before edits,
`git branch --show-current`, `git rev-parse HEAD`, and `git status --short`
confirmed that state. Protocol: **v1/r19** (`totipo-v1`, manifest revision `r19`);
**92 manifest cases**. No revision bump or vector generation.

Classification: **BUILD_NIX_CI** only. No SPEC_TEXT, REQUIREMENTS, VECTOR_CORPUS,
CONFORMANCE_GO, PYTHON_VERIFIER, or GENERATOR change. No executable guard changed
or was added to the repository. Only two workflow pins and this report change.

Read AGENTS.md, CONTRIBUTING.md, README.md, every workflow YAML, flake.nix,
Makefile, conformance/go.mod and go.sum. Also inspected go.work, Dependabot
configuration, current Python structural consumers/tests, and historical
qualification reports. Historical reports are preserved.

## Complete workflow and action inventory

The sole workflow is `.github/workflows/conformance.yml`, triggered by push,
pull_request, and the pre-existing workflow_dispatch entry. There are five
external `uses:` occurrences, four distinct actions, no local actions, reusable
workflow calls, Docker actions, or standalone cache action.

| Workflow/job | Action | Type | Old runtime | Final runtime |
| --- | --- | --- | --- | --- |
| conformance.yml / qualification | actions/checkout | JavaScript | Node 24 | Node 24 |
| conformance.yml / qualification | cachix/install-nix-action | Composite | Bash | Bash |
| conformance.yml / portability | actions/checkout | JavaScript | Node 24 | Node 24 |
| conformance.yml / portability | actions/setup-python | JavaScript | Node 20 | Node 24 |
| conformance.yml / portability | actions/setup-go | JavaScript | Node 20 | Node 24 |

| Action | Starting release / full SHA | Final release / full SHA |
| --- | --- | --- |
| actions/checkout (both jobs) | v7.0.1 / `3d3c42e5aac5ba805825da76410c181273ba90b1` | v7.0.1 / `3d3c42e5aac5ba805825da76410c181273ba90b1` |
| cachix/install-nix-action | v31.11.1 / `13d8dd58da0234aa297dedd986986ccb8e7f3e24` | v31.11.1 / `13d8dd58da0234aa297dedd986986ccb8e7f3e24` |
| actions/setup-python | v5.6.0 / `a26af69be951a213d495a4c3e4e4022e16d87065` | v7.0.0 / `5fda3b95a4ea91299a34e894583c3862153e4b97` |
| actions/setup-go | v5.6.0 / `40f1582b2485089dde7abd97c1529aa768e1baff` | v7.0.0 / `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` |

All four canonical upstream repositories are active, neither archived nor
disabled, and their latest stable GitHub releases are the final versions above.
Read canonical `/releases/latest`, release histories, tag commit resolutions,
README files and action.yml metadata. Each final tag resolved through GitHub's
commit endpoint to the full SHA shown; metadata was inspected at those releases.
Already-current checkout and installer pins remain unchanged. All uses retain
`owner/action@<40 lowercase hex> # vX.Y.Z`; no floating version was introduced.

## CI authority and portability map

`qualification` remains the authoritative Ubuntu 24.04 Linux Nix job, with:

```sh
nix flake check --print-build-logs path:.
```

Unchanged flake derivation coverage: Go formatting, make check (structure,
Python tests, Go tests, manifest conformance and integrity), go vet, all-package
race, and the three bounded fuzz targets. FuzzDispatch, FuzzOpen and
FuzzArrivalAndDisappearance retain `-run '^$'`, 10s and parallelism 2. There is
no packages.default and no Linux setup-go substitute for Nix qualification.

`portability` remains supplemental: macos-latest and windows-latest each run
Go 1.23.x and stable, plus ubuntu-24.04 with Go 1.23.x (five matrix jobs).
Python remains 3.9. Fail-fast remains false. All structure, formatting, vet,
test, conformance and integrity commands and working directories are unchanged.
Formatting retains explicit Bash; other steps retain runner-default shells,
including PowerShell on Windows. No release authority is assigned to these jobs.

## Upstream release and behavior review

- [checkout v7.0.1 release](https://github.com/actions/checkout/releases/tag/v7.0.1):
  unsafe-PR default-check optimization, ASCII branch whitespace handling,
  escaping unset values, and dependency updates. Already current.
  [Pinned metadata](https://github.com/actions/checkout/blob/3d3c42e5aac5ba805825da76410c181273ba90b1/action.yml)
  uses Node 24. Both jobs retain `persist-credentials: false`, default event ref,
  fetch-depth 1, no tags, submodules or LFS. Current checks read working-tree
  files, including archived history, rather than Git ancestry; there is no CI
  diff or generation step requiring deeper history. Checkout v6 moved persisted
  credentials under RUNNER_TEMP; opting out remains appropriate. The v7 unsafe
  fork checkout protections remain enabled; no privileged PR trigger is added.
- [setup-go v7.0.0 release](https://github.com/actions/setup-go/releases/tag/v7.0.0):
  ESM migration and toolkit dependencies, including @actions/cache 6.2.0.
  [v6.0.0 release](https://github.com/actions/setup-go/releases/tag/v6.0.0)
  introduced Node 24 and toolchain handling. Reviewed intervening v6.1–v6.5
  notes: download fallback/custom URL support, default cache-file switch to
  go.mod, and toolkit/security updates.
  [Pinned metadata](https://github.com/actions/setup-go/blob/b7ad1dad31e06c5925ef5d2fc7ad053ef454303e/action.yml)
  confirms Node 24. Node 24 requires runner 2.327.1 or newer; current GitHub-hosted
  runners are the intended environment, with actual platform execution pending CI.
- [setup-python v7.0.0 release](https://github.com/actions/setup-python/releases/tag/v7.0.0):
  ESM migration, toolkit/cache 6.2.0, manifest-fetch validation/retry, warning
  classification fixes, removal of pip-install, and fixture/dependency updates.
  [v6.0.0 release](https://github.com/actions/setup-python/releases/tag/v6.0.0)
  introduced Node 24; intervening v6.1–v6.3 notes cover optional pip-install,
  additional interpreter platforms, Windows pip-cache handling, Linux cache-key
  scoping and dependency fixes. The removed input was never used here; removal
  of EOL versions in upstream test fixtures does not change this workflow's
  explicit Python 3.9 selection. No pip-version or package-install input is set.
  [Pinned metadata](https://github.com/actions/setup-python/blob/5fda3b95a4ea91299a34e894583c3862153e4b97/action.yml)
  confirms Node 24; Python dependency caching remains disabled by default.
- [install-nix-action v31.11.1 release](https://github.com/cachix/install-nix-action/releases/tag/v31.11.1):
  Nix 2.35.1 to 2.35.2 fixes an assertion crash that could abort builds. Already
  current; neither the action nor installed Nix selection changes here.

No repository-controlled JavaScript action remains on Node 20. No runtime-warning
suppression or Node override was introduced.

## Go and cache boundary review

Compared setup-go v5.6.0 and v7.0.0 main, cache-restore, cache-save and cache-utils
sources; inspected final installer and package-manager definitions. Caching remains
on by default. The workflow explicitly retains
`cache-dependency-path: conformance/go.sum`, so the new default go.mod lookup is
not used. Cache key construction still includes runner OS, architecture, Linux
image where applicable, resolved Go version and dependency-file hash. Cached
paths remain `go env GOMODCACHE` and `go env GOCACHE`. Cache-save remains a
successful-job post step; denied/colliding saves do not become qualification.
No standalone cache action or additional cache privilege is needed.

New setup-go exports `GOTOOLCHAIN=local`, making commands use the selected Go
installation. This is compatible with the intended explicit matrix: go.work and
conformance/go.mod both require Go 1.23.0, neither has a toolchain directive, and
the dependencies require no higher Go minimum (x/crypto v0.36.0, x/sys v0.31.0).
The explicit `go-version` takes precedence over file inference. Thus there is no
change to the effective intended Go versions or dependency graph. No automatic
module update, tidy, dependency installation policy, or custom download URL is
introduced. Supplemental jobs retain normal checksum-verified host module use;
Linux's fixed vendorHash and offline vendored qualification remain authoritative.

## Third-party action and permissions review

The only non-GitHub action is the maintained canonical
[cachix/install-nix-action](https://github.com/cachix/install-nix-action).
Its [pinned composite metadata](https://github.com/cachix/install-nix-action/blob/13d8dd58da0234aa297dedd986986ccb8e7f3e24/action.yml)
contains one Bash step invoking the bundled
[install-nix.sh](https://github.com/cachix/install-nix-action/blob/13d8dd58da0234aa297dedd986986ccb8e7f3e24/install-nix.sh),
with no nested external `uses:` or JavaScript runtime. Read the script without
executing it. It downloads the versioned official Nix 2.35.2 installer, invokes
it with sh, configures multi-user installation on systemd Linux, enables flakes,
and sets environment paths. Existing sudo/KVM/trusted-user behavior and default
GitHub token access for Nix GitHub fetches remain unchanged. The action SHA pins
the action code; its installer download is an existing network trust boundary,
not a newly claimed independent immutable installer digest.

Workflow permissions remain `contents: read`; no write, OIDC, secrets, privileged
trigger, replacement project, extra cache service or permission scope is added.
`accept-flake-config = true` is unchanged, as are flake substituters/trust keys.

## Protected-byte evidence and invalidation

Before edits, captured SHA-256 values for 167 tracked protected/build/history
and process files. Final `sha256sum -c` passed for all 167. Separately, a temporary
standard-library Go validator compared exact bytes against starting HEAD and
path membership for **135** files in spec/, vectors/, requirements/, conformance/
(including generator and tests), tools/ source, go.work, Makefile, flake.nix,
flake.lock and the included archived history. No added/deleted source paths.
The existing ignored tools/__pycache__ directory is excluded from source
membership; it predates this milestone and is not a qualification input.
Root go.mod, go.sum and go.work.sum are absent and remain absent; the only module
files are conformance/go.mod and conformance/go.sum, both identical. All previous
reports, protocol bytes, cases, manifest, schemas, outcomes, pins, generator
outputs, dependency graph, vendorHash and lock are unchanged.

Reinspected qualificationSource: spec/, vectors/, requirements/, conformance
Go/module files, tools Python files, Makefile, go.work, and
review/V1_PRE_R16_REVISION_HISTORY.md are included. Workflows and this new report
are excluded. Current structure checks/tests consume no workflow YAML or this
report. No executable qualification input or Nix evaluator input changes.

## Focused validation and execution counts

PASS: actionlint **v1.7.12**, downloaded outside the repository and verified
against its upstream release checksum, with
`actionlint -shellcheck='' .github/workflows/conformance.yml`. ShellCheck is
not installed; embedded shell static analysis was disabled, while workflow,
expression and action-input linting ran. Shell commands themselves are unchanged.
Two actionlint executions passed.

PASS: temporary Go static audit, enforcing the reviewed allowlist, exact 40-hex
SHA/release comments, all five uses across every workflow, expected downloaded
metadata runtimes, starting HEAD, protected bytes/path membership, absent root
module files, and an exact baseline workflow comparison allowing only the two
intended pin substitutions. No existing workflow pin guard was found, and no
permanent validator or new Python surface was necessary.

The first Go audit attempt failed on pre-existing ignored tools/__pycache__
path membership. After inspecting its older timestamps and ignored status, the
temporary audit was corrected to exclude that directory; the second attempt
passed. This was an audit-scope error, not a source or qualification failure.
Temporary audit/upstream artifacts are outside the repository under
`/tmp/totipo-actions-review/`; they are session evidence, not committed tooling.

PASS: final `git diff --check` and new report whitespace check; manual diff review
confirms only two action pins change in executable workflow content.

| Execution | Count / result |
| --- | --- |
| Baseline make check | 0; excluded workflow-only milestone |
| Focused Python unit tests | 0 |
| Focused Go tests | 0 |
| Temporary Go static audits | 2 attempts: 1 scope failure, 1 pass |
| actionlint | 2 passes |
| Focused fuzz by target | 0 each |
| Final complete host suite / incomplete attempts | 0 / 0 |
| Final make check / vet / race / fuzz | 0 / 0 / 0 / 0; not invalidated |
| Python/Go package results this milestone | Not executed; no fresh test pass claimed |
| Human Nix / agent Nix | 0 / 0 |
| Generator / remote CI dispatch | 0 / 0 |

## Nix and remote CI disposition

No fresh human Nix request: qualification source, evaluator files, lock and fixed
vendorHash are identical. Workflow pin maintenance does not invalidate existing
Nix derivation evidence. The inherited SPEC_NIX_CI_UNIFICATION_REPORT records
passing host/isolated suites (92 cases, 14 Python tests, seven tested Go packages,
two command packages without tests, vet/race and all three 10s fuzz targets),
but explicitly leaves its human Nix gate **pending** and remote CI **not run**.
SPEC_QUALIFICATION_LADDER_REPORT carries that status forward. Those are inherited
observations, not fresh executions, cached Nix results or a Nix pass here.
The older r19 report's user-reported Nix pass predates the later Nix/CI
unification and does not resolve that later pending gate.

No workflow was dispatched, pushed or otherwise run remotely. Updated action
execution, including macOS/Windows Python availability and shell integration,
remains unverified until normal CI after user review/commit/push. Remote CI is
supplemental portability evidence alongside the authoritative Linux Nix job;
static validation does not claim remote execution success.

## Final Git state

`main` remains at `b77f4f2a5f73c6edc781268e51bbfcfcbf1f40db`.
One tracked modification: `.github/workflows/conformance.yml`.
One untracked addition: `review/SPEC_GITHUB_ACTIONS_REFRESH_REPORT.md`.
The index is unchanged and empty of staged changes. Everything is unstaged and
uncommitted. No commit, push, tag, release, workflow dispatch, vector regeneration,
Nix execution, dependency update, or file change outside this repository occurred
except temporary review/validation artifacts under /tmp.
