#!/bin/bash
set -e

echo "🚀 Starting Flask app locally..."

VENV_DIR="venv"

if [ ! -d "$VENV_DIR" ]; then
    echo "📦 Creating virtual environment..."
    python3 -m venv "$VENV_DIR"
fi

echo "✅ Activating virtual environment..."
source "$VENV_DIR/bin/activate"

echo "📥 Installing dependencies..."
pip install -r requirements.txt

echo "🎯 Starting Flask application..."
python app.py
