GOFLAGS ?= $(GOFLAGS:)

default: all

all: install test

install: get-deps
	@go build $(GOFLAGS) ./

test: install
	@go test $(GOFLAGS) ./

get-deps:
	@go get github.com/stretchr/testify
	@go get github.com/op/go-logging
	@go get gopkg.in/yaml.v2
	@go get github.com/awslabs/aws-sdk-go
	@go get github.com/go-martini/martini
	@go get github.com/martini-contrib/encoder

clean:
	@go clean $(GOFLAGS) -i ./
