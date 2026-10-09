# v1/r19 cross-vault transfer notes

**Application/workflow design.**
**Not part of Totipo v1 conformance.**
**Not normative protocol state.**

Cross-vault copy can be useful independently of migration. Migration can be
modeled as new-vault creation followed by copy, with independent source and
destination stores and a fresh destination K_root. The source stays unchanged;
its old vault may remain useful for later recovery or another copy.

Ordinary copy should normally generate fresh destination TOKEN_IDs. Applications
may offer “Copy current state” or “Preserve available history” and should describe
what validated state/history each mode copies. History, conflict, and tombstone
policy belong to the application. Copied assertions are re-authored and encrypted
under destination keys; encrypted source objects are not destination objects.

A resumable transfer may retain application-local receipts. Source and destination
VAULT_IDs can bind that local job, and a stable per-job mapping from selected source
tokens to fresh destination TOKEN_IDs can prevent duplicate destination identities
on partial/retried copy. Receipts should not contain passwords, roots, or token
secrets unnecessarily. Their protection, cleanup, retry behavior, and lifetime are
application design, not a new protocol requirement.

The observed source snapshot may change while a transfer is paused. A resumed job
should describe the snapshot/selection it actually copies and handle newly observed
source data deliberately. Neither migration nor ordinary copy can prove unseen
remote data is absent or that devices/providers have fully converged. Applications
should disclose possible omission of not-yet-synchronized data and retain the old
vault until the user verifies the new one. Copying an exposed TOTP seed does not
revoke the attacker's knowledge; issuer-side credential rotation remains necessary.

No transfer protocol, schema, vectors, permanent migration ancestry, or source/
destination binding is added to VAULT or TOKEN. Local job receipts are not
synchronized protocol-authoritative state, and a destination need not remember its
source permanently.
