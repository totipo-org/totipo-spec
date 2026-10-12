# Vendor path portability fix

Started on clean `main` at `ed86f433aeafcd234e1531e012ebdd0ee466f55a`,
protocol **v1/r19**. Classification: **CONFORMANCE_GO** (repository
qualification logic and tests). Historical reports remain unchanged.

## Diagnosis and implementation

The user reported macOS and Windows CI failures in `TestDependencyDirVendor`:
`vendor source symlink escape: golang.org/x/crypto/hkdf`. The check compared a
symlink-resolved package directory with an unresolved repository root. An ancestor
alias, such as macOS `/var` resolving to `/private/var`, therefore caused false
escape detection. The exact Windows runner path was not supplied; platform
runtime confirmation remains pending remote CI.

The helper now resolves the repository root before comparing resolved package
containment. It deliberately does not resolve the vendor boundary itself, so
links carrying vendored sources outside that boundary remain rejected.
Standard-library handling and all other dependency policy checks are unchanged.

A new regression creates a repository directory alias. It reproduced the exact
reported failure on Linux before the implementation change. With the fix it
accepts ordinary vendored sources through that alias and rejects a vendor link
to an external directory. The existing ordinary vendor and escape checks remain.
Windows symlink creation retains the existing privilege-unavailable skip policy.

## Inherited evidence

The starting revision's Wave 3B report records host qualification PASS and one
human-reported Nix PASS, with cache reuse unknown and remote CI not yet run at
handoff. The user's subsequent macOS/Windows failure is failure evidence, not a
remote pass. Those prior qualifications do not cover the changed Go inputs.

## Execution evidence

Agent host: Linux amd64, Go 1.26.8. Local logs remain outside Git under
`/tmp/totipo-vendor-*.log`, including the deliberate pre-fix reproduction failure.
No vectors were generated; no Nix command was run by the agent.

| Execution | Count and result |
| --- | --- |
| Baseline `make check` | 1, PASS |
| Baseline race/fuzz | 0 |
| Focused Python test invocations | 0 |
| Focused Go test invocations | 2: vendor regression before fix FAIL; complete repopolicy package after fix PASS |
| Focused fuzz outside final suite | 0 per target |
| Final complete host suites | 1, PASS; incomplete attempts 0 |
| Final `make check` | 1, PASS: 92 manifest cases, 92 case files, 3 Python tests OK |
| Final conformance and generator vet | 1 each, PASS |
| Final `make race` | 1, all packages PASS |
| Final `make fuzz` | 1, all three targets PASS |
| Formatting and `git diff --check` | PASS |
| Agent Nix | 0 |
| Human Nix for this milestone | 1, PASS, human-reported |
| Cached Nix for this milestone | Cache reuse versus fresh execution not reported |
| Remote CI for this milestone | PENDING; no commit/push/dispatch |

Ordinary and race tests PASS for generator `cmd/generate-vectors`,
`internal/fixture`, `internal/wire`; conformance `cmd/check-spec`,
`cmd/totipo-conformance`, `internal/cryptov1`, `internal/graph`, `internal/object`,
`internal/repopolicy`, `internal/schema`, `internal/specpolicy`, `internal/storage`,
`internal/tlv`, `internal/totp`, `internal/vectors`. `cmd/check-architecture`
compiles with no test files. Architecture, spec structure, full conformance, and
verify-only checks PASS.

All targets retained `-run '^$' -fuzztime 10s -parallel 2`, including seeds
and failure semantics:

| Target and package | Result and reported package duration |
| --- | --- |
| `FuzzDispatch`, `internal/object` | PASS, 10.075s |
| `FuzzOpen`, `internal/cryptov1` | PASS, 10.940s |
| `FuzzArrivalAndDisappearance`, `internal/graph` | PASS, 11.022s |

## Byte invariance and human handoff

Compared every tracked byte against the starting revision for spec (1 path),
vectors (97), requirements (2), conformance (35), generator (8), tools (2),
Makefile, go.work, flake.nix, flake.lock, the CI workflow, and consumed revision
history (one each). Only the two intended files under
`conformance/internal/repopolicy/` differ. No original paths were added/deleted;
the only new intended path is this report. All Go module/sum files are covered;
root go.mod/go.sum and generator/go.sum are absent. Protocol, corpus,
requirements, generator, Python verification, build/filter inputs, and CI matrix
remain byte-identical.

Inspected the actual flake fileset and structure consumers: changed Go files are
qualification-source inputs; this report is excluded and unconsumed. Prior Nix
results cannot qualify this change. All qualification inputs are frozen for the
human command `nix flake check`. The report is intent-added only, existing edits
remain unstaged, and `git diff --cached --quiet` passes. No commit or push.

Per AGENTS.md, “The agent does not run Nix. The final routine pinned-environment
gate is the human command: nix flake check.” The user subsequently reported
“nix flake check passed” for this milestone. Human execution count is 1, PASS;
agent Nix count remains 0. Cache reuse versus fresh execution was not reported.

Post-gate inspection confirmed the same main/starting HEAD, the exact two tested
Go diffs, unchanged remaining tracked inputs, no unrelated untracked files, and
this report's intent-to-add membership. The cached diff remains empty. Only this
excluded report was updated after the gate; qualification inputs are unchanged,
so no repeat host suite or Nix gate is needed. Remote macOS/Windows CI remains
pending; Linux execution does not establish those platforms' runtime success.
