import pdfplumber
import pandas as pd

pdf_path = "/Users/husaini.rosdi/Documents/working_project/pdf-extraction/PPDM Data Model 3.9 Documentation.pdf"

def analyze_pdf(path):
    print(f"Analyzing {path}...")
    try:
        with pdfplumber.open(path) as pdf:
            print(f"Total pages: {len(pdf.pages)}")
            # Check first 20 pages for tables to see headers
            for i, page in enumerate(pdf.pages[:20]):
                tables = page.extract_tables()
                if tables:
                    print(f"\nPage {i+1} Tables:")
                    for table in tables:
                        # Print first row (headers)
                        if table:
                            cleaned_headers = [str(h).replace('\n', ' ').strip() for h in table[0] if h]
                            print(f"Headers: {cleaned_headers}")
                            return # Found a table, let's stop and verify
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    analyze_pdf(pdf_path)

