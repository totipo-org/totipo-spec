package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"totipo/conformance/internal/specpolicy"
)

func run(args []string, out, errout io.Writer) int {
	flags := flag.NewFlagSet("check-spec", flag.ContinueOnError)
	flags.SetOutput(errout)
	root := flags.String("root", ".", "repository root (only spec and archived history are read)")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(errout, "unexpected arguments")
		return 2
	}
	if err := specpolicy.Check(*root); err != nil {
		fmt.Fprintln(errout, "FAIL:", err)
		return 1
	}
	fmt.Fprintln(out, "PASS: Totipo v1/r19 structural checks")
	return 0
}
func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
