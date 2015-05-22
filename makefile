GOFLAGS ?= $(GOFLAGS:)

default: all

all: install test

install:
	@go build $(GOFLAGS) ./

test: install
	@go test $(GOFLAGS) ./

clean:
	@go clean $(GOFLAGS) -i ./
