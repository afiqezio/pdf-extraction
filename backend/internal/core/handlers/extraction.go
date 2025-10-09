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
	"strings"
	"time"

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
	extraction.GET("/status/:id", h.GetExtractionStatus)
	extraction.GET("/debug", h.DebugFiles)
	extraction.GET("/latest-json", h.GetLatestJson)
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

	// Clean up uploaded file after processing
	os.Remove(filePath)

	// Return extraction results
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "PDF processed successfully",
		"results": extractionResult,
		"filename": file.Filename,
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
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"filename": fileInfo.Name(),
		"path": mostRecentFile,
		"data": jsonData,
		"size": fileInfo.Size(),
		"modified": fileInfo.ModTime().Format(time.RFC3339),
	})
}
