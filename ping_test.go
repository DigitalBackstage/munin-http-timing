package main

import (
	"io"
	"log"
	"net"
	"reflect"
	"sort"
	"strconv"
	"testing"
)
import "net/http"

func TestPing(t *testing.T) {
	pings := make([]string, 0)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		pings = append(pings, req.RequestURI)
	})

	closer, port, err := listenAndServeWithClose("127.0.0.1:0", nil)
	baseURI := "http://127.0.0.1:" + strconv.Itoa(port)
	defer closer.Close()
	if err != nil {
		panic(err)
	}

	err = DoPing(map[string]string{
		"test1": baseURI + "/test1",
		"test2": baseURI + "/test2",
		"test3": "http://localhost:" + strconv.Itoa(port) + "/test3",
	})
	if err != nil {
		t.Error(err)
	}

	sort.Strings(pings)
	expected := []string{"/test1", "/test2", "/test3"}
	if !reflect.DeepEqual(pings, expected) {
		t.Errorf("DoPing did not request the server, got %v expected %v.", pings, expected)
	}

	err = DoPing(map[string]string{})
	if err == nil {
		t.Error("Ping should error out when no URIs provided.")
	}
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
