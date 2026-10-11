# Totipo Spec support/tooling architecture audit

## 1. Starting revision and state

Audit date: 2026-10-11. Classification: **DOCS_ONLY — architecture / support-surface audit**. No semantic or executable qualification input changes.

The required first commands returned:

```text
git branch --show-current
main
git rev-parse HEAD
ad4e9000b6cdd4c7eebbaeeac06178750984f836
git status --short
(empty)
```

Clean `main` was established before inspection. Committed prerequisites are `b77f4f2` (qualification ladder), `31b79c9` (Nix/CI unification), and `ad4e900` (GitHub Actions refresh). Protocol: **v1/r19**. Requirements: `totipo-requirements-v1`, `moving-pre-rc`, `spec_revision: r19`, in [v1-pre-rc.json](../requirements/v1-pre-rc.json). Manifest: **92 cases**, all required; 92 physical case JSON files. This is a moving draft, not an RC or release.

No pruning is performed. The minimum future surface is Go for corpus validation, conformance and deliberate generation; Nix for fixed dependencies and Linux qualification; YAML for CI; a small Make interface; and explicit revision/qualification documentation. Python currently supplies some distinct evidence that must survive any consolidation.

## 2. Methodology

Read tracked tooling, both Go command implementations, all seven Python files, all seven internal Go packages and their tests, Make, flake/workspace/module inputs, workflow and Dependabot configuration, schemas, manifest/profile, maintenance docs and relevant historical reports. Used `rg` for callers and references, Git history for origins, and raw line counts for size. No tool under audit, generator, test suite, Nix, or remote workflow was executed. No external research was needed: recommendations concern this committed repository, not current upstream product versions.

Caller claims distinguish executable callers from documentation and historical invocations. Searches covered active source, build files, CI and docs plus historical evidence; absence means no in-repository caller found, not proof of no external users. Origin classifications are evidence-based where possible; a hardening revision alone does not establish a real failure. Historical reports are evidence, never prune targets.

LOC means physical lines including comments/blanks, not complexity or executable statements. Tables group a package's source and tests where their disposition is identical. All estimates describe the starting HEAD. The report is outside the inspected Nix source filter and outside structure-check dependencies.

## 3. Language inventory

| Language / representation | Files | LOC approx | Entry points / tests | Contributor relevance |
| --- | ---: | ---: | --- | --- |
| Go | 20 | 3,044 | Two commands; seven internal packages; eight test files, 29 Test functions and three Fuzz functions | Protocol implementation, generation and semantic evaluation |
| Python | 7 | 707 | Two active check scripts; three unittest modules with 14 test methods; two historical audits | Prose anchors, schema contract, independent selected wire checks |
| Make | 1 | 35 | Eight named targets | Human and Linux qualification orchestration |
| Nix | 1 | 135 | `checks.spec` on Linux, dev shell, formatter | Fixed dependency closure, toolchains, full Linux gate |
| YAML | 2 | 73 | One workflow (63); Dependabot (10) | Qualification, portability, update suggestions |
| Shell | 0 standalone `.sh`; 1 `.envrc` | 2 standalone environment lines | `use flake`, `layout go`; embedded Make/Nix/CI shell | Environment setup, cache/path handling, formatting test |
| JSON tooling contracts | 4 selected | 4,241 | Case schema 2,475; manifest schema 80; manifest 1,306; requirements 380 | Language-neutral corpus contract, inventory and pins |
| Other dependency metadata | 4 | 293 | flake.lock 279; go.work 3; go.mod 7; go.sum 4 | Reproducible dependency/workspace boundary |
| Markdown process docs | 3 selected | 306 | AGENTS 153; CONTRIBUTING 62; README 91 | Revision and qualification instructions |

Go: 2,138 non-test lines and 906 test lines. Python: 214 historical-audit lines, 232 active checker lines, 261 test lines. Active Python qualification is therefore **493 lines in five files**; only 103 of those lines are the independent wire test module. The 3,044 Go lines include the 508-line generator and 41-line consumer CLI. JSON corpus contents, spec prose, design assets and the 30 existing review files are not implementation LOC totals. The four selected JSON files are tooling inputs, not candidates to discard because they dominate line count.

Running ordinary qualification requires both Go and Python installations; editing every component does not require fluency in both. A semantic change spanning wire/corpus contracts commonly requires Go, Python, JSON and normative Markdown. Spec-only clarification commonly encounters Python anchors and Go-driven pin maintenance even when Go semantics do not change. That coupling is the principal maintenance friction.

## 4. Executable entry points

| Entry point | Actual callers | Result / mutation boundary |
| --- | --- | --- |
| `make` / `make help` | Human; Make default | Prints seven advertised commands; advertises seven checking commands; no obsolete aliases |
| `make check` | Human docs; Nix build | Serial `spec-check`, `test`, `conformance`, `verify` |
| `make spec-check` | `check`; human help/docs | `python3 tools/check_spec.py` |
| `make test` | `check`; human help | Python discovery, then `go -C conformance test -count=1 ./...` |
| `make verify` | `check`; human generator-maintenance docs | Python schemas, then Go CLI `-verify-only` |
| `make conformance` | `check`; human docs | Go CLI full manifest execution |
| `make race` | Human docs; Nix | All Go packages, `-race -count=1 ./...` |
| `make fuzz` | Human docs; Nix | Three bounded Go targets; exact policy in section 12 |
| `go run ./conformance/cmd/totipo-conformance -root .` | Make; portability CI; human | Read + profile check + expected-result execution; no case-selection flag |
| Same command with `-verify-only` | Make; portability CI; human | Read + profile check, no expected-result execution |
| `go run ./conformance/cmd/generate-vectors -root .` | Human, documented in conformance README / vectors FORMAT | Explicit writes to cases, manifest, profile; removes undeclared physical cases; never normal qualification |
| `python3 tools/check_spec.py` | Make; portability CI; test_check_spec subprocess | Reads spec/archive/profile/artifacts; asserts anchors and pins |
| `python3 tools/check_vectors.py` | Make; portability CI | Validates manifest and listed case JSON against schema subset |
| `python3 -m unittest discover -s tools -p 'test_*.py'` | Make; portability CI | 14 declared tests; no generation |
| Focused `python3 -m unittest …`, `go -C conformance test …` | Human/agent ladder | Narrow diagnostics; not independent new authorities |
| `go -C conformance vet ./...` | Human docs; Nix; portability CI uses working directory | Static Go analysis |
| `test -z "$(gofmt -l conformance)"` | Nix; portability CI (`shell: bash`); human docs | Formatting; no rewriting |
| `nix flake check path:.` | Human docs; Linux CI adds `--print-build-logs` | Linux pinned full qualification; agent forbidden to run |
| `go work vendor` | Nix dependency derivation | Produces fixed Go dependency material, not corpus evidence |
| `nix develop`, direnv `.envrc` | Human environment setup | Development tools, not qualification evidence |
| Flake formatter (`nix fmt` interface) | Exposed for human convenience; no CI caller found | `nixpkgs-fmt`; not a protocol check |
| `python3 tools/audit_r15.py`, `python3 tools/audit_r16.py` | Historical review instructions only | Git-baseline/worktree migration evidence; stdout JSON, not live checks |

No release, tag, version-bump, revision-update, publishing or packaging command was found. Release/revision discipline is currently human review plus moving pins, schema/profile revision literals, history and qualification policy. No frozen-release mechanism is implied by these tools.

## 5. Call graph

```text
human -> Make check -> check_spec -> spec + archived history + four profile pins
                    -> test -> Python discovery -> test_check_vectors -> check_vectors helpers
                                               -> test_check_spec -> isolated check_spec CLI
                                               -> test_r16_vectors -> check_vectors.read + independent formulas
                            -> Go ./... -> vectors.TestCorpus -> Read + VerifyProfile + Run
                                        -> package regressions / fuzz seeds
                    -> conformance -> CLI -> Read + VerifyProfile + Run
                    -> verify -> check_vectors -> JSON schemas + manifest-listed cases
                              -> CLI -verify-only -> Read + VerifyProfile
human -> explicit generator -> Go codecs/types + source known answers -> cases/manifest/profile
human -> race/fuzz/vet/formatting
Linux CI -> nix flake check -> checks.spec -> fixed goDependencies + environment
                                       -> formatting + Make check + vet + Make race + Make fuzz
portability CI -> direct Python + formatting + Go vet/test + both consumer CLI modes
historical human instructions -> audit_r15/r16 -> git show fixed BASE + historical worktree
Dependabot -> weekly Go-module / Actions update suggestions
```

Package edges: vectors calls cryptov1, object, graph, storage, totp; graph calls object; storage calls cryptov1 and object; cryptov1 calls object and x/crypto; object calls tlv; totp uses stdlib. The generator imports vectors for contracts/Decode/Hash and graph/object/tlv/cryptov1, but does not call `vectors.Run` or query `graph.Evaluate` to obtain expected graph results. Test-only closure and integration oracles remain separate. No live executable calls the historical audit scripts. Nix includes their bytes through its broad tools `.py` filter nonetheless.

## 6. Component disposition

Disposition vocabulary is restricted to KEEP, KEEP_BUT_SIMPLIFY, MERGE_CANDIDATE, DELETE_CANDIDATE, HISTORICAL_ONLY, MIGRATION_CANDIDATE, UNKNOWN. Origin vocabulary is DURABLE_DESIGN_REQUIREMENT, REGRESSION_FROM_REAL_FAILURE, TEMPORARY_DIAGNOSTIC_THAT_BECAME_PERMANENT, LEGACY_FROM_OLDER_WORKFLOW, CONVENIENCE, UNKNOWN. Confidence concerns the recommendation, not a new runtime pass.

| Component | Language | LOC approx | Callers | Responsibility | Unique invariant | Origin | Disposition | Confidence |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| internal/tlv, incl. tests | Go | 65 | object, generator, package tests | Fixed framing | Exact big-endian tag/length, bounded truncation handling | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| internal/object, incl. tests | Go | 273 | crypto capacity, graph, storage, vectors, generator, tests | TOKEN grammar/value projection | Canonical fields/order/widths; complete value; metadata distinction; Dispatch fuzz | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| internal/cryptov1, incl. tests | Go | 264 | storage, vectors, generator, tests | Envelope/bootstrap crypto | AEAD then length/padding/keyed-ID; exact VAULT; Open fuzz | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| internal/graph, incl. tests | Go | 410 | vectors, generator types, tests | Observed graph and folds | SCC heads/missing edges; exact snapshots; bounded arrival/disappearance fuzz | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| internal/storage, incl. tests | Go | 316 | vectors, tests | Modeled namespaces/publication/create | No overwrite, exact retry, create-only VAULT, trusted persistence outcomes | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| internal/totp, incl. tests | Go | 83 | vectors, integration/package tests | RFC computation | Known-answer codes, unsigned time and credential bounds | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| vectors/types.go | Go | 148 | runner, generator, tests | Shared corpus types | Stable consumer-facing contract, optional presence | DURABLE_DESIGN_REQUIREMENT | KEEP_BUT_SIMPLIFY | MEDIUM |
| vectors/runner.go | Go | 594 | CLI, tests; generator Decode/Hash | Integrity/shape/evaluation | Exact physical/profile accounting and expected outcomes | DURABLE_DESIGN_REQUIREMENT | MERGE_CANDIDATE | HIGH |
| vectors/post_aead.go | Go | 90 | Run, runner tests | Authenticated-negative fixture proof | Defect must actually pass AEAD and reach named boundary | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| vectors/*_test.go (three files) | Go | 375 | Go tests/race | Corpus, strict input, integration regressions | All cases + metadata/fold preservation + real negative boundary | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| cmd/totipo-conformance | Go | 41 | Make, portability CI, human | Read/verify/execute CLI | Failing case returns failure; every required case executed | DURABLE_DESIGN_REQUIREMENT | KEEP_BUT_SIMPLIFY | HIGH |
| cmd/generate-vectors | Go | 508 | Human only | Deterministic fixture/pin maintenance | Reviewed bytes and explicit expectations; preserves RFC inputs | DURABLE_DESIGN_REQUIREMENT | KEEP_BUT_SIMPLIFY | HIGH |
| tools/check_spec.py | Python | 131 | Make, CI, subprocess tests | Prose/history anchors and artifact pins | Live structure plus stale-byte-pin detection | REGRESSION_FROM_REAL_FAILURE | MERGE_CANDIDATE | HIGH |
| tools/check_vectors.py | Python | 101 | Make, CI, Python tests/helpers | Schema-subset interpreter / strict JSON | Required fields/types including false/zero presence; schema contract | DURABLE_DESIGN_REQUIREMENT | MIGRATION_CANDIDATE | HIGH |
| tools/test_check_spec.py | Python | 69 | unittest discovery | Isolated CLI pin regressions | Actual byte edits, missing/stale pins fail | REGRESSION_FROM_REAL_FAILURE | MERGE_CANDIDATE | HIGH |
| tools/test_check_vectors.py | Python | 89 | unittest discovery | Schema negative regressions | Required false result, exact u64, nested unknowns, no rewrap | DURABLE_DESIGN_REQUIREMENT | MIGRATION_CANDIDATE | HIGH |
| tools/test_r16_vectors.py | Python | 103 | unittest discovery | Independent wire/intermediate checks | Independent HMAC/HKDF/TLV/VAULT-ID and defect-input evidence | DURABLE_DESIGN_REQUIREMENT | MIGRATION_CANDIDATE | MEDIUM |
| tools/audit_r15.py | Python | 125 | Historical report command | Fixed r14→r15 migration audit | Historical evidence only; no r19 invariant | TEMPORARY_DIAGNOSTIC_THAT_BECAME_PERMANENT | DELETE_CANDIDATE | HIGH |
| tools/audit_r16.py | Python | 89 | Historical report command/link | Fixed r15→r16 migration mapping | Historical evidence only; no r19 invariant | TEMPORARY_DIAGNOSTIC_THAT_BECAME_PERMANENT | DELETE_CANDIDATE | HIGH |
| Makefile | Make/shell | 35 | Human, Nix | Portable command policy | Full check, all-package race, exact bounded fuzz coverage | CONVENIENCE | KEEP_BUT_SIMPLIFY | HIGH |
| flake.nix qualification/dependencies | Nix/shell | Part of 135 | Linux CI, human | Fixed vendor input + hermetic full suite | vendorHash closure, offline Go, pinned Linux gate | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| flake.nix dev tools/jailed agents | Nix | Part of 135 | Human shell | Development composition | No protocol invariant; optional authoring convenience | CONVENIENCE | KEEP_BUT_SIMPLIFY | MEDIUM |
| flake.lock | JSON/Nix | 279 | Nix | Input revision/content pins | Fixed environment/dependency graph inputs | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| go.work | Go metadata | 3 | Root Go commands, Nix vendoring | Workspace membership | Root entry points resolve one conformance module | CONVENIENCE | KEEP_BUT_SIMPLIFY | MEDIUM |
| conformance/go.mod, go.sum | Go metadata | 11 | Go, Nix, Dependabot/cache | Dependency declarations/checksums | Exact module dependency boundary | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| workflows/conformance.yml | YAML/shell | 63 | GitHub push/PR/manual dispatch | Linux authority + portability | Pinned actions, full Linux gate, five supplemental combinations | DURABLE_DESIGN_REQUIREMENT | KEEP_BUT_SIMPLIFY | HIGH |
| .github/dependabot.yml | YAML | 10 | GitHub service | Weekly Go/Actions proposals | No semantic invariant; dependency upkeep convenience | CONVENIENCE | KEEP_BUT_SIMPLIFY | MEDIUM |
| .envrc | Shell/direnv | 2 | Human direnv | Enter development environment | No independent qualification invariant | CONVENIENCE | KEEP_BUT_SIMPLIFY | MEDIUM |
| vectors/{case,manifest}.schema.json | JSON | 2,555 | Python checker, tests; pinned by Go/profile | Language-neutral contract | Exact presence/type/operation shapes | DURABLE_DESIGN_REQUIREMENT | KEEP_BUT_SIMPLIFY | HIGH |
| vectors/manifest.json | JSON | 1,306 | Go/Python, generator, tests | Enumerate evidence | Each physical case once; exact ID/hash/expected result | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| requirements/v1-pre-rc.json | JSON | 380 | Go, check_spec, generator/tests | Moving requirements pins | All-and-only cases and exact four artifact bytes | DURABLE_DESIGN_REQUIREMENT | KEEP | HIGH |
| AGENTS/CONTRIBUTING/README instructions | Markdown | 306 | Human/agent | Revision/review/qualification discipline | Avoid skipped final coverage and generation-as-verification | DURABLE_DESIGN_REQUIREMENT | KEEP_BUT_SIMPLIFY | HIGH |
| vectors/FORMAT, conformance/README, requirements/README | Markdown | Not in implementation totals | Human | Contract/model limits and maintenance | Explicit semantic expectations and modeling limits | DURABLE_DESIGN_REQUIREMENT | KEEP_BUT_SIMPLIFY | HIGH |
| review/V1_PRE_R16_REVISION_HISTORY.md | Markdown | Historical | check_spec; included by Nix | Archived revision entries | Historical revision continuity is currently a live guard | LEGACY_FROM_OLDER_WORKFLOW | HISTORICAL_ONLY | HIGH |
| Other existing review reports / audit JSON / status text | Markdown/JSON/text | Historical | Human evidence, links | Decisions and execution records | Preserve evidence; no live qualification authority | LEGACY_FROM_OLDER_WORKFLOW | HISTORICAL_ONLY | HIGH |

The archived-history row is HISTORICAL_ONLY as content/disposition; its current executable dependency is explicitly real. Do not remove it or silently drop that dependency. Core packages are small and have distinct invariants; merging them to reduce directory count would give little benefit.

## 7. Python inventory

All `.py` files are in `tools/`; all imports are standard library or local `check_vectors`. No pip dependencies, virtualenv, requirements.txt or Python packaging machinery is required.

| File | Responsibility and callers | Go overlap | Tests and semantic importance | Still needed? Could live in Go? |
| --- | --- | --- | --- | --- |
| check_spec.py | Top-level assertions, reads relative to working directory; Make/CI/isolated test subprocess | Four artifact pins exactly overlap VerifyProfile; revision literals overlap schema/Go/generator | Three CLI tests cover pins, not most prose anchors; guards policy/spec scope as well as syntax | Structure/revision/byte discipline yes; many exact prose fingerprints need review first. Merge pins; migrate reduced structural contract to Go |
| check_vectors.py | Duplicate-member reader, schema subset, main validates listed paths; Make/CI; imports from both vector test modules | Go strict Decode/ValidateShape cover some shape, not all required presence; Go alone covers exact physical/profile accounting | Eight schema tests; no dedicated end-to-end CLI test; semantically important omission/type checks | Yes until complete equivalent contract validation exists in Go; use Go schema validation, not only current structs |
| test_check_spec.py | Temp copies, raw-byte perturbations, stale/missing fields, current pass; discovery | Go tests currently lack equivalently exhaustive four-artifact pin mutations | Three methods, many subtests; regression from actual r18 stale spec pin | Keep behaviors as consolidated verifier tests; CLI harness can disappear with old entry point |
| test_check_vectors.py | Mutates corpus copies/schema inputs; discovery; imports read/validate/pairs | Strict decoding and shape checks overlap; Go bool zero values do not prove required false field presence | Eight methods; durable negative contract cases | Port cases after merged authority, discard Python-specific helper coupling |
| test_r16_vectors.py | Corpus wire/key intermediates, VAULT_ID/max, post-AEAD input construction; discovery | Recomputes selected values also checked by Go consumer/generator | Three methods; independent formulas; no AES/Argon2 independence claim | Responsibility still needed despite r16 filename. Can live in Go tests only if separately implemented oracle remains independent of production/generator helpers |
| audit_r15.py | Hard-coded BASE `0b63f778…`, compares worktree and Git bytes; historical r15 report | Hashing/JSON inventory only; no live responsibility | No tests/callers in qualification; frozen-source/history assertions describe superseded migration | No live r19 need; do not port. Retrieve at historical commit for reproduction |
| audit_r16.py | Hard-coded BASE `0a01fc1f…`, replacement mapping and audit JSON; historical r16 report | Generic hash/JSON reading only | No tests; depends on historical corpus shape and current worktree | No live r19 need; do not port. Preserve committed audit evidence and historical source retrieval |

Python removal is feasible, but removing its 103-line independent oracle by calling Go production helpers would lose evidence. Language diversity itself is not the essential contract; independent computation is. A separately implemented Go oracle is possible, with review of correlated-error risk and explicit preservation of exact constants/formulas. It still would not provide independent AES/Argon2 interoperability evidence.

## 8. Go command inventory

Exactly two `package main` commands exist. `totipo-conformance` accepts `-root` and `-verify-only`, rejects positional arguments, exits nonzero on load/profile or case failure, and prints r19 summaries. There is no single-case CLI selection. Focused cases can be exercised through `TestCorpus` subtests; no invented manifest selection interface should be documented.

`generate-vectors` accepts `-root`, owns public fixed fixture inputs, writes indented JSON with newline, sorts cases, preserves source RFC TOTP bytes, generates malformed authenticated fixtures deliberately, removes undeclared case paths and updates profile pins. It has no direct `_test.go` file or normal regenerate-and-compare check. Go test lists both command packages as without tests; compilation is not generator behavioral qualification. Current generator accepts flag parsing without the consumer's explicit positional-argument rejection.

Do not combine these commands. Verification-only and full consumer modes can share one read/profile pipeline, as they already do. Their separate invocations inside `make check` duplicate integrity work, but provide independently callable interfaces and CLI smoke evidence. Evaluate reducing repeated invocation only with replacement CLI coverage and unchanged qualification contracts, in a later approved change.

## 9. Generator/verifier boundary

| Function | Current authority | Intentional independence / limits |
| --- | --- | --- |
| Generation | generate-vectors; source known answers | Explicit maintenance, never called by tests/check/Nix/CI |
| Verification of bytes/intermediates | Go vectors runner + Python wire tests | Go shares crypto with generator; Python independently computes selected intermediates only |
| Manifest validation | Go Read; Python schema validator | Go checks identity/hash/path/physical accounting; Python validates schema presence/types |
| Structure validation | check_spec; schemas/check_vectors | Prose/history guards distinct from execution; do not infer semantics from prose substring success |
| Expected-result evaluation | vectors.Run + model packages | Graph/workflow expectations authored explicitly; not computed by evaluator in generation |
| Fixture authenticity/defect verification | runPostAEAD + tests | Raw AES-GCM open proves negative fixture reaches intended post-authentication rejection |
| Requirements verification | Go VerifyProfile; repeated artifact hashes in check_spec | Checks actual raw bytes, not generator self-certification |

Generator and consumer share contract types, hashes and production crypto. This is reproducibility plus consumer evidence, not two independent cryptographic implementations. RFC TOTP answers are external fixed inputs. Bootstrap source known answers are retained, including Unicode fixture. Graph expected heads/values are authored; generator attaches declared head records from input facts. Fold stage expectations are authored separately from Fold evaluation. Independent graph matrix closure tests cover all 512 directed three-node graphs. None of these should be replaced with generator-derived expected outputs.

## 10. Make audit

| Target | Logic / wrapper classification | Consumers | Future judgment |
| --- | --- | --- | --- |
| help | Human convenience; printf command catalogue | Default Make goal, human | KEEP_BUT_SIMPLIFY; derive/check advertised revision only if worthwhile |
| spec-check | Python wrapper | check, human | MERGE_CANDIDATE after Go structural/check authority exists |
| test | Orchestration of Python + Go | check, human | KEEP_BUT_SIMPLIFY; eventual Go-only test orchestration |
| conformance | Go CLI wrapper | check, human | KEEP; invariant: execute all required expected outcomes |
| verify | Python + Go wrappers | check, maintenance docs/human | KEEP_BUT_SIMPLIFY; one merged read-only verifier command |
| check | Dependency orchestration; `.NOTPARALLEL` | Nix, human | KEEP; invariant: complete portable normal qualification |
| race | Go wrapper with all-package flags | Nix, human | KEEP; invariant: all-package race coverage |
| fuzz | Policy-bearing orchestration of exact three targets | Nix, human | KEEP; invariant: seeds/failures and bounded 10s/parallel-2 policy |

There are **eight** targets including help (seven checking/orchestration targets). No obsolete alias was found. `GO`/`PYTHON` override and writable `GOCACHE` export are environment ergonomics, not protocol semantics. Make earns its 35 lines by giving stable short commands and central race/fuzz flags; deleting it would relocate orchestration into Nix/CI/docs. Do not implement schema/protocol logic in Make to eliminate Python. Later consolidations should preserve target meanings unless explicitly reviewed as interface changes.

## 11. Nix audit

[flake.nix](../flake.nix) owns environment composition, fixed Go dependencies and Linux qualification, not protocol semantics. `buildGoModule.goModules` plus custom `go work vendor` fixes the workspace dependency boundary with vendorHash `sha256-637XdgBok7hRhMw5s9hdY3Mr5jlPXqc1viSYoYh9Iqw=`. The qualification derivation copies that dependency tree and uses `GOFLAGS=-mod=vendor`, `GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local`, `GOENV=off`, fresh temporary caches, strict dependencies and CGO/compiler support for race. No Nix execution was used to validate these static observations.

`qualificationSource` includes spec/, vectors/, requirements/, conformance Go/module files, tools `.py`, Makefile, go.work and the pre-r16 archive. README, AGENTS, CONTRIBUTING and this report are excluded. Historical audit Python files and maintenance-only generator Go are included even though no qualification command executes them; this broadens invalidation, not semantic coverage. Narrowing that filter is a future build change requiring source-boundary review and human Nix qualification, not a docs-only cleanup.

Nix runs formatting, Make check, explicit vet, Make race/fuzz. Protocol checks remain tools, not Nix expressions. Repeated r19 labels in pname versions/output are descriptive metadata, not a revision authority. No package/default app or artificial build target is needed. Linux-only checks and Darwin development shells are intentional.

Dev shells include gopls, golangci-lint, its language server, and jailed Codex/Pi packages in addition to Go/Make/Python; jail configurations repeat extra package lists and forward GOPATH/GOBIN. These optional editing tools are not qualification inputs/tools used by the derivation. Future simplification may share the package list or split optional agent tooling from the minimal contributor shell. Do not remove their lock closure as supposedly dead without owner confirmation of usage.

## 12. Qualification overlap matrix

| Layer | Exactly what it establishes when passing | Overlap / distinct evidence |
| --- | --- | --- |
| make check | Prose/history/pins; 14 Python methods; seven Go test packages; all 92 cases; schema and Go integrity CLI | TestCorpus and full CLI both evaluate corpus; full CLI and verify-only both load/pin-check; Python repeats four hashes |
| go vet | Static Go diagnostics across `./...` | Distinct from execution; not proof of runtime semantics |
| make race | `go -C conformance test -race -count=1 ./...` | Reexecutes package tests with runtime race instrumentation; not equivalent to ordinary tests or proof all races absent |
| make fuzz | Object FuzzDispatch, crypto FuzzOpen, graph FuzzArrivalAndDisappearance | Seeded mutation/property exploration beyond deterministic tests; each `-run '^$' -fuzztime 10s -parallel 2` |
| nix flake check | Linux formatting + complete check/vet/race/fuzz in pinned dependencies/environment | Reuses host suite intentionally; adds reproducible dependency/tool boundary, not new protocol semantics |
| portability jobs | Python checks/tests, formatting, vet, Go tests, both CLI modes across five OS/Go combinations | Deliberate alternate host/toolchain evidence; no race/fuzz here |

Fuzz seed counts in source: Dispatch 2, Open 2, graph 1. Preserve seed contents, failure propagation and reproduction inputs. Ordinary Go tests also run seed cases; they do not replace bounded fuzz exploration. Fuzz properties: canonical Dispatch roundtrip/no invalid object leakage; canonical Open/Seal identity; graph arrival-order/disappearance invariance. No instruction here changes flags, Make targets, CI matrix or final coverage.

Portability matrix expands to macOS × {1.23.x, stable}, Windows × {1.23.x, stable}, plus Ubuntu 24.04 × 1.23.x. Authoritative Linux qualification is the separate pinned Nix job. All-package race covers the seven test packages; the two commands have no test bodies. Formatting is absent from `make check` itself and supplied explicitly by Nix/CI/docs.

## 13. Structural/revision guard audit

check_spec currently guards 21 consecutive section numbers, protocol/revision literals, pre-r16 archive link and r1–r19 entry presence, retired-term exclusions in active prose, registry tags/names, numeric section references, storage layout, metadata/value separation, fold preservation, storage entry types, application/core scope, create-only VAULT, and four raw artifact pins. These are live rules, but exact sentences/line breaks and section-specific regex splits are implementation choices rather than required architecture.

The pre-r16 archive check protects continuity of recorded entries, not the correctness of historical protocol semantics. Historical report hashes/output should remain immutable evidence. The archive is deliberately read and filtered into Nix; other reports are not. A simpler structural validator could validate section uniqueness/reference resolution, registry data, revision/history coverage and hash pins, while selected policy anchors remain focused tests or review requirements. Retired-term rejection is not proof of obsolete semantics being absent and can reject harmless explanatory prose; review each rule before porting. Never silently discard a guard just because it uses a substring.

Revision is encoded as literals across spec, Go Read/VerifyProfile, consumer output, generator, Python checker, manifest schema, manifest/profile, Make help, Nix labels and docs. There is no shared revision parser. Machine consistency can be consolidated around explicit validated metadata without deriving expected requirements pins from the current artifacts during verification. Historical r16 titles/comments in case schema, object package and Python test name are stale labels, not unused semantic implementations. TLV's comment refers to section 41, another documentation cleanup candidate.

The pin block and its three tests have a documented real origin: [r18 integrity repair](V1_R18_INTEGRITY_REPAIR_REPORT.md), commit `4623a7e`, found a spec pin left stale after final editorial wording. check_spec previously passed; Go VerifyProfile already provided full integrity checking. Preserve the regression behavior when merging, rather than retaining a second hashing authority simply because it fixed that incident locally.

## 14. Corpus/vector tooling audit

The manifest is the explicit portable inventory; the requirements profile requires exactly that inventory. Cases contain expected outcomes; schemas define the interchange shape; specification defines normative semantics. Go consumer is conformance evidence, not an authority allowed to redefine expectations.

Category accounting: bootstrap 5, crypto 5, encoding 30, fold 6, graph 13, metadata 7, size 1, storage 13, totp 3, vault 9 = **92**. All are portable required cases; no capability/conditional selector remains. Exact accounting needs unique IDs and paths, no unlisted physical files, matching case identity/outcome, per-file raw hash, exact profile membership/hash and four artifact pins. Go Read/VerifyProfile uniquely supply this complete accounting today. Python main only enumerates manifest-listed cases; it does not prove physical completeness or required-set equality.

Python supplies schema presence/types/enums/ranges, including required false-valued booleans and exact integer limits. Go `DisallowUnknownFields`, uniqueJSON, typed decoding and ValidateShape reject many malformed inputs but cannot distinguish omitted bools from explicit false in ordinary structs. Therefore deleting the schema validator because Go can unmarshal is unsafe. Go pins schema bytes but does not interpret JSON Schema. The Python interpreter is intentionally a checked-in subset, not a general draft-2020-12 validator. Its unknown-keyword rejection applies when schema nodes are visited; it is not a full independent schema meta-validation pass.

Schema repetition in the 2,475-line case contract could later use reviewed data reuse, but changing schemas changes qualification inputs and pins. Keep the language-neutral contract and every negative case. Do not turn generator structs into the sole definition and sole verifier of accepted corpus shapes.

## 15. Tooling-test audit

| Tests | Durable evidence | Implementation-specific portion / missing coverage |
| --- | --- | --- |
| test_check_spec (3) | Raw byte drift including trailing newline; all four missing/stale pins | Temp copy/subprocess/current-directory harness can change; most prose anchors lack mutation regressions |
| test_check_vectors (8) | Required false, exact u64/metadata, nested unknowns/mixed payload, revision/retired extension, genuine defect shape, no replace/rewrap, duplicate JSON | Calls Python helpers directly; no dedicated checker-main path traversal/CLI test |
| test_r16_vectors (3) | Independently encoded TLV and HMAC/HKDF-derived intermediates; VAULT hash/max padding; post-AEAD construction inputs | Local reader import can change; r16 filename can change; no AES/Argon2 independent verification |
| vectors runner tests (5) | Corpus, strict decoding, manifest revision tampering, metadata-loss rejection, authenticated negative fixtures | Single revision tamper test is not exhaustive pin/path coverage; no direct consumer CLI test |
| vectors integration (1) | Real encrypted folds, metadata exact presence, retained facts and TOTP independence from status | Keep outcomes/oracle, not test layout |
| Core package tests (23) | Framing/grammar/crypto/graph/storage/TOTP invariants, snapshot independence, matrix closure | Package layout may change, semantic cases must remain |
| Three Go fuzz functions | Canonical parsers and observed-set order properties | Seed runs alone are insufficient; preserve bounded mutation execution |
| Generator command | Output is covered by normal verification and independent known answers | No generator unit tests or automatic temporary regenerate-and-compare target; compilation does not establish deterministic regeneration |

29 Go Test functions = 5 runner + 1 integration + 23 core; seven test packages, eight files. Test counts are static inventory and inherited reports, not fresh executions. Regression cases should move with responsibilities; Python helper/harness tests can disappear only when replaced by equivalent new-interface evidence. Do not add tests mirroring a literal port instead of testing invariants.

## 16. Duplication matrix

| Responsibility | Implementations | Classification | Consolidation boundary |
| --- | --- | --- | --- |
| Revision/header parsing | Python literals, Go literals, schema const, generator literals | SHARED_LIBRARY_CANDIDATE | One validated metadata contract for tools, separate expected-version/pin assertions |
| Manifest parsing | Go typed Read/Decode; Python read/schema; generator source RFC reads | SHARED_LIBRARY_CANDIDATE | Go strict reader plus schema contract; do not weaken external interchange |
| JSON duplicate/unknown/shape rejection | Python pairs/validate; Go uniqueJSON/Decode/ValidateShape | SHARED_LIBRARY_CANDIDATE | Union of protections, especially required false/zero and nested fields |
| Case enumeration | Go Read/WalkDir, Python listed iteration, generator declared corpus | INTENTIONAL_INDEPENDENT_CHECK | Writer declaration and verifier physical accounting remain separate; share read-only inventory only where safe |
| Artifact checksum/profile pins | Python check_spec; Go VerifyProfile | ACCIDENTAL_DUPLICATION | One verifier implementation plus focused regression and independent audit when needed |
| Case checksum creation vs checking | Generator Hash/pins vs Read/VerifyProfile | INTENTIONAL_INDEPENDENT_CHECK | Same hash primitive is fine; verification reads actual bytes and external pins |
| Requirements exact set | Go VerifyProfile; generator writes full set | INTENTIONAL_INDEPENDENT_CHECK | Never let verification repair expected set |
| Spec structure | Python only; Go has version/header checks | UNKNOWN | Most prose semantics are not encoded twice; distinguish necessary anchors from review fingerprints |
| Wire/key intermediates | Generator/Go runner vs Python formulas | INTENTIONAL_INDEPENDENT_CHECK | Independent oracle must not reuse production helpers |
| Expected semantic results | Generator explicit declarations vs Go evaluators and closure test | INTENTIONAL_INDEPENDENT_CHECK | Do not query evaluator to author expectations |
| Corpus execution | TestCorpus + consumer CLI; race repeats tests | INTENTIONAL_INDEPENDENT_CHECK | Same oracle, different integration/instrumentation evidence; possible future CLI smoke refinement |
| Integrity CLI after full CLI | Full CLI + verify-only in Make/CI | ACCIDENTAL_DUPLICATION | Repeated work, not independent verification; retain standalone verify interface and current suite until reviewed |
| Formatting checks | Nix/portability YAML/docs same gofmt command | ACCIDENTAL_DUPLICATION | Central small formatting entry point if useful; no new semantic validator |
| Vet/test orchestration | Make/Nix + direct portability YAML | SHARED_LIBRARY_CANDIDATE | Common portable check command, if Windows invocation stays reliable |
| Dependency boundary | go.mod/sum, vendorHash, lock, offline flags | INTENTIONAL_INDEPENDENT_CHECK | Different supply-chain layers; not interchangeable hashes |
| Developer package lists | Default shell / two jailed shells | ACCIDENTAL_DUPLICATION | Share environment list or separate optional extras |

## 17. Intentional independence

Preserve writer vs reader; authored expected results vs evaluator; schema contract vs permissive language defaults; independent wire/key formulas vs production helpers; RFC known answers vs generated values; graph matrix closure vs traversal; authenticated-defect verification vs generic rejection; and normal/race/fuzz/toolchain-platform evidence. Independence can exist within one language, but sharing all algorithm code removes it. Conversely two invocations of the same Go Read/VerifyProfile are repetition, not independent implementations.

No claim of independent full crypto validation is warranted today. Generator and consumer share AES/Argon2/HKDF implementations; Python covers selected SHA/HMAC/TLV intermediates, not independent AEAD/KDF. Before RC freeze, the documented interoperability requirement remains outstanding, regardless of language reduction.

## 18. Dependency/tool surface

Required direct local tooling: Go 1.23+, Python 3.9+, GNU Make for its entry points, and a compiler/CGO-capable environment for race. Go external modules: x/crypto v0.36.0 (Argon2/HKDF) and indirect x/sys v0.31.0, with go.sum checksums. Go stdlib supplies AES-GCM, hashes, JSON, paths, flags and tests. Python is stdlib-only. Bash is explicitly used for portable formatting; Make/Nix contain ordinary shell orchestration.

Pinned Linux gate uses Nix/nixpkgs plus flake-utils and the fixed dependency derivation. llm-agents and jailed-agents are optional authoring integrations in the flake, not protocol requirements. gopls and golangci tools are editor conveniences; no golangci-lint qualification invocation was found. Git is needed for repository/history discipline and old audit scripts, not normal verification inside the filtered Nix source.

CI pins checkout, Nix installer, setup-go and setup-python by SHA, uses read-only contents permission and disables persisted checkout credentials. Setup-go caches by conformance/go.sum. Dependabot proposes Go-module and Actions updates weekly; it does not update Nix lock/vendorHash automatically. No pip, Node, npm, shell-script library, release action or artifact-packaging surface was found for this qualification scope.

## 19. Contributor cognitive surface

| Task | Must understand | Python/Go interaction / avoidable friction |
| --- | --- | --- |
| A. Spec clarification | Normative Markdown, revision/history discipline, JSON artifact pins, qualification commands | Exact Python prose anchors may fail after editorial change; Go generator currently refreshes pins together with corpus |
| B. Protocol semantic change | Spec, Go model/crypto, JSON corpus/schema/profile, Python checks/oracles, Make/Nix gates | Largest cross-language surface: repeated revision/header, schema/struct shapes, wire formulas and anchors |
| C. New conformance vector | JSON contract/expectations, Go generator or deliberate fixture source, consumer and exact requirements set | Python schema regression may be needed; Go defaults are not equivalent to schema presence checks |
| D. Generator change | Go deterministic inputs/codecs, source known answers, independent verifier, temporary regeneration comparison | Python independent evidence must be interpreted; generator success cannot establish correctness |
| E. Dependency/Go update | go.mod/sum, Go compatibility, go.work, Nix vendorHash/lock, portability matrix, race/fuzz | Python normally unchanged; installing it is still required for full gate |
| F. CI maintenance | YAML, pinned actions, Nix Linux authority, shell commands, Go/Python setup, Make mapping | Portable YAML repeats check command list and Python version floor; preserve actual five combinations |

Contributors can run gates without reading Nix internals; only dependency/environment changes require maintaining those details. Consolidating structural/hash/schema validation into Go reduces two-language maintenance for A–D. Moving optional authoring tooling out of the minimal environment reduces setup cognitive load without reducing protocol evidence.

## 20. Dead/unreferenced tooling

No active Make target, Go command or internal package is confidently dead. Every core package has semantic consumers/tests. Generator has explicit human maintenance callers even though CI never invokes it. `-verify-only` has real Make/CI/doc callers. test_r16_vectors is actively discovered despite its old name. No superseding generator or unused standalone shell script was found.

The two audit scripts have **historical callers only**, not zero references. Their fixed old BASE commits, retired-field assumptions and worktree comparisons mean they are not general r19 auditors. r15 checks old applicability/history/wire invariance; r16 maps removed concepts to replacement IDs. They should not be invoked on r19 as current gates. Their outputs and source at historical commits preserve their purpose.

Nix nevertheless includes all tools `.py` bytes, so these inactive scripts currently affect qualification source identity. That is a source-filter issue, not evidence they execute. Optional editor/agent packages and Dependabot cannot be proven unused from repository references; usage outside tracked commands is UNKNOWN. Do not delete them on age or lack of CI execution alone.

## 21. Low-risk delete shortlist

**HIGH confidence DELETE_CANDIDATE:** `tools/audit_r15.py` (125 lines) and `tools/audit_r16.py` (89 lines), total **214 Python lines / two live-tree files**, after preserving reproducibility through their existing historical commits and audit outputs. Their responsibility no longer needs to exist as active maintenance tooling; no Python→Go port is warranted.

This is not an unconditional link-safe deletion. Existing immutable reports link to or instruct running these paths. Future deletion must provide a historical retrieval convention that leaves reports untouched (for example, documenting the archived commit/path in current maintenance docs), verify Git objects/source availability, and make clear that historical reproduction requires the historical checkout. Existing review reports/JSON remain verbatim. If the project requires every historical relative source link to resolve in current main, archive/move or keep these two files instead; no unconditional HIGH-confidence file deletion then remains. No other file qualifies for a low-risk deletion shortlist.

## 22. Merge candidates before migration

First merge the **responsibility**, not syntax: one read-only Go integrity/profile verifier with the r18 four-artifact negative regressions replaces the duplicate Python hash block. Keep a clear spec-structure responsibility, but review every prose fingerprint before implementing it in Go. Consolidate JSON validation around the union of schema and Go protections; merely deleting ValidateShape or the schema interpreter would be premature.

Separate shared corpus contracts/utilities from writer declarations and evaluator logic where it makes dependencies clearer. The generator currently imports vectors, which contains both wire types and runtime verification. Extracting a small read-only contract/integrity layer may improve that relationship, but excessive package splitting would increase surface. The existing 148-line type file is a reasonable starting point, not evidence of required new packages.

Do not merge generation and verification, explicit expected results and evaluator output, or independent formula tests and production codec helpers. Centralize portable command orchestration only after preserving equivalent platform invocation and CLI evidence.

## 23. Migration candidates

After consolidation, migrate schema-subset/strict JSON validation and its eight negative tests to Go, then the reduced structural/revision guard. Preserve the checked-in schema contract, required field presence, integer exactness, nested unknown and duplicate rejection, safe paths, artifact pins and archived revision continuity. A standard Go schema dependency would add supply-chain surface; a small reviewed subset may be enough, but evaluate correctness/maintenance rather than assuming fewer dependencies always wins.

Migrate the independent Python wire tests last, with independent stdlib Go HMAC/SHA/binary implementation that does **not** call cryptov1/object/generator helpers. Test observable intermediates and negative inputs; preserve existing modeling limits. This migration has medium confidence because same-language oracle review must avoid hidden sharing. No Python responsibility is intrinsically irreducible to Go; independence and schema completeness are the irreducible requirements.

## 24. Non-negotiable contracts

- Normative spec bytes, revision/history discipline and raw artifact pins remain explicit and reviewed.
- Exact portable-case accounting: all physical cases exactly once, all-and-only required IDs/hashes, matching expected outcomes, safe paths and strict JSON/schema contract.
- Vector expected outcomes remain authored evidence, never changed to make implementation pass.
- Generation remains explicit, deterministic where applicable, reviewed via temporary generation/compare and then normal independent verification.
- Generator never becomes its only verifier; preserve independent wire formulas, fixed known answers, graph oracle and authenticated-negative boundary proofs.
- All-package race coverage, all three bounded fuzz targets, original seeds/failures, 10s / parallel 2 / `-run '^$'` remain intact.
- Go module checksums, fixed vendorHash/dependency closure, offline pinned Linux Nix authority and existing portability CI remain intact.
- Historical review evidence stays preserved; modeled storage outcomes do not become claims of host crash/race durability or application conformance.
- Human-only Nix-last policy and truthful evidence statuses persist; pending is not passed.

## 25. Minimal future architecture

| Owner | Minimum responsibility | Exclude from that responsibility |
| --- | --- | --- |
| Normative Markdown + JSON contracts | Spec/revision/history, portable cases, schemas, manifest/profile pins | No automatic expected-result repair |
| Go conformance utility | Strict schema/inventory/pin validation; reduced spec structure; evaluate every required case | No writes during verification |
| Go generator command | Deterministic fixture construction and explicit reviewed pin refresh | No invocation from checks, tests or CI |
| Go tests/oracles | Semantic regressions, independent formulas/known answers, CLI behavior, race/fuzz | No evaluator-derived expected fixtures |
| Nix | Fixed dependency closure, hermetic Linux suite, minimal tool environment | No protocol parsing or expectations in expressions |
| Make | Thin stable human/check/race/fuzz entry points if retained | No validation logic |
| YAML | Linux gate, supplemental portability, pinned update policy | No alternative semantic authority |
| Shell | Environment/filesystem/process wiring and formatting predicate | No JSON/revision/protocol parsing |
| Current contributor docs | One qualification/revision procedure and cross-links | Historical reports are preserved separately |

Two Go commands remain justified. No release framework, Python packaging, Node tooling, new shell framework or additional generator is required. A Go-only executable maintenance surface is feasible without changing the normative protocol or weakening evidence.

## 26. Actionable pruning waves

These are proposals for separately authorized milestones. Coverage flags/targets/matrix remain unchanged in this audit.

| Wave | Files / action | Risk | Invariant and required evidence | Expected reduction |
| --- | --- | --- | --- | --- |
| 1 — dead/superseded deletion | Retire tools/audit_r15.py and audit_r16.py from active tree, subject to historical retrieval/link policy | LOW semantic; MEDIUM archival/link risk | Preserve historical reports/outputs/source retrieval; caller search; compare all live inputs. Because tools `.py` are Nix-filtered inputs, deletion invalidates Nix; human final flake check required. Do not infer unchanged Nix result | 214 Python lines, 2 files; no language removal yet |
| 2 — merge duplicate responsibilities | check_spec pin block/test_check_spec behavior into vectors integrity tests; review redundant full/verify-only work; consolidate shape requirements before port | MEDIUM | Mutation tests for all pins, exact case accounting/presence, CLI failure behavior; baseline once/focused/final complete host suite; human Nix after frozen sources. No suite reduction by default | Roughly 15–30 net lines from duplicate pin code; any invocation reduction requires separate evidence review |
| 3 — migrate remaining necessary Python | check_vectors + test_check_vectors; reduced check_spec + remaining tests; independent test_r16_vectors last | HIGH for schema/oracle equivalence; MEDIUM for structure | Differential accepted/rejected contract inputs, current corpus, raw byte/pin mutations, independent intermediate oracle review; generator change comparison only if generator touched; complete check/vet/race/fuzz, human Nix, remote portability after commit/push | Remove remaining 493 Python lines / 5 files before Go replacements; add roughly 350–550 Go lines; net change may range from 143 fewer to 57 more lines |
| 4 — simplify orchestration/environment/docs | Makefile, flake.nix shell lists/source filter, workflow command mapping/Python setup, current docs, optional .envrc boundary | MEDIUM platform/build; LOW prose | Preserve eight-target meanings unless explicitly changed, all qualification/race/fuzz, five supplemental combinations, fixed dependencies/Linux authority. Validate source filter/static workflow; affected host suite and human Nix; actual remote CI remains supplemental | Perhaps 10–30 orchestration lines; Python runtime/setup/dev-shell entries can go only once no active Python remains |

Wave 1 does not justify deleting historical reports. A source-filter-only alternative can remove historical scripts from qualification inputs while retaining source links, but reduces no repository LOC and still needs Nix qualification. Wave 2 must precede literal ports. Wave 3 may be split by responsibility to keep review tractable. Wave 4 should simplify current docs, not rewrite historical execution records.

## 27. Expected LOC/language reduction

Baseline executable/orchestration inventory: Go 3,044 + Python 707 + Make 35 + Nix 135 + workflow YAML 63 + Dependabot YAML 10 + .envrc 2 = **4,096 physical lines** (excluding locks/contracts/prose). Wave 1 saves exactly 214 lines if deletion policy permits. Python elimination then removes 493 existing lines but adds an estimated 350–550 Go lines; this is language/tool reduction more reliably than LOC reduction. Wave 2 savings overlap Wave 3; do not add them twice.

With optional 10–30 orchestration lines saved, the combined indicative net reduction is **167–387 lines**, around **4–9%** of this baseline. These are design estimates, not promised changes. A schema library or equivalent independent oracle can increase Go/dependency LOC. The reliable simplification is seven Python files and Python/setup-python/runtime knowledge leaving the active maintenance path (two historical files may instead remain archived), with Go/Nix/Make/YAML/minimal shell still justified. No projected reduction assumes deleting final semantic/race/fuzz coverage, schemas or evidence.

## 28. Unknowns and limits

External callers of Make/CLI paths and use of editor/jailed-agent tools are not observable in this checkout. Historical source-link retention policy is unresolved; it gates unconditional script deletion. No generator run was performed, so current deterministic regeneration is not newly established. No executable test or schema differential audit was run, so migration equivalence is unproven. No independent full AES/Argon2 interoperability evidence is introduced. Actual remote workflow runtime compatibility and newer human Nix results cannot be inferred from local static inspection.

Some guards were introduced during rXX hardening from design requirements, not documented observed failures. Only the r18 stale-pin incident is specifically classified here as a real-failure regression. Other historical tests are not mislabeled as real failures merely because they defend bugs. Optional authoring dependencies may be separable, but owner usage evidence is needed before deletion.

## 29. Qualification disposition and validation

This is DOCS_ONLY. [AGENTS.md](../AGENTS.md), the actual flake filter, check_spec and its tests were inspected: this report is excluded and not consumed. No fresh human Nix request is needed. Inherited [Nix/CI report](SPEC_NIX_CI_UNIFICATION_REPORT.md) records passing host/isolated-source suites: 92 cases, 14 Python tests, seven Go test packages, vet/race and all three bounded fuzz targets. Those are inherited observations, not agent executions here. The [Actions refresh report](SPEC_GITHUB_ACTIONS_REFRESH_REPORT.md) preserves the later human Nix gate as **PENDING** and remote CI as **NOT RUN**. The earlier r19 human-reported Nix pass predates unification and does not resolve it. No cached Nix pass is claimed.

| Execution | This milestone |
| --- | --- |
| Baseline make check | 0 |
| Focused Python / Go test invocations | 0 / 0 |
| Focused fuzz: Dispatch / Open / graph | 0 / 0 / 0 |
| Final complete suites / incomplete attempts | 0 / 0 |
| Final make check / vet / race / fuzz | 0 / 0 / 0 / 0 |
| Generator | 0 |
| Human Nix / agent Nix | 0 / 0; not requested / not run |
| Remote CI / cached Nix observations | 0 / none newly observed |

Validation: all **260 existing tracked files** are byte-identical to starting HEAD; the explicit protected subset contains **131 files**. Protected raw bytes and path membership compared against starting HEAD for spec/, vectors/, requirements/, conformance/, tools/, go.work and every go.mod/go.sum; root go.mod/go.sum absent at both boundaries. Makefile, flake.nix, flake.lock, .github configuration, process docs and all existing history/reports also remain byte-identical. No protected additions/deletions. Changed-path allowlist permits only this report. Local Markdown links resolve, code fences balance, new-file whitespace checks and `git diff --check` pass. Static inventory/byte/link checks are not Python verifier or unit-test executions.

## 30. Final Git state

`main` remains at `ad4e9000b6cdd4c7eebbaeeac06178750984f836`. Index and tracked files are unchanged. The only added path is untracked, unstaged `review/SPEC_SUPPORT_TOOLING_AUDIT.md`. No executable content changed; no spec/vector regeneration, stage, commit, push, tag, release, Nix command or CI dispatch occurred. Present this report and stop; implementation belongs to future authorized milestones.
