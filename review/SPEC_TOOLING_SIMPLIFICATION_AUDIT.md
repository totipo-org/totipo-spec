# Specification tooling simplification audit

Date: 2026-10-14. Classification: **DOCS_ONLY**; design-only, read-only source audit.
Only this report is added. No recommendation is implemented.

## 1. Executive summary

The protocol and fixture architecture is already reasonably small. The repository
has two Go modules, three small commands in the consumer, one deliberate writer,
one Makefile, one qualification derivation, and one CI workflow. There is no live
Python tooling, release state machine, artifact publisher, general task framework,
or separate binary distribution pipeline to retire.

Useful simplification remains at the orchestration boundary: full conformance
already performs all verify-only checks, yet Make runs both, and portability CI
runs verify-only twice around full conformance. Dedicated CLI tests now cover the
distinct verify-only behavior that earlier migration reports required before
deduplication. The largest optional reduction is deleting the exact Makefile
allowlist; it duplicates commands and freezes innocuous spelling and help text,
but also provides a real accidental-regeneration tripwire. That is an owner
policy choice, not an assurance-free deletion.

Recommend **Option A**, a bounded orchestration/documentation pass, preserving
all semantic, malformed-input, schema, identity, independence, race and fuzz
coverage. Option B is credible under stronger reliance on code review, but its
small additional savings do not justify making it the default recommendation.
Do not merge modules, share protocol helpers, replace the schema interpreter,
or change the corpus contract to save lines.

Explicit answers:

| Question | Finding |
| --- | --- |
| Does generator/verifier separation achieve useful semantic independence? | Yes. Independent models, constants, TLV/TOKEN composition and HKDF/envelope construction; no implementation imports or process delegation between roles. Scope is Totipo composition, not independent AES/Argon2 implementations or independent human interpretation. |
| Which custom checks protect correctness? | Schema presence/types/ambiguity, strict TOKEN grammar, crypto intermediate and envelope checks, post-AEAD defect isolation/rejection, authored graph/fold/storage outcomes, RFC TOTP answers, required-case/profile equality, raw pins, and independent oracle/mutation checks. |
| Which enforce preferred structure? | Exact module/workspace names and Go directive, exact dependency-version duplication in policy, source-root inventory, exact allowed subprocess filename, and especially the complete Makefile text allowlist. Some enforce useful independence boundaries as well; they are not all interchangeable. |
| Are fixture identities repeatedly checked without new confidence? | Yes across repeated normal CLI invocations. Within one acceptance pass, manifest case hashes and profile artifact/set checks establish different relationships; the profile does not rehash every case. Read-only snapshots and tampering tests establish different failure properties. |
| Can Go/Nix/Git replace support machinery? | They already handle compilation, package loading, formatting, dependency checksums, fixed vendoring and history. Remove repeated acceptance invocations; optionally rely on review instead of Make spelling/version duplication. They do not replace protocol acceptance or semantic independence checks. |
| Is the repository sufficiently simple? | Its semantic and distribution architecture is. A small cleanup is warranted; another broad migration or new framework is not. |

## 2. Starting HEAD, revision and source state

Required initial inspection established:

```text
Repository: /home/niki/Sources/totipo-spec
Branch: main
Full HEAD: d5b26b315b44f37df8cd793c1d458f3f8602e2af
git status --short: empty
Tracked starting files: 288
```

AGENTS.md was read completely. Current README, CONTRIBUTING, conformance README,
requirements guidance, vectors/FORMAT.md, active revision/scopes/history,
Make/Nix/CI definitions and current module implementations/tests were inspected.
Wave 3 architecture and A/B/C/D reports supply historical context, not a current
implementation inventory. No sibling repository, sibling worktree, historical
checkout or delegated agent was inspected or used.

Current normative specification: `spec/totipo-vault-format-v1.md`, **protocol 1,
v1/r19**, design draft. Manifest: `totipo-vector-manifest-v1`, `totipo-v1`, r19.
Profile: `totipo-requirements-v1`, `moving-pre-rc`, r19. All **92 manifest cases,
92 physical JSON cases, and 92 required profile entries** remain present.
Categories: bootstrap 5, crypto 5, encoding 30, fold 6, graph 13, metadata 7,
size 1, storage 13, TOTP 3, vault 9. Protected artifact inventory is 100 files:
spec 1, vectors 97, requirements 2. Counts are static observations, not a fresh
conformance pass.

Both modules and `go.work` declare Go **1.23.0**; no toolchain directive.
`totipo/generator` is stdlib-only. `totipo/conformance` requires
`golang.org/x/crypto v0.36.0` and indirect `golang.org/x/sys v0.31.0`.
Its four-line go.sum supplies normal module/content metadata checksums.
Reviewed imported external primitives are argon2, hkdf, their blake2b dependency,
and sys/cpu. This is source inspection plus inherited resolved-closure evidence;
no new `go list` qualification was executed.
Observed executable: `/nix/store/iqq9yz2yrcrzsn52h3vkkz72i92h1zds-go-1.26.8/bin/go`;
`go version` reports **go1.26.8 linux/amd64**. This is the host observation,
not a newly evaluated version for every flake system.

Absent: root go.mod/go.sum/go.work.sum, generator/go.sum, tracked `tools/`, and
any current Python support command. `.envrc` is tracked and reads `use flake`
and `layout go`. Existing workflow: `.github/workflows/conformance.yml` only;
`.github/dependabot.yml` schedules weekly conformance gomod and Actions proposals.

Nix root inputs resolve through these lock nodes (commit identity is the version):

| Input/node | Locked revision |
| --- | --- |
| Root nixpkgs → nixpkgs_3 | 39ad350a0602fa0a58a544344e3e9187526ea45c |
| flake-utils and flake-utils_2 | 11707dc2f618dd54ca8739b309ec4fc024de578b |
| jailed-agents | 82b9d99d454bd752ae028e8a4bf0c90daf1c23ef |
| llm-agents | 4bb57cff45b5554a02dba0cf7a8c7f4f010a864d |
| jailed-agents nixpkgs → nixpkgs | a799d3e3886da994fa307f817a6bc705ae538eeb |
| llm-agents nixpkgs → nixpkgs_2 | f45c6f04c2f013f004bf94e284e95d72898d9393 |
| jail-nix | 42b355c38ca63dab4904acc5c0d95f17954a8c9b |
| bun2nix | 07a5bfc8ac36c5370199343fa620b69f63186e33 |
| flake-parts | 024633cd702b10285db5cb19b40ad48d2399ba60 |
| treefmt-nix | 4cfe8d7065cc8ca91c0a50361b6dd74920eeb1f5 |
| systems and systems_2 | da67096a3b9bf56a91d16901293e51ba5b49a27e |
| systems_3 | d36eeed142f88b5a636f2e3856a53c61654cf29b |

All locked nodes carry narHash values in unchanged `flake.lock`. No Nix command
was run. Workspace vendorHash is
`sha256-637XdgBok7hRhMw5s9hdY3Mr5jlPXqc1viSYoYh9Iqw=`.

Inherited W3D report records passing host suites and one human-reported
`nix flake check` PASS on 2026-10-12; cached versus fresh provenance is unknown.
It leaves exact-commit Linux/macOS/Windows CI pending. This audit does not query
or upgrade remote status, and executes no replacement qualification suite.

## 3. Current generator/verifier architecture

```text
go.work: developer convenience, two sibling modules
  generator/cmd/generate-vectors
    -> generator/internal/fixture (writer JSON records, source reader, hash)
    -> generator/internal/wire (writer-owned TOKEN/crypto composition)
    -> stdlib
        writes serialized cases, manifest, moving profile only deliberately
  conformance/cmd/totipo-conformance
    -> internal/vectors.Read + VerifyProfile
       -> internal/schema, serialized artifact checks
    -> internal/vectors.Run
       -> tlv/object/cryptov1/graph/storage/totp
  conformance/cmd/check-spec -> internal/specpolicy -> spec/archive only
  conformance/cmd/check-architecture -> internal/repopolicy -> source + Go metadata
  conformance/internal/wireoracle tests -> own reader/formulas -> stdlib
```

Entrypoints: `make generate` / `go run ./generator/cmd/generate-vectors -root .`;
`make conformance` / consumer full mode; `make verify` / consumer `-verify-only`;
`make spec-check`; `make architecture`. Normal qualification never invokes
generator main. Generator tests compile the command and exercise writer helpers;
that does not execute main or regenerate committed artifacts.

No shared repository semantic package, cross-role import or subprocess invocation
was found. `os/exec` occurs only in repopolicy's Go inspection helper. Schema and
spec policy are verifier-owned; the generator neither imports nor invokes them.

## 4. Tooling inventory and LOC

Physical LOC includes comments/blanks; approximate component costs below were
counted from current files. Total Go: **6,619 LOC**, production **3,786**, tests
**2,833**. Generator: 903 production + 294 test; consumer: 2,883 + 2,539.
Normative prose, schemas and ordinary vectors are not charged as wasteful tooling.
Test cost means retained fixture/work scope, not measured runtime in this audit.

| Source path(s); LOC production/test | Callers; inputs → outputs | Concrete failure/boundary; overlap; fixture cost |
| --- | --- | --- |
| generator/cmd/generate-vectors/main.go, main_test.go; 487/38 | make generate, explicit human command; retained answers + authored declarations + spec/schema bytes → cases/manifest/profile | Deterministic construction and authored outcomes; no independent acceptance. No general generator framework; two helper tests, temporary full comparison during writer qualification. |
| generator/internal/fixture/{records,json}.go, json_test.go; 226/24 | generator command/writer tests; JSON records/source bytes → records/hash | Writer owns models and serialization order; duplicate-safe TOTP source reading. Overlaps mechanical JSON/hash primitives, not verifier internals; source refusal/field-order tests. |
| generator/internal/wire/writer.go, writer_test.go; 190/232 | generator/its tests; explicit fields/root → TOKEN/envelope/ID and intentional defects | Independent composition, authoring limits and single-defect negatives. Same standard primitives as consumer; reads committed corpus, constraint and mutation tests. |
| conformance/cmd/totipo-conformance/{main,main_test}.go; 51/205 | Make, Nix through Make, CI, users; artifact root → status/diagnostics | Real CLI acceptance/refusal, mode difference, read-only contract. Four tests include temporary artifact copies and before/after snapshots; no binary packaging framework. |
| conformance/internal/vectors/{types,runner,post_aead}.go + four test files; 834/565 | CLI/tests; schemas/manifest/profile/cases → validated records and semantic results | Identity, coverage, containment, exact outcomes, authenticated defect isolation. Protocol-aware acceptance; repeated CLI calls overlap this same work. Full corpus, pin/set/path mutations, real encrypted folds. |
| conformance/internal/schema/schema.go + schema_test.go/contracts_test.go; 459/329 | vectors.Read, schema tests; authored schema + raw JSON → acceptance/refusal | Required-member presence, duplicate rejection, exact integer tokens, conditions, fail-closed subset. Go decoder alone cannot establish presence of false/zero/empty values. Corpus required-member removals and focused negatives. |
| conformance/internal/{tlv,object,cryptov1,graph,storage,totp} production and tests; 757/654 | vectors.Run, cross-package integration/tests; bytes/facts → protocol results | Exact encoding, hostile input rejection, KDF/envelope/ID, graph/fold/store/TOTP semantics. Uses standard primitives, not compiler emulation. Negative/unit/integration tests, 512 three-node graph oracle, three bounded fuzz targets. |
| conformance/internal/wireoracle/oracle_test.go; 0/364 | go test/race; committed serialized cases → selected independent formula checks | Own duplicate-safe map reader, field/key/padding/VAULT_ID/capacity checks; no production imports. Three corpus traversals plus reader/mutation tests; no independent AES-GCM/Argon2. |
| conformance/internal/specpolicy/policy.go, policy_test.go; 283/145 | check-spec/tests; active spec + archived history → prose policy result | Metadata, registry, history, stale-rule/wording regressions. Not protocol authority or application behavioral proof; exact wording and Python compatibility costs. Five policy test families. |
| conformance/cmd/check-spec/{main,main_test}.go; 32/47 | Make/CI; root/args → policy status | Two-input CLI/refusal evidence; does not need corpus or evaluator. One CLI family in a minimal temporary tree. |
| conformance/internal/repopolicy/{policy,wiring}.go + tests; 444/230 | check-architecture/tests; source/Make/modules/workspace/resolved test deps → policy status | Cross-role/shared implementation and accidental generation restrictions. Relies on Go parser/list/edit commands; additional layout/version/spelling policy. 32 source/metadata mutations, 13 wiring mutations, closure/symlink/vendor positives/negatives. |
| conformance/cmd/check-architecture/main.go; 23/0 | Make/CI; root → guard status | Small adapter; no command tests, policy tested below it. |
| Makefile; 40 | developer/CI indirectly/Nix; Go commands → test/verification results | Discoverable orchestration and exact race/fuzz scope; check redundantly calls verify after conformance. No custom test runner. |
| flake.nix; 133, flake.lock; 279 | direnv/developer/human/Linux CI; filtered source + locked packages → vendor closure/qualification marker/shell | Fixed dependencies/offline build and complete Linux suite. Standard buildGoModule/stdenv, no vector artifact production. One dependency derivation, one qualification derivation; three repeated dev/jail utility lists. |
| .github/workflows/conformance.yml; 71 | push/PR/manual; commit + platform/toolchain → CI evidence | Authoritative Linux Nix + five supplemental combinations. Repeated workspace/standalone tests and CLI checks; standard GitHub jobs, no workflow parser. |
| .github/dependabot.yml; 10, .envrc; 2 | GitHub/direnv | Standard dependency proposals and shell activation; no custom integrity layer or fixture cost. |
| go.work, generator/go.mod, conformance/go.mod/go.sum | Go/Nix/guard | Standard module identities/checksums; independently owned module requirements, root workspace. Policy repeats exact pins/directives. |
| README.md 109, CONTRIBUTING.md 153, conformance/README.md 76, AGENTS.md 175 | humans/agents | Development/qualification contracts; excluded from qualificationSource and not read by current structure checks. Historical reports are preserved context. |
| vectors/FORMAT.md 125, requirements/README.md 13 | implementation consumers/humans; included directory inputs | Corpus/profile guidance. FORMAT has stale generator path/shared-composition description; directory inclusion means editing it invalidates Nix source. |

There is no current `tools/` script inventory or release/revision executable.
`design/icons/` has its own branding manifest/SHA256SUMS and image assets; no
Make/CI/Nix qualification caller references them. They are not protocol corpus
checks and are not candidates for a conformance cleanup.

## 5. Normative versus support responsibilities

The active spec is normative. Section 20 distinguishes core, store/writer and
application scopes, requires appropriate evidence, and limits portable corpus
claims. The moving profile defines all-and-only required cases for this repository
revision. Serialized schemas define fixture shape/presence. Cases supply exact
bytes and authored outcomes. The Go consumer is reference evidence, not an
alternative authority that can silently redefine prose or repair expectations.

Make, architecture policy, Nix, CI, history prose and developer shells support
qualification. A policy test passing does not authorize a semantic revision.
Git commit/tag identity and reviewed history anchor intended contents externally
to the files being checked; locally rewritable pins and guards are regressions
against accidental drift, not authorization against an author controlling both.

## 6. Semantic independence evaluation

**A — Semantic independence.** Writer owns Token, JSON records, constants and
ordered field construction; reader owns Object, cursor/TLV parsing, validation
and models. Writer validates authoring fields before encoding; it has no accepting
TOKEN parser. `SplitField` is a bounded fixture mutation utility, not grammar
acceptance. Writer uses explicit HMAC extract/single-block expand, append-built
padding and domains; consumer uses x/crypto/hkdf and separately selected domains,
sizes, nonce/AAD and acceptance checks. Both use stdlib AES/GCM/HMAC/SHA-256.

Generator graph heads/values/missing parents and storage/workflow outcomes are
authored declarations. `graphCase` attaches records for already declared heads;
it does not compute reachability. Conflict is computed from declared distinct
values. Fold's 4/3 staging arithmetic is independently implemented, though its
algorithm resembles the simple specified consumer algorithm. Bootstrap records
are retained source answers, with VAULT_ID recomputed; RFC TOTP answers are retained
and their files are not rewritten. Generator does not compute Argon2 or TOTP.

**B — Repository/module separation.** Sibling internal packages prevent ordinary
cross-module imports of internal helpers. They do not prevent public packages,
shared third modules, copied algorithms or process delegation. The source guard
and closure restrictions add meaningful regression protection here. `go.work`
itself is convenience and is not isolation.

**C — Build-system separation.** Direct verification builds only the consumer.
Ordinary repository tests also test the writer without executing main. The shared
workspace and Nix vendor derivation do not share Totipo implementation semantics.
Standalone GOWORK=off checks prove each module declares enough dependencies for
its own build; they do not add a second independent semantic evaluation.

**D — Layout fixtures.** Synthetic two-module trees, exact roots/directives,
Make text and the precise repopolicy filename enforce policy conventions. Their
value depends on a concrete coupling they reject. Oracle import tests, public
cross-role imports, shared external semantic package/replace bridges and
process/native escape tests have one. Wrong Go directive and changed help text
do not establish additional semantic independence.

Behavioral evidence: writer TestCommittedComposition compares independent
construction to committed fields; oracle mutation tests reject wrong domains,
fields and repaired defects; consumer metadata/defect tests reject lost semantics;
graph tests compare traversal against matrix closure over 512 graphs. These are
useful, bounded evidence. No guard proves independent intellectual ancestry,
full cryptographic interoperability, live backend crash safety or GUI compliance.
Both future options preserve independent implementations and serialized interfaces.

## 7. Generator determinism

Case IDs/categories are deliberate strings, output cases are sorted by ID, parent
sets are deterministically sorted where constructed, JSON uses ordered writer
structs with two-space indentation and a final newline, and crypto uses fixed
public roots plus explicit retained salt/nonce answers. There is no clock/random
input or timestamp in output content. Filesystem mtimes are not part of identity.
Only RFC cases bypass rewriting to preserve their original bytes.

Manifest hashes are calculated from actual written/read file bytes. Required pins
reuse the manifest entry hashes; cases are not hashed again for each profile pin.
Spec, manifest and both schema bytes each get an artifact pin. The keep-set walk
removes undeclared physical cases during explicit generation. These responsibilities
are direct and small; no unnecessary abstract provider system was found.

Limits: `must` panics on failure, writes are incremental and cleanup is destructive
within the requested root, so an interrupted run may leave mixed artifacts. This
is acceptable for deliberate temporary-tree maintenance followed by independent
verification and reviewed Git differences; generation is not an atomic release
transaction. Do not create recovery/state machinery for it.

Bootstrap source cases use permissive `json.Unmarshal`, unlike duplicate-safe
`fixture.ReadSource` for RFC cases. This is a consistency gap in trusted source
ingestion, not evidence of nondeterminism in current committed inputs. A future
writer-hardening choice could use its own strict helper; do not import verifier
acceptance or silently change retained answers. It is outside the recommended
orchestration pass.

The deterministic-write unit test covers serialization, not complete regeneration.
Historical W3A temporary full generation achieved byte equality after diagnosing
a case-prefix bug. This audit neither repeats nor upgrades that result. Writer
changes still require temporary full generation/inventory comparison, retained
source answers, reviewed intended differences and final committed-corpus verification.

## 8. Independent conformance verification

`vectors.Read` compiles both schemas, validates bytes before typed decoding,
rejects duplicates/trailing values/unknown typed fields, requires r19, validates
case paths and symlink containment, hashes case bytes, compares ID/outcome, checks
operation shape, and rejects unlisted physical files. `VerifyProfile` independently
requires r19/moving-pre-rc and four artifact pins plus exact required ID/hash set.
Full CLI runs both before semantic execution; verify-only skips only execution.

Semantic execution checks TOKEN fields/order/widths/ranges/UTF-8, canonical
re-encoding, all declared crypto intermediates, authenticated envelope opening,
bootstrap wrap/unwrap/VAULT_ID, RFC counter/code results, graph/SCC outcomes,
fold stage parents/value/metadata, storage observation and immutable publication.
Raw post-AEAD opening proves negatives really authenticate and isolate the named
length/padding/ID defect before strict reader and storage rejection. A simple
"Open failed" assertion would lose that important guarantee.

Apparent duplication across writer and reader is intentional independent
acceptance. Schema validation and typed Decode/ValidateShape overlap some fields
but have distinct APIs: schema establishes exact member presence; Decode protects
direct calls and typed interpretation; shape establishes operation relationships.
Read currently parses each manifest/case for schema and again in Decode before
typed decoding. This is a small performance opportunity, not a priority: changing
it risks weakening the direct Decode contract for negligible corpus cost. Keep it.

Unit/integration tests that use consumer encode/open are local regressions, not
independent producer evidence; committed writer outputs and oracle formulas add
that evidence. Negative tests are not redundant just because they reuse pins:
repinning temporary mutations deliberately isolates schema/semantic failures
past integrity gates. Preserve them.

## 9. Corpus and manifest integrity

| Property | Current evidence | Recommendation |
| --- | --- | --- |
| Intended normative revision | Reader expectation r19, spec header/policy, profile exact spec/schema/manifest pins, external Git revision | Keep. Incoming generator output cannot select a supported revision automatically. |
| Complete required coverage | Listed IDs/paths unique; physical inventory equal; required profile set exactly matches manifest; Section 20 obligations | Keep. Set equality proves declared coverage, not that humans authored a case for every prose requirement. |
| Deterministic generation | Stable writer inputs/order/serialization; writer regressions; deliberate temporary full comparison | Keep separately from normal acceptance; no routine regeneration prerequisite. |
| Independent correctness | Reader evaluates committed bytes/outcomes, oracle/formulas, authored graph/storage answers, external interoperability before RC | Keep; locally rewritable pins alone cannot supply this. |
| Identical distributed bytes | Case raw hashes + profile raw artifact hashes + Git commit/tag pin | Keep existing formats. No new attestation or binary artifact needed. |

Per normal CLI invocation there are 92 case hashes plus four artifact hashes.
Profile case entries compare existing manifest hashes without rereading cases.
Hashing spec/schema/manifest is not redundant with hashing individual cases: it
also pins revision prose, acceptance contract, classifications and required order.
Repeating Read/VerifyProfile in another mode on the same stable tree gains no new
identity trust boundary. Under `make check`, TestCorpus, CLI unit tests, full CLI
and verify CLI all revisit artifacts. Keep behavior-focused tests; remove only
unnecessary top-level invocation repeats.

CLI before/after path/hash snapshots test read-only behavior, including rejected
input; they are not competing corpus identity authority. They do not detect a
transient write restored before return or every possible future input. Keep the
test, and optionally make its fixed 100-file count derive from starting inventory
in a future revision cleanup. Hardcoded 92/100 counts currently flag corpus drift;
they are not an exact implementation-layout mandate worth a separate project.

Git distributes source-controlled conformance material. Nix's qualification marker
is a build result, not a separate released vector archive. No additional fixture
hashing or generated-binary-release identity service is necessary.

## 10. Go, Nix, CI and toolchain acquisition

Nix uses root nixpkgs `go` for nativeBuildInputs and devShell tools, sets
GOTOOLCHAIN=local in vendoring/qualification, supplies Make, enables CGO for race,
and uses fixed workspace vendoring with proxy/sumdb disabled in qualification.
Standard Go module sums and Nix vendorHash address dependency integrity. The
architecture policy's duplicate version maps express reviewed dependencies;
they are not independent integrity verifiers.

Dev-shell utilities are repeated in outer buildInputs and both jail extraPkgs;
Go is the same root pkgs expression in each. Factoring one local support-tool
list would reduce drift without changing acquisition. Different jailed/agent
inputs have separate nixpkgs nodes; do not assert that following all to root
preserves their compatibility. Development does not explicitly set every offline
qualification environment variable, and `layout go` supplies developer GOPATH.
These are different responsibilities, not inconsistent qualification pinning.

Linux Nix is authoritative. CI installs Nix once for that job and runs
`nix flake check --print-build-logs`. Five portability combinations use setup-go:
macOS/Windows × 1.23.x/stable, plus Ubuntu 24.04 × 1.23.x. Stable is deliberately
moving supplemental evidence; setup-go on those hosts is not duplicate acquisition
of the authoritative Nix compiler. Preserve the matrix and pinned Actions.

Portability currently runs verify-only in "Spec policy and schemas", workspace
vet/tests for each module, architecture, standalone tests for each module, full
CLI, then verify-only again. Each ordinary all-package run also exercises CLI
and corpus tests. Keep standalone checks: they establish workspace independence,
especially undeclared dependency reliance. There is no need to test Go compilation
once more merely to confirm the compiler works. Remove duplicate verify-only
steps; retain full CLI as explicit real command evidence. No custom workflow
step-name checker/parser exists to remove.

## 11. Specification release process

Current v1/r19 is a moving design draft, not a frozen release. Local tag inventory
contains only `v0-rc1`; this audit makes no claim about its publication state or
remote GitHub releases. No release workflow, tag utility, publisher, recovery
database or Maven-like staging process exists in the current repository.

Normative authority is the spec; distributed identity is a consumer-pinned Git
commit or reviewed release tag resolving to that commit. The profile binds exact
spec, schemas, manifest and required cases within the tree; its revision label
alone is insufficient because the profile is moving. Revision preparation changes
prose, relevant authored cases and pins under review, with deliberate generation
only when needed, separate acceptance and qualification. Existing requirements
guidance forbids RC freeze until interoperability and required evidence exist.

For a future frozen tag, preserve tagged spec/case/schema/manifest/profile bytes
and tag-to-commit identity; corrections need a new reviewed identity. Git tags
and a GitHub release pointing to the qualified commit are sufficient distribution
for current source material. This is a proposed minimal release convention, not
an observed implemented immutability enforcement service. CI validates checked-out
contents but does not establish owner intent, forbid moving remote tags, generate
release fixtures or independently authorize a coordinated pin rewrite. No custom
release machinery needs deletion, and none should be invented.

## 12. Structural and self-referential checks

Source scan covers inactive tags, generated-marked files and tests, where ordinary
Go dependency loading may not. Cross-role/public/shared imports, external semantic
dependencies, linkname/native/plugin/exec bridges and oracle production imports
are real potential semantic coupling. Metadata/no-replace and resolved closure
checks complement those restrictions. Keep them in the recommended pass.

The complete Makefile copy in `repopolicy/wiring.go` is different. It rejects
generation prerequisites, hidden shell expansion/includes and removal of checks,
but also rejects harmless command/help changes. Its thirteen negatives mostly
prove that changed executable lines fail an exact allowlist, rather than inspect
Make semantics. It cannot authorize against edits to both Makefile and policy;
Nix/Actions wiring and caller-supplied GO/toolchains remain reviewed trust inputs.
Removing it loses an accidental-wiring review gate, not independent TOKEN/crypto
implementations. Go's internal visibility cannot replace that gate by itself.

Metadata repeats exact versions, indirect status, Go 1.23.0 and directive bans.
Version/indirect checks can be delegated to committed module/lock review in Option B;
no-replace, ownership, stdlib generator and oracle restrictions still matter.
Vendor containment guards additionally detect wrong physical source resolution;
historical Windows/macOS fixes show their portability maintenance cost. Keep
containment for now: this audit cannot demonstrate a cheaper equally protective
replacement using only current inputs.

Prose policy anchors are review tripwires for normative requirements, including
application behaviors absent from executable corpus. Exact line/newline wording,
section layout and retired-term lists are conventions; they are not proof of the
meaning of prose. Unicode whitespace/digit and bool/number/final-newline schema
compatibility code preserves measured W3B behavior. Removing compatibility would
change validation behavior; no high-confidence deletion is recommended.

There are no hashes of local implementation source used as authorization in the
current gate, no test of fixed workflow step names, and no bespoke gofmt verifier.
The tiny Hash("abc") assertion in generator fixture tests mostly verifies stdlib
SHA-256/hex plus a one-line wrapper. Removing just that assertion is safe if wrapper
correctness remains covered by full writer byte comparison; maintenance savings
are three lines and do not justify a milestone by themselves.

## 13. KEEP / SIMPLIFY / REMOVE classification

| Component | Classification | Reason |
| --- | --- | --- |
| Spec, schemas, case expectations, profile formats/pins | KEEP | Normative/implementation contract; not unnecessary tooling LOC. |
| Separate writer/reader models/constants/composition/modules | KEEP | Useful semantic independence. |
| Core package negative/integration/oracle tests | KEEP | Behavior and hostile-input guarantees; shared primitives do not make them redundant. |
| Exact inventory/profile identity, path safety, schema presence | KEEP | Detect incomplete/ambiguous/wrong corpus without trusting generator. |
| CLI mode/refusal/read-only tests | KEEP | Dedicated coverage enables invocation deduplication. |
| Make test/race/fuzz/verify/conformance public targets | KEEP | Understandable interface and unchanged final coverage. |
| verify as make check prerequisite | REMOVE duplicate invocation | Full CLI already verifies; public make verify remains. |
| Both portability verify-only runs | REMOVE duplicate invocations | Full CLI verifies and module tests exercise verify-only semantics. |
| Make allowlist/wiring tests | SIMPLIFY only with owner policy approval | Optional deletion of duplicated executable text; automatic wiring protection lost. |
| Exact dependency/Go directive policy duplication | SIMPLIFY optional | Standard module/lock review can own versions; retain anti-coupling restrictions. |
| Nix dev/jail utility lists | SIMPLIFY optional | Same root pkgs tools repeated; factor list, preserve compiler extras. |
| Stale FORMAT generator guidance | SIMPLIFY/correct | Current path/ownership should be directly understandable to implementers. |
| Hash("abc") test fragment | REMOVE optional | Tiny redundant mechanical check; preserve serialization/presence test. |
| Release tooling/obsolete current scripts | None to remove | Already absent; history retrieval instructions are not active compatibility tooling. |

## 14. Option A — Conservative

Change `Makefile` check prerequisites to retain architecture, spec-check, both
module tests and full conformance, leaving verify available separately. Update
the corresponding check line in `repopolicy/wiring.go` and mutation fixture strings
in `wiring_test.go` so generation/architecture/full-verification weakening still
fails. Do not delete boundary tests or replace the allowlist with a new parser.

In `.github/workflows/conformance.yml`, the first step runs check-spec only
(rename its descriptive label); remove the final integrity step. Retain full
CLI, both workspace/standalone test sets, vets, formatting and architecture.
Correct stale command/composition guidance in `vectors/FORMAT.md`; remove the
supersession workaround in CONTRIBUTING and align conformance/requirements
maintenance references if needed. Do not alter fixture schemas or formats.

Retains all semantics, pins, required-case coverage, independent formulas,
module/closure/wiring policy, CLI modes, read-only and malformed-input tests,
Linux race/fuzz and five portability combinations. Loses only repeated top-level
verify-only smoke execution and repeated identity reads on the same inputs.
Assumes dedicated CLI tests plus full command execution are sufficient mode
coverage; those tests are present now. No independent implementation risk from
changed corpus bytes because they must remain identical.

Estimated net code reduction: roughly **5–12 lines**, plus clearer documentation;
one fewer Go CLI process per make check and two fewer per portability combination
(ten per complete matrix). No timing claim. Hash assertion cleanup is optional;
leave test functions otherwise intact. This small gain is worthwhile only as one
bounded pass. Qualification is detailed in §20.

## 15. Option B — Minimal sufficient

Include Option A's invocation/doc cleanup, then delete
`conformance/internal/repopolicy/wiring.go` and `wiring_test.go` (140 LOC), remove
the Wiring call in `policy.go` and the Makefile fixture in `policy_test.go`.
Retain the Makefile and its public targets as reviewed orchestration. No substitute
test framework, YAML validator or Make parser.

Optionally relax only duplicate version/Go-directive assertions in policy.go,
its version maps and the wrong-Go-version fixture; module/lock files still pin
inputs. Keep roles, no cross-role imports/delegation/shared semantic dependency,
no replacement bridges, stdlib-only writer/oracle, source/tag checks and resolved
ownership/containment tests. Factor the common developer tool list in flake.nix
without removing agent packages, changing jail inputs or following nixpkgs nodes.
No go.mod/go.sum/flake.lock update is needed for the design.

Retains all corpus/semantic/independence checks and final qualification scope.
Loses automatic Make executable-line approval and, if version checks are relaxed,
the independent local allowlist of exact primitive versions/Go directive.
New assumption: owner review and GitHub review/CI, standard Go sums and locked
Nix dependencies sufficiently control those ordinary infrastructure edits.
Coordinated changes to guards/fixtures were already trusted; this additionally
allows an accidental Make edit to escape an automatic spelling tripwire. A
CLI snapshot cannot catch orchestration regenerating fixtures before it runs.
That residual risk must be accepted explicitly in a future implementation request.

Estimated net reduction **150–200 LOC** including optional narrow policy/list
cleanup, about **2–3% of Go-plus-small-orchestration source**, with thirteen
wiring mutation fixtures removed. No claim of large runtime reduction beyond
Option A. Exact totals require implementation measurement; no 400–600-line guard
deletion is justified by this audit. Same final coverage and extra independence
negative checks are required after policy edits.

## 16. Recommended target architecture

Select **A**: two independently owned Go modules, serialized artifacts as their
interface, one schema-aware independent consumer, one selected stdlib oracle,
small prose and architecture gates, standard Make/Go/Nix/Git/GitHub support.
Retain fixed corpus identity. Remove redundant normal CLI verification only after
reviewing the concrete tests already in place. Keep B as an explicit later owner
choice, not a mandatory migration. Neither option requires shared constants,
code generation of consumer models, a third module or a new language/format.

## 17. Exact file-level recommendations and callers

| Priority/confidence | Exact files | Caller impact / bounded future action |
| --- | --- | --- |
| 1 high | Makefile; conformance/internal/repopolicy/wiring.go; wiring_test.go | Nix invokes make check, humans invoke Make. Remove check's redundant verify dependency; adjust copied line and weakening/generation fixtures together. Full conformance stays mandatory. |
| 1 high | .github/workflows/conformance.yml | GitHub portability jobs: remove both verify-only calls; keep check-spec and full CLI and existing tests. Linux Nix command/matrix unchanged. |
| 1 optional tiny | generator/internal/fixture/json_test.go | Both module test/race paths: remove only Hash("abc") assertion, keep byte-order/presence test; full temp generation comparison covers wrapper. |
| 2 optional | flake.nix | direnv/dev/jails: factor repeated Go/Make/gopls/lint tool list; preserve libgcc/gcc jail extras and qualification inputs. No compiler update. |
| 3 owner choice B | conformance/internal/repopolicy/wiring.go, wiring_test.go, policy.go, policy_test.go; CONTRIBUTING.md | check-architecture through Make/CI: retire exact Make allowlist and fixtures/call; revise documented guarantee. Keep all semantic isolation checks. |
| 5 owner choice B | conformance/internal/repopolicy/policy.go, policy_test.go; CONTRIBUTING.md | Metadata/resolved closure checking: remove duplicated exact versions/directive convention only; preserve no-replace/semantic dependency boundaries. |
| 6 high documentation | vectors/FORMAT.md; CONTRIBUTING.md; conformance/README.md; requirements/README.md | Human maintenance/implementation contract: correct old generator command and shared-Totipo-composition implication; remove workaround guidance where superseded. FORMAT is a flake source input. |

No obsolete scripts, custom release orchestration, repeated vendor derivations,
or workflow validators were found. Do not fabricate deletions for those categories.
No core semantic/test file or manifest/profile/schema format change is recommended.

## 18. Guarantees lost and risks

A removes repeated top-level execution, not acceptance conditions. CLI argument,
mode, schema failure, evaluation failure and unchanged-artifact evidence remain
in normal tests. Separate verify command remains usable by implementations.

B trades Make/pin convention enforcement for reviewed infrastructure ownership.
Source/module separation alone cannot stop hidden generation in Make or copied
algorithms. Retained source/closure policy protects cross-role coupling but is
still repository-local regression evidence. Shared stdlib crypto and human
specification interpretation remain correlated risks under both options.

Changing schema compatibility or prose-policy anchors would risk acceptance or
governance drift and is excluded from the recommended pass. Removing fixture pins
would weaken the implementation contract and is not acceptable. Portable modeled
store/application limits remain explicit. Neither option qualifies a release or
proves independent external interoperability.

## 19. Estimated maintenance reduction

Most support complexity pays for actual acceptance/independence. A saves commands
and removes confusing stale documentation with little LOC change. B could save
roughly 150–200 LOC and the need to mirror benign Make edits, at the cost above.
Factoring shell tool lists saves repetition, not dependency closure size. No
credible multi-thousand-line reduction exists without deleting useful evidence.
Prefer one bounded A pass; do not launch another broad support-tooling wave.

## 20. Implementation qualification requirements

This audit has no implementation qualification dependency. For a future A/B pass,
start again on clean main and record HEAD/revision/classifications. A changes
BUILD_NIX_CI and CONFORMANCE_GO (wiring policy tests) plus documentation. B also
changes qualification logic. These are not docs-only once implemented.

Run one baseline make check. Use focused repopolicy tests while adjusting wiring;
preserve mutation evidence for prohibited coupling and weakened full acceptance.
Inspect final orchestration to confirm generation remains deliberate. Once inputs
stabilize, run complete make check, both module vets, all-package make race and
make fuzz once, formatting and git diff --check. Preserve all three existing
fuzz targets with `-run '^$' -fuzztime 10s -parallel 2`; no race/fuzz reduction.
Record per-package/case counts, Python zero, fuzz durations and all attempts.

Compare spec/cases/schemas/manifest/profile bytes and inventories to baseline.
If documentation in vectors is corrected, record that sole prose delta separately;
do not repin or regenerate cases to accommodate it. Generator behavior changes
(including optional source-reader hardening) need separate temporary full
generation/compare and normal corpus verification; do not execute generator merely
because unrelated CI lines change. The optional hash-test deletion should retain
applicable temporary writer comparison evidence as stated above.

Freeze intended files, use intent-to-add for new membership only, confirm no
staged content, and request human `nix flake check` last. Agent never runs Nix.
Preserve Linux authority and all five remote portability combinations after an
owner-authorized commit/push. Do not weaken expected outcomes or claim pending CI
passed. Report-only edits outside source inputs do not require repeating suites.

## 21. Remaining owner decisions

No permission is needed to complete this design report. Future implementation
choices: authorize bounded A, decide whether B's loss of automatic Make/pin
tripwires is acceptable, and optionally factor development utility lists. Writer
strict bootstrap ingestion and schema/prose compatibility changes should remain
separate, explicitly scoped decisions. No release, tag, push or module merge is
proposed for immediate execution.

## 22. Git and source-invariance evidence

This report is excluded by the actual flake fileset: it is outside spec/vectors/
requirements, is not Go/native/module source, Makefile/go.work, or the explicitly
included `review/V1_PRE_R16_REVISION_HISTORY.md`. Current specpolicy reads only
spec and that archive; repopolicy does not read Markdown report contents;
schema/consumer/oracle read artifact inputs. No Go/Nix run is required for its
addition. Historical reports and qualification/build/filter inputs remain unchanged.

Verification uses starting Git tree membership plus content comparisons, including
additions/deletions, rather than mtimes. All original tracked files are compared
against starting HEAD; protected spec/vectors/requirements/conformance/generator,
workspace/all module/sum files, Makefile/flake/lock/workflows/tools inventory are
also checked explicitly. Absent root module/sum/workspace-sum and generator sum
remain absent. Final results are recorded below after Git-aware membership setup.

Execution accounting for this audit: baseline make check **0**; focused Go/Python
test invocations **0/0**; Python test executions **0**; focused fuzz
Dispatch/Open/ArrivalAndDisappearance **0/0/0**; final complete suites **0**;
incomplete final attempts **0**; explicit vet/race/fuzz **0/0/0**; generation **0**;
agent Nix **0**, human Nix **0 (not required for this excluded report)**;
remote CI **not queried/executed**. No cached result is claimed as current audit
execution. One optional read-only inventory attempt found python3 unavailable;
counts were collected with jq/rg/wc instead. No test/generation command failed.

No source edit, Go dependency update, vector regeneration, expectation change,
protocol update, Nix invocation, permanent tool, content staging, commit, tag,
push, release or delegation occurred. The report remains uncommitted.

Final source/membership checks: **PASS**. `git ls-tree -r HEAD` enumerated all
288 original paths; `git hash-object` for each working file matched its starting
Git blob identity. `git diff --quiet HEAD -- spec vectors requirements conformance
generator tools go.work Makefile flake.nix flake.lock .github .envrc` exited 0.
The whole-tree changed-path list contains exactly this report, so there are no
original path deletions or other intended additions. Nonignored untracked file
inventory is empty after intent-to-add; ignored local state is not added.
The five absent paths listed above remain absent.

```text
git branch --show-current: main
git rev-parse HEAD: d5b26b315b44f37df8cd793c1d458f3f8602e2af
git status --short:
 A review/SPEC_TOOLING_SIMPLIFICATION_AUDIT.md
git diff --name-only HEAD:
review/SPEC_TOOLING_SIMPLIFICATION_AUDIT.md
git diff --cached --quiet: exit 0
git diff --check: exit 0
```

The leading blank in ` A` indicates an unstaged new file with intent-to-add,
not staged report contents. Membership remains in place for owner review.
No human gate is pending for this audit. Work stops with this single report.
