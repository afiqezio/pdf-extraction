import fitz
import logging
import re

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

doc = fitz.open("final_extraction_system/input_pdfs/test.pdf")

# Step 1: Find and collect all watermark-related XRefs
watermark_xrefs = set()
ocg_xrefs = set()

logger.info("Step 1: Scanning for watermark XObjects and OCGs...")

# Scan all pages for XObjects
for page_num, page in enumerate(doc):
    xobject_list = page.get_xobjects()
    for xref_info in xobject_list:
        xref = xref_info[0]
        try:
            xref_obj = doc.xref_object(xref, compressed=False)
            
            # Check multiple watermark indicators
            is_watermark = False
            
            # Check 1: Direct watermark marker
            if "/Watermark" in xref_obj:
                is_watermark = True
                logger.info(f"  Found watermark XObject (direct) at xref {xref}")
            
            # Check 2: OC reference (Optional Content)
            if "/OC" in xref_obj:
                oc_ref = xref_obj.split("/OC")[1].split()[0]
                # Extract the OCG xref number
                oc_match = re.search(r'(\d+)\s+0\s+R', oc_ref)
                if oc_match:
                    ocg_xref = int(oc_match.group(1))
                    ocg_xrefs.add(ocg_xref)
                    is_watermark = True
                    logger.info(f"  Found watermark XObject (OC reference) at xref {xref}, OCG at {ocg_xref}")
            
            # Check 3: PieceInfo with Watermark
            if "/PieceInfo" in xref_obj and "/Watermark" in xref_obj:
                is_watermark = True
                logger.info(f"  Found watermark XObject (PieceInfo) at xref {xref}")
            
            if is_watermark:
                watermark_xrefs.add(xref)
                
        except Exception as e:
            logger.warning(f"  Error checking xref {xref}: {e}")
            continue

# Step 2: Check OCG objects for watermark names
logger.info("\nStep 2: Scanning OCG objects...")
for xref in range(1, doc.xref_length()):
    try:
        obj_type = doc.xref_get_key(xref, "Type")
        if obj_type[1] == "/OCG":
            obj_str = doc.xref_object(xref, compressed=False)
            if "/Name(Watermark)" in obj_str or "/Name/Watermark" in obj_str:
                ocg_xrefs.add(xref)
                logger.info(f"  Found watermark OCG at xref {xref}")
    except:
        continue

# Step 3: Clear watermark XObject streams (safer than deleting)
logger.info(f"\nStep 3: Clearing {len(watermark_xrefs)} watermark XObject stream(s)...")
for xref in watermark_xrefs:
    try:
        # Clear the XObject stream content (makes it empty/invisible)
        doc.update_stream(xref, b"")
        logger.info(f"  Cleared XObject stream at xref {xref}")
    except Exception as e:
        logger.warning(f"  Could not clear xref {xref}: {e}")

# Step 4: Clean content streams (conservative approach)
logger.info("\nStep 4: Cleaning content streams (conservative)...")
for page_num, page in enumerate(doc):
    try:
        page.clean_contents()
        
        # Get contents xrefs properly
        contents_xrefs = page.get_contents()
        if not contents_xrefs:
            continue
        
        # Handle both single xref and array of xrefs
        if isinstance(contents_xrefs, list):
            xref_list = contents_xrefs
        else:
            xref_list = [contents_xrefs]
        
        for xref in xref_list:
            try:
                # Read the content stream
                content_bytes = doc.xref_stream(xref)
                if not content_bytes:
                    continue
                
                content_str = content_bytes.decode('latin-1', errors='ignore')
                original_str = content_str
                
                # Only remove watermark artifact blocks (conservative)
                lines = content_str.split('\n')
                filtered_lines = []
                skip_until_emc = False
                in_watermark_artifact = False
                
                for i, line in enumerate(lines):
                    # Only skip complete watermark artifact blocks
                    if "/Artifact" in line and "/Watermark" in line:
                        in_watermark_artifact = True
                        skip_until_emc = True
                        continue
                    
                    if skip_until_emc:
                        if line.strip() == "EMC":
                            skip_until_emc = False
                            in_watermark_artifact = False
                        continue
                    
                    # Keep all other content (don't remove OC commands - they might be needed)
                    filtered_lines.append(line)
                
                # Only update if we actually removed something
                new_content = '\n'.join(filtered_lines)
                if new_content != original_str and in_watermark_artifact == False:
                    # Use doc.update_stream() - the correct method
                    doc.update_stream(xref, new_content.encode('latin-1'))
                    logger.info(f"  ✓ Page {page_num + 1}: Cleaned content stream xref {xref}")
                    
            except Exception as e:
                logger.warning(f"  Page {page_num + 1}: Could not process content xref {xref}: {e}")
                continue
                
    except Exception as e:
        logger.warning(f"  Page {page_num + 1}: Could not clean content: {e}")

# Step 5: Remove OCG properties from document catalog
logger.info("\nStep 5: Removing OCG properties from document...")
try:
    root_xref = doc.pdf_catalog()
    root_obj = doc.xref_object(root_xref, compressed=False)
    if "/OCProperties" in root_obj:
        logger.info("  Found OCProperties in catalog (keeping for compatibility)")
except Exception as e:
    logger.warning(f"  Could not check OCProperties: {e}")

# Save the cleaned PDF
output_path = "no_watermark.pdf"
doc.save(output_path, garbage=4, deflate=True)
doc.close()

logger.info(f"\n✅ Saved cleaned PDF: {output_path}")
logger.info(f"   Cleared {len(watermark_xrefs)} watermark XObject stream(s)")
logger.info(f"   Found {len(ocg_xrefs)} watermark OCG(s)")