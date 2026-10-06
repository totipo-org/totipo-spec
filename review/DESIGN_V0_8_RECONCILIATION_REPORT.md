# Design v0.8 reconciliation report

## Starting state and source

- Branch: `main`.
- Full starting HEAD: `4569230c10645f28d215d0b86515bb603b3bd20e`.
- Starting `git status --short`: empty (clean tree), verified before editing.
- Source: committed [Totipo Design Guidelines](../docs/design/DESIGN.md) at that HEAD, headed `Draft v0.7`; working copy now says `Draft v0.8`.
- Scope remains the cross-platform v0 product and visual-design baseline; protocol status remains Non-normative.

## Revision summary and rationale

v0.8 corrects two product-model assumptions:

1. **Change Vault Password requires current-password reauthentication.** An unlocked session permits ordinary vault use and ordinary Add/Edit/Delete/Resolve operations, but alone does not authorize replacement of the credential controlling future vault access. Requiring the current vault password prevents someone with temporary access to an unlocked desktop from silently replacing its future unlock password. It avoids retaining the original unlock password merely for later reuse and aligns with the existing Java/password-rewrap security model. Biometric/device unlock is a local convenience for regaining access to a known vault and does not substitute for current vault password entry during password change. The screen collects current password, new password, and new-password confirmation; wrong current password is a local authentication failure, and an empty new password still needs explicit confirmation.
2. **v0 has no inferred/global Read-only vault mode.** Failed writes, definite publication failure, publication uncertainty, invalid/corrupt data, and unavailable data do not establish a persistent global storage capability. They retain their truthful distinctions, including blocking only where a usable vault cannot safely be established and localized evidence where safe. Publication outcomes remain operation-specific. Explicit read-only sessions or affirmatively reported backend read-only capabilities require separate future product design; neither is introduced here. About This Vault reports only available capability/format facts and never fabricates writable status or uses an installed library version as vault-format metadata.

Password change still protects/rewraps the same vault key/root, without root rotation or re-encryption of all contents. Retained historical copies are not promised revoked or securely erased.

**Protocol impact: NONE.** Protocol semantics changed: **NO**. Protocol revision changed: **NO**. Wire format changed: **NO**. Schemas, vectors, conformance expectations, and cryptographic behavior are unchanged. The normative specification remains v1/r18.

DESIGN.md has no established revision-history section; this report records the revision summary without introducing one there.

## Design sections changed

| Section | Correction |
| --- | --- |
| Header | Draft v0.8; existing scope and non-normative status preserved. |
| §1, §2.5 | Constrained/unavailable conditions and textual action state replace global read-only references. |
| §14, §§14.1–14.4 | Current-password authorization, three-field form, validation, empty-new-password confirmation, same-root rewrap truth, and capability-based availability. |
| §16.6 | Local unlock does not substitute for vault-password reauthentication; §16.7 reviewed and existing invalidation/re-establishment rule preserved unchanged. |
| §17, §§17.1–17.3, §17.5 | Remove required global mode/banner/mutation list; preserve truthful distinctions, blocking boundary, and operation-specific publication outcomes. §17.4 remains byte-identical. |
| §18 | Truthful optional storage/access capability and actual vault-format metadata. |
| §20 | Remove Read-only shared state; keep Disabled as an action state. |
| §23 | Genuine capability/session unavailability replaces read-only-vault example. |
| §§25.3–25.5 | Explicit ReadOnlyValue styling; selection preserves warning/conflict meaning. |
| §29 | Inferred global classification and Open Read-Only mode are explicit non-goals. |
| §§30.7–30.9 | Corrected vault/password/local-unlock checklist and explicit read-only value distinction. |
| §31 | Immediate blocking/unavailable work and password reauthentication; optional future explicit read-only feature design. |

§24 accessibility and §28 ReadOnlyValue catalogue remain intact.

## Remaining read-only occurrences

Every case-insensitive `read-only`, `read only`, `writable`, and `ReadOnlyValue` occurrence in DESIGN.md was reviewed:

| Location | Why retained |
| --- | --- |
| §9.6 (existing value fields) | Non-editable, selectable conflict-resolution values, distinct from disabled controls. |
| §17.1 heading and three paragraphs | Explicit denial of an inferred v0 mode, prohibition on relabeling failures, and hypothetical future sessions/capabilities. |
| §18 capability paragraph | Prohibits inferred writable/read-only status and requires omission or explanation when capability is unreported; defines no global mode. |
| §24 accessibility | Read-only displayed values remain readable rather than disabled-looking. |
| §25.3, §25.4, §25.5 | ReadOnlyValue component styling and distinction from disabled controls. |
| §28 | ReadOnlyValue component catalogue entry, with no ReadOnlyVault component. |
| §29 (two items) | Explicit non-goals for inferred classification and Open Read-Only sessions. |
| §30.7 | Checklist prohibition on inferred global state, not a state requirement. |
| §30.9 | Read-only UI values distinguished from editable inputs and disabled controls. |
| §31 | Hypothetical future feature, conditional on a concrete requirement and separate design. |

There are no `read only` (space-separated) occurrences in DESIGN.md. The only remaining `writable` occurrence prohibits fabricated status in §18. The two-field new-password/confirmation collection in §13 is **vault creation**, not password change, and remains appropriate.

## Repository cross-references

Repository-wide searches, including hidden files while excluding Git metadata and the local build environment, found no other live v0.7 baseline references, duplicated Change Vault Password forms, or required global read-only product states. No cross-reference files needed changes.

Untouched legitimate non-product uses include the normative specification's read-only exact comparison (§18), historical storage/rewrite reports, audit scripts describing non-mutating audits, and historical JSON audit evidence. Historical reports were not rewritten. New v0.7 references in this report describe the starting state only.

## Validation

- `make check`: PASS; v1/r18 structural/profile hashes, 13 Python tests, Go tests, all 90 conformance cases, JSON schemas, manifest/profile verification.
- `make verify`: PASS independently; schemas and all 90 manifest/profile case files verified.
- `gofmt -l conformance`: PASS; no output.
- `go -C conformance vet ./...`: PASS.
- `make race`: PASS.
- `make fuzz`: PASS; bounded FuzzDispatch, FuzzOpen, and FuzzArrivalAndDisappearance runs (10 seconds each).
- `git diff --check`: PASS.
- Repository `rg` searches for `v0.7`, `read-only`, `Change Vault Password`, `New password`, and `writable/read-only`: completed before editing. Follow-up case-insensitive searches also included `read only`, current/new password, ReadOnlyValue, hidden files, and old `Confirm password` examples: no stale live requirements found.
- Direct document assertions: PASS for Draft v0.8, non-normative header, exact three-field password screen, removed global-state row/availability rule, retained ReadOnlyValue, unchanged localized-failure subsection, and design-only tracked diff.
- Protocol source, vectors, schemas, tooling, requirements, and historical evidence have no working-tree changes. No vectors or metadata were regenerated.

## Files changed and final review state

- `docs/design/DESIGN.md`: product/design correction.
- `review/DESIGN_V0_8_RECONCILIATION_REPORT.md`: new reconciliation evidence.

Final `git status --short`:

```text
 M docs/design/DESIGN.md
?? review/DESIGN_V0_8_RECONCILIATION_REPORT.md
```

Final `git diff --stat` (Git excludes the untracked report):

```text
 docs/design/DESIGN.md | 118 +++++++++++++++++++++++++-------------------------
 1 file changed, 58 insertions(+), 60 deletions(-)
```

HEAD and branch remain unchanged. The index is empty of changes. All work remains unstaged/uncommitted; no commit, tag, release, or push was performed.
