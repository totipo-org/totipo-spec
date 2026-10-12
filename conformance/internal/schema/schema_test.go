package schema

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func contract(t *testing.T, name string) *Schema {
	t.Helper()
	b, e := os.ReadFile("../../../vectors/" + name + ".schema.json")
	if e != nil {
		t.Fatal(e)
	}
	s, e := Compile(b)
	if e != nil {
		t.Fatal(e)
	}
	return s
}
func parsed(t *testing.T, b []byte) any {
	t.Helper()
	v, e := Parse(b)
	if e != nil {
		t.Fatal(e)
	}
	return v
}
func accepts(t *testing.T, s *Schema, raw string, want bool) {
	t.Helper()
	e := s.ValidateBytes([]byte(raw))
	if (e == nil) != want {
		t.Fatalf("%s: acceptance %v, want %v: %v", raw, e == nil, want, e)
	}
}
func TestParse(t *testing.T) {
	for _, raw := range []string{`{"x":1,"x":2}`, `{"a":[{"x":1,"x":2}]}`, `{"x":1,"\u0078":2}`, `{} {}`, `[1,]`, `{"x":`, `]`, string([]byte{'"', 255, '"'})} {
		if _, e := Parse([]byte(raw)); e == nil {
			t.Fatal("accepted", raw)
		}
	}
	v := parsed(t, []byte(`{"false":false,"zero":0,"empty":"","array":[],"null":null,"X":1,"x":2}`)).(map[string]any)
	if len(v) != 7 || v["false"] != false || v["zero"] != json.Number("0") || v["empty"] != "" || v["null"] != nil {
		t.Fatal(v)
	}
}
func TestSchemaSource(t *testing.T) {
	for _, raw := range []string{
		`null`, `[]`, `true`, `{"type":"number"}`, `{"type":["string","string"]}`, `{"type":[]}`, `{"type":false}`,
		`{"unknown":1}`, `{"properties":{"unused":{"$ref":"x"}}}`, `{"if":{"const":0},"then":{"unknown":1}}`,
		`{"type":"object","type":"array"}`, `{"properties":{"x":{},"x":{}}}`, `{"properties":[]}`, `{"properties":{"x":true}}`,
		`{"required":"x"}`, `{"required":[1]}`, `{"required":["x","x"]}`, `{"additionalProperties":{}}`,
		`{"items":[]}`, `{"minItems":-1}`, `{"maxItems":1.0}`, `{"minLength":false}`, `{"maxLength":null}`,
		`{"minimum":1.5}`, `{"maximum":true}`, `{"pattern":"["}`, `{"pattern":".*"}`, `{"pattern":1}`,
		`{"enum":[]}`, `{"enum":{}}`, `{"enum":[true,1]}`, `{"allOf":[]}`, `{"allOf":[false]}`, `{"allOf":{}}`,
		`{"if":false}`, `{"then":{}}`, `{"else":{}}`, `{"title":1}`, `{"description":false}`, `{"$id":null}`, `{"$schema":[]}`,
	} {
		if _, e := Compile([]byte(raw)); e == nil {
			t.Fatal("accepted schema", raw)
		}
	}
	if _, e := Compile([]byte(`{"$schema":"annotation","$id":"id","title":"t","description":"d","type":["array","null"],"items":{},"required":[]}`)); e != nil {
		t.Fatal(e)
	}
}
func TestExactNumbersAndPythonEquality(t *testing.T) {
	s := contract(t, "case").props["input"].props["client_time"]
	for _, raw := range []string{"0", "18446744073709551615"} {
		accepts(t, s, raw, true)
	}
	for _, raw := range []string{"-1", "18446744073709551616", "1.5", "1.0", "1e0", "0e0", "true", "false", "null", "\"0\""} {
		accepts(t, s, raw, false)
	}
	// Measured Python quirk is limited to equality: type=integer still rejects bool.
	s = contract(t, "case").props["totp"].props["t0"]
	for _, raw := range []string{"0", "-0", "0.0", "0e0", "false"} {
		accepts(t, s, raw, true)
	}
	accepts(t, s, "true", false)
	s = contract(t, "manifest").props["cases"].items.props["normative"]
	for _, raw := range []string{"true", "1", "1.0", "1e0"} {
		accepts(t, s, raw, true)
	}
	accepts(t, s, "false", false)
	for _, bound := range []struct{ field, max string }{{"status", "2"}, {"algorithm", "3"}, {"digits", "8"}, {"period", "4294967295"}} {
		s = contract(t, "case").props["input"].props[bound.field]
		accepts(t, s, bound.max, true)
		accepts(t, s, "0", false)
	}
}
func TestPatternsAndUnicode(t *testing.T) {
	valid := map[string]string{`^[0-9a-f]{64}$`: strings.Repeat("a", 64), `^r[0-9]+$`: "r19", `^v1\.[a-z0-9-]+(?:\.[a-z0-9-]+)+$`: "v1.x.y", `^(?:[0-9a-f]{2})*$`: "00ff", `^[0-9]{6,8}$`: "123456", `^[0-9a-f]+$`: "af", `^[0-9a-f]{174}$`: strings.Repeat("0", 174)}
	for p, v := range valid {
		b, _ := json.Marshal(map[string]any{"type": "string", "pattern": p})
		s, e := Compile(b)
		if e != nil {
			t.Fatal(e)
		}
		for _, x := range []string{v, v + "\n"} {
			b, _ := json.Marshal(x)
			accepts(t, s, string(b), true)
		}
		for _, x := range []string{"!" + v, v + "\n\n", v + "\r\n", v + "é"} {
			b, _ := json.Marshal(x)
			accepts(t, s, string(b), false)
		}
	}
	s, e := Compile([]byte(`{"type":"string","minLength":1,"maxLength":2}`))
	if e != nil {
		t.Fatal(e)
	}
	for _, raw := range []string{`"é"`, `"😀"`, `"e\u0301"`, `"\ud800"`} {
		accepts(t, s, raw, true)
	}
	accepts(t, s, `""`, false)
	accepts(t, s, `"abc"`, false)
}
func TestConditionsAndArrays(t *testing.T) {
	s, e := Compile([]byte(`{"type":"object","allOf":[{"required":["x"]}],"if":{"properties":{"x":{"const":0}}},"then":{"required":["yes"]},"else":{"required":["no"]}}`))
	if e != nil {
		t.Fatal(e)
	}
	for _, raw := range []string{`{"x":0,"yes":false}`, `{"x":1,"no":[]}`} {
		accepts(t, s, raw, true)
	}
	for _, raw := range []string{`{}`, `{"x":0,"no":true}`, `{"x":1,"yes":true}`} {
		accepts(t, s, raw, false)
	}
	s, e = Compile([]byte(`{"type":"array","minItems":1,"maxItems":2,"items":{"type":"boolean"}}`))
	if e != nil {
		t.Fatal(e)
	}
	for _, raw := range []string{`[]`, `[true,false,true]`, `[0]`, `null`} {
		accepts(t, s, raw, false)
	}
	accepts(t, s, `[false]`, true)
}

// Walk each populated checked-in contract, including selected conditional
// branches. Every required presence check is exercised independently of Go
// structs. Mutations are in memory and restored before walking further.
func TestCorpusRequiredPresence(t *testing.T) {
	ms, cs := contract(t, "manifest"), contract(t, "case")
	b, e := os.ReadFile("../../../vectors/manifest.json")
	if e != nil {
		t.Fatal(e)
	}
	m := parsed(t, b)
	count := 0
	zero, falseValue, emptyArray, emptyString := false, false, false, false
	var walk func(*Schema, any, *Schema, any, string)
	walk = func(s *Schema, v any, rootSchema *Schema, root any, path string) {
		if obj, ok := v.(map[string]any); ok {
			for _, r := range required(s) {
				name := r.(string)
				original, ok := obj[name]
				if !ok {
					t.Fatalf("%s.%s absent in valid corpus", path, name)
				}
				delete(obj, name)
				if e := rootSchema.Validate(root); e == nil {
					t.Fatalf("accepted missing %s.%s", path, name)
				}
				obj[name] = original
				count++
				switch x := original.(type) {
				case bool:
					falseValue = falseValue || !x
				case json.Number:
					zero = zero || x == "0"
				case []any:
					emptyArray = emptyArray || len(x) == 0
				case string:
					emptyString = emptyString || x == ""
				}
			}
			for _, name := range sorted(obj) {
				if child := s.props[name]; child != nil {
					walk(child, obj[name], rootSchema, root, path+"."+name)
				}
			}
		}
		if a, ok := v.([]any); ok && s.items != nil {
			for i, x := range a {
				walk(s.items, x, rootSchema, root, path+"["+string(rune('0'+i%10))+"]")
			}
		}
		for _, sub := range s.all {
			walk(sub, v, rootSchema, root, path)
		}
		if s.condition != nil {
			branch := s.then
			if s.condition.Validate(v) != nil {
				branch = s.otherwise
			}
			if branch != nil {
				walk(branch, v, rootSchema, root, path)
			}
		}
	}
	if e := ms.Validate(m); e != nil {
		t.Fatal(e)
	}
	walk(ms, m, ms, m, "manifest")
	entries := m.(map[string]any)["cases"].([]any)
	for _, entry := range entries {
		path := entry.(map[string]any)["path"].(string)
		b, e := os.ReadFile(filepath.Join("../../../vectors", path))
		if e != nil {
			t.Fatal(e)
		}
		v := parsed(t, b)
		if e := cs.Validate(v); e != nil {
			t.Fatal(path, e)
		}
		walk(cs, v, cs, v, path)
	}
	if !zero || !falseValue || !emptyArray || !emptyString {
		t.Fatal("presence representatives missing", zero, falseValue, emptyArray, emptyString)
	}
	t.Logf("manifest + %d cases; %d required-member removals rejected", len(entries), count)
}
func required(s *Schema) []any { a, _ := s.raw["required"].([]any); return a }
