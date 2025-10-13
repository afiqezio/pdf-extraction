# 🗺️ Complete Column Name Mapping - Based on Database Schema

## 📊 **Database Tables:**
- `petrography_carbonate` - For carbonate rock data
- `petrography_clastic` - For clastic/sandstone data

---

## 🌍 **COMMON FIELDS (Both Tables)**

### **Geographic Location Fields**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `country` | country, Country, COUNTRY, nation, state, country_name |
| `region` | region, Region, REGION, area, province, regional |
| `sub_region` | sub_region, subregion, sub-region, sub area, subarea |
| `business_regions` | business_regions, business region, business area, business_regions |
| `basin` | basin, Basin, BASIN, sedimentary basin, basin_name |
| `sub_basin` | sub_basin, subbasin, sub-basin, sub basin, subbasin_name |

### **Well & Field Information**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `well_name_field_name` | well, Well, WELL, well_name, well name, field, Field, FIELD, field_name, field name, well_name_field_name, well field name, well_field_name, wellname, fieldname |
| `uwi` | uwi, UWI, unique well identifier, unique_well_identifier, well_id, wellid |

### **Coordinates**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `latitude` | lat, latitude, Latitude, LAT, lat., north, north_lat, n_lat, y_coord, y_coordinate |
| `longitude` | lon, longitude, Longitude, LON, long., east, east_lon, e_lon, x_coord, x_coordinate |

### **Geological Information**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `formation_name` | formation, Formation, FORMATION, form, unit, formation_name, formation name, geologic_formation |
| `reservoir_name` | reservoir, Reservoir, RESERVOIR, reservoir_name, reservoir name, reservoir_field |
| `period` | period, Period, PERIOD, geological period, geo_period, time_period |
| `epoch` | epoch, Epoch, EPOCH, geological epoch, geo_epoch, time_epoch |
| `age` | age, Age, AGE, geological age, geo_age, time_age, rock_age |
| `onshore_offshore` | onshore, offshore, onshore_offshore, location, location_type, on_off, onshore_offshore |

### **Water Depth**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `water_depth_m` | water_depth, water depth, water_depth_m, water depth m, water_depth_meters, wd, WD, water, sea depth, water_depth_m |
| `water_depth_ft` | water_depth_ft, water depth ft, water depth feet, water_depth_feet, wd_ft, water_ft |

### **Depth Reference Fields**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `depth_reference_type` | depth_reference_type, depth reference type, depth_ref_type, ref_type, reference_type |
| `depth_reference_elevation_m` | depth_reference_elevation_m, depth reference elevation m, ref_elevation_m, elevation_m |
| `depth_reference_elevation_ft` | depth_reference_elevation_ft, depth reference elevation ft, ref_elevation_ft, elevation_ft |
| `ground_level_elevation_m` | ground_level_elevation_m, ground level elevation m, ground_elevation_m, gnd_elevation_m |
| `ground_level_elevation_ft` | ground_level_elevation_ft, ground level elevation ft, ground_elevation_ft, gnd_elevation_ft |

### **Top Depth (Metric)**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `top_depth_mmddf` | top_depth, top depth, top, depth_top, start depth, top_depth_mmddf, top depth mmddf, top mmddf, top_md, top_mddf |
| `top_depth_mtvddf` | top_depth_mtvddf, top depth mtvddf, top mtvddf, top_tvddf, top_tvd |
| `top_depth_mtvdss` | top_depth_mtvdss, top depth mtvdss, top mtvdss, top_tvdss, top_tvdss |
| `top_depth_mbml` | top_depth_mbml, top depth mbml, top mbml, top_bml, top_bml |

### **Bottom Depth (Metric)**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `bottom_depth_mmddf` | bottom_depth, bottom depth, bottom, depth_bottom, end depth, bottom_depth_mmddf, bottom depth mmddf, bottom mmddf, bottom_md, bottom_mddf |
| `bottom_depth_mtvddf` | bottom_depth_mtvddf, bottom depth mtvddf, bottom mtvddf, bottom_tvddf, bottom_tvd |
| `bottom_depth_mtvdss` | bottom_depth_mtvdss, bottom depth mtvdss, bottom mtvdss, bottom_tvdss, bottom_tvdss |
| `bottom_depth_mbml` | bottom_depth_mbml, bottom depth mbml, bottom mbml, bottom_bml, bottom_bml |

### **Top Depth (Imperial)**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `top_depth_ftmddf` | top_depth_ft, top depth ft, top ft, top feet, top_depth_ftmddf, top depth ftmddf, top ftmddf, top_ft, top_ftmddf |
| `top_depth_fttvddf` | top_depth_fttvddf, top depth fttvddf, top fttvddf, top_fttvddf, top_fttvd |
| `top_depth_fttvdss` | top_depth_fttvdss, top depth fttvdss, top fttvdss, top_fttvdss, top_fttvdss |
| `top_depth_ftbml` | top_depth_ftbml, top depth ftbml, top ftbml, top_ftbml, top_ftbml |

### **Bottom Depth (Imperial)**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `bottom_depth_ftmddf` | bottom_depth_ft, bottom depth ft, bottom ft, bottom feet, bottom_depth_ftmddf, bottom depth ftmddf, bottom ftmddf, bottom_ft, bottom_ftmddf |
| `bottom_depth_fttvddf` | bottom_depth_fttvddf, bottom depth fttvddf, bottom fttvddf, bottom_fttvddf, bottom_fttvd |
| `bottom_depth_fttvdss` | bottom_depth_fttvdss, bottom depth fttvdss, bottom fttvdss, bottom_fttvdss, bottom_fttvdss |
| `bottom_depth_ftbml` | bottom_depth_ftbml, bottom depth ftbml, bottom ftbml, bottom_ftbml, bottom_ftbml |

### **Porosity & Permeability (Common)**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `visible_porosity_percent` | porosity, Porosity, POROSITY, por, phi, visible_porosity, visible porosity, vis porosity, porosity_percent, porosity % |
| `permeability_md` | permeability, Permeability, PERMEABILITY, perm, k, permeability_md, permeability md, perm md, perm_md, k_md |

---

## 🪨 **CARBONATE-SPECIFIC FIELDS**

### **Facies Description**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `lithofacies_core` | lithofacies, lithofacies_core, lithology, facies, core_facies, lithofacies_core, lithofacies core, core_lithofacies |
| `microfacies_thin_section` | microfacies, microfacies_thin_section, thin_section, thin section, microfacies_thin_section, microfacies thin section |
| `depofacies` | depofacies, depofacies, depositional_facies, depositional facies, depo_facies |

### **Porosity (Carbonate)**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `he_porosity_percent` | he_porosity, he porosity, helium_porosity, helium porosity, he_porosity_percent, he porosity percent, he_porosity % |

### **Matrix Mineralogy**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `calcite` | calcite, Calcite, CALCITE, cal, calcite_percent, calcite %, calcite_percent |
| `dolomite` | dolomite, Dolomite, DOLOMITE, dol, dolomite_percent, dolomite %, dolomite_percent |
| `micrite` | micrite, Micrite, MICRITE, mic, micrite_percent, micrite %, micrite_percent |
| `micrite_envelopes` | micrite_envelopes, micrite envelopes, micrite_envelope, micrite envelope, micrite_envelopes_percent |
| `microspar_pseudospar` | microspar_pseudospar, microspar pseudospar, microspar_pseudospar_percent, microspar pseudospar percent |
| `kaolinite` | kaolinite, Kaolinite, KAOLINITE, kaol, kaolinite_percent, kaolinite %, kaolinite_percent |
| `clay` | clay, Clay, CLAY, clay_percent, clay %, clay_percent, total_clay |
| `total_mineralogy_matrix_percent` | total_mineralogy_matrix_percent, total mineralogy matrix percent, total_matrix_percent, total matrix percent |

### **Bioclasts - General and Specific Foraminifera**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `bioclasts` | bioclast, bioclasts, Bioclasts, BIOCLASTS, bioclast_percent, bioclasts_percent, bioclast %, bioclasts % |
| `lepido` | lepido, Lepido, LEPIDO, lepidocyclina, lepido_percent, lepido %, lepido_percent |
| `coral` | coral, Coral, CORAL, corals, coral_percent, coral %, coral_percent |
| `rhodolith` | rhodolith, Rhodolith, RHODOLITH, rhodolith_percent, rhodolith %, rhodolith_percent |
| `red_algae` | red_algae, red algae, red_algae_percent, red algae percent, red_algae %, red algae % |
| `red_algae_enc` | red_algae_enc, red algae enc, red_algae_enc_percent, red algae enc percent, red_algae_enc % |
| `green_algae` | green_algae, green algae, green_algae_percent, green algae percent, green_algae %, green algae % |
| `echinoderms` | echinoderms, Echinoderms, ECHINODERMS, echinoderm, echinoderms_percent, echinoderms %, echinoderms_percent |
| `miliolid` | miliolid, Miliolid, MILIOLID, miliolid_percent, miliolid %, miliolid_percent |
| `lepidocyclina` | lepidocyclina, Lepidocyclina, LEPIDOCYCLINA, lepidocyclina_percent, lepidocyclina %, lepidocyclina_percent |
| `cycloclypeus` | cycloclypeus, Cycloclypeus, CYCLOCLYPEUS, cycloclypeus_percent, cycloclypeus %, cycloclypeus_percent |
| `operculina` | operculina, Operculina, OPERCULINA, operculina_percent, operculina %, operculina_percent |
| `other_rotaliids` | other_rotaliids, other rotaliids, other_rotaliids_percent, other rotaliids percent, other_rotaliids % |
| `gypsinid` | gypsinid, Gypsinid, GYPSINID, gypsinid_percent, gypsinid %, gypsinid_percent |
| `planorbulinella` | planorbulinella, Planorbulinella, PLANORBULINELLA, planorbulinella_percent, planorbulinella %, planorbulinella_percent |
| `hemotremid` | hemotremid, Hemotremid, HEMOTREMID, hemotremid_percent, hemotremid %, hemotremid_percent |
| `heterostegina` | heterostegina, Heterostegina, HETEROSTEGINA, heterostegina_percent, heterostegina %, heterostegina_percent |
| `enc_frm` | enc_frm, enc frm, enc_frm_percent, enc frm percent, enc_frm % |
| `planktonic` | planktonic, Planktonic, PLANKTONIC, planktonic_percent, planktonic %, planktonic_percent |
| `bryozoans` | bryozoans, Bryozoans, BRYOZOANS, bryozoan, bryozoans_percent, bryozoans %, bryozoans_percent |
| `amphistegina` | amphistegina, Amphistegina, AMPHISTEGINA, amphistegina_percent, amphistegina %, amphistegina_percent |
| `gastropods` | gastropods, Gastropods, GASTROPODS, gastropod, gastropods_percent, gastropods %, gastropods_percent |
| `bivalve` | bivalve, Bivalve, BIVALVE, bivalves, bivalve_percent, bivalve %, bivalve_percent |
| `ostracod` | ostracod, Ostracod, OSTRACOD, ostracods, ostracod_percent, ostracod %, ostracod_percent |
| `oncoids` | oncoids, Oncoids, ONCOIDS, oncoid, oncoids_percent, oncoids %, oncoids_percent |
| `undiff_molluscs` | undiff_molluscs, undiff molluscs, undifferentiated_molluscs, undiff_molluscs_percent, undiff molluscs percent |
| `undiff_benthonic` | undiff_benthonic, undiff benthonic, undifferentiated_benthonic, undiff_benthonic_percent, undiff benthonic percent |
| `undiff_skeletal` | undiff_skeletal, undiff skeletal, undifferentiated_skeletal, undiff_skeletal_percent, undiff skeletal percent |
| `undiff_foram` | undiff_foram, undiff foram, undifferentiated_foram, undiff_foram_percent, undiff foram percent |
| `total_skeletal_percent` | total_skeletal_percent, total skeletal percent, total_skeletal, total skeletal, total_skeletal % |

### **Non-Skeletal Components**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `organic` | organic, Organic, ORGANIC, organic_percent, organic %, organic_percent, organic_matter |
| `peloids` | peloids, Peloids, PELLETS, pellets, pellet, peloids_percent, peloids %, peloids_percent |
| `micritised_grains` | micritised_grains, micritised grains, micritized_grains, micritised_grains_percent, micritised grains percent |
| `pseudoclasts` | pseudoclasts, Pseudoclasts, PSEUDOCLASTS, pseudoclast, pseudoclasts_percent, pseudoclasts %, pseudoclasts_percent |
| `intraclast` | intraclast, Intraclast, INTRACLAST, intraclasts, intraclast_percent, intraclast %, intraclast_percent |
| `quartz` | quartz, Quartz, QUARTZ, qtz, quartz_percent, quartz %, quartz_percent |
| `total_non_skeletal_percent` | total_non_skeletal_percent, total non skeletal percent, total_non_skeletal, total non skeletal, total_non_skeletal % |

### **Porosity Types**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `interparticle` | interparticle, Interparticle, INTERPARTICLE, inter_particle, interparticle_percent, interparticle %, interparticle_percent |
| `intraparticle` | intraparticle, Intraparticle, INTRAPARTICLE, intra_particle, intraparticle_percent, intraparticle %, intraparticle_percent |
| `intercrystalline` | intercrystalline, Intercrystalline, INTERCRYSTALLINE, inter_crystalline, intercrystalline_percent, intercrystalline %, intercrystalline_percent |
| `matrix_intercrystalline` | matrix_intercrystalline, matrix intercrystalline, matrix_intercrystalline_percent, matrix intercrystalline percent |
| `mouldic` | mouldic, Mouldic, MOULDIC, moldic, mouldic_percent, mouldic %, mouldic_percent |
| `vuggy` | vuggy, Vuggy, VUGGY, vug, vuggy_percent, vuggy %, vuggy_percent |
| `fractures` | fractures, Fractures, FRACTURES, fracture, fractures_percent, fractures %, fractures_percent |
| `micro` | micro, Micro, MICRO, micro_percent, micro %, micro_percent, micro_porosity |
| `total_porosity_percent` | total_porosity_percent, total porosity percent, total_porosity, total porosity, total_porosity % |

### **Cement Types**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `fringing` | fringing, Fringing, FRINGING, fringing_percent, fringing %, fringing_percent |
| `meniscus` | meniscus, Meniscus, MENISCUS, meniscus_percent, meniscus %, meniscus_percent |
| `blocky` | blocky, Blocky, BLOCKY, blocky_percent, blocky %, blocky_percent |
| `sparry` | sparry, Sparry, SPARRY, sparry_percent, sparry %, sparry_percent |
| `micritic` | micritic, Micritic, MICRITIC, micritic_percent, micritic %, micritic_percent |
| `pendant` | pendant, Pendant, PENDANT, pendant_percent, pendant %, pendant_percent |
| `syntax` | syntax, Syntax, SYNTAX, syntax_percent, syntax %, syntax_percent |
| `calcite_syntaxial` | calcite_syntaxial, calcite syntaxial, calcite_syntaxial_percent, calcite syntaxial percent |
| `calcite_fringing` | calcite_fringing, calcite fringing, calcite_fringing_percent, calcite fringing percent |
| `calcite_mosaic` | calcite_mosaic, calcite mosaic, calcite_mosaic_percent, calcite mosaic percent |
| `calcite_blocky` | calcite_blocky, calcite blocky, calcite_blocky_percent, calcite blocky percent |
| `calcite_ferroan` | calcite_ferroan, calcite ferroan, calcite_ferroan_percent, calcite ferroan percent |
| `pyrite` | pyrite, Pyrite, PYRITE, pyrite_percent, pyrite %, pyrite_percent |
| `fluorite` | fluorite, Fluorite, FLUORITE, fluorite_percent, fluorite %, fluorite_percent |
| `total_cement_percent` | total_cement_percent, total cement percent, total_cement, total cement, total_cement % |

### **Replacement and Accessories**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `replacement` | replacement, Replacement, REPLACEMENT, replacement_percent, replacement %, replacement_percent |
| `saddle` | saddle, Saddle, SADDLE, saddle_percent, saddle %, saddle_percent |
| `total_dolomite_percent` | total_dolomite_percent, total dolomite percent, total_dolomite, total dolomite, total_dolomite % |
| `stylolite` | stylolite, Stylolite, STYLOLITE, stylolite_percent, stylolite %, stylolite_percent |
| `bioturbation` | bioturbation, Bioturbation, BIOTURBATION, bioturbation_percent, bioturbation %, bioturbation_percent |
| `total_accessories_percent` | total_accessories_percent, total accessories percent, total_accessories, total accessories, total_accessories % |
| `total_percent` | total_percent, total percent, total, total %, total_percent |

### **Analysis Type (Carbonate)**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `analysis_types` | analysis_types, analysis types, analysis_type, analysis type, analysis_types, analysis types |

---

## 🏔️ **CLASTIC-SPECIFIC FIELDS**

### **Textural Analysis**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `grain_size` | grain_size, grain size, size, grain_size, grain size, grain_size, grain size, grain_size |
| `grain_shape` | grain_shape, grain shape, shape, form, grain_shape, grain shape, grain_shape, grain shape |
| `grain_contact` | grain_contact, grain contact, contact, contacts, grain_contact, grain contact, grain_contact, grain contact |
| `sedimentary_structure` | sedimentary_structure, sedimentary structure, structure, sedimentary_structure, sedimentary structure, sedimentary_structure, sedimentary structure |
| `sorting` | sorting, Sorting, SORTING, sort, sorting, sorting, sorting, sorting |
| `lithofacies` | lithofacies, lithofacies, lithology, facies, lithofacies, lithofacies, lithofacies, lithofacies |

### **Porosity (Clastic)**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `ambient_he_porosity_percent` | ambient_he_porosity, ambient he porosity, ambient_he_porosity_percent, ambient he porosity percent, ambient_he_porosity % |
| `grain_density_g_cc` | grain_density, grain density, grain_density_g_cc, grain density g/cc, grain_density_g_cc, grain density g/cc |

### **Quartz Content**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `monocrystalline_quartz` | monocrystalline_quartz, monocrystalline quartz, mono_quartz, mono quartz, monocrystalline_quartz_percent, monocrystalline quartz percent |
| `polycrystalline_quartz` | polycrystalline_quartz, polycrystalline quartz, poly_quartz, poly quartz, polycrystalline_quartz_percent, polycrystalline quartz percent |
| `total_quartz_percent` | total_quartz_percent, total quartz percent, total_quartz, total quartz, total_quartz %, total quartz % |

### **Feldspar Content**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `potassium_feldspar` | potassium_feldspar, potassium feldspar, k_feldspar, k feldspar, k_fsp, k fsp, potassium_feldspar_percent, potassium feldspar percent |
| `plagioclase` | plagioclase, Plagioclase, PLAGIOCLASE, plagioclase_percent, plagioclase %, plagioclase_percent |
| `feldspar_undifferentiated` | feldspar_undifferentiated, feldspar undifferentiated, undiff_feldspar, undiff feldspar, feldspar_undifferentiated_percent, feldspar undifferentiated percent |
| `total_feldspar_percent` | total_feldspar_percent, total feldspar percent, total_feldspar, total feldspar, total_feldspar %, total feldspar % |

### **Mica Content**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `muscovite` | muscovite, Muscovite, MUSCOVITE, muscovite_percent, muscovite %, muscovite_percent |
| `biotite` | biotite, Biotite, BIOTITE, biotite_percent, biotite %, biotite_percent |
| `mica_undifferentiated` | mica_undifferentiated, mica undifferentiated, undiff_mica, undiff mica, mica_undifferentiated_percent, mica undifferentiated percent |
| `total_mica_percent` | total_mica_percent, total mica percent, total_mica, total mica, total_mica %, total mica % |

### **Heavy Minerals**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `zircon` | zircon, Zircon, ZIRCON, zircon_percent, zircon %, zircon_percent |
| `tourmaline` | tourmaline, Tourmaline, TOURMALINE, tourmaline_percent, tourmaline %, tourmaline_percent |
| `heavy_minerals_undifferentiated` | heavy_minerals_undifferentiated, heavy minerals undifferentiated, undiff_heavy_minerals, undiff heavy minerals, heavy_minerals_undifferentiated_percent |
| `total_heavy_minerals_percent` | total_heavy_minerals_percent, total heavy minerals percent, total_heavy_minerals, total heavy minerals, total_heavy_minerals % |

### **Igneous Rock Fragments**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `plutonic_rock_fragments` | plutonic_rock_fragments, plutonic rock fragments, plutonic_rf, plutonic rf, plutonic_rock_fragments_percent, plutonic rock fragments percent |
| `mafic_intermediate_volcanic_fragment` | mafic_intermediate_volcanic_fragment, mafic intermediate volcanic fragment, mafic_volcanic, mafic volcanic, mafic_intermediate_volcanic_fragment_percent |
| `volcanic_rock_fragment` | volcanic_rock_fragment, volcanic rock fragment, volcanic_rf, volcanic rf, volcanic_rock_fragment_percent, volcanic rock fragment percent |
| `total_igneous_rf_percent` | total_igneous_rf_percent, total igneous rf percent, total_igneous_rf, total igneous rf, total_igneous_rf % |

### **Metamorphic Rock Fragments**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `quartzose_rock_fragment` | quartzose_rock_fragment, quartzose rock fragment, quartzose_rf, quartzose rf, quartzose_rock_fragment_percent, quartzose rock fragment percent |
| `schistose_rock_fragment` | schistose_rock_fragment, schistose rock fragment, schistose_rf, schistose rf, schistose_rock_fragment_percent, schistose rock fragment percent |
| `metamorphic_rock_fragment_undifferentiated` | metamorphic_rock_fragment_undifferentiated, metamorphic rock fragment undifferentiated, meta_rf_undiff, meta rf undiff, metamorphic_rock_fragment_undifferentiated_percent |
| `total_metamorphic_rf_percent` | total_metamorphic_rf_percent, total metamorphic rf percent, total_metamorphic_rf, total metamorphic rf, total_metamorphic_rf % |

### **Sedimentary Rock Fragments**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `sandstone_siltstone_rock_fragments` | sandstone_siltstone_rock_fragments, sandstone siltstone rock fragments, ss_silt_rf, ss silt rf, sandstone_siltstone_rock_fragments_percent |
| `argillaceous_rock_fragments` | argillaceous_rock_fragments, argillaceous rock fragments, argillaceous_rf, argillaceous rf, argillaceous_rock_fragments_percent |
| `siliciclastic_rock_fragments_undifferentiated` | siliciclastic_rock_fragments_undifferentiated, siliciclastic rock fragments undifferentiated, siliciclastic_rf_undiff, siliciclastic rf undiff, siliciclastic_rock_fragments_undifferentiated_percent |
| `limestone_rock_fragments` | limestone_rock_fragments, limestone rock fragments, limestone_rf, limestone rf, limestone_rock_fragments_percent |
| `dolostone_rock_fragments` | dolostone_rock_fragments, dolostone rock fragments, dolostone_rf, dolostone rf, dolostone_rock_fragments_percent |
| `chert` | chert, Chert, CHERT, chert_percent, chert %, chert_percent |
| `total_sedimentary_rf_percent` | total_sedimentary_rf_percent, total sedimentary rf percent, total_sedimentary_rf, total sedimentary rf, total_sedimentary_rf % |
| `total_rock_fragments_percent` | total_rock_fragments_percent, total rock fragments percent, total_rock_fragments, total rock fragments, total_rock_fragments % |

### **Other Grains**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `rip_up_clast` | rip_up_clast, rip up clast, rip_up, rip up, rip_up_clast_percent, rip up clast percent |
| `glauconite` | glauconite, Glauconite, GLAUCONITE, glauconite_percent, glauconite %, glauconite_percent |
| `bioclast` | bioclast, Bioclast, BIOCLAST, bioclast_percent, bioclast %, bioclast_percent |
| `foraminifera_grains` | foraminifera_grains, foraminifera grains, foram_grains, foram grains, foraminifera_grains_percent, foraminifera grains percent |
| `undifferentiated_other_grains` | undifferentiated_other_grains, undifferentiated other grains, undiff_other_grains, undiff other grains, undifferentiated_other_grains_percent |
| `total_other_grains_percent` | total_other_grains_percent, total other grains percent, total_other_grains, total other grains, total_other_grains % |

### **Matrix**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `clay_matrix` | clay_matrix, clay matrix, clay_matrix_percent, clay matrix percent, clay_matrix % |
| `mixed_clay_silt_fine_matrix` | mixed_clay_silt_fine_matrix, mixed clay silt fine matrix, mixed_clay_silt_matrix, mixed clay silt matrix, mixed_clay_silt_fine_matrix_percent |
| `silt_very_fine_matrix` | silt_very_fine_matrix, silt very fine matrix, silt_fine_matrix, silt fine matrix, silt_very_fine_matrix_percent |
| `organic_matrix` | organic_matrix, organic matrix, organic_matrix_percent, organic matrix percent, organic_matrix % |
| `matrix_undifferentiated` | matrix_undifferentiated, matrix undifferentiated, undiff_matrix, undiff matrix, matrix_undifferentiated_percent |
| `total_matrix_percent` | total_matrix_percent, total matrix percent, total_matrix, total matrix, total_matrix % |

### **Authigenic Clay**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `kaolinite` | kaolinite, Kaolinite, KAOLINITE, kaol, kaolinite_percent, kaolinite %, kaolinite_percent |
| `kaolinite_replaces_k_feldspar` | kaolinite_replaces_k_feldspar, kaolinite replaces k feldspar, kaol_replaces_k_fsp, kaol replaces k fsp, kaolinite_replaces_k_feldspar_percent |
| `illite_pore_grain_lining` | illite_pore_grain_lining, illite pore grain lining, illite_pore_lining, illite pore lining, illite_pore_grain_lining_percent |
| `illite_pore_filling` | illite_pore_filling, illite pore filling, illite_pore_fill, illite pore fill, illite_pore_filling_percent |
| `illite_replaces_k_feldspar` | illite_replaces_k_feldspar, illite replaces k feldspar, illite_replaces_k_fsp, illite replaces k fsp, illite_replaces_k_feldspar_percent |
| `total_authigenic_clay_percent` | total_authigenic_clay_percent, total authigenic clay percent, total_authigenic_clay, total authigenic clay, total_authigenic_clay % |

### **Authigenic Non-Clay**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `syntaxial_quartz_overgrowths` | syntaxial_quartz_overgrowths, syntaxial quartz overgrowths, syntaxial_quartz, syntaxial quartz, syntaxial_quartz_overgrowths_percent |
| `feldspar_overgrowths` | feldspar_overgrowths, feldspar overgrowths, feldspar_overgrowth, feldspar overgrowth, feldspar_overgrowths_percent |
| `fe_calcite` | fe_calcite, fe calcite, fe_calcite_percent, fe calcite percent, fe_calcite % |
| `fe_dolomite` | fe_dolomite, fe dolomite, fe_dolomite_percent, fe dolomite percent, fe_dolomite % |
| `siderite` | siderite, Siderite, SIDERITE, siderite_percent, siderite %, siderite_percent |
| `mn_siderite` | mn_siderite, mn siderite, mn_siderite_percent, mn siderite percent, mn_siderite % |
| `pyrite` | pyrite, Pyrite, PYRITE, pyrite_percent, pyrite %, pyrite_percent |
| `iron_oxide_minerals` | iron_oxide_minerals, iron oxide minerals, iron_oxide, iron oxide, iron_oxide_minerals_percent |
| `total_authigenic_non_clay_percent` | total_authigenic_non_clay_percent, total authigenic non clay percent, total_authigenic_non_clay, total authigenic non clay, total_authigenic_non_clay % |

### **Primary Porosity**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `intergranular` | intergranular, Intergranular, INTERGRANULAR, inter_granular, intergranular_percent, intergranular %, intergranular_percent |
| `intercrystalline` | intercrystalline, Intercrystalline, INTERCRYSTALLINE, inter_crystalline, intercrystalline_percent, intercrystalline %, intercrystalline_percent |
| `pri_porosity_intragranular` | pri_porosity_intragranular, pri porosity intragranular, primary_porosity_intragranular, primary porosity intragranular, pri_porosity_intragranular_percent |
| `total_primary_porosity_percent` | total_primary_porosity_percent, total primary porosity percent, total_primary_porosity, total primary porosity, total_primary_porosity % |

### **Secondary Porosity**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `sec_porosity_intragranular` | sec_porosity_intragranular, sec porosity intragranular, secondary_porosity_intragranular, secondary porosity intragranular, sec_porosity_intragranular_percent |
| `intracrystalline` | intracrystalline, Intracrystalline, INTRACRYSTALLINE, intra_crystalline, intracrystalline_percent, intracrystalline %, intracrystalline_percent |
| `mouldic` | mouldic, Mouldic, MOULDIC, moldic, mouldic_percent, mouldic %, mouldic_percent |
| `fracture` | fracture, Fracture, FRACTURE, fractures, fracture_percent, fracture %, fracture_percent |
| `total_secondary_porosity_percent` | total_secondary_porosity_percent, total secondary porosity percent, total_secondary_porosity, total secondary porosity, total_secondary_porosity % |

### **Total and Analysis (Clastic)**
| Database Field | Possible User Input Variations |
|----------------|-------------------------------|
| `total_percent` | total_percent, total percent, total, total %, total_percent |
| `analysis_types` | analysis_types, analysis types, analysis_type, analysis type, analysis_types, analysis types |

---

## 🎯 **USAGE EXAMPLES:**

### **Example 1: Carbonate Table**
```
User Headers: ["Well Name", "Depth", "Calcite", "Dolomite", "Porosity", "Bioclasts"]
Maps to: ["well_name_field_name", "top_depth_mmddf", "calcite", "dolomite", "visible_porosity_percent", "bioclasts"]
Table: petrography_carbonate
```

### **Example 2: Clastic Table**
```
User Headers: ["Field", "Top Depth", "Quartz", "Feldspar", "Grain Size", "Sorting"]
Maps to: ["well_name_field_name", "top_depth_mmddf", "total_quartz_percent", "total_feldspar_percent", "grain_size", "sorting"]
Table: petrography_clastic
```

### **Example 3: Mixed Terms**
```
User Headers: ["Country", "Basin", "Formation", "Latitude", "Longitude", "Depth"]
Maps to: ["country", "basin", "formation_name", "latitude", "longitude", "top_depth_mmddf"]
Table: Both (auto-detected based on other fields)
```

---

## 🤖 **Fuzzy Matching Rules:**

1. **Case Insensitive**: `Well` = `well` = `WELL`
2. **Space Variations**: `well name` = `well_name` = `wellname`
3. **Abbreviations**: `lat` = `latitude`, `depth` = `top_depth_mmddf`
4. **Synonyms**: `field` = `well_name_field_name`, `por` = `porosity`
5. **60% Similarity Threshold**: Close matches are accepted
6. **Underscore Variations**: `grain_size` = `grain size` = `grain-size`
7. **Percent Variations**: `calcite` = `calcite_percent` = `calcite %` = `calcite%`

---

## 📝 **Notes:**

- **Both Tables**: Fields that exist in both carbonate and clastic tables
- **Carbonate Only**: Fields specific to carbonate rock analysis
- **Clastic Only**: Fields specific to clastic/sandstone analysis
- **Auto-Detection**: System automatically chooses table based on field content
- **Flexible Matching**: Multiple variations of the same field are supported
- **Total Fields**: Many fields have corresponding "total" fields for percentages
- **Undifferentiated Fields**: Many fields have "undifferentiated" variants for mixed categories
