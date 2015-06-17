GOFLAGS ?= $(GOFLAGS:)

TAG := $(VERSION)
ifeq ($(TAG),)
  BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
  DT := $(shell date '+%F::%T')
  VSN := $(BRANCH)-$(DT)
else
  VSN := $(TAG)
endif

GOFLAGS = -ldflags "-X main.version $(VSN)"

print-%  : ; @echo $* = $($*)

default: all

all: test install

install: get-deps
	@go build $(GOFLAGS) *.go

test:
	@go test $(GOFLAGS) ./...

get-deps:
	@go get ./...

clean:
	@go clean $(GOFLAGS) -i ./
