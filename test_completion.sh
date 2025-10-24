#!/bin/bash

echo "Testing DSH Tab Completion"
echo "=========================="
echo ""
echo "Try these commands in DSH:"
echo "1. Type 'l' and press Tab (should show ls, ln, etc.)"
echo "2. Type 'cd ' and press Tab (should show directories)"
echo "3. Type 'echo test' and press Tab (should show files)"
echo ""
echo "Starting DSH shell..."
echo ""

./dsh
