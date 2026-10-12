package fixture

import (
	"encoding/json"
	"testing"
)

func TestRetainedSourceRefusals(t *testing.T) {
	for _, s := range []string{`{"id":"a","id":"b"}`, `{"totp":{"rows":[],"rows":[]}}`, `{"unknown":0}`, `{} {}`, `{"id":`} {
		var c Case
		if ReadSource([]byte(s), &c) == nil {
			t.Fatal("accepted", s)
		}
	}
}
func TestSerializationOrderAndPresence(t *testing.T) {
	b, e := json.MarshalIndent(Manifest{Format: "f", Protocol: "p", Revision: "r", Cases: []Entry{}}, "", "  ")
	if e != nil || string(b) != "{\n  \"format\": \"f\",\n  \"protocol\": \"p\",\n  \"spec_revision\": \"r\",\n  \"cases\": []\n}" {
		t.Fatal(string(b), e)
	}
	if Hash([]byte("abc")) != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Fatal("sha256")
	}
}
