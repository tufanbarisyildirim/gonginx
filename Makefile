PWD := $(shell pwd)
export GO111MODULE=on

test:
	go test -race ${PWD}/parser
	go test -race ${PWD}/config

fmt:
	find . -name "*.go" | xargs gofmt -w -s

deps:
	go get -v all

.PHONY: fmt deps test