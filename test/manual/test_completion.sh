#!/bin/bash

# Manual testing script for DSH interactive features
# Run this to manually verify tab completion and rendering

set -e

echo "=== DSH Interactive Feature Testing ==="
echo

# Build DSH
echo "Building DSH..."
cd "$(dirname "$0")/../.."
make build
echo "✓ Build complete"
echo

# Create test environment
TEST_DIR=$(mktemp -d)
cd "$TEST_DIR"
echo "Test directory: $TEST_DIR"

# Create test files for completion
echo "Creating test files..."
touch test_file.txt another_file.go script.sh
mkdir test_dir
echo "✓ Test files created"
echo

echo "=== Manual Test Instructions ==="
echo
echo "Starting DSH. Please test the following:"
echo
echo "1. TAB COMPLETION:"
echo "   - Type 'ec' and press TAB → should complete to 'echo'"
echo "   - Type 'test_' and press TAB → should show file completions"
echo "   - Type 'cd test' and press TAB → should complete directory"
echo
echo "2. RENDERING:"
echo "   - Check prompt displays correctly: 'dsh> '"
echo "   - Type long command and verify cursor positioning"
echo "   - Use arrow keys to navigate command line"
echo "   - Press Ctrl+L to clear screen"
echo
echo "3. HISTORY:"
echo "   - Run some commands, then use UP/DOWN arrows"
echo "   - Verify history navigation works"
echo
echo "4. LINE EDITING:"
echo "   - Type command, use Ctrl+A/E for home/end"
echo "   - Use Ctrl+W to delete word"
echo "   - Use Ctrl+K to kill to end of line"
echo
echo "Type 'exit' when done testing."
echo
echo "Press ENTER to start DSH..."
read

# Start DSH
exec ./dsh
