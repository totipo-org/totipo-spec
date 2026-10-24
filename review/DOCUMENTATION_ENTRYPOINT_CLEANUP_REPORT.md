# Documentation entry point cleanup

## Scope and starting evidence

Date: 2026-10-24. Repository: `/home/niki/Sources/totipo-spec` only.
Starting branch: `main`; starting HEAD:
`5c5f28ac0f85fcd8cb64d20ddbd8a2da9381e7ac`.
Initial `git status --short` was empty; no unrelated changes were present.
The index contained no staged content. Protocol remains **v1/r19**, a design
draft with a **moving-pre-rc** profile and **92 required cases**.

Classification: **DOCS_ONLY**. No specification text, qualification logic,
requirements, corpus, implementation, build, or dependency change.

Before editing, read AGENTS, README, CONTRIBUTING, the active specification,
requirements/vector/conformance READMEs, and the vector case contract. There is
no generator README or index; its existing documentation is in CONTRIBUTING,
the case contract, and the conformance README. Inspected the generator file
inventory, flake source filter, workflow, Go spec-policy and architecture
consumers, and historical qualification reports. Existing documentation provided
enough navigation; no new implementer guide was needed.

Inherited evidence is historical: the tooling simplification report records
host qualification and an owner-reported Nix PASS dated 2026-10-14, with unknown
cache/fresh provenance and remote CI unverified. Later commits include flake
changes. No fresh qualification or exact-starting-HEAD Nix/remote result was
established here; earlier evidence is not promoted to a current pass.

## README and navigation

README retains its overview and repository map while reducing 112 lines to 48.
It now explains the TOTP vault's purpose and concurrent-change interpretation,
keeps the normative v1/r19 link prominent, and explicitly distinguishes the moving
design from a frozen release or RC.

The trust summary covers incomplete/stale views, external retention and rollback
guarantees, unlocked-client compromise, offline password guessing, and deletion
versus erasure, with direct normative section links. Detailed envelope sizes,
field semantics, graph rules, and schema implementation notes were removed from
README in favor of existing normative and technical documentation.

The linked map distinguishes requirements pins, committed vectors, independent
fixture generation, and read-only verification/semantic execution. The implementer
reading path starts with goals/trust boundaries, identifies Section 20 scopes,
then points to applicable normative text, vectors, the case contract, and consumer
usage. It preserves the limits of portable conformance evidence.

Recent review-report enumeration was replaced by a single archive entry explaining
historical decisions and qualification evidence. Contribution, qualification,
maintenance, agent guidance, and security reporting now have direct destinations.

## Consolidation and operational preservation

CONTRIBUTING is the authoritative qualification-policy home. The previous AGENTS
policy was moved there verbatim apart from heading depth, with the explicit human
command example and report/Git-aware membership clarification retained from the
previous contributor guidance. A normalized text comparison passed for the entire
moved policy. This preserves classifications, baseline/focused/final boundaries,
all-package race coverage, every fuzz command with `10s` and `-parallel 2`, failure
handling, independent generation/verification, invariance rules, invalidation,
execution accounting, and the human-owned Nix boundary.

AGENTS now requires reading and following that authoritative policy and retains
the main agent guardrails. Duplicate README qualification recipes were removed.
Contributor platform/environment details, historical audit retrieval, temporary
generation comparison, module independence, schema and wire-oracle limitations,
and standalone-module commands remain. Added headings make the existing historical
audit and fixture-maintenance material addressable; no procedure was weakened.

## Exact protected-scope verification

For each starting tracked path except README.md, CONTRIBUTING.md, and AGENTS.md,
compared current file bytes directly with `git show <starting-HEAD>:<path>` using
`cmp`: **284 files byte-identical**. This includes every protected file, all prior
review reports, spec, vectors, requirements, generator, conformance, design assets,
flake/lock, environment, CI, go.work, and module metadata. Existing historical
reports were neither rewritten nor deleted.

An independent Git-visible inventory comparison against `git ls-tree` at the
starting HEAD permits exactly one addition: this report. There are no other path
additions or deletions. Root go.mod, go.sum, and go.work.sum were absent and remain
absent; generator/go.mod and conformance/go.mod/go.sum remain byte-identical.
No current tools directory exists. No wire definitions, revision identifiers,
pins, hashes, integrity metadata, fixtures, expected outcomes, or release artifacts
were changed, and no generation occurred.

The inspected flake includes spec, vectors, requirements, repository Go/native
sources and module metadata, go.work, and archived pre-r16 history. All four
milestone Markdown files are excluded. The spec checker reads only the active
spec and archived history; architecture checks inspect source/module metadata,
not these documents' content. Artifact verification reads committed contracts,
manifest, cases, and profile. Thus no qualification input changed and no new human
Nix request is required. This preserves the starting evidence disposition without
asserting a new pass.

## Validation and execution accounting

Passed all 26 local link existence and heading-fragment checks across the changed/new
documents; balanced fenced blocks, heading spacing, trailing whitespace, and final
newlines; manual Markdown/table and editorial review; spec/revision references
against the active spec/profile; policy-text preservation; complete protected-byte
and Git-visible inventory comparisons; `git diff --check`; and
`git diff --cached --quiet`.

No dedicated Markdown linter was available; used a temporary standard-library Go
static audit outside Git membership. No external link/network validation was
needed for the changed local navigation.

| Execution | Count / disposition |
| --- | --- |
| Baseline Go tests / full conformance | 0 / 0; documentation-only path |
| Focused Go package tests / Python tests | 0 / 0 |
| Final complete host suite / incomplete host attempts | 0 / 0 |
| Vet / race | 0 / 0; protected code unchanged |
| Dispatch / Open / graph fuzz | 0 / 0 / 0; no duration or result claimed |
| Temporary Go documentation/inventory audit | 1 PASS; static checks only |
| Generator / agent Nix / human Nix in this milestone | 0 / 0 / 0 |
| Remote CI dispatch / fresh observation | 0 / 0 |

An initial Python-based editing command exited 127 because python3 was unavailable;
it made no edits and ran no tests. Continued with shell editing and the temporary
Go audit. Exploratory searches also encountered absent tools/architecture paths;
the actual repopolicy consumer was subsequently inspected. No qualification
failure, bounded-fuzz reproduction, or incomplete qualification run occurred.

## Remaining issues and final Git state

The non-normative v0.9 product-design baseline still contains password-change
guidance predating r19's immutable VAULT model, as already recorded in the r19
report. Updating that historical product-design baseline or downstream behavior
requires separate authorization. No normative inconsistency requiring a protocol
edit was identified. Independent interoperability, security/platform evidence,
and any eventual RC freeze remain separate work; this milestone claims none.

Branch and HEAD remain unchanged. Three tracked documents are modified; this new
report has intent-to-add membership only. No content is staged, all changes remain
uncommitted, and no Nix command, commit, tag, release, push, or remote mutation was
performed. Temporary audit code is outside repository Git membership.
