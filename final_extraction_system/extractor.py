#!/usr/bin/env python3
"""
Smart Table Extraction System
1. pdfplumber detects pages with tables
2. Extract those pages to separate PDFs
3. Reducto pipeline processes each PDF
4. Save complete Reducto response with all metadata
"""
import os
from pathlib import Path
from typing import Dict, List, Any, Optional, Set
import logging
from dotenv import load_dotenv
import pdfplumber
import json
from datetime import datetime
from reducto import Reducto
import PyPDF2

load_dotenv()

# Setup logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Configuration
REDUCTO_API_KEY = os.getenv("REDUCTO_API_KEY")
REDUCTO_ENABLED = os.getenv("REDUCTO_ENABLED", "false").lower() == "true"
OUTPUT_DIR = "./output"
TEMP_PAGES_DIR = "./temp_pages"
PDF_DIR = "./input_pdfs"
os.makedirs(f"{OUTPUT_DIR}/reducto", exist_ok=True)
os.makedirs(TEMP_PAGES_DIR, exist_ok=True)


# ============================================
# STEP 1: TABLE DETECTION with pdfplumber
# ============================================

def detect_table_pages_with_pdfplumber(pdf_path: str) -> Set[int]:
    """
    Detect which pages contain tables using pdfplumber (FREE & FAST)
    """
    logger.info("=" * 60)
    logger.info("🔍 STEP 1: Detecting pages with tables (pdfplumber)")
    logger.info("=" * 60)
    
    table_pages = set()
    
    try:
        with pdfplumber.open(pdf_path) as pdf:
            total_pages = len(pdf.pages)
            logger.info(f"📄 Scanning {total_pages} pages...")
            
            for page_num, page in enumerate(pdf.pages, start=1):
                tables = page.find_tables()
                
                if tables:
                    table_pages.add(page_num)
                    logger.info(f"  ✓ Page {page_num}: Found {len(tables)} table(s)")
        
    except Exception as e:
        logger.error(f"❌ Table detection failed: {e}")
        return set()
    
    if table_pages:
        logger.info(f"✅ Detected tables on {len(table_pages)} pages: {sorted(table_pages)}")
    else:
        logger.warning("⚠️ No tables detected in PDF")
    
    logger.info("=" * 60)
    return table_pages


# ============================================
# STEP 2: EXTRACT PAGES TO SEPARATE PDFs
# ============================================

def extract_pages_to_folder(pdf_path: str, page_numbers: Set[int], output_folder: str = None) -> Dict[int, str]:
    """
    Extract specific pages into separate PDF files
    
    Args:
        pdf_path: Path to source PDF
        page_numbers: Set of page numbers to extract (1-indexed)
        output_folder: Folder to save extracted PDFs
    
    Returns:
        Dictionary mapping page_number -> pdf_path
    """
    if output_folder is None:
        output_folder = TEMP_PAGES_DIR
        
    logger.info("=" * 60)
    logger.info("📄 STEP 2: Extracting pages to separate PDFs")
    logger.info("=" * 60)
    
    page_pdfs = {}
    pdf_name = Path(pdf_path).stem
    
    try:
        reader = PyPDF2.PdfReader(pdf_path)
        
        for page_num in sorted(page_numbers):
            writer = PyPDF2.PdfWriter()
            writer.add_page(reader.pages[page_num - 1])  # 0-indexed
            
            page_pdf_path = f"{output_folder}/{pdf_name}_page_{page_num}.pdf"
            with open(page_pdf_path, 'wb') as f:
                writer.write(f)
            
            page_pdfs[page_num] = page_pdf_path
            logger.info(f"  ✓ Extracted page {page_num} → {Path(page_pdf_path).name}")
        
        logger.info(f"✅ Extracted {len(page_pdfs)} pages to {output_folder}")
        
    except Exception as e:
        logger.error(f"❌ Page extraction failed: {e}")
        return {}
    
    logger.info("=" * 60)
    return page_pdfs


# ============================================
# STEP 3: REDUCTO PIPELINE EXTRACTION
# ============================================

class ReductoExtractor:
    """
    Reducto SDK extractor using pipeline
    """
    
    def __init__(self, api_key: str = None):
        self.api_key = api_key or REDUCTO_API_KEY
        
        if not self.api_key:
            raise ValueError("REDUCTO_API_KEY not configured")
        
        self.client = Reducto(api_key=self.api_key)
    
    def extract_table_from_pdf_page(self, page_pdf_path: str, page_number: int) -> Optional[Dict[str, Any]]:
        """
        Process a single-page PDF with Reducto Pipeline
        
        Args:
            page_pdf_path: Path to single-page PDF
            page_number: Original page number from source PDF
        
        Returns:
            Complete Reducto response as dict
        """
        logger.info(f"🚀 Processing page {page_number} using pipeline...")
        
        try:
            # Upload the single-page PDF
            upload_result = self.client.upload(file=Path(page_pdf_path))
            
            # Use Reducto pipeline (no page_range needed - it's one page)
            result = self.client.pipeline.run(
                input=upload_result,
                pipeline_id="k9716dd7e2nwycd7q7wrpmpv3h7sy3a1"
            )
            
            # Convert result object to dict
            result_dict = self._convert_result_to_dict(result)
            
            # Add page metadata
            result_dict['page_number'] = page_number
            
            logger.info(f"✅ Successfully processed page {page_number}")
            return result_dict
            
        except Exception as e:
            logger.error(f"❌ Failed to process page {page_number}: {e}")
            return None
    
    def _convert_result_to_dict(self, obj) -> Any:
        """Recursively convert Reducto objects to dictionaries"""
        if obj is None:
            return None
        
        # Handle objects with __dict__
        if hasattr(obj, '__dict__'):
            result = {}
            for key, value in obj.__dict__.items():
                if not key.startswith('_'):  # Skip private attributes
                    result[key] = self._convert_result_to_dict(value)
            return result
        
        # Handle lists
        elif isinstance(obj, list):
            return [self._convert_result_to_dict(item) for item in obj]
        
        # Handle dictionaries
        elif isinstance(obj, dict):
            return {key: self._convert_result_to_dict(value) for key, value in obj.items()}
        
        # Return primitives as-is
        else:
            return obj
    
    def _count_tables(self, result_dict: Dict) -> int:
        """Count number of tables in result"""
        count = 0
        try:
            if result_dict and 'result' in result_dict:
                result_obj = result_dict['result']
                if 'chunks' in result_obj:
                    for chunk in result_obj['chunks']:
                        if 'blocks' in chunk:
                            for block in chunk['blocks']:
                                if block.get('type') == 'Table':
                                    count += 1
        except Exception as e:
            logger.warning(f"Could not count tables: {e}")
        return count


# ============================================
# MAIN EXTRACTION FUNCTION
# ============================================

def extract_tables_smart(pdf_path: str, save_output: bool = True) -> List[Dict[str, Any]]:
    """
    Smart extraction pipeline:
    1. pdfplumber detects table pages
    2. Extract pages to separate PDFs
    3. Reducto pipeline processes each PDF
    4. Cleanup temp files
    5. Save complete Reducto responses
    """
    logger.info("🎯 SMART TABLE EXTRACTION PIPELINE")
    logger.info(f"📄 PDF: {Path(pdf_path).name}")
    logger.info("")
    
    # STEP 1: Detect table pages
    table_pages = detect_table_pages_with_pdfplumber(pdf_path)
    
    if not table_pages:
        logger.warning("❌ No tables detected!")
        return []
    
    # Check Reducto configuration
    if not REDUCTO_ENABLED or not REDUCTO_API_KEY:
        logger.warning("⚠️ Reducto not enabled or API key missing")
        logger.info("💡 Set REDUCTO_ENABLED=true and REDUCTO_API_KEY in .env")
        return []
    
    # STEP 2: Extract pages to separate PDFs
    page_pdfs = extract_pages_to_folder(pdf_path, table_pages)
    
    if not page_pdfs:
        logger.error("❌ Failed to extract pages!")
        return []
    
    # STEP 3: Process each PDF with Reducto Pipeline
    logger.info("=" * 60)
    logger.info("🚀 STEP 3: Processing with Reducto Pipeline")
    logger.info("=" * 60)

    page_results = []
    extractor = ReductoExtractor()

    for page_num, page_pdf_path in page_pdfs.items():
        logger.info(f"📄 Processing page {page_num}...")
        
        try:
            result = extractor.extract_table_from_pdf_page(page_pdf_path, page_num)
            if result:
                page_results.append(result)
        except Exception as e:
            logger.error(f"  ❌ Error: {e}")
            continue

    # STEP 4: Cleanup temporary PDFs
    logger.info("=" * 60)
    logger.info("🗑️  STEP 4: Cleaning up temporary files...")
    logger.info("=" * 60)
    
    for page_pdf_path in page_pdfs.values():
        try:
            if os.path.exists(page_pdf_path):
                os.remove(page_pdf_path)
                logger.info(f"  ✓ Deleted: {Path(page_pdf_path).name}")
        except Exception as e:
            logger.warning(f"  ⚠️ Could not delete {page_pdf_path}: {e}")

    # Summary
    total_tables = sum(extractor._count_tables(r) for r in page_results)
    
    logger.info("=" * 60)
    logger.info(f"✅ EXTRACTION COMPLETE!")
    logger.info(f"   🔍 Pages detected: {len(table_pages)}")
    logger.info(f"   📋 Pages extracted: {len(page_results)}")
    logger.info(f"   📊 Total tables found: {total_tables}")
    logger.info("=" * 60)
    
    # Save output
    if save_output and page_results:
        save_results_to_json(page_results, pdf_path)
    
    return page_results


def save_results_to_json(results: List[Dict], pdf_path: str):
    """Save complete Reducto results to JSON"""
    try:
        pdf_name = Path(pdf_path).stem
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        
        output_data = {
            "filename": Path(pdf_path).name,
            "extraction_date": datetime.now().isoformat(),
            "extraction_method": "pdfplumber_detection + reducto_pipeline_extraction",
            "total_pages_processed": len(results),
            "pages": results  # Array of complete Reducto responses, one per page
        }
        
        json_file = f"{OUTPUT_DIR}/reducto/{pdf_name}_reducto_{timestamp}.json"
        with open(json_file, 'w') as f:
            json.dump(output_data, f, indent=2)
        
        logger.info(f"💾 Saved: {json_file}")
        
    except Exception as e:
        logger.error(f"Failed to save: {e}")


# ============================================
# MAIN
# ============================================

if __name__ == "__main__":
    import sys
    
    if len(sys.argv) < 2:
        print("Usage: python extractor.py <pdf_file>")
        sys.exit(1)
    
    # Get just the filename
    pdf_filename = Path(sys.argv[1]).name
    
    # Use PDF_DIR constant
    pdf_path = Path(PDF_DIR) / pdf_filename
    
    if not pdf_path.exists():
        logger.error(f"❌ Specified PDF not found: {pdf_path}")
        sys.exit(1)
    
    logger.info(f"📁 Processing specific PDF: {pdf_path.name}")
    
    results = extract_tables_smart(str(pdf_path))
    
    print(f"\n✅ Extracted data from {len(results)} pages!")