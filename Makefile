all:
.PHONY: cover release clean

clean:
	git clean -fdX

BUILD=go build -ldflags '-s -w'
release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 $(BUILD) -o release/http-timing_amd64
	GOOS=linux GOARCH=arm GOARM=6 $(BUILD) -o release/http-timing_arm

cover:
	go test -v -coverprofile=.coverage
	go tool cover -html=.coverage
