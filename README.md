# Totipo

Totipo is an encrypted, append-only vault format for time-based one-time password
(TOTP) credentials. It defines how clients store and interpret credential changes
in a configured durable store, preserving concurrent disagreements as conflicts.
The same format supports local-only and externally synchronized stores;
synchronization is optional and outside the protocol.

The authoritative current specification is **[Totipo Vault Format v1,
revision r19](spec/totipo-vault-format-v1.md)**. It is a **design draft** with a
moving pre-release-candidate requirements profile, not a frozen release or RC.

The trust boundary matters: Totipo authenticates and validates observed content,
but does not prove that a store presents complete or fresh history. Store loss,
withheld objects, or rollback can hide history; stronger freshness, rollback
detection, and retention guarantees require storage or application mechanisms.
An unlocked vault root grants authoring authority, and a compromised unlocked
client can read secrets and author valid state. The bootstrap permits offline
password guessing; deletion does not securely erase credentials from history.
See specification [Section 2](spec/totipo-vault-format-v1.md#2-threat-model-and-storage-boundary),
[Section 5](spec/totipo-vault-format-v1.md#5-password-bytes), and
[Section 16](spec/totipo-vault-format-v1.md#16-complete-values-and-truthful-descriptions)
for the limitations and applicable requirements.

| Entry point | Role |
| --- | --- |
| [Specification](spec/totipo-vault-format-v1.md) | Normative protocol and conformance scopes |
| [Requirements](requirements/README.md) | Moving pre-RC profile pinning the spec, schemas, manifest, and required cases |
| [Vectors](vectors/README.md) and [case contract](vectors/FORMAT.md) | Committed byte, negative, graph, fold, and modeled storage cases |
| [Verifier / Go consumer](conformance/README.md) | Read-only artifact verification and semantic execution; conformance evidence rather than a production client |
| [Generator maintenance](CONTRIBUTING.md#fixture-maintenance-and-implementation-independence) | Deliberate fixture writing in the independent [generator module](generator/); never part of routine verification |
| [Review archive](review/) | Historical decisions, reviews, and qualification evidence; commands and results belong to their recorded revisions |

A new implementer can start with the specification's goals and trust model
(Sections 1–2), then identify the intended conformance scopes and operations in
[Section 20](spec/totipo-vault-format-v1.md#20-conformance-scopes-and-evidence).
Read the applicable normative sections alongside the vector overview and case
contract, and use the requirements profile to identify the exact current inputs.
The Go consumer documents how to verify those artifacts and execute every manifest
case. Passing the portable corpus does not establish application conformance,
independent interoperability, or a live backend's crash durability.

[CONTRIBUTING.md](CONTRIBUTING.md) documents contribution review, the supported
Nix/Go environment, qualification, and deliberate fixture maintenance.
Its [qualification policy](CONTRIBUTING.md#qualification-policy) defines the
human-owned `nix flake check` gate and the documentation-only validation path.
[AGENTS.md](AGENTS.md) directs agents to that same policy.
Report suspected vulnerabilities through [SECURITY.md](SECURITY.md).
