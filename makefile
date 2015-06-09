GOFLAGS ?= $(GOFLAGS:)

TAG := $(shell echo $$TRAVIS_TAG)
ifeq ($(TAG),)
  BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
  DT := $(shell date '+%F::%T')
  VERSION := $(BRANCH)-$(DT)
else
  VERSION := $(TAG)
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
