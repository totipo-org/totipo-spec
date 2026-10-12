// Package schema interprets only the checked-in schema keyword and pattern subset.
// Contracts remain independently authored JSON; no Go model defines acceptance.
package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// Parse detects duplicate members before constructing each generic object. Numbers
// keep their lexical spelling: 0, 0.0 and 0e0 are distinct for integer typing.
func Parse(b []byte) (any, error) {
	if !utf8.Valid(b) {
		return nil, fmt.Errorf("invalid UTF-8 JSON")
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	var read func() (any, error)
	read = func() (any, error) {
		t, e := d.Token()
		if e != nil {
			return nil, e
		}
		delim, ok := t.(json.Delim)
		if !ok {
			return t, nil
		}
		switch delim {
		case '{':
			obj := map[string]any{}
			for d.More() {
				k, e := d.Token()
				if e != nil {
					return nil, e
				}
				key, ok := k.(string)
				if !ok {
					return nil, fmt.Errorf("invalid member")
				}
				if _, ok := obj[key]; ok {
					return nil, fmt.Errorf("duplicate JSON member %q", key)
				}
				v, e := read()
				if e != nil {
					return nil, e
				}
				obj[key] = v
			}
			if _, e := d.Token(); e != nil {
				return nil, e
			}
			return obj, nil
		case '[':
			arr := []any{}
			for d.More() {
				v, e := read()
				if e != nil {
					return nil, e
				}
				arr = append(arr, v)
			}
			if _, e := d.Token(); e != nil {
				return nil, e
			}
			return arr, nil
		default:
			return nil, fmt.Errorf("unexpected delimiter")
		}
	}
	v, e := read()
	if e != nil {
		return nil, e
	}
	if _, e := d.Token(); e != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	return v, nil
}

var integer = regexp.MustCompile(`^-?(?:0|[1-9][0-9]*)$`)

func intValue(v any) (*big.Int, bool) {
	n, ok := v.(json.Number)
	if !ok || !integer.MatchString(string(n)) {
		return nil, false
	}
	i, ok := new(big.Int).SetString(string(n), 10)
	return i, ok
}
func sorted(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// The allowlist is intentional: future regex syntax needs an explicit parity
// review. Go supports these classes/repetitions/noncapturing groups directly.
var patterns = map[string]bool{
	`^[0-9a-f]{64}$`: true, `^r[0-9]+$`: true,
	`^v1\.[a-z0-9-]+(?:\.[a-z0-9-]+)+$`: true,
	`^(?:[0-9a-f]{2})*$`:                true, `^[0-9]{6,8}$`: true,
	`^[0-9a-f]+$`: true, `^[0-9a-f]{174}$`: true,
}

// Schema is a validated, immutable-by-convention contract tree.
type Schema struct {
	raw                        map[string]any
	props                      map[string]*Schema
	items                      *Schema
	all                        []*Schema
	condition, then, otherwise *Schema
	pattern                    *regexp.Regexp
}

func Compile(b []byte) (*Schema, error) {
	v, e := Parse(b)
	if e != nil {
		return nil, fmt.Errorf("schema source: %w", e)
	}
	return compile(v, "$")
}
func compile(v any, path string) (*Schema, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: schema must be an object", path)
	}
	s := &Schema{raw: m, props: map[string]*Schema{}}
	bad := func(k string) error { return fmt.Errorf("%s: invalid/unsupported schema keyword %s", path, k) }
	child := func(v any, k string) (*Schema, error) { return compile(v, path+"."+k) }
	for _, k := range sorted(m) {
		v := m[k]
		switch k {
		case "$schema", "$id", "title", "description":
			if _, ok := v.(string); !ok {
				return nil, bad(k)
			}
		case "type":
			list, ok := v.([]any)
			if str, yes := v.(string); yes {
				list = []any{str}
				ok = true
			}
			if !ok || len(list) == 0 {
				return nil, bad(k)
			}
			seen := map[string]bool{}
			for _, t := range list {
				str, ok := t.(string)
				if !ok || seen[str] {
					return nil, bad(k)
				}
				switch str {
				case "object", "array", "string", "integer", "boolean", "null":
				default:
					return nil, bad(k)
				}
				seen[str] = true
			}
		case "const": // Any JSON value; equality below preserves Python's bool/number quirk.
		case "enum":
			a, ok := v.([]any)
			if !ok || len(a) == 0 {
				return nil, bad(k)
			}
			for i, x := range a {
				for _, y := range a[:i] {
					if equal(x, y) {
						return nil, bad(k)
					}
				}
			}
		case "properties":
			p, ok := v.(map[string]any)
			if !ok {
				return nil, bad(k)
			}
			for _, name := range sorted(p) {
				c, e := child(p[name], k+"."+name)
				if e != nil {
					return nil, e
				}
				s.props[name] = c
			}
		case "required":
			a, ok := v.([]any)
			if !ok {
				return nil, bad(k)
			}
			seen := map[string]bool{}
			for _, x := range a {
				name, ok := x.(string)
				if !ok || seen[name] {
					return nil, bad(k)
				}
				seen[name] = true
			}
		case "additionalProperties":
			if _, ok := v.(bool); !ok {
				return nil, bad(k)
			}
		case "items":
			c, e := child(v, k)
			if e != nil {
				return nil, e
			}
			s.items = c
		case "minItems", "maxItems", "minLength", "maxLength":
			i, ok := intValue(v)
			if !ok || i.Sign() < 0 {
				return nil, bad(k)
			}
		case "minimum", "maximum":
			if _, ok := intValue(v); !ok {
				return nil, bad(k)
			}
		case "pattern":
			p, ok := v.(string)
			if !ok || !patterns[p] {
				return nil, bad(k)
			}
			r, e := regexp.Compile(p)
			if e != nil {
				return nil, bad(k)
			}
			s.pattern = r
		case "allOf":
			a, ok := v.([]any)
			if !ok || len(a) == 0 {
				return nil, bad(k)
			}
			for i, x := range a {
				c, e := child(x, fmt.Sprintf("allOf[%d]", i))
				if e != nil {
					return nil, e
				}
				s.all = append(s.all, c)
			}
		case "if", "then", "else":
			c, e := child(v, k)
			if e != nil {
				return nil, e
			}
			switch k {
			case "if":
				s.condition = c
			case "then":
				s.then = c
			case "else":
				s.otherwise = c
			}
		default:
			return nil, bad(k)
		}
	}
	if s.condition == nil && (s.then != nil || s.otherwise != nil) {
		return nil, bad("then/else without if")
	}
	return s, nil
}

func numeric(v any) (*big.Rat, bool) {
	switch x := v.(type) {
	case bool:
		if x {
			return big.NewRat(1, 1), true
		}
		return new(big.Rat), true
	case json.Number:
		r, ok := new(big.Rat).SetString(string(x))
		return r, ok
	}
	return nil, false
}

// Python uses == for const/enum, including True == 1 and False == 0. Type
// checking remains strict (bool is never integer). Keep that measured behavior.
func equal(a, b any) bool {
	if x, ok := numeric(a); ok {
		y, ok := numeric(b)
		return ok && x.Cmp(y) == 0
	}
	switch x := a.(type) {
	case nil:
		return b == nil
	case string:
		y, ok := b.(string)
		return ok && x == y
	case []any:
		y, ok := b.([]any)
		if !ok || len(x) != len(y) {
			return false
		}
		for i, v := range x {
			if !equal(v, y[i]) {
				return false
			}
		}
		return true
	case map[string]any:
		y, ok := b.(map[string]any)
		if !ok || len(x) != len(y) {
			return false
		}
		for k, v := range x {
			w, ok := y[k]
			if !ok || !equal(v, w) {
				return false
			}
		}
		return true
	}
	return false
}
func matchesType(v any, t string) bool {
	switch t {
	case "object":
		_, ok := v.(map[string]any)
		return ok
	case "array":
		_, ok := v.([]any)
		return ok
	case "string":
		_, ok := v.(string)
		return ok
	case "integer":
		_, ok := intValue(v)
		return ok
	case "boolean":
		_, ok := v.(bool)
		return ok
	case "null":
		return v == nil
	}
	return false
}
func (s *Schema) ValidateBytes(b []byte) error {
	v, e := Parse(b)
	if e != nil {
		return e
	}
	return s.Validate(v)
}
func (s *Schema) Validate(v any) error { return s.validate(v, "$") }
func (s *Schema) validate(v any, path string) error {
	fail := func(rule string) error { return fmt.Errorf("%s: %s", path, rule) }
	m := s.raw
	if t, ok := m["type"]; ok {
		a, ok := t.([]any)
		if !ok {
			a = []any{t}
		}
		match := false
		for _, x := range a {
			match = match || matchesType(v, x.(string))
		}
		if !match {
			return fail("type")
		}
	}
	if c, ok := m["const"]; ok && !equal(v, c) {
		return fail("const")
	}
	if c, ok := m["enum"]; ok {
		match := false
		for _, x := range c.([]any) {
			match = match || equal(v, x)
		}
		if !match {
			return fail("enum")
		}
	}
	bound := func(n *big.Int, lower, upper string) error {
		if b, ok := m[lower]; ok {
			limit, _ := intValue(b)
			if n.Cmp(limit) < 0 {
				return fail(lower)
			}
		}
		if b, ok := m[upper]; ok {
			limit, _ := intValue(b)
			if n.Cmp(limit) > 0 {
				return fail(upper)
			}
		}
		return nil
	}
	switch x := v.(type) {
	case map[string]any:
		if a, ok := m["required"].([]any); ok {
			for _, k := range a {
				if _, ok := x[k.(string)]; !ok {
					return fail("required member " + k.(string))
				}
			}
		}
		for _, k := range sorted(x) {
			c, known := s.props[k]
			if !known {
				if ap, ok := m["additionalProperties"].(bool); ok && !ap {
					return fail("unknown member " + k)
				}
				continue
			}
			if e := c.validate(x[k], path+"."+k); e != nil {
				return e
			}
		}
	case []any:
		if e := bound(big.NewInt(int64(len(x))), "minItems", "maxItems"); e != nil {
			return e
		}
		if s.items != nil {
			for i, y := range x {
				if e := s.items.validate(y, fmt.Sprintf("%s[%d]", path, i)); e != nil {
					return e
				}
			}
		}
	case string:
		// Python re.search('$') also matches immediately before one final newline.
		if s.pattern != nil && !s.pattern.MatchString(x) && !s.pattern.MatchString(strings.TrimSuffix(x, "\n")) {
			return fail("pattern")
		}
		if e := bound(big.NewInt(int64(utf8.RuneCountInString(x))), "minLength", "maxLength"); e != nil {
			return e
		}
	case json.Number:
		if i, ok := intValue(x); ok {
			if e := bound(i, "minimum", "maximum"); e != nil {
				return e
			}
		}
	}
	for _, c := range s.all {
		if e := c.validate(v, path); e != nil {
			return e
		}
	}
	if s.condition != nil {
		branch := s.then
		if s.condition.validate(v, path) != nil {
			branch = s.otherwise
		}
		if branch != nil {
			return branch.validate(v, path)
		}
	}
	return nil
}
