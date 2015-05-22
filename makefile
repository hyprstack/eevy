GOFLAGS ?= $(GOFLAGS:)

default: all

all: install test

install: get-deps
	@go build $(GOFLAGS) ./

test: install
	@go test $(GOFLAGS) ./

get-deps:
	@go get ./

clean:
	@go clean $(GOFLAGS) -i ./
