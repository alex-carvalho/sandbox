#!/bin/bash

mkdir -p bin
echo "Compiling file_operations_examples.c and running it"
gcc file_operations_examples.c -o bin/file_operations_examples
./bin/file_operations_examples

echo ""
echo "Running Python file_operations_examples.py"
python3 file_operations_examples.py 