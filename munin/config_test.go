package munin

import (
	"testing"

	"github.com/DigitalBackstage/munin-http-timing/config"
)

func TestConfigWithoutURIs(t *testing.T) {
	var config config.Config
	err := DoConfig(config)
	if err == nil {
		t.Error("DoConfig should fail when given no URIs.")
	}
}
