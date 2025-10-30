package models

import (
	"gorm.io/gorm"
	"time"
)

// PetrographyUnified represents a unified petrography table with ALL possible fields
// This single table can handle both carbonate and clastic data
type PetrographyUnified struct {
	// Primary Key
	ID uint `gorm:"primaryKey" json:"id"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// ========================================
	// BASE FIELDS (Essential identifiers only)
	// ========================================
	WellNameFieldName *string  `gorm:"type:varchar(255);index" json:"well_name_field_name,omitempty"`
	Depth             *float64 `gorm:"type:decimal(10,2);index" json:"depth,omitempty"` // Single depth column

	// ========================================
	// FACIES DESCRIPTION
	// ========================================
	LithofaciesCore        *string `gorm:"type:text" json:"lithofacies_core,omitempty"`
	MicrofaciesThinSection *string `gorm:"type:text" json:"microfacies_thin_section,omitempty"`
	Depofacies             *string `gorm:"type:text" json:"depofacies,omitempty"`

	// ========================================
	// TEXTURAL ANALYSIS (Clastic specific)
	// ========================================
	GrainSize            *string `gorm:"type:varchar(255)" json:"grain_size,omitempty"`
	GrainShape           *string `gorm:"type:varchar(255)" json:"grain_shape,omitempty"`
	GrainContact         *string `gorm:"type:varchar(255)" json:"grain_contact,omitempty"`
	SedimentaryStructure *string `gorm:"type:varchar(255)" json:"sedimentary_structure,omitempty"`
	Sorting              *string `gorm:"type:varchar(255)" json:"sorting,omitempty"`

	// ========================================
	// POROSITY & PERMEABILITY
	// ========================================
	VisiblePorosityPercent   *float64 `gorm:"type:decimal(10,4)" json:"visible_porosity_percent,omitempty"`
	HePorosityPercent        *float64 `gorm:"type:decimal(10,4)" json:"he_porosity_percent,omitempty"`
	AmbientHePorosityPercent *float64 `gorm:"type:decimal(10,4)" json:"ambient_he_porosity_percent,omitempty"`
	PermeabilityMd           *float64 `gorm:"type:decimal(12,4)" json:"permeability_md,omitempty"`
	GrainDensityGCc          *float64 `gorm:"type:decimal(10,4)" json:"grain_density_g_cc,omitempty"`

	// ========================================
	// MATRIX MINERALOGY (Carbonate)
	// ========================================
	Calcite                      *float64 `gorm:"type:decimal(10,4)" json:"calcite,omitempty"`
	Dolomite                     *float64 `gorm:"type:decimal(10,4)" json:"dolomite,omitempty"`
	Micrite                      *float64 `gorm:"type:decimal(10,4)" json:"micrite,omitempty"`
	MicriteEnvelopes             *float64 `gorm:"type:decimal(10,4)" json:"micrite_envelopes,omitempty"`
	MicrosparPseudospar          *float64 `gorm:"type:decimal(10,4)" json:"microspar_pseudospar,omitempty"`
	Kaolinite                    *float64 `gorm:"type:decimal(10,4)" json:"kaolinite,omitempty"`
	Clay                         *float64 `gorm:"type:decimal(10,4)" json:"clay,omitempty"`
	TotalMineralogyMatrixPercent *float64 `gorm:"type:decimal(10,4)" json:"total_mineralogy_matrix_percent,omitempty"`

	// ========================================
	// BIOCLASTS - FORAMINIFERA & SKELETAL (Carbonate)
	// ========================================
	Bioclasts            *float64 `gorm:"type:decimal(10,4)" json:"bioclasts,omitempty"`
	Lepido               *float64 `gorm:"type:decimal(10,4)" json:"lepido,omitempty"`
	Coral                *float64 `gorm:"type:decimal(10,4)" json:"coral,omitempty"`
	Rhodolith            *float64 `gorm:"type:decimal(10,4)" json:"rhodolith,omitempty"`
	RedAlgae             *float64 `gorm:"type:decimal(10,4)" json:"red_algae,omitempty"`
	RedAlgaeEnc          *float64 `gorm:"type:decimal(10,4)" json:"red_algae_enc,omitempty"`
	GreenAlgae           *float64 `gorm:"type:decimal(10,4)" json:"green_algae,omitempty"`
	Echinoderms          *float64 `gorm:"type:decimal(10,4)" json:"echinoderms,omitempty"`
	Miliolid             *float64 `gorm:"type:decimal(10,4)" json:"miliolid,omitempty"`
	Lepidocyclina        *float64 `gorm:"type:decimal(10,4)" json:"lepidocyclina,omitempty"`
	Cycloclypeus         *float64 `gorm:"type:decimal(10,4)" json:"cycloclypeus,omitempty"`
	Operculina           *float64 `gorm:"type:decimal(10,4)" json:"operculina,omitempty"`
	OtherRotaliids       *float64 `gorm:"type:decimal(10,4)" json:"other_rotaliids,omitempty"`
	Gypsinid             *float64 `gorm:"type:decimal(10,4)" json:"gypsinid,omitempty"`
	Planorbulinella      *float64 `gorm:"type:decimal(10,4)" json:"planorbulinella,omitempty"`
	Hemotremid           *float64 `gorm:"type:decimal(10,4)" json:"hemotremid,omitempty"`
	Heterostegina        *float64 `gorm:"type:decimal(10,4)" json:"heterostegina,omitempty"`
	EncFrm               *float64 `gorm:"type:decimal(10,4)" json:"enc_frm,omitempty"`
	Planktonic           *float64 `gorm:"type:decimal(10,4)" json:"planktonic,omitempty"`
	Bryozoans            *float64 `gorm:"type:decimal(10,4)" json:"bryozoans,omitempty"`
	Amphistegina         *float64 `gorm:"type:decimal(10,4)" json:"amphistegina,omitempty"`
	Gastropods           *float64 `gorm:"type:decimal(10,4)" json:"gastropods,omitempty"`
	Bivalve              *float64 `gorm:"type:decimal(10,4)" json:"bivalve,omitempty"`
	Ostracod             *float64 `gorm:"type:decimal(10,4)" json:"ostracod,omitempty"`
	Oncoids              *float64 `gorm:"type:decimal(10,4)" json:"oncoids,omitempty"`
	UndiffMolluscs       *float64 `gorm:"type:decimal(10,4)" json:"undiff_molluscs,omitempty"`
	UndiffBenthonic      *float64 `gorm:"type:decimal(10,4)" json:"undiff_benthonic,omitempty"`
	UndiffSkeletal       *float64 `gorm:"type:decimal(10,4)" json:"undiff_skeletal,omitempty"`
	UndiffForam          *float64 `gorm:"type:decimal(10,4)" json:"undiff_foram,omitempty"`
	TotalSkeletalPercent *float64 `gorm:"type:decimal(10,4)" json:"total_skeletal_percent,omitempty"`

	// ========================================
	// NON-SKELETAL COMPONENTS (Carbonate)
	// ========================================
	Organic                 *float64 `gorm:"type:decimal(10,4)" json:"organic,omitempty"`
	Peloids                 *float64 `gorm:"type:decimal(10,4)" json:"peloids,omitempty"`
	MicritisedGrains        *float64 `gorm:"type:decimal(10,4)" json:"micritised_grains,omitempty"`
	Pseudoclasts            *float64 `gorm:"type:decimal(10,4)" json:"pseudoclasts,omitempty"`
	Intraclast              *float64 `gorm:"type:decimal(10,4)" json:"intraclast,omitempty"`
	Quartz                  *float64 `gorm:"type:decimal(10,4)" json:"quartz,omitempty"`
	TotalNonSkeletalPercent *float64 `gorm:"type:decimal(10,4)" json:"total_non_skeletal_percent,omitempty"`

	// ========================================
	// QUARTZ CONTENT (Clastic)
	// ========================================
	MonocrystallineQuartz *float64 `gorm:"type:decimal(10,4)" json:"monocrystalline_quartz,omitempty"`
	PolycrystallineQuartz *float64 `gorm:"type:decimal(10,4)" json:"polycrystalline_quartz,omitempty"`
	TotalQuartzPercent    *float64 `gorm:"type:decimal(10,4)" json:"total_quartz_percent,omitempty"`

	// ========================================
	// FELDSPAR CONTENT (Clastic)
	// ========================================
	PotassiumFeldspar        *float64 `gorm:"type:decimal(10,4)" json:"potassium_feldspar,omitempty"`
	Plagioclase              *float64 `gorm:"type:decimal(10,4)" json:"plagioclase,omitempty"`
	FeldsparUndifferentiated *float64 `gorm:"type:decimal(10,4)" json:"feldspar_undifferentiated,omitempty"`
	TotalFeldsparPercent     *float64 `gorm:"type:decimal(10,4)" json:"total_feldspar_percent,omitempty"`

	// ========================================
	// MICA CONTENT (Clastic)
	// ========================================
	Muscovite            *float64 `gorm:"type:decimal(10,4)" json:"muscovite,omitempty"`
	Biotite              *float64 `gorm:"type:decimal(10,4)" json:"biotite,omitempty"`
	MicaUndifferentiated *float64 `gorm:"type:decimal(10,4)" json:"mica_undifferentiated,omitempty"`
	TotalMicaPercent     *float64 `gorm:"type:decimal(10,4)" json:"total_mica_percent,omitempty"`

	// ========================================
	// HEAVY MINERALS (Clastic)
	// ========================================
	Zircon                        *float64 `gorm:"type:decimal(10,4)" json:"zircon,omitempty"`
	Tourmaline                    *float64 `gorm:"type:decimal(10,4)" json:"tourmaline,omitempty"`
	HeavyMineralsUndifferentiated *float64 `gorm:"type:decimal(10,4)" json:"heavy_minerals_undifferentiated,omitempty"`
	TotalHeavyMineralsPercent     *float64 `gorm:"type:decimal(10,4)" json:"total_heavy_minerals_percent,omitempty"`

	// ========================================
	// IGNEOUS ROCK FRAGMENTS (Clastic)
	// ========================================
	PlutonicRockFragments             *float64 `gorm:"type:decimal(10,4)" json:"plutonic_rock_fragments,omitempty"`
	MaficIntermediateVolcanicFragment *float64 `gorm:"type:decimal(10,4)" json:"mafic_intermediate_volcanic_fragment,omitempty"`
	VolcanicRockFragment              *float64 `gorm:"type:decimal(10,4)" json:"volcanic_rock_fragment,omitempty"`
	TotalIgneousRFPercent             *float64 `gorm:"type:decimal(10,4)" json:"total_igneous_rf_percent,omitempty"`

	// ========================================
	// METAMORPHIC ROCK FRAGMENTS (Clastic)
	// ========================================
	QuartzoseRockFragment                   *float64 `gorm:"type:decimal(10,4)" json:"quartzose_rock_fragment,omitempty"`
	SchistoseRockFragment                   *float64 `gorm:"type:decimal(10,4)" json:"schistose_rock_fragment,omitempty"`
	MetamorphicRockFragmentUndifferentiated *float64 `gorm:"type:decimal(10,4)" json:"metamorphic_rock_fragment_undifferentiated,omitempty"`
	TotalMetamorphicRFPercent               *float64 `gorm:"type:decimal(10,4)" json:"total_metamorphic_rf_percent,omitempty"`

	// ========================================
	// SEDIMENTARY ROCK FRAGMENTS (Clastic)
	// ========================================
	SandstoneSiltstoneRockFragments            *float64 `gorm:"type:decimal(10,4)" json:"sandstone_siltstone_rock_fragments,omitempty"`
	ArgillaceousRockFragments                  *float64 `gorm:"type:decimal(10,4)" json:"argillaceous_rock_fragments,omitempty"`
	SiliciclasticRockFragmentsUndifferentiated *float64 `gorm:"type:decimal(10,4)" json:"siliciclastic_rock_fragments_undifferentiated,omitempty"`
	LimestoneRockFragments                     *float64 `gorm:"type:decimal(10,4)" json:"limestone_rock_fragments,omitempty"`
	DolostoneRockFragments                     *float64 `gorm:"type:decimal(10,4)" json:"dolostone_rock_fragments,omitempty"`
	Chert                                      *float64 `gorm:"type:decimal(10,4)" json:"chert,omitempty"`
	TotalSedimentaryRFPercent                  *float64 `gorm:"type:decimal(10,4)" json:"total_sedimentary_rf_percent,omitempty"`
	TotalRockFragmentsPercent                  *float64 `gorm:"type:decimal(10,4)" json:"total_rock_fragments_percent,omitempty"`

	// ========================================
	// OTHER GRAINS (Clastic)
	// ========================================
	RipUpClast                  *float64 `gorm:"type:decimal(10,4)" json:"rip_up_clast,omitempty"`
	Glauconite                  *float64 `gorm:"type:decimal(10,4)" json:"glauconite,omitempty"`
	ForaminiferaGrains          *float64 `gorm:"type:decimal(10,4)" json:"foraminifera_grains,omitempty"`
	UndifferentiatedOtherGrains *float64 `gorm:"type:decimal(10,4)" json:"undifferentiated_other_grains,omitempty"`
	TotalOtherGrainsPercent     *float64 `gorm:"type:decimal(10,4)" json:"total_other_grains_percent,omitempty"`

	// ========================================
	// MATRIX (Both Carbonate & Clastic)
	// ========================================
	ClayMatrix              *float64 `gorm:"type:decimal(10,4)" json:"clay_matrix,omitempty"`
	MixedClaySiltFineMatrix *float64 `gorm:"type:decimal(10,4)" json:"mixed_clay_silt_fine_matrix,omitempty"`
	SiltVeryFineMatrix      *float64 `gorm:"type:decimal(10,4)" json:"silt_very_fine_matrix,omitempty"`
	OrganicMatrix           *float64 `gorm:"type:decimal(10,4)" json:"organic_matrix,omitempty"`
	CarbonateMatrix         *float64 `gorm:"type:decimal(10,4)" json:"carbonate_matrix,omitempty"`
	MatrixUndifferentiated  *float64 `gorm:"type:decimal(10,4)" json:"matrix_undifferentiated,omitempty"`
	TotalMatrixPercent      *float64 `gorm:"type:decimal(10,4)" json:"total_matrix_percent,omitempty"`

	// ========================================
	// AUTHIGENIC CLAY (Clastic)
	// ========================================
	KaoliniteReplacesKFeldspar *float64 `gorm:"type:decimal(10,4)" json:"kaolinite_replaces_k_feldspar,omitempty"`
	IllitePoreGrainLining      *float64 `gorm:"type:decimal(10,4)" json:"illite_pore_grain_lining,omitempty"`
	IllitePoreFilling          *float64 `gorm:"type:decimal(10,4)" json:"illite_pore_filling,omitempty"`
	IlliteReplacesKFeldspar    *float64 `gorm:"type:decimal(10,4)" json:"illite_replaces_k_feldspar,omitempty"`
	TotalAuthigenicClayPercent *float64 `gorm:"type:decimal(10,4)" json:"total_authigenic_clay_percent,omitempty"`

	// ========================================
	// AUTHIGENIC NON-CLAY (Clastic)
	// ========================================
	SyntaxialQuartzOvergrowths    *float64 `gorm:"type:decimal(10,4)" json:"syntaxial_quartz_overgrowths,omitempty"`
	FeldsparOvergrowths           *float64 `gorm:"type:decimal(10,4)" json:"feldspar_overgrowths,omitempty"`
	FeCalcite                     *float64 `gorm:"type:decimal(10,4)" json:"fe_calcite,omitempty"`
	FeDolomite                    *float64 `gorm:"type:decimal(10,4)" json:"fe_dolomite,omitempty"`
	Siderite                      *float64 `gorm:"type:decimal(10,4)" json:"siderite,omitempty"`
	MnSiderite                    *float64 `gorm:"type:decimal(10,4)" json:"mn_siderite,omitempty"`
	Pyrite                        *float64 `gorm:"type:decimal(10,4)" json:"pyrite,omitempty"`
	IronOxideMinerals             *float64 `gorm:"type:decimal(10,4)" json:"iron_oxide_minerals,omitempty"`
	Chlorite                      *float64 `gorm:"type:decimal(10,4)" json:"chlorite,omitempty"`
	TotalAuthigenicNonClayPercent *float64 `gorm:"type:decimal(10,4)" json:"total_authigenic_non_clay_percent,omitempty"`

	// ========================================
	// POROSITY TYPES
	// ========================================
	// Carbonate porosity
	Interparticle          *float64 `gorm:"type:decimal(10,4)" json:"interparticle,omitempty"`
	Intraparticle          *float64 `gorm:"type:decimal(10,4)" json:"intraparticle,omitempty"`
	Intercrystalline       *float64 `gorm:"type:decimal(10,4)" json:"intercrystalline,omitempty"`
	MatrixIntercrystalline *float64 `gorm:"type:decimal(10,4)" json:"matrix_intercrystalline,omitempty"`
	Mouldic                *float64 `gorm:"type:decimal(10,4)" json:"mouldic,omitempty"`
	Vuggy                  *float64 `gorm:"type:decimal(10,4)" json:"vuggy,omitempty"`
	Fractures              *float64 `gorm:"type:decimal(10,4)" json:"fractures,omitempty"`
	Fracture               *float64 `gorm:"type:decimal(10,4)" json:"fracture,omitempty"`
	Micro                  *float64 `gorm:"type:decimal(10,4)" json:"micro,omitempty"`

	// Clastic porosity
	Intergranular            *float64 `gorm:"type:decimal(10,4)" json:"intergranular,omitempty"`
	PriPorosityIntragranular *float64 `gorm:"type:decimal(10,4)" json:"pri_porosity_intragranular,omitempty"`
	SecPorosityIntragranular *float64 `gorm:"type:decimal(10,4)" json:"sec_porosity_intragranular,omitempty"`
	Intracrystalline         *float64 `gorm:"type:decimal(10,4)" json:"intracrystalline,omitempty"`

	// Total porosity
	TotalPorosityPercent          *float64 `gorm:"type:decimal(10,4)" json:"total_porosity_percent,omitempty"`
	TotalPrimaryPorosityPercent   *float64 `gorm:"type:decimal(10,4)" json:"total_primary_porosity_percent,omitempty"`
	TotalSecondaryPorosityPercent *float64 `gorm:"type:decimal(10,4)" json:"total_secondary_porosity_percent,omitempty"`

	// ========================================
	// CEMENT TYPES (Carbonate)
	// ========================================
	Fringing           *float64 `gorm:"type:decimal(10,4)" json:"fringing,omitempty"`
	Meniscus           *float64 `gorm:"type:decimal(10,4)" json:"meniscus,omitempty"`
	Blocky             *float64 `gorm:"type:decimal(10,4)" json:"blocky,omitempty"`
	Sparry             *float64 `gorm:"type:decimal(10,4)" json:"sparry,omitempty"`
	Micritic           *float64 `gorm:"type:decimal(10,4)" json:"micritic,omitempty"`
	Pendant            *float64 `gorm:"type:decimal(10,4)" json:"pendant,omitempty"`
	Syntax             *float64 `gorm:"type:decimal(10,4)" json:"syntax,omitempty"`
	CalciteSyntaxial   *float64 `gorm:"type:decimal(10,4)" json:"calcite_syntaxial,omitempty"`
	CalciteFringing    *float64 `gorm:"type:decimal(10,4)" json:"calcite_fringing,omitempty"`
	CalciteMosaic      *float64 `gorm:"type:decimal(10,4)" json:"calcite_mosaic,omitempty"`
	CalciteBlocky      *float64 `gorm:"type:decimal(10,4)" json:"calcite_blocky,omitempty"`
	CalciteFerroan     *float64 `gorm:"type:decimal(10,4)" json:"calcite_ferroan,omitempty"`
	Fluorite           *float64 `gorm:"type:decimal(10,4)" json:"fluorite,omitempty"`
	TotalCementPercent *float64 `gorm:"type:decimal(10,4)" json:"total_cement_percent,omitempty"`

	// ========================================
	// REPLACEMENT & ACCESSORIES (Carbonate)
	// ========================================
	Replacement             *float64 `gorm:"type:decimal(10,4)" json:"replacement,omitempty"`
	Saddle                  *float64 `gorm:"type:decimal(10,4)" json:"saddle,omitempty"`
	TotalDolomitePercent    *float64 `gorm:"type:decimal(10,4)" json:"total_dolomite_percent,omitempty"`
	Stylolite               *float64 `gorm:"type:decimal(10,4)" json:"stylolite,omitempty"`
	Bioturbation            *float64 `gorm:"type:decimal(10,4)" json:"bioturbation,omitempty"`
	TotalAccessoriesPercent *float64 `gorm:"type:decimal(10,4)" json:"total_accessories_percent,omitempty"`

	// ========================================
	// TOTALS
	// ========================================
	TotalPercent *float64 `gorm:"type:decimal(10,4)" json:"total_percent,omitempty"`

	// ========================================
	// ANALYSIS TYPE
	// ========================================
	AnalysisTypes *string `gorm:"type:varchar(255)" json:"analysis_types,omitempty"`

	// ========================================
	// METADATA
	// ========================================
	SourcePDF      *string    `gorm:"type:varchar(500)" json:"source_pdf,omitempty"`
	PageNumber     *int       `gorm:"type:int" json:"page_number,omitempty"`
	TableID        *string    `gorm:"type:varchar(100)" json:"table_id,omitempty"`
	ExtractionDate *time.Time `json:"extraction_date,omitempty"`
	Notes          *string    `gorm:"type:text" json:"notes,omitempty"`
	UnmappedFields *string    `gorm:"type:jsonb" json:"unmapped_fields,omitempty"` // Store any fields that couldn't be mapped
}

// TableName specifies the table name for GORM
func (PetrographyUnified) TableName() string {
	return "petrography_unified"
}
