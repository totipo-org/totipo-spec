# Git-aware Nix worktree qualification policy

## Repository and starting state

Date: 2026-10-11. This milestone applies only to `/home/niki/Sources/totipo-spec`,
identified by its spec/vector/requirements/Go-conformance layout and origin
`git@github.com:totipo-org/totipo-spec.git`. No sibling repository was inspected.
Starting branch: `main`; starting HEAD:
`4cc5b74d55c5777268600cf83e6dde83d9d84981`.
The required initial branch/HEAD/status/index checks ran before inspection:
working tree clean, no untracked files, `git diff --cached --quiet` PASS.
Protocol remains **v1/r19**, manifest protocol `totipo-v1`, revision `r19`,
**92 cases**.

Task classifications: **BUILD_TOOLING**, **QUALIFICATION_POLICY**.
Repository ladder classifications: **BUILD_NIX_CI**, **DOCS_ONLY** (process text).
No semantic, verifier, generator, or conformance implementation change.

## Source-reference inventory and changes

Repository-wide search covered all non-ignored files, including hidden workflow
and environment files, for `path:`, `nix flake check`, `nix build`, `nix run`,
`nix develop`, `nix eval`, and `nix flake metadata`. The following is the complete
starting explicit Nix `path:` inventory; line numbers refer to starting HEAD.
Historical reports remain byte-identical and their old commands are historical
examples/evidence (classification C), rather than current instructions.

| Starting location | Prior form | Classification / disposition |
| --- | --- | --- |
| AGENTS.md:132 | `nix flake check path:.` | A CURRENT_GIT_WORKTREE; changed to `nix flake check` |
| CONTRIBUTING.md:16 | `nix flake check path:.` | A; changed to `nix flake check` |
| README.md:43 | `nix flake check path:.` | A; changed to `nix flake check` |
| .github/workflows/conformance.yml:23 | `nix flake check --print-build-logs path:.` | A; changed to `nix flake check --print-build-logs` |
| review/SPEC_GITHUB_ACTIONS_REFRESH_REPORT.md:55 | CI flake check with `path:.` | C EXAMPLE_OR_TEST_FIXTURE; historical execution description retained |
| review/SPEC_NIX_CI_UNIFICATION_REPORT.md:247 | CI flake check with `path:.` | C; historical command retained |
| review/SPEC_NIX_CI_UNIFICATION_REPORT.md:285 | Human flake check with `path:.` | C; historical documentation description retained |
| review/SPEC_NIX_CI_UNIFICATION_REPORT.md:331 | Human flake check with `path:.` | C; historical pending handoff retained |
| review/SPEC_QUALIFICATION_LADDER_REPORT.md:107 | Human flake check with `path:.` | C; historical policy description retained |
| review/SPEC_SUPPORT_TOOLING_AUDIT.md:69 | Human flake check with `path:.` | C; historical audit retained |
| review/SPEC_SUPPORT_TOOLING_WAVE1_REPORT.md:212 | Human flake check with `path:.` | C; historical handoff retained |
| review/SPEC_SUPPORT_TOOLING_WAVE2_REPORT.md:309 | Human flake check with `path:.` | C; historical handoff retained |
| review/SPEC_WAVE2_PORTABILITY_FIX_REPORT.md:263 | Human flake check with `path:.` | C; historical handoff retained |
| review/V1_R19_IMMUTABLE_VAULT_REPORT.md:444 | `nix build path:.` | C; historical explanation of inapplicable build retained |

There are no B INTENTIONAL_PATH_SOURCE or D UNKNOWN Nix occurrences.
The two other starting textual `path:` matches are the YAML field
`cache-dependency-path:` in the workflow and its historical description at
review/SPEC_GITHUB_ACTIONS_REFRESH_REPORT.md:115; neither is a flake reference.
New generic `path:` policy wording and old forms quoted in this report are
explanatory examples, not active commands.

No release workflow, release app, Nix wrapper, package.nix, or command-string
test exists here. No active build/run/eval/metadata command needs conversion.
The existing `.envrc` is already `use flake` followed by `layout go`; it remains
unchanged, as does the documented Git-aware `nix develop`.

## Authoritative guidance and source membership

AGENTS and CONTRIBUTING now specify Git-aware local pre-commit qualification,
unstaged modified tracked files, and authorized
`git add --intent-to-add <file>` for each new intended milestone file. README
links to this policy and names the canonical human `nix flake check` gate.
Intent-to-add establishes membership only: **no staged file content** is the
invariant, enforced by `git diff --cached --quiet`. Actual staging/commit still
requires explicit user authorization. Intent-to-add entries remain for human
review/commit and must not be reset after qualification.

Ignored/local build, cache, and secret state stays outside Git source membership.
No cache, generated output, secret, or ignored file was intent-added. The sole
new intended commit/source-membership path is:

- `review/NIX_GIT_WORKTREE_QUALIFICATION_POLICY_REPORT.md`

This report is intent-added before the human gate. Temporary validation tools,
checksums, and the frozen file-hash inventory reside outside the repository at
`/tmp/totipo-spec-nix-git-policy/`. No source-copy/filter mechanism was introduced.
Explicit raw/non-Git path references remain permitted when intentional.

## CI and substantive invariance

The qualification job performs checkout and Nix installation before its check.
There is no source-generation step and no required untracked repository input.
Only the source-reference suffix was removed; the command count, logging flag,
job order, matrix, action pins, permissions, credential policy, trust configuration,
and network behavior are unchanged. No application or release behavior changes.

Exact byte comparison against starting HEAD passed for **169** tracked
protected/build/history files: spec/, vectors/, requirements/, conformance/,
tools/, every prior review report, go.work, Makefile, flake.nix, flake.lock,
and .envrc. Protected source path membership, including additions/deletions,
also matches. Root go.mod, go.sum, go.work.sum, and package.nix were absent and
remain absent; conformance/go.mod and go.sum are unchanged.

Protocol/spec/schema/vector bytes, requirements pins, expected outcomes,
generator/verifier independence, Go/Python architecture, dependencies, locks,
vendorHash, and qualification semantics remain unchanged. No vectors were
regenerated and no language migration occurred.

The existing qualificationSource filter remains unchanged: spec, vectors,
requirements, conformance Go/module files, tools Python files, Makefile, go.work,
and review/V1_PRE_R16_REVISION_HISTORY.md. README, CONTRIBUTING, AGENTS, workflow
YAML, and this report are excluded from the derivation source. Inspection of
current structure checks/tests confirms they do not consume these changed files.
This report may therefore be finalized after human Nix without invalidating
qualification. All other files are frozen in the temporary SHA-256 inventory.

## Direct validation and execution accounting

PASS:

- actionlint v1.7.12, using the upstream release-checksum-verified binary outside
  the repository: `actionlint -shellcheck='' .github/workflows/conformance.yml`.
  ShellCheck is unavailable; workflow/expression/action-input checks ran. The
  changed shell command is a simple argument removal.
- Exact workflow comparison allowing only the intended command substitution;
  all five full 40-hex action pins and release annotations, read-only permission,
  both credential opt-outs, and absence of privileged triggers/write additions.
- `.envrc` shell syntax (`bash -n .envrc`).
- Markdown balanced fences and local link existence in changed guidance;
  canonical current-checkout gate forms; manual diff and source-consumer review.
- Protected byte/path invariance and `git diff --check`.

No existing static/pin/security test script was found; the focused temporary
audit enforces the existing workflow invariants. No permanent test or dependency
was added. No affected structural consumer requires a protocol/corpus test run.

| Execution | Count / status |
| --- | --- |
| Baseline make check | 0; infrastructure/process source-reference change |
| Focused Python / Go unit-test invocations | 0 / 0 |
| Focused fuzz Dispatch / Open / graph | 0 / 0 / 0 |
| Final complete host suites / incomplete attempts | 0 / 0 |
| Final make check / vet / race / fuzz | 0 / 0 / 0 / 0; substantive inputs unchanged |
| actionlint | 1 PASS |
| Temporary static/byte/document audit | 1 PASS |
| Agent Nix / generator executions | 0 / 0 |
| Human Nix | 1; PASS, human-reported |
| Remote CI dispatch / fresh observations | 0 / 0 |

No qualification failure or reproduction input occurred. Python test count and
Go package/race/fuzz results are not fresh executions here. The committed
SPEC_WAVE2_PORTABILITY_FIX_REPORT records inherited host qualification and one
human-reported Nix PASS for its frozen fix. The older Nix/CI-unification pending
gate is historical; no new pass or cached execution is inferred from it.

## Human gate and remote evidence

**PASS — human-reported.** After the frozen source-membership/index handoff,
the user reported “flake check passed” for the requested command from this
repository root:

```sh
nix flake check
```

Human execution count: 1; agent Nix execution count: 0. Cache reuse versus fresh
derivation execution was not reported, so neither is inferred. No remote job was run, pushed, or dispatched; remote
Linux qualification and supplemental portability remain unverified for these
uncommitted changes. All frozen file hashes and Git path membership were rechecked after human PASS
and remain identical to the handoff inventory. Only this excluded report was
finalized after Nix; no included effective input changed. The post-update frozen
inventory and `git diff --cached --quiet` checks also pass.

## Git/index handoff state and commit suggestion

Branch/HEAD remain `main` at the starting revision. Four tracked modifications:
AGENTS.md, CONTRIBUTING.md, README.md, .github/workflows/conformance.yml.
One new intent-to-add report, listed above. `git diff --name-only` contains these
five paths; `git diff --cached --name-only` is empty and
`git diff --cached --quiet` passes. No file contents are staged. All changes
remain uncommitted; no commit, tag, push, or intent-to-add reset occurred.

Suggested commit subject: **Adopt Git-aware Nix worktree qualification**.
