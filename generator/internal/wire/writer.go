// Package wire composes fixture bytes from v1/r19 sections 9–12. It contains no
// accepting TOKEN reader, graph evaluator, bootstrap KDF, or TOTP implementation.
package wire

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"unicode/utf8"
)

type Token struct {
	Identity   []byte   `json:"identity"`
	Parents    [][]byte `json:"parents"`
	Status     byte     `json:"status"`
	Issuer     string   `json:"issuer"`
	Account    string   `json:"account"`
	Algorithm  byte     `json:"algorithm"`
	Digits     byte     `json:"digits"`
	Period     uint32   `json:"period"`
	Secret     []byte   `json:"secret"`
	ClientName *string  `json:"client_name,omitempty"`
	ClientTime *uint64  `json:"client_time,omitempty"`
}

// Integer emits an unsigned big-endian value of the requested fixture width.
func Integer(n uint64, width int) []byte {
	b := make([]byte, width)
	for i := width - 1; i >= 0; i-- {
		b[i] = byte(n)
		n >>= 8
	}
	return b
}
func U16(n uint16) []byte { return Integer(uint64(n), 2) }
func U32(n uint32) []byte { return Integer(uint64(n), 4) }
func U64(n uint64) []byte { return Integer(n, 8) }
func Encode(tag uint16, value []byte) ([]byte, error) {
	if len(value) > 65535 {
		return nil, fmt.Errorf("fixture TLV exceeds u16 length")
	}
	b := append(U16(tag), U16(uint16(len(value)))...)
	return append(b, value...), nil
}

// SplitField is only a mutation utility for already authored fixture TLVs.
// It does not recognize or accept TOKEN grammar.
type Field struct {
	Tag   uint16
	Value []byte
}

func SplitField(b []byte) (Field, []byte, error) {
	if len(b) < 4 {
		return Field{}, nil, fmt.Errorf("short fixture field")
	}
	size := int(b[2])<<8 | int(b[3])
	if size > len(b)-4 {
		return Field{}, nil, fmt.Errorf("short fixture value")
	}
	return Field{uint16(b[0])<<8 | uint16(b[1]), b[4 : 4+size]}, b[4+size:], nil
}
func (t Token) Encode() ([]byte, error) {
	textOK := func(s string, limit int) bool { return utf8.ValidString(s) && len(s) <= limit }
	if len(t.Identity) != 32 || len(t.Parents) > 4 || t.Status < 1 || t.Status > 2 ||
		!textOK(t.Issuer, 256) || !textOK(t.Account, 256) || t.Algorithm < 1 || t.Algorithm > 3 ||
		t.Digits < 6 || t.Digits > 8 || t.Period == 0 || len(t.Secret) == 0 || len(t.Secret) > 128 ||
		(t.ClientName != nil && !textOK(*t.ClientName, 128)) {
		return nil, fmt.Errorf("invalid authoring fields")
	}
	for i, p := range t.Parents {
		if len(p) != 32 || (i > 0 && bytes.Compare(t.Parents[i-1], p) >= 0) {
			return nil, fmt.Errorf("invalid authoring parents")
		}
	}
	// Build a field sequence first, preserving absence independently of zero values.
	fields := []Field{{1, t.Identity}, {2, U16(uint16(len(t.Parents)))}}
	for _, p := range t.Parents {
		fields = append(fields, Field{3, p})
	}
	fields = append(fields, []Field{{4, []byte{t.Status}}, {5, []byte(t.Issuer)}, {6, []byte(t.Account)},
		{7, []byte{t.Algorithm}}, {8, []byte{t.Digits}}, {9, U32(t.Period)}, {10, t.Secret}}...)
	if t.ClientName != nil {
		fields = append(fields, Field{11, []byte(*t.ClientName)})
	}
	if t.ClientTime != nil {
		fields = append(fields, Field{12, U64(*t.ClientTime)})
	}
	size := 0
	for _, f := range fields {
		size += 4 + len(f.Value)
	}
	if size > 1006 {
		return nil, fmt.Errorf("semantic capacity exceeded")
	}
	out := make([]byte, 0, size)
	for _, f := range fields {
		b, err := Encode(f.Tag, f.Value)
		if err != nil {
			return nil, err
		}
		out = append(out, b...)
	}
	return out, nil
}

// All HKDF outputs used here are one SHA-256 block. RFC 5869 expand is
// HMAC(PRK, info || 0x01); extract is HMAC(zero salt, root).
func MAC(key, data []byte) []byte { h := hmac.New(sha256.New, key); h.Write(data); return h.Sum(nil) }
func block(prk []byte, domain string, suffix []byte) []byte {
	info := append([]byte(domain), suffix...)
	return MAC(prk, append(info, 1))
}

type Keys struct{ ID, ObjectRoot []byte }

func Derive(root []byte) (Keys, error) {
	if len(root) != 32 {
		return Keys{}, fmt.Errorf("root must be 32 bytes")
	}
	extracted := MAC(make([]byte, 32), root)
	return Keys{block(extracted, "totipo/v1/object-id", nil), block(extracted, "totipo/v1/object-key-root", nil)}, nil
}
func (k Keys) ObjectKey(id []byte) []byte { return block(k.ObjectRoot, "totipo/v1/object-key", id) }
func AAD(id []byte) []byte                { return append([]byte("totipo/v1/object"), id...) }
func Padded(p []byte) ([]byte, error) {
	if len(p) > 1006 {
		return nil, fmt.Errorf("semantic capacity exceeded")
	}
	return append(append(U16(uint16(len(p))), p...), make([]byte, 1006-len(p))...), nil
}
func (k Keys) Seal(p []byte) (string, []byte, error) {
	plain, err := Padded(p)
	if err != nil {
		return "", nil, err
	}
	id := MAC(k.ID, p)
	b, err := k.Encrypt(id, plain)
	return hex.EncodeToString(id), b, err
}

// Encrypt accepts intentionally malformed padded input for post-AEAD fixtures.
func (k Keys) Encrypt(id, plain []byte) ([]byte, error) {
	if len(id) != 32 || len(plain) != 1008 {
		return nil, fmt.Errorf("fixture envelope width")
	}
	c, err := aes.NewCipher(k.ObjectKey(id))
	if err != nil {
		return nil, err
	}
	a, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	return a.Seal(nil, id[:12], plain, AAD(id)), nil
}
func VaultID(record []byte) ([]byte, error) {
	if len(record) != 87 || string(record[:11]) != "TOTIPO-VLT\x01" {
		return nil, fmt.Errorf("fixture VAULT layout")
	}
	h := sha256.Sum256(record)
	return h[:], nil
}

// Malform changes exactly one post-authentication predicate (sections 11,13).
func Malform(k Keys, p []byte, defect string) (id, plain, encrypted []byte, err error) {
	plain, err = Padded(p)
	if err != nil {
		return
	}
	id = MAC(k.ID, p)
	switch defect {
	case "nonzero-padding":
		plain[1007] = 1
	case "object-id-mismatch":
		id[0] ^= 1
	case "semantic-length-invalid":
		binary.BigEndian.PutUint16(plain, 1007)
	default:
		err = fmt.Errorf("unknown defect")
		return
	}
	encrypted, err = k.Encrypt(id, plain)
	return
}
