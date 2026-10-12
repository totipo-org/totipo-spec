// Package wireoracle preserves the selected independent wire evidence formerly
// supplied by Python. It reads serialized artifacts and imports only stdlib;
// it does not use either implementation's models, constants, or composition.
// It supplies no independent AES-GCM or Argon2 evidence.
package wireoracle

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

type record = map[string]any

// readValue checks duplicate members at every depth before assembling raw maps.
func readValue(d *json.Decoder) (any, error) {
	tok, err := d.Token()
	if err != nil {
		return nil, err
	}
	switch tok {
	case json.Delim('{'):
		m := record{}
		for d.More() {
			k, err := d.Token()
			if err != nil {
				return nil, err
			}
			key := k.(string)
			if _, exists := m[key]; exists {
				return nil, fmt.Errorf("duplicate member %q", key)
			}
			v, err := readValue(d)
			if err != nil {
				return nil, err
			}
			m[key] = v
		}
		_, err := d.Token()
		return m, err
	case json.Delim('['):
		a := []any{}
		for d.More() {
			v, err := readValue(d)
			if err != nil {
				return nil, err
			}
			a = append(a, v)
		}
		_, err := d.Token()
		return a, err
	default:
		return tok, nil
	}
}
func read(path string) (record, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	v, err := readValue(d)
	if err != nil {
		return nil, err
	}
	if _, err := d.Token(); err != io.EOF {
		return nil, fmt.Errorf("trailing JSON: %v", err)
	}
	m, ok := v.(record)
	if !ok {
		return nil, fmt.Errorf("expected object")
	}
	return m, nil
}

// Access failures and predicate failures share a local error path; recover only
// this private failure type, so implementation panics still fail the test.
type failure struct{ message string }

func require(ok bool, what string) {
	if !ok {
		panic(failure{what})
	}
}
func obj(m record, key string) record {
	v, ok := m[key].(record)
	require(ok, key+": object required")
	return v
}
func str(m record, key string) string {
	v, ok := m[key].(string)
	require(ok, key+": string required")
	return v
}
func num(m record, key string) uint64 {
	v, ok := m[key].(json.Number)
	require(ok, key+": integer required")
	n, err := strconv.ParseUint(string(v), 10, 64)
	require(err == nil, key+": invalid integer")
	return n
}
func unhex(m record, key string) []byte {
	b, err := hex.DecodeString(str(m, key))
	require(err == nil, key+": invalid hex")
	return b
}
func unbase(s any) []byte {
	v, ok := s.(string)
	require(ok, "base64 string required")
	b, err := base64.StdEncoding.Strict().DecodeString(v)
	require(err == nil, "invalid base64")
	return b
}
func equal(got, want any, what string) { require(reflect.DeepEqual(got, want), what) }
func mac(k, p []byte) []byte           { h := hmac.New(sha256.New, k); h.Write(p); return h.Sum(nil) }
func expand(k, info []byte) []byte     { return mac(k, append(append([]byte{}, info...), 1)) }
func integer(n uint64, width int) []byte {
	b := make([]byte, width)
	for i := width - 1; i >= 0; i-- {
		b[i] = byte(n)
		n >>= 8
	}
	require(n == 0, "integer width")
	return b
}
func padded(p []byte) []byte {
	require(len(p) <= 1006, "capacity")
	b := make([]byte, 1008)
	binary.BigEndian.PutUint16(b, uint16(len(p)))
	copy(b[2:], p)
	return b
}
func guarded(fn func()) (err error) {
	defer func() {
		if v := recover(); v != nil {
			if f, ok := v.(failure); ok {
				err = fmt.Errorf("%s", f.message)
			} else {
				panic(v)
			}
		}
	}()
	fn()
	return nil
}

func fieldsAndKeys(c record) error {
	return guarded(func() {
		if _, ok := c["crypto"]; !ok {
			return
		}
		p, root := unhex(c, "semantic_hex"), unhex(c, "root_hex")
		prk := mac(make([]byte, 32), root)
		kid := expand(prk, []byte("totipo/v1/object-id"))
		kroot := expand(prk, []byte("totipo/v1/object-key-root"))
		oid := mac(kid, p)
		x := obj(c, "crypto")
		for k, want := range map[string][]byte{
			"id_key_hex": kid, "object_root_key_hex": kroot, "object_id": oid,
			"object_key_hex": expand(kroot, append([]byte("totipo/v1/object-key"), oid...)),
			"nonce_hex":      oid[:12], "aad_hex": append([]byte("totipo/v1/object"), oid...), "padded_plaintext_hex": padded(p),
		} {
			equal(unhex(x, k), want, k)
		}
		require(len(unhex(x, "object_hex")) == 1024, "object size")
		if _, ok := c["input"]; !ok {
			return
		}
		type field struct {
			tag   uint16
			value string
		}
		fields := []field{}
		for offset := 0; offset < len(p); {
			require(len(p)-offset >= 4, "TLV header")
			tag, n := binary.BigEndian.Uint16(p[offset:]), int(binary.BigEndian.Uint16(p[offset+2:]))
			require(n <= len(p)-offset-4, "TLV value")
			fields = append(fields, field{tag, string(p[offset+4 : offset+4+n])})
			offset += 4 + n
		}
		v := obj(c, "input")
		parents, ok := v["parents"].([]any)
		require(ok && len(parents) <= 4, "parents")
		expected := []field{{1, string(unbase(v["identity"]))}, {2, string(integer(uint64(len(parents)), 2))}}
		for _, parent := range parents {
			expected = append(expected, field{3, string(unbase(parent))})
		}
		expected = append(expected, field{4, string(integer(num(v, "status"), 1))}, field{5, str(v, "issuer")}, field{6, str(v, "account")}, field{7, string(integer(num(v, "algorithm"), 1))}, field{8, string(integer(num(v, "digits"), 1))}, field{9, string(integer(num(v, "period"), 4))}, field{10, string(unbase(v["secret"]))})
		if _, ok := v["client_name"]; ok {
			expected = append(expected, field{11, str(v, "client_name")})
		}
		if _, ok := v["client_time"]; ok {
			expected = append(expected, field{12, string(integer(num(v, "client_time"), 8))})
		}
		equal(fields, expected, "TOKEN fields/order/widths")
	})
}
func vault(c record) error {
	return guarded(func() {
		if _, ok := c["bootstrap"]; ok {
			x := obj(c, "bootstrap")
			sum := sha256.Sum256(unhex(x, "record_hex"))
			equal(str(x, "vault_id_hex"), hex.EncodeToString(sum[:]), "VAULT_ID")
		}
		if str(c, "id") == "v1.size.token-max-4.001" {
			x := obj(c, "crypto")
			// Twelve singleton fields plus three extra PARENT_ID records.
			maximum := 15*4 + 32 + 2 + 4*32 + 1 + 256 + 256 + 1 + 1 + 4 + 128 + 128 + 8
			equal(maximum, 1005, "independent maximum sum")
			equal(len(unhex(c, "semantic_hex")), maximum, "maximum semantic bytes")
			equal(num(x, "semantic_length"), uint64(1005), "maximum length")
			b := unhex(x, "padded_plaintext_hex")
			require(len(b) > 0 && b[len(b)-1] == 0, "maximum padding")
		}
	})
}
func postAEAD(c record) error {
	return guarded(func() {
		if _, ok := c["post_aead"]; !ok {
			return
		}
		x := obj(c, "post_aead")
		p := unhex(c, "semantic_hex")
		plain := unhex(x, "encryption_plaintext_hex")
		require(len(plain) == 1008, "post-AEAD plaintext size")
		require(len(unhex(x, "object_hex")) == 1024, "post-AEAD object size")
		kid := expand(mac(make([]byte, 32), unhex(c, "root_hex")), []byte("totipo/v1/object-id"))
		correct := hex.EncodeToString(mac(kid, p))
		n := int(binary.BigEndian.Uint16(plain))
		canonical := padded(p)
		switch str(x, "defect") {
		case "nonzero-padding":
			equal(str(x, "object_id"), correct, "padding ID")
			equal(n, len(p), "padding length")
			equal(plain[:2+n], canonical[:2+n], "padding prefix")
			require(!bytes.Equal(plain[2+n:], canonical[2+n:]), "padding defect")
		case "object-id-mismatch":
			require(str(x, "object_id") != correct, "ID defect")
			equal(plain, canonical, "ID plaintext")
		case "semantic-length-invalid":
			equal(str(x, "object_id"), correct, "length ID")
			equal(n, 1007, "length defect")
			equal(plain[2:], canonical[2:], "length body")
		default:
			require(false, "unknown post-AEAD defect")
		}
	})
}
func root() string {
	if r := os.Getenv("TOTIPO_WIRE_ORACLE_ROOT"); r != "" {
		return r
	}
	return "../../.."
}
func corpus(t *testing.T, check func(record) error) {
	t.Helper()
	r := root()
	m, err := read(filepath.Join(r, "vectors/manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	entries, ok := m["cases"].([]any)
	if !ok {
		t.Fatal("cases array required")
	}
	for _, entry := range entries {
		e, ok := entry.(record)
		if !ok {
			t.Fatal("entry object required")
		}
		path, ok := e["path"].(string)
		if !ok {
			t.Fatal("path required")
		}
		t.Run(path, func(t *testing.T) {
			c, err := read(filepath.Join(r, "vectors", path))
			if err != nil {
				t.Fatal(err)
			}
			if err = check(c); err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestFieldsAndKeyIntermediates(t *testing.T) { corpus(t, fieldsAndKeys) }
func TestVaultIDAndMaximum(t *testing.T)         { corpus(t, vault) }
func TestPostAEADFixtureInputs(t *testing.T)     { corpus(t, postAEAD) }

func TestReaderRejectsAmbiguity(t *testing.T) {
	for _, body := range []string{`{"x":1,"x":2}`, `{"x":{"a":1,"\u0061":2}}`, `{"x":[{"a":1,"a":2}]}`, `{} {}`, `{"x":`, `[]`} {
		p := filepath.Join(t.TempDir(), "case.json")
		if err := os.WriteFile(p, []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := read(p); err == nil {
			t.Fatalf("accepted %s", body)
		}
	}
}
func TestIndependentMutationSensitivity(t *testing.T) {
	for _, key := range []string{"id_key_hex", "object_root_key_hex", "object_id", "object_key_hex", "nonce_hex", "aad_hex", "padded_plaintext_hex", "object_hex"} {
		t.Run(key, func(t *testing.T) {
			c, err := read(filepath.Join(root(), "vectors/cases/encoding/v1.encoding.parents-4.001.json"))
			if err != nil {
				t.Fatal(err)
			}
			obj(c, "crypto")[key] = "00"
			if fieldsAndKeys(c) == nil {
				t.Fatal("wrong key/domain/nonce/AAD/padding/size accepted")
			}
		})
	}
	for _, key := range []string{"identity", "status", "period", "client_name", "client_time"} {
		t.Run("field-"+key, func(t *testing.T) {
			c, err := read(filepath.Join(root(), "vectors/cases/encoding/v1.encoding.parents-4.001.json"))
			if err != nil {
				t.Fatal(err)
			}
			v := obj(c, "input")
			switch key {
			case "identity":
				v[key] = "AA=="
			case "client_name":
				v[key] = "wrong"
			default:
				v[key] = json.Number("65535")
			}
			if fieldsAndKeys(c) == nil {
				t.Fatal("wrong field/tag/width accepted")
			}
		})
	}
	for _, defect := range []string{"nonzero-padding", "object-id-mismatch", "semantic-length-invalid"} {
		t.Run(defect, func(t *testing.T) {
			c, err := read(filepath.Join(root(), "vectors/cases/crypto/v1.crypto."+defect+".001.json"))
			if err != nil {
				t.Fatal(err)
			}
			x := obj(c, "post_aead")
			if defect == "object-id-mismatch" {
				kid := expand(mac(make([]byte, 32), unhex(c, "root_hex")), []byte("totipo/v1/object-id"))
				x["object_id"] = hex.EncodeToString(mac(kid, unhex(c, "semantic_hex")))
			} else {
				x["encryption_plaintext_hex"] = hex.EncodeToString(padded(unhex(c, "semantic_hex")))
			}
			if postAEAD(c) == nil {
				t.Fatal("repaired defect accepted")
			}
		})
	}
}
