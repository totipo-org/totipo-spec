# Windows architecture checker qualification

The architecture checker now skips filesystem symlink resolution for packages
identified by Go as standard library. Non-standard dependency checks remain in
place, including rejection of vendored sources resolving outside the vendor tree.
Directory-resolution errors now identify the dependency and path; repository-root
resolution errors also include the path.

## Starting revision and scope

Started on clean `main` at `82fd338cff84e3924b94e4486730eb2b60ff702a`, protocol
**v1/r19**. Classification: **CONFORMANCE_GO** (repository qualification logic and
tests). This report is supporting documentation. Historical reports are preserved.

The starting revision's W3A report records host qualification PASS and one
human-reported Nix PASS, with cache reuse unknown and remote CI not yet run at
that report's handoff. The current user supplied a subsequent Windows CI failure
at `go run ./conformance/cmd/check-architecture -root .`:
`FAIL architecture: The system cannot find the path specified.` That failure
is not upgraded to a pass by local Linux results.

## Diagnosis and change

The checker previously applied `filepath.EvalSymlinks` to every loaded package
directory before skipping standard-library packages. The likely Windows trigger
is a toolchain installation behind a junction that Go can use but this filesystem
operation cannot traverse. Go 1.23's
[symlink walker](https://github.com/golang/go/blob/go1.23.12/src/path/filepath/symlink.go)
can return a bare `ENOTDIR` while walking a non-directory path component. The
original error lacks a path, so the exact failing directory is not established;
Windows runtime confirmation remains a remote CI gate.

The Go tool already owns and classifies standard-library packages; directory
ownership policy applies to non-standard sources. The fix avoids this unnecessary
filesystem operation for standard packages and retains dependency load-error,
module, version, import, replacement, cache, and vendor policy checks. Tests model
an unresolvable standard-library directory, require contextual failure for the
same non-standard directory, accept a normal vendor directory, and reject a vendor
symlink escape. Existing source symlink and architecture mutation tests remain.
Windows symlink fixtures retain the existing privilege-unavailable skip policy.

## Qualification results and execution counts

Agent host: Linux amd64, Go 1.26.8. Logs are local evidence outside the checkout at
`/tmp/totipo-windows-architecture-evidence/`.

| Execution | Count and result |
| --- | --- |
| Baseline `make check` | 1, PASS |
| Baseline race or fuzz | 0 |
| Focused Python tests | 0 |
| Focused Go runtime tests | 1, `go -C conformance test -count=1 ./internal/repopolicy`, PASS |
| Windows Go compile check | 1, `GOOS=windows GOARCH=amd64 go -C conformance test -c -o /tmp/totipo-windows-architecture-evidence/repopolicy-windows.test.exe ./internal/repopolicy`, PASS; no Windows execution |
| Focused fuzz runs outside final suite | 0 for each target |
| Final complete host suites | 1, PASS; incomplete attempts 0 |
| Final `make check` | 1, PASS; 92 manifest cases, 92 case files, 15 Python tests OK |
| Final generator and conformance vet | 1 each, PASS |
| Final `make race` | 1, all packages PASS |
| Final `make fuzz` | 1, all three targets PASS |
| Agent Nix | 0 |
| Human Nix for this change | 1, PASS, human-reported |
| Cached Nix for this change | Cache reuse versus fresh execution not reported |
| Remote CI for this change | PENDING; no commit, push, or dispatch |

The focused runtime test preceded a Windows-safe diagnostic assertion correction
(`strconv.Quote` for backslash escaping); the corrected test passed in the final
suite. No test failure occurred.

Final ordinary and race package results: generator `cmd/generate-vectors`,
`internal/fixture`, and `internal/wire` PASS; conformance `internal/cryptov1`,
`internal/graph`, `internal/object`, `internal/repopolicy`, `internal/storage`,
`internal/tlv`, `internal/totp`, and `internal/vectors` PASS. Conformance command
packages compile with no test files. Architecture, spec structure, schema checks,
full manifest conformance, and verify-only all PASS. Go formatting and
`git diff --check` PASS.

Each fuzz command retained `-run '^$' -fuzztime 10s -parallel 2`:

| Target and package | Result and reported package duration |
| --- | --- |
| `FuzzDispatch`, `internal/object` | PASS, 10.071s |
| `FuzzOpen`, `internal/cryptov1` | PASS, 11.013s |
| `FuzzArrivalAndDisappearance`, `internal/graph` | PASS, 11.033s |

## Source invariance and handoff

Byte comparison against the starting revision, including path membership, found
only the two intended changed files under `conformance/internal/repopolicy/`.
Compared path counts: spec 1, vectors 97, requirements 2, conformance 27,
generator 8, tools 5; also Makefile, go.work, flake.nix, flake.lock, the CI workflow,
and consumed revision history, one each. All module and sum files are covered;
root go.mod/go.sum and generator/go.sum are absent. No vectors were regenerated.
Build/filter inputs and supplemental portability matrix are unchanged.

Inspection of flake.nix and structure-check consumers confirms the two Go files
are qualification-source inputs; this report is excluded and unconsumed. Prior
Nix evidence cannot qualify the changed Go inputs. All qualification inputs are
frozen for the human command `nix flake check`. This report is Git-known through
intent-to-add only; `git diff --cached --quiet` passed at handoff and after the
human gate.
Existing tracked changes remain unstaged. No commit or push is performed.

The user reported “nix flake check passed” for this change. Human execution count
is 1, PASS; agent Nix count remains 0. Cache reuse versus fresh derivation
execution was not reported. After that report, both changed source blobs matched
the handoff diff, every other file from the starting revision remained byte
identical, Git membership contained only the intended new report, and the cached
diff remained empty. Only this excluded report was updated after the gate;
qualification inputs remain unchanged, so no repeat suite or Nix gate is needed.

Remote Windows execution remains pending. Linux qualification and Windows
cross-compilation do not establish Windows runtime success.
