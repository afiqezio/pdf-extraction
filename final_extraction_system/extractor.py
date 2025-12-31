#!/usr/bin/env python3
"""
Smart Table Extraction System
1. pdfplumber detects pages with tables
2. Convert those pages to PNG images (high-res)
3. Reducto pipeline processes each PNG
4. Save complete Reducto response (all metadata)
"""
import os
from pathlib import Path
from typing import Dict, List, Any, Optional, Set
import logging
from dotenv import load_dotenv
import pdfplumber
import fitz  # PyMuPDF for PDF to image conversion
import json
from datetime import datetime
from reducto import Reducto

load_dotenv()

# Logging setup
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Config
REDUCTO_API_KEY = os.getenv("REDUCTO_API_KEY")
REDUCTO_ENABLED = os.getenv("REDUCTO_ENABLED", "false").lower() == "true"
OUTPUT_DIR = "./output"
TEMP_IMG_DIR = "./temp_images"
PDF_DIR = "./input_pdfs"
os.makedirs(f"{OUTPUT_DIR}/reducto", exist_ok=True)
os.makedirs(TEMP_IMG_DIR, exist_ok=True)

def detect_table_pages_with_pdfplumber(pdf_path: str) -> Set[int]:
    """Detect which pages contain tables using pdfplumber."""
    logger.info("=" * 60)
    logger.info("🔍 STEP 1: Detecting pages with tables (pdfplumber)")
    logger.info("=" * 60)
    table_pages = set()
    try:
        with pdfplumber.open(pdf_path) as pdf:
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

def convert_pages_to_images(pdf_path: str, page_numbers: Set[int], dpi: int = 300) -> Dict[int, str]:
    """
    Convert specific PDF pages to high-quality PNG images.
    """
    logger.info("=" * 60)
    logger.info("🖼️  STEP 2: Converting pages to PNG images")
    logger.info("=" * 60)
    page_images = {}
    pdf_name = Path(pdf_path).stem
    try:
        pdf_document = fitz.open(pdf_path)
        for page_num in sorted(page_numbers):
            logger.info(f"📸 Converting page {page_num} at {dpi} DPI...")
            page = pdf_document[page_num - 1]
            zoom = dpi / 72
            mat = fitz.Matrix(zoom, zoom)
            pix = page.get_pixmap(matrix=mat)
            image_path = f"{TEMP_IMG_DIR}/{pdf_name}_page_{page_num}.png"
            pix.save(image_path)
            page_images[page_num] = image_path
            logger.info(f"  ✓ Saved: {Path(image_path).name} ({pix.width}x{pix.height})")
        pdf_document.close()
        logger.info(f"✅ Converted {len(page_images)} pages to images")
        logger.info("=" * 60)
        return page_images
    except Exception as e:
        logger.error(f"❌ Image conversion failed: {e}")
        return {}

class ReductoExtractor:
    """
    Reducto SDK extractor (pipeline mode)
    """
    def __init__(self, api_key: str = None):
        self.api_key = api_key or REDUCTO_API_KEY
        if not self.api_key:
            raise ValueError("REDUCTO_API_KEY not configured")
        self.client = Reducto(api_key=self.api_key)

    def extract_table_from_image(self, image_path: str, page_number: int) -> Optional[Dict[str, Any]]:
        """
        Process a PNG image with Reducto Pipeline.
        """
        logger.info(f"🚀 Processing page {page_number} from image: {image_path} ...")
        try:
            # Upload the PNG image
            upload_result = self.client.upload(file=Path(image_path))
            # Use Reducto pipeline
            result = self.client.pipeline.run(
                input=upload_result, pipeline_id="k9716dd7e2nwycd7q7wrpmpv3h7sy3a1"
            )
            result_dict = self._convert_result_to_dict(result)
            result_dict['page_number'] = page_number
            logger.info(f"✅ Successfully processed page {page_number}")
            return result_dict
        except Exception as e:
            logger.error(f"❌ Failed to process page {page_number}: {e}")
            return None

    def _convert_result_to_dict(self, obj) -> Any:
        if obj is None:
            return None
        if hasattr(obj, '__dict__'):
            result = {}
            for key, value in obj.__dict__.items():
                if not key.startswith('_'):
                    result[key] = self._convert_result_to_dict(value)
            return result
        elif isinstance(obj, list):
            return [self._convert_result_to_dict(item) for item in obj]
        elif isinstance(obj, dict):
            return {key: self._convert_result_to_dict(value) for key, value in obj.items()}
        else:
            return obj
    def _count_tables(self, result_dict: Dict) -> int:
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

def extract_tables_smart(pdf_path: str, save_output: bool = True) -> List[Dict[str, Any]]:
    """
    Smart extraction pipeline:
    1. pdfplumber detects table pages
    2. Convert pages to PNG images
    3. Reducto pipeline processes each image
    4. Cleanup temp files
    5. Save Reducto responses to JSON
    """
    logger.info("🎯 SMART TABLE EXTRACTION PIPELINE")
    logger.info(f"📄 PDF: {Path(pdf_path).name}\n")
    # STEP 1: Detect table pages
    table_pages = detect_table_pages_with_pdfplumber(pdf_path)
    if not table_pages:
        logger.warning("❌ No tables detected!")
        return []
    if not REDUCTO_ENABLED or not REDUCTO_API_KEY:
        logger.warning("⚠️ Reducto not enabled or API key missing")
        logger.info("💡 Set REDUCTO_ENABLED=true and REDUCTO_API_KEY in .env")
        return []
    # STEP 2: Convert pages to PNG images (high-res)
    page_images = convert_pages_to_images(pdf_path, table_pages, dpi=300)
    if not page_images:
        logger.error("❌ Failed to convert to images!")
        return []
    # STEP 3: Process each image
    logger.info("=" * 60)
    logger.info("🚀 STEP 3: Processing images with Reducto Pipeline")
    logger.info("=" * 60)
    page_results = []
    extractor = ReductoExtractor()
    for page_num, image_path in page_images.items():
        logger.info(f"📄 Processing page {page_num} ...")
        try:
            result = extractor.extract_table_from_image(image_path, page_num)
            if result:
                page_results.append(result)
        except Exception as e:
            logger.error(f"  ❌ Error: {e}")
            continue
    # Cleanup
    logger.info("🗑️  Cleaning up temporary PNGs...")
    for image_path in page_images.values():
        try:
            if os.path.exists(image_path):
                os.remove(image_path)
                logger.info(f"  ✓ Deleted: {Path(image_path).name}")
        except Exception as e:
            logger.warning(f"⚠️ Could not delete {image_path}: {e}")
    # Summary
    total_tables = sum(extractor._count_tables(r) for r in page_results)
    logger.info("=" * 60)
    logger.info(f"✅ EXTRACTION COMPLETE!")
    logger.info(f"   🔍 Pages detected: {len(table_pages)}")
    logger.info(f"   📋 Pages extracted: {len(page_results)}")
    logger.info(f"   📊 Total tables found: {total_tables}")
    logger.info("=" * 60)
    if save_output and page_results:
        save_results_to_json(page_results, pdf_path)
    return page_results

def save_results_to_json(results: List[Dict], pdf_path: str):
    """Save complete Reducto results to JSON."""
    try:
        pdf_name = Path(pdf_path).stem
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        output_data = {
            "filename": Path(pdf_path).name,
            "extraction_date": datetime.now().isoformat(),
            "extraction_method": "pdfplumber_detection + reducto_pipeline_image_extraction",
            "total_pages_processed": len(results),
            "pages": results,
        }
        json_file = f"{OUTPUT_DIR}/reducto/{pdf_name}_reducto_{timestamp}.json"
        with open(json_file, 'w') as f:
            json.dump(output_data, f, indent=2)
        logger.info(f"💾 Saved: {json_file}")
    except Exception as e:
        logger.error(f"Failed to save: {e}")

if __name__ == "__main__":
    import sys
    if len(sys.argv) < 2:
        print("Usage: python extractor.py <pdf_file>")
        sys.exit(1)
    pdf_filename = Path(sys.argv[1]).name
    pdf_path = Path(PDF_DIR) / pdf_filename
    if not pdf_path.exists():
        logger.error(f"❌ Specified PDF not found: {pdf_path}")
        sys.exit(1)
    logger.info(f"📁 Processing specific PDF: {pdf_path.name}")
    results = extract_tables_smart(str(pdf_path))
    print(f"\n✅ Extracted data from {len(results)} pages!")