GOFLAGS ?= $(GOFLAGS:)

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
ifeq ($(BRANCH),$(filter $(BRANCH), master HEAD))
  VERSION := $(shell git describe --tags)
else
  DT := $(shell date '+%F::%T')
  VERSION := $(BRANCH)-$(DT)
endif

GOFLAGS = -ldflags '-X main.version $(VERSION)'

default: all

all: test install

install: get-deps
	@go build $(GOFLAGS) *.go

test: install
	@go test $(GOFLAGS) ./...

get-deps:
	@go get ./...

clean:
	@go clean $(GOFLAGS) -i ./
