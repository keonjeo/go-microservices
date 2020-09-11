.PHONY: run test dep build
# NAME = $(notdir $(shell pwd))

run:
	go run *.go
dep:
	go mod download
build:
	go build
test:
	go test -count=1
