package pinger

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/DigitalBackstage/munin-http-timing/config"
)

func TestPing(t *testing.T) {
	TestServerPings.Purge()

	errs := doPingTest(map[string]string{
		"test1": TestServerBaseURI + "/test1",
		"test2": TestServerBaseURI + "/test2",
		"test3": "http://localhost:" + strconv.Itoa(TestServerPort) + "/test3",
	})

	if len(errs) != 0 {
		t.Error(errs)
	}

	expected := []string{"/test1", "/test2", "/test3"}
	actual := TestServerPings.Sorted()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("DoParallelPings did not request the server, got %v expected %v.", actual, expected)
	}
}

func TestHTTPErrors(t *testing.T) {
	errs := doPingTest(map[string]string{
		"err301": TestServerBaseURI + "/error/301",
		"err302": TestServerBaseURI + "/error/302",
		"err318": TestServerBaseURI + "/error/318",
		"err404": TestServerBaseURI + "/error/404",
		"err418": TestServerBaseURI + "/error/418",
		"err500": TestServerBaseURI + "/error/500",
	})

	if len(errs) != 6 {
		t.Error("Should have 6 errors.")
		t.Error(errs)
	}
}

// doPingTest pings a set of URIs and return the errors
func doPingTest(uris map[string]string) []error {
	queue := make(chan *RequestInfo, len(uris))
	config := config.Config{
		URIs:               uris,
		RandomDelayEnabled: false,
	}

	DoParallelPings(config, queue)
	var errs []error
	for i := 0; i < len(uris); i++ {
		info := <-queue

		if info.Error != nil {
			errs = append(errs, info.Error)
		}
	}

	return errs
}
