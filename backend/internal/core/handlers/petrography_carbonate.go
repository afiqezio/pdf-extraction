package handlers

import (
	"net/http"
	"strconv"

	"workbench/internal/core/models"
	"workbench/internal/database"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PetrographyCarbonateHandler struct {
	db *gorm.DB
}

func NewPetrographyCarbonateHandler(db *gorm.DB) *PetrographyCarbonateHandler {
	return &PetrographyCarbonateHandler{db: db}
}

func (h *PetrographyCarbonateHandler) PetrographyCarbonateRoutes(g *echo.Group) {
	carbonate := g.Group("/petrography-carbonate")
	carbonate.POST("", h.CreatePetrographyCarbonate)
	carbonate.GET("", h.GetPetrographyCarbonateRecords)
	carbonate.GET("/search", h.SearchPetrographyCarbonateRecords)
	carbonate.GET("/:id", h.GetPetrographyCarbonate)
	carbonate.PUT("/:id", h.UpdatePetrographyCarbonate)
	carbonate.DELETE("/:id", h.DeletePetrographyCarbonate)
}

// CreatePetrographyCarbonate creates a new petrography carbonate record
func (h *PetrographyCarbonateHandler) CreatePetrographyCarbonate(c echo.Context) error {
	var record models.EPBEPetrographyCarbonate

	if err := c.Bind(&record); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if record.Country == "" || record.WellNameFieldName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Country and Well Name/Field Name are required",
		})
	}

	// Create record
	if err := h.db.Create(&record).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create petrography carbonate record",
		})
	}

	return c.JSON(http.StatusCreated, record)
}

// GetPetrographyCarbonate retrieves a petrography carbonate record by ID
func (h *PetrographyCarbonateHandler) GetPetrographyCarbonate(c echo.Context) error {
	id := c.Param("id")

	recordID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid record ID",
		})
	}

	var record models.EPBEPetrographyCarbonate
	if err := h.db.Where("id = ?", uint(recordID)).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Petrography carbonate record not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve petrography carbonate record",
		})
	}

	return c.JSON(http.StatusOK, record)
}

// GetPetrographyCarbonateRecords retrieves all petrography carbonate records with pagination
func (h *PetrographyCarbonateHandler) GetPetrographyCarbonateRecords(c echo.Context) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	pagination := &models.Pagination{
		Page:  page,
		Limit: limit,
	}

	var records []models.EPBEPetrographyCarbonate
	var total int64

	// Count total records
	h.db.Model(&models.EPBEPetrographyCarbonate{}).Count(&total)

	// Get records with pagination
	if err := h.db.Scopes(database.Paginate(pagination)).Find(&records).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve petrography carbonate records",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"records": records,
		"pagination": map[string]interface{}{
			"page":        pagination.GetPage(),
			"limit":       pagination.GetLimit(),
			"total":       total,
			"total_pages": (total + int64(pagination.GetLimit()) - 1) / int64(pagination.GetLimit()),
		},
	})
}

// UpdatePetrographyCarbonate updates a petrography carbonate record by ID
func (h *PetrographyCarbonateHandler) UpdatePetrographyCarbonate(c echo.Context) error {
	id := c.Param("id")

	recordID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid record ID",
		})
	}

	var record models.EPBEPetrographyCarbonate
	if err := h.db.Where("id = ?", uint(recordID)).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Petrography carbonate record not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve petrography carbonate record",
		})
	}

	// Bind update data
	var updateData models.EPBEPetrographyCarbonate
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Update record fields (exclude ID and timestamps)
	record.Country = updateData.Country
	record.Region = updateData.Region
	record.SubRegion = updateData.SubRegion
	record.BusinessRegions = updateData.BusinessRegions
	record.Basin = updateData.Basin
	record.SubBasin = updateData.SubBasin
	record.WellNameFieldName = updateData.WellNameFieldName
	record.UWI = updateData.UWI
	record.Latitude = updateData.Latitude
	record.Longitude = updateData.Longitude
	record.FormationName = updateData.FormationName
	record.ReservoirName = updateData.ReservoirName
	record.Period = updateData.Period
	record.Epoch = updateData.Epoch
	record.Age = updateData.Age
	record.OnshoreOffshore = updateData.OnshoreOffshore
	record.WaterDepthM = updateData.WaterDepthM
	record.WaterDepthFt = updateData.WaterDepthFt

	// Update depth information
	record.DepthReferenceType = updateData.DepthReferenceType
	record.DepthReferenceElevationM = updateData.DepthReferenceElevationM
	record.DepthReferenceElevationFt = updateData.DepthReferenceElevationFt
	record.GroundLevelElevationM = updateData.GroundLevelElevationM
	record.GroundLevelElevationFt = updateData.GroundLevelElevationFt
	record.TopDepthMMDDF = updateData.TopDepthMMDDF
	record.TopDepthMTVDDF = updateData.TopDepthMTVDDF
	record.TopDepthMTVDSS = updateData.TopDepthMTVDSS
	record.TopDepthMBML = updateData.TopDepthMBML
	record.BottomDepthMMDDF = updateData.BottomDepthMMDDF
	record.BottomDepthMTVDDF = updateData.BottomDepthMTVDDF
	record.BottomDepthMTVDSS = updateData.BottomDepthMTVDSS
	record.BottomDepthMBML = updateData.BottomDepthMBML
	record.TopDepthFtMDDF = updateData.TopDepthFtMDDF
	record.TopDepthFtTVDDF = updateData.TopDepthFtTVDDF
	record.TopDepthFtTVDSS = updateData.TopDepthFtTVDSS
	record.TopDepthFtBML = updateData.TopDepthFtBML
	record.BottomDepthFtMDDF = updateData.BottomDepthFtMDDF
	record.BottomDepthFtTVDDF = updateData.BottomDepthFtTVDDF
	record.BottomDepthFtTVDSS = updateData.BottomDepthFtTVDSS
	record.BottomDepthFtBML = updateData.BottomDepthFtBML

	// Update facies description
	record.LithofaciesCore = updateData.LithofaciesCore
	record.MicrofaciesThinSection = updateData.MicrofaciesThinSection
	record.Depofacies = updateData.Depofacies

	// Update porosity and permeability
	record.VisiblePorosityPercent = updateData.VisiblePorosityPercent
	record.HePorosityPercent = updateData.HePorosityPercent
	record.PermeabilityMd = updateData.PermeabilityMd

	// Update matrix mineralogy
	record.Calcite = updateData.Calcite
	record.Dolomite = updateData.Dolomite
	record.Micrite = updateData.Micrite
	record.MicriteEnvelopes = updateData.MicriteEnvelopes
	record.MicrosparPseudospar = updateData.MicrosparPseudospar
	record.Kaolinite = updateData.Kaolinite
	record.Clay = updateData.Clay
	record.TotalMineralogyMatrixPercent = updateData.TotalMineralogyMatrixPercent

	// Update bioclasts - general and specific foraminifera
	record.Bioclasts = updateData.Bioclasts
	record.Lepido = updateData.Lepido
	record.Coral = updateData.Coral
	record.Rhodolith = updateData.Rhodolith
	record.RedAlgae = updateData.RedAlgae
	record.RedAlgaeEnc = updateData.RedAlgaeEnc
	record.GreenAlgae = updateData.GreenAlgae
	record.Echinoderms = updateData.Echinoderms
	record.Miliolid = updateData.Miliolid
	record.Lepidocyclina = updateData.Lepidocyclina
	record.Cycloclypeus = updateData.Cycloclypeus
	record.Operculina = updateData.Operculina
	record.OtherRotaliids = updateData.OtherRotaliids
	record.Gypsinid = updateData.Gypsinid
	record.Planorbulinella = updateData.Planorbulinella
	record.Hemotremid = updateData.Hemotremid
	record.Heterostegina = updateData.Heterostegina
	record.EncFrm = updateData.EncFrm
	record.Planktonic = updateData.Planktonic
	record.Bryozoans = updateData.Bryozoans
	record.Amphistegina = updateData.Amphistegina
	record.Gastropods = updateData.Gastropods
	record.Bivalve = updateData.Bivalve
	record.Ostracod = updateData.Ostracod
	record.Oncoids = updateData.Oncoids
	record.UndiffMolluscs = updateData.UndiffMolluscs
	record.UndiffBenthonic = updateData.UndiffBenthonic
	record.UndiffSkeletal = updateData.UndiffSkeletal
	record.UndiffForam = updateData.UndiffForam
	record.TotalSkeletalPercent = updateData.TotalSkeletalPercent

	// Update non-skeletal components
	record.Organic = updateData.Organic
	record.Peloids = updateData.Peloids
	record.MicritisedGrains = updateData.MicritisedGrains
	record.Pseudoclasts = updateData.Pseudoclasts
	record.Intraclast = updateData.Intraclast
	record.Quartz = updateData.Quartz
	record.TotalNonSkeletalPercent = updateData.TotalNonSkeletalPercent

	// Update porosity types
	record.Interparticle = updateData.Interparticle
	record.Intraparticle = updateData.Intraparticle
	record.Intercrystalline = updateData.Intercrystalline
	record.MatrixIntercrystalline = updateData.MatrixIntercrystalline
	record.Mouldic = updateData.Mouldic
	record.Vuggy = updateData.Vuggy
	record.Fractures = updateData.Fractures
	record.Micro = updateData.Micro
	record.TotalPorosityPercent = updateData.TotalPorosityPercent

	// Update cement types
	record.Fringing = updateData.Fringing
	record.Meniscus = updateData.Meniscus
	record.Blocky = updateData.Blocky
	record.Sparry = updateData.Sparry
	record.Micritic = updateData.Micritic
	record.Pendant = updateData.Pendant
	record.Syntax = updateData.Syntax
	record.CalciteSyntaxial = updateData.CalciteSyntaxial
	record.CalciteFringing = updateData.CalciteFringing
	record.CalciteMosaic = updateData.CalciteMosaic
	record.CalciteBlocky = updateData.CalciteBlocky
	record.CalciteFerroan = updateData.CalciteFerroan
	record.Pyrite = updateData.Pyrite
	record.Fluorite = updateData.Fluorite
	record.TotalCementPercent = updateData.TotalCementPercent

	// Update replacement and accessories
	record.Replacement = updateData.Replacement
	record.Saddle = updateData.Saddle
	record.TotalDolomitePercent = updateData.TotalDolomitePercent
	record.Stylolite = updateData.Stylolite
	record.Bioturbation = updateData.Bioturbation
	record.TotalAccessoriesPercent = updateData.TotalAccessoriesPercent
	record.TotalPercent = updateData.TotalPercent

	// Update analysis type
	record.AnalysisTypes = updateData.AnalysisTypes

	// Update metadata
	record.DataSource = updateData.DataSource
	record.AnalysisType = updateData.AnalysisType
	record.Owner = updateData.Owner
	record.Ownership = updateData.Ownership
	record.Assurance = updateData.Assurance
	record.DataGenerator = updateData.DataGenerator
	record.DataGenerationDate = updateData.DataGenerationDate
	record.Remark = updateData.Remark
	record.DataEntryDate = updateData.DataEntryDate
	record.DataEntryMode = updateData.DataEntryMode
	record.DataEntryFocal = updateData.DataEntryFocal
	record.MetadataDisciplineName = updateData.MetadataDisciplineName
	record.MetadataDataSourceName = updateData.MetadataDataSourceName
	record.SessionID = updateData.SessionID

	if err := h.db.Save(&record).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update petrography carbonate record",
		})
	}

	return c.JSON(http.StatusOK, record)
}

// DeletePetrographyCarbonate soft deletes a petrography carbonate record by ID
func (h *PetrographyCarbonateHandler) DeletePetrographyCarbonate(c echo.Context) error {
	id := c.Param("id")

	recordID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid record ID",
		})
	}

	var record models.EPBEPetrographyCarbonate
	if err := h.db.Where("id = ?", uint(recordID)).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Petrography carbonate record not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve petrography carbonate record",
		})
	}

	// Soft delete
	if err := h.db.Delete(&record).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete petrography carbonate record",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Petrography carbonate record deleted successfully",
	})
}

// SearchPetrographyCarbonateRecords searches petrography carbonate records by various fields
func (h *PetrographyCarbonateHandler) SearchPetrographyCarbonateRecords(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Search query is required",
		})
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	pagination := &models.Pagination{
		Page:  page,
		Limit: limit,
	}

	var records []models.EPBEPetrographyCarbonate
	var total int64

	// Count total records with search
	searchQuery := h.db.Model(&models.EPBEPetrographyCarbonate{}).Scopes(
		database.Search(query,
			"country",
			"region",
			"sub_region",
			"basin",
			"sub_basin",
			"well_name_field_name",
			"uwi",
			"formation_name",
			"reservoir_name",
			"period",
			"epoch",
			"age",
			"lithofacies_core",
			"microfacies_thin_section",
			"depofacies",
			"analysis_types",
		),
	)
	searchQuery.Count(&total)

	// Get records with search and pagination
	if err := searchQuery.Scopes(database.Paginate(pagination)).Find(&records).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to search petrography carbonate records",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"records": records,
		"pagination": map[string]interface{}{
			"page":        pagination.GetPage(),
			"limit":       pagination.GetLimit(),
			"total":       total,
			"total_pages": (total + int64(pagination.GetLimit()) - 1) / int64(pagination.GetLimit()),
		},
		"query": query,
	})
}
