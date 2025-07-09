#!/bin/bash

mkdir -p bin
echo "Compiling file_system_operations_examples.c and running it"
gcc file_system_operations_examples.c -o bin/file_system_operations_examples
./bin/file_system_operations_examples

echo ""
echo "Running Python file_system_operations_examples.py"
python3 file_system_operations_examples.py