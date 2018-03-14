package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

// Options contains general options
type Options struct {
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
