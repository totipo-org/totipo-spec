# Specification Nix/CI qualification unification

Status: implementation and non-Nix validation complete; human Nix gate pending.
No Nix command was run by the agent. Remote CI: **NOT RUN**.

## 1. Starting HEAD/state

- Branch: `main`.
- HEAD: `cdb4e91be1c6d3704874b2b92457ffe7be5e9084` (`flake update`).
- Initial `git status --short`: empty; clean committed checkout was the baseline.
- Current committed specification change: `a4e274c` (immutable VAULT, v1/r19).
- The checkout describes a design draft with moving pre-RC evidence, not a frozen
  release candidate. No released-state claim or older report baseline was substituted.
- Work was limited to this repository. Sibling CI conventions were read remotely.

## 2. Current normative revision and case count

Normative protocol: v1/r19. Manifest: `totipo-vector-manifest-v1`, r19.
Requirements: `requirements/v1-pre-rc.json`, `totipo-requirements-v1`,
`moving-pre-rc`, r19. Exactly **92 manifest entries and 92 physical JSON cases**;
all are required. Current profile artifact pins:

| Artifact | SHA-256 |
| --- | --- |
| Specification | `8bb76b890086eb4bf89271edd5861e02f833f3ffa11a83b5865080ed4e4228cb` |
| Manifest | `3953dbc315b4dcb3d0d31dd399bf82a15cb38c49fde2e2b93ea81c5cfe0ec714` |
| Manifest schema | `f265771f904be57d19602dd6da9a572c16b3f80fdd166c5177721aaa3f9a91ec` |
| Case schema | `99805d442f13872cd4febe9ac8ae4f36fd25607580bc1e457ff5ae2efe199e23` |

## 3. Baseline qualification

Before edits, all current credential-free workflow authorities passed:

| Authority | Baseline result |
| --- | --- |
| `make check`: specification structure and artifact hashes | PASS, v1/r19 |
| Python unit tests | PASS, 14 tests |
| Go `test -count=1 ./...` | PASS, seven test packages; two commands have no tests |
| Manifest-driven conformance | PASS, all 92 cases |
| Manifest/case schema validation | PASS, manifest and 92 cases |
| Manifest hashes, physical coverage, requirements pins | PASS, 92 case files |
| Formatting (`test -z "$(gofmt -l conformance)"`) | PASS |
| `go -C conformance vet ./...` | PASS |
| `make race` | PASS, all seven test packages with `-race -count=1` |
| `make fuzz` | PASS, all three targets at 10s each, parallelism 2 |
| `git diff --check` | PASS |

Baseline fuzz execution counts: FuzzDispatch 1,155,748; FuzzOpen 617,836;
FuzzArrivalAndDisappearance 225,817. Counts are observations, not thresholds.
No separate generator/check-only command exists in current CI. Generation was
not performed against committed evidence.

## 4. Makefile ownership and command inventory

Makefile remains unchanged and owns the semantics of `check`, `race`, and `fuzz`.
`check` serially requires `spec-check`, `test`, `conformance`, and `verify`.
`test` includes Python unit tests and uncached Go tests. `verify` runs schemas
and the Go manifest/profile verifier. Vet and formatting were explicit workflow
requirements and remain explicit alongside `make check` in the derivation.

Tracked-tree searches covered Make check/verify/conformance/race/fuzz, Go test/vet,
Python, generation/regeneration, Nix develop/flake check, and setup-go/setup-python.
Historical review reports were searched as evidence and preserved verbatim.

| Command/reference | Classification |
| --- | --- |
| `make check`, `make spec-check`, `make test`, `make verify`, `make conformance` | PORTABLE NORMAL CHECK |
| Current direct Python structure/schema/unittest commands | PORTABLE NORMAL CHECK |
| `gofmt -l conformance`, Go vet and ordinary Go tests | PORTABLE NORMAL CHECK |
| Authoritative full Linux flake check | LINUX-SPECIFIC CHECK |
| `make race`, Go `test -race -count=1 ./...` | RACE CHECK; current Linux scope |
| `make fuzz` and its three Go fuzz commands | BOUNDED FUZZ CHECK; current Linux scope |
| `go run ./conformance/cmd/generate-vectors -root .` | GENERATOR / MAINTENANCE ONLY |
| `tools/audit_r15.py`, `tools/audit_r16.py` | HISTORICAL EVIDENCE / explicit historical maintenance; not normal checks |
| Qualification commands in existing `review/*.md` | HISTORICAL EVIDENCE |
| `nix develop`, direnv, setup-go/setup-python | Environment setup for portable development/checks; not independent qualification |

## 5. Existing CI matrix/scope

Previously: Ubuntu/macOS/Windows × Go `1.23.x`/`stable` (six combinations).
All combinations ran structure, schemas, Python unit tests, formatting, vet,
Go tests, manifest execution, and manifest/profile verification. Both Ubuntu
combinations additionally ran race and the three bounded fuzz targets.
There was no temporary regenerate-and-compare step or generator invocation.

## 6. Existing flake baseline

The flake exposed a useful development shell and `nixpkgs-fmt` formatter, with
no checks, packages, or apps. Default systems are x86_64-linux, aarch64-linux,
x86_64-darwin, and aarch64-darwin. Shell tooling: Go, GNU Make, gopls,
golangci-lint, golangci-lint-langserver, and Python 3. Jailed Codex/Pi shells
also include libgcc/GCC and forward GOPATH/GOBIN. All behavior is preserved.

Root inputs and current locked revisions:

| Input | Revision |
| --- | --- |
| nixpkgs (`nixpkgs_3` root node) | `39ad350a0602fa0a58a544344e3e9187526ea45c` |
| flake-utils | `11707dc2f618dd54ca8739b309ec4fc024de578b` |
| llm-agents | `4bb57cff45b5554a02dba0cf7a8c7f4f010a864d` |
| jailed-agents (llm-agents follows root input) | `82b9d99d454bd752ae028e8a4bf0c90daf1c23ef` |

Existing Numtide cache/trust configuration is unchanged (see section 24).
Transitive inputs remain exactly as committed in `flake.lock`.

## 7. Final checks design

`checks.<linux-system>.spec = specQualification` exposes one full check on each
existing Linux system. It is a `pkgs.stdenv.mkDerivation`, using the C compiler
provided by stdenv, with Go, Python 3, and GNU Make as native build inputs.
Sequence: formatting assertion; `make check`; `go -C conformance vet ./...`;
`make race`; `make fuzz`. Failure of any command fails qualification.

The only installed artifact is `$out/qualified`, with a constant v1/r19 success
line. No timestamps, source-tree copy, release artifact, or application binary.
Darwin shells and formatter remain; Darwin checks do not claim full qualification.
Nix evaluation, build behavior, and aarch64 execution remain unconfirmed until
human/CI execution on the relevant platform.

## 8. No artificial default package

The product is the normative specification, exact corpus/pins, structural
verification, reference consumer, and conformance/integrity/race/fuzz evidence.
No `packages.default` or app was invented. The normal gate needs no second build.

## 9. Nix toolchain

The unchanged root nixpkgs pin selects Go **1.26.8** (`pkgs.go`) and Python
**3.14.7** (`pkgs.python3`/python314). These were confirmed by reading the pinned
nixpkgs source, not by Nix evaluation. Local non-Nix validation used these same
reported versions: `go version go1.26.8 linux/amd64`, `Python 3.14.7`.
`go.work` and the sole member `conformance/go.mod` declare Go **1.23.0**, with no
separate toolchain directive. Documented portable Python requirement is **3.9+**.
`GOTOOLCHAIN=local` prevents implicit compiler downloads in the check and the
fixed dependency derivation. No toolchain or flake input upgrade was made.

## 10. Source filtering

An explicit `lib.fileset.toSource` includes spec/, vectors/, requirements/,
conformance Go files plus go.mod/go.sum, Python sources in tools/, Makefile,
go.work, and **review/V1_PRE_R16_REVISION_HISTORY.md**. The last file is an actual
structural-check and Python-test dependency and must not be dropped.
There are currently no additional Go modules, embedded assets, or Go testdata
files requiring inclusion. All physical vector files are retained, allowing the
verifier to reject unexpected fixtures.

.git, .direnv, Python bytecode, editor files, unrelated design resources, generated
binaries, and other review reports are excluded. The new report is not a check
input. A matching isolated source copy passed the full suite without .git.
No verifier uses Git metadata. Existing historical audit tools are not invoked.

## 11. Writable caches

HOME, GOCACHE, GOTMPDIR, GOPATH, and GOMODCACHE are explicitly created below
TMPDIR. GOENV is disabled; Python bytecode writes are disabled. The unchanged
Makefile already exports `GOCACHE ?= .../.direnv/go-build` and therefore honors
the caller's temporary cache. No Makefile adaptation was necessary.
The vendored dependency tree is copied into the writable build tree.

## 12. Go module/dependency network boundary

There is exactly one workspace member, `./conformance`. External dependencies:
`golang.org/x/crypto v0.36.0` and indirect `golang.org/x/sys v0.31.0`.
Their versions and go.mod/go.sum bytes remain unchanged. `go mod verify` passed.

`goDependencies` uses the pinned nixpkgs `buildGoModule` fixed-output `goModules`
mechanism, overriding its dependency build command to `go work vendor` so that
root-level Makefile `go run` commands retain workspace semantics. Supply-chain
input `vendorHash` is:

`sha256-637XdgBok7hRhMw5s9hdY3Mr5jlPXqc1viSYoYh9Iqw=`

This hash was computed, not guessed: `GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local
 go work vendor -o /tmp/totipo-spec-vendor` used the existing verified module cache,
then a standalone SHA-256 implementation serialized the tree in canonical NAR
format (length-prefixed little-endian strings with eight-byte padding, sorted
entries, file bytes and executable flags). Nix confirmation of the dependency
output hash is pending the human gate. No Nix hashing command was run.

The fixed-output dependency fetch can obtain the exact pinned modules; its
entire output is hash-checked by Nix. This is the explicit supply-chain boundary,
not online dependency resolution during qualification. The real check sets
`GOFLAGS=-mod=vendor`, `GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local` and uses an
empty temporary module cache. There is no online download or fallback in its
build phase. The isolated full suite passed with this configuration. It was not
an OS-level network-isolation test; it demonstrated no Go proxy/checksum-service
or preexisting build/module cache dependency.

## 13. Python dependency boundary

Executed Python tools and tests use the standard library and repository-local
`check_vectors`; no third-party imports, pip, or package resolution. Pinned
nixpkgs Python is sufficient. Historical audit tools also use standard-library
imports but remain outside qualification.

## 14–16. Normal, vet, and race coverage

The derivation calls existing `make check` intact, preserving structure, artifact
hashes, Python tests, Go tests, every manifest case, schemas, physical coverage,
and requirements checks. It separately calls the existing explicit vet authority.
Race calls `make race` unchanged: `go -C conformance test -race -count=1 ./...`.
CGO is enabled, and compiler-bearing stdenv is used; no stdenvNoCC substitution,
package omission, failure suppression, or race-flag removal.

## 17. Bounded fuzz coverage

Unchanged Makefile policy: FuzzDispatch (internal/object), FuzzOpen
(internal/cryptov1), FuzzArrivalAndDisappearance (internal/graph), each with
`-run '^$' -fuzztime 10s -parallel 2`. Seeds and target implementations are
unchanged. Fresh GOCACHE keeps interesting-input caches transient. Go may write
failure reproductions under the writable unpacked source's testdata/fuzz paths;
that fails the derivation and never updates the checkout, vectors, or pins.
No failures or source corpus additions occurred in validation.

## 18. Manifest/case/hash/requirements integrity

The original Python schema checker and Go `vectors.Read`/`VerifyProfile` remain
unchanged. They validate all 92 cases, unique IDs/paths, exact case hashes,
manifest/case expected-outcome agreement, shape, safe paths, no unexpected
physical cases, exact required-case membership/hashes, revision, and the exact
specification/manifest/schema artifact pins. TestCorpus runs every case as a
subtest. The CLI also executes all 92 cases and independently verifies integrity.
No optional case selection, skip policy, or accepted regenerated hash was added.

## 19. Generation versus verification

No fixture generator is called by normal qualification or either CI job.
`go work vendor` constructs dependency inputs, not Totipo evidence. The explicit
Totipo generator stays maintenance-only, documented in vectors/FORMAT.md and
conformance/README.md. There was no existing regenerate-and-compare authority to
preserve. Generator implementation, fixtures, pins, and expected outcomes are
byte-identical. No temporary successful generation was substituted for verification.

## 20. Byte invariance

Before editing, SHA-256 inventories were captured for all tracked files under
spec/, vectors/, requirements/, conformance/, tools/, plus go.work and flake.lock:
**132 protected files**. Final comparison: **zero changes**. This includes schemas,
all expected outcomes, requirements pins, generator source, Go dependency graph,
workspace membership, and fixture inputs. No revision, revision-history, tag,
release metadata, or protocol wording changes. Historical reports are untouched.
The isolated qualification copy also retained its input bytes after execution.

## 21–22. Authoritative Linux CI and portability

Primary `qualification` job on ubuntu-24.04: checkout; install Nix; exactly one
`nix flake check --print-build-logs path:.`. No setup-go/setup-python, subsequent
Make/Go qualification, redundant Nix build, secrets, or write permissions.

Supplemental `portability` job retains macOS/Windows × Go 1.23.x/stable and adds
Linux Go 1.23.x. All retain the original portable structural/schema/unit,
formatting, vet, Go test, manifest execution, and integrity commands. Python 3.9
is now explicitly set up to exercise the documented minimum. Linux race/fuzz
move to the pinned authoritative check; they are not duplicated outside Nix.
The previous Linux stable portable coverage is supplied by the pinned Linux
qualification toolchain; mutable future stable is still exercised on macOS and
Windows. The supplemental jobs are explicitly named portability evidence.

## 23. Immutable GitHub Action pins

| Action | Immutable commit | Release annotation |
| --- | --- | --- |
| actions/checkout | `3d3c42e5aac5ba805825da76410c181273ba90b1` | v7.0.1 |
| cachix/install-nix-action | `13d8dd58da0234aa297dedd986986ccb8e7f3e24` | v31.11.1 |
| actions/setup-go | `40f1582b2485089dde7abd97c1529aa768e1baff` | v5.6.0, existing pin retained |
| actions/setup-python | `a26af69be951a213d495a4c3e4e4022e16d87065` | v5.6.0 |

Checkout and Nix-installer convention read from:
[desktop CI at 501d9faf](https://github.com/totipo-org/totipo-desktop/blob/501d9faf86cfe80a8b88ff2f57956a0037b7a858/.github/workflows/ci.yml)
and [Android CI at f8d5d450](https://github.com/totipo-org/totipo-android/blob/f8d5d4503a68b003c710c92bf587f99ffe14a56b/.github/workflows/ci.yml).
The setup-python v5.6.0 tag was resolved through the official Actions repository;
setup-go's existing annotation was also checked against its official tag.
Every checkout sets `persist-credentials: false`. Workflow permissions remain
`contents: read`.

## 24. Cache/trust/flake-config policy

Existing `extra-substituters = [ "https://cache.numtide.com" ]` and key
`niks3.numtide.com-1:DTx8wZduET09hRmMtKdQDxNNthLQETkc/yaX7M4qK0g=` are unchanged.
The reviewed installer receives `extra_nix_config: accept-flake-config = true`.
No new cache, signature disabling, credential, or trust bypass.

## 25–26. Documentation and lock invariance

README and CONTRIBUTING name `nix flake check path:.` as the human normal gate,
explain its Linux scope and full coverage, retain direct Make/Go development
commands, distinguish supplemental portability, and keep fixture maintenance
separate. CONTRIBUTING explicitly states that the agent does not run Nix.
No second human build command is required.

flake.lock remains byte-identical to current starting HEAD. SHA-256:
`44728fcfd8529462cb415b186ee4ce3400343dd911d325b6a965bd386b9f7749`.
No flake update or input change was performed.

## 27. Final non-Nix validation

Final ordinary checkout suite and isolated source-filter-equivalent suite both
passed: formatting, `make check`, explicit vet, `make race`, and `make fuzz`.
Both verified 92 cases, 14 Python tests, seven Go test packages (29 declared
Test functions, including TestCorpus's 92 case subtests), and two command packages
without tests. All three fuzz targets retained 10s/parallel-2 policy. The isolated
copy used fresh temporary caches, fixed workspace vendor inputs, and disabled
Go proxy/checksum service access. Exact final fuzz observations are recorded below.
| Final fuzz target | Checkout executions / elapsed | Isolated executions / elapsed |
| --- | --- | --- |
| FuzzDispatch | 805,505 / 10.092s | 523,830 / 10.210s |
| FuzzOpen | 439,845 / 11.010s | 8,010 / 11.014s |
| FuzzArrivalAndDisappearance | 163,581 / 10.217s | 34,859 / 11.043s |

All are PASS. Fuzz elapsed time includes startup/shutdown around the unchanged
10-second policy; execution counts vary and are not acceptance thresholds.
The isolated run started with only the source seeds (2, 2, and 1 respectively),
with no committed or local fuzz-cache dependency. All 130 protected files present
in that filtered copy remained byte-identical; no failure corpus was created.

`git diff --check` passed. These are non-Nix results, not a Nix sandbox claim.

Static review: Linux checks are reachable; tools come from the unchanged pin;
compiler-bearing stdenv and CGO support race; cache paths are writable and
temporary; full qualification has no network commands or .git dependency; fuzz
is bounded; no generator mutates evidence; no package was invented; lock unchanged.
CI has immutable actions, read-only permissions, no persisted credentials, explicit
flake-config acceptance, and one primary flake-check command with no repeated
qualification or secrets. Workflow execution remains pending remote CI.

## 28. Human one-command gate

**PENDING.** The agent did not run Nix. Human command, on Linux:

```sh
nix flake check path:.
```

No successful Nix evaluation, dependency hash acceptance, derivation execution,
or cached qualification result is claimed yet. The check definition reaches the
real suite, including make check, vet, race, and bounded fuzz. If the human gate
reuses a valid store result, record that rather than claiming fresh execution.
If it fails, diagnose and fix the infrastructure boundary without weakening tests,
then repeat this same gate.

## 29. Remote CI

**NOT RUN.** Nothing was pushed or dispatched. Exact committed CI must be
inspected separately after human review/commit/push.

## 30. Final Git state

Branch and HEAD unchanged. Four modified tracked infrastructure/documentation
files: flake.nix, .github/workflows/conformance.yml, README.md, CONTRIBUTING.md.
One untracked report: review/SPEC_NIX_CI_UNIFICATION_REPORT.md.
Index unchanged; everything remains unstaged/uncommitted. No commit, tag, release,
push, CI dispatch, or Nix execution by the agent. Completion remains pending the
human flake-check result.
