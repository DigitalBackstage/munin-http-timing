package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// Requested URIs will be appended in order here
var TestServerPings []string
var TestServerPort int
var TestServerBaseURI string

func init() {
	var buf bytes.Buffer
	stdout.SetOutput(&buf)
	stderr.SetOutput(&buf)

	os.Setenv("RANDOM_DELAY", "1")
}

func TestMain(m *testing.M) {
	closer, port, err := SetupTestServer(&TestServerPings)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	TestServerPort = port
	TestServerBaseURI = "http://127.0.0.1:" + strconv.Itoa(port)

	os.Exit(m.Run())
}

// SetupTestServerTest runs an HTTP serveur for testing with the following routes:
// - /error/:code to return the HTTP error given by :code
// - /panic to call panic()
// - anything else to append the RequestURI to the given pings slice
func SetupTestServer(pings *[]string) (srvCloser io.Closer, port int, err error) {
	http.HandleFunc("/error/", func(w http.ResponseWriter, req *http.Request) {
		status, _ := strconv.Atoi(filepath.Base(req.RequestURI))

		if status >= 300 && status < 400 {
			w.Header().Add("Location", fmt.Sprintf("http://127.0.0.1:%d/panic", port))
		}

		http.Error(w, req.RequestURI, status)
	})
	http.HandleFunc("/panic", func(w http.ResponseWriter, req *http.Request) {
		panic("This should be unreachable: " + req.RequestURI)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		*pings = append(*pings, req.RequestURI)
	})

	srvCloser, port, err = listenAndServeWithClose("127.0.0.1:0", nil)
	return
}

// Adapted from https://stackoverflow.com/a/40041517
func listenAndServeWithClose(addr string, handler http.Handler) (srvCloser io.Closer, port int, err error) {
	srv := &http.Server{Addr: addr, Handler: handler}
	if addr == "" {
		addr = ":http"
	}

	//var listener net.Listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}

	go func() {
		err := srv.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
		if err != nil {
			log.Println("HTTP Server Error - ", err)
		}
	}()

	srvCloser = listener
	port = listener.Addr().(*net.TCPAddr).Port
	return
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}
