package main

import (
	"net/http"
	"net/http/httptrace"
	"time"
)

// TimingInfo contains the different timings involved in sending
// an HTTP request and its response
type TimingInfo struct {
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
}

// Start starts the timer
func (t *TimingInfo) Start() {
	t.start = time.Now()
}

// ConnectDone sets the connection time
func (t *TimingInfo) ConnectDone() {
	t.connectDone = time.Now()
	t.Connecting = t.connectDone.Sub(t.dnsDone)
}

// WroteRequest sets the writing time
func (t *TimingInfo) WroteRequest() {
	t.wroteRequest = time.Now()
	t.Sending = t.wroteRequest.Sub(t.connectDone)
}

// GotFirstResponseByte sets the waiting time
func (t *TimingInfo) GotFirstResponseByte() {
	t.gotFirstResponseByte = time.Now()
	t.Waiting = t.gotFirstResponseByte.Sub(t.wroteRequest)
}

// End sets the receiving time
func (t *TimingInfo) End() {
	t.Receiving = time.Now().Sub(t.gotFirstResponseByte)
	t.Total = time.Now().Sub(t.start)
}

// DNSStart starts the resolution timer
func (t *TimingInfo) DNSStart() {
	t.dnsStart = time.Now()
}

// DNSDone sets the resolving time
func (t *TimingInfo) DNSDone() {
	t.dnsDone = time.Now()
	t.Resolving = t.dnsDone.Sub(t.dnsStart)
}

// Ping gets an HTTP URL and returns the request timing information
func Ping(url string) (info TimingInfo, err error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	trace := getHTTPTrace(&info)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &trace))

	info.Start()
	_, err = client.Do(req)
	info.End()

	return
}

func getHTTPTrace(info *TimingInfo) httptrace.ClientTrace {
	return httptrace.ClientTrace{
		DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
			info.DNSStart()
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			info.DNSDone()
		},
		ConnectDone: func(network, addr string, err error) {
			info.ConnectDone()
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			info.WroteRequest()
		},
		GotFirstResponseByte: func() {
			info.GotFirstResponseByte()
		},
	}
}
