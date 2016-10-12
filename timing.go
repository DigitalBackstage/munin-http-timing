package main

import "time"

// TimingInfo contains the different timings involved in sending
// an HTTP request and its response
type TimingInfo struct {
	start time.Time

	connectDone          time.Time
	wroteRequest         time.Time
	gotFirstResponseByte time.Time
	done                 time.Time

	Connecting time.Duration
	Sending    time.Duration
	Waiting    time.Duration
	Receiving  time.Duration
}

// Start starts the timer
func (t *TimingInfo) Start() {
	t.start = time.Now()
}

// ConnectDone sets the connection time
func (t *TimingInfo) ConnectDone() {
	t.connectDone = time.Now()
	t.Connecting = t.connectDone.Sub(t.start)
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
	t.done = time.Now()
	t.Receiving = t.done.Sub(t.gotFirstResponseByte)
}
