# Totipo Design Guidelines

**Status:** Draft v0  
**Scope:** Cross-platform product and interaction design  
**Normative protocol status:** Non-normative

These guidelines define Totipo's shared product language across desktop and Android, and are intended to remain applicable to future platforms such as iOS.

They do **not** define the Totipo storage or synchronization protocol and do not override protocol requirements. Platform applications should follow native platform conventions where those conventions do not conflict with the product semantics defined here.

> **Consistency of meaning and behavior is more important than consistency of appearance.**

---

## 1. Product model

Totipo is primarily a TOTP retrieval tool backed by a synchronized vault.

The normal user journey is:

1. Open a vault.
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

4. **Exceptional state handling**
   - Resolve conflicts.
   - Handle invalid/corrupt data.
   - Handle read-only or otherwise constrained vault state.

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

## 8. Shared semantic states

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

A reusable `StatusPanel` should represent substantial exceptional state.

A reusable `InlineNotice` should represent contextual information inside a form.

Neither should rely on color or iconography alone.

---

## 9. Choice controls

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

## 10. Forms and sections

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

## 11. Dialog actions

Task dialogs should use a consistent action area.

Typical desktop arrangement:

```text
                              Cancel    Primary Action
```

Long forms may scroll, but the action area should remain readily accessible.

Escape/system Back should normally correspond to the neutral cancellation path unless doing so would discard meaningful work without an appropriate confirmation.

Enter may invoke the primary action only where safe and unambiguous.

---

## 12. Accessibility

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

## 13. Visual tokens

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

## 14. Icons

Icons reinforce text rather than replace important concepts.

Important states and actions must remain understandable without interpreting iconography.

Do not use danger/error imagery for neutral actions such as Cancel.

---

## 15. Platform relationship

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

## 16. Initial component catalogue

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
    TokenCandidate
    ConflictPanel
    AddTOTPFlow
```

Components should be added only when real screens demonstrate a recurring need.

---

## 17. Not yet specified

The following areas remain to be worked through before the guidelines should be considered complete:

- Edit TOTP workflow
- Delete TOTP workflow and history disclosure
- Conflict resolution UI
- Vault selection
- Open/unlock vault
- Change vault password
- Read-only and invalid/corrupt vault presentation
- application-level lock/background behavior
- clipboard lifetime and feedback details beyond TokenRow Copy
- detailed Android navigation structure
- detailed desktop menu structure
- visual token values after implementation testing

These should be designed from concrete application behavior rather than invented in isolation.
