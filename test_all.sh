#!/bin/bash

# Test runner for DSH shell
set -e

echo "🧪 Running DSH Test Suite"
echo "========================="

# Run all tests with coverage
echo "📊 Running tests with coverage..."
gotestsum --format testname -- -race -coverprofile=coverage.out ./internal/...

# Generate coverage report
echo "📈 Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

# Show coverage summary
echo "📋 Coverage Summary:"
go tool cover -func=coverage.out | tail -1

echo ""
echo "✅ All tests completed!"
echo "📄 Coverage report: coverage.html"
