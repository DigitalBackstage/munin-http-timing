package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

const httpGetTimeout = time.Duration(20 * time.Second)

// DoPing does the actual stats gathering (HTTP requests) and prints it for munin
func DoPing(uris map[string]string) {
	totals := map[string]int64{}
	queue := make(chan RequestInfo, len(uris))
	doParallelPings(uris, queue)

	for i := 0; i < len(uris); i++ {
		info := <-queue

		info.Print()
		totals[info.Name] = toMillisecond(info.Total)
	}

	fmt.Println("multigraph timing")
	for name, value := range totals {
		fmt.Printf("%s_total.value %v\n", name, value)
	}
}

// ping gets an HTTP URL and returns the request timing information
func ping(name, uri string) (info RequestInfo, err error) {
	trace := getHTTPTrace(&info)
	client := http.Client{
		Timeout: httpGetTimeout,
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &trace))

	info.RequestStart(name, uri)
	response, err := client.Do(req)
	if err != nil {
		return
	}

	info.BodySize, err = getResponseBodyBodySize(response)

	// Keep this _after_ fetching the whole body because Request.Do returns as
	// soon as the headers are received.
	info.RequestDone(response.StatusCode)

	return
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

func doParallelPings(uris map[string]string, queue chan<- RequestInfo) {
	for name, uri := range uris {
		go func(name, uri string) {
			info, err := ping(name, uri)
			if err != nil || info.StatusCode >= 400 {
				fmt.Fprintln(os.Stderr, "Unable to fetch "+uri)
				fmt.Fprintln(os.Stderr, err)
			} else if info.StatusCode >= 300 && info.StatusCode < 400 {
				fmt.Fprintln(os.Stderr, "Not following redirection given by "+uri)
			}
			queue <- info
		}(name, uri)
	}
}
