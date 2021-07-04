VERSION := $(shell cat VERSION)
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"

all: build

.PHONY: clean
clean:
	if [ -f ct-exporter ]; then rm ct-exporter; fi

.PHONY: lint
lint: clean
	if [ -z "$(shell which pre-commit)" ]; then pip3 install pre-commit; fi
	pre-commit install
	pre-commit run --all-files

.PHONY: test
test: clean build
	go test -race -v ./...

.PHONY: verify
verify:
	go mod download
	go mod verify

.PHONY: build
build: clean
	env CGO_ENABLED=0 go build $(LDFLAGS) .
