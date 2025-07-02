#!/bin/bash

mkdir bin
echo "Compiling memory_management_examples.c and running it"
gcc memory_management_examples.c -o bin/memory_management_examples
./bin/memory_management_examples

echo ""
echo "Running Python memory_management_examples.py"
python3 memory_management_examples.py