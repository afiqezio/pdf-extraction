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

	"workbench/internal/core/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
}

// ProcessPDF handles PDF upload and extraction
func (h *ExtractionHandler) ProcessPDF(c echo.Context) error {
	fmt.Printf("Processing PDF")
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
	inputDir := "../final_extraction_system/input_pdfs/"
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create input directory",
		})
	}

	fmt.Printf("Input directory: %s", inputDir)

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

	fmt.Println("Saving uploaded file to", filePath)

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

	fmt.Println("Starting run python extraction")

	// Run Python extraction script
	extractionResult, err := h.runPythonExtraction()
	if err != nil {
		// Clean up uploaded file on error
		os.Remove(filePath)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Extraction failed: %v", err),
		})
	}

	// Clean up uploaded file after processing
	os.Remove(filePath)

	// Return extraction results
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "PDF processed successfully",
		"results":      extractionResult,
		"filename":     file.Filename,
		"processed_at": time.Now().Format(time.RFC3339),
	})
}

// runPythonExtraction executes the Python extraction script
func (h *ExtractionHandler) runPythonExtraction() (map[string]interface{}, error) {
	log.Println("🚀 Starting Python extraction function")

	// Get the absolute path to the extraction system directory
	dir_temp := "../"
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
				"path":     path,
				"data":     jsonData,
				"size":     info.Size(),
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
				"path":     path,
				"content":  string(content),
				"size":     info.Size(),
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
		"id":      id,
		"status":  "completed",
		"message": "Extraction status endpoint - ready for async implementation",
	})
}

// DebugFiles returns debug information about JSON files
func (h *ExtractionHandler) DebugFiles(c echo.Context) error {
	dir_temp := "../"
	outputDir := filepath.Join(dir_temp, "final_extraction_system", "output", "markdown")

	files, err := filepath.Glob(filepath.Join(outputDir, "*.json"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"output_dir":  outputDir,
		"files_found": len(files),
		"files":       files,
	})
}

// GetLatestJson returns the most recent JSON file
func (h *ExtractionHandler) GetLatestJson(c echo.Context) error {
	outputDir := "../final_extraction_system/output/markdown"

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

	return c.JSON(http.StatusOK, map[string]interface{}{
		"filename": fileInfo.Name(),
		"path":     mostRecentFile,
		"data":     jsonData,
		"size":     fileInfo.Size(),
		"modified": fileInfo.ModTime().Format(time.RFC3339),
	})
}

// SaveToDatabase handles saving extracted tables to database
func (h *ExtractionHandler) SaveToDatabase(c echo.Context) error {
	log.Println("🚀 Starting save to database process")

	// Parse request body
	var request struct {
		Tables   []map[string]interface{} `json:"tables"`
		Filename string                   `json:"filename"`
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
		mappedData, err := h.mapTableToDatabaseFields(table)
		if err != nil {
			log.Printf("❌ Failed to map table %d: %v", i+1, err)
			continue
		}

		// Determine table type and save
		tableType := h.detectTableType(mappedData)
		log.Printf("🎯 Table %d detected as: %s", i+1, tableType)

		records, err := h.saveTableToDatabase(tableType, mappedData)
		if err != nil {
			log.Printf("❌ Failed to save table %d: %v", i+1, err)
			continue
		}

		savedTables++
		totalRecords += records
		log.Printf("✅ Table %d saved: %d records", i+1, records)
	}

	log.Printf("🎉 Save complete: %d tables, %d total records", savedTables, totalRecords)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":       true,
		"saved_tables":  savedTables,
		"total_records": totalRecords,
		"details":       fmt.Sprintf("Successfully saved %d tables with %d total records to database", savedTables, totalRecords),
	})
}

// mapTableToDatabaseFields maps table headers to database field names using fuzzy matching
func (h *ExtractionHandler) mapTableToDatabaseFields(table map[string]interface{}) (map[string]interface{}, error) {
	headers, ok1 := table["headers"].([]interface{})
	rows, ok2 := table["rows"].([]interface{})

	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid table format")
	}

	// Database field mappings for fuzzy matching
	fieldMappings := map[string]string{
		// Geographic fields
		"country":          "country",
		"region":           "region",
		"sub_region":       "sub_region",
		"business_regions": "business_regions",
		"basin":            "basin",
		"sub_basin":        "sub_basin",

		// Well fields
		"well":      "well_name_field_name",
		"well_name": "well_name_field_name",
		"field":     "well_name_field_name",
		"uwi":       "uwi",

		// Coordinates
		"lat":       "latitude",
		"latitude":  "latitude",
		"lon":       "longitude",
		"longitude": "longitude",
		"coord":     "latitude",

		// Formation
		"formation": "formation_name",
		"reservoir": "reservoir_name",
		"period":    "period",
		"epoch":     "epoch",
		"age":       "age",

		// Depth fields
		"depth":        "top_depth_mmddf",
		"top_depth":    "top_depth_mmddf",
		"bottom_depth": "bottom_depth_mmddf",
		"porosity":     "visible_porosity_percent",
		"permeability": "permeability_md",

		// Carbonate specific
		"calcite":  "calcite",
		"dolomite": "dolomite",
		"micrite":  "micrite",
		"bioclast": "bioclasts",
		"coral":    "coral",
		"algae":    "red_algae",

		// Clastic specific
		"quartz":     "monocrystalline_quartz",
		"feldspar":   "potassium_feldspar",
		"mica":       "muscovite",
		"grain_size": "grain_size",
		"sorting":    "sorting",
	}

	// Create mapping from user headers to database fields
	headerMapping := make(map[int]string)

	for i, header := range headers {
		headerStr := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", header)))

		// Try exact match first
		if dbField, exists := fieldMappings[headerStr]; exists {
			headerMapping[i] = dbField
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
		"rows":    rows,
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

// detectTableType determines if table is carbonate or clastic based on content
func (h *ExtractionHandler) detectTableType(mappedData map[string]interface{}) string {
	headers, _ := mappedData["headers"].([]interface{})

	// Check for carbonate-specific terms
	carbonateTerms := []string{"calcite", "dolomite", "micrite", "bioclast", "coral", "algae", "porosity"}
	clasticTerms := []string{"quartz", "feldspar", "mica", "grain", "sorting", "sandstone", "siltstone"}

	carbonateScore := 0
	clasticScore := 0

	for _, header := range headers {
		headerStr := strings.ToLower(fmt.Sprintf("%v", header))

		for _, term := range carbonateTerms {
			if strings.Contains(headerStr, term) {
				carbonateScore++
			}
		}

		for _, term := range clasticTerms {
			if strings.Contains(headerStr, term) {
				clasticScore++
			}
		}
	}

	if carbonateScore > clasticScore {
		return "petrography_carbonate"
	} else if clasticScore > carbonateScore {
		return "petrography_clastic"
	} else {
		// Default to carbonate if unclear
		return "petrography_carbonate"
	}
}

// saveTableToDatabase saves the mapped table data to the appropriate database table
func (h *ExtractionHandler) saveTableToDatabase(tableType string, mappedData map[string]interface{}) (int, error) {
	headers, _ := mappedData["headers"].([]interface{})
	rows, _ := mappedData["rows"].([]interface{})
	mapping, _ := mappedData["mapping"].(map[int]string)

	log.Printf("💾 Saving to table: %s", tableType)
	log.Printf("📋 Headers: %v", headers)
	log.Printf("🔗 Mapping: %v", mapping)
	log.Printf("📊 Rows: %d", len(rows))

	// Get database connection
	db := h.db

	// Insert records based on table type
	switch tableType {
	case "petrography_carbonate":
		return h.insertCarbonateRecords(db, headers, rows, mapping)
	case "petrography_clastic":
		return h.insertClasticRecords(db, headers, rows, mapping)
	default:
		log.Printf("⚠️ Unknown table type: %s", tableType)
		return 0, fmt.Errorf("unknown table type: %s", tableType)
	}
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
			case "well_name_field_name":
				carbonate.WellNameFieldName = cellStr
			case "country":
				carbonate.Country = cellStr
			case "region":
				carbonate.Region = cellStr
			case "basin":
				carbonate.Basin = cellStr
			case "top_depth_mmddf":
				if depth, err := strconv.ParseFloat(cellStr, 64); err == nil {
					carbonate.TopDepthMMDDF = &depth
				}
			case "bottom_depth_mmddf":
				if depth, err := strconv.ParseFloat(cellStr, 64); err == nil {
					carbonate.BottomDepthMMDDF = &depth
				}
			case "formation_name":
				carbonate.FormationName = cellStr
			case "reservoir_name":
				carbonate.ReservoirName = cellStr
			case "period":
				carbonate.Period = cellStr
			case "lithofacies_core":
				carbonate.LithofaciesCore = cellStr
			case "porosity":
				if porosity, err := strconv.ParseFloat(cellStr, 64); err == nil {
					carbonate.VisiblePorosityPercent = &porosity
				}
			case "calcite":
				if calcite, err := strconv.ParseFloat(cellStr, 64); err == nil {
					carbonate.Calcite = &calcite
				}
			case "dolomite":
				if dolomite, err := strconv.ParseFloat(cellStr, 64); err == nil {
					carbonate.Dolomite = &dolomite
				}
			case "micrite":
				if micrite, err := strconv.ParseFloat(cellStr, 64); err == nil {
					carbonate.Micrite = &micrite
				}
			case "bioclasts":
				if bioclasts, err := strconv.ParseFloat(cellStr, 64); err == nil {
					carbonate.Bioclasts = &bioclasts
				}
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
			case "well_name_field_name":
				clastic.WellNameFieldName = cellStr
			case "country":
				clastic.Country = cellStr
			case "region":
				clastic.Region = cellStr
			case "basin":
				clastic.Basin = cellStr
			case "top_depth_mmddf":
				if depth, err := strconv.ParseFloat(cellStr, 64); err == nil {
					clastic.TopDepthMMDDF = &depth
				}
			case "bottom_depth_mmddf":
				if depth, err := strconv.ParseFloat(cellStr, 64); err == nil {
					clastic.BottomDepthMMDDF = &depth
				}
			case "formation_name":
				clastic.FormationName = cellStr
			case "reservoir_name":
				clastic.ReservoirName = cellStr
			case "period":
				clastic.Period = cellStr
			case "lithofacies_core":
				clastic.Lithofacies = cellStr
			case "porosity":
				if porosity, err := strconv.ParseFloat(cellStr, 64); err == nil {
					clastic.VisiblePorosityPercent = &porosity
				}
			case "quartz":
				if quartz, err := strconv.ParseFloat(cellStr, 64); err == nil {
					clastic.TotalQuartzPercent = &quartz
				}
			case "feldspar":
				if feldspar, err := strconv.ParseFloat(cellStr, 64); err == nil {
					clastic.TotalFeldsparPercent = &feldspar
				}
			case "mica":
				if mica, err := strconv.ParseFloat(cellStr, 64); err == nil {
					clastic.TotalMicaPercent = &mica
				}
			case "grain_size":
				clastic.GrainSize = cellStr
			case "sorting":
				clastic.Sorting = cellStr
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
