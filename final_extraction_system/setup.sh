#!/bin/bash
# Setup script for Final Extraction System

echo "🚀 Setting up Final Extraction System..."

# Create directories
mkdir -p input_pdfs
mkdir -p output/markdown

# Create virtual environment
echo "📦 Creating virtual environment..."
python3 -m venv extraction_env

# Activate virtual environment
echo "🔧 Activating virtual environment..."
source extraction_env/bin/activate

# Install requirements
echo "📥 Installing requirements..."
pip install --upgrade pip
pip install -r requirements.txt

echo "✅ Setup complete!"
echo ""
echo "To use the system:"
echo "1. Place your PDF files in the 'input_pdfs' folder"
echo "2. Run: source extraction_env/bin/activate"
echo "3. Run: python better_markdown_extractor.py"
echo "4. Check results in 'output/markdown' folder"
echo ""
echo "🎉 Ready to extract tables!"
