package pinger

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/DigitalBackstage/munin-http-timing/config"
)

const httpGetTimeout = time.Duration(20 * time.Second)

// ping gets an HTTP URL and returns the request timing information
func ping(name, uri, userAgent string) (*RequestInfo, error) {
	var err error

	info := NewRequestInfo()
	trace := getHTTPTrace(info)
	client := http.Client{
		Timeout: httpGetTimeout,
		// Disable redirect, https://stackoverflow.com/a/38150816
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return info, err
	}
	req.Header.Set("User-Agent", userAgent)

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &trace))

	info.RequestStart(name, uri)
	response, err := client.Do(req)
	if err != nil {
		return info, err
	}

	info.BodySize, err = getResponseBodyBodySize(response)
	if err != nil {
		return info, err
	}

	// Keep this _after_ fetching the whole body because Request.Do returns as
	// soon as the headers are received.
	info.RequestDone(response.StatusCode)

	if info.StatusCode >= 400 {
		err = fmt.Errorf("Got a %d, unable to fetch %s\n", info.StatusCode, uri)
	} else if info.StatusCode >= 300 && info.StatusCode < 400 {
		err = fmt.Errorf("Not following redirection given by %s\n", uri)
	}

	return info, err
}

func getResponseBodyBodySize(r *http.Response) (int, error) {
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

// DoParallelPings calls ping on the given URIs and pushes the result in the
// given queue
func DoParallelPings(config config.Config, queue chan<- *RequestInfo) {
	for name, uri := range config.URIs {
		go func(name, uri string) {
			// Avoid sending all requests at the exact same time
			if config.RandomDelayEnabled {
				time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			}

			info, err := ping(name, uri, config.UserAgent)
			info.Error = err
			queue <- info
		}(name, uri)
	}
}
