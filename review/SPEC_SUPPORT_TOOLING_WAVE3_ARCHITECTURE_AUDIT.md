# Totipo support tooling Wave 3 architecture audit

## 1. Starting state and scope

Date: 2026-10-11. Repository: `/home/niki/Sources/totipo-spec` only. Before inspection, the required commands established:

```text
git branch --show-current: main
git rev-parse HEAD: 1f78b574f0855a515dfadaccd7975e513bd6c962
git status --short: empty
git diff --cached --quiet: exit 0
```

Protocol remains **v1/r19**, moving pre-RC. Static inventory: **92 manifest cases, 92 physical case JSON files**, all required by the profile. Categories: bootstrap 5, crypto 5, encoding 30, fold 6, graph 13, metadata 7, size 1, storage 13, totp 3, vault 9. No protocol revision, corpus/schema/dependency change, regeneration, implementation, or semantic refactor is authorized by this audit.

Read AGENTS.md, README.md, CONTRIBUTING.md, Makefile, flake.nix, go.work, the five requested committed support/policy reports, all remaining Python files, generator/consumer implementations, package relationships, schemas, tests, module metadata and CI configuration. Inventory and dependency conclusions are static source observations, not fresh qualification passes. Official Go documentation was consulted for module/workspace/internal behavior; no repository tools under audit were executed.

Inherited evidence: Wave 1, Wave 2 and the portability-fix reports record passing host suites and human-reported Nix passes. The latest Git-worktree-policy report records one human-reported `nix flake check` PASS, without cache/fresh-build attribution. The portability fix explicitly leaves remote macOS/Windows success unverified; the policy report also claims no new remote evidence. These statuses are preserved, not upgraded by this audit. Historical `path:.` commands in reports remain historical; current policy uses Git-aware commands.

## 2. Classifications and decision

Requested classifications: **SUPPORT_TOOLING_ARCHITECTURE**, **LANGUAGE_SIMPLIFICATION**, **GENERATOR_VERIFIER_INDEPENDENCE**, **READ_ONLY_AUDIT**. Goal labels also include LANGUAGE_SURFACE_AUDIT and READ_ONLY. Under the interactive ladder, the sole output is **DOCS_ONLY**, outside qualificationSource and current structure consumers.

Recommendation: **A. TWO_MODULE_SPLIT_RECOMMENDED**, using option C below: two modules with the existing workspace retained for convenience, no shared Go module. The split must remove existing semantic ancestry, not merely relocate the command. Retain an import-restricted independent oracle test package in the verifier module; a third module is unnecessary for its current small role. Keep Python until replacement evidence exists.

## 3. Executable/tool responsibility matrix

Y means an actual responsibility; I means indirect through called code; T means test fixtures/temporary files only; — means none. S/V/W/Vf = reads spec / reads vectors / writes vectors / verifies vectors. Sem/E/H/R = protocol semantics / expected-value construction / hashes or pins / repository structure. O/P/C = independent oracle / generated-artifact producer / release-or-CI-only. Package rows describe all production files named; test rows enumerate the remaining executable code.

| Path | Language; entry; callers | Responsibility | S/V/W/Vf | Sem/E/H/R | O/P/C |
| --- | --- | --- | --- | --- | --- |
| conformance/cmd/generate-vectors/main.go | Go; main; deliberate human command documented in conformance/README and vectors/FORMAT | Fixtures, manifest, moving profile; removes undeclared case files | Y/Y/Y/— | Y/Y/Y/— | —/Y/— |
| conformance/cmd/totipo-conformance/main.go | Go; main, -root, -verify-only; Make, portability CI, human | Read/profile validation; full mode executes every case; rejects positional arguments | I/Y/—/Y | I/—/I/I | —/—/— |
| conformance/internal/vectors/{types,runner,post_aead}.go | Go; Read, VerifyProfile, Run; CLI/tests; generator uses types/Decode/Hash | Contracts, strict JSON, inventory/paths/hashes/pins, semantic dispatch, authenticated-defect proof | Y/Y/—/Y | Y/—/Y/Y | —/—/— |
| conformance/internal/tlv/tlv.go | Go; Take, Encode, U16/U32/U64; object/generator/tests | u16be framing and integer bytes | —/I/—/I | Y/ I /—/— | —/—/— |
| conformance/internal/object/object.go | Go; Dispatch, Object.Encode, ValueBytes; crypto/graph/storage/vectors/generator/tests | TOKEN grammar, validation, canonical bytes, value projection | —/I/—/I | Y/I/—/— | —/—/— |
| conformance/internal/cryptov1/crypto.go | Go; Derive/Seal/Open/Wrap/Unwrap etc.; vectors/storage/generator/tests | Object and bootstrap crypto composition | —/I/—/I | Y/I/Y/— | —/—/— |
| conformance/internal/graph/graph.go | Go; Store, Evaluate, Fold/FoldObjects, FromObject; vectors/tests; generator uses Node/Result only | Observed graph, SCC-equivalent heads, unresolved edges, fold snapshots | —/I/—/I | Y/—/—/— | —/—/— |
| conformance/internal/storage/{storage,publication}.go | Go; Scan, Candidate, OpenBootstrap, Install/Publish/Create; vectors/tests | Modeled namespace, authentication boundary, publication/create outcomes; no live host store adapter | —/I/—/I | Y/—/—/— | —/—/— |
| conformance/internal/totp/totp.go | Go; Code; vectors/tests | RFC 6238 with credential bounds and unsigned time | —/I/—/I | Y/—/—/— | —/—/— |
| conformance/internal/{tlv,object,cryptov1,graph,storage,totp}/*_test.go | Go testing; test/race/fuzz | Package regressions, three fuzz targets; graph includes independent matrix closure | —/T/—/T | Y/T/T/— | graph subset/—/— |
| conformance/internal/vectors/{runner,integrity,path,integration}_test.go | Go testing; test/race | Corpus execution, strict decode, pins/set mutations, path portability, encrypted folds, authenticated negatives | Y/Y/T/Y | Y/T/Y/Y | authored outcomes/—/— |
| tools/check_spec.py | Python; top-level CLI; Make/CI/test_check_spec | Spec structure/history/registry and policy anchors | Y/—/—/— | prose assertions/—/—/Y | independent literals/—/— |
| tools/check_vectors.py | Python; main/read/validate/pairs; Make/CI/two test modules | JSON schema subset, duplicate rejection, listed-path containment | —/Y/—/Y | shape only/—/—/Y | schema contract/—/— |
| tools/test_check_spec.py | Python unittest; Make/CI | Two-file isolated structure CLI and mutations | Y/—/—/— | prose mutations/T/—/Y | —/—/— |
| tools/test_check_vectors.py | Python unittest; Make/CI | Shape/presence/type/operation negatives | —/Y/T/Y | shape only/T/—/Y | schema evidence/—/— |
| tools/test_r16_vectors.py | Python unittest; Make/CI | Independent selected wire/key/VAULT/defect inputs | —/Y/—/Y | Y/—/Y/— | Y/—/— |
| Makefile | Make + shell; help/spec-check/test/conformance/verify/check/race/fuzz; human/Nix | Stable orchestration, cache default, exact fuzz policy | I/I/—/I | I/—/I/I | —/—/— |
| flake.nix | Nix + shell; checks.spec, goDependencies, devShells.default, formatter | Linux pinned qualification, workspace vendoring, editor/jail environment; **no apps or application package** | I/I/—/I | I/—/dependency hash/I | —/vendor tree/qualification is CI + human |
| .github/workflows/conformance.yml | YAML + shell; push/PR/manual | Linux Nix gate; five supplemental OS/Go combinations | I/I/—/I | I/—/I/I | —/—/CI-only |
| .github/dependabot.yml | YAML service configuration; weekly | Go module and Actions update proposals | —/—/—/— | —/—/—/— | —/—/CI service only |
| .envrc | direnv shell; human | use flake; layout go | —/—/—/— | —/—/—/— | —/—/— |

No standalone shell script, release executable, publishing tool or remaining historical audit executable exists. JSON schemas/manifest/profile are inputs, not executables. go.work, go.mod/go.sum and flake.lock are dependency metadata. Nix vendoring produces dependency artifacts, not vectors. Generator command compilation during Go tests does not execute main.

Current exact orchestration: `check: spec-check test conformance verify`, serial; spec-check invokes Python structure; test invokes Python discovery then `go -C conformance test -count=1 ./...`; conformance invokes root workspace CLI; verify invokes Python schema then CLI -verify-only. Race is all-package `-race -count=1 ./...`. Fuzz runs FuzzDispatch/internal/object, FuzzOpen/internal/cryptov1, FuzzArrivalAndDisappearance/internal/graph, each `-run '^$' -fuzztime 10s -parallel 2`.

Nix calls `go work vendor` for dependencies, then formatting, make check, vet, race and fuzz with offline dependency access and temporary caches. Portability YAML calls both Python scripts/discovery, gofmt, conformance vet/test and both CLI modes. Matrix: macOS/Windows × Go 1.23.x/stable, plus Ubuntu 24.04 Go 1.23.x. Linux Nix is a separate authoritative job. No qualification workflow runs generation.

## 4. Every remaining Python file

Physical LOC includes comments/blanks. Exactly five files, **486 LOC**; no external dependencies. Executable means directly invoked checker rather than unittest module. All local imports are listed.

| File; LOC; kind | Imports | Inputs → outputs; trust boundary; Go overlap | Responsibility classification and reason |
| --- | --- | --- | --- |
| tools/check_spec.py; 115; executable | pathlib.Path, re | CWD spec + archived history → PASS or assertion/error; authored prose/policy validation; Go only asserts artifact metadata/pins, not these rules | MIGRATE_TO_VERIFIER_GO for repository validation family, specifically separate policy command, then DELETE_AFTER_MIGRATION; no language-based reason to retain Python |
| tools/test_check_spec.py; 56; test | pathlib.Path, shutil, subprocess, sys, tempfile, unittest | Copies only spec/archive into temp root; launches checker; current-pass + eight mutation failures → unittest results; does not need vectors/profile | MIGRATE_TO_VERIFIER_GO policy-command tests, then DELETE_AFTER_MIGRATION; retain observable CLI behavior and two-file dependency boundary |
| tools/check_vectors.py; 101; executable/library | json, re, pathlib.Path | Manifest/schema/case bytes → parsed values or Invalid/PASS; untrusted JSON and resolved paths; overlaps Go strict decoding/shape, not full contract | MIGRATE_TO_VERIFIER_GO, then DELETE_AFTER_MIGRATION; required presence and schema evaluation belong to artifact acceptance |
| tools/test_check_vectors.py; 111; test | copy, json, unittest, pathlib.Path; local check_vectors.Invalid/read/validate/pairs | Committed schemas and cases, in-memory mutations → acceptance/rejection; tests schema trust boundary; some Go negatives overlap | MIGRATE_TO_VERIFIER_GO, then DELETE_AFTER_MIGRATION; port all ten behavioral tests, expand field-presence mutation coverage |
| tools/test_r16_vectors.py; 103; test | base64, hashlib, hmac, struct, unittest, pathlib.Path; local check_vectors.read | Manifest/listed cases → independent arithmetic/encoding comparisons; consumes generator-produced expectations but recomputes from inputs | KEEP_INDEPENDENT, eventually DELETE_AFTER_MIGRATION into guarded stdlib-only Go oracle tests; never implement by delegation to current helpers |

No responsibility is MIGRATE_TO_GENERATOR_GO: Python currently constructs no persistent fixture outputs. KEEP_PYTHON is a valid temporary execution state for all five until gates pass, not an intrinsic architectural requirement. UNKNOWN: unobserved external callers; no repository-local responsibility is unidentified. Output files from tests are disposable temporary fixtures only. Python removal reduces language/runtime surface; it does not itself prove simplification or independence.

## 5. Current generator graph and ownership

All paths below are within conformance today:

```text
cmd/generate-vectors/main.go
  -> internal/vectors: Case + all nested wire structs; Manifest/Entry/Profile/Pin;
                       Decode(RFC source cases), Hash(case/artifact bytes)
       -> graph + object (types)
       -> crypto + graph + object + storage + totp (runner package closure)
  -> internal/object: Object, Encode -> Dispatch -> tlv.Take
  -> internal/tlv: Encode, Take, U16/U32/U64
  -> internal/cryptov1: Derive, Keys.Seal/ObjectKey, Padded, AAD, MAC, VaultID
       -> object.Capacity; x/crypto/hkdf; stdlib crypto
  -> internal/graph: Node, Result (data only)
  -> stdlib aes/cipher directly for intentionally malformed authenticated fixtures
```

Package dependency is broader than functions actually executed: importing vectors brings its semantic runner dependencies into the build even though generator never calls Run; importing graph brings its object dependency even though only structs are used.

Ownership now: `token`/`encoded` construct object inputs using Object.Encode, which validates with the same Dispatch parser used by consumer. `field` and malformed-field mutation loop use shared TLV Encode/Take. `envelope` uses shared Derive/Seal/Padded/ObjectKey/AAD; OBJECT_ID is shared MAC over P. `postAEAD` reuses that derivation/padding, mutates length/padding/ID, and directly seals with AES-GCM. Bootstrap records/keys are preserved source known answers; generator recomputes VAULT_ID (including changed-record ID), not Argon2 wrapping. RFC TOTP files are read with vectors.Decode, retained byte-for-byte, and not recalculated or written.

Graph expected heads/values/unresolved IDs are explicit main-source declarations, **not graph.Evaluate**. `result` computes Conflicting from declared values. `graphCase` tracks input records/removals/integrity failures and attaches records for explicitly declared heads; this is expected snapshot construction, not SCC/head evaluation. Fold expectation loop independently stages 4 then 3 original IDs plus carry, sorts parents, and supplies symbolic stage IDs; it does not call graph.Fold. Storage/workflow outcomes are explicit declarations; generator does not call storage.Scan/Create/Publish. Supported/Invalid strings come from object constants.

End of main sorts cases, writes all non-TOTP case JSON with indentation/newline, removes physical files not in its keep set, hashes files into manifest, and hashes spec/manifest/two schemas into moving profile plus required ID/hash list. It reads spec bytes only to pin them. It neither generates nor validates schemas. Packages participating in expected-content production are **vectors, object, tlv, cryptov1, graph types**, plus generator-local declarations/arithmetic and standard primitives. storage/totp are dependency-closure ancestry through vectors, not generators of expectations.

## 6. Current verifier/consumer graph

```text
cmd/totipo-conformance
  -> vectors.Read -> Decode(uniqueJSON + strict typed decode)
                  -> lexical path validation + symmetric Abs/EvalSymlinks containment
                  -> case raw hash + identity/outcome + ValidateShape + physical inventory
  -> vectors.VerifyProfile -> independent expected r19; four raw artifact pins;
                              exact all-and-only required IDs/hashes
  -> vectors.Run (unless -verify-only)
       -> crypto/dispatch: cryptov1 Derive/Seal/Open + intermediates;
                           object Dispatch/Encode
       -> bootstrap: Wrap/Unwrap/WrapKey/VaultID + changed-tag proof
       -> post-aead: raw AES-GCM open; verify isolated defect;
                     strict Open rejection -> storage.Scan contributes no state
       -> graph: New/Add/Remove/Evaluate -> exact expected graph.Result
       -> fold: FoldObjects -> Encode/Dispatch/ValueBytes + exact metadata/stage parents
       -> storage: Scan -> cryptov1.Open -> object.Dispatch
       -> workflow: storage.Store.Publish or Create; existing/vault preservation
       -> totp: totp.Code + independently declared counter/time/code fields
object -> tlv; crypto -> object.Capacity + hkdf/argon2;
graph -> object; storage -> crypto/object; totp -> stdlib only
```

Every internal package participates in committed-artifact verification. Future sibling families/noncandidate paths are ignored by storage.Candidate; exact objects-v1 regular filenames alone enter v1. Authentic envelope with invalid TOKEN grammar is Invalid; authentication/size failures are INVALID_STORAGE. There is no signature subsystem or future-family semantic decoder. Storage models observed entry kinds/backend outcomes, not live syscall/crash proofs. Verify-only validates schema **pins**, not schema contents; Python is currently necessary for actual schema acceptance.

## 7. Shared implementation ancestry

Complete repository package intersection: **tlv, object, cryptov1, graph, vectors, storage, totp** in the compiled generator/consumer closures. Last two are incidental generator dependencies. Direct generator imports: five packages (all except storage/totp). The following enumerates shared repository symbols used to construct/check artifacts; functions reached transitively are included. stdlib and x/crypto primitive overlap is addressed in section 18.

| Shared item | Class | Independence effect |
| --- | --- | --- |
| tlv.Field, Take, Encode, U16/U32/U64; ErrEncoding | SEMANTIC_SHARED (errors mechanical) | Yes: same framing and malformed-field machinery can bless a common mistake |
| object.Object (including optional presence), Encode, Dispatch, cursor.get/u8/next, validText; Supported/Invalid/Capacity/MaxParents; ErrGrammar | SEMANTIC_SHARED; JSON field shape also structural | Yes: writer validates via verifier parser; same tag order, widths, limits and enums; type carries behavior |
| cryptov1.Keys, MAC, expand, Derive, ObjectKey, AAD, gcm, Padded, Seal, VaultID; ErrCrypto | SEMANTIC_SHARED (error mechanical) | Yes: common key domains, HKDF composition, object-ID, nonce/AAD, padding and hash representation checks |
| graph.Node, graph.Result | STRUCTURAL_SHARED with semantic fields | Conservative yes: same projection/expectation representation can normalize away omitted information; no shared Evaluate/Fold execution by generator |
| vectors.Case, Crypto, Bootstrap, TOTP/TOTPRow, GraphCase/Step, FoldCase, StorageEntry/StorageCase, Workflow, PostAEAD, Manifest/Entry/Profile/Pin | STRUCTURAL_SHARED | Lower but real risk: writer structs and consumer structs coupled; protocol types embedded; no neutral shared-type extraction recommended |
| vectors.Decode -> uniqueJSON | SEMANTIC_SHARED validation | Yes for source-RFC ingestion; generator depends on acceptance logic residing with evaluator, even if only JSON validation is exercised |
| vectors.Hash | MECHANICAL_SHARED SHA-256 hex | Low alone: standard checksum; still duplicate tiny wrapper to simplify policy |
| vectors runner's remaining functions; cryptov1 Open/Wrap/Unwrap/WrapKey/PasswordValid; object.ValueBytes; graph evaluators; storage/totp code | Package closure shared, not generator call overlap | Does not mean generator computes these results; creates an easy future convergence path. Remove dependency edges rather than claiming all these functions currently author expectations |

No custom shared filesystem or hex helper package exists. Generator hx/raw/write/id and verifier unhex/equalHex/resolvedPath/containedPath are already separate. Graph canonical strings in abstract fixtures are generator-local, while real graph projection uses object encoding. Package identity is evidence of common ancestry; repeated invocations of one function are not independent implementations. Copying these helpers into a new module can break edges but does not immediately establish independent derivation.

## 8. Exact independent Python oracle

`test_r16_vectors.py` has three tests, no shared Go constants or generated source, and only a local duplicate-safe JSON reader. It iterates manifest-listed cases; it does not independently establish manifest completeness or pin relationships.

For every case with `crypto`, it consumes semantic_hex and root_hex and computes:

```text
PRK = HMAC-SHA256(32 zero bytes, root)
K_id = HMAC-SHA256(PRK, literal "totipo/v1/object-id" || 0x01)
K_object_root = HMAC-SHA256(PRK, literal "totipo/v1/object-key-root" || 0x01)
OBJECT_ID = HMAC-SHA256(K_id, P)
K_object = HMAC-SHA256(K_object_root, literal "totipo/v1/object-key" || OBJECT_ID || 0x01)
nonce = OBJECT_ID[0:12]
AAD = literal "totipo/v1/object" || OBJECT_ID
padded = u16be(len(P)) || P || zero[1006-len(P)]
```

Compares id_key_hex, object_root_key_hex, object_id, object_key_hex, nonce_hex, aad_hex, padded_plaintext_hex; checks object_hex decodes to 1024 bytes. HKDF is independently expressed as extract plus single-block expand, not via Go's hkdf package. It does not compare ciphertext_hex/gcm_tag_hex or cryptographically open object_hex.

Where input exists, independently walks P using `>HH` tag/length and exact final offset; base64-decodes identity/parents/secret, enforces at most four parents, and constructs the expected field list: tags 1–10 for identity/u16 parent count/repeated parents/status/UTF-8 issuer/account/algorithm/digits/u32 period/secret, optional tag 11 UTF-8 client_name and 12 u64 client_time by key presence. Compares complete ordered (tag,value) lists. It is not a full negative grammar parser; e.g. it is not an independent parent-order validator.

For every bootstrap it compares SHA-256(exact record_hex bytes) to vault_id_hex; it does not independently validate header/87-byte grammar or changed_record_hex hash. For maximum case it asserts semantic_length=1005 and final padded byte zero. The general padding computation and actual field encoding give selected capacity evidence, but the maximum test does **not** independently sum every maximum field length.

For three post_aead cases it consumes semantic_hex/root_hex and post_aead.{object_id,encryption_plaintext_hex,object_hex}; checks 1008/1024 byte lengths and independently recomputes correct ID/canonical padding. Nonzero-padding: correct ID, correct length/prefix, differing padding. ID mismatch: differing ID, exact canonical plaintext. Length-invalid: correct ID, prefix length exactly 1007, unchanged remainder. It does not prove those declared plaintexts are encrypted by the supplied ciphertext: Go runPostAEAD supplies that separate proof.

No AES-GCM, Argon2id, TOTP, graph or signature operation occurs in this oracle. It reads implementation-produced expected fields as **claims to compare**, not as authoritative formulas/constants. A generator-only or shared-Go-helper bug in covered derivation/encoding can be detected. Jointly wrong inputs/expectations outside its invariants, shared spec misunderstanding, unlisted cases or AEAD-only bugs can escape. Rewriting it with cryptov1/object/tlv would lose the independent covered composition/encoding/literal checks, even with a different command/module name.

## 9. Required independence contract

1. Generator and verifier have sibling module paths `totipo/generator` and `totipo/conformance`; all implementation/data types are below each module's own internal directory. Neither imports the other, including tests and tools. Neither depends on another local/shared semantic module.
2. No shared Totipo-specific TLV codec, canonicalization, crypto composition/opening, OBJECT_ID, TOTP formula, graph/state/fold algorithm, size/capacity formula, expected-value code, or semantic constant source. Primitives remain shareable (section 18).
3. Only serialized files cross the implementation boundary. No subprocess delegation to the other implementation, linked plugin, generated shared Go source, source include/symlink, common code generator/template, go:linkname bridge or test import bypass. No automatic regeneration in ordinary verification.
4. Generator expectations retain provenance: explicitly authored outcomes, source known answers, and computed intermediates remain distinct. Verifier recomputes from input bytes/fields, compares assertions, never repairs them.
5. Keep selected independent formula tests separately implemented and import-restricted to stdlib. Do not derive their literals from production constants, generated schemas or writer data. Independent schema contract remains authored JSON; no schema generation from writer structs.
6. Review semantic changes on each side against spec/known answers independently. A source copy with rewritten import paths is a lineage-preserving extraction, not evidence of independent implementation. Rework generator-owned composition from the spec and verify different intermediate failures with mutation tests; leave consumer behavior stable where practical.

Two mutual-import prohibitions alone are insufficient. Dependency guards enforce ancestry limits; independently authored algorithms, negative evidence and review establish meaningful computation independence. No static guard proves two authors did not copy an algorithm or misread the same spec. Same-language independence is valuable but does not replace eventual external interoperability evidence.

## 10. Module-boundary options

| Option | Isolation and cross-import risk | Ergonomics, CI/Nix and maintenance | Judgment |
| --- | --- | --- | --- |
| A: one module, package guards | Strong if exhaustive graph/source guard; current conformance/internal permits both commands; sibling internal subtrees could add visibility restrictions | One mod/sum/vendor closure, easiest tooling; guard must own every layer restriction | Sufficient in principle; less obvious ownership and easier direct reuse under current layout |
| B: two modules, no workspace | Separate build/dependency roots plus own internal trees; public package/third-module reuse still possible | Explicit -C commands everywhere; both module tests/vet/race; either two vendor derivations/hashes or changed vendoring design; two mod/sum maintenance | Viable isolation, unnecessary root-command/Nix churn here |
| C: two modules + go.work | Same internal restrictions; workspace makes public local cross-imports easy unless prohibited | Preserves root go run and current go work vendor model; add module source/cache/CI coverage; no independent module release/version scheme needed | Recommended: smallest practical split for existing workspace |
| D: generator/verifier/tiny shared module | Third module immediately creates exception channel; types/helpers can acquire semantics | Extra mod/version/hash/source/guard maintenance; small duplicated helpers save little | Reject; no justified shared module |

All options need explicit enforcement and independent implementations. Separate go.mod files are not a substitute for removing current shared code. C has two modules total, not an oracle module. Workspace minimal-version selection combines dependency constraints during workspace builds; isolated module checks still matter for undeclared dependency masking. Primitive version duplication is acceptable; no dependency upgrades are proposed.

## 11. What separate modules enforce and do not enforce

A module supplies a separate package/dependency root; `./...` from one module does not supply tests for the sibling module. Import resolution needs dependency/workspace configuration for external packages. Visibility is supplied by **internal directory ancestry**, not a blanket module firewall. With proposed sibling module paths, generator cannot import conformance/internal and vice versa; choosing nested module paths or exporting a bridging package would defeat the intended pattern. [Go command internal-directory documentation](https://pkg.go.dev/cmd/go#hdr-Internal_Directories).

Public packages remain importable. `require` plus local `replace` resolves sibling/shared code; go.work can resolve public workspace packages without sibling require entries. Workspace replace directives can override module replacements. Local replacements point to a directory with a module, not relative package imports: `import "../conformance/..."` is not the module-mode escape mechanism. Both modules can require a third module, execute another binary, or compile copied/generated semantic source; none of these creates an automatic independent implementation. Ban replacements/local dependencies and additional workspace members in this repository's policy. [Go modules reference: replacements and workspaces](https://go.dev/ref/mod).

No new release tags/version lifecycle is required for two unpublished local commands. Independent `GOWORK=off` checks catch workspace dependency masking, but cannot detect copied code or subprocess delegation. The internal restriction plus guard is the enforceable structural boundary; policy and parity evidence address remaining risks.

## 12. Smallest mechanical guard design (proposal only)

Use one stdlib-only repository-policy command/test in conformance (e.g. `cmd/check-architecture` with a small internal repository-policy implementation). It is not a third module or shared semantic library. Invoke through a Make boundary target included in check, and Linux/supplemental CI. It may run the Go tool for inspection; no generation occurs. Inspection is deliberately narrower than semantic verification.

Exact inputs and predicates:

- Walk **all** .go files under both module roots, including _test.go, generated files and build-tag-excluded source. Parse imports using go/parser; reject imports of the sibling module and any nonallowlisted external package. Resolve root inventory without relying on .git (Nix filtered source has no repository history). Require exactly the approved two go.mod roots; root module absent. Reject symlinked Go source and source directories escaping their module.
- Inspect declared module paths/go versions, require/replace/exclude/tool directives and go.work membership/replacements through `go mod edit -json` / `go work edit -json` output. Require only the two sibling module paths, two workspace use entries, no replace/local dependency bridge, no new third local module. Unknown directives/dependency additions require explicit policy review. Examine all metadata, not just selected dependency graph.
- Run `go list -deps -test -json ./...` **in each module**, with the qualification environment's actual workspace/vendor mode. Fail any load/DepsErrors, including partial output. Inspect ImportPath, Imports, TestImports, XTestImports, ForTest, Standard, Dir, Module.Path/Dir and Module.Replace (including replacement destinations). Nonstandard code must be owned by that module or the precise approved primitive dependency packages. Check real directories as well as module names to prevent alias/replacement bypass. Exclude only go-tool synthetic test wrappers, never their imports.
- Current primitive allowlist begins with x/crypto/{hkdf,argon2} and required dependency implementation packages such as x/crypto/{blake2b,internal/*} and x/sys/cpu: implementation milestone must enumerate the actual transitive primitive closure and pin that exact list. Do not blanket-allow all x/crypto or a third Totipo library. Standard library is shareable. A new common external semantic package fails closed. Assert intersection of nonstandard closures contains only these approved primitives, with no repository-owned shared package.
- Oracle test package has a stricter rule: **stdlib only**, no other local production package, no production constants. Its own minimal duplicate-safe artifact reader is local to oracle tests. Guard all source files, not only current compiled tests.
- Source-policy scan rejects go:linkname, go:generate/shared-source embedding, cgo/assembly bridges in repository implementation, and production os/exec/plugin imports; only the guard and explicitly scoped CLI test harness may spawn inspection/test processes. Make/Nix/CI command review must ensure no verifier path delegates to generator. A Go import graph cannot prove runtime execution or template provenance; these small additional source constraints and workflow review are necessary.

Normal graph inspection honors existing vendoring, offline flags and caches; do not force GOWORK=off inside current workspace-vendored Nix layout. Supplemental standalone-module Go tests/listing with GOWORK=off and -mod=readonly use ordinary module dependencies/cache and catch missing declarations. Keep guards compatible with Go 1.23; avoid a dependency just to parse module files. Exact dependency-list and path policy will have small synthetic negative tests (cross-import, local replacement, third semantic module, tagged/test import, oracle production import). No output regeneration, network-dependent guard behavior or broad source-copy detector is needed.

This is a modest import/metadata/source policy, not a sandbox against a malicious author editing the guard itself. Review remains required for code-generation provenance and changed enforcement rules.

## 13. Crisp shared-code policy

**No shared repository Go module/package at all between generator and verifier.** Share documented serialized contracts and standard/library primitives, not executable Totipo logic. This makes the guard simpler than maintaining semantic exceptions inside shared utilities.

| Candidate | Assessment | Application |
| --- | --- | --- |
| JSON structs | BETTER_DUPLICATED | Plain structs alone could be safe, but embedding Object/Node or custom unmarshal introduces behavior; duplicate small local models |
| Manifest/profile structs | BETTER_DUPLICATED | Independent writer/read acceptance; profile expectation never sourced from writer |
| Filesystem helpers | SAFE_TO_SHARE in abstract; BETTER_DUPLICATED here | Use stdlib; containment/canonical path validation is acceptance logic and FORBIDDEN_TO_SHARE |
| Hex parsing | BETTER_DUPLICATED | encoding/hex primitive safe; canonical lowercase/rejection policy independently owned |
| SHA helpers | SAFE_TO_SHARE primitive; BETTER_DUPLICATED wrapper | crypto/sha256 safe; domain/representation/OBJECT_ID composition FORBIDDEN_TO_SHARE |
| Test helpers | BETTER_DUPLICATED | No common expected-value builder or oracle fixture normalizer; semantic test helpers FORBIDDEN_TO_SHARE |
| Schema constants | FORBIDDEN_TO_SHARE as generated Go acceptance/semantic constants | Both may read authored schema files; reader interprets them; no shared generated decoder or constants package |

Trivial errors, file writing and hex wrappers do not justify a third module. Duplication here buys clear policy and avoids future helper growth becoming hidden semantic reuse.

## 14. Serialized-artifact interface and generator authority

```text
generator's own construction implementation
   -> files in temporary deliberate output tree / reviewed committed artifacts
   -> verifier opens files independently from disk
   -> integrity, schema and semantic comparison
```

No Go API crosses the boundary. Committed verification must work without generator execution or generated build prerequisites. Shared inputs (normative spec, authored schemas, preserved known answers) are evidence sources, not a shared semantic implementation.

Current violations: generator main imports internal/vectors structs/Decode/Hash; object.Object.Encode (which calls Dispatch); tlv framing helpers; cryptov1 derivation/sealing/padding/ID helpers; graph.Node/Result. These are all in `conformance/cmd/generate-vectors/main.go`; vectors/types.go embeds object/graph types, and vectors/runner.go joins data with evaluator dependencies. Go module itself permits the cross-role internal imports. Runtime orchestration already follows files and does not invoke generator, so no Make/Nix/CI generation API/subprocess violation was found. Consumer-side integration tests constructing bytes with consumer writers are legitimate verifier self-tests, not independent cross-generator evidence.

Generator may own fixture construction, canonical writer, malformed fixture encryption, deterministic serialization, explicit graph/fold/storage declarations, intermediate outputs, manifest and currently generated profile updates. It must not own acceptance authority or rewrite expected outcomes from verifier results.

Three provenance categories must survive without changing JSON schemas: (1) computed intermediates such as padded bytes/keys/IDs; (2) independently authored expectations such as graph heads/storage result/negative classifications; (3) preserved externally sourced RFC TOTP known answers and existing bootstrap source known answers. Do not overclaim bootstrap external provenance from source retention alone. Record provenance in source/tests/docs rather than generating every expectation anew. Generator must continue preserving RFC case bytes and existing bootstrap password/root/salt/nonce/wrap-key/header/record fields.

## 15. Structural-spec migration and verifier authority

Verifier independently owns artifact integrity, exact manifest/profile set, shapes, TLV parsing, crypto opening, object-ID recomputation, graph/fold/storage/TOTP evaluation and expected comparison. Repository pin validation is already its responsibility. Python schema checks fit naturally into vectors.Read before typed acceptance. Prose governance should use a **separate tiny `check-spec` command in the same conformance module**, without semantic-package imports, rather than polluting Run. No third policy module is needed; Make remains orchestration, not regex/JSON logic.

Complete check_spec rule categories and destinations (all retained initially):

| Individual rule family | Category | Destination |
| --- | --- | --- |
| Split history at section 21; exact sections 1–21; active Section references 1–20 | Markdown structure | check-spec policy command |
| Header protocol 1/revision r19 | Protocol metadata | check-spec with independently authored expectations |
| Exact archive link; r1–r15 archive and r16–r19 current headings | Revision history | check-spec; only spec/archive inputs |
| Token-boundary exclusion list DEVICE through K_signature_context; lowercase provenance/opaque exclusions | Repository governance | check-spec; retain full current list, not only sample terms |
| Sole TOKEN, MAX_PARENTS/count bound, 1005/1006/1024 | Protocol metadata + prose guard | check-spec literal/registry assertions |
| VAULT_ID/causal-equivalence formula anchors | Protocol metadata prose guard | check-spec; not runtime formula authority |
| Create-only state, optional external synchronization, no exhaustive enumeration, no parent publication prerequisite | Repository governance/prose guard | check-spec |
| Exact ordered 12 tags/names 0x0001–0x000c | Protocol metadata registry | check-spec independent table |
| Layout begins vault; no uppercase VAULT entry/path claim | Markdown/storage metadata | check-spec |
| Fold table exact three fields/presence; TokenValue excludes metadata; heads retain all metadata/value | Protocol metadata/prose guard | check-spec |
| Entry type table vault regular/objects-v1 directory; symlink follow/interpret/traverse bans | Protocol metadata/prose guard | check-spec |
| No host-specific syscall prescriptions; optional hardening | Repository governance | check-spec |
| Represented metadata preservation; no persistent/advisory/cache/cross-run requirement | Repository governance/prose guard | check-spec |
| Informed authorship with unavailable history; CSPRNG new ID/empty parents | Repository governance/prose guard | check-spec |
| Empty-password explicit confirmation | Repository/application governance | check-spec |
| No silent creation; explicit recovery workflow; orphan preservation; reopen regular vault; NO-REPLACE; persistence/revalidation before success | Repository/store governance | check-spec |
| No secure-erasure claim for tombstone/delete | Repository/application governance | check-spec |
| No password change/rewrap/root reuse; INFORMATIVE/incomplete sync migration notes | Repository governance/prose guard | check-spec |
| Retired fingerprint/replacement terms | Repository governance | check-spec |
| CSPRNG 32-byte lifetime-immutable root; no rotation; bootstrap 0x01, never newer, alternate same-root representation anchor | Protocol metadata/prose guard | check-spec |
| No VAULT replacement even exact-identical | Repository governance/prose guard | check-spec |
| Core/Store-writer/Application scope table + metadata; sections 12/17 scope; GUI/corpus/application boundary; no stronger atomic primitive | Repository governance/prose guard | check-spec |
| Cross-file artifact pins / generated-artifact integrity | **No longer in check_spec** | Existing vectors.Read/VerifyProfile, retain mutations |

The full 115-line script is the migration rule inventory; table groups related assertions without authorizing deletion. Regex wording, newline sensitivity and section splits are brittle implementations, but this migration should preserve their behavior first. Any later prose-lint reduction requires separate review. Shell is unsuitable for these rules. Python need not remain for regex availability; Go regex/strings suffice. check-spec tests should preserve success with only the two documented inputs and all eight mutation failures. Retain active/history separation; historical text must not enter retired-term rejection.

## 16. Vector schema and exact presence migration

Go currently rejects duplicates/trailing JSON/unknown typed fields, validates operation payload count/prerequisites, unsafe paths, hashes, identities/profile relationships, and executes semantics. It does **not** interpret either JSON Schema. Missing scalars, case-insensitive struct-key matches, nulls mapping to zero/nil, unconstrained strings/arrays, and operation-nested required lists are not equivalent to schema acceptance. Required false/zero and empty arrays are especially dangerous.

Schema responsibilities lacking comprehensive Go coverage: exact key presence and case at every nesting; type/null distinction including optional present metadata; required graph head_objects even empty; conflicting/diagnostics false; five workflow booleans; TOTP t0=0 and row unix time/counter zero; required issuer/account empty strings; all pattern/length/enum/numeric/array constraints; all nine operation conditionals and action-specific add/remove/evaluate contracts; supported-keyword rejection; schema source validation with duplicate-safe reading. Existing Go incidental rejection of some malformed values is not parity.

Required-field inventory:

- Manifest: format/protocol/spec_revision/cases; entries id/category/kind/normative/path/expected/spec_sections/sha256.
- Case: format/id/operation/expected; crypto additionally input/semantic_hex/root_hex/crypto, dispatch semantic_hex/root_hex/crypto, bootstrap root_hex/bootstrap, totp/graph/fold respective payload, storage root_hex/storage, workflow workflow, post-aead semantic_hex/root_hex/post_aead.
- Input and fold token: identity/parents/status/issuer/account/algorithm/digits/period/secret. Optional client_name/time must have valid exact types when present.
- Crypto: object_id/id_key_hex/object_root_key_hex/object_key_hex/nonce_hex/aad_hex/padded_plaintext_hex/ciphertext_hex/gcm_tag_hex/object_hex/semantic_length. Bootstrap: password_hex/salt_hex/nonce_hex/wrap_key_hex/header_hex/record_hex/vault_id_hex.
- TOTP: source/notes/algorithm/digits/period/t0/secret_hex/rows; each row unix_time_seconds/counter/counter_hex/code.
- Graph: steps; step action; add node, remove id, evaluate identity/expect. Node/head records id/identity/parents/canonical/value. Expect heads/values/unresolved/conflicting/head_objects.
- Fold: frontier/stage_ids/parents/token. Storage: namespace_kind/entries/classes/diagnostics; entry path/kind. Workflow: action/kind/existing_hex/intended_hex/readable/complete/durable/orphan_objects/parents_available/result. PostAEAD: object_id/encryption_plaintext_hex/object_hex/defect.

Recommended design: **two-pass structural validation**. First duplicate-safe parse with json.Decoder.UseNumber into generic exact-key maps/arrays (or raw messages preserving explicit null); apply the checked-in schema subset before any typed struct decoding. Then existing strict typed Decode + ValidateShape + semantic execution. Required means map membership, never nonzero value. Do not attempt to fix hundreds of fields only by making bool pointers: pointers still conflate absent/null and do not enforce all schema restrictions.

Implement the existing closed keyword set: $schema/$id/title/description annotations; type, const, enum, properties, required, additionalProperties, items, minItems/maxItems, minimum/maximum, pattern, minLength/maxLength, allOf, if/then/else. Integer comparisons must be exact through u64max, exclude bool, preserve current Python distinction of integer tokens from float/exponent spellings if parity requires it; no float64 intermediate. Stdlib math/big is available. String lengths count Unicode code points as Python does, not UTF-8 bytes; wire text byte bounds remain separate semantic checks. Current patterns use anchors, character classes, repetition and noncapturing groups; test Go/Python matching parity explicitly, including Python end-anchor behavior before a final newline. Inspect regex dialect compatibility rather than assuming future schemas are portable; unsupported regex/schema extensions fail closed. Preserve current const/enum behavior for actual contracts; differential tests should expose Python bool/int equality quirks instead of silently changing acceptance. Reject schema duplicate members and unsupported visited keywords; a complete schema-tree preflight can strengthen rejection of unknown unvisited keywords, but any acceptance change must be explicitly recorded.

Parity evidence before removal: all ten existing schema methods; remove every required field from representative operation/nested records including false/zero/empty arrays; explicit null, wrong-case keys, bool-as-integer, negative/fraction/overflow; min/max/pattern/enum boundaries; duplicate nested members; mixed/unknown payloads; safe-root alias and escaping symlinks; current manifest/all 92 cases. Differential old/new accepted/rejected results must match or document deliberate verifier hardening separately. No corpus/schema/pin edits to obtain parity. Add actual CLI integration coverage; Go struct success alone is never a deletion gate.

## 17. Oracle migration options and formula placement

Current verifier cannot absorb Python merely by calling its existing packages: its computations are identical to generator's. Two independently implemented sides after the split supply useful A/B evidence, but keeping the small third selected formula implementation preserves current oracle evidence at low cost.

| Formula/responsibility | Generator A | Verifier B | Independent tests C |
| --- | --- | --- | --- |
| TLV/canonical writer and field widths | Own writer/negative constructor | Own parser/canonical writer for consumer checks | Direct binary field walk + independently assembled expected fields |
| HKDF domains, ID, object key, nonce/AAD/padding | Own composition from spec | Existing independently maintained composition | Direct HMAC extract/single-block expand with own literals |
| VAULT_ID | Own SHA over source record | Own representation validation/hash | Direct SHA over record bytes |
| Size/capacity | Own writer sizing | Own accepting bounds | Preserve 1005/1006/1008/1024 assertions; add explicit independently summed maximum field sizes |
| Post-AEAD defect inputs | Own malformed constructor | Raw AEAD verification + strict rejection | Independent plaintext/ID/length predicates, including n=1007 |
| TOTP | No formula required: preserve RFC files | Existing Code/counter validation | External RFC answers remain; no new duplicate TOTP implementation needed for migration |
| Graph heads/state | Explicit authored outcomes/record attachment | Evaluate | Keep existing independent matrix-closure test |
| Fold/staging | Own expectation planner | Fold/FoldObjects | Explicit expected lists and encrypted integration evidence |
| Argon2/bootstrap wrap | Preserve known-answer input fields | Existing Wrap/Unwrap/WrapKey | No Python evidence to port; do not claim new KDF independence |

Place C in `conformance/internal/wireoracle/` as independent tests with their own tiny strict reader and no imports of conformance production or generator. The import guard is essential; a directory name alone is cosmetic. Raw maps avoid type/model normalizations masking fields. Keep Python and C running in parallel through parity and mutation review before deletion. Test mutations must demonstrate a wrong domain/tag/width/nonce/padding/ID or expectation is independently rejected; checking a copied formula against itself is insufficient. No third command/module is needed unless later external interoperability or larger independent evaluator scope justifies it.

## 18. Crypto independence line

| Primitive/operation | Current overlap | Allowed boundary |
| --- | --- | --- |
| SHA-256 | stdlib; shared MAC/Hash/VaultID wrappers; Python hashlib | Same primitive allowed; Totipo domain/bytes selection independently specified |
| HKDF | x/crypto/hkdf inside common Derive/expand; Python direct HMAC block | Primitive package allowed; extract salt, labels, output lengths and composition cannot be shared |
| AES-GCM | stdlib in common crypto; raw generator/post_aead calls | Same aes/cipher allowed; envelope sizing, key/nonce/AAD/header construction/open policy independent |
| Argon2id | x/crypto/argon2, fixed parameters in verifier; generator currently preserves bootstrap answers | Primitive allowed; password rules, parameters, salt slicing, wrapping composition owned by verifier and any future writer independently |
| P-256 signature verification | **Absent**; retired from active protocol | No migration or new dependency; do not confuse historical reports with live crypto |
| DER parsing | **Absent** | No replacement needed |
| TOTP HMAC | stdlib HMAC SHA1/256/512 in totp.Code; generator preserves RFC cases | Primitives allowed; counter/truncation/modulus/format/credential rules must not be common writer-reader implementation if generator later computes codes |

Independence is Totipo composition and acceptance, not reimplementing AES/Argon2/hash primitives. This audit does not establish independent full cryptographic interoperability. Existing Python covers no AEAD/Argon2; retaining C does not change that limitation.

## 19. Constant independence

Current shared protocol literals/constants include object.Capacity=1006, MaxParents=4, Supported/Invalid; tag numbers/order/widths and bounds in Object.Encode/Dispatch; domains in Derive/ObjectKey/AAD; object/AEAD sizes 1008/1024, nonce 12, ID/root 32; VAULT magic TOTIPO-VLT plus version 0x01, length 87/header 39; Argon2 parameters 3/65536/4/32; status/algorithm enums, digits 6–8, secret/text limits. Graph/fold limits and canonicalization rules also belong to semantics. Some generator values already use independent literals (maximum 1005, fold 4/3, explicit invalid tags/values); Python has its own domains/tags/sizes. TOTP algorithm constants reside only in verifier computation and external case data, not a generator evaluator.

Duplicate wrong-shared-value-sensitive protocol constants in A/B/C where each implements that rule. Do not export one constants module or generate constants from the writer/schema for all implementations. Preserve explicit reader expected revision r19 independently of incoming manifest/writer. Revision/version metadata in manifest/profile/schema/spec is compared, not accepted as a new supported revision automatically. Mechanical descriptive revision labels (Make help/Nix pname/printed summaries) can reuse a local validated expectation within one role; they need no cross-module semantic constants package. Independent literals are useful evidence, not waste to optimize away.

## 20. Test topology and generation policy

Current static count: **22 Go files, 3,379 LOC, 37 Test functions and 3 Fuzz functions**; two command packages without tests, seven tested internal packages, ten Go test files. Python: **15 methods** (2 structure, 10 schema, 3 wire). These are declarations, not freshly run results.

| Test family | Current locations | Intended topology |
| --- | --- | --- |
| Generator self-tests | None directly; command compiles | Generator module local writer/serialization tests and deliberate temp regenerate/compare evidence |
| Verifier self-tests | tlv/object/crypto/graph/storage/totp tests, strict decode and path tests | Stay conformance-local; all-package test/vet/race coverage preserved |
| Cross-artifact tests | vectors.TestCorpus, integrity/path mutations, consumer CLI; real-fold integration is consumer-local roundtrip | Committed artifacts → independent reader; explicit generated temporary artifacts → independent verifier only during generator qualification |
| Independent oracle | Python wire tests; graph closure; RFC known answers | Guarded C tests + existing closure + immutable known answers; no generator prerequisite |
| Repository/governance | Python structure + schema tests; Go integrity tests | policy command tests separate from semantic Run; schemas in vectors acceptance; guard tests |

Keep full/verify-only CLI smoke behavior until dedicated CLI tests prove argument errors, load/profile failure, full evaluator failure and verify-only skipping evaluation. Do not remove duplicate invocations as part of a module migration. Preserve all graph/storage modeling limits and metadata snapshots.

Generation remains a human deliberate maintenance action. No automatic generation from make check/test/verify/conformance, Nix or CI. Generator changes require a temporary complete fixture tree seeded with retained known answers, explicit regenerate-and-compare against committed bytes, reviewed intended differences, then ordinary verification of committed artifacts. Committed-vector verification never requires running generator first. For a boundary-only phase there must be **zero** intended generated differences.

Proposed contributor targets: keep make check/spec-check/test/verify/conformance/race/fuzz; optionally add make generate solely as explicit wrapper. test/vet/race cover both modules; fuzz retains the exact three consumer targets/flags/seeds. spec-check runs conformance/cmd/check-spec; verify runs schema-aware conformance CLI -verify-only; conformance runs full CLI; check includes architecture guard + existing responsibilities. generate invokes generator/cmd/generate-vectors. Do not implement invocation deduplication during this work. Module boundaries remain hidden behind familiar Make UX where possible.

## 21. go.work recommendation

Retain and extend existing go.work to the two modules. **go.work is convenience, NOT isolation.** It preserves root go run commands and existing workspace-vendor strategy. Guards must inspect active workspace resolution **and** all declared source imports/module metadata, so accidental public cross-imports cannot hide behind workspace membership. Standalone GOWORK=off checks in portability/developer qualification catch undeclared dependency reliance; they supplement workspace checks and must account for different vendoring lookup.

Avoid nested module path ancestry, extra workspace members and replace directives. Workspace dependency selection can align primitive versions during joint builds but does not enforce equal module declarations; keep declarations explicit/pinned. No root go.mod is needed. Nix source must include both module trees and go.work; root workspace sums, if actually generated/required, must be reviewed rather than silently added.

## 22. Dependencies, Nix and CI impact

Current dependencies: conformance module Go 1.23.0; direct x/crypto v0.36.0, indirect x/sys v0.31.0; four checksum lines. hkdf/argon2 are the live external primitive packages. Generator can implement single-block HKDF using stdlib HMAC as Python does and preserve bootstrap inputs, so a stdlib-only generator is feasible. Alternatively duplicating current pinned primitive module declarations is acceptable if genuinely used; do not copy unused dependencies or upgrade them. No JSON-schema, CLI, assertion or crypto library is needed for migration. Exact Go LOC/dependency closure must be measured during implementation.

Required future changes:

- flake qualificationSource must include generator .go/go.mod/go.sum, conformance policy/oracle/guard Go, and retained Python until deletion. Keep archived pre-r16 history; exclude unrelated reports. Add formatting and tests/vet/race for generator. Make changes alone do not add the missing generator source to the filter.
- Preserve current workspace `go work vendor` fixed dependency derivation, offline GOPROXY/GOSUMDB, GOFLAGS=-mod=vendor, local toolchain, temporary build/mod caches, compiler/CGO race environment. Module split can change vendor/modules.txt even if primitive versions do not; current vendorHash `sha256-637XdgBok7hRhMw5s9hdY3Mr5jlPXqc1viSYoYh9Iqw=` must be requalified, not assumed valid. Compute/update only if vendor bytes differ through the human workflow. One workspace hash suffices; two separate module hashes are unnecessary for option C.
- Include module count/import guard in filtered-source qualification; it must not need .git. Resolve root from explicit -root rather than brittle CWD assumptions. Workspace vendor vs standalone-module mode requires explicit CI configuration; root vendor cannot simply serve an arbitrary GOWORK=off module command.
- CI retains authoritative Linux `nix flake check --print-build-logs` and all five supplemental combinations. Add generator formatting/vet/test and module cache-dependency paths (only sums that exist), plus guard. Keep current Python setup/checks until final deletion. Dependabot adds generator gomod entry only if it has managed external dependencies. No lock/action/toolchain upgrade is required.
- Go cache keys naturally incorporate different package/source paths; source moves invalidate relevant cached builds. Setup-go dependency cache should include both modules' actual sums; this is no qualification evidence by itself. No generator run in CI normal checks.
- After all Python deletion gates, remove runtime from qualification nativeBuildInputs, portable setup and directly needed development/jail package lists; remove tools .py filter only once no active Python remains. Optional authoring tools are otherwise out of scope. No Nix apps/packages need adding.

Every implementation phase changing included Go/Python/Make/workspace/filter input invalidates its prior Nix gate. New files, including reports, must be Git-known via intent-to-add with cached diff empty before frozen human `nix flake check`; agent never runs Nix. Linux host qualification cannot claim remote portability. No dependency updates solely for this migration.

## 23. Minimal safe migration phases

Three reviewable phases are sufficient; a separate fourth cleanup phase is optional only if deletion evidence is not ready.

**W3A — establish real boundary, keep Python.** Move generator command to sibling module, replace all conformance imports with generator-owned structs/writer/composition from spec, preserve authored results/known answers, establish internal visibility and import/metadata/source guard. Keep verifier and all Python checks active. Update workspace/filter/orchestration/CI for both modules without eliminating coverage. Require temporary generation comparison byte-for-byte, independent existing Python/Go verification and focused generator mutations. A mechanical extraction alone must be labeled as shared ancestry until reviewed independent derivation is established; do not call it completion of independence.

**W3B — structural/schema replacement.** Port check-spec and its two-input CLI tests separately from semantic reader. Add exact-presence two-pass schema validation and all negatives/CLI tests. Run Python/Go differential acceptance before deleting check_spec/test_check_spec/check_vectors/test_check_vectors. Keep wire Python oracle and a local minimal reader it needs; easiest keep check_vectors.py until W3C rather than leave a hidden import failure. In that case remove its Make schema CLI only after parity, but postpone deleting its file. Freeze each approved responsibility before full final qualification.

**W3C — independent oracle and justified Python removal.** Port selected independent formulas to guarded stdlib-only C tests; prove parity and mutation sensitivity against retained Python. Delete last Python files only after their imports/callers are gone. Remove Python normal-tool/runtime/CI setup/filter references, preserve contributor targets and complete qualification. No protocol/corpus/schema/profile edit is part of any phase. If B parity or C independence fails, keep Python; zero-Python is not a deadline.

Each phase is a separately authorized implementation milestone: baseline make check once as applicable, focused loop, final full check/vet/race/three fuzz once after stabilization, human Nix last, remote CI after user commit/push. Generator comparison is deliberate evidence, never a replacement for normal verification.

## 24. Exact eventual Python deletion gates

| File | Replacement | Parity/independence evidence | Deletion gate |
| --- | --- | --- | --- |
| tools/check_spec.py | conformance/cmd/check-spec + local policy code | All current assertions; same two-input success; all mutation failures; no semantic imports | Policy CLI active in Make/CI; tests pass; no external scope assumed |
| tools/test_check_spec.py | policy Go tests | Both methods, eight mutations, argument/exit behavior as applicable | Replacement regression evidence retained; checker no longer Python |
| tools/check_vectors.py | schema-aware Go vectors.Read/verify + oracle's own reader | Exact schema/presence/type/path differential; duplicate-safe schema/case reading; no generated contract | No Python wire test imports remain; all CLI callers migrated; retained until C if needed |
| tools/test_check_vectors.py | Go schema contract/mutation/CLI tests | All ten methods plus exhaustive required zero/false and null/case/range tests | Union of old schema and Go protections established; no corpus changes to pass |
| tools/test_r16_vectors.py | independent stdlib-only wireoracle tests | All three methods/formulas/constants, post-AEAD predicates; independent-reader/import guard; wrong-intermediate mutation rejection | A/B boundary verified; C parity reviewed; no production helper reuse; existing AEAD/Argon2 limits documented |

At B some files may remain dormant temporarily for C's comparison/import needs; record actual deletions/LOC and do not claim Python removed early. Historical audits already removed in Wave 1 are not further deletion targets; preserve every historical report.

## 25. Measurable success metrics

| Metric | Current | Successful Wave 3 target |
| --- | --- | --- |
| Executable implementation languages | Go + Python | Go; Nix/Make/YAML/shell orchestration retained |
| Active Python | 5 files / 486 LOC | 0 only after all five deletion gates |
| Go | 22 files / 3,379 LOC | Measure actual added/moved/deleted LOC per phase; no net-LOC promise |
| Planning size budget | Existing generator 508 LOC; Python 486 | Rough incremental Go 800–1,400 LOC for independent generator internals, policy/schema/guard/oracle/tests; estimate, not measured result or fixed cap |
| Modules / workspace | 1 / 1 | 2 / existing workspace; no shared module |
| Generator→consumer direct repo package edges | 5 imports | 0; transitive opposite-module edges 0 in source/tests/build closures |
| Shared repository semantic packages | object/tlv/crypto; structural graph/vectors; incidental storage/totp closure | 0; common external packages limited to reviewed primitives |
| Independent selected wire implementations | shared Go A/B + Python C | independent generator A, consumer B, guarded oracle C; no shared composition/constants |
| Graph/fold/TOTP expectation sources | authored outcomes, independent closure, fixed RFC cases | Preserved; no evaluator-produced expected answers |
| Ordinary commands | 8 Make targets, two Go commands, Python checks | Same checking UX, optional explicit generate wrapper; policy/guard commands internal plumbing |
| External dependencies | x/crypto 0.36.0 + x/sys 0.31.0 | No new external dependency; generator may use stdlib only |
| Artifact interface | Files plus shared Go implementation/model APIs | Files only; no cross-role API/subprocess/shared-generated-source edge |
| Qualification coverage | 92 cases; all-package race; three bounded fuzz targets | Same cases/seeds/flags; extend all-package checks across modules; no automatic regeneration |
| Artifact bytes | Committed spec/schemas/cases/manifest/profile | Byte-identical throughout tooling migration |

LOC growth is acceptable only where it purchases explicit independence or retained validation. Report exact totals, not earlier audit estimates (the original 493 active Python LOC predates Wave 2; current is 486). Zero Python without presence coverage or independent computation is failure.

## 26. Final recommendation and proposed tree

**A. TWO_MODULE_SPLIT_RECOMMENDED**, option C:

```text
go.work                         # use ./generator and ./conformance
generator/
  go.mod                        # module totipo/generator; Go 1.23-compatible
  cmd/generate-vectors/main.go
  internal/fixture/             # writer-owned models, declared expectations, JSON/pins
  internal/wire/                # independently implemented writer/TLV/crypto composition
conformance/
  go.mod                        # retain totipo/conformance and pinned primitives
  go.sum
  cmd/totipo-conformance/
  cmd/check-spec/                # prose/history policy only
  cmd/check-architecture/        # small stdlib guard
  internal/{tlv,object,cryptov1,graph,storage,totp,vectors}/
  internal/repopolicy/           # only if command size justifies package extraction
  internal/wireoracle/          # independent tests; stdlib-only imports
spec/ vectors/ requirements/    # unchanged serialized/normative interface
```

Tree is a responsibility sketch, not a requirement to create one package per formula. Keep generator internals small; no graph evaluator or TOTP calculator is needed there today. No root module or third shared mechanical module. Optional generator go.sum exists only if external dependencies actually used.

What forces structural separation: sibling internal visibility, exhaustive import/dependency/metadata policy, frozen allowed primitive dependencies, independent oracle restrictions, both module builds, serialized output tests. What documents intent: tree names, go.mod count, README and this contract. What establishes computation independence beyond structure: separately derived composition/constants, authored expectation provenance, differential/mutation evidence and reviewed code ancestry. Static imports cannot establish that alone.

## 27. Next implementation milestone preview (not executed)

Scope only W3A after explicit authorization: new generator go.mod/command/internal source/tests, removal of old generator command, small conformance architecture guard/tests, go.work, Make orchestration, flake.nix source/check coverage, portability cache/test/format commands, current contributor docs and new excluded implementation report. No Python removal yet; no consumer semantic refactor, schema rewrite, dependency upgrade or protocol change.

First establish the import guard/internal boundary and independent generator writer/composition; do not export consumer helpers or extract a shared library. Preserve spec/, vectors/, requirements/ **including all raw hashes, field presence and serialized bytes**, all historical reports, existing verifier behavior, Python oracle, primitive versions, race/fuzz coverage and matrix combinations. Source movement necessarily changes generator Go bytes; it does not permit artifact differences.

Checks: starting clean main/HEAD/revision; baseline make check once; focused generator/guard tests; deliberate temp regeneration seeded with preserved known answers and exact comparison; normal independent corpus/schema/oracle verification; final make check, vet in both modules, all-package race in both, exact three bounded consumer fuzz targets, formatting and diff checks. Preserve failure output/reproduction inputs and narrow diagnosis. New files intent-to-add, require no staged content, freeze all included sources, then **human `nix flake check`**; no agent Nix. Final remote Linux/macOS/Windows matrix only after user's commit/push, never inferred from host results.

## 28. Validation, execution counts and final Git state

This audit's only authorized output is `review/SPEC_SUPPORT_TOOLING_WAVE3_ARCHITECTURE_AUDIT.md`. Actual flake filter selects spec/vectors/requirements, conformance Go/module files, tools .py, Makefile, go.work and only review/V1_PRE_R16_REVISION_HISTORY.md. Structure tools consume spec/archive; schema/oracle tools consume vectors. None consumes this report. No fresh human Nix request is needed for this excluded report.

Execution accounting: baseline make check **0**; focused Python/Go tests **0/0**; focused fuzz Dispatch/Open/graph **0/0/0**; final complete suites/incomplete attempts **0/0**; final check/vet/race/fuzz **0/0/0/0**; generator **0**; human/agent Nix **0/0** for this audit; remote CI **0**; cached Nix observations **none**. Python used for static inventory and byte comparison is not checker/test/oracle execution. Official-document browsing is not repository qualification.

Final validation: raw bytes and path inventories compared against starting HEAD for all **265 pre-existing tracked files**; protected scopes spec/, vectors/, requirements/, conformance/, tools/, go.work and every go.mod/go.sum are unchanged, with no additions/deletions. Protected scope counts are spec 1, vectors 97, requirements 2, conformance 25 and tools 5 files. Makefile, flake.nix, flake.lock, .github inputs, process docs and all previous review/history evidence are also unchanged. Root go.mod/go.sum/go.work.sum remain absent. Nonignored path additions are allowlisted to this report only. `git diff --check` and `git diff --cached --quiet` pass. Report Markdown local links/fences were checked.

Final branch/HEAD remain main / `1f78b574f0855a515dfadaccd7975e513bd6c962`. The report is Git-known via the required `git add --intent-to-add review/SPEC_SUPPORT_TOOLING_WAVE3_ARCHITECTURE_AUDIT.md`, with **no staged content**, and remains uncommitted:

```text
 A review/SPEC_SUPPORT_TOOLING_WAVE3_ARCHITECTURE_AUDIT.md
```

No other repository edit, implementation, generation, qualification execution, commit, push, tag or workflow dispatch occurred. Implementation awaits a separate authorized milestone.
