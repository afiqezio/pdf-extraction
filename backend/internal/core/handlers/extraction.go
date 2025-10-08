package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

	fmt.Println("in run pytohn extraction function");
	// Get the absolute path to the extraction system directory
	// Use project root as extractionDir (assuming cwd is project root)
	extractionDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}
	dir_temp := "/Users/husaini.rosdi/Documents/working_project/pdf-extraction"
	// Use bash to activate the virtual environment and run the Python script
	activateScript := filepath.Join(dir_temp, "temp_env", "bin", "activate")
	pythonScriptPath := filepath.Join(dir_temp, "final_extraction_system", "better_markdown_extractor.py")
	
	// Run: source activate && python script
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && python %s", activateScript, pythonScriptPath))
	cmd.Dir = filepath.Join(extractionDir, "final_extraction_system")
	
	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("python script failed: %v, output: %s", err, string(output))
	}

	// Read the generated markdown files
	outputDir := filepath.Join(extractionDir, "output", "markdown")
	markdownFiles, err := h.readMarkdownFiles(outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read markdown files: %v", err)
	}

	return map[string]interface{}{
		"extraction_output": string(output),
		"markdown_files": markdownFiles,
		"files_count": len(markdownFiles),
	}, nil
}

// readMarkdownFiles reads all markdown files from the output directory
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
