package pinger

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"testing"
)

// Requested URIs will be appended in order here
var TestServerPings *Pings
var TestServerPort int
var TestServerBaseURI string

type Pings struct {
	lock  sync.Mutex
	pings []string
}

func NewPings() *Pings {
	p := &Pings{}
	p.pings = make([]string, 0)
	return p
}

func (p *Pings) Purge() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.pings = []string{}
}

func (p *Pings) Sorted() []string {
	p.lock.Lock()
	defer p.lock.Unlock()

	sort.Strings(p.pings)
	return p.pings
}

func (p *Pings) Push(uri string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.pings = append(p.pings, uri)
}

func init() {
	os.Setenv("RANDOM_DELAY", "1")
}

func TestMain(m *testing.M) {
	TestServerPings = NewPings()
	closer, port, err := SetupTestServer(TestServerPings)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	TestServerPort = port
	TestServerBaseURI = "http://127.0.0.1:" + strconv.Itoa(port)

	os.Exit(m.Run())
}

// SetupTestServerTest runs an HTTP server for testing with the following routes:
// - /error/:code to return the HTTP error given by :code
// - /panic to call panic()
// - anything else to append the RequestURI to the given pings slice
func SetupTestServer(pings *Pings) (srvCloser io.Closer, port int, err error) {
	http.HandleFunc("/error/", func(w http.ResponseWriter, req *http.Request) {
		status, _ := strconv.Atoi(filepath.Base(req.RequestURI))

		if status >= 300 && status < 400 {
			// As we should never follow redirects, redirecting to /panic gives
			// us the test for free
			w.Header().Add("Location", fmt.Sprintf("http://127.0.0.1:%d/panic", port))
		}

		http.Error(w, req.RequestURI, status)
	})
	http.HandleFunc("/panic", func(w http.ResponseWriter, req *http.Request) {
		panic("This should be unreachable: " + req.RequestURI)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		pings.Push(req.RequestURI)
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
