all:

.PHONY: cover
cover:
	go test -coverprofile=.coverage
	go tool cover -html=.coverage
