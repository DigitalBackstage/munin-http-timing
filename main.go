package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DigitalBackstage/munin-http-timing/config"
	"github.com/DigitalBackstage/munin-http-timing/munin"
)

var stdout = log.New(os.Stdout, "", 0)
var stderr = log.New(os.Stderr, "", 0)

// Filled during build
var version string

func main() {
	var err error
	var out string

	config := config.NewConfigFromEnv()
	config.SetSuffixFromArg0(os.Args[0])

	switch {
	case len(os.Args) > 2:
		fallthrough
	default:
		stderr.Print(usage())
		os.Exit(1)
	case len(os.Args) == 1:
		out, err = munin.DoPing(config)
	case os.Args[1] == "config":
		err = munin.DoConfig(config)
		if config.ConfigAndPing && err == nil {
			out, err = munin.DoPing(config)
		}
	case os.Args[1] == "autoconf":
		out = "no" +
			" (This module is meant to run outside of the node hosting the URIs" +
			" and is to be configured manually.)\n"
	case os.Args[1] == "version":
		out = version + "\n"
	}

	if err != nil {
		stderr.Println(err)
		os.Exit(1)
	}

	stdout.Print(out)
	os.Exit(0)
}

// usage returns the usage string (help)
func usage() string {
	return fmt.Sprintf("Usage: %s [config|autoconf|version]\n", os.Args[0])
}
