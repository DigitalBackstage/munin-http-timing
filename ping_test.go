package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"testing"
)
import "net/http"

func TestPing(t *testing.T) {
	pings := make([]string, 0)
	closer, port, err := setupServer(&pings)
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
	pings = make([]string, 0)

	err = DoPing(map[string]string{})
	if err == nil {
		t.Error("Ping should error out when no URIs provided.")
	}

	err = DoPing(map[string]string{
		"err500": baseURI + "/error/500",
		"err404": baseURI + "/error/404",
		"err302": baseURI + "/error/302",
	})
	if err != nil {
		t.Error(err)
	}
}

func setupServer(pings *[]string) (srvCloser io.Closer, port int, err error) {
	http.HandleFunc("/error/", func(w http.ResponseWriter, req *http.Request) {
		*pings = append(*pings, req.RequestURI)
		status, _ := strconv.Atoi(filepath.Base(req.RequestURI))

		if status >= 300 && status < 400 {
			w.Header().Add("Location", fmt.Sprintf("http://127.0.0.1:%d/panic", port))
		}

		http.Error(w, req.RequestURI, status)
	})
	http.HandleFunc("/panic", func(w http.ResponseWriter, req *http.Request) {
		panic("This should be unreachable.")
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
