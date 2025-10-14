package handlers

// GetFieldMappings returns the comprehensive field mappings for database field matching
func GetFieldMappings() map[string]string {
	return map[string]string{
		// Well and location fields
		"well_name_field_name": "well_name_field_name", "well name field name": "well_name_field_name", "well_name": "well_name_field_name", "well name": "well_name_field_name", "well": "well_name_field_name", "well_id": "well_name_field_name", "well id": "well_name_field_name", "field_name": "well_name_field_name", "field name": "well_name_field_name", "field": "well_name_field_name",
		"country": "country", "country_name": "country", "country name": "country",
		"region": "region", "region_name": "region", "region name": "region", "sub_region": "sub_region", "sub region": "sub_region", "subregion": "sub_region",
		"business_regions": "business_regions", "business regions": "business_regions", "business_region": "business_regions", "business region": "business_regions",
		"basin": "basin", "basin_name": "basin", "basin name": "basin", "sub_basin": "sub_basin", "sub basin": "sub_basin", "subbasin": "sub_basin",
		"uwi": "uwi", "unique_well_identifier": "uwi", "unique well identifier": "uwi", "well_identifier": "uwi", "well identifier": "uwi",
		
		// Additional geological field mappings from console output
		"depth (m)": "top_depth_mmddf", "depth_m": "top_depth_mmddf", "depth m": "top_depth_mmddf",
		"botdepth_m": "bottom_depth_mmddf", "bot depth m": "bottom_depth_mmddf", "bottom depth m": "bottom_depth_mmddf",
		"classification": "lithofacies_core", "lithologic summary": "lithofacies_core", "lithologic": "lithofacies_core",
		"plagioclase feldspar": "plagioclase", "plag": "plagioclase",
		"feldspar, microcline": "potassium_feldspar", "microcline": "potassium_feldspar", "k-feldspar": "potassium_feldspar",
		"untwinned feldspar": "feldspar_undifferentiated", "feldspar undiff": "feldspar_undifferentiated",
		"chert %": "chert",
		"clay matrix, undiff.": "clay_matrix", "clay matrix undiff": "clay_matrix",
		"siliceous matrix": "silt_very_fine_matrix", "siliceous": "silt_very_fine_matrix",
		"carbonate matrix": "carbonate_matrix",
		"organic matrix": "organic_matrix",
		"siderite cement": "siderite",
		"iron oxide cement": "iron_oxide_minerals",
		"chlorite cement": "chlorite",
		"kaolinite cement": "kaolinite",
		"rock frag., undiff": "siliciclastic_rock_fragments_undifferentiated", "rock fragments undiff": "siliciclastic_rock_fragments_undifferentiated",
		"rip-up clasts": "rip_up_clast", "rip up clasts": "rip_up_clast",
		"micas, undiff": "mica_undifferentiated", "micas undiff": "mica_undifferentiated", "micas": "mica_undifferentiated",
		"mica muscovite": "muscovite",
		"mica biotite": "biotite",
		"fossils\\phosphatic": "bioclast", "fossils phosphatic": "bioclast", "fossils": "bioclast",
		"secondary pores": "total_secondary_porosity_percent", "secondary porosity": "total_secondary_porosity_percent",
		"replacement, undiff": "replacement", "replacement undiff": "replacement",
		
		// Coordinates
		"lat": "latitude", "latitude": "latitude", "lat.": "latitude", "north": "latitude", "north_lat": "latitude", "n_lat": "latitude", "y_coord": "latitude", "y_coordinate": "latitude",
		"lon": "longitude", "longitude": "longitude", "long.": "longitude", "east": "longitude", "east_lon": "longitude", "e_lon": "longitude", "x_coord": "longitude", "x_coordinate": "longitude",
		"coord": "latitude", "coordinate": "latitude", "coordinates": "latitude", "position": "latitude",
		"x": "longitude", "y": "latitude", "easting": "longitude", "northing": "latitude",
		
		// Geological information
		"formation": "formation_name", "form": "formation_name", "unit": "formation_name", "formation_name": "formation_name", "formation name": "formation_name", "geologic_formation": "formation_name",
		"reservoir": "reservoir_name", "reservoir name": "reservoir_name", "reservoir_name": "reservoir_name", "reservoir_field": "reservoir_name",
		"period": "period", "geological period": "period", "geo_period": "period", "time_period": "period",
		"epoch": "epoch", "geological epoch": "epoch", "geo_epoch": "epoch",
		"age": "age", "geological age": "age", "geo_age": "age", "ma": "age", "million years": "age",
		"onshore_offshore": "onshore_offshore", "onshore offshore": "onshore_offshore", "location": "onshore_offshore", "onshore": "onshore_offshore", "offshore": "onshore_offshore",
		
		// Water depth
		"water_depth": "water_depth_m", "water depth": "water_depth_m", "wd": "water_depth_m", "wd_m": "water_depth_m", "water_depth_m": "water_depth_m",
		"water_depth_ft": "water_depth_ft", "water depth ft": "water_depth_ft", "wd_ft": "water_depth_ft", "water_depth_feet": "water_depth_ft",
		
		// Depth information
		"top_depth": "top_depth_mmddf", "top depth": "top_depth_mmddf", "td": "top_depth_mmddf", "top": "top_depth_mmddf", "start_depth": "top_depth_mmddf",
		"bottom_depth": "bottom_depth_mmddf", "bottom depth": "bottom_depth_mmddf", "bd": "bottom_depth_mmddf", "bottom": "bottom_depth_mmddf", "end_depth": "bottom_depth_mmddf",
		"depth": "top_depth_mmddf", "md": "top_depth_mmddf", "measured_depth": "top_depth_mmddf", "measured depth": "top_depth_mmddf",
		"tvd": "top_depth_mtvddf", "true_vertical_depth": "top_depth_mtvddf", "true vertical depth": "top_depth_mtvddf",
		"tvdss": "top_depth_mtvdss", "tvd_ss": "top_depth_mtvdss", "true_vertical_depth_ss": "top_depth_mtvdss",
		"bml": "top_depth_mbml", "below_mudline": "top_depth_mbml", "below mudline": "top_depth_mbml",
		
		// Lithofacies
		"lithofacies": "lithofacies_core", "lithofacies_core": "lithofacies_core", "lithofacies core": "lithofacies_core", "facies": "lithofacies_core", "lithology": "lithofacies_core",
		"microfacies": "microfacies_thin_section", "microfacies_thin_section": "microfacies_thin_section", "microfacies thin section": "microfacies_thin_section",
		"depofacies": "depofacies", "depositional_facies": "depofacies", "depositional facies": "depofacies",
		
		// Analysis types
		"analysis": "analysis_types", "analysis_types": "analysis_types", "analysis types": "analysis_types", "method": "analysis_types", "technique": "analysis_types",
		
		// Porosity and permeability
		"porosity": "visible_porosity_percent", "visible_porosity": "visible_porosity_percent", "visible porosity": "visible_porosity_percent", "phi": "visible_porosity_percent", "phie": "visible_porosity_percent",
		"he_porosity": "he_porosity_percent", "he porosity": "he_porosity_percent", "helium_porosity": "he_porosity_percent", "helium porosity": "he_porosity_percent",
		"permeability": "permeability_md", "perm": "permeability_md", "k": "permeability_md", "kh": "permeability_md", "kv": "permeability_md",
		"grain_density": "grain_density_g_cc", "density": "grain_density_g_cc", "rhob": "grain_density_g_cc",
		
		// Carbonate specific - Matrix mineralogy
		"calcite": "calcite", "calcite_percent": "calcite", "calcite %": "calcite", "cal": "calcite",
		"dolomite": "dolomite", "dolomite_percent": "dolomite", "dolomite %": "dolomite", "dol": "dolomite",
		"micrite": "micrite", "micrite_percent": "micrite", "micrite %": "micrite", "mic": "micrite",
		"micrite_envelopes": "micrite_envelopes", "micrite envelopes": "micrite_envelopes", "micrite_env": "micrite_envelopes",
		"microspar": "microspar_pseudospar", "microspar_pseudospar": "microspar_pseudospar", "microspar pseudospar": "microspar_pseudospar", "pseudospar": "microspar_pseudospar",
		"kaolinite": "kaolinite", "kaolinite_percent": "kaolinite", "kaolinite %": "kaolinite", "kaol": "kaolinite",
		"clay": "clay", "clay_percent": "clay", "clay %": "clay",
		"total_matrix": "total_mineralogy_matrix_percent", "total matrix": "total_mineralogy_matrix_percent", "total_matrix_percent": "total_mineralogy_matrix_percent", "total matrix percent": "total_mineralogy_matrix_percent",
		
		// Carbonate specific - Bioclasts
		"bioclasts": "bioclasts", "bioclast": "bioclasts", "bioclast_percent": "bioclasts", "bioclast %": "bioclasts", "skeletal": "bioclasts",
		"lepido": "lepido", "lepidocyclina": "lepido", "lepidocyclina_percent": "lepido", "lepidocyclina %": "lepido",
		"coral": "coral", "coral_percent": "coral", "coral %": "coral",
		"rhodolith": "rhodolith", "rhodolith_percent": "rhodolith", "rhodolith %": "rhodolith",
		"red_algae": "red_algae", "red algae": "red_algae", "red_algae_percent": "red_algae", "red algae percent": "red_algae",
		"red_algae_enc": "red_algae_enc", "red algae enc": "red_algae_enc", "red_algae_encrusting": "red_algae_enc",
		"green_algae": "green_algae", "green algae": "green_algae", "green_algae_percent": "green_algae", "green algae percent": "green_algae",
		"echinoderms": "echinoderms", "echinoderm": "echinoderms", "echinoderm_percent": "echinoderms", "echinoderm %": "echinoderms",
		"miliolid": "miliolid", "miliolid_percent": "miliolid", "miliolid %": "miliolid",
		"cycloclypeus": "cycloclypeus", "cycloclypeus_percent": "cycloclypeus", "cycloclypeus %": "cycloclypeus",
		"operculina": "operculina", "operculina_percent": "operculina", "operculina %": "operculina",
		"other_rotaliids": "other_rotaliids", "other rotaliids": "other_rotaliids", "rotaliids": "other_rotaliids",
		"gypsinid": "gypsinid", "gypsinid_percent": "gypsinid", "gypsinid %": "gypsinid",
		"planorbulinella": "planorbulinella", "planorbulinella_percent": "planorbulinella", "planorbulinella %": "planorbulinella",
		"hemotremid": "hemotremid", "hemotremid_percent": "hemotremid", "hemotremid %": "hemotremid",
		"heterostegina": "heterostegina", "heterostegina_percent": "heterostegina", "heterostegina %": "heterostegina",
		"enc_frm": "enc_frm", "encrusting_foram": "enc_frm", "encrusting foram": "enc_frm",
		"planktonic": "planktonic", "planktonic_percent": "planktonic", "planktonic %": "planktonic",
		"bryozoans": "bryozoans", "bryozoan": "bryozoans", "bryozoan_percent": "bryozoans", "bryozoan %": "bryozoans",
		"amphistegina": "amphistegina", "amphistegina_percent": "amphistegina", "amphistegina %": "amphistegina",
		"gastropods": "gastropods", "gastropod": "gastropods", "gastropod_percent": "gastropods", "gastropod %": "gastropods",
		"bivalve": "bivalve", "bivalve_percent": "bivalve", "bivalve %": "bivalve",
		"ostracod": "ostracod", "ostracod_percent": "ostracod", "ostracod %": "ostracod",
		"oncoids": "oncoids", "oncoid": "oncoids", "oncoid_percent": "oncoids", "oncoid %": "oncoids",
		"undiff_molluscs": "undiff_molluscs", "undiff molluscs": "undiff_molluscs", "molluscs": "undiff_molluscs",
		"undiff_benthonic": "undiff_benthonic", "undiff benthonic": "undiff_benthonic", "benthonic": "undiff_benthonic",
		"undiff_skeletal": "undiff_skeletal", "undiff skeletal": "undiff_skeletal",
		"undiff_foram": "undiff_foram", "undiff foram": "undiff_foram", "foram": "undiff_foram",
		"total_skeletal": "total_skeletal_percent", "total skeletal": "total_skeletal_percent", "total_skeletal_percent": "total_skeletal_percent", "total skeletal percent": "total_skeletal_percent",
		
		// Carbonate specific - Non-skeletal components
		"organic": "organic", "organic_percent": "organic", "organic %": "organic",
		"peloids": "peloids", "peloid": "peloids", "peloid_percent": "peloids", "peloid %": "peloids",
		"micritised_grains": "micritised_grains", "micritised grains": "micritised_grains", "micritized_grains": "micritised_grains",
		"pseudoclasts": "pseudoclasts", "pseudoclast": "pseudoclasts", "pseudoclast_percent": "pseudoclasts", "pseudoclast %": "pseudoclasts",
		"intraclast": "intraclast", "intraclast_percent": "intraclast", "intraclast %": "intraclast",
		"quartz": "quartz", "quartz_percent": "quartz", "quartz %": "quartz", "qtz": "quartz",
		"total_non_skeletal": "total_non_skeletal_percent", "total non skeletal": "total_non_skeletal_percent", "total_non_skeletal_percent": "total_non_skeletal_percent",
		
		// Carbonate specific - Porosity types
		"interparticle": "interparticle", "interparticle_percent": "interparticle", "interparticle %": "interparticle", "inter_particle": "interparticle",
		"intraparticle": "intraparticle", "intraparticle_percent": "intraparticle", "intraparticle %": "intraparticle", "intra_particle": "intraparticle",
		"intercrystalline": "intercrystalline", "inter_crystalline": "intercrystalline", "intercrystalline_percent": "intercrystalline", "intercrystalline %": "intercrystalline",
		"matrix_intercrystalline": "matrix_intercrystalline", "matrix intercrystalline": "matrix_intercrystalline", "matrix_inter_crystalline": "matrix_intercrystalline",
		"mouldic": "mouldic", "moldic": "mouldic", "mouldic_percent": "mouldic", "mouldic %": "mouldic",
		"vuggy": "vuggy", "vuggy_percent": "vuggy", "vuggy %": "vuggy", "vugs": "vuggy",
		"fracture": "fractures", "fractures": "fractures", "fracture_percent": "fractures", "fracture %": "fractures",
		"micro": "micro", "micro_percent": "micro", "micro %": "micro", "micro_porosity": "micro",
		"total_porosity": "total_porosity_percent", "total porosity": "total_porosity_percent", "total_porosity_percent": "total_porosity_percent",
		
		// Carbonate specific - Cement types
		"fringing": "fringing", "fringing_percent": "fringing", "fringing %": "fringing",
		"meniscus": "meniscus", "meniscus_percent": "meniscus", "meniscus %": "meniscus",
		"blocky": "blocky", "blocky_percent": "blocky", "blocky %": "blocky",
		"sparry": "sparry", "sparry_percent": "sparry", "sparry %": "sparry",
		"micritic": "micritic", "micritic_percent": "micritic", "micritic %": "micritic",
		"pendant": "pendant", "pendant_percent": "pendant", "pendant %": "pendant",
		"syntax": "syntax", "syntax_percent": "syntax", "syntax %": "syntax", "syntaxial": "syntax",
		"calcite_syntaxial": "calcite_syntaxial", "calcite syntaxial": "calcite_syntaxial", "calcite_syntax": "calcite_syntaxial",
		"calcite_fringing": "calcite_fringing", "calcite fringing": "calcite_fringing",
		"calcite_mosaic": "calcite_mosaic", "calcite mosaic": "calcite_mosaic",
		"calcite_blocky": "calcite_blocky", "calcite blocky": "calcite_blocky",
		"calcite_ferroan": "calcite_ferroan", "calcite ferroan": "calcite_ferroan", "ferroan_calcite": "calcite_ferroan",
		"pyrite": "pyrite", "pyrite_percent": "pyrite", "pyrite %": "pyrite",
		"fluorite": "fluorite", "fluorite_percent": "fluorite", "fluorite %": "fluorite",
		"total_cement": "total_cement_percent", "total cement": "total_cement_percent", "total_cement_percent": "total_cement_percent",
		
		// Carbonate specific - Replacement and accessories
		"replacement": "replacement", "replacement_percent": "replacement", "replacement %": "replacement",
		"saddle": "saddle", "saddle_percent": "saddle", "saddle %": "saddle", "saddle_dolomite": "saddle",
		"total_dolomite": "total_dolomite_percent", "total dolomite": "total_dolomite_percent", "total_dolomite_percent": "total_dolomite_percent",
		"stylolite": "stylolite", "stylolite_percent": "stylolite", "stylolite %": "stylolite",
		"bioturbation": "bioturbation", "bioturbation_percent": "bioturbation", "bioturbation %": "bioturbation",
		"total_accessories": "total_accessories_percent", "total accessories": "total_accessories_percent", "total_accessories_percent": "total_accessories_percent",
		"total_percent": "total_percent", "total percent": "total_percent", "total": "total_percent",
		
		// Clastic specific - Textural Analysis
		"grain_size": "grain_size", "grain size": "grain_size", "size": "grain_size", "gs": "grain_size",
		"grain_shape": "grain_shape", "grain shape": "grain_shape", "shape": "grain_shape", "grain_form": "grain_shape",
		"grain_contact": "grain_contact", "grain contact": "grain_contact", "contact": "grain_contact", "contacts": "grain_contact",
		"sedimentary_structure": "sedimentary_structure", "sedimentary structure": "sedimentary_structure", "structure": "sedimentary_structure",
		"sorting": "sorting", "sort": "sorting",
		
		// Clastic specific - Porosity
		"ambient_he_porosity_percent": "ambient_he_porosity_percent", "ambient he porosity": "ambient_he_porosity_percent", "ambient_he_porosity": "ambient_he_porosity_percent", "ambient he porosity percent": "ambient_he_porosity_percent", "ambient_he_porosity %": "ambient_he_porosity_percent",
		"grain_density_g_cc": "grain_density_g_cc", "grain density g/cc": "grain_density_g_cc",
		
		// Clastic specific - Quartz Content
		"monocrystalline_quartz": "monocrystalline_quartz", "monocrystalline quartz": "monocrystalline_quartz", "mono_quartz": "monocrystalline_quartz", "mono quartz": "monocrystalline_quartz", "monocrystalline_quartz_percent": "monocrystalline_quartz", "monocrystalline quartz percent": "monocrystalline_quartz",
		"polycrystalline_quartz": "polycrystalline_quartz", "polycrystalline quartz": "polycrystalline_quartz", "poly_quartz": "polycrystalline_quartz", "poly quartz": "polycrystalline_quartz", "polycrystalline_quartz_percent": "polycrystalline_quartz", "polycrystalline quartz percent": "polycrystalline_quartz",
		"total_quartz_percent": "total_quartz_percent", "total quartz percent": "total_quartz_percent", "total_quartz": "total_quartz_percent", "total quartz": "total_quartz_percent", "total_quartz %": "total_quartz_percent", "total quartz %": "total_quartz_percent",
		
		// Clastic specific - Feldspar Content
		"potassium_feldspar": "potassium_feldspar", "potassium feldspar": "potassium_feldspar", "k_feldspar": "potassium_feldspar", "k feldspar": "potassium_feldspar", "k_fsp": "potassium_feldspar", "k fsp": "potassium_feldspar", "potassium_feldspar_percent": "potassium_feldspar", "potassium feldspar percent": "potassium_feldspar",
		"plagioclase": "plagioclase", "plagioclase_percent": "plagioclase", "plagioclase %": "plagioclase",
		"feldspar_undifferentiated": "feldspar_undifferentiated", "feldspar undifferentiated": "feldspar_undifferentiated", "undiff_feldspar": "feldspar_undifferentiated", "undiff feldspar": "feldspar_undifferentiated", "feldspar_undifferentiated_percent": "feldspar_undifferentiated", "feldspar undifferentiated percent": "feldspar_undifferentiated",
		"total_feldspar_percent": "total_feldspar_percent", "total feldspar percent": "total_feldspar_percent", "total_feldspar": "total_feldspar_percent", "total feldspar": "total_feldspar_percent", "total_feldspar %": "total_feldspar_percent", "total feldspar %": "total_feldspar_percent",
		
		// Clastic specific - Mica Content
		"muscovite": "muscovite", "muscovite_percent": "muscovite", "muscovite %": "muscovite",
		"biotite": "biotite", "biotite_percent": "biotite", "biotite %": "biotite",
		"mica_undifferentiated": "mica_undifferentiated", "mica undifferentiated": "mica_undifferentiated", "undiff_mica": "mica_undifferentiated", "undiff mica": "mica_undifferentiated", "mica_undifferentiated_percent": "mica_undifferentiated", "mica undifferentiated percent": "mica_undifferentiated",
		"total_mica_percent": "total_mica_percent", "total mica percent": "total_mica_percent", "total_mica": "total_mica_percent", "total mica": "total_mica_percent", "total_mica %": "total_mica_percent", "total mica %": "total_mica_percent",
		
		// Clastic specific - Heavy Minerals
		"zircon": "zircon", "zircon_percent": "zircon", "zircon %": "zircon",
		"tourmaline": "tourmaline", "tourmaline_percent": "tourmaline", "tourmaline %": "tourmaline",
		"heavy_minerals_undifferentiated": "heavy_minerals_undifferentiated", "heavy minerals undifferentiated": "heavy_minerals_undifferentiated", "undiff_heavy_minerals": "heavy_minerals_undifferentiated", "undiff heavy minerals": "heavy_minerals_undifferentiated", "heavy_minerals_undifferentiated_percent": "heavy_minerals_undifferentiated",
		"total_heavy_minerals_percent": "total_heavy_minerals_percent", "total heavy minerals percent": "total_heavy_minerals_percent", "total_heavy_minerals": "total_heavy_minerals_percent", "total heavy minerals": "total_heavy_minerals_percent", "total_heavy_minerals %": "total_heavy_minerals_percent",
		
		// Clastic specific - Igneous Rock Fragments
		"plutonic_rock_fragments": "plutonic_rock_fragments", "plutonic rock fragments": "plutonic_rock_fragments", "plutonic_rf": "plutonic_rock_fragments", "plutonic rf": "plutonic_rock_fragments", "plutonic_rock_fragments_percent": "plutonic_rock_fragments", "plutonic rock fragments percent": "plutonic_rock_fragments",
		"mafic_intermediate_volcanic_fragment": "mafic_intermediate_volcanic_fragment", "mafic intermediate volcanic fragment": "mafic_intermediate_volcanic_fragment", "mafic_volcanic": "mafic_intermediate_volcanic_fragment", "mafic volcanic": "mafic_intermediate_volcanic_fragment", "mafic_intermediate_volcanic_fragment_percent": "mafic_intermediate_volcanic_fragment",
		"volcanic_rock_fragment": "volcanic_rock_fragment", "volcanic rock fragment": "volcanic_rock_fragment", "volcanic_rf": "volcanic_rock_fragment", "volcanic rf": "volcanic_rock_fragment", "volcanic_rock_fragment_percent": "volcanic_rock_fragment", "volcanic rock fragment percent": "volcanic_rock_fragment",
		"total_igneous_rf_percent": "total_igneous_rf_percent", "total igneous rf percent": "total_igneous_rf_percent", "total_igneous_rf": "total_igneous_rf_percent", "total igneous rf": "total_igneous_rf_percent", "total_igneous_rf %": "total_igneous_rf_percent",
		
		// Clastic specific - Metamorphic Rock Fragments
		"quartzose_rock_fragment": "quartzose_rock_fragment", "quartzose rock fragment": "quartzose_rock_fragment", "quartzose_rf": "quartzose_rock_fragment", "quartzose rf": "quartzose_rock_fragment", "quartzose_rock_fragment_percent": "quartzose_rock_fragment", "quartzose rock fragment percent": "quartzose_rock_fragment",
		"schistose_rock_fragment": "schistose_rock_fragment", "schistose rock fragment": "schistose_rock_fragment", "schistose_rf": "schistose_rock_fragment", "schistose rf": "schistose_rock_fragment", "schistose_rock_fragment_percent": "schistose_rock_fragment", "schistose rock fragment percent": "schistose_rock_fragment",
		"metamorphic_rock_fragment_undifferentiated": "metamorphic_rock_fragment_undifferentiated", "metamorphic rock fragment undifferentiated": "metamorphic_rock_fragment_undifferentiated", "meta_rf_undiff": "metamorphic_rock_fragment_undifferentiated", "meta rf undiff": "metamorphic_rock_fragment_undifferentiated", "metamorphic_rock_fragment_undifferentiated_percent": "metamorphic_rock_fragment_undifferentiated",
		"total_metamorphic_rf_percent": "total_metamorphic_rf_percent", "total metamorphic rf percent": "total_metamorphic_rf_percent", "total_metamorphic_rf": "total_metamorphic_rf_percent", "total metamorphic rf": "total_metamorphic_rf_percent", "total_metamorphic_rf %": "total_metamorphic_rf_percent",
		
		// Clastic specific - Sedimentary Rock Fragments
		"sandstone_siltstone_rock_fragments": "sandstone_siltstone_rock_fragments", "sandstone siltstone rock fragments": "sandstone_siltstone_rock_fragments", "ss_silt_rf": "sandstone_siltstone_rock_fragments", "ss silt rf": "sandstone_siltstone_rock_fragments", "sandstone_siltstone_rock_fragments_percent": "sandstone_siltstone_rock_fragments",
		"argillaceous_rock_fragments": "argillaceous_rock_fragments", "argillaceous rock fragments": "argillaceous_rock_fragments", "argillaceous_rf": "argillaceous_rock_fragments", "argillaceous rf": "argillaceous_rock_fragments", "argillaceous_rock_fragments_percent": "argillaceous_rock_fragments",
		"siliciclastic_rock_fragments_undifferentiated": "siliciclastic_rock_fragments_undifferentiated", "siliciclastic rock fragments undifferentiated": "siliciclastic_rock_fragments_undifferentiated", "siliciclastic_rf_undiff": "siliciclastic_rock_fragments_undifferentiated", "siliciclastic rf undiff": "siliciclastic_rock_fragments_undifferentiated", "siliciclastic_rock_fragments_undifferentiated_percent": "siliciclastic_rock_fragments_undifferentiated",
		"limestone_rock_fragments": "limestone_rock_fragments", "limestone rock fragments": "limestone_rock_fragments", "limestone_rf": "limestone_rock_fragments", "limestone rf": "limestone_rock_fragments", "limestone_rock_fragments_percent": "limestone_rock_fragments",
		"dolostone_rock_fragments": "dolostone_rock_fragments", "dolostone rock fragments": "dolostone_rock_fragments", "dolostone_rf": "dolostone_rock_fragments", "dolostone rf": "dolostone_rock_fragments", "dolostone_rock_fragments_percent": "dolostone_rock_fragments",
		"chert": "chert", "chert_percent": "chert",
		"total_sedimentary_rf_percent": "total_sedimentary_rf_percent", "total sedimentary rf percent": "total_sedimentary_rf_percent", "total_sedimentary_rf": "total_sedimentary_rf_percent", "total sedimentary rf": "total_sedimentary_rf_percent", "total_sedimentary_rf %": "total_sedimentary_rf_percent",
		"total_rock_fragments_percent": "total_rock_fragments_percent", "total rock fragments percent": "total_rock_fragments_percent", "total_rock_fragments": "total_rock_fragments_percent", "total rock fragments": "total_rock_fragments_percent", "total_rock_fragments %": "total_rock_fragments_percent",
		
		// Clastic specific - Other Grains
		"rip_up_clast": "rip_up_clast", "rip up clast": "rip_up_clast", "rip_up": "rip_up_clast", "rip up": "rip_up_clast", "rip_up_clast_percent": "rip_up_clast", "rip up clast percent": "rip_up_clast",
		"glauconite": "glauconite", "glauconite_percent": "glauconite", "glauconite %": "glauconite",
		"foraminifera_grains": "foraminifera_grains", "foraminifera grains": "foraminifera_grains", "foram_grains": "foraminifera_grains", "foram grains": "foraminifera_grains", "foraminifera_grains_percent": "foraminifera_grains", "foraminifera grains percent": "foraminifera_grains",
		"undifferentiated_other_grains": "undifferentiated_other_grains", "undifferentiated other grains": "undifferentiated_other_grains", "undiff_other_grains": "undifferentiated_other_grains", "undiff other grains": "undifferentiated_other_grains", "undifferentiated_other_grains_percent": "undifferentiated_other_grains",
		"total_other_grains_percent": "total_other_grains_percent", "total other grains percent": "total_other_grains_percent", "total_other_grains": "total_other_grains_percent", "total other grains": "total_other_grains_percent", "total_other_grains %": "total_other_grains_percent",
		
		// Clastic specific - Matrix
		"clay_matrix": "clay_matrix", "clay matrix": "clay_matrix", "clay_matrix_percent": "clay_matrix", "clay matrix percent": "clay_matrix", "clay_matrix %": "clay_matrix",
		"mixed_clay_silt_fine_matrix": "mixed_clay_silt_fine_matrix", "mixed clay silt fine matrix": "mixed_clay_silt_fine_matrix", "mixed_clay_silt_matrix": "mixed_clay_silt_fine_matrix", "mixed clay silt matrix": "mixed_clay_silt_fine_matrix", "mixed_clay_silt_fine_matrix_percent": "mixed_clay_silt_fine_matrix",
		"silt_very_fine_matrix": "silt_very_fine_matrix", "silt very fine matrix": "silt_very_fine_matrix", "silt_fine_matrix": "silt_very_fine_matrix", "silt fine matrix": "silt_very_fine_matrix", "silt_very_fine_matrix_percent": "silt_very_fine_matrix",
		"organic_matrix": "organic_matrix", "organic_matrix_percent": "organic_matrix", "organic matrix percent": "organic_matrix", "organic_matrix %": "organic_matrix",
		"matrix_undifferentiated": "matrix_undifferentiated", "matrix undifferentiated": "matrix_undifferentiated", "undiff_matrix": "matrix_undifferentiated", "undiff matrix": "matrix_undifferentiated", "matrix_undifferentiated_percent": "matrix_undifferentiated",
		
		// Clastic specific - Authigenic Clay
		"kaolinite_replaces_k_feldspar": "kaolinite_replaces_k_feldspar", "kaolinite replaces k feldspar": "kaolinite_replaces_k_feldspar", "kaol_replaces_k_fsp": "kaolinite_replaces_k_feldspar", "kaol replaces k fsp": "kaolinite_replaces_k_feldspar", "kaolinite_replaces_k_feldspar_percent": "kaolinite_replaces_k_feldspar",
		"illite_pore_grain_lining": "illite_pore_grain_lining", "illite pore grain lining": "illite_pore_grain_lining", "illite_pore_lining": "illite_pore_grain_lining", "illite pore lining": "illite_pore_grain_lining", "illite_pore_grain_lining_percent": "illite_pore_grain_lining",
		"illite_pore_filling": "illite_pore_filling", "illite pore filling": "illite_pore_filling", "illite_pore_fill": "illite_pore_filling", "illite pore fill": "illite_pore_filling", "illite_pore_filling_percent": "illite_pore_filling",
		"illite_replaces_k_feldspar": "illite_replaces_k_feldspar", "illite replaces k feldspar": "illite_replaces_k_feldspar", "illite_replaces_k_fsp": "illite_replaces_k_feldspar", "illite replaces k fsp": "illite_replaces_k_feldspar", "illite_replaces_k_feldspar_percent": "illite_replaces_k_feldspar",
		"total_authigenic_clay_percent": "total_authigenic_clay_percent", "total authigenic clay percent": "total_authigenic_clay_percent", "total_authigenic_clay": "total_authigenic_clay_percent", "total authigenic clay": "total_authigenic_clay_percent", "total_authigenic_clay %": "total_authigenic_clay_percent",
		
		// Clastic specific - Authigenic Non-Clay
		"syntaxial_quartz_overgrowths": "syntaxial_quartz_overgrowths", "syntaxial quartz overgrowths": "syntaxial_quartz_overgrowths", "syntaxial_quartz": "syntaxial_quartz_overgrowths", "syntaxial quartz": "syntaxial_quartz_overgrowths", "syntaxial_quartz_overgrowths_percent": "syntaxial_quartz_overgrowths",
		"feldspar_overgrowths": "feldspar_overgrowths", "feldspar overgrowths": "feldspar_overgrowths", "feldspar_overgrowth": "feldspar_overgrowths", "feldspar overgrowth": "feldspar_overgrowths", "feldspar_overgrowths_percent": "feldspar_overgrowths",
		"fe_calcite": "fe_calcite", "fe calcite": "fe_calcite", "fe_calcite_percent": "fe_calcite", "fe calcite percent": "fe_calcite", "fe_calcite %": "fe_calcite",
		"fe_dolomite": "fe_dolomite", "fe dolomite": "fe_dolomite", "fe_dolomite_percent": "fe_dolomite", "fe dolomite percent": "fe_dolomite", "fe_dolomite %": "fe_dolomite",
		"siderite": "siderite", "siderite_percent": "siderite", "siderite %": "siderite",
		"mn_siderite": "mn_siderite", "mn siderite": "mn_siderite", "mn_siderite_percent": "mn_siderite", "mn siderite percent": "mn_siderite", "mn_siderite %": "mn_siderite",
		"iron_oxide_minerals": "iron_oxide_minerals", "iron oxide minerals": "iron_oxide_minerals", "iron_oxide": "iron_oxide_minerals", "iron oxide": "iron_oxide_minerals", "iron_oxide_minerals_percent": "iron_oxide_minerals",
		"total_authigenic_non_clay_percent": "total_authigenic_non_clay_percent", "total authigenic non clay percent": "total_authigenic_non_clay_percent", "total_authigenic_non_clay": "total_authigenic_non_clay_percent", "total authigenic non clay": "total_authigenic_non_clay_percent", "total_authigenic_non_clay %": "total_authigenic_non_clay_percent",
		
		// Clastic specific - Primary Porosity
		"intergranular": "intergranular", "inter_granular": "intergranular", "intergranular_percent": "intergranular", "intergranular %": "intergranular",
		"pri_porosity_intragranular": "pri_porosity_intragranular", "pri porosity intragranular": "pri_porosity_intragranular", "primary_porosity_intragranular": "pri_porosity_intragranular", "primary porosity intragranular": "pri_porosity_intragranular", "pri_porosity_intragranular_percent": "pri_porosity_intragranular",
		"total_primary_porosity_percent": "total_primary_porosity_percent", "total primary porosity percent": "total_primary_porosity_percent", "total_primary_porosity": "total_primary_porosity_percent", "total primary porosity": "total_primary_porosity_percent", "total_primary_porosity %": "total_primary_porosity_percent",
		
		// Clastic specific - Secondary Porosity
		"sec_porosity_intragranular": "sec_porosity_intragranular", "sec porosity intragranular": "sec_porosity_intragranular", "secondary_porosity_intragranular": "sec_porosity_intragranular", "secondary porosity intragranular": "sec_porosity_intragranular", "sec_porosity_intragranular_percent": "sec_porosity_intragranular",
		"intracrystalline": "intracrystalline", "intra_crystalline": "intracrystalline", "intracrystalline_percent": "intracrystalline", "intracrystalline %": "intracrystalline",
		"total_secondary_porosity_percent": "total_secondary_porosity_percent", "total secondary porosity percent": "total_secondary_porosity_percent", "total_secondary_porosity": "total_secondary_porosity_percent", "total secondary porosity": "total_secondary_porosity_percent", "total_secondary_porosity %": "total_secondary_porosity_percent",
		
		// Analysis Type (Common)
		"analysis_type": "analysis_types", "analysis type": "analysis_types",
	}
}