# Totipo Design Guidelines

**Status:** Draft v0.3  
**Scope:** Cross-platform product and interaction design  
**Protocol status:** Non-normative

These guidelines define Totipo's shared product language across desktop and Android, and are intended to remain applicable to future platforms such as iOS.

They do **not** define the Totipo storage or synchronization protocol and do not override protocol requirements. Platform applications should follow native platform conventions where those conventions do not conflict with the product semantics defined here.

> **Consistency of meaning and behavior is more important than consistency of appearance.**

---

## 1. Product model

Totipo is primarily a TOTP retrieval tool backed by a synchronized vault.

The normal user journey is:

1. Open or unlock the remembered vault.
2. Find the required TOTP.
3. Reveal its current code.
4. Copy or use the code.

The interface should optimize for this path.

The product hierarchy is:

1. **Primary product task**
   - Find a TOTP.
   - Show its current code.
   - Copy/use the code.

2. **Primary collection action**
   - Add TOTP.

3. **Secondary management**
   - Edit TOTP.
   - Delete TOTP.
   - Change vault.
   - Refresh when manual refresh is required.
   - Change vault password.

4. **Exceptional state handling**
   - Resolve conflicts.
   - Handle invalid/corrupt data.
   - Handle read-only or otherwise constrained vault state.
   - Handle uncertain publication or newly arrived conflict information.

Exceptional operations may become visually prominent when they are relevant, but they should not dominate the normal interface.

---

## 2. Design principles

### 2.1 Calm

Totipo is a security utility. Screens should contain only what supports the current task.

Avoid unnecessary decoration, animation, permanent success messages, excessive color, and controls that are present only because the underlying toolkit makes them easy to expose.

### 2.2 Explicit

Security-relevant consequences and exceptional states should be explained in plain language.

Do not rely on vague confirmation text such as "Are you sure?" when the important consequence can be stated directly.

### 2.3 Predictable

Equivalent actions should have equivalent semantics across screens and platforms.

The same concept should use the same terminology throughout Totipo.

### 2.4 Native

Totipo should not attempt to make desktop and Android pixel-identical.

Desktop should behave like a desktop application: keyboard navigation, focus states, menus, compact layouts, pointer interaction, resizable windows, and platform-appropriate dialogs.

Android should behave like an Android application: touch-sized controls, system back behavior, camera-based QR acquisition where appropriate, and platform-standard accessibility.

### 2.5 State is visible

Conflict, error, read-only operation, pending state, and similar conditions must be represented textually.

Color, icons, borders, and animation may reinforce meaning but must not carry the meaning alone.

### 2.6 Security without unnecessary friction

Totipo should minimize unnecessary exposure of TOTP values and credential material, but should not make ordinary authentication workflows cumbersome.

Security-sensitive behavior should be explicit and predictable rather than clever or surprising.

---

## 3. Interaction hierarchy

A task context should have one obvious next action.

Examples:

- Hidden TOTP: **Show Code**
- Revealed TOTP: **Copy**
- Vault collection: **Add TOTP**
- Create/add flow: **Add TOTP**
- Edit flow: **Save**
- Conflict: **Resolve**
- Vault selection: **Select Folder**
- Unlock: **Open**
- New vault: **Create Vault**
- Password change: **Change Password**

This does not imply one global primary button per screen. The main TOTP list contains multiple independently actionable rows, each of which may have its own primary action.

Rare actions should not consume permanent visual space merely because they are important when used.

---

## 4. Action roles

### 4.1 Primary action

Advances or completes the current task.

Examples:

- Show Code
- Copy
- Add TOTP
- Save
- Resolve
- Open
- Select Folder
- Create Vault
- Change Password

### 4.2 Secondary action

A normal but less important action.

Examples:

- Change Vault
- Refresh
- Edit when explicitly exposed

### 4.3 Quiet action

A low-emphasis action that should not compete with the primary task.

Examples:

- Cancel
- Clear Search
- ancillary navigation

Cancel is neutral. It must not use danger styling merely because it closes a dialog.

### 4.4 Destructive action

An action with destructive or difficult-to-reverse consequences.

Examples:

- Delete TOTP
- destructive filesystem action, where such an action is intentionally exposed

Danger styling is reserved for destructive semantics.

---

## 5. TokenRow

`TokenRow` is the central Totipo component. It represents one usable TOTP credential in the main vault view.

Its priorities are:

1. Make the TOTP easy to recognize.
2. Make revealing the current code obvious and fast.
3. Once revealed, make the code easy to read and copy.
4. Communicate how long the current code remains valid.
5. Keep uncommon management actions out of the primary path.

### 5.1 Anatomy

Conceptually:

```text
TokenRow
    Identity
        PrimaryLabel
        SecondaryLabel?

    TOTP
        Hidden
            ShowCode
        or
        Revealed
            Code
            RemainingTime
            TimeIndicator?
            Copy

    Management
        MoreActions?
```

A typical hidden desktop row may look like:

```text
GitHub
niki@example.com                              Show Code    ⋮
```

A revealed row may look like:

```text
GitHub                              213 933
niki@example.com                    25 s    ◯    Copy    ⋮
```

Exact layout is platform-specific.

### 5.2 Identity

For a conventional TOTP:

- issuer is the primary label;
- account is the secondary label.

If only one is present, promote the available value rather than leaving empty visual space.

If neither is available, display a stable fallback such as `Unnamed TOTP`.

### 5.3 Hidden state

Tokens are concealed by default.

`Show Code` is the primary row action.

Management operations such as Edit and Delete must not compete visually with `Show Code`.

### 5.4 Revealed state

Revealing happens in place. It is a state change of the row, not navigation to a new screen or dialog.

The code should become one of the strongest typographic elements in the row.

For readability, visually group digits with whitespace:

```text
6 digits    213 933
7 digits    213 9334
8 digits    2139 9334
```

Do not add punctuation to the credential value merely for visual grouping.

Use tabular digits where supported. A monospaced font is acceptable but not required.

Accessible output should make the digits understandable as a code rather than as one large numeric quantity.

### 5.5 Countdown

Remaining validity must be available textually, for example:

```text
25 s
```

A ring or progress indicator may supplement the text but must never be the only indication of remaining validity.

Near expiration, Totipo may visually emphasize that the code is expiring, but should avoid distracting flashing or alarm-like behavior.

### 5.6 Bounded reveal lifetime

Revealing a TOTP is a bounded disclosure.

Normal reveal:

```text
HIDDEN --Show Code--> REVEALED
                         |
                      rollover
                         |
                         v
                       HIDDEN
```

If the user reveals a code with **10 seconds or less remaining** in the current TOTP period, Totipo provides one additional period of visibility:

```text
HIDDEN --Show Code--> REVEALED_CURRENT
                         |
                      rollover
                         |
                         v
                    REVEALED_GRACE
                         |
                      rollover
                         |
                         v
                       HIDDEN
```

Therefore:

> A reveal normally lasts through the current TOTP period. If revealed with 10 seconds or less remaining, Totipo also reveals the immediately following period, then automatically conceals the row.

Copying a code does not extend its reveal lifetime.

The 10-second threshold is a Totipo UX policy. It is not a protocol rule.

### 5.7 Independent rows

Reveal state is tracked independently for every token.

> Revealing, copying, rolling over, or concealing one TokenRow must not alter the reveal state of any other TokenRow.

Several TOTP codes may therefore be visible simultaneously. This supports workflows where a user authenticates to multiple services in quick succession.

### 5.8 Copy

Once a row is revealed, `Copy` is the immediate primary action for that row.

Copy feedback should be:

- immediate;
- local to the row;
- transient;
- accessible;
- non-modal.

For example, `Copy` may temporarily become `Copied`.

Copy must not be disabled merely because the code is near expiration.

> Copy always copies exactly the code currently displayed by the row.

### 5.9 Manual concealment

A dedicated persistent `Hide` control is not required for v0 because disclosure is already bounded by TOTP rollover.

Platforms may provide an unobtrusive explicit conceal action if it improves usability, but it must not clutter the main reveal/copy interaction.

### 5.10 Management actions

Edit and Delete are secondary management operations.

They should normally live behind a low-emphasis management affordance such as an overflow/context menu, and may additionally be exposed through desktop application menus or keyboard shortcuts.

Management affordances require accessible names that identify the affected token.

### 5.11 Pointer/touch behavior

Android may use a larger portion of the row as the reveal target where that improves touch ergonomics.

Desktop should prefer an explicit `Show Code` control and should avoid revealing credential values merely because identity text was clicked accidentally.

### 5.12 Keyboard behavior

Desktop must provide keyboard access to row actions.

The row should not create unnecessary focus stops. In the normal hidden state, the principal focusable controls should be approximately:

- Show Code
- More Actions, if present

In the revealed state:

- Copy
- optional explicit conceal action
- More Actions

### 5.13 Focus, hover, and reveal are distinct

Do not use the same visual treatment ambiguously for:

- pointer hover;
- keyboard focus;
- selection;
- revealed state.

Totipo should avoid introducing a persistent "selected token" state unless a concrete feature requires one.

### 5.14 Conflicts

Conflict is not a TokenRow state.

A conflict may contain multiple candidate versions of the same logical token. The conflict therefore belongs to a containing component, for example:

```text
ConflictPanel
    Resolve
    TokenCandidate
    TokenCandidate
```

A token candidate may reuse the normal reveal/copy interaction where safe and appropriate.

---

## 6. Add TOTP

`Add TOTP` is the primary collection-level action.

Use `Add TOTP`, not `Create Token`, because Totipo normally receives authenticator credentials created by another service.

### 6.1 Acquisition and review are separate

Credential acquisition is separate from credential review.

Conceptually:

```text
Acquire
    |
    +-- Scan QR
    +-- Paste setup URI
    +-- Manual entry
    |
    v
Parsed TOTP draft
    |
    v
Review
    |
    v
Add / Update
```

Scanning or pasting must not immediately commit a credential.

### 6.2 Platform acquisition methods

Android should support local enrollment without requiring desktop.

The preferred Android path is expected to be QR scanning, with manual entry always available.

Desktop should support local enrollment without requiring Android.

For desktop v0, the preferred acquisition methods are:

- paste/import an `otpauth://` setup URI;
- manual entry.

Desktop camera scanning is not required for v0.

A credential added on any device is simply a normal vault change and is expected to propagate through Totipo synchronization. No special transfer mechanism is required.

### 6.3 Raw secret presentation

The raw secret is acquisition material, not normal review information.

After QR or URI acquisition, the normal review screen should not redisplay the raw secret unnecessarily.

Manual entry exposes the secret because the user must enter it.

### 6.4 Manual entry

Manual setup should collect:

- issuer/service;
- account;
- secret.

Authenticator parameters should be secondary or advanced:

- algorithm;
- digits;
- period.

The common configuration should not dominate the form.

### 6.5 Review

The review screen should present the user-meaningful identity and authenticator configuration.

Imported non-default parameters must be visible before commit. Totipo should not silently hide unusual imported configuration.

### 6.6 Validation

Validation should be local and specific.

Examples include:

- malformed Base32;
- malformed `otpauth://` URI;
- unsupported algorithm;
- unsupported digit count;
- invalid period.

Friendly normalization, such as accepting spaces in a manually entered Base32 secret, is acceptable. Totipo should not silently guess ambiguous credential data.

### 6.7 Existing issuer/account

Totipo does not scan the vault for duplicate secrets.

Instead, after acquisition/review, Totipo checks the current user-visible token set for an existing active TOTP with the same issuer and account.

If exactly one match exists, offer:

- **Update existing TOTP**
- **Add as another TOTP**
- Cancel

`Update existing TOTP` should be the preferred/default action when the UI needs a default.

Updating modifies the existing logical token rather than deleting it and creating a new unrelated token.

If multiple active matches exist, Totipo must not guess. The user chooses which existing TOTP to update or chooses to add another.

Matching is based on user-visible identity, not on secret comparison.

### 6.8 After commit

After a successful add:

1. Return to the main vault view.
2. Bring the new token into view where possible.
3. Leave the code concealed.
4. Make `Show Code` the obvious next action.

Totipo should not automatically leave newly added codes visible.

---

## 7. Main vault view

The main vault view exists to get the user from an open vault to the required TOTP quickly.

Conceptually:

```text
Vault identity / navigation

Search                                      Add TOTP

Token list
    TokenRow
    TokenRow
    TokenRow

Quiet utility/status
```

### 7.1 Main hierarchy

Normal usage order:

1. Search or visually find a token.
2. Show Code.
3. Copy/use the code.
4. Add TOTP when needed.
5. Use management and exceptional actions only when required.

### 7.2 Search

Search is navigation within the current vault, not a separate results screen.

Search only user-facing identity fields:

- issuer;
- account.

Do not search:

- secret;
- current TOTP code;
- algorithm;
- period;
- internal object identifiers;
- historical protocol metadata.

> Searching the vault must not require inspection of credential secret material.

For v0, matching should be:

- case-insensitive;
- substring-based;
- whitespace-tolerant;
- immediate while typing.

For multiple search terms, all terms should match somewhere across issuer/account.

No query language, fuzzy ranking, or advanced search syntax is required for v0.

### 7.3 Search results

Filtering hides nonmatching rows in place.

If no result matches:

```text
No TOTPs match "foobar".

Clear search
```

This must be distinct from an actually empty vault.

An empty vault should instead encourage the collection action:

```text
No TOTPs yet.

Add your first TOTP to this vault.

Add TOTP
```

### 7.4 Search and reveal state

Search filtering is independent of TokenRow reveal state.

If a revealed token is filtered out, its reveal timer continues normally. If the filter is cleared before disclosure expires, the token may reappear in its still-valid revealed state.

Search text must never match the currently displayed TOTP code.

### 7.5 Search keyboard behavior

Desktop should support a conventional shortcut such as `Ctrl+F` / `Cmd+F` to focus search.

When search has text, Escape should clear the active search before performing broader navigation behavior.

### 7.6 Ordering

The default user-facing order should be stable and based on displayed identity rather than protocol/storage order.

For v0, sort case-insensitively by:

1. displayed primary identity;
2. displayed secondary identity.

User-configurable sorting is not required for v0.

### 7.7 Add TOTP placement

`Add TOTP` should be a persistent, prominent collection-level action.

Desktop may place a textual `Add TOTP` button near search.

Android may use an appropriate platform-native persistent action such as an app-bar action or floating action button.

The shared requirement is semantic prominence, not identical geometry.

### 7.8 Refresh

Refresh is a utility action and must not compete visually with `Add TOTP`.

Desktop may expose Refresh through the Vault menu, a shortcut, or a quiet toolbar control if manual refresh is operationally useful.

### 7.9 Vault identity

The full filesystem path should not dominate the normal main window.

Prefer a concise vault identity such as the vault directory/display name, with the full path available in vault details or another secondary surface.

If ambiguity exists between identically named vaults, provide sufficient disambiguation without turning the full path into the primary visual identity.

### 7.10 After Add or Update

If the resulting token matches the active search, bring it into view.

If it does not match the current filter, do not silently clear the user's search. Instead, provide a subtle indication that the newly added/updated TOTP is hidden by the current search, with an option to clear the search.

---

## 8. Edit and delete TOTP

### 8.1 Edit

Edit is a secondary management operation.

The normal edit screen should focus on user-facing identity:

```text
Edit TOTP

Service
[ GitHub                         ]

Account
[ niki@example.com               ]

Authenticator setup
SHA1 · 6 digits · 30 seconds
Change setup…

                         Cancel   Save
```

The raw secret should not normally appear as an editable field.

Changing authenticator setup is an explicit sub-flow.

### 8.2 Change setup

`Change setup…` reuses the same acquisition/review model as Add TOTP:

- QR scan where supported;
- setup URI;
- manual entry.

The user reviews the replacement setup before applying it.

Algorithm, digits, and period remain secondary details unless the user explicitly enters or reviews advanced setup.

### 8.3 Delete

Delete should be a direct destructive management action, not a lifecycle radio button buried in Edit.

Conceptually:

```text
Delete TOTP?

GitHub
niki@example.com

This removes the TOTP from the active vault.

Deleting a TOTP does not erase previous versions from vault history.

                    Cancel   Delete TOTP
```

One clear confirmation is sufficient.

Deleted entries disappear from the normal active TOTP list.

The UI must not promise secure erasure.

---

## 9. Conflict presentation and resolution

Conflict is exceptional but should be presented in user terms, not protocol internals.

### 9.1 Alternatives are user-facing versions

A conflict may have multiple causal Heads, and multiple Heads may converge on the same semantic value.

The normal UI presents distinct **Alternatives** as user-facing **versions**.

Heads are provenance behind those versions, not separate choices merely because they are separate lineages.

Conceptually:

```text
Conflict
    Version A
        one or more Heads
    Version B
        one or more Heads
```

The number of Heads supporting a version is not a vote and must not imply preference.

### 9.2 Conflict panel

A conflict panel may look like:

```text
Conflict                                      Resolve…

Two versions of this TOTP need your attention.

GitHub
personal@example.com                         Show Code

GitHub
work@example.com                             Show Code
```

Three-way and larger conflicts use the same model: more versions, not a different interaction.

Versions may support the ordinary bounded `Show Code → Copy` interaction when safe and meaningful.

### 9.3 Simple resolution first

The resolver should first offer complete versions:

```text
Resolve TOTP conflict

Which version should Totipo keep?

○ GitHub
  personal@example.com
  SHA1 · 6 digits · 30 seconds

○ GitHub
  work@example.com
  SHA1 · 6 digits · 30 seconds

○ Combine details…
```

No version is silently preferred.

### 9.4 Versions that differ only in secret material

Distinct Alternatives may have identical visible issuer/account/algorithm/digits/period but different authenticator secrets.

The UI must keep such Alternatives distinguishable without exposing existing secret bytes.

It may say, for example:

```text
GitHub
niki@example.com
SHA1 · 6 digits · 30 seconds
Authenticator key: Version 1
```

and:

```text
GitHub
niki@example.com
SHA1 · 6 digits · 30 seconds
Authenticator key: Version 2
```

`Show Code` may help the user identify which credential is currently valid.

### 9.5 Heads as provenance

Head-level information belongs behind an optional Details affordance.

Useful provenance may include safe client/writer metadata where available.

The UI must not label a version as "newest", "most likely", or "recommended" merely because of Head count or client-reported time.

### 9.6 Combine details

`Combine details…` presents only fields that actually differ.

Agreed fields need no decision.

For example:

```text
Service
GitHub

Account
○ personal@example.com
○ work@example.com
○ Other…

Authenticator setup
○ Setup from Version 1
○ Setup from Version 2
Advanced…
```

The detailed merge UI may expose independently resolvable fields where the protocol/API permits this.

Distinct equal values should be shown once, not once per Head.

### 9.7 Deleted alternatives

A deleted Alternative should be expressed in user terms, for example:

```text
○ Keep GitHub / niki@example.com
○ Keep this TOTP deleted
```

Do not expose protocol terms such as `TOMBSTONE` in normal UI.

### 9.8 Newly arrived conflict information

If new relevant conflict information arrives while the user is resolving:

- do not silently publish a now-partial resolution;
- do not flatten the event into a generic save failure.

Present an explicit state such as:

```text
Another version appeared while you were resolving this conflict.

Review the updated conflict before continuing.

Review Updated Conflict
```

Any advanced partial-resolution path must be explicit.

### 9.9 Publication uncertainty

If Totipo cannot affirm whether a conflict resolution was persisted, do not say simply `Save failed`.

Use wording that communicates uncertainty, for example:

```text
Totipo couldn't confirm whether the resolution was saved.

Refresh or reopen the vault before deciding what to do next.
```

The UI must not automatically retry in a way that obscures the uncertain outcome.

---

## 10. Vault selection

The task is to choose the folder or platform storage location containing a Totipo vault.

The selector should optimize for navigation and selection, not behave like a miniature file manager.

### 10.1 Desktop

Conceptually:

```text
Select Vault Folder

/home/niki/Sources
────────────────────────────────────

../
projects/
totipo-vault/

Selected folder
/home/niki/Sources/totipo-vault

                         Cancel   Select Folder
```

Directories are the primary selectable objects.

The Totipo vault selector should not prominently expose unrelated destructive filesystem operations such as Delete or Rename.

`New Folder` may be available where platform-appropriate, but should remain a utility action.

A generic file filter is not required for v0.

### 10.2 Native chooser policy

Prefer the platform's normal directory-selection experience when it is sufficiently usable.

A custom picker is acceptable where the platform/toolkit chooser provides a materially worse experience.

The shared contract is semantic, not widget-specific.

### 10.3 Selection and validation

Folder selection and vault validation are separate:

```text
Select folder
    ↓
Validate
    ├── valid vault → Unlock
    └── not a vault → Explain
```

Selecting a directory must never implicitly initialize protocol state.

### 10.4 Open and create are separate intents

Opening an existing vault and creating a new vault are separate flows.

An empty folder selected during `Open Existing Vault` must not silently become a new vault.

---

## 11. Startup and remembered vault

Totipo should remember the last successfully opened vault location.

Normal startup:

```text
Launch Totipo
   ↓
remembered vault?
   ├── yes → unlock that vault
   └── no  → Open Existing / Create New
```

Rules:

- remember the vault location, not the password;
- only replace the remembered location after successful vault validation/open;
- a failed attempt to open another folder must not overwrite the last known-good location;
- if the remembered location is unavailable, fall back gracefully to the initial Open/Create flow;
- after a successful Change Vault operation, the new location becomes the remembered location;
- platforms that require durable access grants/bookmarks should persist the platform-appropriate access reference rather than assuming a raw path is enough.

The remembered location is local application metadata and should not become synchronized vault content.

---

## 12. Application shell and vault unlock

Totipo should treat startup, unlock, normal use, and vault-level failure as states of one application shell rather than as unrelated modal workflows.

Conceptually:

```text
Application shell
   │
   ├── NO_VAULT
   │     Open Existing Vault…
   │     Create New Vault…
   │
   ├── LOCKED
   │     remembered vault identity
   │     password unlock
   │     biometric unlock, when enabled and available
   │     Change Vault…
   │
   ├── UNLOCKED
   │     Search
   │     Add TOTP
   │     TokenRows
   │
   └── VAULT_ERROR
         explanatory/recovery state
```

Desktop should normally present these states in the main application window instead of opening a separate password dialog over an otherwise empty shell. Android may use the equivalent platform-native screen structure.

### 12.1 Initial state

If there is no remembered vault:

```text
Totipo

Open an existing vault or create a new one.

Open Existing Vault…
Create New Vault…
```

These are peer entry points.

### 12.2 Locked state

Once a vault is known but locked:

```text
Totipo

totipo-vault
/home/niki/Sources/totipo-vault

Password
[                              ]

                    Change Vault…   Unlock
```

Vault name is primary identity; full path is secondary context.

`Unlock` is primary.

`Change Vault…` is secondary.

`Exit` should not be a task action when normal application/window close behavior already exists.

### 12.3 Password behavior

The password field should:

- receive initial focus when password unlock is the active path;
- obscure input by default;
- support platform-standard password behavior;
- never leak the password into logs or error text;
- submit with Enter when unambiguous.

A platform-standard show-password affordance is acceptable.

### 12.4 Empty password

Opening with an empty password requires explicit confirmation.

The confirming action should clearly describe the exceptional choice, e.g. `Open Anyway`.

### 12.5 Failed unlock vs invalid vault

Wrong/unusable credentials and invalid/corrupt vault state must not collapse into one message.

If Totipo can establish that the vault itself cannot safely be interpreted, do not encourage repeated password retries.

### 12.6 Successful unlock

On success, transition the same application shell to the main vault view with all TOTP codes concealed.

Unlocking must not automatically reveal any TOTP.

---

## 13. Create new vault

Creation is separate from opening an existing vault.

Conceptually:

```text
Create New Vault
      ↓
Choose/create empty vault location
      ↓
Set vault password
      ↓
Create
      ↓
Empty main vault view
```

### 13.1 Location

For v0, create a vault in:

- a newly created directory; or
- an existing empty directory.

Do not overwrite an existing Totipo vault or unrelated files.

If an existing Totipo vault is found, direct the user toward Open Existing Vault.

### 13.2 Password

Collect:

- new password;
- confirmation.

Do not impose arbitrary password composition rules.

Platform-standard reveal-password behavior is acceptable.

### 13.3 Empty password

An empty password is allowed only through explicit confirmation.

The confirming action should say something like `Create Without Password`, not merely `Create`.

### 13.4 Do not expose protocol machinery

Normal vault creation should not ask the user about:

- cryptographic algorithms;
- KDF parameters;
- root/vault keys;
- protocol revision;
- writer/device IDs;
- object padding;
- synchronization internals.

Use supported implementation defaults.

### 13.5 Success

Only enter the main vault view after creation is affirmed.

A newly created vault should land on the empty state:

```text
No TOTPs yet.

Add your first TOTP to this vault.

Add TOTP
```

No additional success dialog is required.

---

## 14. Change vault password

Changing the vault password does **not** re-encrypt vault contents and does **not** rotate the vault encryption key.

The UI must not imply otherwise.

### 14.1 Core explanation

The screen should include user-facing explanatory text such as:

> Changing the password does not re-encrypt your TOTP data. Totipo keeps the vault's encryption key and changes how that key is protected by your password.

A shorter form may be used inline:

> Your new password will protect access to the existing vault encryption key. The contents of the vault are not re-encrypted.

### 14.2 Historical copies

Where relevant, also disclose that old retained copies may remain usable with the old password:

> Previous copies may still exist in backups, synchronization history, or other retained storage and may still be accessible using the old password.

Do not promise revocation or erasure of historical copies.

### 14.3 Interaction

Collect:

- new password;
- confirmation.

Do not impose arbitrary password composition rules.

Empty new password requires explicit confirmation.

The action should be called **Change Vault Password**, not `Re-encrypt Vault`, `Rotate Encryption`, or similar terminology.

Password change should not present fake progress implying that every TOTP object is being rewritten.

---

## 15. Clipboard policy

Copying a TOTP is part of the primary product workflow, but clipboard exposure should be bounded.

### 15.1 Copied value

Copy the canonical credential digits, not their visual grouping.

For example, a visually rendered `213 933` is copied as `213933`.

### 15.2 Clipboard lifetime

A copied TOTP remains eligible for clipboard use only through the validity period of that exact code.

> Totipo should make a best-effort attempt to clear a copied TOTP when that copied code reaches its rollover boundary.

If the row is in late-reveal grace and rolls to a new code, Totipo must not silently replace the clipboard with the new code.

Copying does not extend TokenRow reveal lifetime.

### 15.3 Do not destroy newer clipboard contents

Clipboard cleanup must never blindly overwrite content that the user copied after the TOTP.

Where the platform provides ownership, change tokens, provenance, or another safe mechanism, use it to determine whether the clipboard still contains the value Totipo placed there.

If Totipo cannot safely establish that the clipboard still contains its copied value, leave the clipboard unchanged.

### 15.4 Sensitive clipboard handling

Where a platform provides a supported way to mark clipboard data as sensitive or suppress clipboard previews, Totipo should use it.

Clipboard values must not be copied into logs, notifications, diagnostics, or other incidental persistent surfaces.

Clipboard cleanup is a best-effort privacy measure, not a secure-erasure guarantee.

---

## 16. Locking, background behavior, and biometric unlock

Concealment and locking are separate concepts.

- **Conceal** hides currently displayed TOTP values.
- **Lock** discards the active unlocked-vault session and requires authentication before vault contents can be accessed again.

### 16.1 Background transition

When Totipo leaves the foreground, all displayed TOTP values should be concealed immediately.

Backgrounding does not by itself need to clear a copied TOTP; clipboard lifetime follows the specific copied code as described in §15.

The vault may remain unlocked for a bounded grace period:

```text
FOREGROUND + UNLOCKED
        │
     background
        ▼
BACKGROUND + UNLOCKED
all codes concealed
lock timer running
        │
        ├── return before timeout
        │       ↓
        │   UNLOCKED
        │   codes remain concealed
        │
        └── timeout
                ↓
             LOCKED
```

Returning before the lock timeout must not automatically restore previously revealed codes.

The exact default lock timeout is a product policy value to choose and validate before release. It should not be inferred from the TOTP period.

### 16.2 Android resume behavior

After the vault has locked, returning to Totipo requires authentication.

If biometric unlock is enabled and available, Totipo may invoke the platform biometric prompt. Canceling that prompt leaves Totipo locked and must not create a repeated prompt loop. Password unlock remains available.

If the user returns before the lock timeout, do not prompt for authentication; return to the unlocked main view with all TokenRows concealed.

### 16.3 Biometric unlock

Biometric unlock is a local device convenience. It is not a vault password and is not part of the Totipo synchronization protocol.

Use platform biometric terminology and APIs rather than assuming a particular modality such as fingerprint or face recognition.

Biometric unlock should only be enabled after a successful vault-password unlock.

Local biometric unlock material must be protected using the platform's supported secure authentication/keystore facilities. Totipo must not persist the plaintext vault password merely to replay it after biometric authentication.

Password unlock must remain available as a fallback.

Enabling biometric unlock on one device must not modify synchronized vault state or enable biometric unlock on any other device.

### 16.4 Password-wrapper changes and local biometric unlock

When Totipo observes that the vault's password protection has changed, previously established local biometric unlock material should be treated as stale and must not silently bypass the new password state.

Require a successful unlock using the current vault password before re-establishing local biometric unlock.

This preserves the distinction between a synchronized vault password and device-local convenience authentication.

### 16.5 Desktop locking

Desktop should support explicit Lock.

The application should also lock in response to an operating-system session lock and after the selected Totipo inactivity policy expires.

Minimizing or hiding the desktop window should at least conceal displayed TOTP values. It need not immediately lock the vault unless the product lock policy requires it.

### 16.6 App-switcher and background privacy

On platforms that provide a supported mechanism, Totipo should prevent sensitive unlocked content from appearing in system task/app-switcher previews.

Notifications must not contain TOTP values unless a future feature explicitly designs and reviews such behavior.

---

## 17. Read-only, unavailable, and invalid vault states

Vault-level failure modes should preserve distinctions the implementation can establish. Do not flatten every problem into `Error opening vault`.

### 17.1 Read-only vault

Read-only is a persistent vault-level state, not a reason to block ordinary TOTP retrieval.

When Totipo can safely read the vault but cannot modify it, the user may still:

- search;
- reveal TOTP codes;
- copy TOTP codes;
- view conflict/version details.

Mutation actions are unavailable, including:

- Add TOTP;
- Edit;
- Delete;
- Change Setup;
- Resolve conflict;
- Change Vault Password;
- other operations that require publication.

Present a persistent explanatory state, for example:

```text
Read-only

You can view and copy TOTP codes, but changes cannot be saved to this vault.
```

Important unavailable actions may remain disabled-but-discoverable when that improves understanding, with an accessible explanation of why they are unavailable.

A generic failed write must not automatically be relabeled `read-only`; transient I/O failure and publication uncertainty are different states.

### 17.2 Preserve known distinctions

Examples of distinct user-facing conditions include:

- no Totipo vault at the selected location;
- unsupported vault/protocol version;
- failed unlock;
- required data unavailable;
- invalid/integrity-failing vault data;
- publication uncertainty.

Use the most specific truthful condition available from the API.

### 17.3 Blocking vault-level integrity failure

If Totipo cannot safely establish a usable vault state, do not show a normal-looking token list.

A blocking state may look like:

```text
Totipo can't safely open this vault

Some required vault data is missing or invalid.

Try Again    Choose Another Vault
Details…
```

Diagnostic details must not expose secrets, passwords, or decrypted credential material.

### 17.4 Localized failures

If the API can safely establish the rest of the vault while identifying a localized unavailable or invalid item, keep the failure local rather than blocking the entire vault.

### 17.5 Recovery actions

Do not offer a generic destructive `Repair Vault` action.

Recovery actions must correspond to explicit, protocol/application-supported operations whose consequences Totipo can state truthfully.

---

## 18. Desktop navigation

Desktop should remain a compact utility application centered on one main vault view. A sidebar or multi-destination shell is not required for v0.

A simple menu structure is sufficient, approximately:

```text
File
    Open / Change Vault…
    Exit

Vault
    Lock
    Refresh
    Change Vault Password…
    Vault Information…        optional/later

Token
    Add TOTP…
```

Edit/Delete should primarily remain contextual to the affected TokenRow rather than occupy permanent global UI.

Useful shortcuts may include:

```text
Ctrl/Cmd+F       Search
Ctrl/Cmd+N       Add TOTP
Ctrl/Cmd+R       Refresh, if manual refresh remains meaningful
```

Do not add a global reveal shortcut whose target is ambiguous. Secret disclosure should happen in an explicit token context.

The main desktop view remains approximately:

```text
totipo-vault

Search                                      Add TOTP

TokenRow
TokenRow
ConflictPanel
TokenRow
```

---

## 19. Android navigation

Android should also center on one primary destination: the current vault's TOTP list. Bottom navigation is not required for v0.

A typical unlocked structure may be:

```text
┌────────────────────────────────────┐
│ totipo-vault        Search      ⋮  │
├────────────────────────────────────┤
│                                    │
│ GitHub                  Show Code  │
│ niki@example.com                   │
│                                    │
│ AWS                     Show Code  │
│ work@example.com                   │
│                                    │
│                              [+]   │
└────────────────────────────────────┘
```

The persistent Add action means `Add TOTP`.

Infrequent vault operations belong in the app-bar overflow or an equivalent platform-native surface, for example:

- Lock;
- Refresh;
- Change Vault;
- Change Vault Password.

### 19.1 Search

Search may use the standard Android pattern of temporarily replacing the app bar with a search field rather than permanently consuming vertical space.

### 19.2 TokenRow touch interaction

The main body of a TokenRow may act as a generous Show Code target on Android. Management remains behind an overflow/context affordance.

### 19.3 Add TOTP

The persistent Add action opens the Add TOTP flow. QR scanning should be prominent on Android, with manual setup always available.

Scanning returns to a review step and must not commit immediately.

### 19.4 Conflict resolution

A conflict may remain inline in the main list, while `Resolve` navigates to a full-screen resolver.

Leaving the resolver without publication leaves the conflict unresolved.

### 19.5 System Back

Use normal Android navigation semantics:

- active search → leave/clear search;
- sub-flow → return to the previous screen;
- main vault view → normal Android app/background behavior.

Do not add an `Are you sure you want to exit?` prompt to normal Back behavior.

### 19.6 Settings

A top-level Settings destination is not required for v0 unless genuine user preferences emerge.

Security/interaction policies such as TokenRow reveal lifetime and clipboard lifetime should not become preferences merely because they are configurable in code.

---

## 20. Shared semantic states

Totipo applications should share a small vocabulary of semantic states.

| State | Meaning |
| --- | --- |
| Normal | No exceptional condition |
| Informational | Useful non-problem state |
| Pending | State/result is incomplete |
| Warning | User attention is advisable |
| Conflict | Explicit resolution is required |
| Error | An operation failed |
| Invalid / corrupt | Data cannot safely be interpreted |
| Read-only | Content is visible but modification is unavailable |
| Deleted | Logical deletion/tombstone state |
| Disabled | Action is currently unavailable |
| Publication uncertain | Totipo cannot affirm whether a change persisted |
| Locked | Vault location is known, but authentication is required |

A reusable `StatusPanel` should represent substantial exceptional state.

A reusable `InlineNotice` should represent contextual information inside a form.

Neither should rely on color or iconography alone.

---

## 21. Choice controls

Finite exclusive choices are not ordinary command buttons.

Examples include:

- SHA1 / SHA256 / SHA512
- 6 / 7 / 8 digits

Use a semantic exclusive-choice component.

Requirements:

- selected state is obvious without relying only on color;
- keyboard navigation is appropriate on desktop;
- accessible selection semantics are exposed;
- platform implementations may map to native segmented/choice controls.

---

## 22. Forms and sections

Forms should use a consistent semantic structure:

```text
Section
    Field label
    Control
    Supporting text?
    Validation/error?
```

Desktop may use compact two-column layouts where appropriate.

Android may stack fields vertically.

Labels must not rely solely on placeholder text.

Sections have:

- section title;
- optional description;
- contents.

The exact Swing fieldset appearance is not normative.

---

## 23. Dialog actions

Task dialogs should use a consistent action area.

Typical desktop arrangement:

```text
                              Cancel    Primary Action
```

Long forms may scroll, but the action area should remain readily accessible.

Escape/system Back should normally correspond to the neutral cancellation path unless doing so would discard meaningful work without an appropriate confirmation.

Enter may invoke the primary action only where safe and unambiguous.

---

## 24. Accessibility

At minimum:

- every actionable control has an accessible name;
- icons are labelled or decorative;
- state is never communicated solely by color;
- conflict and error text is available to assistive technologies;
- focus order follows task order;
- increased text size does not truncate security-relevant information;
- desktop remains usable at HiDPI scaling;
- Android uses appropriate touch target sizes;
- ordinary desktop workflows are keyboard accessible.

Accessibility is part of the component contract, not a later polish pass.

---

## 25. Visual tokens

Totipo should use semantic tokens rather than hard-coded ad hoc styling.

A minimal spacing scale may use platform-equivalent values around:

```text
4   micro
8   small
12  compact
16  normal
24  large
32  section
```

Typography roles:

```text
windowTitle
dialogTitle
sectionTitle
body
secondary
fieldLabel
button
status
code
```

Color roles:

```text
surface
surfaceRaised
text
textSecondary
border
accent
focus
warning
danger
success
disabled
```

The exact values may differ by platform and theme.

Use platform/system fonts. Totipo does not require a custom bundled font.

Light and dark environments should both be supported.

---

## 26. Icons

Icons reinforce text rather than replace important concepts.

Important states and actions must remain understandable without interpreting iconography.

Do not use danger/error imagery for neutral actions such as Cancel.

---

## 27. Platform relationship

The design hierarchy is:

```text
Totipo Design Guidelines
    |
    +-- semantic tokens
    +-- states
    +-- components
    +-- terminology
    +-- behavior
    +-- accessibility
           |
       +---+---+
       |       |
    Desktop  Android
     native    native
     norms     norms
```

Android may use Material 3 internally.

Material 3 is not the normative Totipo design system.

Likewise, Swing's current look-and-feel is not normative.

---

## 28. Initial component catalogue

The initial shared component vocabulary is:

```text
Action
    PrimaryAction
    SecondaryAction
    QuietAction
    DestructiveAction

Input
    TextField
    SecretField
    ExclusiveChoice

Structure
    Section
    DialogActions

Status
    InlineNotice
    StatusPanel

Domain
    TokenRow
    TokenAlternativeView
    ConflictPanel
    AddTOTPFlow
    VaultSelector
    VaultUnlock
    LockedVaultView
    ReadOnlyNotice
```

Components should be added only when real screens demonstrate a recurring need.

---

## 29. Not yet fully specified

The following areas still need a dedicated pass before these guidelines should be considered complete:

- the exact default inactivity/background lock timeout and whether any platforms need different defaults;
- detailed accessibility announcements for countdown rollover, copy feedback, lock transitions, and conflict updates;
- visual token values after implementation testing;
- history/deleted-item inspection, if exposed at all;
- device admission/management UX, if exposed to end users;
- exact platform-specific secure-storage implementation for local biometric unlock;
- any explicit recovery workflows supported by future Java/API capabilities beyond the states defined here.

These should be designed from concrete application behavior and API guarantees rather than invented in isolation.
