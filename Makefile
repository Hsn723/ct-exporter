BUILD_DIR=/tmp/ct-exporter/artifacts
VERSION := $(shell cat VERSION)
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"

all: build

.PHONY: clean
clean:
	rm -rf ${BUILD_DIR}

.PHONY: setup
setup:
	mkdir -p ${BUILD_DIR}

.PHONY: lint
lint: clean setup
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
build: clean setup
	env CGO_ENABLED=0 go build $(LDFLAGS) .
