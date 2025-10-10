# Connect to your database
```bash
psql -h localhost -p 5432 -U postgres -d pdf_extraction
```


# Once connected, you can run these commands:
```bash
\dt                    # List all tables
\d table_name          # Describe a specific table
\q                     # Quit psql
```

-- List all tables
\dt

-- Show table structure
\d table_name

-- Count records in each table
```bash
SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'epbe_bases', COUNT(*) FROM epbe_bases
UNION ALL
SELECT 'depth_infos', COUNT(*) FROM depth_infos
UNION ALL
SELECT 'metadata_infos', COUNT(*) FROM metadata_infos
UNION ALL
SELECT 'petrography_carbonate', COUNT(*) FROM petrography_carbonate
UNION ALL
SELECT 'petrography_clastic', COUNT(*) FROM petrography_clastic;
```

-- View sample data from a table
```bash
SELECT * FROM epbe_bases LIMIT 5;
```

-- Show all columns and their types
```bash
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'epbe_bases';
```
