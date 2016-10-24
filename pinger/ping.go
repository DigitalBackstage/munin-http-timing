package pinger

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

var stderr = log.New(os.Stderr, "", 0)

const httpGetTimeout = time.Duration(20 * time.Second)

// DoPing does the actual stats gathering (HTTP requests) and prints it for munin
func DoPing(uris map[string]string) (string, error) {
	rand.Seed(time.Now().Unix())

	if len(uris) <= 0 {
		return "", errors.New("No URIs provided.")
	}

	totals := map[string]string{}
	queue := make(chan *RequestInfo, len(uris))
	doParallelPings(uris, queue)

	buf := &bytes.Buffer{}

	for i := 0; i < len(uris); i++ {
		info := <-queue

		fmt.Fprint(buf, info)

		if info.IsOk() {
			totals[info.Name] = fmt.Sprintf("%v", toMillisecond(info.Total))
		} else {
			totals[info.Name] = "U"
		}
	}

	fmt.Fprint(buf, "multigraph timing\n")
	for name, value := range totals {
		fmt.Fprintf(buf, "%s_total.value %v\n", name, value)
	}
	fmt.Fprint(buf, "\n")

	return buf.String(), nil
}

// ping gets an HTTP URL and returns the request timing information
func ping(name, uri string) (*RequestInfo, error) {
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
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &trace))

	info.RequestStart(name, uri)
	response, err := client.Do(req)
	if err != nil {
		return info, err
	}

	info.BodySize, err = getResponseBodyBodySize(response)

	// Keep this _after_ fetching the whole body because Request.Do returns as
	// soon as the headers are received.
	info.RequestDone(response.StatusCode)

	if info.StatusCode >= 400 {
		stderr.Printf("Got a %d, unable to fetch %s\n", info.StatusCode, uri)
	} else if info.StatusCode >= 300 && info.StatusCode < 400 {
		stderr.Printf("Not following redirection given by %s\n", uri)
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

func doParallelPings(uris map[string]string, queue chan<- *RequestInfo) {
	for name, uri := range uris {
		go func(name, uri string) {
			// Avoid sending all requests at the exact same time
			if randomDelayEnabled() {
				time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			}

			info, err := ping(name, uri)
			if err != nil {
				stderr.Println(err)
			}
			queue <- info
		}(name, uri)
	}
}

// RandomDelayEnabled returns true if we should delay parallel requests by a
// random delay
func randomDelayEnabled() bool {
	return os.Getenv("RANDOM_DELAY") == "1"
}
