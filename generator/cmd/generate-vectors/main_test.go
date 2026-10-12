package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestDeterministicWrite(t *testing.T) {
	p := filepath.Join(t.TempDir(), "nested/fixture.json")
	v := struct {
		Value string `json:"value"`
	}{"é\n"}
	write(p, v)
	a, e := os.ReadFile(p)
	if e != nil {
		t.Fatal(e)
	}
	write(p, v)
	b, e := os.ReadFile(p)
	if e != nil || !bytes.Equal(a, b) || string(b) != "{\n  \"value\": \"é\\n\"\n}\n" {
		t.Fatal("serialization", string(b), e)
	}
}

func TestAuthoredGraphIdentity(t *testing.T) {
	old := cases
	defer func() { cases = old }()
	cases = nil
	graphCase("sequential", add(node("A", "T", "X")), result("T", []string{"A"}, []string{"X"}, nil))
	if len(cases) != 1 || cases[0].ID != "v1.graph.sequential.001" || cases[0].Operation != "graph" {
		t.Fatal("authored graph identity changed")
	}
	if got := cases[0].Graph.Steps[1].Want; len(got.HeadObjects) != 1 || got.HeadObjects[0].ID != "A" || got.Conflicting {
		t.Fatal("authored expectation record changed")
	}
}
