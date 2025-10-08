#!/usr/bin/env python3
"""
Final Better Markdown Extractor - Clean, readable table output
Produces high-quality Markdown with properly formatted tables
"""
import camelot
import pandas as pd
import os
import json
from pathlib import Path
from typing import List, Dict, Any, Optional
import logging
from tqdm import tqdm
import hashlib
from datetime import datetime

# -----------------------------
# CONFIG
# -----------------------------
PDF_DIR = "/Users/husaini.rosdi/Documents/working_project/pdf-extraction/final_extraction_system/input_pdfs"
OUTPUT_DIR = "/Users/husaini.rosdi/Documents/working_project/pdf-extraction/final_extraction_system/output"

# Camelot Settings
camelot_config = {
    "LATTICE_THRESHOLD": 0.5,
    "STREAM_THRESHOLD": 0.5,
    "LINE_SCALE": 40,
}

# Quality Filtering
quality_config = {
    "MIN_TABLE_CONFIDENCE": 50,
    "MIN_TABLE_ROWS": 3,
    "MIN_TABLE_COLS": 2,
    "MAX_DUPLICATE_SIMILARITY": 0.8,
    "ENABLE_DUPLICATE_FILTERING": True,
}


# -----------------------------
# SETUP
# -----------------------------
os.makedirs(f"{OUTPUT_DIR}/markdown", exist_ok=True)

# Setup logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# -----------------------------
# BETTER CAMELOT EXTRACTOR
# -----------------------------
class BetterCamelotExtractor:
    """Better Camelot extractor that produces clean, readable output"""
    
    def __init__(self):
        self.lattice_threshold = camelot_config["LATTICE_THRESHOLD"]
        self.stream_threshold = camelot_config["STREAM_THRESHOLD"]
        self.line_scale = camelot_config["LINE_SCALE"]
    
    def extract_tables(self, pdf_path: str, pages: Optional[List[int]] = None, flavors: List[str] = None) -> List[Dict[str, Any]]:
        """Extract tables using Camelot with better filtering"""
        if flavors is None:
            flavors = ['stream']
        
        try:
            pdf_path = Path(pdf_path)
            if not pdf_path.exists():
                raise FileNotFoundError(f"PDF file not found: {pdf_path}")
            
            logger.info(f"Processing PDF with Better Camelot: {pdf_path}")
            
            extracted_tables = []
            pages_str = self._format_pages(pages)
            
            for flavor in flavors:
                try:
                    logger.info(f"Trying Camelot with {flavor} flavor")
                    tables = self._extract_with_flavor(pdf_path, pages_str, flavor)
                    
                    for table in tables:
                        table_data = self._process_camelot_table(table, flavor)
                        if table_data and self._is_quality_table(table_data):
                            extracted_tables.append(table_data)
                            
                except Exception as e:
                    logger.warning(f"Camelot {flavor} failed: {e}")
                    continue
            
            logger.info(f"Extracted {len(extracted_tables)} quality tables using Camelot")
            return extracted_tables
            
        except Exception as e:
            logger.error(f"Error extracting tables with Camelot: {e}")
            return []
    
    def _format_pages(self, pages: Optional[List[int]]) -> str:
        """Convert page list to Camelot format"""
        if pages is None:
            return "all"
        return ",".join(map(str, pages))
    
    def _extract_with_flavor(self, pdf_path: Path, pages_str: str, flavor: str):
        """Extract tables with specific flavor"""
        if flavor == 'lattice':
            return camelot.read_pdf(
                str(pdf_path),
                pages=pages_str,
                flavor='lattice',
                line_scale=self.line_scale
            )
        elif flavor == 'stream':
            return camelot.read_pdf(
                str(pdf_path),
                pages=pages_str,
                flavor='stream',
                edge_tol=500
            )
        else:
            raise ValueError(f"Unknown flavor: {flavor}")
    
    def _process_camelot_table(self, table, flavor: str) -> Optional[Dict[str, Any]]:
        """Process Camelot table with better cleaning"""
        try:
            # Get the dataframe
            df = table.df
            
            # Clean the dataframe
            df = self._clean_dataframe(df)
            
            if df.empty or len(df) < 2:
                return None
            
            # Convert to list of lists
            table_data = [df.columns.tolist()] + df.values.tolist()
            
            # Calculate confidence based on accuracy
            confidence = getattr(table, 'accuracy', 0.8)
            if confidence < 0.5:
                confidence = 0.5
            
            return {
                'source': 'camelot',
                'page': getattr(table, 'page', 1),
                'method': f'camelot_{flavor}',
                'data': table_data,
                'confidence': confidence,
                'metadata': {
                    'accuracy': getattr(table, 'accuracy', None),
                    'whitespace': getattr(table, 'whitespace', None),
                    'order': getattr(table, 'order', None),
                    'flavor': flavor,
                    'num_rows': len(df),
                    'num_cols': len(df.columns)
                }
            }
            
        except Exception as e:
            logger.error(f"Error processing Camelot table: {e}")
            return None
    
    def _clean_dataframe(self, df: pd.DataFrame) -> pd.DataFrame:
        """Clean dataframe with better logic"""
        # Remove completely empty rows
        df = df.dropna(how='all')
        
        # Remove completely empty columns
        df = df.dropna(axis=1, how='all')
        
        # Fill NaN values with empty strings
        df = df.fillna('')
        
        # Remove rows where all values are empty strings
        df = df[~(df == '').all(axis=1)]
        
        # Clean up column names
        df.columns = [str(col).strip() for col in df.columns]
        
        return df
    
    def _is_quality_table(self, table_data: Dict[str, Any]) -> bool:
        """Better quality criteria for tables"""
        try:
            data = table_data['data']
            metadata = table_data['metadata']
            
            # Must have at least 3 rows and 2 columns
            if metadata['num_rows'] < 3 or metadata['num_cols'] < 2:
                return False
            
            # Must have reasonable confidence
            if table_data['confidence'] < 0.5:
                return False
            
            # Must have substantial content (at least 40% non-empty cells)
            non_empty_cells = 0
            total_cells = 0
            
            for row in data:
                for cell in row:
                    total_cells += 1
                    if str(cell).strip() != '':
                        non_empty_cells += 1
            
            if total_cells == 0 or non_empty_cells / total_cells < 0.4:
                return False
            
            # Check if it's a real data table (not just headers or formatting)
            has_data_rows = False
            for row in data[1:]:  # Skip header
                non_empty_in_row = sum(1 for cell in row if str(cell).strip() != '')
                if non_empty_in_row >= 2:  # At least 2 non-empty cells per row
                    has_data_rows = True
                    break
            
            if not has_data_rows:
                return False
            
            return True
            
        except Exception as e:
            logger.error(f"Error checking table quality: {e}")
            return False
    

def remove_duplicates(tables: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
    """Remove duplicate tables based on content hash"""
    seen_hashes = set()
    unique_tables = []
    
    for table in tables:
        # Create hash of table content (excluding metadata)
        content_str = str(table['data'])
        content_hash = hashlib.md5(content_str.encode()).hexdigest()
        
        if content_hash not in seen_hashes:
            seen_hashes.add(content_hash)
            unique_tables.append(table)
            logger.info(f"  ✅ Keeping unique table: {table['metadata']['num_rows']}x{table['metadata']['num_cols']}")
        else:
            logger.info(f"  ❌ Removing duplicate table: {table['metadata']['num_rows']}x{table['metadata']['num_cols']}")
    
    logger.info(f"📊 Removed {len(tables) - len(unique_tables)} duplicates, kept {len(unique_tables)} unique tables")
    return unique_tables

def format_table_as_markdown(table_data: List[List[str]]) -> str:
    """Format table data as clean, readable Markdown table"""
    if not table_data:
        return "No data"
    
    # Clean up the data
    cleaned_data = []
    for row in table_data:
        cleaned_row = []
        for cell in row:
            # Clean up cell content
            cell_str = str(cell).strip()
            # Remove excessive whitespace and newlines
            cell_str = ' '.join(cell_str.split())
            cleaned_row.append(cell_str)
        cleaned_data.append(cleaned_row)
    
    # Create header
    header = "| " + " | ".join(cleaned_data[0]) + " |"
    separator = "| " + " | ".join(["---"] * len(cleaned_data[0])) + " |"
    
    # Create rows
    rows = []
    for row in cleaned_data[1:]:
        rows.append("| " + " | ".join(row) + " |")
    
    return "\n".join([header, separator] + rows)

def save_better_markdown(tables: List[Dict[str, Any]], pdf_name: str, pdf_path: Path):
    """Save tables as better formatted Markdown file"""
    try:
        # Get file stats
        file_size = pdf_path.stat().st_size
        extraction_date = datetime.now().isoformat()
        
        # Create better markdown content
        md_content = f"""# {pdf_name}

**Extraction Date:** {extraction_date}
**File Size:** {file_size:,} bytes
**Pages Processed:** all
**Extractors Used:** camelot

## Extraction Summary

- **Total Tables:** {len(tables)}
- **Processing Time:** N/A

## Document Content

This document contains extracted tables from the PDF file. Each table has been processed and cleaned for readability.

## Extracted Tables

"""
        
        # Add each table with better formatting
        for i, table in enumerate(tables):
            table_num = i + 1
            page = table['page']
            method = table['method']
            confidence = table['confidence']
            rows = table['metadata']['num_rows']
            cols = table['metadata']['num_cols']
            
            md_content += f"""### Table {table_num} (Page {page})

**Extraction Method:** {method}  
**Confidence Score:** {confidence:.2f}%  
**Dimensions:** {rows} rows × {cols} columns

{format_table_as_markdown(table['data'])}

---

"""
        
        # Save markdown file
        md_file = f"{OUTPUT_DIR}/markdown/{pdf_name}_extracted.md"
        with open(md_file, 'w', encoding='utf-8') as f:
            f.write(md_content)
        
        logger.info(f"💾 Saved better markdown: {md_file}")
        
    except Exception as e:
        logger.error(f"❌ Failed to save better markdown: {e}")

def main():
    """Main function"""
    logger.info("🚀 Starting Final Better Markdown Extractor...")
    logger.info("=" * 60)
    
    # Find PDF files
    pdf_files = list(Path(PDF_DIR).glob("*.pdf"))
    if not pdf_files:
        logger.error(f"❌ No PDF files found in {PDF_DIR}")
        return
    
    logger.info(f"📁 Found {len(pdf_files)} PDF files to process")
    
    # Initialize extractor
    extractor = BetterCamelotExtractor()
    
    # Process each PDF
    for pdf_path in pdf_files:
        try:
            logger.info(f"🚀 Processing PDF: {pdf_path.name}")
            
            # Extract quality tables
            tables = extractor.extract_tables(str(pdf_path))
            
            if not tables:
                logger.info(f"  No quality tables found in {pdf_path.name}")
                continue
            
            # Remove duplicates
            unique_tables = remove_duplicates(tables)
            
            logger.info(f"✅ {pdf_path.name} completed:")
            logger.info(f"   📊 Total tables found: {len(tables)}")
            logger.info(f"   🎯 Quality unique tables: {len(unique_tables)}")
            
            # Save as better markdown
            pdf_name = pdf_path.stem
            save_better_markdown(unique_tables, pdf_name, pdf_path)
            
        except Exception as e:
            logger.error(f"❌ Failed to process {pdf_path.name}: {e}")
    
    logger.info("🎉 Processing complete!")
    logger.info(f"📁 Check the following directory for extracted markdown:")
    logger.info(f"   📊 Markdown files: {OUTPUT_DIR}/markdown/")

if __name__ == "__main__":
    main()
