package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [config|autoconf]\n", os.Args[0])
		os.Exit(1)
	}

	url := "http://example.com/"
	info, err := Ping(url)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Resolving: %v\n", info.Resolving)
	fmt.Printf("Connecting: %v\n", info.Connecting)
	fmt.Printf("Sending: %v\n", info.Sending)
	fmt.Printf("Waiting: %v\n", info.Waiting)
	fmt.Printf("Receiving: %v\n", info.Receiving)
	fmt.Printf("Total: %v\n", info.Total)
}
