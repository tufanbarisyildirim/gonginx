PWD := $(shell pwd)
export GO111MODULE=on

test:
	go test -race -cover ${PWD}/parser/token
	go test -race -cover ${PWD}/parser
	go test -race -cover ${PWD}

lint:
	golangci-lint run ./...
	golint ./...

check:
	staticcheck ./...

example:
	go run ${PWD}/examples/$(example)

bench:
	go test -bench=. -benchmem ${PWD}/parser
	
fmt:
	find . -name "*.go" | xargs gofmt -w -s

deps:
	go get -v all

.PHONY: fmt deps test lint
