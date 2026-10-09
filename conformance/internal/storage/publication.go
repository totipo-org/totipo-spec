package storage

import "bytes"

// Install classifies backend outcomes, not host syscalls. On failure the
// returned observation does not assert absence or describe effects of an
// ambiguous backend operation. Durable means required persistence succeeded.
func Install(existing, intended []byte, kind string, durable bool) ([]byte, string) {
	if len(intended) != 1024 {
		return existing, "FAILED"
	}
	if kind != "absent" {
		if kind != "regular" || !bytes.Equal(existing, intended) {
			return existing, "FAILED"
		}
		return existing, "ALREADY_PRESENT_EXACT"
	}
	if !durable {
		return existing, "FAILED"
	}
	return bytes.Clone(intended), "PUBLISHED_NEW"
}

// Store models canonical VAULT separately from immutable object publication.
// It does not model filesystem races, persistence primitives, or remote completeness.
type Store struct{ Vault []byte }

func (s *Store) Publish(existing, intended []byte, kind string, durable bool) ([]byte, string) {
	return Install(existing, intended, kind, durable)
}

// Create is absent -> initial-create only. Revalidated denotes a successful
// exact canonical reread and root-wrap validation after durable publication.
// Failure leaves effects unknown; the returned bytes do not assert absence.
func Create(existing, intended []byte, kind string, complete, durable, revalidated, orphans bool) ([]byte, string) {
	if kind != "absent" || orphans || !complete || !durable || !revalidated || len(intended) != 87 {
		return existing, "FAILED"
	}
	return bytes.Clone(intended), "CREATED"
}
