package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rjeczalik/bsondiff"
)

func main() {
	p := &bsondiff.Program{
		Stdout: os.Stdout,
	}

	if err := p.Run(flag.CommandLine, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "bsondiff:", err)
		os.Exit(1)
	}
}
