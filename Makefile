BINARY_NAME=dsh

.PHONY: build test coverage lint fmt clean check deps install

build:
	go build -o $(BINARY_NAME) .

clean:
	go clean
	rm -f $(BINARY_NAME) coverage.out coverage.html

test:
	go test -v -race ./internal/...

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
