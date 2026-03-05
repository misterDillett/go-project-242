.PHONY: build run lint lint-fix test test-all fmt clean

build:
	go build -o bin/hexlet-path-size cmd/hexlet-path-size/main.go

run: build
	./bin/hexlet-path-size

test:
	go test -v ./...

test-all: test

test-size:
	@echo "=== Без флагов ==="
	./bin/hexlet-path-size project/
	@echo "\n=== Только -H ==="
	./bin/hexlet-path-size -H project/
	@echo "\n=== Только -a ==="
	./bin/hexlet-path-size -a project/
	@echo "\n=== -H и -a вместе ==="
	./bin/hexlet-path-size -H -a project/

fmt:
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

clean:
	rm -rf bin/
	go clean -cache -testcache