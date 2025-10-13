package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"workbench/internal/core/models"
)

type ExtractionHandler struct {
	db *gorm.DB
}

func NewExtractionHandler(db *gorm.DB) *ExtractionHandler {
	return &ExtractionHandler{db: db}
}

func (h *ExtractionHandler) ExtractionRoutes(g *echo.Group) {
	extraction := g.Group("/extraction")
	extraction.POST("/process-pdf", h.ProcessPDF)
	extraction.POST("/save-to-db", h.SaveToDatabase)
	extraction.GET("/status/:id", h.GetExtractionStatus)
	extraction.GET("/debug", h.DebugFiles)
	extraction.GET("/latest-json", h.GetLatestJson)
	extraction.GET("/pdf/:filename", h.ServePDF)
	extraction.GET("/pdf/:filename/page/:page", h.ServePDFPage)
}

// ProcessPDF handles PDF upload and extraction
func (h *ExtractionHandler) ProcessPDF(c echo.Context) error {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No file uploaded",
		})
	}

	fmt.Printf("DEBUG: Received file upload: %s, size: %d\n", file.Filename, file.Size)

	// Validate file type
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".pdf") {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Only PDF files are supported",
		})
	}

	// Create input directory if it doesn't exist
	inputDir := "/Users/husaini.rosdi/Documents/working_project/pdf-extraction/final_extraction_system/input_pdfs"
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create input directory",
		})
	}

	// Generate unique filename to avoid conflicts
	timestamp := time.Now().Format("20060102_150405")
	uniqueFilename := fmt.Sprintf("%s_%s", timestamp, file.Filename)
	filePath := filepath.Join(inputDir, uniqueFilename)

	// Save uploaded file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create file",
		})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save file",
		})
	}

	fmt.Println("before run pytohn extraction");

	// Run Python extraction script
	extractionResult, err := h.runPythonExtraction()
	if err != nil {
		// Clean up uploaded file on error
		os.Remove(filePath)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Extraction failed: %v", err),
		})
	}

	// Keep the uploaded PDF file for the viewer (don't delete it)
	// os.Remove(filePath) // Commented out to keep PDF for viewer

	// Return extraction results
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "PDF processed successfully",
		"results": extractionResult,
		"filename": uniqueFilename, // Return the timestamped filename that's actually stored
		"original_filename": file.Filename, // Keep original for reference
		"processed_at": time.Now().Format(time.RFC3339),
	})
}

// runPythonExtraction executes the Python extraction script
func (h *ExtractionHandler) runPythonExtraction() (map[string]interface{}, error) {
	log.Println("🚀 Starting Python extraction function")
	
	// Get the absolute path to the extraction system directory
	dir_temp := "/Users/husaini.rosdi/Documents/working_project/pdf-extraction"
	log.Printf("📁 Working directory: %s", dir_temp)
	
	// Use bash to activate the virtual environment and run the Python script
	activateScript := filepath.Join(dir_temp, "temp_env", "bin", "activate")
	pythonScriptPath := filepath.Join(dir_temp, "final_extraction_system", "better_markdown_extractor.py")
	
	log.Printf("🐍 Running Python script: %s", pythonScriptPath)
	
	// Run: source activate && python script
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && python %s", activateScript, pythonScriptPath))
	cmd.Dir = filepath.Join(dir_temp, "final_extraction_system")
	
	// Explicitly redirect stdout and stderr for better debugging
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Run the command and wait for completion
	err := cmd.Run()
	if err != nil {
		log.Printf("❌ Python script error: %v", err)
		log.Printf("📤 Stdout: %s", stdout.String())
		log.Printf("📤 Stderr: %s", stderr.String())
		return nil, fmt.Errorf("python script failed: %v, stderr: %s", err, stderr.String())
	}
	
	output := stdout.String() + stderr.String()
	log.Printf("✅ Python script completed successfully")

	// Read the generated JSON files
	outputDir := filepath.Join(dir_temp, "final_extraction_system", "output", "markdown")
	log.Printf("📂 Looking for JSON files in: %s", outputDir)
	
	// Check if directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		log.Printf("❌ Output directory does not exist: %s", outputDir)
		return nil, fmt.Errorf("output directory does not exist: %s", outputDir)
	}
	
	// Use os.ReadDir instead of filepath.Glob for more reliable file reading
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		log.Printf("❌ Failed to read output directory: %v", err)
		return nil, fmt.Errorf("failed to read output directory: %v", err)
	}
	
	log.Printf("📊 Found %d entries in output directory", len(entries))
	
	// Filter and collect JSON files with their modification times
	type fileWithTime struct {
		path    string
		modTime time.Time
		info    os.FileInfo
	}
	var jsonFileList []fileWithTime
	
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			fullPath := filepath.Join(outputDir, entry.Name())
			info, err := os.Stat(fullPath)
			if err != nil {
				log.Printf("⚠️ Warning: could not stat file %s: %v", fullPath, err)
				continue
			}
			jsonFileList = append(jsonFileList, fileWithTime{
				path:    fullPath,
				modTime: info.ModTime(),
				info:    info,
			})
			log.Printf("📄 Found JSON file: %s (modified: %s)", entry.Name(), info.ModTime().Format(time.RFC3339))
		}
	}
	
	log.Printf("✅ Found %d JSON files total", len(jsonFileList))
	
	if len(jsonFileList) == 0 {
		log.Printf("⚠️ No JSON files found in output directory")
		return map[string]interface{}{
			"extraction_output": output,
			"json_files":        nil,
			"files_count":       0,
		}, nil
	}
	
	// Sort by modification time (most recent last)
	sort.Slice(jsonFileList, func(i, j int) bool {
		return jsonFileList[i].modTime.Before(jsonFileList[j].modTime)
	})
	
	// Get the most recent file
	mostRecent := jsonFileList[len(jsonFileList)-1]
	log.Printf("📖 Reading most recent JSON file: %s", mostRecent.info.Name())
	
	content, err := os.ReadFile(mostRecent.path)
	if err != nil {
		log.Printf("❌ Failed to read JSON file: %v", err)
		return nil, fmt.Errorf("failed to read most recent json file: %v", err)
	}
	
	var jsonData map[string]interface{}
	if err := json.Unmarshal(content, &jsonData); err != nil {
		log.Printf("❌ Failed to parse JSON: %v", err)
		return nil, fmt.Errorf("failed to parse json: %v", err)
	}
	
	log.Printf("✅ Successfully parsed JSON data from %s", mostRecent.info.Name())
	
	jsonFiles := []map[string]interface{}{
		{
			"filename": mostRecent.info.Name(),
			"path":     mostRecent.path,
			"data":     jsonData,
			"size":     mostRecent.info.Size(),
			"modified": mostRecent.info.ModTime().Format(time.RFC3339),
		},
	}
	
	log.Printf("🎉 Extraction complete! Returning %d JSON file(s)", len(jsonFiles))
	
	return map[string]interface{}{
		"extraction_output": output,
		"json_files":        jsonFiles,
		"files_count":       len(jsonFiles),
	}, nil
}

// readMarkdownFiles reads all markdown files from the output directory
func (h *ExtractionHandler) readJsonFiles(outputDir string) ([]map[string]interface{}, error) {
	var files []map[string]interface{}
	
	fmt.Printf("DEBUG: Looking for JSON files in: %s\n", outputDir)
	
	// Check if output directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		fmt.Printf("DEBUG: Output directory does not exist: %s\n", outputDir)
		return files, nil
	}

	// Read all .json files
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if strings.HasSuffix(strings.ToLower(path), ".json") {
			fmt.Printf("DEBUG: Found JSON file: %s\n", path)
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("DEBUG: Error reading file %s: %v\n", path, err)
				return err
			}
			
			// Parse JSON content
			var jsonData map[string]interface{}
			if err := json.Unmarshal(content, &jsonData); err != nil {
				fmt.Printf("DEBUG: Error parsing JSON from %s: %v\n", path, err)
				return err
			}
			
			fmt.Printf("DEBUG: Successfully parsed JSON from %s\n", path)
			files = append(files, map[string]interface{}{
				"filename": info.Name(),
				"path": path,
				"data": jsonData,
				"size": info.Size(),
				"modified": info.ModTime().Format(time.RFC3339),
			})
		}
		
		return nil
	})
	
	return files, err
}

func (h *ExtractionHandler) readMarkdownFiles(outputDir string) ([]map[string]interface{}, error) {
	var files []map[string]interface{}
	
	// Check if output directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return files, nil
	}

	// Read all .md files
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if strings.HasSuffix(strings.ToLower(path), ".md") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			
			files = append(files, map[string]interface{}{
				"filename": info.Name(),
				"path": path,
				"content": string(content),
				"size": info.Size(),
				"modified": info.ModTime().Format(time.RFC3339),
			})
		}
		
		return nil
	})
	
	return files, err
}

// GetExtractionStatus returns the status of an extraction (placeholder for future async processing)
func (h *ExtractionHandler) GetExtractionStatus(c echo.Context) error {
	id := c.Param("id")
	
	// For now, just return a simple status
	// In the future, this could check a database for extraction status
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
		"status": "completed",
		"message": "Extraction status endpoint - ready for async implementation",
	})
}

// DebugFiles returns debug information about JSON files
func (h *ExtractionHandler) DebugFiles(c echo.Context) error {
	dir_temp := "/Users/husaini.rosdi/Documents/working_project/pdf-extraction"
	outputDir := filepath.Join(dir_temp, "final_extraction_system", "output", "markdown")
	
	files, err := filepath.Glob(filepath.Join(outputDir, "*.json"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"output_dir": outputDir,
		"files_found": len(files),
		"files": files,
	})
}

// GetLatestJson returns the most recent JSON file
func (h *ExtractionHandler) GetLatestJson(c echo.Context) error {
	outputDir := "/Users/husaini.rosdi/Documents/working_project/pdf-extraction/final_extraction_system/output/markdown"
	
	files, err := filepath.Glob(filepath.Join(outputDir, "*.json"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	if len(files) == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "No JSON files found",
		})
	}
	
	// Get the most recent file
	mostRecentFile := files[len(files)-1]
	content, err := os.ReadFile(mostRecentFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	var jsonData map[string]interface{}
	if err := json.Unmarshal(content, &jsonData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	fileInfo, err := os.Stat(mostRecentFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	// Extract timestamped PDF filename from JSON filename
	// JSON filename format: "20251013_102142_PMO3ANDUPP03 - PETROGRAPHIC STUDY OF ANDING UTARA-1ST1 DITCH CUTTINGS_extracted.json"
	// Timestamped PDF filename: "20251013_102142_PMO3ANDUPP03 - PETROGRAPHIC STUDY OF ANDING UTARA-1ST1 DITCH CUTTINGS.pdf"
	jsonFilename := fileInfo.Name()
	timestampedPdfFilename := jsonFilename
	
	// Convert JSON filename to PDF filename
	if strings.Contains(jsonFilename, "_extracted.json") {
		timestampedPdfFilename = strings.Replace(jsonFilename, "_extracted.json", ".pdf", 1)
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"filename": timestampedPdfFilename,
		"path": mostRecentFile,
		"data": jsonData,
		"size": fileInfo.Size(),
		"modified": fileInfo.ModTime().Format(time.RFC3339),
	})
}

// ServePDF handles serving PDF files for the frontend viewer
func (h *ExtractionHandler) ServePDF(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Filename parameter is required",
		})
	}

	// Look for the PDF file in the input_pdfs directory (where uploaded PDFs are stored)
	// Backend runs from backend/ directory, so we need to go up one level
	inputPdfsPath := filepath.Join("..", "final_extraction_system", "input_pdfs", filename)
	
	// Check if file exists
	if _, err := os.Stat(inputPdfsPath); os.IsNotExist(err) {
		log.Printf("❌ PDF file not found: %s", inputPdfsPath)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "PDF file not found",
		})
	}
	
	pdfPath := inputPdfsPath

	// Set appropriate headers for PDF viewing
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	
	// Serve the file
	return c.File(pdfPath)
}

// ServePDFPage handles serving specific PDF pages for the frontend viewer
func (h *ExtractionHandler) ServePDFPage(c echo.Context) error {
	filename := c.Param("filename")
	pageStr := c.Param("page")
	
	if filename == "" || pageStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Filename and page parameters are required",
		})
	}

	// Convert page to integer
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid page number",
		})
	}

	// Look for the PDF file in the input_pdfs directory
	// Backend runs from backend/ directory, so we need to go up one level
	inputPdfsPath := filepath.Join("..", "final_extraction_system", "input_pdfs", filename)
	
	// Check if file exists
	if _, err := os.Stat(inputPdfsPath); os.IsNotExist(err) {
		log.Printf("❌ PDF file not found: %s", inputPdfsPath)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "PDF file not found",
		})
	}

	// For now, serve the full PDF (we can implement page extraction later)
	// Set appropriate headers for PDF viewing
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	
	// Serve the file
	return c.File(inputPdfsPath)
}

// SaveToDatabase handles saving extracted tables to database
func (h *ExtractionHandler) SaveToDatabase(c echo.Context) error {
	log.Println("🚀 Starting save to database process")
	
	// Parse request body
	var request struct {
		Tables   []map[string]interface{} `json:"tables"`
		Filename string                  `json:"filename"`
	}
	
	if err := c.Bind(&request); err != nil {
		log.Printf("❌ Failed to parse request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}
	
	log.Printf("📊 Received %d tables to save", len(request.Tables))
	
	// Filter out empty tables
	validTables := make([]map[string]interface{}, 0)
	for i, table := range request.Tables {
		headers, ok1 := table["headers"].([]interface{})
		rows, ok2 := table["rows"].([]interface{})
		
		if ok1 && ok2 && len(headers) > 0 && len(rows) > 0 {
			validTables = append(validTables, table)
			log.Printf("✅ Table %d: %d headers, %d rows", i+1, len(headers), len(rows))
		} else {
			log.Printf("⚠️ Skipping empty table %d", i+1)
		}
	}
	
	if len(validTables) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No valid tables to save",
		})
	}
	
	log.Printf("💾 Processing %d valid tables", len(validTables))
	
	// Process each table
	savedTables := 0
	totalRecords := 0
	
	for i, table := range validTables {
		log.Printf("🔄 Processing table %d/%d", i+1, len(validTables))
		
		// Map headers to database fields using fuzzy matching
		log.Printf("🔍 Mapping headers for table %d", i+1)
		headers, _ := table["headers"].([]interface{})
		log.Printf("📋 Original headers: %v", headers)
		
		mappedData, err := h.mapTableToDatabaseFields(table)
		if err != nil {
			log.Printf("❌ Failed to map table %d: %v", i+1, err)
			continue
		}
		
		log.Printf("🎯 Mapped data keys: %v", getMapKeys(mappedData))
		if mapping, ok := mappedData["mapping"].(map[int]string); ok {
			log.Printf("🔗 Column mapping: %v", mapping)
		}
		
		// Save to appropriate tables based on mapped fields
		records, err := h.saveTableToDatabase(mappedData)
		if err != nil {
			log.Printf("❌ Failed to save table %d: %v", i+1, err)
			continue
		}
		
		log.Printf("✅ Table %d saved: %d records", i+1, records)
		savedTables++
		totalRecords += records
	}
	
	log.Printf("🎉 Save complete: %d tables, %d total records", savedTables, totalRecords)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"saved_tables": savedTables,
		"total_records": totalRecords,
		"details": fmt.Sprintf("Successfully saved %d tables with %d total records to database", savedTables, totalRecords),
	})
}

// mapTableToDatabaseFields maps table headers to database field names using fuzzy matching
func (h *ExtractionHandler) mapTableToDatabaseFields(table map[string]interface{}) (map[string]interface{}, error) {
	headers, ok1 := table["headers"].([]interface{})
	rows, ok2 := table["rows"].([]interface{})
	
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid table format")
	}
	
	// Get field mappings from separate file
	fieldMappings := GetFieldMappings()
	
	// Create mapping from user headers to database fields
	headerMapping := make(map[int]string)
	
	for i, header := range headers {
		headerStr := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", header)))
		log.Printf("🔍 Processing header %d: '%s' -> '%s'", i, header, headerStr)
		
		// Try exact match first
		if dbField, exists := fieldMappings[headerStr]; exists {
			headerMapping[i] = dbField
			log.Printf("✅ Exact match found: '%s' -> '%s'", headerStr, dbField)
			continue
		}
		
		// Try fuzzy matching
		bestMatch := ""
		bestScore := 0.0
		
		for userField, dbField := range fieldMappings {
			score := h.calculateSimilarity(headerStr, userField)
			if score > 0.6 && score > bestScore { // 60% similarity threshold
				bestMatch = dbField
				bestScore = score
			}
		}
		
		if bestMatch != "" {
			headerMapping[i] = bestMatch
			log.Printf("🔗 Mapped '%s' -> '%s' (score: %.2f)", headerStr, bestMatch, bestScore)
		} else {
			log.Printf("⚠️ No mapping found for header: '%s'", headerStr)
		}
	}
	
	// Create mapped data structure
	mappedData := map[string]interface{}{
		"headers": headers,
		"rows": rows,
		"mapping": headerMapping,
	}
	
	return mappedData, nil
}

// calculateSimilarity calculates string similarity using simple algorithm
func (h *ExtractionHandler) calculateSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	
	// Simple similarity based on common substrings
	common := 0
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}
	
	for i := 0; i < len(s1) && i < len(s2); i++ {
		if s1[i] == s2[i] {
			common++
		}
	}
	
	return float64(common) / float64(maxLen)
}


// saveTableToDatabase saves the mapped table data to the appropriate database table
func (h *ExtractionHandler) saveTableToDatabase(mappedData map[string]interface{}) (int, error) {
	headers, _ := mappedData["headers"].([]interface{})
	rows, _ := mappedData["rows"].([]interface{})
	mapping, _ := mappedData["mapping"].(map[int]string)
	
	log.Printf("📋 Headers: %v", headers)
	log.Printf("🔗 Mapping: %v", mapping)
	log.Printf("📊 Rows: %d", len(rows))
	
	// Get database connection
	db := h.db
	
	// Determine which tables to save to based on mapped fields
	totalRecords := 0
	
	// Check if we have carbonate fields
	carbonateFields := h.getCarbonateFields(mapping)
	if len(carbonateFields) > 0 {
		log.Printf("💾 Saving to petrography_carbonate table with fields: %v", carbonateFields)
		records, err := h.insertCarbonateRecords(db, headers, rows, mapping)
		if err != nil {
			log.Printf("❌ Failed to save to carbonate table: %v", err)
		} else {
			totalRecords += records
			log.Printf("✅ Saved %d records to carbonate table", records)
		}
	}
	
	// Check if we have clastic fields
	clasticFields := h.getClasticFields(mapping)
	if len(clasticFields) > 0 {
		log.Printf("💾 Saving to petrography_clastic table with fields: %v", clasticFields)
		records, err := h.insertClasticRecords(db, headers, rows, mapping)
		if err != nil {
			log.Printf("❌ Failed to save to clastic table: %v", err)
		} else {
			totalRecords += records
			log.Printf("✅ Saved %d records to clastic table", records)
		}
	}
	
	if totalRecords == 0 {
		log.Printf("⚠️ No records saved - no matching fields found")
		return 0, fmt.Errorf("no matching fields found for any table")
	}
	
	return totalRecords, nil
}

// getCarbonateFields returns the fields that belong to the carbonate table
func (h *ExtractionHandler) getCarbonateFields(mapping map[int]string) []string {
	carbonateFields := []string{}
	for _, field := range mapping {
		if h.isCarbonateField(field) {
			carbonateFields = append(carbonateFields, field)
		}
	}
	return carbonateFields
}

// getClasticFields returns the fields that belong to the clastic table
func (h *ExtractionHandler) getClasticFields(mapping map[int]string) []string {
	clasticFields := []string{}
	for _, field := range mapping {
		if h.isClasticField(field) {
			clasticFields = append(clasticFields, field)
		}
	}
	return clasticFields
}

// isCarbonateField checks if a field belongs to the carbonate table
func (h *ExtractionHandler) isCarbonateField(field string) bool {
	// Common fields that exist in both tables
	commonFields := map[string]bool{
		"well_name_field_name": true, "country": true, "region": true, "sub_region": true,
		"business_regions": true, "basin": true, "sub_basin": true, "uwi": true,
		"latitude": true, "longitude": true, "formation_name": true, "reservoir_name": true,
		"period": true, "epoch": true, "age": true, "onshore_offshore": true,
		"water_depth_m": true, "water_depth_ft": true, "top_depth_mmddf": true,
		"bottom_depth_mmddf": true, "top_depth_mtvddf": true, "top_depth_mtvdss": true,
		"top_depth_mbml": true, "lithofacies_core": true, "microfacies_thin_section": true,
		"depofacies": true, "analysis_types": true, "visible_porosity_percent": true,
		"he_porosity_percent": true, "permeability_md": true, "grain_density_g_cc": true,
	}
	
	if commonFields[field] {
		return true
	}
	
	// Carbonate-specific fields
	carbonateSpecific := map[string]bool{
		"calcite": true, "dolomite": true, "micrite": true, "micrite_envelopes": true,
		"microspar_pseudospar": true, "kaolinite": true, "clay": true, "total_mineralogy_matrix_percent": true,
		"bioclasts": true, "lepido": true, "coral": true, "rhodolith": true, "red_algae": true,
		"red_algae_enc": true, "green_algae": true, "echinoderms": true, "miliolid": true,
		"cycloclypeus": true, "operculina": true, "other_rotaliids": true, "gypsinid": true,
		"planorbulinella": true, "hemotremid": true, "heterostegina": true, "enc_frm": true,
		"planktonic": true, "bryozoans": true, "amphistegina": true, "gastropods": true,
		"bivalve": true, "ostracod": true, "oncoids": true, "undiff_molluscs": true,
		"undiff_benthonic": true, "undiff_skeletal": true, "undiff_foram": true,
		"total_skeletal_percent": true, "organic": true, "peloids": true, "micritised_grains": true,
		"pseudoclasts": true, "intraclast": true, "quartz": true, "total_non_skeletal_percent": true,
		"interparticle": true, "intraparticle": true, "intercrystalline": true, "matrix_intercrystalline": true,
		"mouldic": true, "vuggy": true, "fractures": true, "micro": true, "total_porosity_percent": true,
		"fringing": true, "meniscus": true, "blocky": true, "sparry": true, "micritic": true,
		"pendant": true, "syntax": true, "calcite_syntaxial": true, "calcite_fringing": true,
		"calcite_mosaic": true, "calcite_blocky": true, "calcite_ferroan": true, "pyrite": true,
		"fluorite": true, "total_cement_percent": true, "replacement": true, "saddle": true,
		"total_dolomite_percent": true, "stylolite": true, "bioturbation": true,
		"total_accessories_percent": true, "total_percent": true,
	}
	
	return carbonateSpecific[field]
}

// isClasticField checks if a field belongs to the clastic table
func (h *ExtractionHandler) isClasticField(field string) bool {
	// Common fields that exist in both tables
	commonFields := map[string]bool{
		"well_name_field_name": true, "country": true, "region": true, "sub_region": true,
		"business_regions": true, "basin": true, "sub_basin": true, "uwi": true,
		"latitude": true, "longitude": true, "formation_name": true, "reservoir_name": true,
		"period": true, "epoch": true, "age": true, "onshore_offshore": true,
		"water_depth_m": true, "water_depth_ft": true, "top_depth_mmddf": true,
		"bottom_depth_mmddf": true, "top_depth_mtvddf": true, "top_depth_mtvdss": true,
		"top_depth_mbml": true, "lithofacies_core": true, "microfacies_thin_section": true,
		"depofacies": true, "analysis_types": true, "visible_porosity_percent": true,
		"he_porosity_percent": true, "permeability_md": true, "grain_density_g_cc": true,
	}
	
	if commonFields[field] {
		return true
	}
	
	// Clastic-specific fields
	clasticSpecific := map[string]bool{
		"grain_size": true, "grain_shape": true, "grain_contact": true, "sedimentary_structure": true,
		"sorting": true, "ambient_he_porosity_percent": true, "monocrystalline_quartz": true,
		"polycrystalline_quartz": true, "total_quartz_percent": true, "potassium_feldspar": true,
		"plagioclase": true, "feldspar_undifferentiated": true, "total_feldspar_percent": true,
		"muscovite": true, "biotite": true, "mica_undifferentiated": true, "total_mica_percent": true,
		"zircon": true, "tourmaline": true, "heavy_minerals_undifferentiated": true,
		"total_heavy_minerals_percent": true, "plutonic_rock_fragments": true,
		"mafic_intermediate_volcanic_fragment": true, "volcanic_rock_fragment": true,
		"total_igneous_rf_percent": true, "quartzose_rock_fragment": true,
		"schistose_rock_fragment": true, "metamorphic_rock_fragment_undifferentiated": true,
		"total_metamorphic_rf_percent": true, "sandstone_siltstone_rock_fragments": true,
		"argillaceous_rock_fragments": true, "siliciclastic_rock_fragments_undifferentiated": true,
		"limestone_rock_fragments": true, "dolostone_rock_fragments": true, "chert": true,
		"total_sedimentary_rf_percent": true, "total_rock_fragments_percent": true,
		"rip_up_clast": true, "glauconite": true, "foraminifera_grains": true,
		"undifferentiated_other_grains": true, "total_other_grains_percent": true,
		"clay_matrix": true, "mixed_clay_silt_fine_matrix": true, "silt_very_fine_matrix": true,
		"organic_matrix": true, "matrix_undifferentiated": true, "kaolinite_replaces_k_feldspar": true,
		"illite_pore_grain_lining": true, "illite_pore_filling": true, "illite_replaces_k_feldspar": true,
		"total_authigenic_clay_percent": true, "syntaxial_quartz_overgrowths": true,
		"feldspar_overgrowths": true, "fe_calcite": true, "fe_dolomite": true, "siderite": true,
		"mn_siderite": true, "iron_oxide_minerals": true, "total_authigenic_non_clay_percent": true,
		"intergranular": true, "pri_porosity_intragranular": true, "total_primary_porosity_percent": true,
		"sec_porosity_intragranular": true, "intracrystalline": true, "total_secondary_porosity_percent": true,
	}
	
	return clasticSpecific[field]
}

// Helper function to convert string to float64 pointer
func stringToFloat64Ptr(s string) *float64 {
	if s == "" {
		return nil
	}
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return &val
	}
	return nil
}

// Helper function to convert string to int pointer
func stringToIntPtr(s string) *int {
	if s == "" {
		return nil
	}
	if val, err := strconv.Atoi(s); err == nil {
		return &val
	}
	return nil
}

// Helper function to get map keys for debugging
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// insertCarbonateRecords inserts data into the petrography_carbonate table
func (h *ExtractionHandler) insertCarbonateRecords(db *gorm.DB, headers []interface{}, rows []interface{}, mapping map[int]string) (int, error) {
	recordCount := 0
	
	for _, row := range rows {
		rowSlice, ok := row.([]interface{})
		if !ok {
			continue
		}
		
		// Check if row has data
		hasData := false
		for _, cell := range rowSlice {
			if str, ok := cell.(string); ok && str != "" {
				hasData = true
				break
			}
		}
		if !hasData {
			continue
		}
		
		// Create carbonate record
		carbonate := models.EPBEPetrographyCarbonate{}
		
		// Map data to struct fields
		log.Printf("🔍 Processing carbonate row %d with %d cells", recordCount+1, len(rowSlice))
		for colIndex, cell := range rowSlice {
			cellStr, ok := cell.(string)
			if !ok {
				log.Printf("⚠️ Cell %d is not a string: %v", colIndex, cell)
				continue
			}
			
			fieldName, exists := mapping[colIndex]
			if !exists {
				log.Printf("⚠️ No mapping for column %d (value: %s)", colIndex, cellStr)
				continue
			}
			
			log.Printf("📝 Mapping column %d: '%s' -> %s", colIndex, cellStr, fieldName)
			
			// Map field names to struct fields
			switch fieldName {
			// String fields
			case "well_name_field_name":
				carbonate.WellNameFieldName = cellStr
			case "country":
				carbonate.Country = cellStr
			case "region":
				carbonate.Region = cellStr
			case "sub_region":
				carbonate.SubRegion = cellStr
			case "business_regions":
				carbonate.BusinessRegions = cellStr
			case "basin":
				carbonate.Basin = cellStr
			case "sub_basin":
				carbonate.SubBasin = cellStr
			case "uwi":
				carbonate.UWI = cellStr
			case "formation_name":
				carbonate.FormationName = cellStr
			case "reservoir_name":
				carbonate.ReservoirName = cellStr
			case "period":
				carbonate.Period = cellStr
			case "epoch":
				carbonate.Epoch = cellStr
			case "age":
				carbonate.Age = cellStr
			case "onshore_offshore":
				carbonate.OnshoreOffshore = cellStr
			case "lithofacies_core":
				carbonate.LithofaciesCore = cellStr
			case "microfacies_thin_section":
				carbonate.MicrofaciesThinSection = cellStr
			case "depofacies":
				carbonate.Depofacies = cellStr
			case "analysis_types":
				carbonate.AnalysisTypes = cellStr
			
			// Float64 fields
			case "latitude":
				carbonate.Latitude = stringToFloat64Ptr(cellStr)
			case "longitude":
				carbonate.Longitude = stringToFloat64Ptr(cellStr)
			case "water_depth_m":
				carbonate.WaterDepthM = stringToFloat64Ptr(cellStr)
			case "water_depth_ft":
				carbonate.WaterDepthFt = stringToFloat64Ptr(cellStr)
			case "top_depth_mmddf":
				carbonate.TopDepthMMDDF = stringToFloat64Ptr(cellStr)
			case "top_depth_mtvddf":
				carbonate.TopDepthMTVDDF = stringToFloat64Ptr(cellStr)
			case "top_depth_mtvdss":
				carbonate.TopDepthMTVDSS = stringToFloat64Ptr(cellStr)
			case "top_depth_mbml":
				carbonate.TopDepthMBML = stringToFloat64Ptr(cellStr)
			case "bottom_depth_mmddf":
				carbonate.BottomDepthMMDDF = stringToFloat64Ptr(cellStr)
			case "bottom_depth_mtvddf":
				carbonate.BottomDepthMTVDDF = stringToFloat64Ptr(cellStr)
			case "bottom_depth_mtvdss":
				carbonate.BottomDepthMTVDSS = stringToFloat64Ptr(cellStr)
			case "bottom_depth_mbml":
				carbonate.BottomDepthMBML = stringToFloat64Ptr(cellStr)
			case "top_depth_ftmddf":
				carbonate.TopDepthFtMDDF = stringToFloat64Ptr(cellStr)
			case "top_depth_fttvddf":
				carbonate.TopDepthFtTVDDF = stringToFloat64Ptr(cellStr)
			case "top_depth_fttvdss":
				carbonate.TopDepthFtTVDSS = stringToFloat64Ptr(cellStr)
			case "top_depth_ftbml":
				carbonate.TopDepthFtBML = stringToFloat64Ptr(cellStr)
			case "bottom_depth_ftmddf":
				carbonate.BottomDepthFtMDDF = stringToFloat64Ptr(cellStr)
			case "bottom_depth_fttvddf":
				carbonate.BottomDepthFtTVDDF = stringToFloat64Ptr(cellStr)
			case "bottom_depth_fttvdss":
				carbonate.BottomDepthFtTVDSS = stringToFloat64Ptr(cellStr)
			case "bottom_depth_ftbml":
				carbonate.BottomDepthFtBML = stringToFloat64Ptr(cellStr)
			case "visible_porosity_percent":
				carbonate.VisiblePorosityPercent = stringToFloat64Ptr(cellStr)
			case "he_porosity_percent":
				carbonate.HePorosityPercent = stringToFloat64Ptr(cellStr)
			case "permeability_md":
				carbonate.PermeabilityMd = stringToFloat64Ptr(cellStr)
			
			// Matrix mineralogy
			case "calcite":
				carbonate.Calcite = stringToFloat64Ptr(cellStr)
			case "dolomite":
				carbonate.Dolomite = stringToFloat64Ptr(cellStr)
			case "micrite":
				carbonate.Micrite = stringToFloat64Ptr(cellStr)
			case "micrite_envelopes":
				carbonate.MicriteEnvelopes = stringToFloat64Ptr(cellStr)
			case "microspar_pseudospar":
				carbonate.MicrosparPseudospar = stringToFloat64Ptr(cellStr)
			case "kaolinite":
				carbonate.Kaolinite = stringToFloat64Ptr(cellStr)
			case "clay":
				carbonate.Clay = stringToFloat64Ptr(cellStr)
			case "total_mineralogy_matrix_percent":
				carbonate.TotalMineralogyMatrixPercent = stringToFloat64Ptr(cellStr)
			
			// Bioclasts
			case "bioclasts":
				carbonate.Bioclasts = stringToFloat64Ptr(cellStr)
			case "lepido":
				carbonate.Lepido = stringToFloat64Ptr(cellStr)
			case "coral":
				carbonate.Coral = stringToFloat64Ptr(cellStr)
			case "rhodolith":
				carbonate.Rhodolith = stringToFloat64Ptr(cellStr)
			case "red_algae":
				carbonate.RedAlgae = stringToFloat64Ptr(cellStr)
			case "red_algae_enc":
				carbonate.RedAlgaeEnc = stringToFloat64Ptr(cellStr)
			case "green_algae":
				carbonate.GreenAlgae = stringToFloat64Ptr(cellStr)
			case "echinoderms":
				carbonate.Echinoderms = stringToFloat64Ptr(cellStr)
			case "miliolid":
				carbonate.Miliolid = stringToFloat64Ptr(cellStr)
			case "lepidocyclina":
				carbonate.Lepidocyclina = stringToFloat64Ptr(cellStr)
			case "cycloclypeus":
				carbonate.Cycloclypeus = stringToFloat64Ptr(cellStr)
			case "operculina":
				carbonate.Operculina = stringToFloat64Ptr(cellStr)
			case "other_rotaliids":
				carbonate.OtherRotaliids = stringToFloat64Ptr(cellStr)
			case "gypsinid":
				carbonate.Gypsinid = stringToFloat64Ptr(cellStr)
			case "planorbulinella":
				carbonate.Planorbulinella = stringToFloat64Ptr(cellStr)
			case "hemotremid":
				carbonate.Hemotremid = stringToFloat64Ptr(cellStr)
			case "heterostegina":
				carbonate.Heterostegina = stringToFloat64Ptr(cellStr)
			case "enc_frm":
				carbonate.EncFrm = stringToFloat64Ptr(cellStr)
			case "planktonic":
				carbonate.Planktonic = stringToFloat64Ptr(cellStr)
			case "bryozoans":
				carbonate.Bryozoans = stringToFloat64Ptr(cellStr)
			case "amphistegina":
				carbonate.Amphistegina = stringToFloat64Ptr(cellStr)
			case "gastropods":
				carbonate.Gastropods = stringToFloat64Ptr(cellStr)
			case "bivalve":
				carbonate.Bivalve = stringToFloat64Ptr(cellStr)
			case "ostracod":
				carbonate.Ostracod = stringToFloat64Ptr(cellStr)
			case "oncoids":
				carbonate.Oncoids = stringToFloat64Ptr(cellStr)
			case "undiff_molluscs":
				carbonate.UndiffMolluscs = stringToFloat64Ptr(cellStr)
			case "undiff_benthonic":
				carbonate.UndiffBenthonic = stringToFloat64Ptr(cellStr)
			case "undiff_skeletal":
				carbonate.UndiffSkeletal = stringToFloat64Ptr(cellStr)
			case "undiff_foram":
				carbonate.UndiffForam = stringToFloat64Ptr(cellStr)
			case "total_skeletal_percent":
				carbonate.TotalSkeletalPercent = stringToFloat64Ptr(cellStr)
			
			// Non-skeletal components
			case "organic":
				carbonate.Organic = stringToFloat64Ptr(cellStr)
			case "peloids":
				carbonate.Peloids = stringToFloat64Ptr(cellStr)
			case "micritised_grains":
				carbonate.MicritisedGrains = stringToFloat64Ptr(cellStr)
			case "pseudoclasts":
				carbonate.Pseudoclasts = stringToFloat64Ptr(cellStr)
			case "intraclast":
				carbonate.Intraclast = stringToFloat64Ptr(cellStr)
			case "quartz":
				carbonate.Quartz = stringToFloat64Ptr(cellStr)
			case "total_non_skeletal_percent":
				carbonate.TotalNonSkeletalPercent = stringToFloat64Ptr(cellStr)
			
			// Porosity types
			case "interparticle":
				carbonate.Interparticle = stringToFloat64Ptr(cellStr)
			case "intraparticle":
				carbonate.Intraparticle = stringToFloat64Ptr(cellStr)
			case "intercrystalline":
				carbonate.Intercrystalline = stringToFloat64Ptr(cellStr)
			case "matrix_intercrystalline":
				carbonate.MatrixIntercrystalline = stringToFloat64Ptr(cellStr)
			case "mouldic":
				carbonate.Mouldic = stringToFloat64Ptr(cellStr)
			case "vuggy":
				carbonate.Vuggy = stringToFloat64Ptr(cellStr)
			case "fractures":
				carbonate.Fractures = stringToFloat64Ptr(cellStr)
			case "micro":
				carbonate.Micro = stringToFloat64Ptr(cellStr)
			case "total_porosity_percent":
				carbonate.TotalPorosityPercent = stringToFloat64Ptr(cellStr)
			
			// Cement types
			case "fringing":
				carbonate.Fringing = stringToFloat64Ptr(cellStr)
			case "meniscus":
				carbonate.Meniscus = stringToFloat64Ptr(cellStr)
			case "blocky":
				carbonate.Blocky = stringToFloat64Ptr(cellStr)
			case "sparry":
				carbonate.Sparry = stringToFloat64Ptr(cellStr)
			case "micritic":
				carbonate.Micritic = stringToFloat64Ptr(cellStr)
			case "pendant":
				carbonate.Pendant = stringToFloat64Ptr(cellStr)
			case "syntax":
				carbonate.Syntax = stringToFloat64Ptr(cellStr)
			case "calcite_syntaxial":
				carbonate.CalciteSyntaxial = stringToFloat64Ptr(cellStr)
			case "calcite_fringing":
				carbonate.CalciteFringing = stringToFloat64Ptr(cellStr)
			case "calcite_mosaic":
				carbonate.CalciteMosaic = stringToFloat64Ptr(cellStr)
			case "calcite_blocky":
				carbonate.CalciteBlocky = stringToFloat64Ptr(cellStr)
			case "calcite_ferroan":
				carbonate.CalciteFerroan = stringToFloat64Ptr(cellStr)
			case "pyrite":
				carbonate.Pyrite = stringToFloat64Ptr(cellStr)
			case "fluorite":
				carbonate.Fluorite = stringToFloat64Ptr(cellStr)
			case "total_cement_percent":
				carbonate.TotalCementPercent = stringToFloat64Ptr(cellStr)
			
			// Replacement and accessories
			case "replacement":
				carbonate.Replacement = stringToFloat64Ptr(cellStr)
			case "saddle":
				carbonate.Saddle = stringToFloat64Ptr(cellStr)
			case "total_dolomite_percent":
				carbonate.TotalDolomitePercent = stringToFloat64Ptr(cellStr)
			case "stylolite":
				carbonate.Stylolite = stringToFloat64Ptr(cellStr)
			case "bioturbation":
				carbonate.Bioturbation = stringToFloat64Ptr(cellStr)
			case "total_accessories_percent":
				carbonate.TotalAccessoriesPercent = stringToFloat64Ptr(cellStr)
			case "total_percent":
				carbonate.TotalPercent = stringToFloat64Ptr(cellStr)
			}
		}
		
		// Insert record
		if err := db.Create(&carbonate).Error; err != nil {
			log.Printf("❌ Failed to insert carbonate record: %v", err)
			continue
		}
		
		recordCount++
	}
	
	log.Printf("✅ Inserted %d carbonate records", recordCount)
	return recordCount, nil
}

// insertClasticRecords inserts data into the petrography_clastic table
func (h *ExtractionHandler) insertClasticRecords(db *gorm.DB, headers []interface{}, rows []interface{}, mapping map[int]string) (int, error) {
	recordCount := 0
	
	for _, row := range rows {
		rowSlice, ok := row.([]interface{})
		if !ok {
			continue
		}
		
		// Check if row has data
		hasData := false
		for _, cell := range rowSlice {
			if str, ok := cell.(string); ok && str != "" {
				hasData = true
				break
			}
		}
		if !hasData {
			continue
		}
		
		// Create clastic record
		clastic := models.EPBEPetrographyClastic{}
		
		// Map data to struct fields
		for colIndex, cell := range rowSlice {
			cellStr, ok := cell.(string)
			if !ok {
				continue
			}
			
			fieldName, exists := mapping[colIndex]
			if !exists {
				continue
			}
			
			// Map field names to struct fields
			switch fieldName {
			// String fields
			case "well_name_field_name":
				clastic.WellNameFieldName = cellStr
			case "country":
				clastic.Country = cellStr
			case "region":
				clastic.Region = cellStr
			case "sub_region":
				clastic.SubRegion = cellStr
			case "business_regions":
				clastic.BusinessRegions = cellStr
			case "basin":
				clastic.Basin = cellStr
			case "sub_basin":
				clastic.SubBasin = cellStr
			case "uwi":
				clastic.UWI = cellStr
			case "formation_name":
				clastic.FormationName = cellStr
			case "reservoir_name":
				clastic.ReservoirName = cellStr
			case "period":
				clastic.Period = cellStr
			case "epoch":
				clastic.Epoch = cellStr
			case "age":
				clastic.Age = cellStr
			case "onshore_offshore":
				clastic.OnshoreOffshore = cellStr
			case "lithofacies_core":
				clastic.Lithofacies = cellStr
			case "analysis_types":
				clastic.AnalysisTypes = cellStr
			case "grain_size":
				clastic.GrainSize = cellStr
			case "grain_shape":
				clastic.GrainShape = cellStr
			case "grain_contact":
				clastic.GrainContact = cellStr
			case "sedimentary_structure":
				clastic.SedimentaryStructure = cellStr
			case "sorting":
				clastic.Sorting = cellStr
			
			// Float64 fields
			case "latitude":
				clastic.Latitude = stringToFloat64Ptr(cellStr)
			case "longitude":
				clastic.Longitude = stringToFloat64Ptr(cellStr)
			case "water_depth_m":
				clastic.WaterDepthM = stringToFloat64Ptr(cellStr)
			case "water_depth_ft":
				clastic.WaterDepthFt = stringToFloat64Ptr(cellStr)
			case "top_depth_mmddf":
				clastic.TopDepthMMDDF = stringToFloat64Ptr(cellStr)
			case "top_depth_mtvddf":
				clastic.TopDepthMTVDDF = stringToFloat64Ptr(cellStr)
			case "top_depth_mtvdss":
				clastic.TopDepthMTVDSS = stringToFloat64Ptr(cellStr)
			case "top_depth_mbml":
				clastic.TopDepthMBML = stringToFloat64Ptr(cellStr)
			case "bottom_depth_mmddf":
				clastic.BottomDepthMMDDF = stringToFloat64Ptr(cellStr)
			case "bottom_depth_mtvddf":
				clastic.BottomDepthMTVDDF = stringToFloat64Ptr(cellStr)
			case "bottom_depth_mtvdss":
				clastic.BottomDepthMTVDSS = stringToFloat64Ptr(cellStr)
			case "bottom_depth_mbml":
				clastic.BottomDepthMBML = stringToFloat64Ptr(cellStr)
			case "top_depth_ftmddf":
				clastic.TopDepthFtMDDF = stringToFloat64Ptr(cellStr)
			case "top_depth_fttvddf":
				clastic.TopDepthFtTVDDF = stringToFloat64Ptr(cellStr)
			case "top_depth_fttvdss":
				clastic.TopDepthFtTVDSS = stringToFloat64Ptr(cellStr)
			case "top_depth_ftbml":
				clastic.TopDepthFtBML = stringToFloat64Ptr(cellStr)
			case "bottom_depth_ftmddf":
				clastic.BottomDepthFtMDDF = stringToFloat64Ptr(cellStr)
			case "bottom_depth_fttvddf":
				clastic.BottomDepthFtTVDDF = stringToFloat64Ptr(cellStr)
			case "bottom_depth_fttvdss":
				clastic.BottomDepthFtTVDSS = stringToFloat64Ptr(cellStr)
			case "bottom_depth_ftbml":
				clastic.BottomDepthFtBML = stringToFloat64Ptr(cellStr)
			case "visible_porosity_percent":
				clastic.VisiblePorosityPercent = stringToFloat64Ptr(cellStr)
			case "he_porosity_percent":
				clastic.AmbientHePorosityPercent = stringToFloat64Ptr(cellStr)
			case "permeability_md":
				clastic.PermeabilityMd = stringToFloat64Ptr(cellStr)
			
			// Clastic mineralogy - Quartz
			case "monocrystalline_quartz":
				clastic.MonocrystallineQuartz = stringToFloat64Ptr(cellStr)
			case "polycrystalline_quartz":
				clastic.PolycrystallineQuartz = stringToFloat64Ptr(cellStr)
			case "quartz":
				clastic.TotalQuartzPercent = stringToFloat64Ptr(cellStr)
			
			// Feldspar
			case "potassium_feldspar":
				clastic.PotassiumFeldspar = stringToFloat64Ptr(cellStr)
			case "plagioclase":
				clastic.Plagioclase = stringToFloat64Ptr(cellStr)
			case "feldspar_undifferentiated":
				clastic.FeldsparUndifferentiated = stringToFloat64Ptr(cellStr)
			case "feldspar":
				clastic.TotalFeldsparPercent = stringToFloat64Ptr(cellStr)
			
			// Mica
			case "muscovite":
				clastic.Muscovite = stringToFloat64Ptr(cellStr)
			case "biotite":
				clastic.Biotite = stringToFloat64Ptr(cellStr)
			case "mica_undifferentiated":
				clastic.MicaUndifferentiated = stringToFloat64Ptr(cellStr)
			case "mica":
				clastic.TotalMicaPercent = stringToFloat64Ptr(cellStr)
			
			// Heavy Minerals
			case "zircon":
				clastic.Zircon = stringToFloat64Ptr(cellStr)
			case "tourmaline":
				clastic.Tourmaline = stringToFloat64Ptr(cellStr)
			case "heavy_minerals_undifferentiated":
				clastic.HeavyMineralsUndifferentiated = stringToFloat64Ptr(cellStr)
			case "total_heavy_minerals_percent":
				clastic.TotalHeavyMineralsPercent = stringToFloat64Ptr(cellStr)
			
			// Rock Fragments
			case "plutonic_rock_fragments":
				clastic.PlutonicRockFragments = stringToFloat64Ptr(cellStr)
			case "mafic_intermediate_volcanic_fragment":
				clastic.MaficIntermediateVolcanicFragment = stringToFloat64Ptr(cellStr)
			case "volcanic_rock_fragment":
				clastic.VolcanicRockFragment = stringToFloat64Ptr(cellStr)
			case "total_igneous_rf_percent":
				clastic.TotalIgneousRFPercent = stringToFloat64Ptr(cellStr)
			
			// Sedimentary Rock Fragments
			case "sandstone_siltstone_rock_fragments":
				clastic.SandstoneSiltstoneRockFragments = stringToFloat64Ptr(cellStr)
			case "argillaceous_rock_fragments":
				clastic.ArgillaceousRockFragments = stringToFloat64Ptr(cellStr)
			case "siliciclastic_rock_fragments_undifferentiated":
				clastic.SiliciclasticRockFragmentsUndifferentiated = stringToFloat64Ptr(cellStr)
			case "limestone_rock_fragments":
				clastic.LimestoneRockFragments = stringToFloat64Ptr(cellStr)
			case "dolostone_rock_fragments":
				clastic.DolostoneRockFragments = stringToFloat64Ptr(cellStr)
			case "chert":
				clastic.Chert = stringToFloat64Ptr(cellStr)
			case "total_sedimentary_rf_percent":
				clastic.TotalSedimentaryRFPercent = stringToFloat64Ptr(cellStr)
			case "total_rock_fragments_percent":
				clastic.TotalRockFragmentsPercent = stringToFloat64Ptr(cellStr)
			
			// Matrix
			case "clay_matrix":
				clastic.ClayMatrix = stringToFloat64Ptr(cellStr)
			case "mixed_clay_silt_fine_matrix":
				clastic.MixedClaySiltFineMatrix = stringToFloat64Ptr(cellStr)
			case "silt_very_fine_matrix":
				clastic.SiltVeryFineMatrix = stringToFloat64Ptr(cellStr)
			case "organic_matrix":
				clastic.OrganicMatrix = stringToFloat64Ptr(cellStr)
			case "matrix_undifferentiated":
				clastic.MatrixUndifferentiated = stringToFloat64Ptr(cellStr)
			case "total_matrix_percent":
				clastic.TotalMatrixPercent = stringToFloat64Ptr(cellStr)
			
			// Authigenic Clay
			case "kaolinite":
				clastic.Kaolinite = stringToFloat64Ptr(cellStr)
			case "kaolinite_replaces_k_feldspar":
				clastic.KaoliniteReplacesKFeldspar = stringToFloat64Ptr(cellStr)
			case "illite_pore_grain_lining":
				clastic.IllitePoreGrainLining = stringToFloat64Ptr(cellStr)
			case "illite_pore_filling":
				clastic.IllitePoreFilling = stringToFloat64Ptr(cellStr)
			case "illite_replaces_k_feldspar":
				clastic.IlliteReplacesKFeldspar = stringToFloat64Ptr(cellStr)
			case "total_authigenic_clay_percent":
				clastic.TotalAuthigenicClayPercent = stringToFloat64Ptr(cellStr)
			
			// Porosity
			case "intergranular":
				clastic.Intergranular = stringToFloat64Ptr(cellStr)
			case "intercrystalline":
				clastic.Intercrystalline = stringToFloat64Ptr(cellStr)
			case "mouldic":
				clastic.Mouldic = stringToFloat64Ptr(cellStr)
			case "fracture":
				clastic.Fracture = stringToFloat64Ptr(cellStr)
			case "total_primary_porosity_percent":
				clastic.TotalPrimaryPorosityPercent = stringToFloat64Ptr(cellStr)
			case "total_secondary_porosity_percent":
				clastic.TotalSecondaryPorosityPercent = stringToFloat64Ptr(cellStr)
			
			// Other
			case "pyrite":
				clastic.Pyrite = stringToFloat64Ptr(cellStr)
			case "bioclast":
				clastic.Bioclast = stringToFloat64Ptr(cellStr)
			case "total_percent":
				clastic.TotalPercent = stringToFloat64Ptr(cellStr)
			}
		}
		
		// Insert record
		if err := db.Create(&clastic).Error; err != nil {
			log.Printf("❌ Failed to insert clastic record: %v", err)
			continue
		}
		
		recordCount++
	}
	
	log.Printf("✅ Inserted %d clastic records", recordCount)
	return recordCount, nil
}
