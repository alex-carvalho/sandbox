#!/bin/bash

mkdir -p bin
echo "Compiling inter_process_communication_examples.c and running it"
gcc inter_process_communication_examples.c -o bin/inter_process_communication_examples
./bin/inter_process_communication_examples

echo ""
echo "Running Python inter_process_communication_examples.py"
python3 inter_process_communication_examples.py