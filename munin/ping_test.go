package munin

import "testing"

func TestEmptyURIList(t *testing.T) {
	out, err := DoPing(map[string]string{})

	if err == nil {
		t.Error("Ping should error out when given no URIs.")
	}
	if out != "" {
		t.Error("Should get empty response from DoPing.")
	}
}
