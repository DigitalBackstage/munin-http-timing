package main

import (
	"fmt"
	"time"
)

// RequestInfo contains the different timings involved in sending
// an HTTP request and its response
type RequestInfo struct {
	Name       string
	URI        string
	StatusCode int

	start                time.Time
	dnsStart             time.Time
	dnsDone              time.Time
	connectDone          time.Time
	wroteRequest         time.Time
	gotFirstResponseByte time.Time

	Resolving  time.Duration
	Connecting time.Duration
	Sending    time.Duration
	Waiting    time.Duration
	Receiving  time.Duration
	Total      time.Duration

	BodySize int
}

// RequestStart starts the timer
func (t *RequestInfo) RequestStart(name, uri string) {
	t.start = time.Now()
	t.Name = name
	t.URI = uri
}

// Print prints the timings following the Munin multigraph protocol
// It prints the fields in a specific order, it must match the one in
// graphOrder in config.go
func (t *RequestInfo) Print() {
	fmt.Printf("multigraph timing.%s\n", t.Name)

	if t.StatusCode > 400 {
		fmt.Println("resolving.value U")
		fmt.Println("connecting.value U")
		fmt.Println("sending.value U")
		fmt.Println("waiting.value U")
		fmt.Println("receiving.value U")
	} else {
		fmt.Printf("resolving.value %v\n", toMillisecond(t.Resolving))
		fmt.Printf("connecting.value %v\n", toMillisecond(t.Connecting))
		fmt.Printf("sending.value %v\n", toMillisecond(t.Sending))
		fmt.Printf("waiting.value %v\n", toMillisecond(t.Waiting))
		fmt.Printf("receiving.value %v\n", toMillisecond(t.Receiving))
	}

	fmt.Println("")
}

// ConnectDone sets the connection time
func (t *RequestInfo) ConnectDone() {
	t.connectDone = time.Now()

	// If there was no DNS request (eg. IP), use start time
	if t.dnsDone.IsZero() {
		t.Connecting = t.connectDone.Sub(t.start)
	} else {
		t.Connecting = t.connectDone.Sub(t.dnsDone)
	}
}

// WroteRequest sets the writing time
func (t *RequestInfo) WroteRequest() {
	t.wroteRequest = time.Now()
	t.Sending = t.wroteRequest.Sub(t.connectDone)
}

// GotFirstResponseByte sets the waiting time
func (t *RequestInfo) GotFirstResponseByte() {
	t.gotFirstResponseByte = time.Now()
	t.Waiting = t.gotFirstResponseByte.Sub(t.wroteRequest)
}

// RequestDone sets the receiving time
func (t *RequestInfo) RequestDone(statusCode int) {
	t.Receiving = time.Now().Sub(t.gotFirstResponseByte)
	t.Total = time.Now().Sub(t.start)
	t.StatusCode = statusCode
}

// DNSStart starts the resolution timer
func (t *RequestInfo) DNSStart() {
	t.dnsStart = time.Now()
}

// DNSDone sets the resolving time
func (t *RequestInfo) DNSDone() {
	t.dnsDone = time.Now()
	t.Resolving = t.dnsDone.Sub(t.dnsStart)
}

func toMillisecond(d time.Duration) int64 {
	return int64(d / time.Millisecond)
}
