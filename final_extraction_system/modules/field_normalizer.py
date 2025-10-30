"""
Field Normalizer Module

Purpose: Normalize column names from PDFs to standardized field names
- Handles exact matches (e.g., "Quartz" → "quartz")
- Handles fuzzy matches for typos (e.g., "Qaurtz" → "quartz")
- Tracks unmapped fields for review
"""

import json
import logging
from pathlib import Path
from typing import Dict, List, Tuple, Optional
from rapidfuzz import fuzz, process

logger = logging.getLogger(__name__)


class FieldNormalizer:
    """Normalizes field names using exact and fuzzy matching"""
    
    def __init__(self, mappings_file: str = "field_mappings.json"):
        """
        Initialize the field normalizer
        
        Args:
            mappings_file: Path to the JSON file containing field mappings
        """
        self.mappings_file = Path(mappings_file)
        self.mappings = {}
        self.fuzzy_threshold = 85
        self.unmapped_fields = []
        self.fuzzy_matched_fields = []
        
        # Load mappings from JSON file
        self._load_mappings()
        
    def _load_mappings(self):
        """Load field mappings from JSON file"""
        try:
            with open(self.mappings_file, 'r', encoding='utf-8') as f:
                data = json.load(f)
                
            self.mappings = data.get('mappings', {})
            self.fuzzy_threshold = data.get('fuzzy_matching', {}).get('threshold', 85)
            
            logger.info(f"✅ Loaded {len(self.mappings)} field mappings from {self.mappings_file}")
            logger.info(f"🎯 Fuzzy matching threshold: {self.fuzzy_threshold}%")
            
        except FileNotFoundError:
            logger.error(f"❌ Mappings file not found: {self.mappings_file}")
            raise
        except json.JSONDecodeError as e:
            logger.error(f"❌ Invalid JSON in mappings file: {e}")
            raise
    
    def normalize_field(self, field_name: str) -> Tuple[str, str, float]:
        """
        Normalize a single field name
        
        Args:
            field_name: Original field name from PDF
            
        Returns:
            Tuple of (normalized_name, match_type, confidence)
            - normalized_name: The standardized field name
            - match_type: 'exact', 'fuzzy', or 'unmapped'
            - confidence: Match confidence (0-100)
        """
        if not field_name or not isinstance(field_name, str):
            return field_name, 'unmapped', 0.0
        
        # Clean the field name
        cleaned = field_name.strip().lower()
        
        # Try exact match first
        if cleaned in self.mappings:
            normalized = self.mappings[cleaned]
            logger.debug(f"✅ Exact match: '{field_name}' → '{normalized}'")
            return normalized, 'exact', 100.0
        
        # Try fuzzy match
        normalized, match_type, confidence = self._fuzzy_match(cleaned, field_name)
        
        # Track unmapped fields
        if match_type == 'unmapped':
            if field_name not in self.unmapped_fields:
                self.unmapped_fields.append(field_name)
                logger.warning(f"⚠️ Unmapped field: '{field_name}'")
        
        return normalized, match_type, confidence
    
    def _fuzzy_match(self, cleaned_name: str, original_name: str) -> Tuple[str, str, float]:
        """
        Attempt fuzzy matching for typos and variations
        
        Args:
            cleaned_name: Cleaned lowercase field name
            original_name: Original field name from PDF
            
        Returns:
            Tuple of (normalized_name, match_type, confidence)
        """
        # Get all possible field name variations from mappings (the keys)
        possible_matches = list(self.mappings.keys())
        
        # Find best match using rapidfuzz
        result = process.extractOne(
            cleaned_name,
            possible_matches,
            scorer=fuzz.ratio,
            score_cutoff=self.fuzzy_threshold
        )
        
        if result:
            matched_key, confidence, _ = result
            normalized = self.mappings[matched_key]
            
            # Track fuzzy matches for user review
            fuzzy_info = {
                'original': original_name,
                'matched_to': matched_key,
                'normalized': normalized,
                'confidence': confidence
            }
            if fuzzy_info not in self.fuzzy_matched_fields:
                self.fuzzy_matched_fields.append(fuzzy_info)
                logger.info(f"🔍 Fuzzy match ({confidence:.1f}%): '{original_name}' → '{normalized}'")
            
            return normalized, 'fuzzy', confidence
        
        # No match found
        logger.debug(f"❌ No match found for: '{original_name}'")
        return original_name, 'unmapped', 0.0
    
    def normalize_headers(self, headers: List[str]) -> Dict:
        """
        Normalize a list of column headers
        
        Args:
            headers: List of original column names from PDF
            
        Returns:
            Dictionary containing:
            - normalized_headers: List of standardized field names
            - original_headers: Original headers (unchanged)
            - mapping: Dictionary mapping original → normalized
            - match_info: List of match information for each field
        """
        normalized_headers = []
        mapping = {}
        match_info = []
        
        for header in headers:
            normalized, match_type, confidence = self.normalize_field(header)
            
            normalized_headers.append(normalized)
            mapping[header] = normalized
            match_info.append({
                'original': header,
                'normalized': normalized,
                'match_type': match_type,
                'confidence': confidence
            })
        
        logger.info(f"📊 Normalized {len(headers)} headers:")
        logger.info(f"   ✅ Exact matches: {sum(1 for m in match_info if m['match_type'] == 'exact')}")
        logger.info(f"   🔍 Fuzzy matches: {sum(1 for m in match_info if m['match_type'] == 'fuzzy')}")
        logger.info(f"   ⚠️ Unmapped: {sum(1 for m in match_info if m['match_type'] == 'unmapped')}")
        
        return {
            'normalized_headers': normalized_headers,
            'original_headers': headers,
            'mapping': mapping,
            'match_info': match_info
        }
    
    def normalize_table(self, table: Dict) -> Dict:
        """
        Normalize all field names in a table
        
        Args:
            table: Table dictionary with 'headers' and 'rows'
            
        Returns:
            Enhanced table dictionary with normalization information
        """
        if 'headers' not in table:
            logger.warning("⚠️ Table has no 'headers' field")
            return table
        
        # Normalize headers
        normalization_result = self.normalize_headers(table['headers'])
        
        # Add normalization info to table
        table['normalized_headers'] = normalization_result['normalized_headers']
        table['original_headers'] = normalization_result['original_headers']
        table['header_mapping'] = normalization_result['mapping']
        table['normalization_info'] = normalization_result['match_info']
        
        return table
    
    def normalize_tables(self, tables: List[Dict]) -> List[Dict]:
        """
        Normalize field names in multiple tables
        
        Args:
            tables: List of table dictionaries
            
        Returns:
            List of enhanced table dictionaries with normalization info
        """
        logger.info(f"🚀 Normalizing {len(tables)} tables...")
        
        normalized_tables = []
        for i, table in enumerate(tables):
            logger.info(f"📋 Normalizing table {i+1}/{len(tables)}...")
            normalized_table = self.normalize_table(table)
            normalized_tables.append(normalized_table)
        
        logger.info(f"✅ Normalization complete!")
        logger.info(f"   Total unmapped fields: {len(self.unmapped_fields)}")
        logger.info(f"   Total fuzzy matches: {len(self.fuzzy_matched_fields)}")
        
        return normalized_tables
    
    def get_unmapped_fields(self) -> List[str]:
        """Get list of fields that couldn't be mapped"""
        return self.unmapped_fields
    
    def get_fuzzy_matches(self) -> List[Dict]:
        """Get list of fields that were fuzzy matched (for user review)"""
        return self.fuzzy_matched_fields
    
    def save_unmapped_fields(self, output_file: str = "unmapped_fields.json"):
        """
        Save unmapped fields to JSON file for review
        
        Args:
            output_file: Path to output JSON file
        """
        if not self.unmapped_fields and not self.fuzzy_matched_fields:
            logger.info("✅ No unmapped or fuzzy matched fields to save")
            return
        
        output_data = {
            'unmapped_fields': self.unmapped_fields,
            'fuzzy_matches': self.fuzzy_matched_fields,
            'total_unmapped': len(self.unmapped_fields),
            'total_fuzzy': len(self.fuzzy_matched_fields)
        }
        
        output_path = Path(output_file)
        with open(output_path, 'w', encoding='utf-8') as f:
            json.dump(output_data, f, indent=2, ensure_ascii=False)
        
        logger.info(f"💾 Saved unmapped/fuzzy fields to: {output_path}")
        logger.info(f"   📊 {len(self.unmapped_fields)} unmapped fields")
        logger.info(f"   🔍 {len(self.fuzzy_matched_fields)} fuzzy matches")


# Example usage
if __name__ == "__main__":
    # Setup logging
    logging.basicConfig(level=logging.INFO)
    
    # Example table data
    example_table = {
        'headers': ['WELL', 'Depth (m)', 'Qaurtz', 'Poly Quartz', 'Felspar', 'Unknown Column'],
        'rows': [
            ['Well-1', '2880', '5', '7', '2', '10'],
            ['Well-2', '2900', '6', '8', '3', '12']
        ]
    }
    
    # Initialize normalizer
    normalizer = FieldNormalizer()
    
    # Normalize table
    normalized = normalizer.normalize_table(example_table)
    
    # Print results
    print("\n" + "="*60)
    print("NORMALIZATION RESULTS")
    print("="*60)
    print(f"\nOriginal headers: {normalized['original_headers']}")
    print(f"Normalized headers: {normalized['normalized_headers']}")
    print(f"\nMapping:")
    for orig, norm in normalized['header_mapping'].items():
        print(f"  '{orig}' → '{norm}'")
    
    print(f"\nMatch Info:")
    for info in normalized['normalization_info']:
        print(f"  {info['original']}: {info['match_type']} ({info['confidence']:.1f}%)")
    
    # Save unmapped fields
    normalizer.save_unmapped_fields()




