package main

import (
	"fmt"
	"strings"
)

// DoConfig prints the munin plugin configuration to stdout
func DoConfig(uris map[string]string) {
	printMainGraph(uris)

	for name, uri := range uris {
		printURIGraph(name, uri)
	}
}

// One serie per URI showing total time on the main graph
func printMainGraph(uris map[string]string) {
	fmt.Println(`multigraph timing
graph_title Total time
graph_category network
graph_args --base 1000 -l 0
graph_scale no
graph_info This graph show the timing of the different parts of an HTTP request in miliseconds.
graph_order total resolving connecting sending waiting receiving
graph_vlabel Time (ms)`)

	for name, url := range uris {
		fmt.Printf("%s_total.label %s\n", name, url)
	}

	fmt.Println("")
}

// One serie per timing category per URI
func printURIGraph(name, uri string) {
	p := fmt.Printf
	p("multigraph timing.%s\n", name)
	p("graph_title Timings for %s\n", uri)
	p("graph_vlabel Time (ms)\n")
	printFields(name)
}

func printFields(name string) {
	labels := map[string]string{
		"Total":      "Time spent performing the whole request.",
		"Resolving":  "Time spent resolving the domain name.",
		"Connecting": "Time spent initiating the TCP connection.",
		"Sending":    "Time spent sending the HTTP request.",
		"Waiting":    "Time spent waiting for the first byte of the HTTP response.",
		"Receiving":  "Time spend receiving the request body.",
	}

	fmt.Println("total.warning 1000:2000")
	fmt.Println("total.critical 2000:")

	for label, info := range labels {
		field := strings.ToLower(label)
		fmt.Printf("%s.label %s\n", field, label)
		fmt.Printf("%s.draw LINE\n", field)
		fmt.Printf("%s.info %s\n", field, info)
	}

	fmt.Println("")
}
