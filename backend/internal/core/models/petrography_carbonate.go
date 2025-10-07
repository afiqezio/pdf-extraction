package models

type EPBEPetrographyCarbonate struct {
	// Base fields - Country, Region, Sub-Region, Business Region, Basin, Sub-Basin, Well Name/Field Name, UWI, Latitude, Longitude, Onshore/Offshore, Water Depth (m), Water Depth (ft)
	EPBEBase

	// Depth Reference and Ground Level
	DepthInfo

	// Facies Description (columns 41-43)
	LithofaciesCore        string `json:"lithofacies_core" gorm:"column:lithofacies_core;size:255"`
	MicrofaciesThinSection string `json:"microfacies_thin_section" gorm:"column:microfacies_thin_section;size:255"`
	Depofacies             string `json:"depofacies" gorm:"column:depofacies;size:255"`

	// Porosity and Permeability (columns 44-46)
	VisiblePorosityPercent *float64 `json:"visible_porosity_percent" gorm:"column:visible_porosity_percent;type:decimal(8,4)"`
	HePorosityPercent      *float64 `json:"he_porosity_percent" gorm:"column:he_porosity_percent;type:decimal(8,4)"`
	PermeabilityMd         *float64 `json:"permeability_md" gorm:"column:permeability_md;type:decimal(12,4)"`

	// Matrix Mineralogy (columns 47-52)
	Calcite                      *float64 `json:"calcite" gorm:"column:calcite;type:decimal(8,4)"`
	Dolomite                     *float64 `json:"dolomite" gorm:"column:dolomite;type:decimal(8,4)"`
	Micrite                      *float64 `json:"micrite" gorm:"column:micrite;type:decimal(8,4)"`
	MicriteEnvelopes             *float64 `json:"micrite_envelopes" gorm:"column:micrite_envelopes;type:decimal(8,4)"`
	MicrosparPseudospar          *float64 `json:"microspar_pseudospar" gorm:"column:microspar_pseudospar;type:decimal(8,4)"`
	Kaolinite                    *float64 `json:"kaolinite" gorm:"column:kaolinite;type:decimal(8,4)"`
	Clay                         *float64 `json:"clay" gorm:"column:clay;type:decimal(8,4)"`
	TotalMineralogyMatrixPercent *float64 `json:"total_mineralogy_matrix_percent" gorm:"column:total_mineralogy_matrix_percent;type:decimal(8,4)"`

	// Bioclasts - General and Specific Foraminifera (columns 53-82)
	Bioclasts            *float64 `json:"bioclasts" gorm:"column:bioclasts;type:decimal(8,4)"`
	Lepido               *float64 `json:"lepido" gorm:"column:lepido;type:decimal(8,4)"`
	Coral                *float64 `json:"coral" gorm:"column:coral;type:decimal(8,4)"`
	Rhodolith            *float64 `json:"rhodolith" gorm:"column:rhodolith;type:decimal(8,4)"`
	RedAlgae             *float64 `json:"red_algae" gorm:"column:red_algae;type:decimal(8,4)"`
	RedAlgaeEnc          *float64 `json:"red_algae_enc" gorm:"column:red_algae_enc;type:decimal(8,4)"`
	GreenAlgae           *float64 `json:"green_algae" gorm:"column:green_algae;type:decimal(8,4)"`
	Echinoderms          *float64 `json:"echinoderms" gorm:"column:echinoderms;type:decimal(8,4)"`
	Miliolid             *float64 `json:"miliolid" gorm:"column:miliolid;type:decimal(8,4)"`
	Lepidocyclina        *float64 `json:"lepidocyclina" gorm:"column:lepidocyclina;type:decimal(8,4)"`
	Cycloclypeus         *float64 `json:"cycloclypeus" gorm:"column:cycloclypeus;type:decimal(8,4)"`
	Operculina           *float64 `json:"operculina" gorm:"column:operculina;type:decimal(8,4)"`
	OtherRotaliids       *float64 `json:"other_rotaliids" gorm:"column:other_rotaliids;type:decimal(8,4)"`
	Gypsinid             *float64 `json:"gypsinid" gorm:"column:gypsinid;type:decimal(8,4)"`
	Planorbulinella      *float64 `json:"planorbulinella" gorm:"column:planorbulinella;type:decimal(8,4)"`
	Hemotremid           *float64 `json:"hemotremid" gorm:"column:hemotremid;type:decimal(8,4)"`
	Heterostegina        *float64 `json:"heterostegina" gorm:"column:heterostegina;type:decimal(8,4)"`
	EncFrm               *float64 `json:"enc_frm" gorm:"column:enc_frm;type:decimal(8,4)"`
	Planktonic           *float64 `json:"planktonic" gorm:"column:planktonic;type:decimal(8,4)"`
	Bryozoans            *float64 `json:"bryozoans" gorm:"column:bryozoans;type:decimal(8,4)"`
	Amphistegina         *float64 `json:"amphistegina" gorm:"column:amphistegina;type:decimal(8,4)"`
	Gastropods           *float64 `json:"gastropods" gorm:"column:gastropods;type:decimal(8,4)"`
	Bivalve              *float64 `json:"bivalve" gorm:"column:bivalve;type:decimal(8,4)"`
	Ostracod             *float64 `json:"ostracod" gorm:"column:ostracod;type:decimal(8,4)"`
	Oncoids              *float64 `json:"oncoids" gorm:"column:oncoids;type:decimal(8,4)"`
	UndiffMolluscs       *float64 `json:"undiff_molluscs" gorm:"column:undiff_molluscs;type:decimal(8,4)"`
	UndiffBenthonic      *float64 `json:"undiff_benthonic" gorm:"column:undiff_benthonic;type:decimal(8,4)"`
	UndiffSkeletal       *float64 `json:"undiff_skeletal" gorm:"column:undiff_skeletal;type:decimal(8,4)"`
	UndiffForam          *float64 `json:"undiff_foram" gorm:"column:undiff_foram;type:decimal(8,4)"`
	TotalSkeletalPercent *float64 `json:"total_skeletal_percent" gorm:"column:total_skeletal_percent;type:decimal(8,4)"`

	// Non-Skeletal Components (columns 83-88)
	Organic                 *float64 `json:"organic" gorm:"column:organic;type:decimal(8,4)"`
	Peloids                 *float64 `json:"peloids" gorm:"column:peloids;type:decimal(8,4)"`
	MicritisedGrains        *float64 `json:"micritised_grains" gorm:"column:micritised_grains;type:decimal(8,4)"`
	Pseudoclasts            *float64 `json:"pseudoclasts" gorm:"column:pseudoclasts;type:decimal(8,4)"`
	Intraclast              *float64 `json:"intraclast" gorm:"column:intraclast;type:decimal(8,4)"`
	Quartz                  *float64 `json:"quartz" gorm:"column:quartz;type:decimal(8,4)"`
	TotalNonSkeletalPercent *float64 `json:"total_non_skeletal_percent" gorm:"column:total_non_skeletal_percent;type:decimal(8,4)"`

	// Porosity Types (columns 89-97)
	Interparticle          *float64 `json:"interparticle" gorm:"column:interparticle;type:decimal(8,4)"`
	Intraparticle          *float64 `json:"intraparticle" gorm:"column:intraparticle;type:decimal(8,4)"`
	Intercrystalline       *float64 `json:"intercrystalline" gorm:"column:intercrystalline;type:decimal(8,4)"`
	MatrixIntercrystalline *float64 `json:"matrix_intercrystalline" gorm:"column:matrix_intercrystalline;type:decimal(8,4)"`
	Mouldic                *float64 `json:"mouldic" gorm:"column:mouldic;type:decimal(8,4)"`
	Vuggy                  *float64 `json:"vuggy" gorm:"column:vuggy;type:decimal(8,4)"`
	Fractures              *float64 `json:"fractures" gorm:"column:fractures;type:decimal(8,4)"`
	Micro                  *float64 `json:"micro" gorm:"column:micro;type:decimal(8,4)"`
	TotalPorosityPercent   *float64 `json:"total_porosity_percent" gorm:"column:total_porosity_percent;type:decimal(8,4)"`

	// Cement Types (columns 98-121)
	Fringing           *float64 `json:"fringing" gorm:"column:fringing;type:decimal(8,4)"`
	Meniscus           *float64 `json:"meniscus" gorm:"column:meniscus;type:decimal(8,4)"`
	Blocky             *float64 `json:"blocky" gorm:"column:blocky;type:decimal(8,4)"`
	Sparry             *float64 `json:"sparry" gorm:"column:sparry;type:decimal(8,4)"`
	Micritic           *float64 `json:"micritic" gorm:"column:micritic;type:decimal(8,4)"`
	Pendant            *float64 `json:"pendant" gorm:"column:pendant;type:decimal(8,4)"`
	Syntax             *float64 `json:"syntax" gorm:"column:syntax;type:decimal(8,4)"`
	CalciteSyntaxial   *float64 `json:"calcite_syntaxial" gorm:"column:calcite_syntaxial;type:decimal(8,4)"`
	CalciteFringing    *float64 `json:"calcite_fringing" gorm:"column:calcite_fringing;type:decimal(8,4)"`
	CalciteMosaic      *float64 `json:"calcite_mosaic" gorm:"column:calcite_mosaic;type:decimal(8,4)"`
	CalciteBlocky      *float64 `json:"calcite_blocky" gorm:"column:calcite_blocky;type:decimal(8,4)"`
	CalciteFerroan     *float64 `json:"calcite_ferroan" gorm:"column:calcite_ferroan;type:decimal(8,4)"`
	Pyrite             *float64 `json:"pyrite" gorm:"column:pyrite;type:decimal(8,4)"`
	Fluorite           *float64 `json:"fluorite" gorm:"column:fluorite;type:decimal(8,4)"`
	TotalCementPercent *float64 `json:"total_cement_percent" gorm:"column:total_cement_percent;type:decimal(8,4)"`

	// Replacement and Accessories (columns 122-128)
	Replacement             *float64 `json:"replacement" gorm:"column:replacement;type:decimal(8,4)"`
	Saddle                  *float64 `json:"saddle" gorm:"column:saddle;type:decimal(8,4)"`
	TotalDolomitePercent    *float64 `json:"total_dolomite_percent" gorm:"column:total_dolomite_percent;type:decimal(8,4)"`
	Stylolite               *float64 `json:"stylolite" gorm:"column:stylolite;type:decimal(8,4)"`
	Bioturbation            *float64 `json:"bioturbation" gorm:"column:bioturbation;type:decimal(8,4)"`
	TotalAccessoriesPercent *float64 `json:"total_accessories_percent" gorm:"column:total_accessories_percent;type:decimal(8,4)"`
	TotalPercent            *float64 `json:"total_percent" gorm:"column:total_percent;type:decimal(8,4)"`

	// Analysis Type
	AnalysisTypes string `json:"analysis_types" gorm:"column:analysis_types;size:255"`

	// Metadata fields
	MetadataInfo
}

func (EPBEPetrographyCarbonate) TableName() string {
	return "petrography_carbonate"
}
