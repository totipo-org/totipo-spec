# Agent workflow

[CONTRIBUTING.md — Qualification policy](CONTRIBUTING.md#qualification-policy)
is the authoritative home for baseline evidence, classifications, focused
iteration, final host checks, documentation invariance, source membership,
and execution accounting. Read and follow that policy before starting a milestone.

Start on clean `main` unless the user authorizes continuing identified changes.
Record the starting HEAD, Git state, protocol revision, and inherited evidence.
Work only in the authorized repository and scope; preserve historical reports.

Agents do not execute Nix. Use ordinary Go commands in the activated Nix
environment. Follow the contributor policy's documentation-only validation path
when applicable; preserve its complete final suite for semantic, conformance,
or qualification-logic changes. Generation requires deliberate authorization
and remains separate from verification.

Keep tracked edits unstaged. New intended files may use intent-to-add only;
`git diff --cached --quiet` must pass. Request any required human Nix gate only
after intended files are Git-known and included inputs are frozen. Leave commits,
pushes, and release operations to the explicitly authorized user workflow.
