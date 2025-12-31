#!/usr/bin/env python3
"""
PyMuPDF Watermark Removal Test Script
Uses PyMuPDF's content stream manipulation to remove watermarks.
"""
import os
import sys
from pathlib import Path
from typing import Dict, List, Optional, Set
import logging
import fitz  # PyMuPDF
import argparse
import re

# Logging setup
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Output directories
OUTPUT_DIR = "./watermark_test_output"
ORIGINAL_DIR = f"{OUTPUT_DIR}/original"
PROCESSED_DIR = f"{OUTPUT_DIR}/processed"
CLEANED_PDF_DIR = f"{OUTPUT_DIR}/cleaned_pdfs"

def setup_directories():
    """Create output directories."""
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    os.makedirs(ORIGINAL_DIR, exist_ok=True)
    os.makedirs(PROCESSED_DIR, exist_ok=True)
    os.makedirs(CLEANED_PDF_DIR, exist_ok=True)
    logger.info(f"📁 Output directory: {OUTPUT_DIR}")

def method_clean_contents(pdf_path: str, page_num: int, dpi: int = 300) -> Optional[str]:
    """
    Method 1: Use PyMuPDF's clean_contents() to remove redundant/watermark content.
    """
    try:
        pdf_document = fitz.open(pdf_path)
        page = pdf_document[page_num - 1]
        
        # Clean page contents (removes redundant operations)
        page.clean_contents()
        
        # Render
        zoom = dpi / 72
        mat = fitz.Matrix(zoom, zoom)
        pix = page.get_pixmap(matrix=mat, alpha=False, annots=False)
        
        pdf_name = Path(pdf_path).stem
        image_path = f"{PROCESSED_DIR}/{pdf_name}_page_{page_num}_clean_contents.png"
        pix.save(image_path)
        pdf_document.close()
        return image_path
    except Exception as e:
        logger.error(f"❌ Method 1 failed: {e}")
        return None

def method_remove_images_by_size(pdf_path: str, page_num: int, dpi: int = 300, 
                                  min_size: float = 0.3, max_size: float = 0.9) -> Optional[str]:
    """
    Method 2: Remove images that are likely watermarks based on size.
    Watermarks often cover a large portion of the page.
    """
    try:
        pdf_document = fitz.open(pdf_path)
        page = pdf_document[page_num - 1]
        
        # Get page dimensions
        page_rect = page.rect
        page_area = page_rect.width * page_rect.height
        
        # Get all images on the page
        image_list = page.get_images()
        removed_count = 0
        
        for img_index, img in enumerate(image_list):
            try:
                xref = img[0]
                base_image = pdf_document.extract_image(xref)
                image_bytes = base_image["image"]
                
                # Get image rectangles
                image_rects = page.get_image_rects(xref)
                
                for rect in image_rects:
                    # Calculate image area as percentage of page
                    img_area = rect.width * rect.height
                    img_percentage = img_area / page_area
                    
                    # If image is large (likely watermark), remove it
                    if min_size <= img_percentage <= max_size:
                        logger.info(f"    🗑️ Removing large image ({img_percentage:.1%} of page)")
                        # Delete the image from the page
                        page.delete_image(xref)
                        removed_count += 1
            except Exception as e:
                logger.warning(f"    ⚠️ Could not process image {img_index}: {e}")
                continue
        
        if removed_count > 0:
            logger.info(f"  ✓ Removed {removed_count} potential watermark image(s)")
        
        # Render
        zoom = dpi / 72
        mat = fitz.Matrix(zoom, zoom)
        pix = page.get_pixmap(matrix=mat, alpha=False, annots=False)
        
        pdf_name = Path(pdf_path).stem
        image_path = f"{PROCESSED_DIR}/{pdf_name}_page_{page_num}_remove_large_images.png"
        pix.save(image_path)
        pdf_document.close()
        return image_path
    except Exception as e:
        logger.error(f"❌ Method 2 failed: {e}")
        return None

def method_filter_content_stream(pdf_path: str, page_num: int, dpi: int = 300) -> Optional[str]:
    """
    Method 3: Filter content stream to remove watermark-related commands.
    This removes text/images with low opacity or specific patterns.
    """
    try:
        pdf_document = fitz.open(pdf_path)
        page = pdf_document[page_num - 1]
        
        # Get page contents
        contents = page.read_contents()
        
        # Convert to string for processing
        content_str = contents.decode('latin-1') if isinstance(contents, bytes) else str(contents)
        
        # Patterns that might indicate watermarks
        watermark_patterns = [
            r'/GS\d+\s+gs',  # Graphics state with transparency
            r'rg\s+rg\s+rg',  # RGB color commands (repeated)
            r'0\.\d+\s+g',    # Gray color (low values = light/transparent)
        ]
        
        # Try to filter out watermark-related commands
        # This is a simplified approach - you may need to adjust based on your PDFs
        filtered_content = content_str
        
        # Remove low-opacity graphics states (common for watermarks)
        # This is a heuristic - adjust based on your needs
        lines = filtered_content.split('\n')
        filtered_lines = []
        skip_next = False
        
        for i, line in enumerate(lines):
            # Skip lines with very low opacity (watermarks are often semi-transparent)
            if '/ca' in line or '/CA' in line:
                # Check opacity value
                opacity_match = re.search(r'(\d+\.?\d*)\s+ca', line)
                if opacity_match:
                    opacity = float(opacity_match.group(1))
                    if opacity < 0.3:  # Very transparent = likely watermark
                        skip_next = True
                        continue
            
            if not skip_next:
                filtered_lines.append(line)
            skip_next = False
        
        filtered_content = '\n'.join(filtered_lines)
        
        # Write filtered content back (if we made changes)
        if filtered_content != content_str:
            page.write_contents(filtered_content.encode('latin-1'))
            logger.info(f"  ✓ Filtered content stream")
        
        # Render
        zoom = dpi / 72
        mat = fitz.Matrix(zoom, zoom)
        pix = page.get_pixmap(matrix=mat, alpha=False, annots=False)
        
        pdf_name = Path(pdf_path).stem
        image_path = f"{PROCESSED_DIR}/{pdf_name}_page_{page_num}_filtered_stream.png"
        pix.save(image_path)
        pdf_document.close()
        return image_path
    except Exception as e:
        logger.error(f"❌ Method 3 failed: {e}")
        return None

def method_remove_text_by_properties(pdf_path: str, page_num: int, dpi: int = 300) -> Optional[str]:
    """
    Method 4: Remove text that has watermark-like properties (large, centered, low opacity).
    """
    try:
        pdf_document = fitz.open(pdf_path)
        page = pdf_document[page_num - 1]
        
        # Get all text blocks
        text_dict = page.get_text("dict")
        page_rect = page.rect
        page_center_x = page_rect.width / 2
        page_center_y = page_rect.height / 2
        
        removed_count = 0
        
        for block in text_dict.get("blocks", []):
            if "lines" not in block:
                continue
            
            for line in block["lines"]:
                for span in line.get("spans", []):
                    # Check if text has watermark characteristics
                    bbox = span.get("bbox", [])
                    if len(bbox) != 4:
                        continue
                    
                    x0, y0, x1, y1 = bbox
                    text_width = x1 - x0
                    text_height = y1 - y0
                    text_center_x = (x0 + x1) / 2
                    text_center_y = (y0 + y1) / 2
                    
                    # Check if text is large and centered (watermark characteristics)
                    is_large = text_width > page_rect.width * 0.3 or text_height > page_rect.height * 0.1
                    is_centered = (abs(text_center_x - page_center_x) < page_rect.width * 0.2 and
                                 abs(text_center_y - page_center_y) < page_rect.height * 0.2)
                    
                    # Check font size (watermarks are often large)
                    font_size = span.get("size", 0)
                    is_large_font = font_size > 20
                    
                    # If it matches watermark characteristics, try to remove
                    if (is_large and is_centered) or is_large_font:
                        try:
                            # Create a rectangle to redact
                            rect = fitz.Rect(bbox)
                            # Redact (remove) the text
                            page.add_redact_annot(rect, fill=(1, 1, 1))  # White fill
                            removed_count += 1
                        except Exception as e:
                            logger.warning(f"    ⚠️ Could not remove text: {e}")
        
        if removed_count > 0:
            page.apply_redactions()  # Apply all redactions
            logger.info(f"  ✓ Removed {removed_count} potential watermark text block(s)")
        
        # Render
        zoom = dpi / 72
        mat = fitz.Matrix(zoom, zoom)
        pix = page.get_pixmap(matrix=mat, alpha=False, annots=False)
        
        pdf_name = Path(pdf_path).stem
        image_path = f"{PROCESSED_DIR}/{pdf_name}_page_{page_num}_remove_text.png"
        pix.save(image_path)
        pdf_document.close()
        return image_path
    except Exception as e:
        logger.error(f"❌ Method 4 failed: {e}")
        return None

def method_remove_all_images(pdf_path: str, page_num: int, dpi: int = 300) -> Optional[str]:
    """
    Method 5: Remove ALL images from the page (nuclear option - use if watermarks are images).
    """
    try:
        pdf_document = fitz.open(pdf_path)
        page = pdf_document[page_num - 1]
        
        # Get all images
        image_list = page.get_images()
        removed_count = 0
        
        for img_index, img in enumerate(image_list):
            try:
                xref = img[0]
                page.delete_image(xref)
                removed_count += 1
            except Exception as e:
                logger.warning(f"    ⚠️ Could not remove image {img_index}: {e}")
        
        if removed_count > 0:
            logger.info(f"  ✓ Removed {removed_count} image(s)")
        
        # Render
        zoom = dpi / 72
        mat = fitz.Matrix(zoom, zoom)
        pix = page.get_pixmap(matrix=mat, alpha=False, annots=False)
        
        pdf_name = Path(pdf_path).stem
        image_path = f"{PROCESSED_DIR}/{pdf_name}_page_{page_num}_remove_all_images.png"
        pix.save(image_path)
        pdf_document.close()
        return image_path
    except Exception as e:
        logger.error(f"❌ Method 5 failed: {e}")
        return None

def method_combined_pymupdf(pdf_path: str, page_num: int, dpi: int = 300) -> Optional[str]:
    """
    Method 6: Combine multiple PyMuPDF techniques.
    """
    try:
        pdf_document = fitz.open(pdf_path)
        page = pdf_document[page_num - 1]
        
        # Step 1: Clean contents
        page.clean_contents()
        
        # Step 2: Remove large images (potential watermarks)
        page_rect = page.rect
        page_area = page_rect.width * page_rect.height
        image_list = page.get_images()
        
        for img_index, img in enumerate(image_list):
            try:
                xref = img[0]
                image_rects = page.get_image_rects(xref)
                for rect in image_rects:
                    img_area = rect.width * rect.height
                    img_percentage = img_area / page_area
                    if img_percentage > 0.2:  # Remove images > 20% of page
                        page.delete_image(xref)
                        break
            except:
                continue
        
        # Step 3: Remove centered large text
        text_dict = page.get_text("dict")
        page_center_x = page_rect.width / 2
        page_center_y = page_rect.height / 2
        
        for block in text_dict.get("blocks", []):
            if "lines" not in block:
                continue
            for line in block["lines"]:
                for span in line.get("spans", []):
                    bbox = span.get("bbox", [])
                    if len(bbox) == 4:
                        x0, y0, x1, y1 = bbox
                        text_center_x = (x0 + x1) / 2
                        text_center_y = (y0 + y1) / 2
                        text_width = x1 - x0
                        
                        if (abs(text_center_x - page_center_x) < page_rect.width * 0.3 and
                            abs(text_center_y - page_center_y) < page_rect.height * 0.3 and
                            text_width > page_rect.width * 0.4):
                            try:
                                rect = fitz.Rect(bbox)
                                page.add_redact_annot(rect, fill=(1, 1, 1))
                            except:
                                pass
        
        page.apply_redactions()
        
        # Render
        zoom = dpi / 72
        mat = fitz.Matrix(zoom, zoom)
        pix = page.get_pixmap(matrix=mat, alpha=False, annots=False)
        
        pdf_name = Path(pdf_path).stem
        image_path = f"{PROCESSED_DIR}/{pdf_name}_page_{page_num}_combined_pymupdf.png"
        pix.save(image_path)
        pdf_document.close()
        return image_path
    except Exception as e:
        logger.error(f"❌ Method 6 failed: {e}")
        return None

def convert_original(pdf_path: str, page_num: int, dpi: int = 300) -> Optional[str]:
    """Convert original page to image for comparison."""
    try:
        pdf_document = fitz.open(pdf_path)
        page = pdf_document[page_num - 1]
        zoom = dpi / 72
        mat = fitz.Matrix(zoom, zoom)
        pix = page.get_pixmap(matrix=mat, alpha=False, annots=False)
        
        pdf_name = Path(pdf_path).stem
        image_path = f"{ORIGINAL_DIR}/{pdf_name}_page_{page_num}_original.png"
        pix.save(image_path)
        pdf_document.close()
        return image_path
    except Exception as e:
        logger.error(f"❌ Original conversion failed: {e}")
        return None

def test_pymupdf_watermark_removal(pdf_path: str, page_numbers: Optional[List[int]] = None, dpi: int = 300):
    """Test all PyMuPDF watermark removal methods."""
    if not os.path.exists(pdf_path):
        logger.error(f"❌ PDF not found: {pdf_path}")
        return
    
    setup_directories()
    
    pdf_document = fitz.open(pdf_path)
    total_pages = len(pdf_document)
    pdf_document.close()
    
    if page_numbers is None:
        pages_to_process = list(range(1, total_pages + 1))
    else:
        pages_to_process = [p for p in page_numbers if 1 <= p <= total_pages]
    
    logger.info("=" * 70)
    logger.info("🧪 PyMuPDF WATERMARK REMOVAL TEST")
    logger.info("=" * 70)
    logger.info(f"📄 PDF: {Path(pdf_path).name}")
    logger.info(f"📊 Total pages: {total_pages}")
    logger.info(f"🎯 Processing pages: {pages_to_process}")
    logger.info(f"🖼️  DPI: {dpi}")
    logger.info("=" * 70)
    
    methods = [
        ("Clean Contents", method_clean_contents),
        ("Remove Large Images", method_remove_images_by_size),
        ("Filter Content Stream", method_filter_content_stream),
        ("Remove Watermark Text", method_remove_text_by_properties),
        ("Remove All Images", method_remove_all_images),
        ("Combined PyMuPDF", method_combined_pymupdf),
    ]
    
    for page_num in pages_to_process:
        logger.info(f"\n📄 Processing page {page_num}/{total_pages}...")
        
        # Original
        original_image = convert_original(pdf_path, page_num, dpi)
        if original_image:
            logger.info(f"  ✓ Original saved: {Path(original_image).name}")
        
        # Apply all methods
        pdf_name = Path(pdf_path).stem
        for method_name, method_func in methods:
            logger.info(f"  🔧 Testing: {method_name}...")
            try:
                result = method_func(pdf_path, page_num, dpi)
                if result and os.path.exists(result):
                    logger.info(f"    ✓ Saved: {Path(result).name}")
                else:
                    logger.warning(f"    ⚠️ Method failed or no output")
            except Exception as e:
                logger.error(f"    ❌ Error: {e}")
    
    logger.info("\n" + "=" * 70)
    logger.info("✅ TEST COMPLETE!")
    logger.info("=" * 70)
    logger.info(f"📁 Check results in: {OUTPUT_DIR}/")
    logger.info(f"   - Original images: {ORIGINAL_DIR}/")
    logger.info(f"   - Processed images: {PROCESSED_DIR}/")
    logger.info("\n💡 Compare the images to see which PyMuPDF method works best!")
    logger.info("=" * 70)

def main():
    parser = argparse.ArgumentParser(description='PyMuPDF watermark removal test')
    parser.add_argument('pdf_file', type=str, help='Path to PDF file')
    parser.add_argument('--pages', type=str, help='Page numbers (e.g., "1,2,3" or "1-5")')
    parser.add_argument('--dpi', type=int, default=300, help='DPI (default: 300)')
    
    args = parser.parse_args()
    
    page_numbers = None
    if args.pages:
        page_numbers = []
        for part in args.pages.split(','):
            part = part.strip()
            if '-' in part:
                start, end = map(int, part.split('-'))
                page_numbers.extend(range(start, end + 1))
            else:
                page_numbers.append(int(part))
        page_numbers = sorted(set(page_numbers))
    
    pdf_path = args.pdf_file
    if not os.path.isabs(pdf_path):
        input_pdfs_dir = Path(__file__).parent / "input_pdfs"
        if (input_pdfs_dir / pdf_path).exists():
            pdf_path = str(input_pdfs_dir / pdf_path)
        elif not os.path.exists(pdf_path):
            logger.error(f"❌ PDF not found: {pdf_path}")
            sys.exit(1)
    
    test_pymupdf_watermark_removal(pdf_path, page_numbers, args.dpi)

if __name__ == "__main__":
    main()