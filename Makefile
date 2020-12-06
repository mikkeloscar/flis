# The MIT License (MIT)

# Copyright (c) 2016, 2017 Mikkel Oscar Lyderik Larsen

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

.PHONY: clean deps test check build init

BINARY        ?= flis
VERSION       ?= $(shell git describe --tags --always --dirty)
SOURCES       = $(shell find . -name '*.go')
GOPKGS        = $(shell go list ./... | grep -v /vendor/)
BUILD_FLAGS   ?= -v
LDFLAGS       ?= -X main.version=$(VERSION) -w -s

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

MAKE_COLOR=\033[33;01m%-20s\033[0m

default: build

.PHONY: help
help:
	@echo -e "$(OK_COLOR)==== $(BINARY) [$(VERSION)] ====$(NO_COLOR)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(MAKE_COLOR) : %s\n", $$1, $$2}'

init: ## Install requirements
	@echo -e "$(OK_COLOR)[$(BINARY)] Install requirements$(NO_COLOR)"
	@go get -u github.com/golang/glog
	@go get -u github.com/kardianos/govendor
	@go get -u github.com/golang/lint/golint
	@go get -u github.com/kisielk/errcheck

deps:
	go get -v -u -t ./...

clean: ## Cleanup
	rm -rf build

test: ## Launch unit tests
	go test -v $(GOPKGS)

check: ## Check go code
	golint ./...
	go vet -v $(GOPKGS)

build: build/$(BINARY) ## Build binary

build/$(BINARY): $(SOURCES)
	go build -o build/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .
