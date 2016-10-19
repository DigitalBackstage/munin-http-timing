package main

import (
	"os"
	"reflect"
	"testing"
)

func assertDeepEqual(t *testing.T, expected, actual interface{}, msg string) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%s\ngot:      %v\nexpected: %v\n", msg, actual, expected)
	}
}

func TestURIsFromEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("TARGET_EXAMPLE1", "https://example.com/?1")
	os.Setenv("TARGET_EXAMPLE2", "https://example.com/?2")
	os.Setenv("TARGET_example3", "https://example.com/?3")

	actual := getURIsFromEnv()
	expected := map[string]string{
		"example1": "https://example.com/?1",
		"example2": "https://example.com/?2",
		"example3": "https://example.com/?3",
	}
	assertDeepEqual(t, expected, actual, "getURIsFromEnv properly parse env vars")
}

func TestBadURIsFromEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("TARGET_", "https://example.com/?noname")
	assertDeepEqual(t, map[string]string{}, getURIsFromEnv(), "blank names are not allowed")

	os.Clearenv()
	assertDeepEqual(t, map[string]string{}, getURIsFromEnv(), "no env means no URIs")

	os.Clearenv()
	os.Setenv("TARGET_BAD_URI", "utter nonsense")
	assertDeepEqual(t, map[string]string{}, getURIsFromEnv(), "bad URIs are not to be returned")

	os.Clearenv()
	os.Setenv("RANDOM_VAR", "https://example.com")
	assertDeepEqual(t, map[string]string{}, getURIsFromEnv(), "only use TARGET_ envs")
}
