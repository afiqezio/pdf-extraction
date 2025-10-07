package models

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

// EPBEBase contains common fields for all ePBE models in exact SQL column order
type EPBEBase struct {
	// Geographic Location Fields (columns 1-10)
	Country           string   `json:"country" gorm:"column:country;size:255"`
	Region            string   `json:"region" gorm:"column:region;size:255"`
	SubRegion         string   `json:"sub_region" gorm:"column:sub_region;size:255"`
	BusinessRegions   string   `json:"business_regions" gorm:"column:business_regions;size:255"`
	Basin             string   `json:"basin" gorm:"column:basin;size:255"`
	SubBasin          string   `json:"sub_basin" gorm:"column:sub_basin;size:255"`
	WellNameFieldName string   `json:"well_name_field_name" gorm:"column:well_name_field_name;size:255"`
	UWI               string   `json:"uwi" gorm:"column:uwi;size:255"`
	Latitude          *float64 `json:"latitude" gorm:"column:latitude;type:decimal(10,7)"`
	Longitude         *float64 `json:"longitude" gorm:"column:longitude;type:decimal(10,7)"`

	// Geological Information Fields (columns 11-18)
	FormationName   string   `json:"formation_name" gorm:"column:formation_name;size:255"`
	ReservoirName   string   `json:"reservoir_name" gorm:"column:reservoir_name;size:255"`
	Period          string   `json:"period" gorm:"column:period;size:255"`
	Epoch           string   `json:"epoch" gorm:"column:epoch;size:255"`
	Age             string   `json:"age" gorm:"column:age;size:255"`
	OnshoreOffshore string   `json:"onshore_offshore" gorm:"column:onshore_offshore;size:255"`
	WaterDepthM     *float64 `json:"water_depth_m" gorm:"column:water_depth_m;type:decimal(10,2)"`
	WaterDepthFt    *float64 `json:"water_depth_ft" gorm:"column:water_depth_ft;type:decimal(10,2)"`
}

// DepthInfo contains depth information fields used across multiple models in SQL column order
type DepthInfo struct {
	// Depth Reference Fields (columns 19-23)
	DepthReferenceType        string   `json:"depth_reference_type" gorm:"column:depth_reference_type;size:255"`
	DepthReferenceElevationM  *float64 `json:"depth_reference_elevation_m" gorm:"column:depth_reference_elevation_m;type:decimal(10,2)"`
	DepthReferenceElevationFt *float64 `json:"depth_reference_elevation_ft" gorm:"column:depth_reference_elevation_ft;type:decimal(10,2)"`
	GroundLevelElevationM     *float64 `json:"ground_level_elevation_m" gorm:"column:ground_level_elevation_m;type:decimal(10,2)"`
	GroundLevelElevationFt    *float64 `json:"ground_level_elevation_ft" gorm:"column:ground_level_elevation_ft;type:decimal(10,2)"`

	// Metric depths (columns vary by table but follow this order)
	TopDepthMMDDF     *float64 `json:"top_depth_mmddf" gorm:"column:top_depth_mmddf;type:decimal(10,2)"`
	TopDepthMTVDDF    *float64 `json:"top_depth_mtvddf" gorm:"column:top_depth_mtvddf;type:decimal(10,2)"`
	TopDepthMTVDSS    *float64 `json:"top_depth_mtvdss" gorm:"column:top_depth_mtvdss;type:decimal(10,2)"`
	TopDepthMBML      *float64 `json:"top_depth_mbml" gorm:"column:top_depth_mbml;type:decimal(10,2)"`
	BottomDepthMMDDF  *float64 `json:"bottom_depth_mmddf" gorm:"column:bottom_depth_mmddf;type:decimal(10,2)"`
	BottomDepthMTVDDF *float64 `json:"bottom_depth_mtvddf" gorm:"column:bottom_depth_mtvddf;type:decimal(10,2)"`
	BottomDepthMTVDSS *float64 `json:"bottom_depth_mtvdss" gorm:"column:bottom_depth_mtvdss;type:decimal(10,2)"`
	BottomDepthMBML   *float64 `json:"bottom_depth_mbml" gorm:"column:bottom_depth_mbml;type:decimal(10,2)"`

	// Imperial depths
	TopDepthFtMDDF     *float64 `json:"top_depth_ftmddf" gorm:"column:top_depth_ftmddf;type:decimal(10,2)"`
	TopDepthFtTVDDF    *float64 `json:"top_depth_fttvddf" gorm:"column:top_depth_fttvddf;type:decimal(10,2)"`
	TopDepthFtTVDSS    *float64 `json:"top_depth_fttvdss" gorm:"column:top_depth_fttvdss;type:decimal(10,2)"`
	TopDepthFtBML      *float64 `json:"top_depth_ftbml" gorm:"column:top_depth_ftbml;type:decimal(10,2)"`
	BottomDepthFtMDDF  *float64 `json:"bottom_depth_ftmddf" gorm:"column:bottom_depth_ftmddf;type:decimal(10,2)"`
	BottomDepthFtTVDDF *float64 `json:"bottom_depth_fttvddf" gorm:"column:bottom_depth_fttvddf;type:decimal(10,2)"`
	BottomDepthFtTVDSS *float64 `json:"bottom_depth_fttvdss" gorm:"column:bottom_depth_fttvdss;type:decimal(10,2)"`
	BottomDepthFtBML   *float64 `json:"bottom_depth_ftbml" gorm:"column:bottom_depth_ftbml;type:decimal(10,2)"`
}

// MetadataInfo contains metadata fields that appear at the end of all tables in SQL column order
type MetadataInfo struct {
	// Data Management Fields (appear near end of each table)
	DataSource             string     `json:"data_source" gorm:"column:data_source;size:255"`
	AnalysisType           string     `json:"analysis_type" gorm:"column:analysis_type;size:255"`
	Owner                  string     `json:"owner" gorm:"column:owner;size:255"`
	Ownership              string     `json:"ownership" gorm:"column:ownership;size:255"`
	Assurance              string     `json:"assurance" gorm:"column:assurance;size:255"`
	DataGenerator          string     `json:"data_generator" gorm:"column:data_generator;size:255"`
	DataGenerationDate     *time.Time `json:"data_generation_date" gorm:"column:data_generation_date"`
	Remark                 string     `json:"remark" gorm:"column:remark;type:text"`
	DataEntryDate          *time.Time `json:"data_entry_date" gorm:"column:data_entry_date"`
	DataEntryMode          string     `json:"data_entry_mode" gorm:"column:data_entry_mode;size:255"`
	DataEntryFocal         string     `json:"data_entry_focal" gorm:"column:data_entry_focal;size:255"`
	MetadataDisciplineName string     `json:"metadata_discipline_name" gorm:"column:metadata_discipline_name;size:100"`
	MetadataDataSourceName string     `json:"metadata_data_source_name" gorm:"column:metadata_data_source_name;size:100"`

	// Session tracking
	SessionID string `json:"session_id" gorm:"column:session_id;size:100;index"`

	UpdatedTimestamp time.Time `json:"updated_timestamp" gorm:"column:updated_timestamp;default:CURRENT_TIMESTAMP"`
	CreatedTimestamp time.Time `json:"created_timestamp" gorm:"column:created_timestamp;default:CURRENT_TIMESTAMP"`

	// Primary Key (comes after metadata fields in SQL)
	ID uint `json:"id" gorm:"column:id;primaryKey;autoIncrement"`

	// Duplicate detection fields (last columns in SQL)
	DuplicateStatus           string     `json:"duplicate_status" gorm:"column:duplicate_status;size:50"`
	DuplicateResolutionAction string     `json:"duplicate_resolution_action" gorm:"column:duplicate_resolution_action;size:50"`
	MasterRecordID            *uint      `json:"master_record_id" gorm:"column:master_record_id"`
	ReviewQueueID             string     `json:"review_queue_id" gorm:"column:review_queue_id;size:50"`
	ResolutionTimestamp       *time.Time `json:"resolution_timestamp" gorm:"column:resolution_timestamp"`
	ResolutionReason          string     `json:"resolution_reason" gorm:"column:resolution_reason;size:500"`

	// Soft delete (not in SQL but needed for GORM)
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
}

// NullFloat64 is a custom type for handling SQL NULL values more elegantly
type NullFloat64 struct {
	Float64 float64
	Valid   bool
}

func (nf *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		nf.Float64, nf.Valid = 0, false
		return nil
	}
	nf.Valid = true
	switch v := value.(type) {
	case float64:
		nf.Float64 = v
	case float32:
		nf.Float64 = float64(v)
	default:
		nf.Valid = false
	}
	return nil
}

func (nf NullFloat64) Value() (driver.Value, error) {
	if !nf.Valid {
		return nil, nil
	}
	return nf.Float64, nil
}
