all:
.PHONY: cover release clean test

clean:
	git clean -fdX

VERSION=$(shell git tag --points-at=HEAD)
FULLVERSION=$(VERSION):$(shell git rev-parse HEAD)
BUILD=go build -ldflags '-s -w -X main.version=$(FULLVERSION) -X config.version=$(FULLVERSION)'
BUILDDIR=munin-http-timing-$(VERSION)
release:
	mkdir -p "$(BUILDDIR)"
	GOOS=linux GOARCH=amd64 $(BUILD) -o "$(BUILDDIR)/http-timing_amd64"
	GOOS=linux GOARCH=arm GOARM=6 $(BUILD) -o "$(BUILDDIR)/http-timing_arm"
	tar czf "munin-http-timing-$(VERSION).tgz" "$(BUILDDIR)"
	rm -rf "$(BUILDDIR)"

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
