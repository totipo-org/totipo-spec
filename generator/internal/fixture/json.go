package fixture

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

func Hash(b []byte) string { v := sha256.Sum256(b); return hex.EncodeToString(v[:]) }

// ReadSource rejects duplicate/unknown members and trailing values in retained
// known answers. It neither evaluates nor recomputes those answers.
func ReadSource(b []byte, dst any) error {
	d := json.NewDecoder(bytes.NewReader(b))
	var value func() error
	value = func() error {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch tok {
		case json.Delim('{'):
			seen := map[string]bool{}
			for d.More() {
				key, err := d.Token()
				if err != nil {
					return err
				}
				s, ok := key.(string)
				if !ok || seen[s] {
					return fmt.Errorf("duplicate source member")
				}
				seen[s] = true
				if err := value(); err != nil {
					return err
				}
			}
			_, err = d.Token()
			return err
		case json.Delim('['):
			for d.More() {
				if err := value(); err != nil {
					return err
				}
			}
			_, err = d.Token()
			return err
		}
		return nil
	}
	if err := value(); err != nil {
		return err
	}
	if _, err := d.Token(); err != io.EOF {
		return fmt.Errorf("trailing source data")
	}
	d = json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	return d.Decode(dst)
}
