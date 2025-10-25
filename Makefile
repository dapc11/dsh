BINARY_NAME=dsh

.PHONY: build test test-unit test-integration test-rendering coverage lint fmt clean check deps install

build:
	go build -o $(BINARY_NAME) .

clean:
	go clean
	rm -f $(BINARY_NAME) coverage.out coverage.html

# Run all tests
test: test-unit test-integration test-rendering

# Run unit tests only
test-unit:
	go test -v -race ./internal/...

# Run integration tests (requires built binary)
test-integration: build
	go test -v ./test/integration/...

# Run rendering tests with diagnostic output
test-rendering: build
	go test -v ./test/rendering/...

# Run tests with coverage report
coverage:
	./test_all.sh

fmt:
	gofmt -w .

lint:
	golangci-lint run

# Run all quality checks
check: fmt lint test

# Install dependencies
deps:
	go mod tidy
	go mod download

install: build
	cp $(BINARY_NAME) /usr/local/bin/
