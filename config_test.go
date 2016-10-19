package main

import "testing"

func TestConfig(t *testing.T) {
	err := DoConfig(map[string]string{"test": TestServerBaseURI + "/panic"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigWithoutURIs(t *testing.T) {
	err := DoConfig(map[string]string{})
	if err == nil {
		t.Error("DoConfig should fail when given no URIs.")
	}
}
