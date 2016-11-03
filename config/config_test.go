package config

import (
	"io/ioutil"
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

	actual := getURIsFromEnv(os.Environ())
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
	assertDeepEqual(t, map[string]string{}, getURIsFromEnv(os.Environ()), "blank names are not allowed")

	os.Clearenv()
	assertDeepEqual(t, map[string]string{}, getURIsFromEnv(os.Environ()), "no env means no URIs")

	stderr.SetOutput(ioutil.Discard)
	os.Clearenv()
	os.Setenv("TARGET_BAD_URI", "utter nonsense")
	assertDeepEqual(t, map[string]string{}, getURIsFromEnv(os.Environ()), "bad URIs are not to be returned")
	stderr.SetOutput(os.Stderr)

	os.Clearenv()
	os.Setenv("RANDOM_VAR", "https://example.com")
	assertDeepEqual(t, map[string]string{}, getURIsFromEnv(os.Environ()), "only use TARGET_ envs")
}

func TestNewConfigFromEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("RANDOM_DELAY", "1")
	os.Setenv("MUNIN_CAP_DIRTYCONFIG", "1")

	config := NewConfigFromEnv()

	assertDeepEqual(t, map[string]string{}, config.URIs, "no target expected")
	if config.RandomDelayEnabled != true {
		t.Error("Expected random delay to be enabled.")
	}
	if config.ConfigAndPing != true {
		t.Error("Expected dirty config to be set.")
	}
}

func TestNewConfigFromEnvWithZeroes(t *testing.T) {
	os.Clearenv()
	os.Setenv("RANDOM_DELAY", "0")
	os.Setenv("MUNIN_CAP_DIRTYCONFIG", "0")

	config := NewConfigFromEnv()

	assertDeepEqual(t, map[string]string{}, config.URIs, "no target expected")
	if config.RandomDelayEnabled != false {
		t.Error("Expected random delay to be disabled.")
	}
	if config.ConfigAndPing != false {
		t.Error("Expected dirty config to be unset.")
	}
}
