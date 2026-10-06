# Design v0.9 Manual Add report

## Starting state and source

- Branch: `main`.
- Full starting HEAD: `279fb74383ed1ee32f7a6c14880b9285386d9356`.
- Starting `git status --short`: empty (clean tree), verified before editing.
- Exact Design v0.8 source commit: `279fb74383ed1ee32f7a6c14880b9285386d9356` (`updated design`), the latest commit touching `docs/design/DESIGN.md`.
- Inspected the committed `docs/design/DESIGN.md`: title `Totipo Design Guidelines`, status `Draft v0.8`.
- Working design now says `Draft v0.9`; cross-platform v0 product/visual-design scope and `Protocol status: Non-normative` remain unchanged.

## Revision summary and rationale

v0.9 removes the redundant Review step after Add TOTP Manual Entry. Manual Entry is the editable review surface and uses `Add` after validation; QR/setup-URI import still requires explicit Review before commit.

Imported acquisition requires review of interpreted identity/configuration. Direct manual entry already exposes issuer/service, account, new secret, algorithm, digits, and period for editing. A second Review screen repeats information just entered and adds friction without increasing understanding: all relevant non-secret configuration is already visible/editable. This shared rule applies across desktop and Android.

- Old Add flow: `Manual → Review → Add`.
- New Add flow: `Manual → Validate → Add/duplicate decision`.
- Imported flow unchanged: `URI/QR → Review → Add/duplicate decision`.
- Invalid Manual Add stays on the same form, does not publish, presents local field/group errors, and focuses/scrolls to the first actionable problem.
- Valid Manual Add checks current active issuer/account matches and publishes normally when no duplicate decision is needed.

The document has no revision-history section; this report records the rationale without introducing one there.

## Exact design sections changed

| Section | Change |
| --- | --- |
| Header | Draft v0.9; same scope and non-normative status. |
| §3 | Mode-specific next action: Manual Add/reviewed setup uses Add; QR/setup-URI acquisition uses Review. |
| §6.1 | Separate imported and manual flow diagrams, explicit review boundary, cross-platform semantics, labels, and concise rationale. |
| §6.4 | New secret field wording; editable review surface, Add primary action, authoritative validation, local invalid-input handling, direct duplicate handling/publication, no secret reconfirmation. |
| §6.5 | Review imported setup heading and explicit exclusion of Add TOTP Manual Entry; retains meaningful identity and non-secret configuration. |
| §6.7 | Identity check follows imported review or validated manual entry; explicit manual Add outcomes, multiple-match cancellation, and preserved setup-replacement path. |
| §8.2 | Explicit narrow scope: Change Authenticator Setup remains unchanged, including manual replacement and Update Existing; imported values remain reviewable before replacement. |
| §§16.6, 17.1, 18, 30.8 | Current-draft references advanced from v0.8 to v0.9 only; existing password/local-unlock and vault-state rules unchanged. |
| §28 | AddTOTPFlow paths and TOTPReview applicability clarified; review component retained. |
| §30.5 | Explicit imported Review and direct validated Manual Add checklist requirements. |
| §31 | Live implementation work now distinguishes imported URI review from direct validated Manual Add. |

§6.2 platform acquisition support, §6.3 secret presentation, §6.6 parsing/validation guidance, and §23.1 validation-on-activation remain byte-identical to the starting document.

## Preserved safety and scope decisions

**Duplicate semantics unchanged.** Matching uses issuer/account in the current active TOTP set, never secrets. Exactly one match offers Update Existing…, Add Another, and Cancel. Multiple matches require choosing the existing TOTP to update, Add Another, or Cancel; Totipo never guesses or silently replaces. Manual Add with a match does not immediately publish. Add Another confirms a new logical TOTP; Update Existing… follows existing setup-replacement semantics; Cancel returns to the manual form/draft where practical.

**Secret presentation unchanged.** Imported QR/URI review does not unnecessarily redisplay raw secret material. Manual entry exposes only the new secret being actively entered. Existing stored secrets are never redisplayed merely for review. No new secret exposure or required re-entry/reconfirmation was added.

**Validation-on-activation unchanged.** Manual Add remains activatable with incomplete/invalid input; activation validates authoritatively and exposes local errors and first-problem focus without publishing invalid input. The generic §23.1 rule is unchanged.

**Change Setup scope: narrow Add TOTP change only.** Change Authenticator Setup retains its existing acquisition/review behavior, including manual replacement. Choosing Update Existing… from a duplicate decision retains review of proposed replacement setup for the existing logical TOTP. Imported QR/URI values remain reviewable before replacement. Direct Manual Add semantics are not silently extended to replacement.

**Protocol impact: NONE.** Protocol semantics changed: **NO**. Protocol revision changed: **NO**. Wire format changed: **NO**. The normative specification remains v1/r18. Requirements profiles, schemas, vector manifests/cases, conformance implementation/expectations, cryptographic fixtures, and cryptographic behavior are unchanged.

## Repository consistency review

Before editing, searched the entire repository (including hidden files, excluding Git metadata) for `v0.8`, `Manual entry`, `Review`, `Add TOTP`, `QR/URI acquisition`, and `Parsed TOTP draft`. No applicable AGENTS.md files were found.

After editing, case-insensitive repository searches covered `Manual entry`, `Manual setup`, `Review`, `Parsed TOTP draft`, `Add / Update`, `QR/URI`, `TOTPReview`, and `Add TOTP`. Reviewed all live design occurrences: no live Add TOTP statement still requires Manual → Review → Add. Android scanning still explicitly returns to Review; shared manual semantics apply on Android as well. Change Setup review is retained as the documented scope exception. Conflict review and operating-system preview references are unrelated and unchanged.

No other current cross-reference files needed edits. Remaining v0.8 occurrences describe historical evidence or the starting state in this report. Historical reports, including DESIGN_V0_8_RECONCILIATION_REPORT.md, remain untouched.

## Validation

- `make check`: PASS; v1/r18 structure/profile hashes, 13 Python tests, Go tests, all 90 conformance cases, schemas, manifest/profile verification.
- `make verify`: PASS independently; schemas and all 90 manifest/profile case files verified.
- `gofmt -l conformance`: PASS; no output.
- `go -C conformance vet ./...`: PASS.
- `make race`: PASS.
- `make fuzz`: PASS on independent rerun of the unchanged standard target; FuzzDispatch, FuzzOpen, and FuzzArrivalAndDisappearance each completed their bounded 10-second run. Initial run stopped at FuzzDispatch's 10-second bound with `context deadline exceeded`; both attempts are recorded here. No source or fixture changes were needed.
- `git diff --check`: PASS on final tracked diff; `git diff --no-index --check /dev/null review/DESIGN_V0_9_MANUAL_ADD_REPORT.md` also PASS for the untracked report.
- Direct document assertions: PASS for Draft v0.9, non-normative header, separate manual/imported paths, Add activation and invalid-input handling, explicit imported Review, narrow Change Setup scope, preserved identity-only duplicates, and byte-identical §§6.2, 6.3, 6.6, 23.1.
- Final scope checks: tracked diff contains only DESIGN.md; the only untracked file is this report. Protocol source, requirements, schemas, vectors, conformance code, tooling, and existing review reports have no changes. Branch and HEAD remain at the starting values; index empty.

## Files changed and final review state

- `docs/design/DESIGN.md`: narrow product/design clarification and current version references.
- `review/DESIGN_V0_9_MANUAL_ADD_REPORT.md`: new review evidence.

Final `git status --short`:

```text
 M docs/design/DESIGN.md
?? review/DESIGN_V0_9_MANUAL_ADD_REPORT.md
```

Final `git diff --stat` (Git excludes the untracked report):

```text
 docs/design/DESIGN.md | 72 +++++++++++++++++++++++++++++++++++----------------
 1 file changed, 50 insertions(+), 22 deletions(-)
```

All work remains unstaged/uncommitted. No commit, tag, release, or push was performed.
