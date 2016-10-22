all:
.PHONY: cover release clean test

clean:
	git clean -fdX

VERSION=$(shell git tag --points-at=HEAD):$(shell git rev-parse HEAD)
BUILD=go build -ldflags '-s -w -X main.version=$(VERSION)'
release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 $(BUILD) -o release/http-timing_amd64
	GOOS=linux GOARCH=arm GOARM=6 $(BUILD) -o release/http-timing_arm

cover:
	go test -coverprofile=.coverage
	go tool cover -html=.coverage

test:
	go test -v
	go test -race
