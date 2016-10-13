package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"time"
)

const httpGetTimeout = time.Duration(20 * time.Second)

// RequestInfo contains the different timings involved in sending
// an HTTP request and its response
type RequestInfo struct {
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

	Size int
}

// RequestStart starts the timer
func (t *RequestInfo) RequestStart() {
	t.start = time.Now()
}

// ConnectDone sets the connection time
func (t *RequestInfo) ConnectDone() {
	t.connectDone = time.Now()
	t.Connecting = t.connectDone.Sub(t.dnsDone)
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
func (t *RequestInfo) RequestDone() {
	t.Receiving = time.Now().Sub(t.gotFirstResponseByte)
	t.Total = time.Now().Sub(t.start)
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

// Ping gets an HTTP URL and returns the request timing information
func Ping(url string) (info RequestInfo, err error) {
	trace := getHTTPTrace(&info)
	client := http.Client{
		Timeout: httpGetTimeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &trace))

	info.RequestStart()
	response, err := client.Do(req)
	info.RequestDone()
	if err != nil {
		return
	}

	info.Size, err = getResponseBodySize(response)
	return
}

func getResponseBodySize(r *http.Response) (int, error) {
	body, err := ioutil.ReadAll(r.Body)
	return len(body), err
}

func getHTTPTrace(info *RequestInfo) httptrace.ClientTrace {
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
