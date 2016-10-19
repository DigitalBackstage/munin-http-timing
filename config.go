package main

import (
	"errors"
	"strings"
)

// Field order for graph_order and field descriptions
// This is also hard-coded in Ping
var graphOrder = []string{
	"resolving",
	"connecting",
	"sending",
	"waiting",
	"receiving",
}

// DoConfig prints the munin plugin configuration to stdout
func DoConfig(uris map[string]string) error {
	if len(uris) <= 0 {
		return errors.New("No URIs provided.")
	}

	printMainGraph(uris)

	for name, uri := range uris {
		printURIGraph(name, uri)
	}

	return nil
}

// One serie per URI showing total time on the main graph
func printMainGraph(uris map[string]string) {
	p := stdout.Println
	p("multigraph timing")
	p("graph_title Total time")
	p("graph_category network")
	p("graph_args --base 1000 -l 0")
	p("graph_scale no")
	p("graph_info This graph show the timing of the different parts of an HTTP request in miliseconds.")
	p("graph_vlabel Time (ms)")
	stdout.Printf("graph_order %s\n", strings.Join(graphOrder, " "))

	for name, url := range uris {
		stdout.Printf("%s_total.label %s\n", name, url)
	}

	p("")
}

// One serie per timing category per URI
func printURIGraph(name, uri string) {
	p := stdout.Printf
	p("multigraph timing.%s\n", name)
	p("graph_title Timings for %s\n", uri)
	p("graph_vlabel Time (ms)\n")
	printFields(name)
}

func printFields(name string) {
	labels := map[string]string{
		"Resolving":  "Time spent resolving the domain name.",
		"Connecting": "Time spent initiating the TCP connection.",
		"Sending":    "Time spent sending the HTTP request.",
		"Waiting":    "Time spent waiting for the first byte of the HTTP response.",
		"Receiving":  "Time spend receiving the request body.",
	}

	for _, field := range graphOrder {
		label := strings.ToUpper(field[0:1]) + field[1:]
		stdout.Printf("%s.label %s\n", field, label)

		if field == "resolving" {
			stdout.Printf("%s.draw AREA\n", field)
		} else {
			stdout.Printf("%s.draw STACK\n", field)
		}
		stdout.Printf("%s.info %s\n", field, labels[label])
	}

	stdout.Println("")
}
