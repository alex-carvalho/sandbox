# Terraform Testing Example with tftest

This project demonstrates how to use the Python `tftest` library to test Terraform configurations. The example creates an S3 bucket with versioning and encryption enabled.

## Prerequisites

- Python 3.x
- Terraform
- AWS credentials configured

## Setup

1. Create and activate a virtual environment:
```bash
python -m venv .venv
source .venv/bin/activate  # On Windows use: .venv\Scripts\activate
```

2. Install required packages:
```bash
pip install tftest pytest
```

## Running Tests

To run the tests, make sure you're in the virtual environment and execute:

```bash
pytest 
```