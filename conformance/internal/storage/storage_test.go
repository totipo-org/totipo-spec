package storage

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"totipo/conformance/internal/cryptov1"
	"totipo/conformance/internal/object"
)

func TestNeverReadIgnoredEntries(t *testing.T) {
	k, _ := cryptov1.Derive(make([]byte, 32))
	id := strings.Repeat("a", 64)
	entries := []Entry{}
	for _, path := range []string{"objects-v2/" + id, "objects-v999/" + id, "objects-v1/nested/" + id, "objects-v1/" + strings.ToUpper(id), "objects-v1/" + id + ".tmp", "objects-v1/../objects-v1/" + id, "/objects-v1/" + id, "objects-v1\\" + id} {
		entries = append(entries, Entry{Path: path, Kind: "regular", Read: func() ([]byte, error) { t.Fatal("ignored path read"); return nil, nil }})
	}
	for _, kind := range []string{"directory", "symlink", "fifo", "socket", "device"} {
		entries = append(entries, Entry{Path: "objects-v1/" + id, Kind: kind, Read: func() ([]byte, error) { t.Fatal("nonregular entry read"); return nil, nil }})
	}
	out, e := Scan("directory", entries, k)
	if e != nil || len(out) != 0 {
		t.Fatal(out, e)
	}
	if _, e = Scan("symlink", entries, k); e == nil {
		t.Fatal("namespace symlink accepted")
	}
}
func TestStorageAuthenticationBoundary(t *testing.T) {
	k, _ := cryptov1.Derive(make([]byte, 32))
	id, b, _ := k.Seal([]byte{0, 1, 0, 1, 2, 0, 2, 0, 1, 99})
	for _, n := range []int{0, 1023, 1025} {
		out, e := Scan("directory", []Entry{{Path: "objects-v1/" + id, Kind: "regular", Read: func() ([]byte, error) { return make([]byte, n), nil }}}, k)
		if e != nil || out[0].Class != InvalidStorage || out[0].Object != nil {
			t.Fatal("wrong-size became evidence", out, e)
		}
	}
	valid, err := Scan("directory", []Entry{{Path: "objects-v1/" + id, Kind: "regular", Read: func() ([]byte, error) { return b, nil }}}, k)
	if err != nil || valid[0].Class != "INVALID" {
		t.Fatal("invalid authenticated grammar accepted", err)
	}
	b[0] ^= 1
	out, e := Scan("directory", []Entry{{Path: "objects-v1/" + id, Kind: "regular", Read: func() ([]byte, error) { return b, nil }}}, k)
	if e != nil || out[0].Class != InvalidStorage {
		t.Fatal("bad AEAD became evidence")
	}
	_, e = Scan("directory", []Entry{{Path: "objects-v1/" + id, Kind: "regular", Read: func() ([]byte, error) { return nil, errors.New("unreadable") }}}, k)
	if e == nil {
		t.Fatal("missing read diagnostic")
	}
}

func TestIncompleteScanPreservesAcceptedObservations(t *testing.T) {
	k, _ := cryptov1.Derive(make([]byte, 32))
	id, b, _ := k.Seal([]byte{0, 1, 0, 1, 2, 0, 2, 0, 1, 99})
	entries := []Entry{
		{Path: "objects-v1/" + id, Kind: "regular", Read: func() ([]byte, error) { return b, nil }},
		{Path: "objects-v1/" + strings.Repeat("b", 64), Kind: "regular"},
	}
	obs, err := Scan("directory", entries, k)
	if err == nil || len(obs) != 1 || obs[0].Class != "INVALID" {
		t.Fatalf("%v %v", obs, err)
	}
}

func TestPublicationRetries(t *testing.T) {
	intended := make([]byte, 1024)
	intended[0] = 1
	// An error can still have installed the bytes. The next ordinary observation
	// sees them, and an exact retry succeeds without fresh persistence.
	_, out := Install(nil, intended, "absent", false)
	if out != "FAILED" {
		t.Fatal(out)
	}
	existing := append([]byte{}, intended...)
	after, out := Install(existing, intended, "regular", false)
	if out != "ALREADY_PRESENT_EXACT" || &after[0] != &existing[0] {
		t.Fatal("exact retry mutated")
	}
	other := append([]byte{}, intended...)
	other[0] = 2
	after, out = Install(other, intended, "regular", true)
	if out != "FAILED" || after[0] != 2 {
		t.Fatal("overwrote different")
	}

}

func TestCanonicalEntryTypes(t *testing.T) {
	k, _ := cryptov1.Derive(make([]byte, 32))
	never := func() ([]byte, error) { t.Fatal("deliberately read wrong-type entry"); return nil, nil }
	for _, kind := range []string{"absent", "symlink", "directory", "fifo", "socket", "device", "other"} {
		if _, e := OpenBootstrap(Entry{Path: "vault", Kind: kind, Read: never}, nil); e == nil {
			t.Fatal("bootstrap type", kind)
		}

	}
	for _, path := range []string{"VAULT", "Vault", "vault.tmp", "vault/conflict", "other/vault"} {
		if _, e := OpenBootstrap(Entry{Path: path, Kind: "regular", Read: never}, nil); e == nil {
			t.Fatal("bootstrap alias", path)
		}
	}
	root := make([]byte, 32)
	record, e := cryptov1.Wrap(nil, root, make([]byte, 16), make([]byte, 12))
	if e != nil {
		t.Fatal(e)
	}
	got, e := OpenBootstrap(Entry{Path: "vault", Kind: "regular", Read: func() ([]byte, error) { return record, nil }}, nil)
	if e != nil || len(got) != 32 {
		t.Fatal("regular bootstrap rejected", e)
	}
	entries := []Entry{{Path: "objects-v1/" + strings.Repeat("a", 64), Kind: "regular", Read: never}}
	for _, kind := range []string{"symlink", "regular", "fifo", "socket", "device", "other"} {
		if obs, e := Scan(kind, entries, k); e == nil || len(obs) != 0 {
			t.Fatal("namespace traversed", kind)
		}
	}
	if obs, e := Scan("missing", entries, k); e != nil || len(obs) != 0 {
		t.Fatal("missing namespace", e)
	}
}

// Exercise real TOKEN authorship and retries against a separately held canonical
// VAULT, including failures and successful observation/unlock afterward.
func TestTokenAuthorshipPreservesCanonicalVault(t *testing.T) {
	root := bytes.Repeat([]byte{0x37}, 32)
	password := []byte("creation password")
	record, err := cryptov1.Wrap(password, root, bytes.Repeat([]byte{0x23}, 16), bytes.Repeat([]byte{0x45}, 12))
	if err != nil {
		t.Fatal(err)
	}
	before := bytes.Clone(record)
	canonical, result := Create(nil, record, "absent", true, true, true, false)
	if result != "CREATED" {
		t.Fatal(result)
	}
	store := Store{Vault: canonical}
	keys, err := cryptov1.Derive(root)
	if err != nil {
		t.Fatal(err)
	}
	token := object.Object{Identity: bytes.Repeat([]byte{0x56}, 32), Parents: [][]byte{}, Status: 1, Algorithm: 1, Digits: 6, Period: 30, Secret: []byte("12345678901234567890")}
	plain, err := token.Encode()
	if err != nil {
		t.Fatal(err)
	}
	id, encrypted, err := keys.Seal(plain)
	if err != nil {
		t.Fatal(err)
	}
	published, result := store.Publish(nil, encrypted, "absent", true)
	if result != "PUBLISHED_NEW" {
		t.Fatal(result)
	}
	if _, err := keys.Open(id, published); err != nil {
		t.Fatal(err)
	}
	for _, trial := range []struct {
		kind     string
		existing []byte
		durable  bool
		result   string
	}{
		{"regular", published, false, "ALREADY_PRESENT_EXACT"},
		{"regular", []byte{0}, true, "FAILED"},
		{"absent", nil, false, "FAILED"},
	} {
		if _, got := store.Publish(trial.existing, encrypted, trial.kind, trial.durable); got != trial.result {
			t.Fatal(got)
		}
		if !bytes.Equal(before, store.Vault) {
			t.Fatal("VAULT changed during ordinary authorship/publication")
		}
	}
	unlocked, err := cryptov1.Unwrap(password, store.Vault)
	if err != nil || !bytes.Equal(unlocked, root) {
		t.Fatal("root changed", err)
	}
}
