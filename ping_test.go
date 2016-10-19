package main

import (
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestPing(t *testing.T) {
	TestServerPings = make([]string, 0)
	err := DoPing(map[string]string{
		"test1": TestServerBaseURI + "/test1",
		"test2": TestServerBaseURI + "/test2",
		"test3": "http://localhost:" + strconv.Itoa(TestServerPort) + "/test3",
	})

	if err != nil {
		t.Error(err)
	}

	sort.Strings(TestServerPings)
	expected := []string{"/test1", "/test2", "/test3"}
	if !reflect.DeepEqual(TestServerPings, expected) {
		t.Errorf("DoPing did not request the server, got %v expected %v.", TestServerPings, expected)
	}
}

func TestHTTPErrors(t *testing.T) {
	err := DoPing(map[string]string{
		"err301": TestServerBaseURI + "/error/301",
		"err302": TestServerBaseURI + "/error/302",
		"err318": TestServerBaseURI + "/error/318",
		"err404": TestServerBaseURI + "/error/404",
		"err418": TestServerBaseURI + "/error/418",
		"err500": TestServerBaseURI + "/error/500",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestEmptyURIList(t *testing.T) {
	err := DoPing(map[string]string{})
	if err == nil {
		t.Error("Ping should error out when given no URIs.")
	}
}
