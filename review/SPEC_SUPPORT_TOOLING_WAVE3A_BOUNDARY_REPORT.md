# Totipo support tooling Wave 3A implementation boundary

## 1. Starting state

Date: 2026-10-11. Work was restricted to `totipo-spec`; no sibling repository
was inspected. Required Git checks preceded inspection: branch `main`, HEAD
`a20aa62543219241b08237ccdbba4b6e72706537`, empty status, cached diff exit 0.
The architecture audit was already committed at that HEAD. Historical reports
remain unchanged. Protocol remains **v1/r19**, moving pre-RC.

**W3A complete: host qualification PASS, human-reported Nix PASS, and post-Nix
frozen-input verification PASS.** Changes remain uncommitted. No W3B/W3C work or Python migration is
claimed. The inherited policy report records a human-reported Nix PASS for its
prior inputs, without cache/fresh-build attribution; that does not qualify W3A.
Remote portability evidence is not inferred.

## 2. Classification

Requested: SUPPORT_TOOLING_ARCHITECTURE, LANGUAGE_SIMPLIFICATION,
GENERATOR_VERIFIER_INDEPENDENCE, WAVE_3A. Ladder: GENERATOR, CONFORMANCE_GO
(new architecture-policy implementation/tests only), BUILD_NIX_CI, DOCS_ONLY
(contributor guidance). No SPEC_TEXT, REQUIREMENTS, VECTOR_CORPUS, or
PYTHON_VERIFIER changes.

## 3. Baseline artifact inventory

There are **92 manifest cases, 92 physical case JSON files, and 92 moving-profile
required cases**, both before and after. The protected inventory comprises
100 files: spec 1, vectors 97, requirements 2. Exact baseline paths and SHA-256
values, the original complete trees (including schemas, all manifest IDs/paths/
hashes/expectations and profile pins), original Makefile, and relevant target
expansions were captured outside the checkout at `/tmp/totipo-w3a-evidence/`
before editing or generation. `protected-baseline.json` is the comparison input.

Unchanged artifact SHA-256 pins:

| Artifact | SHA-256 |
| --- | --- |
| spec/totipo-vault-format-v1.md | `8bb76b890086eb4bf89271edd5861e02f833f3ffa11a83b5865080ed4e4228cb` |
| vectors/manifest.json | `3953dbc315b4dcb3d0d31dd399bf82a15cb38c49fde2e2b93ea81c5cfe0ec714` |
| vectors/manifest.schema.json | `f265771f904be57d19602dd6da9a572c16b3f80fdd166c5177721aaa3f9a91ec` |
| vectors/case.schema.json | `99805d442f13872cd4febe9ac8ae4f36fd25607580bc1e457ff5ae2efe199e23` |

## 4. Baseline qualification

Agent baseline: one `make check`, one conformance vet, one all-package
`go -C conformance test -race -count=1 ./...`, one `make fuzz`; all PASS.
Explicit W3A baseline instructions requested broader evidence than the ladder's
ordinary default. No baseline generation or Nix execution occurred.

Python: **15 tests, OK**, structure PASS, schema validation of 92 cases PASS.
Go: cryptov1, graph, object, storage, tlv, totp, vectors all PASS; the two command
packages compiled with no tests. Both full and verify-only CLIs passed 92 cases.
Vet PASS. Race PASS for all packages. Exact fuzz policy was preserved:

| Target/package | Requested budget/parallelism | Baseline result/runtime | Final result/runtime |
| --- | --- | --- | --- |
| FuzzDispatch / internal/object | 10s / 2 | PASS / 10.096s | PASS / 10.109s |
| FuzzOpen / internal/cryptov1 | 10s / 2 | PASS / 11.010s | PASS / 10.084s |
| FuzzArrivalAndDisappearance / internal/graph | 10s / 2 | PASS / 10.544s | PASS / 10.658s |

Each used `-run '^$' -fuzztime 10s -parallel 2`. Runtime can exceed its fuzz budget
while workers finish. Seeds, persistent corpus, and failure semantics were not
changed; Go's ordinary ignored fuzz cache accumulated discoveries.

## 5. Prior dependency graph

The old conformance generator imported vectors (models, Decode, Hash), object
(Object.Encode), tlv (framing/integer/mutation helpers), cryptov1
(derivation/Seal/padding/ID/AAD/VAULT_ID), and graph (Node/Result models).
Its compiled repository closure therefore included **all seven verifier semantic
packages**, including incidental storage/totp through vectors. The command was
508 lines. Existing authored graph/fold/storage expectations were not evaluator
results; their declarations are preserved.

## 6. Final module/workspace structure

```text
go.work                         # exactly ./generator and ./conformance
generator/go.mod                # totipo/generator, Go 1.23.0, no requires
  cmd/generate-vectors/          # deliberate maintenance command and tests
  internal/fixture/             # writer-owned artifact JSON and source reading
  internal/wire/                # independent writer/crypto composition and tests
conformance/go.mod              # totipo/conformance, unchanged dependencies
  cmd/totipo-conformance/        # unchanged
  cmd/check-architecture/        # new stdlib repository-policy CLI
  internal/repopolicy/           # new guard and disposable-tree negative tests
  internal/{tlv,object,cryptov1,graph,storage,totp,vectors}/ # unchanged
```

No root module, third module, common/shared implementation directory, workspace
replace, or module replace. Generator go.sum, root go.mod/go.sum, and go.work.sum
are absent. Total Go source: 31 files / 4,516 lines (generator 7 / 1,197;
conformance 24 / 3,319), versus audit baseline 22 / 3,379. Growth buys separate
composition/model ownership and guard/regression evidence.

## 7. Generator-local data model

`generator/internal/fixture/records.go` owns Case, every nested artifact record,
Manifest/Entry/Profile/Pin, graph Node/Result/Step, fold, storage/workflow,
bootstrap, post-AEAD, and retained TOTP record shapes. `wire.Token` owns authoring
fields including optional pointers. These JSON declarations deliberately retain
the serialized contract and field order of the old structs, but share no Go
identity or executable model behavior. This is intentional schema-shape
ancestry, not a claim that the serialized format changed. No verifier structs
were exported or extracted into a shared module.

## 8. Generator writer

`wire/writer.go` implements its own integer byte loop and TLV append construction;
TOKEN validation occurs on authoring fields before constructing an ordered field
sequence. It independently sums encoded lengths and allocates the writer output.
It does not validate its output by calling a TOKEN parser. Its bounded SplitField
utility only splices already authored TLVs for negative fixtures. No accepting
TOKEN reader, graph/fold/storage evaluator, TOTP calculator, or bootstrap KDF
was added.

## 9. Independent crypto composition

Root extraction is stdlib HMAC with a 32-zero-byte salt. The only required HKDF
expand length is 32, implemented as HMAC(PRK, domain || suffix || 0x01).
Generator-owned literals select K_id, object-root key, per-object key and AAD.
OBJECT_ID commits to the semantic bytes, nonce is its first 12 bytes, padding
is a separately composed u16 length/semantic/zero concatenation. A fixed-width
Encrypt operation uses stdlib AES/GCM and appends its tag. Malform changes one
post-AEAD predicate before sealing under the filename-derived context.
VAULT_ID validates the retained representation's length/magic and hashes its
exact bytes. No custom crypto primitive or third-party dependency was introduced.

## 10. Constant ownership

Generator owns its tag literals/order, widths, parent bound, enum/range checks,
UTF-8 byte limits, capacity, plaintext/envelope sizes, HKDF domains, AAD label,
nonce slicing and VAULT magic/version/length. Verifier constants and Python
literals remain independently owned and byte-unchanged. Neither side imports a
common constants source, generated template, or shared semantic package.

## 11. Expectation provenance

Computed: canonical semantic bytes, keys, IDs, padding, encrypted bytes,
VAULT_ID, case/artifact hashes and pins. Authored: graph heads/values/unresolved
IDs, storage/workflow classifications, negative classifications, symbolic fold
plans/stage IDs. Head records attach facts to explicitly authored head IDs;
no evaluator chooses those IDs. Retained: RFC TOTP files are read and kept
byte-identical without recalculation; four bootstrap source records preserve
password/root/salt/nonce/wrap-key/header/record values. Existing derived VAULT_ID
and the authored one-bit mutation continue to be computed. No stronger external
provenance for bootstrap answers is claimed.

## 12. Architecture guard design

`conformance/internal/repopolicy/policy.go` and `cmd/check-architecture` use only
stdlib. Source walks parse complete Go files/imports/comments, including tests,
generated-marked source and inactive build tags. Every repository Go implementation
file must be under one approved module root. Third/root go.mod files fail.
Only root .git/.direnv/__pycache__/workspace-vendor directories are excluded as
repository metadata, ignored local state, or dependency material; nested same-name
directories under module roots remain inspected. Loaded vendor dependencies are
checked separately. No .git is needed, so filtered Nix sources work.

The guard rejects cross-role/third Totipo imports, unreviewed external imports,
source symlinks (including directory links), go:linkname/go:generate directives,
plugin/C/unsafe imports, native assembly/C/syso source bridges, and os/exec outside
the exact guard implementation file. There is no current CLI subprocess test
exception. The guard invokes only Go inspection commands, never either role's CLI.
This is an architecture policy, not a sandbox or proof against a malicious author
rewriting policy or disguising copied algorithms.

## 13. Guard negative tests

Disposable stdlib-only trees exercise generator→verifier and verifier→generator
imports, test-only/tagged/generated imports, nested vendor source, third package,
third/root module, extra workspace member, local replace/require bridges,
workspace replace, wrong Go version, unreviewed external dependencies on both
sides, source outside approved roots, linkname after imports, go:generate,
os/exec/plugin/cgo/assembly, and source file/directory symlinks. Windows can skip
symlink fixture creation when privilege is unavailable; other tests still run.
Direct dependency-closure tests additionally exercise sibling imports through
tests, unavailable packages, and a nonstdlib generator dependency. A pristine
two-module temporary tree must pass the full guard. Temporary-tree tests select
readonly module mode because they have no workspace vendor tree.

## 14. Module metadata policy

Go's own `mod edit -json` and `work edit -json` inspect directives without writing
metadata. Exact module names and Go 1.23.0 are required; generator requires nothing,
conformance requires exactly x/crypto v0.36.0 and indirect x/sys v0.31.0.
Replace/exclude/retract/tool/toolchain/godebug/ignore directives are rejected.
Unknown metadata JSON fields fail closed. Workspace uses are exactly the two
relative roots, no aliases/duplicates/replacements. Source inventory also rejects
an otherwise unused third module. No dependency upgrades occurred.

## 15. Dependency closures

Both actual package/test closures are inspected with `go list -deps -test -json
./...`; command errors, incomplete packages and dependency errors fail. Ownership,
module version/replacement, module-cache or workspace-vendor directory membership,
and vendor directory symlink resolution are checked. Inspection is offline with
local toolchain/GOENV=off, preserving only empty GOFLAGS or reviewed vendor/readonly
mode; overlays/modfiles/custom execution flags fail review.

Generator closure is **stdlib-only**. Verifier external closure is exactly
`golang.org/x/crypto/{argon2,blake2b,hkdf}` and `golang.org/x/sys/cpu`, at unchanged
versions. These are reviewed primitive implementation packages, not Totipo
semantics. Common nonstdlib external packages between roles: **0**. Standard
crypto/aes, cipher, hmac, sha256 and ordinary stdlib utilities may overlap.
Final full JSON closure captures are outside the checkout.

## 16. go.work and standalone behavior

The workspace is convenience, not isolation. Sibling internal ancestry and the
source/metadata/closure guard enforce ownership. GOWORK=off, GOPROXY=off,
GOSUMDB=off `list -mod=readonly -deps -test ./...` and `test -mod=readonly -count=1
./...` passed separately for each module using the existing host cache, without
fetching or pointing standalone commands at root workspace vendor material.

## 17. Make and contributor commands

All public targets remain. `test` and `race` cover both modules; Python discovery
remains active. `architecture` is added to `check`; `generate` is an explicit
maintenance wrapper. Check/spec-check/test/conformance/verify/race/fuzz never
execute generation. `make -n check` and source review confirm this.

README/CONTRIBUTING/conformance README explain independent ownership, duplicated
literals, serialized artifacts, workspace limitations and deliberate generation.
The protected `vectors/FORMAT.md` retains its old command byte-for-byte;
CONTRIBUTING explicitly supersedes that path with `make generate` or the sibling
command. No historical report was rewritten.

## 18. Nix/filter/vendor impact

The source filter retains protected artifacts, Python, Makefile, workspace and
required archived history. It includes repository-wide Go/module/native inputs,
including generator and guard tests: a forbidden root/third module must remain
visible to the guard even inside filtered qualification. It does not include
review reports generically. Formatting/vet now cover both modules; Make supplies
both modules' tests/race. Linux gate, offline vendor policy, compiler environment,
flake.lock and dependency versions are retained.

Non-Nix `go work vendor -o <outside-path>` ran once against an outside archive of
starting HEAD and once against final modules. Vendor path/content/mode inventories
are identical, including modules.txt. Hence the existing single vendorHash
`sha256-637XdgBok7hRhMw5s9hdY3Mr5jlPXqc1viSYoYh9Iqw=` is retained for human
confirmation, not assumed Nix success. No second dependency derivation/hash or
package materialization was introduced. An isolated vendored-tree guard and its
tests also passed with GOPROXY/GOSUMDB off and an empty module cache.

## 19. CI

Authoritative Linux Nix remains. All five supplemental combinations, Action pins,
Go/Python setup versions and existing Python/consumer invocations remain unchanged.
Formatting includes both modules; added generator vet/tests, architecture guard,
and standalone readonly tests. No generator go.sum exists, so the dependency-cache
input remains the actual conformance/go.sum. Portable jobs retain their prior
coverage; all-package race and three bounded fuzz targets remain in authoritative
Linux qualification. No ordinary CI step executes the generator.

## 20. Python preservation

All five tools/*.py files are byte-identical to starting HEAD and remain active in
Make/Nix/CI. All 15 tests, including all three test_r16_vectors independent wire
oracle methods, passed at baseline and final qualification. No schema replacement
or oracle migration was attempted.

## 21. Focused evidence and failures

Generator tests compare every committed crypto intermediate/envelope and canonical
input, every post-AEAD fixture, and every bootstrap VAULT_ID against existing
Python-checked artifacts. Focused framing/constraints, source JSON refusal,
deterministic indentation/newline/hash, and graph identity tests pass. Mutation
sensitivity exercises wrong tag/width/endianness/domain/object ID/object-key domain/
nonce/AAD/semantic length/padding/object size against committed fields; no verifier
API or self-generated expected answer is used.

Preserved exceptions in outside evidence and the tool transcript:

- Initial guard execution rejected the existing .direnv module cache as an escape;
  ownership was corrected to recognize Go's actual cache path without permitting
  replacement modules or new dependencies.
- One attempted focused guard invocation used GOFLAGS=-mod=vendor in the host
  checkout, which has no vendor directory. Go rejected loading before tests.
  No dependency or verification policy was relaxed; a properly isolated vendored
  tree with an empty module cache then passed.
- First deliberate generation comparison found 28 path/hash differences because
  an import-qualifier replacement accidentally changed the literal `graph.` ID
  prefix. This was a generator migration bug; no committed artifact was touched.
  The prefix was restored and a focused regression added. A verifier invocation
  also ran on that first temporary tree, but is not accepted as parity evidence.
  Subsequent comparisons/verification passed. Failure tree and comparison log
  remain outside the checkout.

No baseline/final check, vet, race, or fuzz failure occurred. No failing fuzz input
or committed corpus change was required.

## 22. Temporary full generation and serialized boundary

Three deliberate new-generator invocations used separate temporary copies of
spec/vectors/requirements, preserving retained source answers. A disposable extra
case file exercised the unchanged removal/keep inventory behavior. Attempt 1 had
the prefix bug above; attempt 2 and the final stabilized attempt each produced
**zero differences across all 100 protected files**, identical inventories and
92/92/92 case counts. Successful outputs were consumed from disk by the unchanged
verifier, both PASS 92. The final comparison is
`/tmp/totipo-w3a-evidence/comparison-final.log`; final artifact tree is
`/tmp/totipo-w3a-evidence/regenerated-final/`.

The ordinary final `make check` independently verified committed artifacts without
executing generation. Generator does not run the verifier; verifier does not run
the generator. Deliberate shell orchestration supplies disk artifacts between
independent processes. Both roles' tests compile through their own modules.

## 23. Exact committed-artifact comparison

`protected-post-suite.log` compares every protected file path and SHA-256 to the
pre-edit capture after final host qualification: **100 files, zero differences**.
Full raw-byte equality was also checked through Git against starting HEAD.
No additions/deletions, schema-byte changes, protocol revision, vector ID,
classification, crypto/graph/storage/RFC answer, manifest expectation/hash, or
required profile ID/hash/pin changes occurred. This includes protected prose such
as vectors/FORMAT.md. No formatting exception was taken.

## 24. Verifier non-change evidence

All pre-existing conformance/internal/{tlv,object,cryptov1,graph,storage,totp,vectors}
source and tests, cmd/totipo-conformance, go.mod and go.sum match starting HEAD
exactly. Only the generator command was removed from conformance; guard/CLI/tests
were added and current conformance README was updated. Normal committed corpus,
Python oracle, all-package test/race and unchanged consumer fuzz all pass.

## 25. Architecture acceptance matrix

| Predicate | Final host result |
| --- | --- |
| Repository modules / workspace members | 2 / 2 |
| Shared repository Go modules/packages | 0 / 0 |
| Generator → verifier imports (including tests/tags/generated files) | 0 |
| Verifier → generator imports | 0 |
| Common repository package dependency | 0 |
| Local require/replace or workspace replace bridges | 0 |
| Generator subprocess use of verifier | 0 |
| Verifier subprocess use of generator | 0 |
| Generator nonstdlib dependencies | 0 |
| Common external nonstdlib packages | 0 |
| Python checks/oracle active | all retained |
| Ordinary qualification regeneration | none |

## 26. Manual independent-implementation review

Reviewed generator composition against normative v1/r19 source, separately from
byte comparison. Locations below are relative to repository root; W denotes
`generator/internal/wire/writer.go`, V denotes
`conformance/internal/cryptov1/crypto.go`. Corresponding verifier files are unchanged.

| Sensitive area / normative spec section | Generator location | Verifier location | Independent ownership/structure review |
| --- | --- | --- | --- |
| u16 TLV and integers / §12 | W:32 Integer, W:43 Encode | internal/tlv/tlv.go:28,38 | Separate shift loop/append framing; verifier uses binary.Put into allocated header. Own widths/tags. |
| Canonical TOKEN / §12 | W:68 Token.Encode | internal/object/object.go:61,113 | Authoring-field checks and field sequence/sizing; no verifier Dispatch/cursor or parser self-validation. Own enum/text/parent/optional rules. |
| Capacity/semantic length / §§11–12 | W:94 sizing, W:131 Padded | object/object.go:13,61; V:59 | Independently summed writer fields and concatenated length/body/zeros, versus verifier bounds/fixed buffer. Own 1006/1008 literals; committed maximum 1005 tested. |
| HKDF extract/expand/domains / §10 | W:114–127 MAC/block/Derive | V:24–38 | Separate stdlib one-block HMAC composition versus verifier x/crypto HKDF/io.ReadFull. All domain bytes locally declared. |
| OBJECT_ID / §10 | W:142 MAC over P | V:75,105 | Separate semantic input selection and composition; shared standard HMAC/SHA primitives only. |
| Object key / §10 | W:129 ObjectKey | V:48 | Generator domain/suffix/block expansion locally owned; verifier generic hkdf expansion remains unchanged. |
| Nonce / §10 | W:160 id[:12] | V:78,95 | Exact formula necessarily simple and equal; separate slicing call sites, checked against committed nonce fields and wrong-nonce encryption. |
| AAD / §10 | W:130 AAD, W:160 | V:51,78,95 | Separate ASCII literal and ID concatenation. No shared constant/helper; wrong-AAD encryption differs from committed bytes. |
| Padding / §11 | W:131–135 | V:59–67 | Concatenated u16/semantic/zero tail versus fixed buffer Put/copy; separate capacities and zero policy. |
| AES-GCM envelope / §11 | W:137 Seal, W:148 Encrypt | V:52,68,80 | Own width checks/composition, separate cipher construction; stdlib AES/GCM primitives allowed. Full ciphertext/tag/1024-byte known artifacts compared. |
| VAULT_ID / §§6,9 | W:162 VaultID | V:41 | Separate length/string-header selection and exact SHA call; unavoidable simple formula, own magic/version/length. No bootstrap derivation ancestry added. |
| Post-AEAD malformations / §§11,13 | W:171 Malform; generator main:73 postAEAD | vectors/post_aead.go:17 plus V:80 Open | Independent writer mutates one plaintext/ID predicate then seals; unchanged verifier decrypts, checks isolation and strict rejection. All three exact artifacts compared. |

The new writer/composition was authored from the normative formulas rather than
copying/repackaging tlv.go/object.go/crypto.go. Straightforward primitive calls,
ASCII concatenations and required formulas naturally resemble the verifier;
independent ownership does not require obscurity. Artifact declarations and
explicit expectation planning intentionally retain old ancestry/provenance.
This establishes internal independent-implementation evidence, not full external
interoperability or proof against a common specification misunderstanding.

## 27. Final host qualification

Host: Go **1.26.8**, Python **3.14.7**, Linux amd64. One complete final sequence:
`make check`, vet in both modules, `make race`, `make fuzz`, formatting for both
modules, `git diff --check`: all PASS. No incomplete final attempt or rerun.

Go tests/race: generator cmd/generate-vectors, internal/fixture, internal/wire
all PASS; conformance cryptov1, graph, object, repopolicy, storage, tlv, totp,
vectors all PASS. check-architecture and totipo-conformance compile with no unit
tests in their command packages. Architecture guard PASS. Python 15/15 OK;
structure/schema/full conformance/verify-only all PASS; case count exactly 92.
Both vet results PASS. Fuzz budgets/results are recorded in §4. Final generation
comparison and disk-consumer proof PASS; protected comparison repeated after the
host suite PASS. Full execution logs reside outside the checkout.

## 28. Execution accounting

Counts describe top-level agent invocations, not subprocess Go inspections inside
guard tests, static inventory scripts, or Go's cached compilation.

| Execution | Count / status |
| --- | --- |
| Baseline make check / vet / all-package race / make fuzz | 1 / 1 / 1 / 1; all PASS |
| Baseline bounded targets Dispatch / Open / graph | 1 / 1 / 1; all PASS |
| Focused Python unit-test invocations | 0 |
| Focused Go generator test invocations | 3; all PASS (first was compile-only before tests existed) |
| Focused Go guard test invocations | 6 attempts: 5 PASS, 1 pre-test vendor-loading failure |
| Focused bounded fuzz outside baseline/final | 0 / 0 / 0 |
| Standalone readonly dependency listings / test invocations | 2 / 2; one each per module, all PASS |
| Final complete host suites / incomplete attempts | 1 / 0 |
| Final make check / generator vet / conformance vet / make race / make fuzz | 1 / 1 / 1 / 1 / 1; all PASS |
| Final bounded targets Dispatch / Open / graph | 1 / 1 / 1; all PASS |
| Deliberate full generation | 3: first comparison failed, two byte-identical comparisons PASS |
| Temporary serialized consumer CLI runs | 3: first excluded from parity evidence; subsequent two PASS |
| Non-Nix workspace vendor generations | 2; starting/final dependency trees identical |
| Human Nix for W3A | 1; PASS, human-reported |
| Agent Nix | 0 |
| Cached W3A Nix evidence | none observed |
| Remote CI execution/dispatch/push | 0; NOT YET RUN |

Architecture CLI diagnostics additionally included an initial local-cache-path
false-positive and a corrected PASS, later PASS after policy updates, isolated
vendored PASS and final make-check PASS. These are guard runs, not hidden complete
qualification attempts. Earlier human Nix is inherited evidence only.

## 29. Frozen-input inventory

Before the human gate, all intended new files, including this report, are made
Git-known using intent-to-add only. Frozen included and supporting inputs are
recorded outside Git in `/tmp/totipo-w3a-evidence/frozen-inputs.json`: path,
existence/deletion state, mode and SHA-256. This covers every Git-known file except
this excluded report, including generator/conformance/tools, Makefile/go.work,
flake.nix/lock, workflow, current contributor docs, required archived history,
and all protected artifacts. `frozen-membership.json` captures Git path membership;
`frozen-summary.json` records the inventory digest/count: **277 frozen paths**
(including the deleted old command path), **278 Git-known membership paths**, and
inventory SHA-256
`2902fd6be68d7f2b925ee2b74d287509df543cb3d0402ba90eb1f82b9fc523a6`.
The vendor comparison records **75 byte/mode-identical files**. No temp generation,
vendor tree, cache, result, local evidence or secret is intent-added.

This report is excluded from qualificationSource and neither structure checker
nor semantic verifier consumes it. After human PASS, all 277 frozen path/state/
mode/hash records, 278 Git-known membership paths, all 100 protected artifact
files, unchanged branch/HEAD, no unintended untracked files, and the empty cached
diff were verified before updating only this report. All checks PASS. No included
input changed after human PASS; report finalization does not invalidate the gate.

## 30. Human Nix gate

**PASS — human-reported; human execution count 1.** The user reported
“nix flake check passed” for the requested command from the repository root:

```sh
nix flake check
```

Agent Nix execution count remains 0. Cache reuse versus fresh derivation execution
was not reported and is not inferred. No vendor hash correction or included-input
change was needed. Post-Nix frozen-input/protected-byte/Git-membership/index
verification PASS, as recorded in §29. Only this excluded report was finalized
afterward; the same checks passed again after its update.

## 31. Remote CI

**NOT YET RUN** for these uncommitted changes. No push or workflow dispatch.
The user's later commit/push will provide authoritative remote Linux and
supplemental macOS/Windows/Go portability evidence. Local Linux success cannot
establish those outcomes.

## 32. Final Git/index state after human qualification

Branch/HEAD remain main / `a20aa62543219241b08237ccdbba4b6e72706537`.
All existing modifications/deletion remain unstaged. All 12 new intended files
(11 source/module/test files plus this report) are intent-added only.
`git diff --cached --quiet` PASS; cached file list empty; no unintended nonignored
untracked file. Nothing committed, pushed, reset, or tagged. Final status and
membership are captured outside the checkout; post-human verification matches
the handoff capture. No staged content exists after report finalization.

## 33. Limitations

Human Nix PASS is user-reported; cache/fresh-build details were not provided.
Post-PASS freeze verification passed, and no included-source edit followed it.
Remote matrix is pending. No external interoperability claim, independent
Argon2 implementation, Python removal, schema migration, or protocol/corpus
revision is claimed. Temporary evidence is local and is not a durable committed
artifact; this report preserves its conclusions and counts. Architecture checks
constrain ownership and escape mechanisms, not arbitrary malicious runtime reuse
or implementation ancestry beyond reviewed code.

## 34. Suggested commit subject

`Isolate spec generator from verifier`

Leave committing to the human after review. All changes remain uncommitted.
