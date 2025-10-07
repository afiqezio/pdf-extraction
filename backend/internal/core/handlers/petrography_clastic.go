package handlers

import (
	"net/http"
	"strconv"

	"workbench/internal/core/models"
	"workbench/internal/database"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PetrographyClasticHandler struct {
	db *gorm.DB
}

func NewPetrographyClasticHandler(db *gorm.DB) *PetrographyClasticHandler {
	return &PetrographyClasticHandler{db: db}
}

func (h *PetrographyClasticHandler) PetrographyClasticRoutes(g *echo.Group) {
	clastic := g.Group("/petrography-clastic")
	clastic.POST("", h.CreatePetrographyClastic)
	clastic.GET("", h.GetPetrographyClasticRecords)
	clastic.GET("/search", h.SearchPetrographyClasticRecords)
	clastic.GET("/:id", h.GetPetrographyClastic)
	clastic.PUT("/:id", h.UpdatePetrographyClastic)
	clastic.DELETE("/:id", h.DeletePetrographyClastic)
}

// CreatePetrographyClastic creates a new petrography clastic record
func (h *PetrographyClasticHandler) CreatePetrographyClastic(c echo.Context) error {
	var record models.EPBEPetrographyClastic

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
			"error": "Failed to create petrography clastic record",
		})
	}

	return c.JSON(http.StatusCreated, record)
}

// GetPetrographyClastic retrieves a petrography clastic record by ID
func (h *PetrographyClasticHandler) GetPetrographyClastic(c echo.Context) error {
	id := c.Param("id")

	recordID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid record ID",
		})
	}

	var record models.EPBEPetrographyClastic
	if err := h.db.Where("id = ?", uint(recordID)).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Petrography clastic record not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve petrography clastic record",
		})
	}

	return c.JSON(http.StatusOK, record)
}

// GetPetrographyClasticRecords retrieves all petrography clastic records with pagination
func (h *PetrographyClasticHandler) GetPetrographyClasticRecords(c echo.Context) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	pagination := &models.Pagination{
		Page:  page,
		Limit: limit,
	}

	var records []models.EPBEPetrographyClastic
	var total int64

	// Count total records
	h.db.Model(&models.EPBEPetrographyClastic{}).Count(&total)

	// Get records with pagination
	if err := h.db.Scopes(database.Paginate(pagination)).Find(&records).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve petrography clastic records",
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

// UpdatePetrographyClastic updates a petrography clastic record by ID
func (h *PetrographyClasticHandler) UpdatePetrographyClastic(c echo.Context) error {
	id := c.Param("id")

	recordID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid record ID",
		})
	}

	var record models.EPBEPetrographyClastic
	if err := h.db.Where("id = ?", uint(recordID)).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Petrography clastic record not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve petrography clastic record",
		})
	}

	// Bind update data
	var updateData models.EPBEPetrographyClastic
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

	// Update textural analysis
	record.GrainSize = updateData.GrainSize
	record.GrainShape = updateData.GrainShape
	record.GrainContact = updateData.GrainContact
	record.SedimentaryStructure = updateData.SedimentaryStructure
	record.Sorting = updateData.Sorting
	record.Lithofacies = updateData.Lithofacies

	// Update porosity and permeability
	record.VisiblePorosityPercent = updateData.VisiblePorosityPercent
	record.AmbientHePorosityPercent = updateData.AmbientHePorosityPercent
	record.PermeabilityMd = updateData.PermeabilityMd
	record.GrainDensityGCc = updateData.GrainDensityGCc

	// Update quartz content
	record.MonocrystallineQuartz = updateData.MonocrystallineQuartz
	record.PolycrystallineQuartz = updateData.PolycrystallineQuartz
	record.TotalQuartzPercent = updateData.TotalQuartzPercent

	// Update feldspar content
	record.PotassiumFeldspar = updateData.PotassiumFeldspar
	record.Plagioclase = updateData.Plagioclase
	record.FeldsparUndifferentiated = updateData.FeldsparUndifferentiated
	record.TotalFeldsparPercent = updateData.TotalFeldsparPercent

	// Update mica content
	record.Muscovite = updateData.Muscovite
	record.Biotite = updateData.Biotite
	record.MicaUndifferentiated = updateData.MicaUndifferentiated
	record.TotalMicaPercent = updateData.TotalMicaPercent

	// Update heavy minerals
	record.Zircon = updateData.Zircon
	record.Tourmaline = updateData.Tourmaline
	record.HeavyMineralsUndifferentiated = updateData.HeavyMineralsUndifferentiated
	record.TotalHeavyMineralsPercent = updateData.TotalHeavyMineralsPercent

	// Update igneous rock fragments
	record.PlutonicRockFragments = updateData.PlutonicRockFragments
	record.MaficIntermediateVolcanicFragment = updateData.MaficIntermediateVolcanicFragment
	record.VolcanicRockFragment = updateData.VolcanicRockFragment
	record.TotalIgneousRFPercent = updateData.TotalIgneousRFPercent

	// Update metamorphic rock fragments
	record.QuartzoseRockFragment = updateData.QuartzoseRockFragment
	record.SchistoseRockFragment = updateData.SchistoseRockFragment
	record.MetamorphicRockFragmentUndifferentiated = updateData.MetamorphicRockFragmentUndifferentiated
	record.TotalMetamorphicRFPercent = updateData.TotalMetamorphicRFPercent

	// Update sedimentary rock fragments
	record.SandstoneSiltstoneRockFragments = updateData.SandstoneSiltstoneRockFragments
	record.ArgillaceousRockFragments = updateData.ArgillaceousRockFragments
	record.SiliciclasticRockFragmentsUndifferentiated = updateData.SiliciclasticRockFragmentsUndifferentiated
	record.LimestoneRockFragments = updateData.LimestoneRockFragments
	record.DolostoneRockFragments = updateData.DolostoneRockFragments
	record.Chert = updateData.Chert
	record.TotalSedimentaryRFPercent = updateData.TotalSedimentaryRFPercent
	record.TotalRockFragmentsPercent = updateData.TotalRockFragmentsPercent

	// Update other grains
	record.RipUpClast = updateData.RipUpClast
	record.Glauconite = updateData.Glauconite
	record.Bioclast = updateData.Bioclast
	record.ForaminiferaGrains = updateData.ForaminiferaGrains
	record.UndifferentiatedOtherGrains = updateData.UndifferentiatedOtherGrains
	record.TotalOtherGrainsPercent = updateData.TotalOtherGrainsPercent

	// Update matrix
	record.ClayMatrix = updateData.ClayMatrix
	record.MixedClaySiltFineMatrix = updateData.MixedClaySiltFineMatrix
	record.SiltVeryFineMatrix = updateData.SiltVeryFineMatrix
	record.OrganicMatrix = updateData.OrganicMatrix
	record.MatrixUndifferentiated = updateData.MatrixUndifferentiated
	record.TotalMatrixPercent = updateData.TotalMatrixPercent

	// Update authigenic clay
	record.Kaolinite = updateData.Kaolinite
	record.KaoliniteReplacesKFeldspar = updateData.KaoliniteReplacesKFeldspar
	record.IllitePoreGrainLining = updateData.IllitePoreGrainLining
	record.IllitePoreFilling = updateData.IllitePoreFilling
	record.IlliteReplacesKFeldspar = updateData.IlliteReplacesKFeldspar
	record.TotalAuthigenicClayPercent = updateData.TotalAuthigenicClayPercent

	// Update authigenic non-clay
	record.SyntaxialQuartzOvergrowths = updateData.SyntaxialQuartzOvergrowths
	record.FeldsparOvergrowths = updateData.FeldsparOvergrowths
	record.FeCalcite = updateData.FeCalcite
	record.FeDolomite = updateData.FeDolomite
	record.Siderite = updateData.Siderite
	record.MnSiderite = updateData.MnSiderite
	record.Pyrite = updateData.Pyrite
	record.IronOxideMinerals = updateData.IronOxideMinerals
	record.TotalAuthigenicNonClayPercent = updateData.TotalAuthigenicNonClayPercent

	// Update primary porosity
	record.Intergranular = updateData.Intergranular
	record.Intercrystalline = updateData.Intercrystalline
	record.PriPorosityIntragranular = updateData.PriPorosityIntragranular
	record.TotalPrimaryPorosityPercent = updateData.TotalPrimaryPorosityPercent

	// Update secondary porosity
	record.SecPorosityIntragranular = updateData.SecPorosityIntragranular
	record.Intracrystalline = updateData.Intracrystalline
	record.Mouldic = updateData.Mouldic
	record.Fracture = updateData.Fracture
	record.TotalSecondaryPorosityPercent = updateData.TotalSecondaryPorosityPercent

	// Update total and analysis
	record.TotalPercent = updateData.TotalPercent
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
			"error": "Failed to update petrography clastic record",
		})
	}

	return c.JSON(http.StatusOK, record)
}

// DeletePetrographyClastic soft deletes a petrography clastic record by ID
func (h *PetrographyClasticHandler) DeletePetrographyClastic(c echo.Context) error {
	id := c.Param("id")

	recordID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid record ID",
		})
	}

	var record models.EPBEPetrographyClastic
	if err := h.db.Where("id = ?", uint(recordID)).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Petrography clastic record not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve petrography clastic record",
		})
	}

	// Soft delete
	if err := h.db.Delete(&record).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete petrography clastic record",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Petrography clastic record deleted successfully",
	})
}

// SearchPetrographyClasticRecords searches petrography clastic records by various fields
func (h *PetrographyClasticHandler) SearchPetrographyClasticRecords(c echo.Context) error {
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

	var records []models.EPBEPetrographyClastic
	var total int64

	// Count total records with search
	searchQuery := h.db.Model(&models.EPBEPetrographyClastic{}).Scopes(
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
			"grain_size",
			"grain_shape",
			"lithofacies",
			"analysis_types",
		),
	)
	searchQuery.Count(&total)

	// Get records with search and pagination
	if err := searchQuery.Scopes(database.Paginate(pagination)).Find(&records).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to search petrography clastic records",
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
