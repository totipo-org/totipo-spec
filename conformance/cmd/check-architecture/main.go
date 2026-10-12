// Command check-architecture checks both source trees and their resolved closures.
package main

import (
	"flag"
	"fmt"
	"os"
	"totipo/conformance/internal/repopolicy"
)

func main() {
	root := flag.String("root", ".", "repository root")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "unexpected arguments")
		os.Exit(2)
	}
	if err := repopolicy.Check(*root); err != nil {
		fmt.Fprintln(os.Stderr, "FAIL architecture:", err)
		os.Exit(1)
	}
	fmt.Println("PASS: two independent repository Go modules; generator stdlib-only")
}
