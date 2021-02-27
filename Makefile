BUILD_DIR=/tmp/ct-exporter/artifacts
VERSION := $(shell cat VERSION)
LDFLAGS=-ldflags "-w -s -X github.com/Hsn723/ct-exporter/cmd.CurrentVersion=${VERSION}"
OS ?= linux
ARCH ?= amd64
ifeq ($(OS), windows)
EXT = .exe
endif

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
test: clean
	go test -race -v ./...

.PHONY: verify
verify:
	go mod download
	go mod verify

.PHONY: build
build: clean setup
	env GOOS=$(OS) GOARCH=$(ARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/ct-exporter-$(OS)-$(ARCH)$(EXT) .
