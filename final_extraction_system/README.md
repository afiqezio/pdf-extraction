# Final Extraction System

A clean, organized system for extracting tables from PDFs using Camelot and producing high-quality Markdown output.

## Features

- **Clean Table Extraction**: Uses Camelot Stream method for reliable table detection
- **Quality Filtering**: Filters out low-quality and duplicate tables
- **Professional Output**: Produces clean, readable Markdown files
- **Easy Setup**: Simple installation and usage

## Quick Start

### 1. Setup
```bash
# Make setup script executable
chmod +x setup.sh

# Run setup
./setup.sh
```

### 2. Usage
```bash
# Activate virtual environment
source extraction_env/bin/activate

# Place PDF files in input_pdfs folder
cp your_file.pdf input_pdfs/

# Run extraction
python better_markdown_extractor.py

# Check results
ls output/markdown/
```

## File Structure

```
final_extraction_system/
├── better_markdown_extractor.py  # Main extraction script
├── requirements.txt              # Python dependencies
├── setup.sh                     # Setup script
├── README.md                    # This file
├── input_pdfs/                  # Place your PDF files here
└── output/
    └── markdown/                # Extracted Markdown files
```

## Output Format

The system produces clean Markdown files with:

- **Document metadata** (extraction date, file size, etc.)
- **Extraction summary** (number of tables found)
- **Quality tables** with proper formatting
- **Table metadata** (method, confidence, dimensions)

## Requirements

- Python 3.7+
- Camelot-py with OpenCV support
- Tabula-py for Java-based extraction
- Pandas for data processing

## Troubleshooting

### Common Issues

1. **Camelot installation issues**: Make sure you have OpenCV and Ghostscript installed
2. **Java issues**: Tabula requires Java 8+ to be installed
3. **Permission errors**: Make sure the setup script is executable (`chmod +x setup.sh`)

### Dependencies

- **Camelot-py**: For vector-based table extraction
- **Tabula-py**: For Java-based table extraction (fallback)
- **Pandas**: For data manipulation
- **Tqdm**: For progress bars

## Example Output

```markdown
# Your Document Name

**Extraction Date:** 2025-10-08T02:25:50.286139
**File Size:** 2,936,025 bytes
**Pages Processed:** all
**Extractors Used:** camelot

## Extraction Summary

- **Total Tables:** 11
- **Processing Time:** N/A

## Extracted Tables

### Table 1 (Page 11)

**Extraction Method:** camelot_stream  
**Confidence Score:** 93.58%  
**Dimensions:** 13 rows × 7 columns

| BotDepth_m | Lithologic Summary | % Sand & Silt | % Meta | % Shale | % Others |
| --- | --- | --- | --- | --- | --- |
| 2830-2840 | Claystone, Sandstone | 4 | 0 | 89 | 7 |
| 2850-2860 | Claystone, Sandstone | 27 | 0 | 73 | 0 |
| 2880-2890 | Sandstone | 100 | 0 | 0 |  |
```

## Support

For issues or questions, check the logs in the terminal output. The system provides detailed logging for troubleshooting.
