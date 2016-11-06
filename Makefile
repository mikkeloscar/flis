.PHONY: clean deps test check build

BINARY        ?= flise
VERSION       ?= $(shell git describe --tags --always --dirty)
SOURCES       = $(shell find . -name '*.go')
GOPKGS        = $(shell go list ./... | grep -v /vendor/)
BUILD_FLAGS   ?= -v
LDFLAGS       ?= -X main.version=$(VERSION) -w -s

default: build

deps:
	go get -v -u -t ./...

clean:
	rm -rf build

test:
	go test -v $(GOPKGS)

check:
	golint ./... | grep -v /vendor/
	go vet -v $(GOPKGS)

build: build/$(BINARY)

build/$(BINARY): $(SOURCES)
	go build -o build/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .
