# Totipo Design Guidelines

**Status:** Draft v0.8 — v0 product and visual-design baseline<br>
**Scope:** Cross-platform product and interaction design<br>
**Protocol status:** Non-normative

These guidelines define Totipo's shared product language across desktop and Android, and are intended to remain applicable to future platforms such as iOS.

They do **not** define the Totipo storage or synchronization protocol and do not override protocol requirements. Platform applications should follow native platform conventions where those conventions do not conflict with the product semantics defined here.

> **Consistency of meaning and behavior is more important than consistency of appearance.**

---

## 1. Product model

Totipo is primarily a TOTP retrieval tool backed by a synchronized vault.

The normal user journey is:

1. Unlock the remembered vault, or select a vault when none is known.
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
   - Add a TOTP.
   - Where the object is unambiguous, the visible control may simply say **Add**.

3. **Secondary management**
   - Edit TOTP identity.
   - Change authenticator setup.
   - Delete TOTP.
   - Change vault.
   - Lock the vault.
   - Refresh when manual refresh is required.
   - Change vault password.
   - View safe diagnostics/details when needed.

4. **Exceptional state handling**
   - Resolve conflicts when the user is ready.
   - Handle invalid/corrupt data.
   - Handle unavailable or otherwise constrained application/vault state.
   - Handle uncertain publication or newly arrived conflict information.

Exceptional operations may become visually prominent when they are relevant, but they should not dominate the normal interface.

Normal user-facing terminology should describe TOTPs rather than implementation tokens. Internal APIs may continue to use token terminology.

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

Conflict, error, unavailable actions, pending state, and similar conditions must be represented textually.

Color, icons, borders, and animation may reinforce meaning but must not carry the meaning alone.

### 2.6 Security without unnecessary friction

Totipo should minimize unnecessary exposure of TOTP values and credential material, but should not make ordinary authentication workflows cumbersome.

Security-sensitive behavior should be explicit and predictable rather than clever or surprising.

### 2.7 Compactly roomy

Totipo should create breathing room through alignment, consistent rhythm, and semantic grouping rather than oversized controls, excessive card padding, or blank filler.

Dense utility surfaces such as the TOTP list may use tighter geometry. Forms and exceptional workflows may use more section separation.

Spend space on hierarchy and clarity, not decoration.

---

## 3. Interaction hierarchy

A task context should have one obvious next action.

Examples:

- Hidden TOTP row: **Show Code**
- Revealed TOTP row: **Copy**
- Vault collection: **Add**
- Add flow: **Add**
- Edit flow: **Save**
- Conflict-resolution flow: **Resolve** or **Save Resolution**
- Vault selection: **Select Folder**
- Locked vault: **Open**
- New vault: **Create Vault**
- Password change: **Change Password**

Use short verbs for controls when the object is unambiguous from context. Use the full noun in titles, explanatory text, confirmations, menus where context may be lost, and places where the object needs to be named. For example, a main-window button may say `Add`, while the resulting screen is titled `Add TOTP`.

Task priority and visual emphasis are related but distinct. A row may have an obvious default action without rendering that action as a filled `PrimaryAction` button.

> A view or dialog should normally have at most one visually primary action. Zero is acceptable when two actions are genuine peers and choosing a visual winner would be arbitrary.

Repeated row-local default actions such as `Show Code` and `Copy` normally use compact neutral/secondary styling so a list of TOTPs does not become a field of competing primary buttons.

Rare actions should not consume permanent visual space merely because they are important when used.

---

## 4. Action roles

### 4.1 Primary action

`PrimaryAction` is the strongest visual action in the current view or dialog. It normally advances or completes that task.

Examples include:

- `Add` in the main unlocked vault view;
- `Save` in Edit TOTP;
- `Resolve` in the whole-version conflict resolver;
- `Save Resolution` in the detailed conflict resolver;
- `Open` in the locked-vault view;
- `Select Folder` in vault selection;
- `Create Vault` in vault creation;
- `Change Password` in the password-change flow.

A view/dialog should normally contain at most one `PrimaryAction`. Primary styling is usually a restrained filled interaction accent with high-contrast text.

Repeated row-local actions do not become `PrimaryAction` merely because they are the default action for that row. `Show Code` and `Copy` normally use a compact neutral/secondary treatment.

A destructive confirmation may have one visually dominant confirming action, but it uses danger styling rather than ordinary accent-primary styling.

### 4.2 Secondary action

A normal visible action that should not dominate the view.

Examples:

- Show Code / Copy in repeated TokenRows;
- Change Vault;
- Refresh;
- Edit when explicitly exposed;
- Cancel in a conventional dialog action row.

Secondary actions normally use neutral surfaces and borders rather than a filled accent.

### 4.3 Quiet action

A low-emphasis action that should not compete with the primary task.

Examples:

- Edit in a dense token row
- Details…
- Clear Search
- ancillary navigation

Quiet does not mean hard to discover or low-contrast.

Cancel is neutral. It must not use danger styling merely because it closes a dialog.

### 4.4 Destructive action

An action with destructive or difficult-to-reverse consequences.

Examples:

- Delete TOTP
- destructive filesystem action, where such an action is intentionally exposed

Danger styling is reserved for destructive/error semantics. A destructive entry action may be restrained until the user reaches the final confirmation; the confirming destructive action may use stronger danger styling.

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
        Edit?
```

A typical hidden desktop row may look like:

```text
GitHub                                      Show Code
niki@example.com                                 Edit
```

A revealed row may look like:

```text
GitHub                       213 933              Copy
niki@example.com              25 sec ◯             Edit
```

Exact geometry is platform-specific. Desktop rows should normally use whitespace and a thin divider rather than a rounded card around every TOTP.

### 5.2 Identity

For a conventional TOTP:

- issuer is the primary label and should have slightly greater typographic weight;
- account is the secondary label and may use lower visual emphasis while remaining comfortably readable.

User data should generally carry more visual weight than labels that describe it.

If only one identity value is present, promote the available value rather than leaving empty visual space.

If neither is available, display a stable fallback such as `Unnamed TOTP`.

### 5.3 Hidden state

TOTPs are concealed by default.

`Show Code` is the default row action. In a repeated list it normally uses compact neutral/secondary styling rather than filled primary styling.

Merely listing, filtering, selecting, focusing, scrolling, refreshing, or opening diagnostics must not derive a TOTP code.

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
25 sec ◯
```

A ring or progress indicator may supplement the text but must never be the only indication of remaining validity.

The numeric text should precede a trailing ring where practical so the ring stays visually stable as the text changes width.

The normal countdown uses the interaction/accent family. When actual remaining duration is **strictly less than 10 seconds**, the indicator uses the warning/amber family. Near expiry is an attention state, not an error; red is reserved for actual error, invalid, or destructive state.

Avoid distracting flashing or alarm-like behavior.

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

If the accepted current code has **strictly less than 10 seconds remaining** in its current TOTP period, the same explicit reveal also authorizes the immediately following period:

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

> A reveal normally lasts through the current TOTP period. If the accepted current code has strictly less than 10 seconds remaining, that same explicit action also authorizes the immediately following period, then Totipo automatically conceals the row.

Exactly 10 seconds does not grant the extra period.

Totipo may prepare the one already-authorized following code before rollover so that the transition can be visually seamless, but the staged code must not be visible, accessible, or copyable before its validity begins. There is never a second automatic rollover without another explicit `Show Code` action.

Copying a code does not extend its reveal lifetime.

The 10-second threshold is a Totipo UX policy. It is not a protocol rule.

### 5.7 Independent rows

Reveal state is tracked independently for every TOTP or semantic conflict version.

> Revealing, copying, rolling over, or concealing one row must not alter another row's reveal state.

Several codes may therefore be visible simultaneously.

Search/filter visibility does not itself revoke an authorized reveal. If a revealed row is filtered out and returns before its authorization expires, it may reappear in its still-valid revealed state. Search must not extend disclosure or derive a code.

### 5.8 Copy

Once a row is revealed, `Copy` is the immediate default action for that row. It normally retains compact neutral/secondary styling rather than becoming a filled primary button.

Copy feedback should be:

- immediate;
- local to the row;
- transient;
- accessible;
- non-modal;
- geometry-preserving.

For example, the button may temporarily change:

```text
Copy → Copied → Copy
```

A separate global toast is not required for ordinary copy confirmation.

Copy must not be disabled merely because the current code is near expiration.

> Copy always copies exactly the valid code currently displayed by the row.

### 5.9 Manual concealment

A dedicated persistent `Hide` control is not required for v0 because disclosure is already bounded by TOTP rollover.

Platforms may provide an unobtrusive explicit conceal action if it improves usability, but it must not clutter the main reveal/copy interaction.

### 5.10 Management actions

Edit and Delete are secondary management operations.

Desktop may expose a compact inline `Edit` action when it remains visually quieter than the row's default `Show Code` / `Copy` action. Android may prefer an overflow/context affordance.

Delete should normally remain behind a management flow and confirmation rather than compete with the normal authentication task.

Management affordances require accessible names that identify the affected TOTP.

### 5.11 Pointer/touch behavior

Android may use a larger portion of the row as the reveal target where that improves touch ergonomics.

Desktop should prefer an explicit `Show Code` control and should avoid revealing credential values merely because identity text was clicked accidentally.

### 5.12 Keyboard behavior

Desktop must provide a coherent keyboard workflow.

At minimum:

- `Ctrl+F` focuses Search and selects an existing query;
- Enter or Down from Search moves focus/selection to a visible result without revealing it;
- Up/Down navigates visible result rows;
- Up from the first visible result returns focus to Search;
- Enter/Space on a hidden focused row invokes `Show Code`;
- Enter/Space on a revealed row with exactly one valid code invokes the same `Copy` behavior as the Copy button;
- a revealed multi-code conflict must not arbitrarily choose one code to copy.

Focus/navigation alone must not derive a code.

Inline buttons retain normal platform button keyboard behavior and must not be double-activated by row-level shortcuts.

### 5.13 Focus, hover, selection, and reveal are distinct

Do not use the same visual treatment ambiguously for:

- pointer hover;
- keyboard focus;
- selection;
- revealed state;
- semantic warning/conflict.

Desktop may maintain a selected row because keyboard navigation, contextual management actions, and diagnostics require one. Selection must remain visually subdued, must not imply reveal, and must not trigger TOTP derivation.

Keyboard focus should be clearly visible and may be stronger than selection.

### 5.14 Conflicts

Conflict is not an ordinary `TokenRow` state.

A conflict belongs to a containing `ConflictGroup` with one child row per distinct semantic Alternative, never one child per causal Head.

Conceptually:

```text
ConflictGroup
    Resolve
    TokenAlternativeRow
    TokenAlternativeRow
```

The conflict may remain unresolved indefinitely. Complete eligible Alternative rows remain independently usable for Show Code / Copy / Edit where safe and meaningful.

---

## 6. Add TOTP

Adding a TOTP is the primary collection-level action.

Where the object is unambiguous, the persistent collection control should normally say **Add**. The resulting screen or flow may be titled **Add TOTP**. Do not use implementation-oriented wording such as `Create Token` in normal UI.

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

The preferred Android path is expected to be QR scanning, with manual entry always available. Android may additionally accept an `otpauth://` URI through appropriate platform sharing/open mechanisms.

Desktop should support local enrollment without requiring Android.

For desktop v0, the preferred acquisition methods are:

- explicit paste/import of an `otpauth://` setup URI;
- manual entry.

Desktop camera scanning is not required for v0.

Totipo must not inspect the clipboard in the background looking for setup material. Clipboard access for URI import must follow an explicit user action.

A credential added on any device is a normal vault change and propagates through Totipo synchronization. No special transfer mechanism is required.

### 6.3 Raw secret presentation

The raw secret is acquisition material, not normal review information.

After QR or URI acquisition, the normal review screen should not redisplay the raw secret unnecessarily.

Manual entry and explicit Change Authenticator Setup expose a secret field because the user must enter a new secret.

Existing stored secrets are never redisplayed merely for review.

### 6.4 Manual entry

Manual setup should collect:

- issuer/service;
- account;
- secret.

Authenticator parameters are visible but visually secondary:

- algorithm;
- digits;
- period.

The common/default configuration should not dominate the form.

Desktop may use full-width exclusive-choice controls for algorithm and digits plus a bounded numeric period control. Android may use equivalent native controls.

### 6.5 Review

The review screen should present the user-meaningful identity and authenticator configuration.

Imported non-default parameters must be visible before commit. Totipo should not silently hide unusual imported configuration.

After QR/URI acquisition, show a setup summary such as:

```text
GitHub
niki@example.com

SHA1 · 6 digits · 30 seconds
```

without unnecessarily redisplaying the raw secret.

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

Instead, after acquisition/review, Totipo checks the current user-visible TOTP set for an existing active TOTP with the same issuer and account.

If exactly one match exists, explicitly offer:

- **Update Existing…**
- **Add Another**
- Cancel

Neither semantic choice should happen silently. If the UI requires a keyboard default, prefer a neutral/cancel path rather than assuming the newly acquired credential replaces the existing one.

`Update Existing…` should review the acquired authenticator setup as a proposed Change Authenticator Setup for the existing logical TOTP before commit.

If multiple active matches exist, Totipo must not guess. The user chooses which existing TOTP to update or chooses to add another.

Matching is based on user-visible identity, not secret comparison.

### 6.8 After commit

After a successful add:

1. Return to the main vault view.
2. Bring the new TOTP into view where possible.
3. Leave the code concealed.
4. Make `Show Code` the obvious next action.

Totipo should not automatically leave newly added codes visible.

---

## 7. Main vault view

The main vault view exists to get the user from an unlocked vault to the required TOTP quickly.

Conceptually:

```text
Vault identity / navigation

Search [................................]   token count   Add

Token list / contextual empty state
    TokenRow
    TokenRow
    ConflictGroup
```

Exact desktop/Android geometry may differ. On desktop, Search should normally absorb spare width while the count and `Add` remain compact trailing controls.

### 7.1 Main hierarchy

Normal usage order:

1. Search or visually find a TOTP.
2. Show Code.
3. Copy/use the code.
4. Add when needed.
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
- digits;
- period;
- internal object identifiers;
- Heads/causal metadata;
- diagnostics text.

> Searching the vault must not require inspection of credential secret material.

For v0, matching should be:

- case-insensitive;
- substring-based;
- whitespace-tolerant;
- immediate while typing.

Split the query on whitespace. Every non-empty term must match somewhere across issuer/account; terms may match different fields. No query language, fuzzy ranking, or advanced search syntax is required.

For a conflict group, match if any semantic Alternative's issuer/account satisfies the query, then keep the complete conflict group visible so the user retains context. Count the group as one logical TOTP.

### 7.3 Search results

Filtering hides nonmatching rows in place.

If no result matches, use the otherwise-empty list area for a contextual empty state:

```text
No TOTPs match "foobar".

Clear Search
```

This must be distinct from an actually empty vault. The normal Search/count/Add header remains visible; the count should communicate filtering, for example `0 of 18`. `Clear Search` is contextual and need not compete visually with the collection-level `Add`.

An empty vault should instead use the list area to encourage the collection action:

```text
No TOTPs yet

Add a TOTP to get started.

Add
```

When the vault is genuinely empty, avoid placing a duplicate visually-primary `Add` immediately beside this call to action. Search may be omitted or de-emphasized because there is nothing to search.

### 7.4 Search and reveal state

Search filtering is independent of reveal state.

If a revealed row is filtered out, its reveal timer continues normally. If the filter is cleared before disclosure expires, the row may reappear in its still-valid revealed state.

Search text must never match the currently displayed code.

Search itself must never derive or extend a TOTP reveal.

### 7.5 Search keyboard behavior

Desktop should support:

- `Ctrl+F` to focus search and select the existing query;
- Enter or Down from Search to move focus/selection into visible results;
- Up from the first visible result to return focus to Search;
- Escape to clear an active search before broader navigation behavior.

No-result Enter/Down should leave focus in Search.

Android may use its platform-standard search/app-bar interaction.

### 7.6 Ordering

The default user-facing order should be stable and based on displayed identity rather than protocol/storage order.

For v0, sort case-insensitively by:

1. displayed primary identity;
2. displayed secondary identity.

Apply the same rule to Alternatives within a ConflictGroup. The ConflictGroup occupies the outer-list position of its first Alternative under that same ordering. This ordering is presentation only and must not imply that the first Alternative is preferred. Exact visible-identity ties may retain a stable presentation order but must not be interpreted as a vote or freshness signal.

User-configurable sorting is not required for v0.

### 7.7 Add placement

The persistent collection action should be prominent but may simply say `Add` where the object is unambiguous.

On desktop, the normal collection header should place `Add` at the trailing/right edge of the Search/count row:

```text
Search [................................]   18 TOTPs   Add
```

Search expands to consume spare width; the count remains secondary; command buttons remain content-sized. `Add` is normally the sole visually primary action in the ordinary populated collection view.

Android may use an appropriate platform-native persistent Add action such as an app-bar action or floating action button.

The shared requirement is semantic prominence, not identical geometry.

### 7.8 Refresh

Refresh is a utility action and must not compete visually with Add.

On desktop, Refresh should not occupy permanent main-view space. Expose it through `Vault → Refresh` and, where useful, the desktop refresh shortcut. Android may follow platform conventions where an explicit refresh control is operationally meaningful.

### 7.9 Vault identity

The full filesystem path should not dominate the normal unlocked main window.

Prefer a concise vault identity such as the vault directory/display name. The full path belongs in the Locked state and in a secondary `About This Vault…` surface.

If ambiguity exists between identically named vaults, provide enough parent-path context to distinguish them without turning the full path into the primary visual identity.

### 7.10 After Add or Update

If the resulting TOTP matches the active search, bring it into view.

If it does not match the current filter, do not silently clear the user's search. Instead, provide a subtle indication that the newly added/updated TOTP is hidden by the current search, with an option to clear the search.

---

## 8. Edit, change setup, and delete TOTP

### 8.1 Edit

Edit is a secondary management operation optimized for the common case.

The normal edit screen focuses on user-facing identity:

```text
Edit TOTP

Issuer
[ GitHub                         ]

Account
[ niki@example.com               ]

Authenticator setup
SHA1 · 6 digits · 30 seconds
Change setup…

Delete TOTP…

                         Cancel   Save
```

Algorithm, digits, period, lifecycle/deletion state, and raw secret should not dominate the ordinary Edit screen.

The setup summary is non-secret. Existing secret material is never displayed.

### 8.2 Change authenticator setup

`Change setup…` is an explicit sub-flow for replacing the credential setup.

It may reuse the same acquisition/review model as Add TOTP:

- QR scan where supported;
- setup URI;
- manual entry.

Manual setup may expose:

- new secret;
- algorithm;
- digits;
- period.

Algorithm/digits/period remain secondary to the actual acquisition task, but unusual imported values must be reviewable before commit.

Changing setup modifies the existing logical TOTP rather than deleting it and creating an unrelated logical token.

### 8.3 Delete

Delete is a direct destructive management action, not a lifecycle/status radio button buried in Edit.

Conceptually:

```text
Delete TOTP?

GitHub
niki@example.com

This removes the TOTP from the active vault.
Previous versions remain in vault history.

                    Cancel   Delete TOTP
```

One clear confirmation is sufficient.

Deleted entries disappear from the normal active TOTP list.

The UI must not promise secure erasure. The internal/protocol status may be `TOMBSTONED`, but ordinary UI should use user-facing terminology such as `Deleted`.

---

## 9. Conflict presentation and resolution

Conflict is exceptional but usable. It should be presented in user terms, not protocol internals, and the user must not be forced to resolve it immediately.

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

### 9.2 Conflict group

A conflict group is expanded by default and keeps each complete semantic version usable:

```text
⚠ Conflict                                      Resolve

    Slie                                     Show Code
    a                                            Edit

    Sile                                     Show Code
    x                                            Edit
```

The group uses one restrained warning/conflict accent. Child versions otherwise reuse normal TokenRow geometry and behavior.

Complete eligible versions may independently support the ordinary bounded `Show Code → Copy` interaction. Editing a child edits that semantic Alternative only and does not resolve the others.

Conflict resolution remains optional. The user may continue using the versions and resolve later when enough information is available.

Search treats the group as one logical TOTP. Matching any child identity keeps the full group visible.

### 9.3 Simple whole-version resolution first

Most conflict resolutions are expected to be "keep this version". `Resolve` should therefore first present the complete semantic versions directly:

```text
Resolve Conflict

Choose the version to keep.

○ Slie
  a
  SHA1 · 6 digits · 30 seconds

○ Sile
  x
  SHA256 · 8 digits · 30 seconds

              Cancel   Combine details…   Resolve
```

No version is silently preferred.

`Resolve` remains available even before a version has been selected. Activating it with no selection performs validation, marks the version choice as incomplete, and moves focus to that choice rather than silently doing nothing or relying on a disabled button to communicate the problem.

Choosing a complete version means keeping that semantic Alternative as the resolved TOTP according to the existing resolution semantics.

`Combine details…` is a quieter secondary path for the uncommon case where no whole version is entirely correct.

### 9.4 Versions that differ only in secret material

Distinct Alternatives may have identical visible issuer/account/algorithm/digits/period but different authenticator secrets.

The UI must keep such Alternatives distinguishable without exposing existing secret bytes.

Use user-facing source identity where possible. If versions remain visually indistinguishable, use a neutral fallback such as `Version 1` / `Version 2`. These labels imply no preference.

`Show Code` may help the user identify which credential is currently valid.

### 9.5 Heads as provenance

Head-level information belongs behind an optional diagnostics/details affordance.

Useful provenance may include safe client/writer metadata where available.

The UI must not label a version as "newest", "most likely", or "recommended" merely because of Head count or client-reported time.

### 9.6 Combine details

`Combine details…` opens a field-by-field resolver. It should not inherit a preference from any version selected on the preceding simple-resolution screen; agreed values may be prefilled, but genuine disagreements begin unresolved.

The detailed resolver uses one visually grouped section per semantic decision.

For a field where all versions agree, show one normal prefilled editable control.

For a field where versions disagree, show one vertical radio row per distinct existing semantic value plus a final editable custom row where the domain permits it.

Example:

```text
┌─ Issuer ───────────────────────────────┐
│ ( ) [ Slie                          ]  │
│ ( ) [ Sile                          ]  │
│ ( ) [                               ]  │
└────────────────────────────────────────┘
```

Existing value fields are read-only/selectable, not disabled-looking. Clicking/focusing an existing row selects its radio. Focusing or typing in the custom field selects the custom radio directly. The user need not first choose a visible `Other` label.

Distinct equal values are shown once, not once per Alternative or Head.

Status is a finite user-facing choice (`Active` / `Deleted`) and has no custom row.

The detailed resolver has an explicit task action area, conceptually:

```text
Back                              Cancel   Save Resolution
```

- `Back` returns to the simple whole-version resolver.
- `Cancel` closes the resolver and leaves the conflict unresolved.
- `Save Resolution` remains available and performs authoritative validation on activation; incomplete decisions are exposed locally and focus/scroll moves to the first actionable problem.
- A detailed draft may be retained when moving Back and re-entering `Combine details…` within the same resolver session.
- Cancelling/closing the resolver discards the detailed draft; it is not durable vault state.
- Returning from the detailed resolver must not silently convert the composed draft into a whole-version selection.

### 9.7 Authenticator setup is an atomic resolution unit

For conflict resolution, the authenticator setup is one indivisible semantic unit:

```text
secret + algorithm + digits + period
```

Do not expose independent existing-value choices that can accidentally synthesize a setup such as secret from Version A + algorithm from Version B.

Two existing setups are equal only when public semantic information proves the same secret-equivalence group, algorithm, digits, and period. Never infer secret equality from coincident current codes.

Existing setup choices never display the secret. Show non-secret setup metadata and source identity where useful:

```text
○ Setup used by Slie · a
  SHA1 · 6 digits · 30 seconds

○ Setup used by Sile · x
  SHA256 · 8 digits · 30 seconds
```

A separate custom setup choice is the explicit path for deliberately constructing a new setup. It exposes a new secret field plus algorithm/digits/period controls. Interacting with any custom setup control selects the custom setup. Abandoned custom secret input must be cleared when the user switches back to an existing setup or cancels.

### 9.8 Deleted and incomplete alternatives

A deleted Alternative should be expressed in user terms such as `Deleted`, not `TOMBSTONED`.

Do not offer Show Code for a deleted/incomplete Alternative when the semantic state is not eligible for TOTP generation.

Do not fabricate missing fields for an incomplete Alternative. Diagnostics may expose technical status/provenance.

### 9.9 Newly arrived conflict information

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

### 9.10 Publication uncertainty

If Totipo cannot affirm whether a conflict resolution was persisted, do not say simply `Save failed`.

Use wording that communicates uncertainty, for example:

```text
Totipo couldn't confirm whether the resolution was saved.

Refresh or reopen the vault before deciding what to do next.
```

Conflict remains usable when resolution is cancelled, fails, or is uncertain.

---

## 10. Vault selection

The task is to choose the folder or platform storage location containing a Totipo vault.

The selector is a temporary platform utility invoked from the persistent application shell. The shell remains the application identity; the chooser must not become a competing Totipo application window.

The selector should optimize for navigation and selection, not behave like a miniature file manager.

### 10.1 Desktop

The desktop existing-vault picker should be deliberately narrower than a general-purpose file manager.

Conceptually:

```text
Select Vault Folder

/home/niki/Sources                                      Up

projects/
totipo-java/
totipo-vault/        ← selected

                                      Cancel   Select Folder
```

Directories are the only ordinary selectable objects.

Interaction:

- single-click selects a displayed child directory;
- `Select Folder` chooses the selected child directory;
- double-click or Enter navigates into the selected child directory;
- when no child is selected, `Select Folder` may choose the directory currently being displayed;
- navigating into another directory clears stale child selection;
- `Up` navigates to the parent directory.

The normal existing-vault picker should omit dot-prefixed/hidden directories by default and should not expose generic filesystem-management operations such as Rename, Delete, New Folder, filename entry, or file-type filters.

The directory list should receive most of the available dialog space and support ordinary keyboard navigation. Totipo does not need a miniature file manager merely to identify an existing vault.

### 10.2 Native chooser policy

Prefer the platform's normal directory-selection experience when it is sufficiently usable.

A custom picker is acceptable where the platform/toolkit chooser provides a materially worse experience.

The shared contract is semantic, not widget-specific.

### 10.3 Selection and validation

Folder selection and vault opening/authentication are separate:

```text
Select folder
    ↓
Inspect target
    ├── recognizable Totipo vault target → make it current/remembered
    │                                  ↓
    │                               LOCKED or blocking state
    └── not a Totipo vault → Explain; keep previous remembered vault
```

A deliberately selected location becomes the current/remembered vault once Totipo can recognize it as a Totipo vault target; successful password authentication is not required. A wrong password therefore leaves the newly selected vault current and `LOCKED` rather than returning to the previous vault.

Likewise, a recognizable selected vault that later proves unsupported, invalid/corrupt, or otherwise blocking remains the current remembered target so Totipo can report and recover against the vault the user actually selected.

A folder that Totipo can establish is not a Totipo vault does not replace the remembered vault.

Selecting a directory must never implicitly initialize protocol state.

### 10.4 Open and create are separate intents

Opening an existing vault and creating a new vault are separate flows.

An empty folder selected through `Select Vault…` must not silently become a new vault.

---

## 11. Startup and remembered vault

Totipo should remember the currently selected recognizable Totipo vault location/access reference. It is not limited to the last vault that was successfully unlocked.

Normal startup:

```text
Launch Totipo
   ↓
remembered vault?
   ├── yes → persistent shell starts on that vault (normally LOCKED)
   └── no  → persistent shell starts NO_VAULT
```

Rules:

- remember the vault location/access reference, not the password;
- accepting a location that Totipo can recognize as a Totipo vault target makes it current and remembered without requiring successful password authentication;
- a wrong password leaves that selected vault current/remembered and keeps the shell `LOCKED`;
- a recognizable vault that enters an unsupported/invalid/corrupt blocking state remains current/remembered;
- selecting a location that Totipo can establish is not a Totipo vault does not replace the remembered location;
- cancelling the chooser before accepting a replacement does not replace the remembered location;
- successful vault creation makes the new vault current and remembered;
- if the remembered location is unavailable at startup, fall back gracefully to `NO_VAULT` with a concise explanation;
- platforms that require durable access grants/bookmarks should persist the platform-appropriate access reference rather than assuming a raw path is enough.

The remembered location is local application metadata and should not become synchronized vault content.

---

## 12. Persistent application shell and vault unlock

Totipo should treat startup, unlock, normal use, lock, and blocking vault failure as states of one persistent application shell rather than separate Totipo application windows.

Conceptually:

```text
ApplicationShell
   │
   ├── NO_VAULT
   │     Choose a vault to continue
   │     Select Vault…
   │     Create New Vault…
   │
   ├── LOCKED
   │     remembered vault identity
   │     password unlock
   │     biometric/device unlock where supported
   │     Change Vault…
   │
   ├── UNLOCKED
   │     Search
   │     Add
   │     TokenList
   │
   └── BLOCKING_VAULT_STATE
         explanatory/recovery state
```

Desktop should normally keep one persistent Totipo application window and replace its content as shell state changes. Android may use equivalent platform-native screens/content states.

Platform utility/modal surfaces such as a filesystem chooser, confirmation, editor, or conflict resolver may temporarily overlay the shell. Totipo should not maintain competing application-content windows or multiple unlocked vault sessions.

### 12.1 No-vault state

If there is no remembered usable vault, present a sparse first-run/return state rather than a form anchored to a corner.

Desktop should use a centered task composition with a modest upward bias:

```text
Choose a vault to continue

      Select Vault

   Create New Vault…
```

`Select Vault` is the visually primary path. `Create New Vault…` is secondary and may appear beneath it. The heading and action group should read as one centered first-impression state.

This is an intentional exception to the normal trailing/right-aligned desktop form-action rule: `NO_VAULT` is an empty/welcome state, not an ordinary data-entry form.

Other platforms may adapt the geometry while preserving the same hierarchy and wording intent.

### 12.2 Locked state

Once a vault is known but locked:

```text
Totipo

totipo-vault
/home/niki/Sources/totipo-vault

Password
[                              ]

                    Change Vault…   Open
```

Vault name is primary identity; full path/access context is secondary confirmation.

`Open` is the visually primary action. `Change Vault…` is secondary.

The shell itself provides ordinary application close/Exit behavior; `Exit` need not appear as an unlock-task button.

### 12.3 Password behavior

The password field should:

- receive initial focus when password unlock is active;
- obscure input by default;
- support platform-standard password behavior;
- never leak the password into logs or error text;
- submit with Enter when unambiguous.

A platform-standard show-password affordance is acceptable.

Failed credentials keep the shell in `LOCKED`, display a local error, and allow another attempt without closing/reopening an unlock window.

### 12.4 Empty password

Opening with an empty password requires explicit confirmation.

The confirming action should clearly describe the exceptional choice, e.g. `Open Anyway`.

### 12.5 Failed unlock vs invalid vault

Wrong/unusable credentials and invalid/corrupt vault state must not collapse into one message.

If Totipo can establish that the vault itself cannot safely be interpreted, transition to a blocking vault state rather than encouraging repeated password retries.

### 12.6 Successful unlock

On success, transition the same application shell to `UNLOCKED` with all TOTP codes concealed.

Unlocking must not automatically reveal any TOTP.

### 12.7 Change Vault

Changing vault is a security boundary, not an unlocked fallback.

From an unlocked vault:

1. retire/close the current unlocked session and its reveal/clipboard/session-owned state immediately;
2. remain in the persistent shell;
3. invoke the platform vault selector;
4. if the user accepts a location recognizable as a Totipo vault target, that location becomes current/remembered and the shell proceeds with that vault in `LOCKED` or an appropriate blocking state;
5. password failure on the newly selected vault leaves that vault current and `LOCKED`.

Cancellation is deterministic:

- cancelling the chooser before accepting a replacement returns to the previous vault in `LOCKED` state when that vault/location is still available;
- if there was no previous vault, cancellation returns to `NO_VAULT`;
- if the previous location has become unavailable, fall back to `NO_VAULT` with a concise explanation.

Cancelling selection never resurrects the retired unlocked session.

At most one unlocked vault session exists at any time.

---

## 13. Create new vault

Creation is separate from opening an existing vault and occurs from the persistent application shell.

Conceptually:

```text
NO_VAULT / LOCKED shell
      ↓
Create New Vault…
      ↓
Choose/create empty vault location
      ↓
Set vault password
      ↓
Create
      ↓
UNLOCKED empty vault view
```

### 13.1 Location

For v0, create a vault in:

- a newly created directory; or
- an existing empty directory.

Do not overwrite an existing Totipo vault or unrelated files.

If an existing Totipo vault is found, direct the user toward selecting/opening it instead.

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

Only transition the shell to the unlocked main vault view after creation is affirmed.

A newly created vault should land on the empty state:

```text
No TOTPs yet.

Add your first TOTP to this vault.

Add
```

No additional success dialog is required.

---

## 14. Change vault password

Changing the vault password does **not** re-encrypt vault contents and does **not** rotate the vault encryption key.

These semantics remain unchanged; the normal UI uses consequence-focused wording rather than implementation mechanics.

Changing the vault password requires successful reauthentication with the **current vault password**. An already-unlocked session alone does not authorize replacement of the credential controlling future vault access. Ordinary vault use and Add/Edit/Delete/Resolve operations do not require this password reauthentication.

### 14.1 User-facing consequence

The normal Change Vault Password screen should stay concise:

```text
Change Vault Password

Current password
[                         ]

New password
[                         ]

Confirm new password
[                         ]

Old backups or retained copies may still be accessible with the previous password.

                 Cancel   Change Password
```

Do not require the user to read an implementation explanation about vault/root keys in order to change the password.

### 14.2 Technical truth and documentation

Documentation/details may explain that Totipo keeps the same vault encryption key/root and rewraps it with the new password. Vault contents are not re-encrypted merely because the password wrapper changes, and the vault root is not rotated.

Previous copies may still exist in backups, synchronization history, or other retained storage and may still be accessible using the old password.

Do not promise revocation or secure erasure of historical copies.

Do not present fake progress implying that every TOTP object is being rewritten.

### 14.3 Interaction

Collect:

- current vault password;
- new password;
- new-password confirmation.

Current-password input is always required, even when the vault is already unlocked. The application should not retain the original unlock password merely to avoid asking for it here.

Do not impose arbitrary password composition rules. Platform-standard reveal-password behavior is acceptable.

`Change Password` remains activatable when input is incomplete or invalid. Activation performs validation and focuses the first actionable problem. A wrong current password is a local authentication failure and must not imply vault corruption.

Empty new password requires explicit confirmation.

The action is called **Change Vault Password** / **Change Password**, not `Re-encrypt Vault`, `Rotate Encryption`, or similar terminology.

### 14.4 Availability

Change Vault Password is available only for an unlocked vault when the current platform/session supports password-wrapper replacement.

If the implementation reports that this operation is genuinely unavailable, the UI may disable it with an accessible explanation. Do not infer that the operation is unavailable from an arbitrary prior failed write.

---

## 15. Clipboard policy

Copying a TOTP is part of the primary product workflow, but clipboard exposure should be bounded.

### 15.1 Copied value

Copy the canonical credential digits, not their visual grouping.

For example, a visually rendered `213 933` is copied as `213933`.

### 15.2 Clipboard lifetime

Clipboard exposure has an independent upper bound from on-screen reveal state.

> Totipo should make a best-effort attempt to clear a copied TOTP at the earlier of (a) that exact code's validity boundary or (b) 30 seconds after the Copy action.

Examples:

- a code copied with 22 seconds of validity remaining is eligible for cleanup in about 22 seconds;
- a code copied with 48 seconds remaining is eligible for cleanup after 30 seconds.

If the row is in late-reveal grace and rolls to a new code, Totipo must not silently replace the clipboard with the new code.

Copying does not extend TokenRow reveal lifetime. Each explicit Copy applies the cleanup policy to the value Totipo placed on the clipboard for that action.

### 15.3 Do not destroy newer clipboard contents

Clipboard cleanup must never blindly overwrite content that the user copied after the TOTP.

Where the platform provides ownership, change tokens, provenance, or another safe mechanism, use it to determine whether the clipboard still contains the value Totipo placed there.

If Totipo cannot safely establish that the clipboard still contains its copied value, leave the clipboard unchanged.

### 15.4 Sensitive clipboard handling

Where a platform provides a supported way to mark clipboard data as sensitive or suppress clipboard previews, Totipo should use it.

Clipboard values must not be copied into logs, notifications, diagnostics, or other incidental persistent surfaces.

Clipboard cleanup is a best-effort privacy measure, not a secure-erasure guarantee.

---

## 16. Locking, inactivity, background behavior, and biometric unlock

Concealment, locking, and application visibility are separate concepts.

- **Conceal** hides a currently displayed TOTP value.
- **Lock** discards the active unlocked-vault session and returns the persistent shell to `LOCKED`.
- **Background/minimized** describes application visibility and does not by itself change reveal state.

TokenRow disclosure remains governed by the bounded reveal rules in §5. Moving Totipo into the background, behind another window, or into a minimized state must not invent a second concealment timer.

### 16.1 Platform policy summary

For v0, use these fixed policies:

| Behavior | Desktop | Android |
| --- | --- | --- |
| TOTP reveal | Bounded by TokenRow rollover rules | Same |
| Focus/background/minimize | No forced concealment or lock | No forced concealment or immediate lock |
| Auto-lock | 15 minutes of Totipo inactivity | 5 minutes of Totipo inactivity |
| OS/device lock | Lock Totipo | Lock Totipo |
| Suspend/sleep | Lock on resume | Follow platform lifecycle; require unlock when app/session has locked |
| Explicit lock | `Vault → Lock`, `Ctrl+L` | App-bar/overflow action where useful |
| Local biometric unlock | Optional/future | Supported where platform facilities are suitable |
| Task/app-switcher preview | Platform-appropriate privacy | Protect sensitive preview |
| User-configurable security timing | None in v0 | None in v0 |

The desktop and Android timeout values are product policy, not vault metadata and not synchronized settings.

### 16.2 What counts as activity

Auto-lock timers are reset only by meaningful direct user interaction with Totipo.

Examples include:

- searching;
- revealing or copying a TOTP;
- opening or using a menu;
- navigating a Totipo flow;
- editing, adding, resolving, or otherwise interacting with vault UI.

The following do **not** reset the inactivity timer:

- countdown repaint/ticks;
- filesystem observation;
- synchronization;
- background refresh;
- automatic UI updates;
- other non-user work.

Use Totipo-local interaction rather than attempting to infer system-wide keyboard or pointer activity.

### 16.3 Desktop locking

Desktop v0 locks after **15 minutes without user interaction with Totipo**.

Desktop also locks:

- on explicit `Vault → Lock` / `Ctrl+L`;
- when the operating-system session locks;
- after suspend/sleep, on resume;
- when the application exits and its unlocked state is discarded.

The following do not by themselves lock or conceal TokenRows:

- ordinary focus loss;
- another window covering Totipo;
- minimizing or hiding the Totipo window.

Any revealed TOTP continues its normal bounded disclosure lifetime and self-conceals according to §5 until lock actually occurs.

Locking retires the active vault session, clears/retires session-owned reveal and clipboard state, closes or safely retires sensitive child flows, and changes the same application shell to `LOCKED`. A subsequent open/unlock requires the vault password or an approved local-device unlock mechanism.

Explicit Lock always wins over an owned Totipo child flow. `Vault → Lock` / `Ctrl+L` while Add, Edit, Change Setup, Resolve, Change Password, or another sensitive child flow is open must retire/cancel that flow, promptly clear temporary secret/password material it owns, and transition the shell to `LOCKED` without an unsaved-draft confirmation. Reopening after authentication starts from ordinary vault state rather than resurrecting the abandoned child flow.

### 16.4 Android locking

Android v0 locks after **5 minutes without user interaction with Totipo**, whether the application remains foregrounded or has moved into the background.

Moving to the background does not itself immediately lock the vault and does not alter reveal state. Revealed codes continue their normal bounded lifetime.

If the user returns before the 5-minute timeout and the device/session has not otherwise locked, Totipo returns directly to the unlocked application state without an authentication prompt.

If the timeout has elapsed, or the device/session has locked, Totipo returns in the `LOCKED` state and requires authentication.

Background work, synchronization, countdown updates, and lifecycle callbacks do not count as user activity.

### 16.5 Android task/app-switcher privacy

Logical TokenRow reveal state and operating-system previews are separate concerns.

Where Android provides a supported mechanism, Totipo should prevent sensitive unlocked content and visible TOTP codes from being retained or shown in task/app-switcher previews.

Protecting the preview must not require changing the underlying TokenRow reveal state.

Notifications must not contain TOTP values unless a future feature explicitly designs and reviews such behavior.

### 16.6 Biometric/device unlock

Biometric/device unlock is a local device convenience for regaining access to an already-known vault. It is not a vault password and is not part of the Totipo synchronization protocol.

In v0.8, local biometric/device authentication does not replace entry of the current vault password when changing the vault password (§14).

Use platform biometric/device-authentication terminology and APIs rather than assuming a particular modality such as fingerprint or face recognition.

Android may store device-protected local material sufficient to regain access to the already-known vault after approved local authentication, while retaining password fallback.

Local biometric/device-unlock material is device-local security state, not synchronized vault content.

### 16.7 Password-wrapper changes and local biometric unlock

After Change Vault Password succeeds, any device-local unlock wrapper/material that depends on the previous password-wrapped vault state must be invalidated or re-established according to the platform security design.

Do not allow stale local biometric/device-unlock material to silently bypass a changed vault wrapper.

---

## 17. Unavailable, invalid, and constrained vault states

Vault-level failure modes should preserve distinctions the implementation can establish. Do not flatten every problem into `Error opening vault`.

The persistent shell should represent these conditions without pretending all of them are the same generic `VAULT_ERROR`.

### 17.1 No inferred read-only vault state

Totipo v0 does not define an automatic global Read-only vault state or provide an `Open Read-Only` mode.

A failed write, definite publication failure, publication uncertainty, invalid/corrupt data, or unavailable data must not be relabeled as Read-only. A filesystem/write failure alone does not establish that the entire vault cannot be modified.

An explicit user-selected read-only session or an affirmatively reported backend/storage read-only capability would require a separate product design driven by a concrete requirement. Neither is part of v0.8.

### 17.2 Preserve known distinctions

Examples of distinct user-facing conditions include:

- no Totipo vault at the selected location;
- failed unlock/wrong password;
- unsupported vault/protocol version;
- required data unavailable;
- invalid/integrity-failing vault data;
- publication uncertainty.

Use the most specific truthful condition available from the API.

### 17.3 Blocking vault state

If Totipo cannot safely establish a usable vault state, do not show a normal-looking token list.

A blocking state may look like:

```text
Totipo can't safely open this vault

Some required vault data is missing or invalid.

Try Again    Change Vault…
Details…
```

Blocking is a shell content state. A definite publication failure or publication uncertainty alone does not establish a blocking vault condition. Diagnostic details must not expose secrets, passwords, or decrypted credential material.

### 17.4 Localized failures

If the API can safely establish the rest of the vault while identifying a localized unavailable or invalid item, keep the failure local rather than blocking the entire vault.

### 17.5 Publication uncertainty

A definite publication failure belongs to the operation that failed. Publication uncertainty is also operation-specific: Totipo cannot affirm whether a requested change persisted.

Do not turn one uncertain publication into a vault-global failure or blocking condition. Present the uncertainty truthfully and use observation/refresh/reopen to establish current state.

### 17.6 Recovery actions

Do not offer a generic destructive `Repair Vault` action.

Recovery actions must correspond to explicit, protocol/application-supported operations whose consequences Totipo can state truthfully.

---

## 18. Desktop navigation

Desktop should remain a compact utility application centered on one persistent shell and one main vault view. A sidebar or multi-destination shell is not required for v0.

A simple menu structure is sufficient, approximately:

```text
File
    Exit

Vault
    Select / Change Vault…
    Lock
    Refresh
    Change Vault Password…
    About This Vault…

Token
    Add…
    Edit…                 contextual/selected item
    Delete…               contextual/selected item
    View Diagnostics…     contextual/selected item
```

Edit/Delete should primarily remain contextual to the affected row rather than occupy permanent global UI.

Useful shortcuts include:

```text
Ctrl+F       Search
Ctrl+L       Lock
Ctrl+N       Add
Ctrl+R       Refresh, if manual refresh remains meaningful
```

Do not add a global reveal shortcut whose target is ambiguous. Secret disclosure happens in an explicit row context.

The unlocked desktop view remains approximately:

```text
totipo-vault

Search [................................]   token count   Add

TokenRow
TokenRow
ConflictGroup
TokenRow
```

Refresh remains in the Vault menu rather than occupying permanent collection space. The exact path does not dominate this view.

`About This Vault…` is a secondary surface for safe vault-level information such as:

- full filesystem/storage location;
- storage/access capability where the application can truthfully report it;
- vault/protocol format/version where useful;
- safe vault-level diagnostics/details.

Do not infer writable/read-only status from an arbitrary operation failure. When the platform/API does not expose storage/access capability, omit the field or explain that the capability is not reported. Such details do not define a global vault mode in v0.8.

Show vault/protocol format metadata only where actually available; do not substitute the installed library version for unknown vault-format metadata.

A future application-level `Help → About Totipo` remains conceptually separate from `About This Vault…`.

---

## 19. Android navigation

Android also centers on one primary destination: the current vault's TOTP list. Bottom navigation is not required for v0.

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

The persistent Add action means add a TOTP; the visible label/icon may simply be `Add` / `+` because no other object type is addable.

Infrequent vault operations belong in the app-bar overflow or an equivalent platform-native surface, for example:

- Lock;
- Refresh;
- Change Vault;
- Change Vault Password;
- About This Vault.

### 19.1 Search

Search may use the standard Android pattern of temporarily replacing the app bar with a search field rather than permanently consuming vertical space.

Matching semantics remain issuer/account-only and follow §7.

### 19.2 TokenRow touch interaction

The main body of a TokenRow may act as a generous Show Code target on Android. Management remains behind an overflow/context affordance.

### 19.3 Add

The persistent Add action opens the Add TOTP flow. QR scanning should be prominent on Android, with manual setup always available.

Scanning returns to a review step and must not commit immediately.

### 19.4 Conflict resolution

A conflict may remain inline in the main list, while `Resolve` navigates to the simple whole-version resolver. `Combine details…` may then navigate to the field-level resolver.

Leaving either resolver without publication leaves the conflict unresolved and usable.

### 19.5 System Back

Use normal Android navigation semantics:

- active search → leave/clear search;
- sub-flow → return to the previous screen;
- main vault view → normal Android app/background behavior.

Do not add an `Are you sure you want to exit?` prompt to normal Back behavior.

### 19.6 Settings and persisted local state

A top-level Settings destination is not required for v0.

Security/interaction policies such as:

- TokenRow reveal lifetime;
- clipboard lifetime;
- desktop inactivity timeout;
- Android inactivity timeout;

are fixed Totipo v0 policy and are not user preferences.

The only preference-like application state required by this design is the currently selected recognizable Totipo vault location/access reference described in §11.

Device-local biometric unlock state is security material/capability state, not a synchronized user preference.

---

## 20. Shared semantic states

Totipo applications should share a small vocabulary of semantic states.

| State | Meaning |
| --- | --- |
| Normal | No exceptional condition |
| Informational | Useful non-problem state |
| Pending | State/result is incomplete or in progress |
| Warning | User attention is advisable |
| Conflict | Multiple semantic versions exist; the user may continue using them and resolve when ready |
| Error | An operation failed |
| Invalid / corrupt | Data cannot safely be interpreted |
| Deleted | Logical deletion/tombstone state |
| Disabled | Action is currently unavailable |
| Publication uncertain | Totipo cannot affirm whether a change persisted |
| Locked | Vault location is known, but authentication is required |

A reusable `StatusPanel` should represent substantial exceptional/application state.

A reusable `InlineNotice` should represent contextual information inside a task/form.

A `FieldMessage` should represent help/validation tied to one field.

None of these should rely on color or iconography alone.

Locked and Deleted are not errors merely because they are non-Normal states.

---

## 21. Choice controls

Finite exclusive choices are not ordinary command buttons.

Examples include:

- SHA1 / SHA256 / SHA512;
- 6 / 7 / 8 digits;
- Active / Deleted;
- one existing conflict value vs another vs a custom value.

Use the semantic control that matches the geometry and meaning:

- **ExclusiveChoice** for compact horizontal finite choices such as algorithm/digits;
- **RadioChoice** for vertically stacked semantic decisions and conflict resolution.

Requirements:

- selected state is obvious without relying only on color;
- exactly one value is selected when the domain requires one;
- keyboard navigation is appropriate on desktop;
- accessible selection semantics are exposed;
- platform implementations may map to native segmented/toggle/radio controls.

For horizontal field-like choice groups on desktop, choices should normally share the available control-column width rather than form a tiny button island.

---

## 22. Forms, sections, and density

Forms should use a consistent semantic structure:

```text
Section
    Field label
    Control
    Supporting text?
    Validation/error?
```

Labels must not rely solely on placeholder text.

### 22.1 Desktop form geometry

Desktop normally uses a compact two-column form where appropriate:

```text
intrinsic label column
16 logical-unit gap
expanding control column
```

Field-like controls in the same form should share left/right edges. Use about 12 logical units between ordinary form rows and 16–24 between semantic sections.

Desktop forms should feel compactly roomy: alignment and rhythm create breathing room; controls should not be made unnecessarily tall merely to create space.

### 22.2 Android form geometry

Android may stack fields vertically and follows platform touch-target conventions while preserving shared semantic grouping.

### 22.3 Sections and FieldGroups

Sections have:

- section title;
- optional description;
- contents.

A `FieldGroup` may use a lightweight title/border for a single semantic decision such as conflict-resolution Issuer or Authenticator Setup.

A typical desktop FieldGroup uses about 12 logical units internal padding and about 16 between sibling groups. The border/title provides grouping; do not add heavy nested cards around every field.

### 22.4 Scrolling

Scrolling is a fallback for limited space or genuinely long content, not a substitute for a useful default window/dialog size.

Prefer:

- one outer vertical scroller for a long form;
- no normal horizontal scrolling;
- action rows outside the scrolling body where practical;
- no nested per-section scroll panes.

Focus should scroll a newly focused control into view where needed.

### 22.5 Stable geometry

Conditional state should avoid unnecessary layout jumps. Where a field is likely to appear/disappear during an ordinary interaction, prefer stable geometry with enabled/disabled state when that preserves clarity and does not expose sensitive material.

---

## 23. Dialog and task actions

Task dialogs/forms should use a consistent action area.

Typical desktop arrangement:

```text
                              Cancel    Primary Action
```

Use about 8 logical units between related action buttons and 16–24 logical units between the task body and action row.

Long forms may scroll, but the action area should remain readily accessible.

A view/dialog should normally have at most one visually primary action. When actions are genuine peers, use ordinary/secondary treatment rather than inventing a primary winner.

### 23.1 Validation-on-activation

Do not disable a form's commit/submit action merely because required input is incomplete or current values are invalid.

Activating the action performs authoritative validation and, when validation fails:

- does not publish/commit;
- exposes specific local validation messages;
- scrolls the first actionable problem into view when necessary;
- moves focus to that problem or its semantic choice group.

Opportunistic validation while editing is allowed, but users must never have to infer what is wrong solely from a disabled primary button.

This rule is about form validity. An action may still be unavailable because the underlying application state genuinely forbids it (for example, a retired/locked session or a genuinely unavailable required operation/capability). Such unavailability must have an accessible explanation rather than relying on a mysteriously disabled control. After a valid submission has started, a transient busy/in-progress state may prevent duplicate submission.

Escape/system Back should normally correspond to the neutral cancellation path unless doing so would discard meaningful work without an appropriate confirmation.

Enter may invoke the primary action only where safe and unambiguous; if invoked on an invalid form, it follows the same validation-on-activation behavior.

Actions should use content-driven widths plus platform-appropriate padding rather than stretching command buttons across large empty areas.

---

## 24. Accessibility

Accessibility is part of the component contract, not a later polish pass.

At minimum:

- every actionable control has an accessible name;
- icons are labelled or decorative;
- state is never communicated solely by color;
- conflict and error text is available to assistive technologies;
- focus order follows task order;
- increased text size does not truncate security-relevant information;
- desktop remains usable at HiDPI scaling;
- Android uses appropriate touch target sizes;
- ordinary desktop workflows are keyboard accessible;
- selection/focus/reveal remain semantically distinct;
- read-only values remain fully readable and are not presented as disabled controls.

Contrast targets across light and dark themes:

```text
normal-sized text                         >= 4.5:1
large text                               >= 3:1
text on filled actions                   >= 4.5:1
focus indicators                         >= 3:1 against adjacent surface
meaningful boundaries/graphics           >= 3:1 where needed to perceive state
```

Disabled controls may be lower-emphasis, but Totipo should still aim for approximately 3:1 legibility where practical so users can understand what is unavailable.

Secondary text must remain comfortably readable. Lower emphasis primarily through tone and hierarchy, not tiny type or extremely low contrast.

---

## 25. Visual system

Totipo uses semantic tokens and component-state contracts rather than hard-coded ad hoc styling.

The overall character is **neutral-first with balanced, purposeful accents**:

- neutral surfaces and high-contrast text form the visual foundation;
- restrained blue identifies interaction, focus, primary action, and normal countdown state;
- amber identifies warning, semantic conflict, and near-expiry attention;
- red identifies destructive, error, and invalid/corrupt state;
- green is reserved for meaningful persistent success where success color adds information;
- color clarifies hierarchy and state rather than decorating ordinary content.

Dark mode is not a more saturated version of light mode, and light mode is not a colorless version of dark mode. Preserve semantic relationships across themes.

### 25.1 Spacing scale

Use this logical, platform-scaled rhythm:

```text
4   micro
8   tight
12  compact
16  normal
24  section
32  major / uncommon
```

Typical use:

- `4`: icon↔text, ring↔seconds, very closely coupled inline pieces;
- `8`: buttons in one group, radio↔value field, label↔helper text;
- `12`: token-row internal rhythm, form-row gap, choices inside one FieldGroup;
- `16`: sibling groups, toolbar/search/list separation;
- `24`: dialog content edge, major form-section separation;
- `32`: rare major structural break.

Most utility UI should live in the 8–24 range. Do not create breathing room with blank filler.

### 25.2 Typography

Use platform/system fonts. Totipo does not require a bundled custom font.

Roles are relative to the platform UI base font:

```text
body            1.00x regular
secondary       ~0.95x regular
fieldLabel      1.00x regular/medium
sectionTitle    ~1.10x semibold
screenTitle     ~1.20–1.25x semibold
button          1.00x platform-appropriate
status          1.00x regular/medium
code            ~1.30x semibold/bold
```

Precise point sizes may adapt to native metrics and accessibility scaling; hierarchy is normative.

User data generally carries more visual weight than descriptive labels. For a TokenRow, issuer/service is medium/semibold; account is regular with secondary emphasis. Secondary does not mean tiny or hard to read.

Do not solve hierarchy by making everything semibold. The TOTP code is normally the clearest enlarged text in the main list.

Use tabular digits for codes where supported without requiring a bundled font.

### 25.3 Color roles

Shared semantic roles:

```text
Surfaces
    surface
    surfaceRaised
    surfaceInput
    surfaceSelected

Text
    text
    textSecondary
    textDisabled

Borders
    border
    borderStrong

Interaction
    accent
    focus
    onAccent

Semantic
    info
    warning
    danger
    success

Semantic surfaces
    infoSurface
    warningSurface
    dangerSurface
    successSurface
```

Do not use one generic `disabled` color for all disabled-state properties; disabled text, surface, and border are distinct state treatments.

`ReadOnlyValue`, Disabled, Selected, Warning, and Conflict treatments must remain visually distinct.

### 25.4 Semantic color behavior

Use semantic color first on icons, edges, outlines, small progress indicators, and concise notices. Do not normally flood large application surfaces with saturated semantic color.

Examples:

- conflict: neutral surface + amber edge/icon/text;
- near-expiry: amber countdown ring/text accent;
- invalid/error: red field/status accent;
- destructive confirmation: danger action;
- `ReadOnlyValue`: readable text on a neutral surface, not warning red/amber or disabled styling;
- locked: ordinary shell state, not error;
- deleted: ordinary semantic value, not persistent danger styling.

Selection is an interaction state, not a semantic state. Use a subtle accent tint/outline and do not erase warning/conflict meaning.

### 25.5 Component states

**Primary action**: the single strongest visual action in the current view/dialog; normally filled accent with high-contrast `onAccent` text. A view/dialog normally has at most one. Repeated row-local default actions such as `Show Code` / `Copy` remain secondary/neutral.

**Secondary action**: neutral surface/border; appropriate for ordinary visible actions and repeated row-local default actions.

**Quiet action**: minimal neutral treatment; still clearly actionable.

**Destructive action**: danger text/border or stronger danger fill at the final destructive confirmation.

**Editable input**: `surfaceInput`, visible ordinary border, normal text; focus adds clear focus treatment.

**Read-only value** (`ReadOnlyValue`): full-strength readable text on a neutral surface; selectable/copyable where useful; must not look disabled.

**Disabled control**: subdued text/surface/border, no interaction, still understandable.

**ExclusiveChoice**: selected option uses restrained accent fill plus non-color selected semantics; unselected options remain neutral.

**Selection**: subtle accent tint/outline, no bright full-row flood.

**Focus**: clearly visible ~2 logical-unit outline where applicable and stronger than ordinary selection without changing component geometry.

### 25.6 Geometry

Desktop target geometry uses logical, scale-aware values rather than physical pixel assumptions:

```text
normal field/button height       ~36
compact inline button            ~32
minimum compact control          ~30
ordinary border                    1
focus outline                      2
semantic accent edge              ~3
main-window outer padding         16
dialog outer padding              24
form-row gap                      12
section gap                       16–24
button-group gap                   8
```

Controls may adapt slightly to platform font metrics.

Main/list surfaces may be denser than forms. Expanding available width should usually go to content/search/input fields rather than command buttons.

### 25.7 Token list geometry

Ordinary TokenRows use whitespace + thin dividers, not one rounded card per row.

A desktop ordinary row targets roughly 64–72 logical units depending on font scale, with approximately:

- 12 horizontal padding;
- 8–10 vertical padding;
- 4 between the two text lines;
- 16 minimum identity↔status separation where width allows;
- 12 minimum status↔action separation where width allows;
- 1-unit divider between ordinary rows.

Long identity text may clip/elide with accessible full text; ordinary use should not require horizontal list scrolling.

Conflict groups may have a lightweight containing boundary/amber semantic edge because they group multiple semantic versions. A thin neutral horizontal divider should separate the group-level `Conflict` / `Resolve` header from the first Alternative row. Child rows reuse normal TokenRow surfaces/geometry.

### 25.8 Radius and elevation

Use modest platform-adapted corner treatment:

```text
radiusSmall     ~6
radiusNormal    ~8
radiusLarge     ~10–12
```

These are intent ranges, not pixel-identical requirements.

Use radius to support structure, not as a decorative theme. Token rows normally do not become individual rounded cards.

Elevation is sparse and semantic:

```text
elevation0   ordinary application/form/list content
elevation1   genuinely raised local surface
elevation2   modal/overlay foreground
```

Elevation may be communicated by tonal shift/border rather than shadow, especially in dark themes. Do not add shadows around every row or create nested stacks of floating cards. Desktop dialogs may rely on native window-manager elevation rather than drawing another fake inner shadow.

### 25.9 Icon policy

The first iteration is intentionally text-first and icon-light. Icons supplement information rather than dominate the interface.

Use text for short actions such as `Add`, `Edit`, `Copy`, `Resolve`, `Save`, `Cancel`, and `Lock`.

Icons may reinforce:

- Conflict;
- Warning;
- Error/invalid state;
- Search where platform-standard.

Important states/actions must remain understandable without interpreting iconography. Do not use danger/error imagery for neutral actions such as Cancel. Do not use icon-only management toolbars, service favicons, or a bundled icon library merely for decoration.

Inline semantic icons should normally remain close to text cap-height / roughly 16–20 logical units on desktop and remain visually subordinate to user data and task text.

### 25.10 Platform adaptation

Light and dark environments should both be supported.

Platform implementations may choose exact colors, native control rendering, corner treatment, shadow/tonal elevation, and font sizes, but must preserve the semantic hierarchy, contrast targets, spacing rhythm, and component-state distinctions defined here.

---

## 26. Icon implementation notes

The normative icon policy is in §25.9.

Platform implementations may use native icons for standard affordances when doing so improves recognition without replacing text or semantic state. Platform-native icon treatment does not justify introducing an icon pack or icon-only management language into the shared v0 design.

---

## 27. Platform relationship and profiles

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

Android may use Material 3 internally. Material 3 is not the normative Totipo design system.

Likewise, Swing's current look-and-feel is not normative.

### 27.1 Desktop visual profile

Desktop realizes the shared design as a compact utility application with deliberate breathing room.

Key profile rules:

- persistent application shell;
- neutral-first balanced color;
- system font, modest hierarchy;
- main-window outer padding around 16 logical units;
- dialog/form outer padding around 24;
- ordinary controls around 34–36 high where platform metrics permit;
- compact inline actions around 30–32;
- two-column forms with intrinsic labels and expanding aligned controls;
- normal populated collection header uses Search + secondary count + trailing Add, with deliberate vertical breathing room above and below;
- token list separated primarily by whitespace/thin dividers;
- one outer vertical scroller for long forms;
- restrained selection plus clear keyboard focus;
- desktop menus use modest additional padding when toolkit defaults are cramped, while preserving mnemonics, accelerators, keyboard navigation, and normal menu semantics;
- native dialog/window behavior and normal resizable-window semantics;
- no card-heavy or icon-heavy visual language.

If the active Swing Look & Feel's default UI font is materially smaller than a comfortable modern utility baseline, Totipo may scale the application UI base modestly (roughly +1 to +2 logical points) while preserving the typography roles in §25.

Default window/dialog sizes must make ordinary tasks usable immediately; `pack()` or native preferred sizing must not reduce a task to an unusably cramped initial state.

### 27.2 Android profile

Android follows platform touch-target, navigation, biometric, camera, app-bar, and accessibility conventions while preserving shared Totipo semantics.

Material components may implement Totipo roles, but semantic colors/actions should still follow the neutral-first balanced model rather than turning ordinary content into highly saturated surfaces.

Touch sizing may be roomier than desktop without changing the shared information hierarchy.

---

## 28. Component catalogue

Components are semantic design contracts, not requirements for one-to-one implementation classes.

```text
Action
    PrimaryAction
    SecondaryAction
    QuietAction
    DestructiveAction

Input
    TextField
    SecretField
    SearchField
    NumberField
    ReadOnlyValue
    ExclusiveChoice
    RadioChoice

Structure
    Section
    FieldGroup
    ScrollableForm
    DialogActions
    Divider
    EmptyState
    ConfirmationDialog

Status / Feedback
    InlineNotice
    FieldMessage
    StatusPanel

Application
    ApplicationShell
    VaultIdentity
    NoVaultView
    LockedVaultView
    UnlockedVaultView
    BlockingVaultState
    VaultSelector

TOTP
    TokenList
    TokenRow
    Countdown
    AddTOTPFlow
    TOTPReview
    EditTOTP
    ChangeAuthenticatorSetup
    DeleteTOTPConfirmation

Conflict
    ConflictGroup
    TokenAlternativeRow
    ConflictResolver
    VersionChoice
    ConflictFieldGroup
    AuthenticatorSetupChoice

Vault
    ChangeVaultPassword
    AboutThisVault

Diagnostics
    TokenDiagnostics
```

Components should be added only when real screens demonstrate a recurring need.

The shared catalogue deliberately does not include generic cards, toast notifications, navigation rails, icon-only action buttons, or a general Settings item because v0 does not currently need them.

Platform-native menus, directory choosers, app bars, and similar toolkit-specific structures are platform implementation details rather than shared Totipo components.

---

## 29. V0 non-goals

The v0 design deliberately does **not** require:

- pixel-identical desktop and Android UI;
- a separate Totipo design-system implementation library;
- a top-level Settings screen;
- user-configurable reveal, clipboard, desktop lock, or Android lock timing;
- desktop camera-based QR scanning;
- automatic vault-wide duplicate-secret scanning;
- generic destructive vault repair;
- an automatic global Read-only vault state inferred from failed writes;
- an explicit `Open Read-Only` session mode;
- a general-purpose file manager inside vault selection;
- a history/deleted-items browser;
- bottom navigation or multiple top-level Android destinations;
- selection as a semantic/reveal state (desktop may use subdued selection for keyboard/navigation/contextual actions);
- a global reveal shortcut whose target is ambiguous;
- synchronized biometric/device-unlock configuration;
- password-change behavior that rotates the vault encryption key or re-encrypts vault contents;
- card-heavy, icon-heavy, or highly saturated visual treatment.

These may be revisited only when concrete product requirements justify them.

---

## 30. V0 implementation review checklist

A desktop or Android implementation can be reviewed against the following baseline.

### 30.1 Application lifecycle

- Startup distinguishes no-vault, locked, unlocked, and blocking vault states in a persistent application shell/screen model.
- Desktop `NO_VAULT` uses the centered `Choose a vault to continue` state with Select Vault as the visually primary path and Create New Vault as secondary.
- The currently selected recognizable Totipo vault location/access reference is remembered locally; successful password authentication is not required for that selection to remain current.
- A location established not to be a Totipo vault does not replace the remembered vault.
- Passwords are not remembered as application preferences.
- Select Vault and Create New Vault are separate intents.
- Desktop existing-vault selection is directory-only, hides dot/hidden directories by default, omits generic file-management controls, and supports selecting a highlighted child directory without first navigating into it.
- Cancelling Change Vault before accepting a replacement returns to the previous vault locked when available; it never resurrects the old unlocked session.
- At most one unlocked vault session exists at a time.
- Locking follows the fixed platform policy in §16.
- Desktop exposes explicit Lock and `Ctrl+L`; explicit Lock retires owned sensitive child flows rather than being blocked by them.

### 30.2 Main product workflow

- The normal path is find TOTP → Show Code → Copy.
- The collection action is semantically Add TOTP but may be labelled simply `Add` in an unambiguous context.
- On desktop, the normal populated collection header is Search + count + trailing Add; Search expands and Add remains the sole visually primary collection action.
- Desktop Refresh stays in the Vault menu rather than consuming permanent main-view space.
- Active filtering reports a concise `M of N` count.
- A genuinely empty vault and a zero-result search use distinct contextual empty states in the list area.
- Edit/Delete/Refresh/Change Vault remain visually secondary.
- Search operates only on issuer/account, using case-insensitive whitespace-split AND substring matching.
- User-facing ordering does not leak protocol/storage order.
- Desktop search/list keyboard navigation follows §5/§7 without deriving codes from focus/selection.

### 30.3 TokenRow

- Codes are concealed by default.
- `Show Code` is the default hidden-row action but normally uses compact neutral/secondary styling rather than a filled primary treatment.
- Reveal lifetime follows current-period / strict-`<10s` one-period grace.
- Any pre-staged grace code is never exposed or copyable before validity begins.
- Rows reveal independently.
- `Copy` copies the currently displayed canonical digits.
- Copy does not extend disclosure.
- Copy feedback is local/transient rather than a global toast.
- Expiry/countdown is not communicated by graphics or color alone.
- Near expiry uses warning/amber, not danger/red.

### 30.4 Clipboard

- Copied values are cleared on a best-effort basis at the earlier of that exact code's expiry or 30 seconds after Copy.
- Cleanup never destroys newer clipboard contents.
- A rollover never silently replaces the clipboard with the next TOTP.
- Platform sensitive-clipboard facilities are used where appropriate.

### 30.5 Add/Edit/Delete

- Form commit actions remain available when input is incomplete/invalid; activation exposes validation errors and focuses/scrolls to the first actionable problem.
- QR/URI acquisition produces a reviewable draft before commit.
- Clipboard setup import is explicit, not background clipboard sniffing.
- Raw secret material is not unnecessarily redisplayed.
- Same issuer/account offers an explicit Update Existing versus Add Another decision without silent replacement.
- Ordinary Edit focuses on issuer/account plus a setup summary.
- Credential replacement uses explicit `Change setup…`.
- Delete is a separate destructive action and discloses that vault history is not securely erased.

### 30.6 Conflicts

- Alternatives/versions are user choices; Heads are provenance.
- Multiple Heads that produce one Alternative do not appear as duplicate choices.
- Conflict remains usable and resolution is optional.
- Each complete eligible Alternative may independently Show Code / Copy / Edit.
- Head count or client time is not treated as a vote/freshness winner.
- Simple whole-version resolution is offered first with no implicit winner; `Resolve` validates on activation rather than being disabled before a selection is made.
- `Combine details…` is the secondary field-level path.
- Detailed resolution provides explicit Back / Cancel / Save Resolution navigation; returning Back does not silently turn a composed draft into a whole-version choice.
- Conflict Alternatives use the ordinary displayed-identity sort rule within their group; group position follows the first Alternative under the same ordering and does not imply preference.
- Conflicting fields use distinct existing-value choices plus a custom editable row where applicable.
- Secret + algorithm + digits + period are one atomic Authenticator Setup resolution unit.
- Existing secret bytes are never displayed.
- Newly arrived conflict information and publication uncertainty are represented truthfully.

### 30.7 Vault states

- v0 does not infer a global Read-only vault state from write/publication failure, publication uncertainty, invalid data, or unavailable data.
- Definite publication failure and publication uncertainty remain operation-specific.
- Wrong location, unlock failure, unsupported/invalid data, unavailable data, and publication uncertainty remain distinct where the API can distinguish them.
- Blocking vault state does not show a normal token list.
- Localized failures remain localized where safe.
- Generic destructive `Repair Vault` behavior is not invented.

### 30.8 Passwords and local unlock

- Empty-password creation/opening requires explicit confirmation.
- Change Vault Password requires the current vault password plus the new password/confirmation; an unlocked session alone does not authorize changing the future unlock credential.
- Change Vault Password keeps the normal UI consequence-focused and rewraps the unchanged vault key/root without rotating it or re-encrypting vault contents.
- Empty new password requires explicit confirmation.
- Historical retained copies are not promised to be revoked by a password change.
- Android biometric/device unlock remains local and retains password fallback; in v0.8 it does not substitute for entry of the current vault password during password change.
- Local biometric/device-unlock material is invalidated/re-established appropriately after password-wrapper changes.

### 30.9 Visual system and accessibility

- Ordinary UI uses neutral surfaces and purposeful semantic accents.
- Desktop collection headers and menus use deliberate spacing rather than relying on cramped toolkit defaults.
- Conflict presentation visually separates the group-level header from Alternative rows with a thin neutral divider.
- A view/dialog normally has at most one visually primary action; repeated TokenRow actions do not create a field of filled primary buttons.
- Normal text/filled-action contrast targets §24.
- Secondary text remains comfortably readable.
- Semantic color does not flood large surfaces by default.
- Selection is subdued and does not replace conflict/warning meaning.
- Editable inputs, read-only values, and disabled controls are visually distinct.
- Desktop forms use aligned expanding controls and consistent spacing rhythm.
- Ordinary token rows use whitespace/thin dividers rather than card-heavy layout.
- Focus remains clearly visible and does not move component geometry.
- Important meaning is not conveyed by color/icon alone.
- Desktop is fully usable by keyboard for ordinary workflows.
- Android uses appropriate touch targets and system Back behavior.
- Sensitive Android task/app-switcher previews are protected.
- Platform-native interaction conventions may differ without changing shared Totipo semantics.

---

## 31. Future design work

The v0 product model, interaction semantics, and visual foundations are sufficiently specified for the next implementation pass. Remaining work should primarily be driven by implementation testing, accessibility review, and concrete product requirements rather than speculative feature expansion.

Immediate implementation work should now center on:

- reshape ordinary Edit around issuer/account plus setup summary, with separate `Change setup…` and `Delete TOTP…` flows;
- complete Add TOTP URI/manual acquisition, review, and explicit same-issuer/account Update Existing / Add Another behavior;
- restore the simple whole-version conflict resolver first, with `Combine details…` for field-level resolution and atomic Authenticator Setup choices;
- align blocking/unavailable vault-state presentation with truthful API distinctions and discoverable genuinely unavailable actions;
- add `About This Vault…` as the secondary vault-identity/details surface;
- complete the concise Change Vault Password experience with current-password reauthentication;
- perform final manual light/dark-theme contrast and density review;
- perform final keyboard/accessibility review under real desktop environments.

The persistent desktop shell, core lock/inactivity behavior, main collection composition, local copy feedback, search behavior, and initial desktop visual foundation are now sufficiently specified to serve as the baseline for that work.

Likely later areas include:

- detailed accessibility announcements for countdown rollover, lock transitions, and conflict updates;
- refinement of exact platform color values after light/dark implementation testing;
- history/deleted-item inspection, if a real user need emerges;
- explicit user-selected read-only sessions or backend-reported read-only capability, if a concrete requirement justifies a separate product design;
- device admission/management UX, if exposed directly to end users;
- platform-specific secure-storage implementation details for local biometric/device unlock;
- explicit recovery workflows supported by future protocol/Java API capabilities;
- future local preferences only when a concrete need justifies adding a Settings surface.

---

