#? all: Build
all: build
	@echo done
.PHONY: all

#? build: Build cl
build:
	@mkdir -p build
	go build -ldflags "-X main.GitCommit=$(shell git rev-parse HEAD)" -o build/hac ./cmd/hac
.PHONY: build

build-mock:
	@mkdir -p build
	go build -ldflags "-X main.GitCommit=$(shell git rev-parse HEAD)" -tags mock -o build/hac-mock ./cmd/hac

#? clean: Clean build
clean:
	rm -rf build
