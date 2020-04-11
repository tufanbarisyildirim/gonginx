PWD := $(shell pwd)
export GO111MODULE=on

test:
	go test -race -cover ${PWD}/parser/token
	go test -race -cover ${PWD}/parser
	go test -race -cover ${PWD}/config
	

fmt:
	find . -name "*.go" | xargs gofmt -w -s

deps:
	go get -v all

.PHONY: fmt deps test