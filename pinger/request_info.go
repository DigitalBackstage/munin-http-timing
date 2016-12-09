package pinger

import (
	"sync"
	"time"
)

// RequestInfo contains the different timings involved in sending
// an HTTP request and its response
type RequestInfo struct {
	Name       string
	URI        string
	StatusCode int
	Error      error

	lock *sync.RWMutex

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

// NewRequestInfo creates a new RequestInfo
func NewRequestInfo() *RequestInfo {
	r := &RequestInfo{}
	r.lock = new(sync.RWMutex)

	return r
}

// IsOk returns true if the request succeeded
func (t *RequestInfo) IsOk() bool {
	return t.StatusCode < 400
}

// RequestStart starts the timer
func (t *RequestInfo) RequestStart(name, uri string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.start = time.Now()
	t.Name = name
	t.URI = uri
}

// Lock locks the RequestInfo mutex
func (t *RequestInfo) Lock() {
	t.lock.Lock()
}

// Unlock unlocks the RequestInfo mutex
func (t *RequestInfo) Unlock() {
	t.lock.Unlock()
}

// ConnectDone sets the connection time
func (t *RequestInfo) ConnectDone() {
	t.lock.Lock()
	defer t.lock.Unlock()

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
	t.lock.Lock()
	defer t.lock.Unlock()

	t.wroteRequest = time.Now()

	// If there was no connection (eg. hitting the same server twice), use start time
	if t.connectDone.IsZero() {
		t.Sending = t.wroteRequest.Sub(t.start)
	} else {
		t.Sending = t.wroteRequest.Sub(t.connectDone)
	}
}

// GotFirstResponseByte sets the waiting time
func (t *RequestInfo) GotFirstResponseByte() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.gotFirstResponseByte = time.Now()
	t.Waiting = t.gotFirstResponseByte.Sub(t.wroteRequest)
}

// RequestDone sets the receiving time
func (t *RequestInfo) RequestDone(statusCode int) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.Receiving = time.Now().Sub(t.gotFirstResponseByte)
	t.Total = time.Now().Sub(t.start)
	t.StatusCode = statusCode
}

// DNSStart starts the resolution timer
func (t *RequestInfo) DNSStart() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.dnsStart = time.Now()
}

// DNSDone sets the resolving time
func (t *RequestInfo) DNSDone() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.dnsDone = time.Now()
	t.Resolving = t.dnsDone.Sub(t.dnsStart)
}
