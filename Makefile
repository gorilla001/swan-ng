.PHONY: default all prepare build golint build-no-golint product integration-test clean

PWD := $(shell pwd)
BIN := $(shell pwd)/bin
BUILD_IMG := "swan-build:latest"
PRODUCT_IMG := "swan:latest"

default: build

all: integration-test push

prepare:
	mkdir -p $(BIN)
	docker build --force-rm -t $(BUILD_IMG) -f Dockerfile.build .

golint: prepare
	docker run --rm -v $(PWD):/src:ro -e GOLINT_ONLY=yes $(BUILD_IMG)

build: prepare
	docker run --rm -v $(PWD):/src:ro -v $(BIN):/product:rw  $(BUILD_IMG)

product: build
	docker build --force-rm -t $(PRODUCT_IMG) -f Dockerfile.product .

integration-test: build
	echo "not implement yet"; exit 1

clean:
	rm -rfv $(BIN)

