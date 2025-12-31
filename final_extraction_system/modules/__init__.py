"""
PDF Extraction Modules

This package contains modules for intelligent PDF table processing:
- field_normalizer: Normalizes column names with fuzzy matching
- table_merger: Intelligently merges related tables
"""

from .field_normalizer import FieldNormalizer
from .table_merger import TableMerger

__all__ = ['FieldNormalizer', 'TableMerger']
__version__ = '1.0.0'












