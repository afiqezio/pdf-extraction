#!/bin/bash
# Run script for Final Extraction System

echo "🚀 Starting Final Extraction System..."

# Check if virtual environment exists
if [ ! -d "extraction_env" ]; then
    echo "❌ Virtual environment not found. Please run setup.sh first."
    exit 1
fi

# Activate virtual environment
echo "🔧 Activating virtual environment..."
source extraction_env/bin/activate

# Check if PDF files exist
if [ ! -d "input_pdfs" ] || [ -z "$(ls -A input_pdfs 2>/dev/null)" ]; then
    echo "❌ No PDF files found in input_pdfs folder."
    echo "Please place your PDF files in the input_pdfs folder and try again."
    exit 1
fi

# Run extraction
echo "📊 Running table extraction..."
python better_markdown_extractor.py

echo "✅ Extraction complete!"
echo "📁 Check the output/markdown folder for results."
