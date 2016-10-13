package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

func main() {
	uris := geturisFromEnv()

	switch {
	case len(os.Args) > 2:
		fallthrough
	default:
		fmt.Fprint(os.Stderr, usage())
		os.Exit(1)
	case len(os.Args) == 1:
		doPing(uris)
	case os.Args[1] == "config":
		doConfig(uris)
	case os.Args[1] == "autoconf":
		fmt.Println("no" +
			" (This module is meant to run outside of the node hosting the URIs" +
			" and is to be configured manually.)",
		)
	}
}

func doConfig(uris map[string]string) {
	panic("TODO")
}

func doPing(uris map[string]string) {
	for _, url := range uris {
		info, err := Ping(url)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\nUrl: %s\n", url)
		fmt.Printf("Resolving: %v\n", info.Resolving)
		fmt.Printf("Connecting: %v\n", info.Connecting)
		fmt.Printf("Sending: %v\n", info.Sending)
		fmt.Printf("Waiting: %v\n", info.Waiting)
		fmt.Printf("Receiving: %v\n", info.Receiving)
		fmt.Printf("Total: %v\n", info.Total)
		fmt.Printf("Size: %v\n", info.Size)
	}
}

func usage() string {
	return fmt.Sprintf("Usage: %s [config|autoconf]\n", os.Args[0])
}

func geturisFromEnv() map[string]string {
	uris := make(map[string]string, 0)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "_", 2)
		// Get all envs that look like TARGET_*
		if len(parts) != 2 || parts[0] != "TARGET" {
			continue
		}
		name := parts[1]
		uri := strings.SplitN(env, "=", 2)[1]

		_, err := url.ParseRequestURI(uri)
		if err != nil {
			continue
		}

		uris[name] = uri
	}

	return uris
}
