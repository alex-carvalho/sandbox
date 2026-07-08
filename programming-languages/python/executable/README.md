# AWS IAM Roles Lister

A Python utility to list AWS IAM roles in a formatted table.

## Features

- Lists all IAM roles in your AWS account
- Displays role name, creation date, and path
- Formatted table output with total count
- Can be built as a standalone executable

## Requirements

- Python 3.x
- AWS credentials configured
- Required packages (see requirements.txt)

## Installation

```bash
pip install -r requirements.txt
```

## Usage

```bash
python list_roles.py
```

## Building Executable

```bash
pyinstaller list_roles.spec
```

The executable will be created in the `dist/` directory.

## AWS Configuration

Ensure your AWS credentials are configured via:
- AWS CLI (`aws configure`)
- Environment variables
- IAM roles (if running on EC2)