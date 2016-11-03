package munin

import "testing"

func TestConfigWithoutURIs(t *testing.T) {
	err := DoConfig(map[string]string{})
	if err == nil {
		t.Error("DoConfig should fail when given no URIs.")
	}
}
