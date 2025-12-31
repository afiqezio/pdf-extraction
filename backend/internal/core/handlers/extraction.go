package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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
	db *gorm.DB // kept for compatibility, not used now
}

func NewExtractionHandler(db *gorm.DB) *ExtractionHandler {
	return &ExtractionHandler{db: db}
}

// ExtractionRoutes registers endpoints under /extraction
func (h *ExtractionHandler) ExtractionRoutes(g *echo.Group) {
	extraction := g.Group("/extraction")
	// New async endpoints
	extraction.POST("/upload", h.UploadForExtraction)
	extraction.GET("/jobs/:id", h.GetExtractionJob)

	// Legacy synchronous endpoint
	extraction.POST("/process-pdf", h.ProcessPDF)
	extraction.POST("/save-qc", h.SaveQC)
	extraction.GET("/pdf/:filename", h.ServePDF)
}

const PythonServiceURL = "http://localhost:8000"

type PythonExtractRequest struct {
	FilePath string `json:"file_path"`
}

type PythonExtractResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

// UploadForExtraction handles file upload and initiates extraction via Python API
func (h *ExtractionHandler) UploadForExtraction(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}

	// Create input directory if it doesn't exist
	// Use absolute path to ensure Python service can find it
	cwd, _ := os.Getwd()
	// Assuming running from backend/, go up one level
	projectRoot := filepath.Dir(cwd)
	inputDir := filepath.Join(projectRoot, "final_extraction_system", "input_pdfs")

	if err := os.MkdirAll(inputDir, 0o755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create input directory"})
	}

	timestamp := time.Now().Format("20060102_150405")
	uniqueFilename := fmt.Sprintf("%s_%s", timestamp, file.Filename)
	filePath := filepath.Join(inputDir, uniqueFilename)

	if err := saveUploadedFile(file, filePath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to save file: %v", err)})
	}

	// Call Python Service
	reqBody := PythonExtractRequest{
		FilePath: filePath,
	}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post(PythonServiceURL+"/extract", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "Failed to connect to extraction service"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return c.JSON(resp.StatusCode, map[string]string{"error": "Extraction service error: " + string(body)})
	}

	var extractResp PythonExtractResponse
	if err := json.NewDecoder(resp.Body).Decode(&extractResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse extraction service response"})
	}

	return c.JSON(http.StatusOK, extractResp)
}

// GetExtractionJob proxies the request to Python service
func (h *ExtractionHandler) GetExtractionJob(c echo.Context) error {
	jobID := c.Param("id")
	if jobID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Job ID required"})
	}

	resp, err := http.Get(fmt.Sprintf("%s/result/%s", PythonServiceURL, jobID))
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "Failed to connect to extraction service"})
	}
	defer resp.Body.Close()

	// Stream response back
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Response().Writer, resp.Body)
	return err
}

// ============================
// Upload + Process Entry Point (Legacy)
// ============================

func (h *ExtractionHandler) ProcessPDF(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".pdf") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Only PDF files are supported"})
	}

	inputDir := "../final_extraction_system/input_pdfs/"
	if err := os.MkdirAll(inputDir, 0o755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create input directory"})
	}

	timestamp := time.Now().Format("20060102_150405")
	uniqueFilename := fmt.Sprintf("%s_%s", timestamp, file.Filename)
	filePath := filepath.Join(inputDir, uniqueFilename)

	if err := saveUploadedFile(file, filePath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to save file: %v", err)})
	}

	extractionResult, err := h.runReductoExtraction(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Extraction failed: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":           "PDF processed successfully",
		"allBlocks":         extractionResult.AllBlocks,
		"filename":          uniqueFilename,
		"original_filename": file.Filename,
		"processed_at":      time.Now().Format(time.RFC3339),
	})
}

func saveUploadedFile(fileHeader *multipart.FileHeader, dstPath string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// ============================
// Run Python + Parse Output
// ============================

type FrontendResponse struct {
	AllBlocks []map[string]interface{} `json:"allBlocks"`
	Filename  string                   `json:"filename"`
}

func (h *ExtractionHandler) runReductoExtraction(pdfPath string) (*FrontendResponse, error) {
	log.Println("🚀 Starting Reducto extraction")

	root := "../"
	pythonBinary := filepath.Join(root, "temp_env", "bin", "python")
	pythonScriptPath := filepath.Join(root, "final_extraction_system", "extractor.py")

	if _, err := os.Stat(pythonBinary); os.IsNotExist(err) {
		return nil, fmt.Errorf("python venv not found at %s", pythonBinary)
	}
	if _, err := os.Stat(pythonScriptPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("extractor.py not found at %s", pythonScriptPath)
	}

	cmd := exec.Command(pythonBinary, pythonScriptPath, pdfPath)
	cmd.Dir = filepath.Join(root, "final_extraction_system")

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start extractor: %w", err)
	}
	go stream("PY-OUT", stdout)
	go stream("PY-ERR", stderr)

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("extractor failed: %w", err)
	}

	outputDir := filepath.Join(root, "final_extraction_system", "output", "reducto")
	latestFile, err := getLatestJSON(outputDir)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(latestFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON: %v", err)
	}

	var resp ReductoResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	var allBlocks []map[string]interface{}
	for _, page := range resp.Pages {
		var chunks []ReductoChunk
		// pipeline.run shape
		if page.Result.Parse != nil && page.Result.Parse.Result.Chunks != nil {
			chunks = page.Result.Parse.Result.Chunks
		} else if page.Result.Chunks != nil {
			// parse.run shape
			chunks = page.Result.Chunks
		}
		for _, ch := range chunks {
			for _, b := range ch.Blocks {
				allBlocks = append(allBlocks, map[string]interface{}{
					"type":       b.Type,
					"content":    b.Content,
					"page":       page.PageNumber,
					"confidence": b.Confidence,
					"bbox":       b.BBox,
				})
			}
		}
	}

	return &FrontendResponse{
		AllBlocks: allBlocks,
		Filename:  resp.Filename,
	}, nil
}

func stream(prefix string, r io.Reader) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		log.Printf("[%s] %s", prefix, sc.Text())
	}
}

// ============================
// Reducto JSON (supports both shapes)
// ============================

type ReductoResponse struct {
	Filename            string               `json:"filename"`
	ExtractionDate      string               `json:"extraction_date"`
	ExtractionMethod    string               `json:"extraction_method"`
	TotalPagesProcessed int                  `json:"total_pages_processed"`
	Pages               []ReductoPageWrapper `json:"pages"`
}

type ReductoPageWrapper struct {
	JobID      string            `json:"job_id"`
	Result     ReductoPageResult `json:"result"`
	PageNumber int               `json:"page_number"`
}

type ReductoPageResult struct {
	Chunks []ReductoChunk `json:"chunks"` // legacy parse.run
	Parse  *ParseWrapper  `json:"parse"`  // pipeline.run
}

type ParseWrapper struct {
	Duration float64     `json:"duration"`
	JobID    string      `json:"job_id"`
	Result   ParseResult `json:"result"`
}

type ParseResult struct {
	Chunks []ReductoChunk `json:"chunks"`
}

type ReductoChunk struct {
	Blocks []ReductoBlock `json:"blocks"`
}

type ReductoBlock struct {
	Type       string                 `json:"type"`
	Content    string                 `json:"content"`
	Confidence interface{}            `json:"confidence"`
	BBox       map[string]interface{} `json:"bbox"`
}

// ============================
// QC Save Endpoint
// ============================

type QCPayload struct {
	BlockID   string                 `json:"block_id"`
	Page      int                    `json:"page"`
	Type      string                 `json:"type"`
	Filename  string                 `json:"filename"`
	Original  string                 `json:"original_html"`
	Edited    string                 `json:"edited_html"`
	TableJSON map[string]interface{} `json:"table_json"`
	QC        map[string]interface{} `json:"qc"`
}

func (h *ExtractionHandler) SaveQC(c echo.Context) error {
	var p QCPayload
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if p.BlockID == "" || p.Filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing block_id or filename"})
	}

	outDir := filepath.Join("../final_extraction_system/output/qc", p.Filename)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "mkdir failed"})
	}
	outPath := filepath.Join(outDir, fmt.Sprintf("page_%d_%s.json", p.Page, sanitize(p.BlockID)))

	b, _ := json.MarshalIndent(p, "", "  ")
	if err := os.WriteFile(outPath, b, 0o644); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "write failed"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "saved"})
}

// ============================
// Serve PDF for viewer
// ============================

func (h *ExtractionHandler) ServePDF(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Filename is required"})
	}
	p := filepath.Join("..", "final_extraction_system", "input_pdfs", filename)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "PDF not found"})
	}

	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
	return c.File(p)
}

// ============================
// Helpers
// ============================

func getLatestJSON(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read dir: %v", err)
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(e.Name()), ".json") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no JSON output file found in %s", dir)
	}
	sort.Slice(files, func(i, j int) bool {
		fi, _ := os.Stat(files[i])
		fj, _ := os.Stat(files[j])
		return fi.ModTime().After(fj.ModTime())
	})
	return files[0], nil
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, string(os.PathSeparator), "_")
	return strings.TrimSpace(s)
}
