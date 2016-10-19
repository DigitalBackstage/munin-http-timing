package main

import "bytes"

func init() {
	var buf bytes.Buffer
	stdout.SetOutput(&buf)
	stderr.SetOutput(&buf)
}
