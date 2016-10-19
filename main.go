package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

var stdout = log.New(os.Stdout, "", 0)
var stderr = log.New(os.Stderr, "", 0)

// Filled during build
var version string

func main() {
	var err error
	uris := getURIsFromEnv()

	// https://munin.readthedocs.io/en/latest/plugin/protocol-dirtyconfig.html#plugin-protocol-dirtyconfig
	dirtyConfig := os.Getenv("MUNIN_CAP_DIRTYCONFIG") == "1"

	switch {
	case len(os.Args) > 2:
		fallthrough
	default:
		stderr.Print(usage())
		os.Exit(1)
	case len(os.Args) == 1:
		err = DoPing(uris)
	case os.Args[1] == "config":
		err = DoConfig(uris)
		if dirtyConfig && err == nil {
			err = DoPing(uris)
		}
	case os.Args[1] == "autoconf":
		stdout.Println("no" +
			" (This module is meant to run outside of the node hosting the URIs" +
			" and is to be configured manually.)",
		)
	case os.Args[1] == "version":
		stdout.Println(version)
	}

	if err != nil {
		stderr.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

// usage returns the usage string (help)
func usage() string {
	return fmt.Sprintf("Usage: %s [config|autoconf|version]\n", os.Args[0])
}

// getURIsFromEnv returns a map associating names to urls from the process env vars
// Only vars prefixed with 'TARGET_' will be used, eg.
// TARGET_EXAMPLE=https://example.com/ will register the URI with "example"
// as the name.
func getURIsFromEnv() map[string]string {
	uris := make(map[string]string, 0)

	for _, env := range os.Environ() {
		// Filter TARGET_*
		parts := strings.SplitN(env, "_", 2)
		if len(parts) != 2 || parts[0] != "TARGET" {
			continue
		}

		// Check for values
		name := strings.ToLower(strings.Split(parts[1], "=")[0])
		uri := strings.SplitN(env, "=", 2)[1]
		if len(name) <= 0 || len(uri) <= 0 {
			continue
		}

		// Check if URI is valid
		_, err := url.ParseRequestURI(uri)
		if err != nil {
			stderr.Printf("Invalid URI: %s\n", env)
			continue
		}

		uris[name] = uri
	}

	return uris
}

// RandomDelayEnabled returns true if we should delay parallel requests by a
// random delay
func RandomDelayEnabled() bool {
	return os.Getenv("RANDOM_DELAY") == "1"
}
