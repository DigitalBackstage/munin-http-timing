LDFLAGS=-s -w
all: munin-http-timing

munin-http-timing:
	go build -ldflags '$(LDFLAGS)'

.PHONY: cover
cover:
	go test -coverprofile=.coverage
	go tool cover -html=.coverage
