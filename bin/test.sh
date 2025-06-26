#!/bin/bash
# Run tests for all Go modules in the repository

# Find all directories containing a go.mod file
modules=$(find . -name "go.mod" -exec dirname {} \;)

# Loop through each module and run tests
for module in $modules; do
    echo "Running tests in $module..."
    (cd $module && go test ./... -v)
done