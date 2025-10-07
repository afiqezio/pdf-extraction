package models

type EPBEPetrographyClastic struct {
	// Base fields - Country, Region, Sub-Region, Business Region, Basin, Sub-Basin, Well Name/Field Name, UWI, Latitude, Longitude, Onshore/Offshore, Water Depth (m), Water Depth (ft)
	EPBEBase

	// Depth Reference and Ground Level
	DepthInfo

	// Textural Analysis (columns 41-46)
	GrainSize            string `json:"grain_size" gorm:"column:grain_size;size:255"`
	GrainShape           string `json:"grain_shape" gorm:"column:grain_shape;size:255"`
	GrainContact         string `json:"grain_contact" gorm:"column:grain_contact;size:255"`
	SedimentaryStructure string `json:"sedimentary_structure" gorm:"column:sedimentary_structure;size:255"`
	Sorting              string `json:"sorting" gorm:"column:sorting;size:255"`
	Lithofacies          string `json:"lithofacies" gorm:"column:lithofacies;size:255"`

	// Porosity and Permeability (columns 47-50)
	VisiblePorosityPercent   *float64 `json:"visible_porosity_percent" gorm:"column:visible_porosity_percent;type:decimal(8,4)"`
	AmbientHePorosityPercent *float64 `json:"ambient_he_porosity_percent" gorm:"column:ambient_he_porosity_percent;type:decimal(8,4)"`
	PermeabilityMd           *float64 `json:"permeability_md" gorm:"column:permeability_md;type:decimal(12,4)"`
	GrainDensityGCc          *float64 `json:"grain_density_g_cc" gorm:"column:grain_density_g_cc;type:decimal(8,4)"`

	// Quartz Content (columns 51-53)
	MonocrystallineQuartz *float64 `json:"monocrystalline_quartz" gorm:"column:monocrystalline_quartz;type:decimal(8,4)"`
	PolycrystallineQuartz *float64 `json:"polycrystalline_quartz" gorm:"column:polycrystalline_quartz;type:decimal(8,4)"`
	TotalQuartzPercent    *float64 `json:"total_quartz_percent" gorm:"column:total_quartz_percent;type:decimal(8,4)"`

	// Feldspar Content (columns 54-57)
	PotassiumFeldspar        *float64 `json:"potassium_feldspar" gorm:"column:potassium_feldspar;type:decimal(8,4)"`
	Plagioclase              *float64 `json:"plagioclase" gorm:"column:plagioclase;type:decimal(8,4)"`
	FeldsparUndifferentiated *float64 `json:"feldspar_undifferentiated" gorm:"column:feldspar_undifferentiated;type:decimal(8,4)"`
	TotalFeldsparPercent     *float64 `json:"total_feldspar_percent" gorm:"column:total_feldspar_percent;type:decimal(8,4)"`

	// Mica Content (columns 58-61)
	Muscovite            *float64 `json:"muscovite" gorm:"column:muscovite;type:decimal(8,4)"`
	Biotite              *float64 `json:"biotite" gorm:"column:biotite;type:decimal(8,4)"`
	MicaUndifferentiated *float64 `json:"mica_undifferentiated" gorm:"column:mica_undifferentiated;type:decimal(8,4)"`
	TotalMicaPercent     *float64 `json:"total_mica_percent" gorm:"column:total_mica_percent;type:decimal(8,4)"`

	// Heavy Minerals (columns 62-65)
	Zircon                        *float64 `json:"zircon" gorm:"column:zircon;type:decimal(8,4)"`
	Tourmaline                    *float64 `json:"tourmaline" gorm:"column:tourmaline;type:decimal(8,4)"`
	HeavyMineralsUndifferentiated *float64 `json:"heavy_minerals_undifferentiated" gorm:"column:heavy_minerals_undifferentiated;type:decimal(8,4)"`
	TotalHeavyMineralsPercent     *float64 `json:"total_heavy_minerals_percent" gorm:"column:total_heavy_minerals_percent;type:decimal(8,4)"`

	// Igneous Rock Fragments (columns 66-69)
	PlutonicRockFragments             *float64 `json:"plutonic_rock_fragments" gorm:"column:plutonic_rock_fragments;type:decimal(8,4)"`
	MaficIntermediateVolcanicFragment *float64 `json:"mafic_intermediate_volcanic_fragment" gorm:"column:mafic_intermediate_volcanic_fragment;type:decimal(8,4)"`
	VolcanicRockFragment              *float64 `json:"volcanic_rock_fragment" gorm:"column:volcanic_rock_fragment;type:decimal(8,4)"`
	TotalIgneousRFPercent             *float64 `json:"total_igneous_rf_percent" gorm:"column:total_igneous_rf_percent;type:decimal(8,4)"`

	// Metamorphic Rock Fragments (columns 70-73)
	QuartzoseRockFragment                   *float64 `json:"quartzose_rock_fragment" gorm:"column:quartzose_rock_fragment;type:decimal(8,4)"`
	SchistoseRockFragment                   *float64 `json:"schistose_rock_fragment" gorm:"column:schistose_rock_fragment;type:decimal(8,4)"`
	MetamorphicRockFragmentUndifferentiated *float64 `json:"metamorphic_rock_fragment_undifferentiated" gorm:"column:metamorphic_rock_fragment_undifferentiated;type:decimal(8,4)"`
	TotalMetamorphicRFPercent               *float64 `json:"total_metamorphic_rf_percent" gorm:"column:total_metamorphic_rf_percent;type:decimal(8,4)"`

	// Sedimentary Rock Fragments (columns 74-82)
	SandstoneSiltstoneRockFragments            *float64 `json:"sandstone_siltstone_rock_fragments" gorm:"column:sandstone_siltstone_rock_fragments;type:decimal(8,4)"`
	ArgillaceousRockFragments                  *float64 `json:"argillaceous_rock_fragments" gorm:"column:argillaceous_rock_fragments;type:decimal(8,4)"`
	SiliciclasticRockFragmentsUndifferentiated *float64 `json:"siliciclastic_rock_fragments_undifferentiated" gorm:"column:siliciclastic_rock_fragments_undifferentiated;type:decimal(8,4)"`
	LimestoneRockFragments                     *float64 `json:"limestone_rock_fragments" gorm:"column:limestone_rock_fragments;type:decimal(8,4)"`
	DolostoneRockFragments                     *float64 `json:"dolostone_rock_fragments" gorm:"column:dolostone_rock_fragments;type:decimal(8,4)"`
	Chert                                      *float64 `json:"chert" gorm:"column:chert;type:decimal(8,4)"`
	TotalSedimentaryRFPercent                  *float64 `json:"total_sedimentary_rf_percent" gorm:"column:total_sedimentary_rf_percent;type:decimal(8,4)"`
	TotalRockFragmentsPercent                  *float64 `json:"total_rock_fragments_percent" gorm:"column:total_rock_fragments_percent;type:decimal(8,4)"`

	// Other Grains (columns 83-88)
	RipUpClast                  *float64 `json:"rip_up_clast" gorm:"column:rip_up_clast;type:decimal(8,4)"`
	Glauconite                  *float64 `json:"glauconite" gorm:"column:glauconite;type:decimal(8,4)"`
	Bioclast                    *float64 `json:"bioclast" gorm:"column:bioclast;type:decimal(8,4)"`
	ForaminiferaGrains          *float64 `json:"foraminifera_grains" gorm:"column:foraminifera_grains;type:decimal(8,4)"`
	UndifferentiatedOtherGrains *float64 `json:"undifferentiated_other_grains" gorm:"column:undifferentiated_other_grains;type:decimal(8,4)"`
	TotalOtherGrainsPercent     *float64 `json:"total_other_grains_percent" gorm:"column:total_other_grains_percent;type:decimal(8,4)"`

	// Matrix (columns 89-94)
	ClayMatrix              *float64 `json:"clay_matrix" gorm:"column:clay_matrix;type:decimal(8,4)"`
	MixedClaySiltFineMatrix *float64 `json:"mixed_clay_silt_fine_matrix" gorm:"column:mixed_clay_silt_fine_matrix;type:decimal(8,4)"`
	SiltVeryFineMatrix      *float64 `json:"silt_very_fine_matrix" gorm:"column:silt_very_fine_matrix;type:decimal(8,4)"`
	OrganicMatrix           *float64 `json:"organic_matrix" gorm:"column:organic_matrix;type:decimal(8,4)"`
	MatrixUndifferentiated  *float64 `json:"matrix_undifferentiated" gorm:"column:matrix_undifferentiated;type:decimal(8,4)"`
	TotalMatrixPercent      *float64 `json:"total_matrix_percent" gorm:"column:total_matrix_percent;type:decimal(8,4)"`

	// Authigenic Clay (columns 95-101)
	Kaolinite                  *float64 `json:"kaolinite" gorm:"column:kaolinite;type:decimal(8,4)"`
	KaoliniteReplacesKFeldspar *float64 `json:"kaolinite_replaces_k_feldspar" gorm:"column:kaolinite_replaces_k_feldspar;type:decimal(8,4)"`
	IllitePoreGrainLining      *float64 `json:"illite_pore_grain_lining" gorm:"column:illite_pore_grain_lining;type:decimal(8,4)"`
	IllitePoreFilling          *float64 `json:"illite_pore_filling" gorm:"column:illite_pore_filling;type:decimal(8,4)"`
	IlliteReplacesKFeldspar    *float64 `json:"illite_replaces_k_feldspar" gorm:"column:illite_replaces_k_feldspar;type:decimal(8,4)"`
	TotalAuthigenicClayPercent *float64 `json:"total_authigenic_clay_percent" gorm:"column:total_authigenic_clay_percent;type:decimal(8,4)"`

	// Authigenic Non-Clay (columns 102-109)
	SyntaxialQuartzOvergrowths    *float64 `json:"syntaxial_quartz_overgrowths" gorm:"column:syntaxial_quartz_overgrowths;type:decimal(8,4)"`
	FeldsparOvergrowths           *float64 `json:"feldspar_overgrowths" gorm:"column:feldspar_overgrowths;type:decimal(8,4)"`
	FeCalcite                     *float64 `json:"fe_calcite" gorm:"column:fe_calcite;type:decimal(8,4)"`
	FeDolomite                    *float64 `json:"fe_dolomite" gorm:"column:fe_dolomite;type:decimal(8,4)"`
	Siderite                      *float64 `json:"siderite" gorm:"column:siderite;type:decimal(8,4)"`
	MnSiderite                    *float64 `json:"mn_siderite" gorm:"column:mn_siderite;type:decimal(8,4)"`
	Pyrite                        *float64 `json:"pyrite" gorm:"column:pyrite;type:decimal(8,4)"`
	IronOxideMinerals             *float64 `json:"iron_oxide_minerals" gorm:"column:iron_oxide_minerals;type:decimal(8,4)"`
	TotalAuthigenicNonClayPercent *float64 `json:"total_authigenic_non_clay_percent" gorm:"column:total_authigenic_non_clay_percent;type:decimal(8,4)"`

	// Primary Porosity (columns 110-113)
	Intergranular               *float64 `json:"intergranular" gorm:"column:intergranular;type:decimal(8,4)"`
	Intercrystalline            *float64 `json:"intercrystalline" gorm:"column:intercrystalline;type:decimal(8,4)"`
	PriPorosityIntragranular    *float64 `json:"pri_porosity_intragranular" gorm:"column:pri_porosity_intragranular;type:decimal(8,4)"`
	TotalPrimaryPorosityPercent *float64 `json:"total_primary_porosity_percent" gorm:"column:total_primary_porosity_percent;type:decimal(8,4)"`

	// Secondary Porosity (columns 114-118)
	SecPorosityIntragranular      *float64 `json:"sec_porosity_intragranular" gorm:"column:sec_porosity_intragranular;type:decimal(8,4)"`
	Intracrystalline              *float64 `json:"intracrystalline" gorm:"column:intracrystalline;type:decimal(8,4)"`
	Mouldic                       *float64 `json:"mouldic" gorm:"column:mouldic;type:decimal(8,4)"`
	Fracture                      *float64 `json:"fracture" gorm:"column:fracture;type:decimal(8,4)"`
	TotalSecondaryPorosityPercent *float64 `json:"total_secondary_porosity_percent" gorm:"column:total_secondary_porosity_percent;type:decimal(8,4)"`

	// Total and Analysis (columns 119-121)
	TotalPercent  *float64 `json:"total_percent" gorm:"column:total_percent;type:decimal(8,4)"`
	AnalysisTypes string   `json:"analysis_types" gorm:"column:analysis_types;size:255"`

	// Metadata fields
	MetadataInfo
}

func (EPBEPetrographyClastic) TableName() string {
	return "petrography_clastic"
}
