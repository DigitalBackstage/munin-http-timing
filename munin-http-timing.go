package main

import (
	"fmt"
	"net/http"
	"net/http/httptrace"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		panic("Only one (required) argument: URL")
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", args[0], nil)
	if err != nil {
		panic(err)
	}

	times := TimingInfo{}

	trace := &httptrace.ClientTrace{
		ConnectDone: func(network, addr string, err error) {
			times.ConnectDone()
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			times.WroteRequest()
		},
		GotFirstResponseByte: func() {
			times.GotFirstResponseByte()
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	times.Start()
	_, err = client.Do(req)
	times.End()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Connecting: %v\n", times.Connecting)
	fmt.Printf("Sending: %v\n", times.Sending)
	fmt.Printf("Waiting: %v\n", times.Waiting)
	fmt.Printf("Receiving: %v\n", times.Receiving)
}
