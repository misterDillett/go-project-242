.PHONY: build run lint lint-fix test test-all test-verbose fmt clean tidy

build:
	go build -o bin/hexlet-path-size cmd/hexlet-path-size/main.go

run: build
	./bin/hexlet-path-size

test:
	go test -v ./...

test-verbose:
	go test -v -cover ./...

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

fmt:
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi

tidy:
	go mod tidy

clean:
	rm -rf bin/
	go clean -cache -testcache

test-all: test lint tidy