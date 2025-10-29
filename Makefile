BINARY_NAME=dsh

.PHONY: build test unit-test lint-test ui-test test-integration test-rendering coverage lint fmt clean check deps install

build:
	go build -o $(BINARY_NAME) .

clean:
	go clean
	rm -f $(BINARY_NAME) coverage.out coverage.html

# Run all tests (without linting for now)
test-no-lint: unit-test ui-test test-integration test-rendering

# Run unit tests only
unit-test:
	gotestsum --format testname -v -- -race ./internal/...

# Run linting tests
lint-test:
	golangci-lint run

# Run UI tests
ui-test:
	gotestsum --format testname -v ./test/ui/...

# Run integration tests (requires built binary)
test-integration: build
	gotestsum --format testname -v ./test/integration/...

# Run rendering tests with diagnostic output
test-rendering: build
	gotestsum --format testname -v ./test/rendering/...

# Run all tests
test:
	gotestsum --format testname ./...

# Legacy aliases for backward compatibility
test-unit: unit-test

# Run tests with coverage report
coverage:
	./test_all.sh

fmt:
	gofmt -w .

lint:
	golangci-lint run

# Run all quality checks
check: fmt lint-test unit-test

# Install dependencies
deps:
	go mod tidy
	go mod download

install: build
	cp $(BINARY_NAME) /usr/local/bin/
