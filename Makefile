all: munin-http-timing

.PHONY: munin-http-timing
munin-http-timing:
	go build -ldflags '-s -w'
