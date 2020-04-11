PWD := $(shell pwd)
export GO111MODULE=on

test:
	go test -race ${PWD}/parser
	go test -race ${PWD}/config
	go test -race ${PWD}/parser/token

fmt:
	find . -name "*.go" | xargs gofmt -w -s

deps:
	go get -v all

.PHONY: fmt deps test