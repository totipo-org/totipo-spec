# v1/r19 immutable VAULT report

**Status:** r19 changes ready for review; ordinary validation passed and the user
reports `nix flake check` passed. Everything is unstaged and uncommitted.

## 1. Starting state

Repository: `totipo-spec`, branch `main`, clean worktree, HEAD
`fbae625f2cad6c89a1b5555de2f2d75d890bb682` (`design resources`). HEAD contains design
commits after the r18 integrity repair (`4623a7e`); the current committed spec,
manifest, moving profile, and README consistently identified v1/r18. The repaired
r18 specification and manifest hashes matched the integrity-repair record. There
was no supplied baseline commit pin beyond committed r18; no checkout/reset occurred.

Starting corpus: **90 cases**. Initial `git status --short` was empty.

| Starting artifact | SHA-256 |
| --- | --- |
| Specification | `8a357e75f3ddd92efa954fde2ffc9af33f40a2c2396d6afbf1d5de5f00bc4f8a` |
| Manifest | `bd2b52adc05b26e09790f5f7367761b2b86ba3cf8d97d10ed187fcf7213fcf02` |
| Moving requirements/profile | `4c7954cd2b59aa0afbbe3c22080cadddf72135d884b186a671266c29d58418df` |

## 2. Baseline validation

Before editing: `make check`, `make verify`, `make conformance`, `make race`,
`make fuzz`, `gofmt -l conformance`, and `go -C conformance vet ./...` all passed.
Baseline Python suite: 13 tests; corpus: 90/90. Fuzz used the Makefile's three
10-second targets, each with two workers. Formatting listed no files.

Read the actual r18 normative specification, Makefile, documented commands,
profile, manifest, both schemas, generator, runner, storage model/tests, crypto
model/tests, current revision checker and integrity tests. A private line inventory
and copies/hashes of all baseline cases and integrity artifacts were captured in
`/tmp/totipo-r19-baseline/` before editing. Historical r15–r18 reports were context,
not a substitute for current text. No baseline failure was bypassed.

## 3. Decision summary

One vault has one immutable canonical VAULT, one immutable random K_root, and one
stable public VAULT_ID. Password/KDF/root changes create an independent vault with
fresh root, salt, and nonce. There is no in-place password change, root rotation,
rewrap, or VAULT replacement. Protocol state is create-only. Migration and transfer
workflows are application choices, not serialized protocol state.

No protocol or BOOTSTRAP_VERSION bump, marker, compatibility bit, legacy reader
branch, dependency update, RC freeze, or release was introduced.

## 4. Exact VAULT byte-format invariance

Section 6's encoding and root-wrap construction code blocks are byte-identical to
r18. Canonical pathname remains exact lowercase `vault`; VAULT names the bytes.

| Offset | Bytes | Field |
| --- | --- | --- |
| 0 | 10 | ASCII `TOTIPO-VLT` |
| 10 | 1 | BOOTSTRAP_VERSION `0x01` |
| 11 | 16 | ARGON2_SALT |
| 27 | 12 | WRAP_NONCE |
| 39 | 32 | WRAPPED_ROOT |
| 71 | 16 | WRAP_TAG |

Total remains **87 bytes**. Argon2id version/parameters, strict UTF-8 domain and
1024-byte maximum, empty-password format semantics, K_wrap construction,
AES-256-GCM wrapping/AAD/tag, pre-KDF checks, and root authentication are unchanged.
The existing wrong-password/tag-failure behavior remains the same `ErrCrypto`
path; no new password oracle or compatibility branch was added. `WrapKey`, `Wrap`,
and `Unwrap` bodies are byte-identical to HEAD. Existing bootstrap crypto bytes
were retained as source known answers rather than recomputed.

## 5. K_root invariance

Creation MUST generate `K_root = CSPRNG(32 bytes)`. The root is permanent for the
vault lifetime and MUST NOT rotate in-place. A different root always denotes a
different cryptographic vault. Migration must generate an independently fresh
root; keeping the old root across a credential change is forbidden.

The HKDF hierarchy, HMAC object addressing, object keys, deterministic nonce/AAD,
and object encryption functions are unchanged. Public fixed roots in test vectors
are fixture inputs, not authority to reuse roots in live creation.

## 6. Immutable VAULT semantics

Section 6 prominently prohibits intentional modification, replacement, rewriting,
rewrapping, or ordinary removal after initial publication. A different canonical
representation is another vault or invalid/substituted evidence, never a newer
version of one vault. No legitimate same-root alternate representation, wrapper
head, bootstrap lineage, or VAULT generation exists.

## 7. Removal of password rewrap

Former normative Section 8's same-root rewrap operation, old/current exact-byte
comparison, replacement durability/races, ambiguous replacement winner, and
historical-wrapper handling are removed. The replacement API and its runner
transition are deleted. The schemas reject `replace`, `rewrap`, and the former
`base_hex` field. A schema regression test and the specification checker enforce
that boundary. No product migration workflow is required by conformance.

## 8. Initial no-replace creation

Section 7 requires absence; an existing exact or different canonical VAULT always
blocks creation without mutation. Wrong type is not absence. Creation constructs
complete bytes separately, validates locally, uses the strongest reasonable
crash-safe NO-REPLACE publication, attempts file/storage and namespace persistence,
then reopens, compares exact canonical bytes, authenticates the wrap, and verifies
the recovered root before success. No particular syscall is prescribed.

Actual r18 had no mandatory post-publication establishment/reread; r19 adds the
requested explicit revalidation while retaining its platform-neutral durability
intent. Failure/ambiguity must not claim success or absence and creates no pending
protocol journal. Later open interprets canonical bytes actually present; no retry
gains replacement authority.

## 9. Create-only store invariant

Canonical VAULT and `objects-v1/<OBJECT_ID>` entries are create-once and immutable.
After initialization no protocol-visible durable entry is intentionally replaced
in-place. Existing object publication behavior, exact-existing read-only success,
size/capacity, namespace, and failure rules remain unchanged. External actors may
still change/delete storage; the protocol cannot prevent that.

## 10. VAULT_ID definition

`VAULT_ID = SHA-256(exact canonical VAULT representation)`: exactly 87 input bytes,
32 raw output bytes, ordinary text lowercase hexadecimal. It is public, non-secret,
and stable because VAULT cannot change. Different bytes imply different identity
except with negligible SHA-256 collision probability. The reference helper checks
canonical width/header and hashes the complete record; hashing is not authentication.

ASCII bootstrap known answer:
`8bef25c2a2b81a0691b4f719eabed9c001a1778af415ddccab380eba9b408351`.

## 11. VAULT_FINGERPRINT disposition

Removed from current normative text, reference crypto helper, bootstrap case
metadata, schema, runner, and current summaries. Actual r18 Section 9 used it for
same-root recognition stable across rewraps and explicitly excluded authorization,
freshness, first-contact authentication, and rollback protection. Inspection found
no independent security property requiring it once VAULT is immutable. Public
SHA-256 recognition replaces that role without unlocking K_root. Historical
revision entries and review/audit artifacts retain their historical terminology.

## 12. Recognition/application configuration

Applications may remember location, expected VAULT_ID, friendly name, and
preferences. This is disposable application configuration, not protocol authority.
For a retained binding the exact canonical VAULT must hash to the expected ID.
Mismatch means a different vault, not a newer/older wrapper; malformed evidence
remains invalid. Configuration loss does not invalidate the vault. VAULT_ID is not
authorization, an additional password verifier, first-contact authenticity,
storage freshness, object completeness, or object-bag rollback protection.

## 13. Removal of wrapper freshness/rollback semantics

No latest-password-wrapper, same-root stale-bootstrap rollback, concurrent
password-change conflict, replacement candidate/winner, or expected-old-wrapper
operation remains. Historical copies of the one immutable source VAULT remain
that source vault; they are not older bootstrap generations. Object-store freshness
limitations remain unchanged.

## 14. Remaining untrusted-store limitations

Section 2 preserves untrusted validation and describes hidden/deleted/substituted/
malformed VAULT, hidden/deleted/malformed objects, stale object subsets, and
incomplete observations. Crypto does not ensure availability or globally current
history. Removing multiple legitimate bootstrap generations does not add object
rollback protection. Independent storage/application mechanisms remain external.
The existing Sections 10–19, including observation, graph, and object publication
semantics, are byte-identical to r18.

## 15. Orphan-object creation safety

With VAULT absent, observed plausible current-family candidates (Section 3's direct
regular-file children with 64 lowercase hex names) block ordinary creation.
An explicit recovery/reconfiguration/new-location workflow is required. Candidates
are unauthenticated contextual evidence and cannot recover a root, but cannot be
declared garbage or deleted because VAULT is missing. Temporary absence and
incomplete/untrusted storage are disclosed. No exhaustive inventory or proof of
unseen-object absence is imposed. This intentionally replaces r18's permissive
creation-with-warning behavior.

## 16. Informative migration guidance

Section 8 normatively requires independent destination creation and leaves the
source VAULT/object bag/namespace unchanged. Section 8.1 is explicitly INFORMATIVE:
a separate destination store normally receives selected validated semantic state
or available history, re-authored and re-encrypted with fresh destination keys.
No transfer IDs, receipts, ancestry, MIGRATED_FROM, or source/destination identifier
fields are added to VAULT or TOKEN. No transfer wire protocol or schema is defined.

## 17. Migration completeness disclosure

The client can copy only observed and validated source data. It cannot prove other
devices have no unsynchronized changes, provider global currency, absence of an
entirely unseen token elsewhere, or full remote convergence. Informative SHOULD
recommendations disclose possible not-yet-synchronized omissions and retain the
old vault until the user verifies the new one; they are explicitly outside
conformance. History preservation, compaction, conflict handling, current-value
selection, and tombstones remain honestly described application policy.

## 18. Cross-vault-transfer note boundary

[V1_R19_CROSS_VAULT_TRANSFER_NOTES.md](V1_R19_CROSS_VAULT_TRANSFER_NOTES.md) is
prominently application/workflow design, outside v1 conformance and normative
protocol state. It records independent copy use, migration as creation plus copy,
fresh destination TOKEN_IDs, current-state/available-history choices, local
resumable receipts, secret minimization, local VAULT_ID job binding, retry identity
mapping, changing source snapshots, unseen-data limitations, and later recovery
from the old vault. It implements no transfer protocol, schema, or vectors.

## 19. Conformance cases removed

Five required case IDs are retired:

- `v1.bootstrap.rewrap.001`: lifecycle meaning retired; all static crypto inputs
  retained under `v1.bootstrap.known-answer-extra.001`.
- `v1.vault.replace-equal.001`: replacement transition removed.
- `v1.vault.replace-stale.001`: old/current wrapper comparison removed.
- `v1.vault.replace-unreadable.001`: replacement read prerequisite removed.
- `v1.vault.replace-ambiguous.001`: replacement outcome removed.

All retired physical case paths and profile/manifest requirements are removed.
Historical reports are unchanged.

## 20. Conformance cases added

Seven required case IDs are added; final corpus **92**:

| Added ID | Evidence |
| --- | --- |
| `v1.bootstrap.known-answer-extra.001` | Retained static crypto known answer, unlock, and VAULT_ID; no lifecycle/root reuse authority |
| `v1.bootstrap.vault-id-mutation.001` | SHA-256 known answer and one-bit tag mutation gives a different ID and fails authentication |
| `v1.vault.create-different.001` | Different canonical existing representation blocks creation unchanged |
| `v1.vault.create-revalidation-failed.001` | Durability alone cannot report success without canonical revalidation |
| `v1.vault.create-wrong-type.001` | Wrong-type canonical entry is not absence |
| `v1.vault.object-publication-unchanged.001` | New ordinary TOKEN publication preserves canonical VAULT exactly |
| `v1.vault.object-exact-retry-unchanged.001` | Exact-existing object retry preserves canonical VAULT |

Existing ASCII/empty/Unicode cases still unlock and now check VAULT_ID.
`create-empty` proves initial creation; `create-existing` now supplies exact existing
canonical bytes and blocks creation; `create-orphans` now expects FAILED.
`TestTokenAuthorshipPreservesCanonicalVault` additionally constructs/authenticates
real TOKEN bytes, exercises new/exact/failed publication, and proves VAULT and its
unlocked root remain unchanged. The normative checker and schema regression test
prove the absence of a rewrap/replacement operation. No application API or
cross-vault transfer behavior was added to the corpus.

## 21. Existing vector byte invariance

All four existing bootstrap cases retain exactly `root_hex`, `password_hex`,
`salt_hex`, `nonce_hex`, `wrap_key_hex`, `header_hex`, and `record_hex`. The three
surviving IDs replace only fingerprint metadata with VAULT_ID; the fourth changes
ID/notes/recognition metadata to retire lifecycle meaning. Public test roots shared
across known answers are not legitimate alternate same-vault bootstrap state.

These baseline and final crypto-field digests are identical. Each digest hashes
UTF-8 compact JSON of those seven named fields, sorted by key (`sort_keys=True`,
`separators=(',', ':')`):

| Baseline bootstrap | Identical baseline/final SHA-256 of crypto fields |
| --- | --- |
| ASCII | `ed1afbe4616b16bdcbb8d38ba4fd53e3c3c5e054ac8947f0abce42f9e092895c` |
| Empty | `b4f63dfa6c22a120847f613e403d3bf6492f688195e51dbc1d593bf239b7a695` |
| Unicode | `9f29157867c9f6fb21401aceba9c592d5e7440e1eef7db1b9a152cadd337d136` |
| Former rewrap → static extra | `14b0b08d6198e03146ab40fbdbcab0f7123d86c6dd1100ec55f2bf2adebdf045` |

**72 surviving complete case files are byte-identical**: every unrelated crypto,
encoding, metadata, size, graph, fold, TOTP, and storage-observation file. Six
object-publication workflow files only remove empty `base_hex`; their intended/
existing object bytes, flags, and outcomes remain identical. Four initial-creation
files intentionally change their modeled inputs; only orphan creation changes its
expected result. Every changed surviving file is listed below.

| Changed surviving case | Reason |
| --- | --- |
| `v1.bootstrap.ascii.001` | Replace fingerprint metadata with SHA-256 VAULT_ID; all root-wrap inputs and bytes identical. |
| `v1.bootstrap.empty.001` | Replace fingerprint metadata with SHA-256 VAULT_ID; all root-wrap inputs and bytes identical. |
| `v1.bootstrap.unicode.001` | Replace fingerprint metadata with SHA-256 VAULT_ID; all root-wrap inputs and bytes identical. |
| `v1.storage.ambiguous-publication.001` | Remove empty replacement-only base_hex field; all object bytes and modeled outcomes unchanged. |
| `v1.storage.exact-existing-publication.001` | Remove empty replacement-only base_hex field; all object bytes and modeled outcomes unchanged. |
| `v1.storage.missing-parent-publication.001` | Remove empty replacement-only base_hex field; all object bytes and modeled outcomes unchanged. |
| `v1.storage.new-publication.001` | Remove empty replacement-only base_hex field; all object bytes and modeled outcomes unchanged. |
| `v1.storage.non-exact-publication.001` | Remove empty replacement-only base_hex field; all object bytes and modeled outcomes unchanged. |
| `v1.storage.unreadable-target.001` | Remove empty replacement-only base_hex field; all object bytes and modeled outcomes unchanged. |
| `v1.vault.create-ambiguous.001` | Remove empty base_hex; supply canonical intended bytes and successful reread; remove orphan flag to isolate durability failure; still FAILED. |
| `v1.vault.create-empty.001` | Remove empty base_hex; supply canonical intended bytes and successful canonical revalidation; still CREATED. |
| `v1.vault.create-existing.001` | Remove empty base_hex; supply exact existing/intended canonical bytes and successful read outcome; still FAILED without replacement. |
| `v1.vault.create-orphans.001` | Remove empty base_hex; supply intended canonical bytes and revalidation outcome; orphan presence changes result CREATED → FAILED. |

Reproducible field/whole-file audit against the unchanged HEAD:

```sh
python3 - <<'AUDIT'
import hashlib, json, subprocess
from pathlib import Path
old = lambda p: subprocess.check_output(['git', 'show', 'HEAD:' + p])
m = json.loads(old('vectors/manifest.json'))
current = json.loads(Path('vectors/manifest.json').read_bytes())
ids = {e['id'] for e in current['cases']}
same = 0
fields = ['password_hex', 'salt_hex', 'nonce_hex', 'wrap_key_hex',
          'header_hex', 'record_hex']
for e in m['cases']:
    p = 'vectors/' + e['path']
    a = json.loads(old(p))
    target = p.replace('bootstrap.rewrap', 'bootstrap.known-answer-extra')
    if a['operation'] == 'bootstrap':
        b = json.loads(Path(target).read_bytes())
        assert a['root_hex'] == b['root_hex']
        assert all(a['bootstrap'][f] == b['bootstrap'][f] for f in fields)
        assert hashlib.sha256(bytes.fromhex(b['bootstrap']['record_hex'])).hexdigest() == b['bootstrap']['vault_id_hex']
    elif e['id'] in ids:
        if a['operation'] != 'workflow':
            assert old(p) == Path(p).read_bytes(), p
        if old(p) == Path(p).read_bytes():
            same += 1
assert same == 72
print('PASS: four bootstrap crypto-field sets and 72 complete surviving files')
AUDIT
```

## 22. Storage/reference implementation changes

`storage.Create` is absent → initial-create only and models complete local
validation, durability, canonical reread/revalidation, orphan safety, and unchanged
existing bytes on failure. `storage.Replace` and replacement tests are removed.
`Store.Publish` keeps VAULT separate and delegates to the unchanged immutable
object installation model. Existing canonical entries can still be read/unlocked
by `OpenBootstrap` and recognized by VAULT_ID. There is no replacement transition.

Backend booleans are trusted modeled outcomes, not an implementation of staging,
crash recovery, syscalls, races, or remote completeness. No provider-specific SAF
semantics were introduced. The bootstrap runner checks all original known-answer
fields, unwraps the record, and checks SHA-256 identity. Python independently
checks identities and all existing object-key intermediates. The mutation fixture
checks one bit differs, both digest answers, and authentication rejection.

## 23. Requirements/profile changes

The moving profile remains `moving-pre-rc`, protocol `totipo-v1`, revision r19,
with all 92 manifest cases required; no RC freeze or conditional capability.
Specification, manifest, manifest-schema revision constant, profile, CLI/runner/
generator/checker constants, Makefile help, README, CONTRIBUTING, requirements
README, vector README/FORMAT, and conformance README now agree on r19.
Case-schema changes replace fingerprint metadata with VAULT_ID, add the focused
changed-record identity fixture and optional held VAULT in publication workflows,
and remove replacement action/results/base field. These are fixture contracts,
not protocol byte fields. Expanded schema copies were updated consistently.
Workflow manifest references now point to Sections 3, 7, and 18 rather than the
removed replacement operation. All physical paths, required IDs, and hashes were
independently verified. The canonical generator was necessary to reconcile these
intentional corpus/schema/pin changes; a second run changed no bytes.

| Final artifact | SHA-256 |
| --- | --- |
| `spec/totipo-vault-format-v1.md` | `8bb76b890086eb4bf89271edd5861e02f833f3ffa11a83b5865080ed4e4228cb` |
| `vectors/manifest.json` | `3953dbc315b4dcb3d0d31dd399bf82a15cb38c49fde2e2b93ea81c5cfe0ec714` |
| `requirements/v1-pre-rc.json` | `7242fc557a7ca56452c934bedc5e7dc2835248f5cd36fa220f05588dc967ec80` |
| `vectors/case.schema.json` | `99805d442f13872cd4febe9ac8ae4f36fd25607580bc1e457ff5ae2efe199e23` |
| `vectors/manifest.schema.json` | `f265771f904be57d19602dd6da9a572c16b3f80fdd166c5177721aaa3f9a91ec` |

## 24. Security implications

An old source vault copy plus its password can still reveal that vault. With a
fresh independent destination root, old VAULT/password does not derive the new
root. Copying an already exposed TOTP seed cannot erase attacker knowledge;
issuer-side seed rotation/reissue is needed for revocation. No new rollback,
availability, first-contact authenticity, or exhaustive migration guarantee is
claimed. Strong initial publication and explicit orphan handling reduce accidental
namespace substitution without requiring proof that an untrusted store is complete.

## 25. Downstream Java impact

Expected work in `totipo-java`: remove password-change API and same-root rewrap
orchestration; remove/simplify replacement-capability SPI, staged replacement,
and replacement-specific storage implementations/tests. Replace fingerprint
recognition with exact immutable VAULT_ID, preserve initial no-replace durability,
add canonical revalidation and orphan-creation safety as needed. Root/object crypto
and wire encoding stay unchanged. No Java repository was modified.

## 26. Downstream desktop impact

Remove in-place Change Password UX and its reauthentication/replacement flow.
A later, separately designed migration workflow may create a new vault and copy
selected visible validated state/history with honest completeness disclosure.
No migration UI is required by r19 conformance or implemented here.

`docs/design/DESIGN.md` labels itself the non-normative v0.9 **v0 product and
visual-design baseline** and contains old Change Vault Password guidance. It was
left byte-identical under the instruction not to edit v0 material. It cannot
override current r19; its password-change flow and related device-local wrapper
wording require downstream product-design revision before use with r19.
No desktop repository was modified.

## 27. Downstream Android impact

Finish provider canonical VAULT bootstrap/binding around immutable VAULT_ID.
Gate existing M3A/M3B object synchronization on matching immutable VAULT; diagnose
substitution/mismatch rather than reconcile mutable provider bootstraps.
No mutable provider VAULT reconciliation is required. Preserve creation durability,
orphan handling, and normal untrusted observation. No SAF-specific semantics or
Android repository changes were made here.

## 28. Validation

On the user-requested follow-up, the agent reran `make check`, `make verify`,
`make conformance`, `make race`, `make fuzz`, vet, formatting, and Git diff checks
directly. All passed again, including 92/92 cases and all three bounded fuzz targets.

| Final check | Result |
| --- | --- |
| `make check` | PASS: spec structure/pins, 14 Python tests, uncached Go tests, 92/92 cases, schemas and profile |
| `make verify` | PASS: both schemas, exact physical coverage, case hashes and all integrity pins |
| `make conformance` | PASS: all 92 required cases |
| `make race` | PASS: uncached Go race suite |
| `make fuzz` | PASS: parser, envelope, and arrival/disappearance graph targets; Makefile 10s bounds, two workers |
| `gofmt -l conformance` | PASS: no files listed |
| `go -C conformance vet ./...` | PASS |
| Independent raw-byte/pin/case audit | PASS: all 92 required IDs/hashes, both schema pins, exact case coverage |
| Independent invariance audit | PASS: four old bootstrap crypto-field sets; 72 whole files; six publication object payloads; unchanged protocol Sections 4–5 and 10–19, VAULT encoding/crypto blocks, and old revision history |
| Independent SHA-256/mutation audit | PASS: all bootstrap IDs and one-bit changed representation ID |
| Crypto source comparison to HEAD | PASS: unchanged Derive, object crypto/auth helpers, WrapKey, Wrap, Unwrap |
| Repeat canonical generation | PASS: no case, manifest, or profile byte changed |
| `git diff --check` | PASS |
| Git status/diff/cached diff | Reviewed; unstaged changes only; cached diff empty |
| `nix flake check` | PASS, reported by the user; not executed by the agent |

No archived v0 material, historical review artifacts, design assets, flake/lock,
implementation repository, object/TOKEN grammar, TOTP behavior, or crypto
construction changed. No new format version, marker, compatibility branch, transfer
protocol, or migration schema was introduced. Abstract cases do not establish live
backend crash behavior or independent interoperability.

## 29. Human Nix

**PASS, user-reported:** the user confirmed that they ran `nix flake check`
successfully. This result is accepted as the human Nix validation evidence; the
agent did not execute the Nix binary. A shell attempt to locate/run Nix found no
executable (`nix: command not found`, exit 127), before the user clarified that
only the Make checks needed to be run directly.

The repository documents `nix develop` (or existing direnv) followed by ordinary
Make/Go checks. The unchanged flake exposes a development shell and formatter,
with no dedicated `checks` suite and no default package/app. Consequently
`nix build path:.` has no declared default build output and is not an applicable
package validation target here. No dependency update is expected or made.

## 30. Remote CI

Not triggered. No remote CI, release, tag, push, or remote mutation was performed.
Local ordinary validation is the evidence reported here.

## 31. Final Git state

Branch and HEAD remain `main` / `fbae625f2cad6c89a1b5555de2f2d75d890bb682`.
All edits, five removed case files, seven new case files, and the two new review
documents are unstaged/uncommitted. `git diff --cached --stat` is empty. No stage,
commit, tag, release, push, or Nix invocation occurred. Only this specification
repository was changed. Final status and diff stat were inspected; Git's ordinary
diff stat excludes the untracked new documents/cases.
