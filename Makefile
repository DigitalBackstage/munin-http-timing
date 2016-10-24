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

PACKAGES = $(shell find ./ -type d -not -path '*/\.*')
cover:
	echo 'mode: count' > .coverage-all
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=.coverage $(pkg);\
		tail -n +2 .coverage >> .coverage-all;\
	)
	go tool cover -html=.coverage-all

test:
	go test -v ./...
	go test -race ./...
