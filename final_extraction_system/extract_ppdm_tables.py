import pdfplumber
import pandas as pd
import re
from tqdm import tqdm
import os
import argparse

# Configuration
DEFAULT_PDF_PATH = "/Users/husaini.rosdi/Documents/working_project/pdf-extraction/PPDM Data Model 3.9 Documentation.pdf"
OUTPUT_EXCEL = "PPDM_Data_Model_3.9_Tables.xlsx"

# Expected headers (based on analysis)
EXPECTED_HEADERS = ['Column Name', 'Null', 'Data Type', 'Length', 'Key', 'Ref Table(s)', 'Column Comment']

def extract_ppdm_tables(pdf_path=DEFAULT_PDF_PATH, output_excel=OUTPUT_EXCEL, limit_pages=None):
    print(f"Opening PDF: {pdf_path}")
    
    all_rows = []
    current_table_name = None
    
    if not os.path.exists(pdf_path):
        print(f"Error: File not found at {pdf_path}")
        return

    try:
        with pdfplumber.open(pdf_path) as pdf:
            total_pages = len(pdf.pages)
            
            if limit_pages:
                total_pages = min(total_pages, limit_pages)
                print(f"Processing first {total_pages} pages (use --all for full document)...")
            else:
                print(f"Processing ALL {total_pages} pages...")
            
            pages_to_process = pdf.pages[:total_pages]
            
            for page_num, page in enumerate(tqdm(pages_to_process, unit="page")):
                # 1. Extract text to find Table Name
                text = page.extract_text()
                if text:
                    # Look for "Table Name: XXXXX"
                    match = re.search(r"Table Name:\s*([A-Z0-9_]+)", text)
                    if match:
                        current_table_name = match.group(1)
                
                # 2. Extract tables
                tables = page.extract_tables()
                
                for table in tables:
                    if not table:
                        continue
                    
                    # Check headers
                    headers = [str(h).replace('\n', ' ').strip() for h in table[0] if h is not None]
                    
                    # Validation: must have Column Name and Data Type
                    if 'Column Name' in headers and 'Data Type' in headers:
                        # Process rows (skip header)
                        for row in table[1:]:
                            clean_row = [str(cell).strip() if cell else "" for cell in row]
                            
                            row_data = {
                                "Table Name": current_table_name,
                                "Page": page_num + 1
                            }
                            
                            if len(clean_row) >= 7:
                                row_data["Column Name"] = clean_row[0]
                                row_data["Null"] = clean_row[1]
                                row_data["Data Type"] = clean_row[2]
                                row_data["Length"] = clean_row[3]
                                row_data["Key"] = clean_row[4]
                                row_data["Ref Table(s)"] = clean_row[5]
                                row_data["Column Comment"] = clean_row[6]
                            else:
                                row_data["Column Name"] = clean_row[0] if len(clean_row) > 0 else ""
                                row_data["Data Type"] = clean_row[2] if len(clean_row) > 2 else ""
                                row_data["Column Comment"] = " ".join(clean_row[1:])
                                
                            all_rows.append(row_data)

    except Exception as e:
        print(f"Error processing PDF: {e}")
        return

    # Save to Excel
    print(f"\nExtracted {len(all_rows)} rows.")
    if all_rows:
        print(f"Saving to {output_excel}...")
        df = pd.DataFrame(all_rows)
        cols = ["Table Name", "Column Name", "Null", "Data Type", "Length", "Key", "Ref Table(s)", "Column Comment", "Page"]
        cols = [c for c in cols if c in df.columns]
        df = df[cols]
        
        df.to_excel(output_excel, index=False)
        print(f"✅ Successfully saved to {output_excel}")
    else:
        print("⚠️ No data extracted.")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Extract PPDM tables to Excel")
    parser.add_argument("--pdf", default=DEFAULT_PDF_PATH, help="Path to PDF file")
    parser.add_argument("--output", default=OUTPUT_EXCEL, help="Output Excel file")
    parser.add_argument("--limit", type=int, default=50, help="Limit number of pages (default: 50)")
    parser.add_argument("--all", action="store_true", help="Process all pages (overrides --limit)")
    
    args = parser.parse_args()
    
    limit = args.limit
    if args.all:
        limit = None
        
    extract_ppdm_tables(args.pdf, args.output, limit)
