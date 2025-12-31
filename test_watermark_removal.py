#!/usr/bin/env python3
"""
Advanced Watermark Removal Test Script
Uses image processing techniques to remove watermarks embedded in PDF content.
"""
import os
import sys
from pathlib import Path
from typing import Dict, List, Optional, Tuple
import logging
import fitz  # PyMuPDF
from PIL import Image, ImageEnhance, ImageFilter, ImageOps, ImageChops
import numpy as np
import argparse

# Try to import OpenCV for advanced processing (optional)
try:
    import cv2
    OPENCV_AVAILABLE = True
except ImportError:
    OPENCV_AVAILABLE = False
    logging.warning("OpenCV not available. Some advanced methods will be skipped.")

# Logging setup
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Output directories
OUTPUT_DIR = "./watermark_test_output"
ORIGINAL_DIR = f"{OUTPUT_DIR}/original"
PROCESSED_DIR = f"{OUTPUT_DIR}/processed"

def setup_directories():
    """Create output directories."""
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    os.makedirs(ORIGINAL_DIR, exist_ok=True)
    os.makedirs(PROCESSED_DIR, exist_ok=True)
    logger.info(f"📁 Output directory: {OUTPUT_DIR}")

def convert_to_image(pdf_path: str, page_num: int, dpi: int = 300) -> Optional[str]:
    """Convert PDF page to image."""
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
        logger.error(f"❌ Conversion failed: {e}")
        return None

def method_high_contrast_brightness(image_path: str, output_path: str) -> str:
    """
    Method 1: Aggressive contrast and brightness adjustment.
    Works for light/semi-transparent watermarks.
    """
    try:
        img = Image.open(image_path).convert('RGB')
        
        # Aggressive contrast increase
        enhancer = ImageEnhance.Contrast(img)
        img = enhancer.enhance(1.5)
        
        # Increase brightness
        enhancer = ImageEnhance.Brightness(img)
        img = enhancer.enhance(1.2)
        
        # Sharpen
        img = img.filter(ImageFilter.SHARPEN)
        
        img.save(output_path, 'PNG', quality=95)
        return output_path
    except Exception as e:
        logger.warning(f"⚠️ Method 1 failed: {e}")
        return image_path

def method_color_separation(image_path: str, output_path: str) -> str:
    """
    Method 2: Separate colors and enhance to reduce watermark visibility.
    """
    try:
        img = Image.open(image_path).convert('RGB')
        img_array = np.array(img)
        
        # Convert to different color spaces and enhance
        # Increase saturation to make watermark less visible
        hsv = cv2.cvtColor(img_array, cv2.COLOR_RGB2HSV) if OPENCV_AVAILABLE else None
        
        if hsv is not None:
            # Adjust saturation
            hsv[:, :, 1] = np.clip(hsv[:, :, 1] * 1.2, 0, 255)
            # Adjust value (brightness)
            hsv[:, :, 2] = np.clip(hsv[:, :, 2] * 1.1, 0, 255)
            img_array = cv2.cvtColor(hsv, cv2.COLOR_HSV2RGB)
            img = Image.fromarray(img_array)
        else:
            # Fallback to PIL
            enhancer = ImageEnhance.Color(img)
            img = enhancer.enhance(1.3)
            enhancer = ImageEnhance.Brightness(img)
            img = enhancer.enhance(1.15)
        
        img.save(output_path, 'PNG', quality=95)
        return output_path
    except Exception as e:
        logger.warning(f"⚠️ Method 2 failed: {e}")
        return image_path

def method_morphological_operations(image_path: str, output_path: str) -> str:
    """
    Method 3: Use morphological operations to remove watermark patterns.
    Works for watermarks with specific patterns/text.
    """
    if not OPENCV_AVAILABLE:
        logger.warning("⚠️ OpenCV not available, skipping Method 3")
        return image_path
    
    try:
        img = cv2.imread(image_path)
        gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
        
        # Create a mask for light areas (watermarks are often light)
        _, mask = cv2.threshold(gray, 200, 255, cv2.THRESH_BINARY)
        
        # Invert mask to get watermark areas
        watermark_mask = cv2.bitwise_not(mask)
        
        # Use inpainting to fill watermark areas
        result = cv2.inpaint(img, watermark_mask, 3, cv2.INPAINT_TELEA)
        
        cv2.imwrite(output_path, result)
        return output_path
    except Exception as e:
        logger.warning(f"⚠️ Method 3 failed: {e}")
        return image_path

def method_adaptive_threshold(image_path: str, output_path: str) -> str:
    """
    Method 4: Adaptive thresholding to separate watermark from content.
    """
    if not OPENCV_AVAILABLE:
        logger.warning("⚠️ OpenCV not available, skipping Method 4")
        return image_path
    
    try:
        img = cv2.imread(image_path)
        gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
        
        # Adaptive threshold to find watermark
        adaptive = cv2.adaptiveThreshold(
            gray, 255, cv2.ADAPTIVE_THRESH_GAUSSIAN_C, 
            cv2.THRESH_BINARY, 11, 2
        )
        
        # Create mask for watermark (light areas)
        _, mask = cv2.threshold(gray, 220, 255, cv2.THRESH_BINARY)
        
        # Inpaint
        result = cv2.inpaint(img, mask, 5, cv2.INPAINT_NS)
        
        cv2.imwrite(output_path, result)
        return output_path
    except Exception as e:
        logger.warning(f"⚠️ Method 4 failed: {e}")
        return image_path

def method_frequency_domain(image_path: str, output_path: str) -> str:
    """
    Method 5: Frequency domain filtering (FFT) to remove periodic watermarks.
    """
    if not OPENCV_AVAILABLE:
        logger.warning("⚠️ OpenCV not available, skipping Method 5")
        return image_path
    
    try:
        img = cv2.imread(image_path)
        gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
        
        # Apply FFT
        f_transform = np.fft.fft2(gray)
        f_shift = np.fft.fftshift(f_transform)
        
        # Create a mask to filter out high frequencies (watermark patterns)
        rows, cols = gray.shape
        crow, ccol = rows // 2, cols // 2
        
        # Create a circular mask
        mask = np.ones((rows, cols), np.uint8)
        r = 30  # Radius - adjust based on watermark pattern
        center = [crow, ccol]
        x, y = np.ogrid[:rows, :cols]
        mask_area = (x - center[0]) ** 2 + (y - center[1]) ** 2 <= r * r
        mask[mask_area] = 0
        
        # Apply mask and inverse FFT
        f_shift = f_shift * mask
        f_ishift = np.fft.ifftshift(f_shift)
        img_back = np.fft.ifft2(f_ishift)
        img_back = np.abs(img_back)
        
        # Convert back to BGR and combine with original color
        img_back = np.uint8(img_back)
        img_back = cv2.cvtColor(img_back, cv2.COLOR_GRAY2BGR)
        
        # Blend with original
        result = cv2.addWeighted(img, 0.7, img_back, 0.3, 0)
        
        cv2.imwrite(output_path, result)
        return output_path
    except Exception as e:
        logger.warning(f"⚠️ Method 5 failed: {e}")
        return image_path

def method_histogram_equalization(image_path: str, output_path: str) -> str:
    """
    Method 6: Histogram equalization to normalize and reduce watermark visibility.
    """
    if not OPENCV_AVAILABLE:
        logger.warning("⚠️ OpenCV not available, skipping Method 6")
        return image_path
    
    try:
        img = cv2.imread(image_path)
        
        # Convert to LAB color space
        lab = cv2.cvtColor(img, cv2.COLOR_BGR2LAB)
        l, a, b = cv2.split(lab)
        
        # Apply CLAHE (Contrast Limited Adaptive Histogram Equalization)
        clahe = cv2.createCLAHE(clipLimit=2.0, tileGridSize=(8, 8))
        l = clahe.apply(l)
        
        # Merge channels
        lab = cv2.merge([l, a, b])
        result = cv2.cvtColor(lab, cv2.COLOR_LAB2BGR)
        
        cv2.imwrite(output_path, result)
        return output_path
    except Exception as e:
        logger.warning(f"⚠️ Method 6 failed: {e}")
        return image_path

def method_denoising(image_path: str, output_path: str) -> str:
    """
    Method 7: Denoising to remove watermark artifacts.
    """
    if not OPENCV_AVAILABLE:
        logger.warning("⚠️ OpenCV not available, skipping Method 7")
        return image_path
    
    try:
        img = cv2.imread(image_path)
        
        # Non-local means denoising
        result = cv2.fastNlMeansDenoisingColored(img, None, 10, 10, 7, 21)
        
        # Additional sharpening
        kernel = np.array([[-1, -1, -1],
                          [-1,  9, -1],
                          [-1, -1, -1]])
        result = cv2.filter2D(result, -1, kernel)
        
        cv2.imwrite(output_path, result)
        return output_path
    except Exception as e:
        logger.warning(f"⚠️ Method 7 failed: {e}")
        return image_path

def method_combine_all(image_path: str, output_path: str) -> str:
    """
    Method 8: Combine multiple techniques for best results.
    """
    try:
        # Start with high contrast/brightness
        temp_path1 = output_path.replace('.png', '_temp1.png')
        method_high_contrast_brightness(image_path, temp_path1)
        
        if OPENCV_AVAILABLE:
            # Apply denoising
            temp_path2 = output_path.replace('.png', '_temp2.png')
            method_denoising(temp_path1, temp_path2)
            
            # Apply histogram equalization
            method_histogram_equalization(temp_path2, output_path)
            
            # Clean up temp files
            for temp in [temp_path1, temp_path2]:
                if os.path.exists(temp):
                    os.remove(temp)
        else:
            # Fallback: just use PIL enhancements
            img = Image.open(temp_path1)
            enhancer = ImageEnhance.Sharpness(img)
            img = enhancer.enhance(1.5)
            img.save(output_path, 'PNG', quality=95)
            if os.path.exists(temp_path1):
                os.remove(temp_path1)
        
        return output_path
    except Exception as e:
        logger.warning(f"⚠️ Method 8 failed: {e}")
        return image_path

def test_watermark_removal(pdf_path: str, page_numbers: Optional[List[int]] = None, dpi: int = 300):
    """Test all advanced watermark removal methods."""
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
    logger.info("🧪 ADVANCED WATERMARK REMOVAL TEST")
    logger.info("=" * 70)
    logger.info(f"📄 PDF: {Path(pdf_path).name}")
    logger.info(f"📊 Total pages: {total_pages}")
    logger.info(f"🎯 Processing pages: {pages_to_process}")
    logger.info(f"🖼️  DPI: {dpi}")
    if OPENCV_AVAILABLE:
        logger.info("✅ OpenCV available - all methods enabled")
    else:
        logger.info("⚠️ OpenCV not available - some methods disabled")
    logger.info("=" * 70)
    
    methods = [
        ("High Contrast/Brightness", method_high_contrast_brightness),
        ("Color Separation", method_color_separation),
        ("Morphological Operations", method_morphological_operations),
        ("Adaptive Threshold", method_adaptive_threshold),
        ("Frequency Domain (FFT)", method_frequency_domain),
        ("Histogram Equalization", method_histogram_equalization),
        ("Denoising", method_denoising),
        ("Combined All Methods", method_combine_all),
    ]
    
    for page_num in pages_to_process:
        logger.info(f"\n📄 Processing page {page_num}/{total_pages}...")
        
        # Convert to image first
        original_image = convert_to_image(pdf_path, page_num, dpi)
        if not original_image:
            logger.error(f"  ❌ Failed to convert page {page_num}")
            continue
        
        logger.info(f"  ✓ Original saved: {Path(original_image).name}")
        
        # Apply all methods
        pdf_name = Path(pdf_path).stem
        for method_name, method_func in methods:
            if not OPENCV_AVAILABLE and method_name in ["Morphological Operations", "Adaptive Threshold", 
                                                          "Frequency Domain (FFT)", "Histogram Equalization", 
                                                          "Denoising"]:
                continue
            
            logger.info(f"  🔧 Testing: {method_name}...")
            try:
                output_path = f"{PROCESSED_DIR}/{pdf_name}_page_{page_num}_{method_name.lower().replace(' ', '_').replace('(', '').replace(')', '')}.png"
                result = method_func(original_image, output_path)
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
    logger.info("\n💡 Compare the images to see which method works best!")
    logger.info("=" * 70)

def main():
    parser = argparse.ArgumentParser(description='Advanced watermark removal test')
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
    
    test_watermark_removal(pdf_path, page_numbers, args.dpi)

if __name__ == "__main__":
    main()