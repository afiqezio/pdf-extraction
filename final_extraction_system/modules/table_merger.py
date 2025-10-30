"""
Table Merger Module

Purpose: Intelligently merge related tables that describe the same samples
- Detects tables without context (missing well/depth columns)
- Finds the previous table with context
- Merges them horizontally row-by-row if they match
- Validates row counts and page proximity before merging
"""

import json
import logging
from pathlib import Path
from typing import Dict, List, Tuple, Optional

logger = logging.getLogger(__name__)


class TableMerger:
    """Merges related tables based on context and proximity"""
    
    def __init__(self, context_fields_config: Dict = None):
        """
        Initialize the table merger
        
        Args:
            context_fields_config: Configuration for context field identification
        """
        # Default context fields if not provided
        self.primary_context_fields = context_fields_config.get('primary', [
            'well_name_field_name',
            'top_depth_mmddf',
            'bottom_depth_mmddf'
        ]) if context_fields_config else [
            'well_name_field_name',
            'top_depth_mmddf',
            'bottom_depth_mmddf'
        ]
        
        self.secondary_context_fields = context_fields_config.get('secondary', [
            'lithofacies_core'
        ]) if context_fields_config else [
            'lithofacies_core'
        ]
        
        self.merge_log = []
        
    def has_context(self, table: Dict) -> bool:
        """
        Check if table has context fields (well, depth)
        
        Args:
            table: Table dictionary with headers
            
        Returns:
            True if table has at least one primary context field
        """
        if 'headers' not in table:
            logger.warning("⚠️ Table missing headers, assuming no context")
            return False
        
        headers = table['headers']
        
        # Check if any primary context field exists (case-insensitive partial match)
        has_primary = False
        for field in self.primary_context_fields:
            # Check for partial matches (e.g., "well" matches "WELL NAME", "Well", etc.)
            field_lower = field.lower()
            for header in headers:
                if field_lower in str(header).lower():
                    has_primary = True
                    break
            if has_primary:
                break
        
        logger.debug(f"Table has context: {has_primary} (Headers: {headers[:3] if len(headers) > 3 else headers}...)")
        return has_primary
    
    def can_merge(self, table1: Dict, table2: Dict, strict: bool = False) -> Tuple[bool, str, float]:
        """
        Determine if two tables can be merged
        
        Args:
            table1: First table (with context)
            table2: Second table (without context)
            strict: If True, use strict validation; if False, use relaxed rules for extraction errors
            
        Returns:
            Tuple of (can_merge, reason, confidence_score)
            confidence_score: 0-100, how confident we are about the merge
        """
        confidence_score = 100.0
        
        # Rule 1: Table 1 must HAVE context, Table 2 must NOT have context
        if not self.has_context(table1):
            return False, "Table 1 lacks context", 0.0
        
        if self.has_context(table2):
            return False, "Table 2 already has context (not a continuation table)", 0.0
        
        # Get basic info
        rows1 = len(table1.get('rows', []))
        rows2 = len(table2.get('rows', []))
        page1 = table1.get('page', 0)
        page2 = table2.get('page', 0)
        
        # Rule 2: Check page proximity (MUST be close)
        page_diff = abs(page1 - page2)
        
        if page_diff > 2:
            return False, f"Pages too far apart: {page1} vs {page2}", 0.0
        
        # Reduce confidence if not on same page
        if page_diff == 1:
            confidence_score -= 10  # 90% confidence if consecutive pages
        elif page_diff == 2:
            confidence_score -= 20  # 80% confidence if 2 pages apart
        
        # Rule 3: Check row count similarity (FLEXIBLE for extraction errors)
        row_diff = abs(rows1 - rows2)
        row_ratio = min(rows1, rows2) / max(rows1, rows2) if max(rows1, rows2) > 0 else 0
        
        if strict:
            # Strict mode: rows must match closely (±2)
            if row_diff > 2:
                return False, f"Row count mismatch: {rows1} vs {rows2} (diff > 2)", 0.0
        else:
            # Relaxed mode: Allow larger differences BUT reduce confidence
            if row_diff == 0:
                # Perfect match
                confidence_score += 0
            elif row_diff <= 2:
                # Close match (within ±2)
                confidence_score -= 5
            elif row_diff <= 5:
                # Moderate difference (extraction errors likely)
                confidence_score -= 15
                logger.warning(f"⚠️ Moderate row mismatch: {rows1} vs {rows2} (likely extraction error)")
            elif row_ratio >= 0.5:
                # Large difference but ratio > 50% (e.g., 10 vs 6)
                confidence_score -= 30
                logger.warning(f"⚠️ Large row mismatch: {rows1} vs {rows2} (ratio: {row_ratio:.1%})")
            else:
                # Too different (ratio < 50%)
                return False, f"Row count too different: {rows1} vs {rows2} (ratio: {row_ratio:.1%})", 0.0
        
        # Rule 4: Check for column overlap (shouldn't have duplicate columns)
        headers1 = set(table1.get('headers', []))
        headers2 = set(table2.get('headers', []))
        overlap = headers1.intersection(headers2)
        
        # Remove context fields from overlap check (it's ok if they overlap on context)
        overlap = overlap - set(self.primary_context_fields) - set(self.secondary_context_fields)
        
        if len(overlap) > 0:
            return False, f"Tables have overlapping data columns: {overlap}", 0.0
        
        # Rule 5: Additional confidence boosts
        # If tables are on same page, boost confidence
        if page_diff == 0:
            confidence_score += 5
        
        # If row counts are very close, boost confidence
        if row_diff <= 1:
            confidence_score += 5
        
        # Ensure confidence is between 0-100
        confidence_score = max(0.0, min(100.0, confidence_score))
        
        # All checks passed!
        reason = f"Merge OK (confidence: {confidence_score:.1f}%, row_diff: {row_diff}, page_diff: {page_diff})"
        return True, reason, confidence_score
    
    def merge_two_tables(self, table1: Dict, table2: Dict, merge_confidence: float = 100.0) -> Dict:
        """
        Merge two tables horizontally (row-by-row)
        
        Args:
            table1: First table (with context)
            table2: Second table (without context)
            
        Returns:
            Merged table dictionary
        """
        logger.info(f"🔗 Merging Table {table1.get('id')} + Table {table2.get('id')}")
        
        # Get headers from both tables
        headers1 = table1.get('headers', [])
        headers2 = table2.get('headers', [])
        
        # Combine headers (table1 first, then table2)
        merged_headers = headers1 + headers2
        
        # Merge rows row-by-row
        rows1 = table1.get('rows', [])
        rows2 = table2.get('rows', [])
        merged_rows = []
        
        # Use the longer table as reference
        max_rows = max(len(rows1), len(rows2))
        
        for i in range(max_rows):
            # Get row from table1 (or empty if out of bounds)
            row1 = rows1[i] if i < len(rows1) else [''] * len(headers1)
            
            # Get row from table2 (or empty if out of bounds)
            row2 = rows2[i] if i < len(rows2) else [''] * len(headers2)
            
            # Combine rows
            merged_row = list(row1) + list(row2)
            merged_rows.append(merged_row)
        
        # Create merged table
        merged_table = {
            'id': f"{table1.get('id')}_merged_{table2.get('id')}",
            'page': table1.get('page'),  # Use first table's page
            'method': table1.get('method'),
            'confidence': min(table1.get('confidence', 0), table2.get('confidence', 0)),  # Use lower confidence
            'dimensions': f"{len(merged_rows)}x{len(merged_headers)}",
            
            # Headers
            'headers': merged_headers,
            
            # Data
            'rows': merged_rows,
            
            # Merge metadata
            'merged_from_tables': [table1.get('id'), table2.get('id')],
            'merge_confidence': merge_confidence,
            'merge_info': {
                'table1_id': table1.get('id'),
                'table2_id': table2.get('id'),
                'table1_rows': len(rows1),
                'table2_rows': len(rows2),
                'merged_rows': len(merged_rows),
                'table1_cols': len(headers1),
                'table2_cols': len(headers2),
                'merged_cols': len(merged_headers),
                'row_difference': abs(len(rows1) - len(rows2)),
                'confidence': merge_confidence
            },
            
            # Keep original table metadata
            'metadata': {
                'table1_metadata': table1.get('metadata', {}),
                'table2_metadata': table2.get('metadata', {}),
                'merge_type': 'horizontal',
                'merge_reason': 'table2_missing_context'
            }
        }
        
        logger.info(f"✅ Merged: {len(rows1)}x{len(headers1)} + {len(rows2)}x{len(headers2)} → {len(merged_rows)}x{len(merged_headers)}")
        
        # Log the merge
        self.merge_log.append({
            'table1_id': table1.get('id'),
            'table2_id': table2.get('id'),
            'merged_id': merged_table['id'],
            'rows': len(merged_rows),
            'cols': len(merged_headers)
        })
        
        return merged_table
    
    def merge_tables(self, tables: List[Dict]) -> List[Dict]:
        """
        Merge related tables in a list
        
        Strategy:
        1. Iterate through tables in order
        2. When we find a table WITHOUT context, look back for the most recent table WITH context
        3. If they can merge, merge them
        4. Continue until all tables processed
        
        Args:
            tables: List of table dictionaries (with headers)
            
        Returns:
            List of tables (some may be merged)
        """
        if not tables:
            return []
        
        logger.info(f"🔗 Starting table merging for {len(tables)} tables...")
        
        merged_tables = []
        skip_indices = set()  # Track which tables have been merged (skip them)
        
        for i, current_table in enumerate(tables):
            # Skip if this table was already merged into another
            if i in skip_indices:
                logger.debug(f"Skipping table {current_table.get('id')} (already merged)")
                continue
            
            # Check if current table has context
            if self.has_context(current_table):
                # This table has context, check if next table(s) can be merged with it
                logger.debug(f"Table {current_table.get('id')} has context, checking for merge candidates...")
                
                merged_table = current_table
                j = i + 1
                
                # Look ahead for tables without context
                while j < len(tables):
                    next_table = tables[j]
                    
                    # Skip if already merged
                    if j in skip_indices:
                        j += 1
                        continue
                    
                    # Check if we can merge
                    can_merge, reason, confidence = self.can_merge(merged_table, next_table, strict=False)
                    
                    if can_merge:
                        logger.info(f"✅ Merging table {merged_table.get('id')} with {next_table.get('id')} (confidence: {confidence:.1f}%)")
                        logger.info(f"   {reason}")
                        merged_table = self.merge_two_tables(merged_table, next_table, merge_confidence=confidence)
                        skip_indices.add(j)  # Mark next table as merged
                        j += 1
                    else:
                        logger.debug(f"❌ Cannot merge table {next_table.get('id')}: {reason}")
                        break  # Stop looking if we can't merge (not a continuation)
                
                merged_tables.append(merged_table)
                
            else:
                # This table has no context and wasn't merged
                logger.warning(f"⚠️ Table {current_table.get('id')} has no context and couldn't be merged")
                merged_tables.append(current_table)  # Keep it anyway
        
        logger.info(f"✅ Merging complete!")
        logger.info(f"   📊 Original tables: {len(tables)}")
        logger.info(f"   🔗 Merged tables: {len(merged_tables)}")
        logger.info(f"   ✂️ Tables consolidated: {len(tables) - len(merged_tables)}")
        
        return merged_tables
    
    def get_merge_summary(self) -> Dict:
        """
        Get summary of merge operations
        
        Returns:
            Dictionary with merge statistics
        """
        return {
            'total_merges': len(self.merge_log),
            'merge_log': self.merge_log
        }
    
    def save_merge_log(self, output_file: str = "merge_log.json"):
        """
        Save merge log to JSON file for review
        
        Args:
            output_file: Path to output JSON file
        """
        if not self.merge_log:
            logger.info("✅ No merges performed, nothing to save")
            return
        
        output_data = {
            'total_merges': len(self.merge_log),
            'merges': self.merge_log
        }
        
        output_path = Path(output_file)
        with open(output_path, 'w', encoding='utf-8') as f:
            json.dump(output_data, f, indent=2, ensure_ascii=False)
        
        logger.info(f"💾 Saved merge log to: {output_path}")
        logger.info(f"   🔗 Total merges: {len(self.merge_log)}")


# Example usage
if __name__ == "__main__":
    # Setup logging
    logging.basicConfig(level=logging.INFO)
    
    # Example: Table 1 with context
    table1 = {
        'id': 1,
        'page': 5,
        'headers': ['WELL', 'DEPTH (m)', 'Quartz', 'Feldspar'],
        'rows': [
            ['Well-1', '2880-2890', '5', '1'],
            ['Well-2', '2900-2910', '7', '2']
        ],
        'method': 'camelot_stream',
        'confidence': 95.0
    }
    
    # Example: Table 2 without context (same page, continuation)
    table2 = {
        'id': 2,
        'page': 5,
        'headers': ['Matrix', 'Calcite Cement', 'Porosity'],
        'rows': [
            ['36', '1', '2'],
            ['40', '2', '3']
        ],
        'method': 'camelot_stream',
        'confidence': 90.0
    }
    
    # Initialize merger
    merger = TableMerger()
    
    # Check if tables can merge
    can_merge, reason, confidence = merger.can_merge(table1, table2)
    print(f"\nCan merge: {can_merge}")
    print(f"Reason: {reason}")
    print(f"Confidence: {confidence:.1f}%")
    
    if can_merge:
        # Merge tables
        merged = merger.merge_two_tables(table1, table2, merge_confidence=confidence)
        
        print("\n" + "="*60)
        print("MERGE RESULTS")
        print("="*60)
        print(f"\nMerged ID: {merged['id']}")
        print(f"Dimensions: {merged['dimensions']}")
        print(f"Headers: {merged['headers']}")
        print(f"\nMerged rows:")
        for row in merged['rows']:
            print(f"  {row}")
        
        # Save merge log
        merger.save_merge_log()

