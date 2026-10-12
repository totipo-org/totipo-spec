package wire_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"totipo/generator/internal/fixture"
	"totipo/generator/internal/wire"
)

func raw(t *testing.T, s string) []byte {
	t.Helper()
	b, e := hex.DecodeString(s)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func corpus(t *testing.T) []fixture.Case {
	t.Helper()
	paths, e := filepath.Glob("../../../vectors/cases/*/*.json")
	if e != nil {
		t.Fatal(e)
	}
	out := []fixture.Case{}
	for _, p := range paths {
		b, e := os.ReadFile(p)
		if e != nil {
			t.Fatal(e)
		}
		var c fixture.Case
		if e = json.Unmarshal(b, &c); e != nil {
			t.Fatal(e)
		}
		out = append(out, c)
	}
	if len(out) != 92 {
		t.Fatalf("case inventory: %d", len(out))
	}
	return out
}
func equal(t *testing.T, label string, got []byte, want string) {
	t.Helper()
	if hex.EncodeToString(got) != want {
		t.Fatalf("%s differs from committed known field", label)
	}
}
func TestCommittedComposition(t *testing.T) {
	for _, c := range corpus(t) {
		t.Run(c.ID, func(t *testing.T) {
			if c.Crypto != nil {
				p := raw(t, c.Semantic)
				k, e := wire.Derive(raw(t, c.Root))
				if e != nil {
					t.Fatal(e)
				}
				id, b, e := k.Seal(p)
				if e != nil {
					t.Fatal(e)
				}
				x := c.Crypto
				if id != x.ObjectID || len(p) != x.SemanticLength || len(b) != 1024 {
					t.Fatal("id/length/size")
				}
				rid := raw(t, id)
				pad, e := wire.Padded(p)
				if e != nil {
					t.Fatal(e)
				}
				for label, pair := range map[string]struct {
					got  []byte
					want string
				}{
					"id key": {k.ID, x.IDKey}, "root key": {k.ObjectRoot, x.ObjectRootKey}, "object key": {k.ObjectKey(rid), x.ObjectKey},
					"nonce": {rid[:12], x.Nonce}, "aad": {wire.AAD(rid), x.AAD}, "padding": {pad, x.Padded}, "envelope": {b, x.Object}, "ciphertext": {b[:1008], x.Ciphertext}, "tag": {b[1008:], x.Tag},
				} {
					equal(t, label, pair.got, pair.want)
				}
				if c.Input != nil {
					got, e := c.Input.Encode()
					if e != nil {
						t.Fatal(e)
					}
					equal(t, "canonical TOKEN", got, c.Semantic)
				}
			}
			if c.PostAEAD != nil {
				k, e := wire.Derive(raw(t, c.Root))
				if e != nil {
					t.Fatal(e)
				}
				id, plain, b, e := wire.Malform(k, raw(t, c.Semantic), c.PostAEAD.Defect)
				if e != nil {
					t.Fatal(e)
				}
				equal(t, "malformed id", id, c.PostAEAD.ObjectID)
				equal(t, "malformed plaintext", plain, c.PostAEAD.Plaintext)
				equal(t, "malformed object", b, c.PostAEAD.Object)
			}
			if c.Bootstrap != nil {
				got, e := wire.VaultID(raw(t, c.Bootstrap.Record))
				if e != nil {
					t.Fatal(e)
				}
				equal(t, "vault hash", got, c.Bootstrap.VaultID)
			}
		})
	}
}
func TestFramingAndAuthoringConstraints(t *testing.T) {
	b, e := wire.Encode(0x1234, []byte{0xab, 0xcd})
	if e != nil || !bytes.Equal(b, []byte{0x12, 0x34, 0, 2, 0xab, 0xcd}) {
		t.Fatal("tag/width/endianness", e)
	}
	if !bytes.Equal(wire.U64(0x0102030405060708), []byte{1, 2, 3, 4, 5, 6, 7, 8}) {
		t.Fatal("u64 width/order")
	}
	if !bytes.Equal(wire.U32(0x01020304), []byte{1, 2, 3, 4}) {
		t.Fatal("u32 width/order")
	}
	if _, e := wire.Encode(1, make([]byte, 65536)); e == nil {
		t.Fatal("TLV overflow")
	}
	for _, c := range corpus(t) {
		if c.Input == nil {
			continue
		}
		for name, change := range map[string]func(*wire.Token){
			"identity": func(o *wire.Token) { o.Identity = nil }, "parents": func(o *wire.Token) { o.Parents = make([][]byte, 5) },
			"duplicate": func(o *wire.Token) { o.Parents = [][]byte{make([]byte, 32), make([]byte, 32)} },
			"status":    func(o *wire.Token) { o.Status = 3 }, "text": func(o *wire.Token) { o.Issuer = "\xff" },
			"text-size": func(o *wire.Token) { o.Account = strings.Repeat("é", 129) }, "secret": func(o *wire.Token) { o.Secret = nil },
			"secret-size": func(o *wire.Token) { o.Secret = make([]byte, 129) }, "algorithm": func(o *wire.Token) { o.Algorithm = 4 },
			"digits": func(o *wire.Token) { o.Digits = 5 }, "period": func(o *wire.Token) { o.Period = 0 },
			"metadata": func(o *wire.Token) { s := strings.Repeat("x", 129); o.ClientName = &s },
		} {
			o := *c.Input
			change(&o)
			if _, e := o.Encode(); e == nil {
				t.Fatal("invalid authoring accepted", name)
			}
		}
		break
	}
	if _, e := wire.Padded(make([]byte, 1007)); e == nil {
		t.Fatal("capacity")
	}
	if _, e := wire.Derive(nil); e == nil {
		t.Fatal("root width")
	}
	if _, e := wire.VaultID(make([]byte, 87)); e == nil {
		t.Fatal("vault magic")
	}
}

// Wrong independent formulas are compared to committed, Python-checked fields;
// no verifier API or generated expected value participates in these tests.
func TestMutationSensitivity(t *testing.T) {
	var c fixture.Case
	for _, x := range corpus(t) {
		if x.ID == "v1.crypto.token-root.001" {
			c = x
		}
	}
	p := raw(t, c.Semantic)
	x := c.Crypto
	k, e := wire.Derive(raw(t, c.Root))
	if e != nil {
		t.Fatal(e)
	}
	id := raw(t, x.ObjectID)
	pad, e := wire.Padded(p)
	if e != nil {
		t.Fatal(e)
	}
	wrong := func(label string, got []byte, want string) {
		t.Helper()
		if hex.EncodeToString(got) == want {
			t.Fatal("mutation escaped", label)
		}
	}
	for name, offset := range map[string]int{"tag": 1, "width": 3, "endianness": 41} {
		q := bytes.Clone(p)
		q[offset] ^= 1
		wrong(name, q, c.Semantic)
	}
	h := hmac.New(sha256.New, make([]byte, 32))
	h.Write(raw(t, c.Root))
	prk := h.Sum(nil)
	h = hmac.New(sha256.New, prk)
	h.Write([]byte("totipo/v1/OBJECT-id\x01"))
	wrong("id domain", h.Sum(nil), x.IDKey)
	badID := bytes.Clone(id)
	badID[0] ^= 1
	wrong("object id", badID, x.ObjectID)
	h = hmac.New(sha256.New, k.ObjectRoot)
	h.Write(append(append([]byte("totipo/v1/object-KEY"), id...), 1))
	wrong("object key domain", h.Sum(nil), x.ObjectKey)
	seal := func(key, nonce, plain, aad []byte) []byte {
		t.Helper()
		b, e := aes.NewCipher(key)
		if e != nil {
			t.Fatal(e)
		}
		a, e := cipher.NewGCM(b)
		if e != nil {
			t.Fatal(e)
		}
		return a.Seal(nil, nonce, plain, aad)
	}
	nonce := bytes.Clone(id[:12])
	nonce[0] ^= 1
	wrong("nonce", seal(k.ObjectKey(id), nonce, pad, wire.AAD(id)), x.Object)
	aad := wire.AAD(id)
	aad[0] ^= 1
	wrong("aad", seal(k.ObjectKey(id), id[:12], pad, aad), x.Object)
	for name, offset := range map[string]int{"semantic length": 1, "padding": 1007} {
		bad := bytes.Clone(pad)
		bad[offset] ^= 1
		wrong(name, bad, x.Padded)
		wrong(name+" envelope", seal(k.ObjectKey(id), id[:12], bad, wire.AAD(id)), x.Object)
	}
	wrong("object size", raw(t, x.Object)[:1023], x.Object)
}
