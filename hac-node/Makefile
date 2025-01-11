#? all: Build
all: build
	@echo done
.PHONY: all

#? build: Build cl
build:
	@mkdir -p build
	go build -ldflags "-X main.GitCommit=$(shell git rev-parse HEAD)" -o build ./cmd/hac
.PHONY: build

#? clean: Clean build
clean:
	rm -rf build
