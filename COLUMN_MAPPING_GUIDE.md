# 🗺️ Complete Column Name Mapping Guide

## 📊 **Database Tables:**
- `petrography_carbonate` - For carbonate rock data
- `petrography_clastic` - For clastic/sandstone data

---

## 🌍 **GEOGRAPHIC FIELDS**

| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `country`, `Country`, `COUNTRY`, `nation`, `state` | `country` | Both |
| `region`, `Region`, `REGION`, `area`, `province` | `region` | Both |
| `sub_region`, `subregion`, `sub-region`, `sub area` | `sub_region` | Both |
| `business_regions`, `business region`, `business area` | `business_regions` | Both |
| `basin`, `Basin`, `BASIN`, `sedimentary basin` | `basin` | Both |
| `sub_basin`, `subbasin`, `sub-basin`, `sub basin` | `sub_basin` | Both |

---

## 🏗️ **WELL & FIELD INFORMATION**

| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `well`, `Well`, `WELL`, `well name`, `well_name` | `well_name_field_name` | Both |
| `field`, `Field`, `FIELD`, `field name`, `field_name` | `well_name_field_name` | Both |
| `well_name_field_name`, `well field name` | `well_name_field_name` | Both |
| `uwi`, `UWI`, `unique well identifier` | `uwi` | Both |

---

## 📍 **COORDINATES**

| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `lat`, `latitude`, `Latitude`, `LAT`, `lat.`, `north` | `latitude` | Both |
| `lon`, `longitude`, `Longitude`, `LON`, `long.`, `east` | `longitude` | Both |
| `coord`, `coordinate`, `coordinates`, `position` | `latitude` | Both |
| `x`, `y`, `easting`, `northing` | `longitude` | Both |

---

## 🏔️ **GEOLOGICAL INFORMATION**

| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `formation`, `Formation`, `FORMATION`, `form`, `unit` | `formation_name` | Both |
| `reservoir`, `Reservoir`, `RESERVOIR`, `reservoir name` | `reservoir_name` | Both |
| `period`, `Period`, `PERIOD`, `geological period` | `period` | Both |
| `epoch`, `Epoch`, `EPOCH`, `geological epoch` | `epoch` | Both |
| `age`, `Age`, `AGE`, `geological age` | `age` | Both |
| `onshore`, `offshore`, `onshore_offshore`, `location` | `onshore_offshore` | Both |

---

## 🌊 **WATER DEPTH**

| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `water_depth`, `water depth`, `water_depth_m`, `water depth m` | `water_depth_m` | Both |
| `water_depth_ft`, `water depth ft`, `water depth feet` | `water_depth_ft` | Both |
| `wd`, `WD`, `water`, `sea depth` | `water_depth_m` | Both |

---

## 📏 **DEPTH MEASUREMENTS**

### **Top Depth (Metric)**
| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `top_depth`, `top depth`, `top`, `depth_top`, `start depth` | `top_depth_mmddf` | Both |
| `top_depth_mmddf`, `top depth mmddf`, `top mmddf` | `top_depth_mmddf` | Both |
| `top_depth_mtvddf`, `top depth mtvddf`, `top mtvddf` | `top_depth_mtvddf` | Both |
| `top_depth_mtvdss`, `top depth mtvdss`, `top mtvdss` | `top_depth_mtvdss` | Both |
| `top_depth_mbml`, `top depth mbml`, `top mbml` | `top_depth_mbml` | Both |

### **Bottom Depth (Metric)**
| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `bottom_depth`, `bottom depth`, `bottom`, `depth_bottom`, `end depth` | `bottom_depth_mmddf` | Both |
| `bottom_depth_mmddf`, `bottom depth mmddf`, `bottom mmddf` | `bottom_depth_mmddf` | Both |
| `bottom_depth_mtvddf`, `bottom depth mtvddf`, `bottom mtvddf` | `bottom_depth_mtvddf` | Both |
| `bottom_depth_mtvdss`, `bottom depth mtvdss`, `bottom mtvdss` | `bottom_depth_mtvdss` | Both |
| `bottom_depth_mbml`, `bottom depth mbml`, `bottom mbml` | `bottom_depth_mbml` | Both |

### **Top Depth (Imperial)**
| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `top_depth_ft`, `top depth ft`, `top ft`, `top feet` | `top_depth_ftmddf` | Both |
| `top_depth_ftmddf`, `top depth ftmddf`, `top ftmddf` | `top_depth_ftmddf` | Both |
| `top_depth_fttvddf`, `top depth fttvddf`, `top fttvddf` | `top_depth_fttvddf` | Both |
| `top_depth_fttvdss`, `top depth fttvdss`, `top fttvdss` | `top_depth_fttvdss` | Both |
| `top_depth_ftbml`, `top depth ftbml`, `top ftbml` | `top_depth_ftbml` | Both |

### **Bottom Depth (Imperial)**
| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `bottom_depth_ft`, `bottom depth ft`, `bottom ft`, `bottom feet` | `bottom_depth_ftmddf` | Both |
| `bottom_depth_ftmddf`, `bottom depth ftmddf`, `bottom ftmddf` | `bottom_depth_ftmddf` | Both |
| `bottom_depth_fttvddf`, `bottom depth fttvddf`, `bottom fttvddf` | `bottom_depth_fttvddf` | Both |
| `bottom_depth_fttvdss`, `bottom depth fttvdss`, `bottom fttvdss` | `bottom_depth_fttvdss` | Both |
| `bottom_depth_ftbml`, `bottom depth ftbml`, `bottom ftbml` | `bottom_depth_ftbml` | Both |

---

## 🧪 **POROSITY & PERMEABILITY**

| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `porosity`, `Porosity`, `POROSITY`, `por`, `phi` | `visible_porosity_percent` | Both |
| `visible_porosity`, `visible porosity`, `vis porosity` | `visible_porosity_percent` | Both |
| `he_porosity`, `he porosity`, `helium porosity` | `he_porosity_percent` | Both |
| `permeability`, `Permeability`, `PERMEABILITY`, `perm`, `k` | `permeability_md` | Both |
| `permeability_md`, `permeability md`, `perm md` | `permeability_md` | Both |

---

## 🪨 **CARBONATE-SPECIFIC FIELDS**

| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `calcite`, `Calcite`, `CALCITE`, `cal` | `calcite` | Carbonate |
| `dolomite`, `Dolomite`, `DOLOMITE`, `dol` | `dolomite` | Carbonate |
| `micrite`, `Micrite`, `MICRITE`, `mic` | `micrite` | Carbonate |
| `bioclast`, `bioclasts`, `Bioclasts`, `BIOCLASTS` | `bioclasts` | Carbonate |
| `coral`, `Coral`, `CORAL`, `corals` | `coral` | Carbonate |
| `algae`, `Algae`, `ALGAE`, `red algae` | `red_algae` | Carbonate |
| `ooids`, `Ooids`, `OOIDS`, `ooid` | `peloids` | Carbonate |
| `pellets`, `Pellets`, `PELLETS`, `pellet` | `peloids` | Carbonate |

---

## 🏔️ **CLASTIC-SPECIFIC FIELDS**

| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `quartz`, `Quartz`, `QUARTZ`, `qtz` | `total_quartz_percent` | Clastic |
| `feldspar`, `Feldspar`, `FELDSPAR`, `fsp` | `total_feldspar_percent` | Clastic |
| `mica`, `Mica`, `MICA`, `muscovite` | `total_mica_percent` | Clastic |
| `grain_size`, `grain size`, `grain size`, `size` | `grain_size` | Clastic |
| `sorting`, `Sorting`, `SORTING`, `sort` | `sorting` | Clastic |
| `roundness`, `Roundness`, `ROUNDNESS`, `round` | `grain_shape` | Clastic |

---

## 🏗️ **FACIES & LITHOLOGY**

| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `lithofacies`, `lithofacies_core`, `lithology`, `facies` | `lithofacies_core` | Carbonate |
| `lithofacies`, `lithofacies_core`, `lithology`, `facies` | `lithofacies` | Clastic |
| `microfacies`, `microfacies_thin_section`, `thin section` | `microfacies_thin_section` | Carbonate |
| `depofacies`, `depofacies`, `depositional facies` | `depofacies` | Carbonate |

---

## 🔬 **TEXTURAL ANALYSIS (Clastic Only)**

| User Input Variations | Database Field | Table |
|----------------------|----------------|-------|
| `grain_shape`, `grain shape`, `shape`, `form` | `grain_shape` | Clastic |
| `grain_contact`, `grain contact`, `contact`, `contacts` | `grain_contact` | Clastic |
| `sedimentary_structure`, `sedimentary structure`, `structure` | `sedimentary_structure` | Clastic |

---

## 🎯 **USAGE EXAMPLES:**

### **Example 1: Carbonate Table**
```
User Headers: ["Well Name", "Depth", "Calcite", "Dolomite", "Porosity"]
Maps to: ["well_name_field_name", "top_depth_mmddf", "calcite", "dolomite", "visible_porosity_percent"]
Table: petrography_carbonate
```

### **Example 2: Clastic Table**
```
User Headers: ["Field", "Top Depth", "Quartz", "Feldspar", "Grain Size"]
Maps to: ["well_name_field_name", "top_depth_mmddf", "total_quartz_percent", "total_feldspar_percent", "grain_size"]
Table: petrography_clastic
```

### **Example 3: Mixed Terms**
```
User Headers: ["Country", "Basin", "Formation", "Latitude", "Longitude"]
Maps to: ["country", "basin", "formation_name", "latitude", "longitude"]
Table: Both (auto-detected based on other fields)
```

---

## 🤖 **Fuzzy Matching Rules:**

1. **Case Insensitive**: `Well` = `well` = `WELL`
2. **Space Variations**: `well name` = `well_name` = `wellname`
3. **Abbreviations**: `lat` = `latitude`, `depth` = `top_depth_mmddf`
4. **Synonyms**: `field` = `well_name_field_name`, `por` = `porosity`
5. **60% Similarity Threshold**: Close matches are accepted

---

## 📝 **Notes:**

- **Both Tables**: Fields that exist in both carbonate and clastic tables
- **Carbonate Only**: Fields specific to carbonate rock analysis
- **Clastic Only**: Fields specific to clastic/sandstone analysis
- **Auto-Detection**: System automatically chooses table based on field content
- **Flexible Matching**: Multiple variations of the same field are supported

