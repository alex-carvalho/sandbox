#!/bin/bash

mkdir bin
echo "Compiling process_management_examples.c and running it"
gcc process_management_examples.c -o bin/process_management_examples
./bin/process_management_examples

echo ""
echo "Running Python process_management_examples.py"
python3 process_management_examples.py