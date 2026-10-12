package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"totipo/conformance/internal/vectors"
)

func run(args []string, out, errout io.Writer) int {
	flags := flag.NewFlagSet("totipo-conformance", flag.ContinueOnError)
	flags.SetOutput(errout)
	root := flags.String("root", ".", "repository root")
	verify := flags.Bool("verify-only", false, "verify manifest, schemas pins and cases without executing")
	if e := flags.Parse(args); e != nil {
		if e == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(errout, "unexpected arguments")
		return 2
	}
	m, cases, e := vectors.Read(*root)
	if e == nil {
		e = vectors.VerifyProfile(*root, m)
	}
	if e != nil {
		fmt.Fprintln(errout, "FAIL:", e)
		return 1
	}
	if *verify {
		fmt.Fprintf(out, "PASS: %s manifest, profile and %d case files verified\n", vectors.ExpectedRevision, len(cases))
		return 0
	}
	failures := 0
	for _, c := range cases {
		if e := vectors.Run(c); e != nil {
			fmt.Fprintf(errout, "FAIL %s: %v\n", c.ID, e)
			failures++
		}
	}
	if failures > 0 {
		return 1
	}
	fmt.Fprintf(out, "PASS: v1/%s moving-pre-rc: %d cases\n", vectors.ExpectedRevision, len(cases))
	return 0
}
func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
