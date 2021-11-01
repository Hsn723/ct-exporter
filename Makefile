VERSION := $(shell cat VERSION)
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"
CST_VERSION = 1.10.0

all: build

.PHONY: clean
clean:
	if [ -f ct-exporter ]; then rm ct-exporter; fi

.PHONY: lint
lint: clean
	if [ -z "$(shell which pre-commit)" ]; then pip3 install pre-commit; fi
	pre-commit install
	pre-commit run --all-files

.PHONY: setup-container-structure-test
setup-container-structure-test:
	if [ -z "$(shell which container-structure-test)" ]; then \
		curl -LO https://storage.googleapis.com/container-structure-test/v$(CST_VERSION)/container-structure-test-linux-amd64 && mv container-structure-test-linux-amd64 container-structure-test && chmod +x container-structure-test && sudo mv container-structure-test /usr/local/bin/; \
	fi

.PHONY: container-structure-test
container-structure-test: setup-container-structure-test
	printf "amd64\narm64" | xargs -n1 -I {} container-structure-test test --image ghcr.io/hsn723/ct-exporter:$(shell git describe --tags --abbrev=0)-next-{} --config cst.yaml

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
