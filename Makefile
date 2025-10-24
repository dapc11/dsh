BINARY_NAME=dsh

build:
	go build -o $(BINARY_NAME) .

clean:
	go clean
	rm -f $(BINARY_NAME)

test:
	go test ./...

fmt:
	gofmt -w .

lint:
	golangci-lint run

install: build
	cp $(BINARY_NAME) /usr/local/bin/

.PHONY: build clean test fmt lint install
