package munin

import (
	"testing"

	"github.com/DigitalBackstage/munin-http-timing/config"
)

func TestEmptyURIList(t *testing.T) {
	config := config.Config{
		URIs: map[string]string{},
	}
	out, err := DoPing(config)

	if err == nil {
		t.Error("Ping should error out when given no URIs.")
	}
	if out != "" {
		t.Error("Should get empty response from DoPing.")
	}
}
