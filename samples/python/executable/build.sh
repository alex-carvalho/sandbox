#!/bin/bash
set -e

python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
pyinstaller list_roles.spec
chmod +x dist/list_roles
echo "Binary created, executi this to run:\n ./dist/list_roles"
