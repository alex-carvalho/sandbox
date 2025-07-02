#!/bin/bash

mkdir bin
echo "Compiling network_examples.c and running it"
gcc network_examples.c -o bin/network_examples
./bin/network_examples

echo ""
echo "Running Python network_examples.py"
python3 network_examples.py 