package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

var stderr = log.New(os.Stderr, "", 0)

// Filled during build
var version string

// Config holds the application configuration
type Config struct {
	URIs map[string]string

	RandomDelayEnabled bool
	ConfigAndPing      bool
	UserAgent          string
}

// NewConfigFromEnv creates and fills a Config from os.Environ()
func NewConfigFromEnv() Config {
	var config Config

	config.URIs = getURIsFromEnv(os.Environ())
	config.RandomDelayEnabled = os.Getenv("RANDOM_DELAY") == "1"
	config.UserAgent = os.Getenv("USER_AGENT")

	if len(config.UserAgent) == 0 {
		config.UserAgent = fmt.Sprintf("http-timing/%s", version)
	}

	// https://munin.readthedocs.io/en/latest/plugin/protocol-dirtyconfig.html#plugin-protocol-dirtyconfig
	config.ConfigAndPing = os.Getenv("MUNIN_CAP_DIRTYCONFIG") == "1"

	return config
}

// getURIsFromEnv returns a map associating names to urls from the process env vars
// Only vars prefixed with 'TARGET_' will be used, eg.
// TARGET_EXAMPLE=https://example.com/ will register the URI with "example"
// as the name.
func getURIsFromEnv(environ []string) map[string]string {
	uris := make(map[string]string, 0)

	for _, env := range environ {
		// Filter TARGET_*
		parts := strings.SplitN(env, "_", 2)
		if len(parts) != 2 || parts[0] != "TARGET" {
			continue
		}

		// Check for values
		name := strings.ToLower(strings.Split(parts[1], "=")[0])
		uri := strings.SplitN(env, "=", 2)[1]
		if len(name) <= 0 || len(uri) <= 0 {
			continue
		}

		// Check if URI is valid
		_, err := url.ParseRequestURI(uri)
		if err != nil {
			stderr.Printf("Invalid URI: %s\n", env)
			continue
		}

		uris[name] = uri
	}

	return uris
}
