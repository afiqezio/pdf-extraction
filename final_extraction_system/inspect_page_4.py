import pdfplumber

pdf_path = "/Users/husaini.rosdi/Documents/working_project/pdf-extraction/PPDM Data Model 3.9 Documentation.pdf"

with pdfplumber.open(pdf_path) as pdf:
    page = pdf.pages[3] # Page 4 (0-indexed is 3)
    text = page.extract_text()
    print("--- TEXT ---")
    print(text[:500]) # First 500 chars
    print("--- END TEXT ---")

