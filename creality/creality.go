package creality

import (
	"encoding/json"
	"fmt"
	"io/ioutil" // Still needed for LoadDefaultDatabase(filePath string) if it remains, otherwise remove.
	"strconv"
	"strings"

	"filament-sync-tool/profiles" // Import the profiles package from the local module
)

// MaterialDatabase represents the top-level structure of the Creality material_database.json.
type MaterialDatabase struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	ReqID  string `json:"reqId"`
	Result struct {
		List    []FilamentProfileEntry `json:"list"`
		Count   int                    `json:"count"`
		Version string                 `json:"version"`
	} `json:"result"`
}

// FilamentProfileEntry represents a single filament profile within the Creality database list.
type FilamentProfileEntry struct {
	EngineVersion  string   `json:"engineVersion"`
	PrinterIntName string   `json:"printerIntName"`
	NozzleDiameter []string `json:"nozzleDiameter"`
	KVParam        KVParam  `json:"kvParam"` // Key-Value Parameters specific to the printer
	Base           BaseInfo `json:"base"`    // Base information for the filament
}

// KVParam holds the detailed technical parameters for a filament.
type KVParam struct {
	ActivateAirFiltration            string `json:"activate_air_filtration"`
	ActivateChamberTempControl       string `json:"activate_chamber_temp_control"`
	AdditionalCoolingFanSpeed        string `json:"additional_cooling_fan_speed"`
	ChamberTemperature               string `json:"chamber_temperature"`
	CloseFanTheFirstXLayers          string `json:"close_fan_the_first_x_layers"`
	CompatiblePrinters               string `json:"compatible_printers"`
	CompatiblePrintersCondition      string `json:"compatible_printers_condition"`
	CompatiblePrints                 string `json:"compatible_prints"`
	CompatiblePrintsCondition        string `json:"compatible_prints_condition"`
	CompletePrintExhaustFanSpeed     string `json:"complete_print_exhaust_fan_speed"`
	CoolCdsFanStartAtHeight          string `json:"cool_cds_fan_start_at_height"`
	CoolPlateTemp                    string `json:"cool_plate_temp"`
	CoolPlateTempInitialLayer        string `json:"cool_plate_temp_initial_layer"`
	CoolSpecialCdsFanSpeed           string `json:"cool_special_cds_fan_speed"`
	DefaultFilamentColour            string `json:"default_filament_colour"`
	DuringPrintExhaustFanSpeed       string `json:"during_print_exhaust_fan_speed"`
	EnableOverhangBridgeFan          string `json:"enable_overhang_bridge_fan"`
	EnablePressureAdvance            string `json:"enable_pressure_advance"`
	EnableSpecialAreaAdditionalCoolingFan string `json:"enable_special_area_additional_cooling_fan"`
	EngPlateTemp                     string `json:"eng_plate_temp"`
	EngPlateTempInitialLayer         string `json:"eng_plate_temp_initial_layer"`
	EpoxyResinPlateTemp              string `json:"epoxy_resin_plate_temp"`
	EpoxyResinPlateTempInitialLayer  string `json:"epoxy_resin_plate_initial_layer"`
	FanCoolingLayerTime              string `json:"fan_cooling_layer_time"`
	FanMaxSpeed                      string `json:"fan_max_speed"`
	FanMinSpeed                      string `json:"fan_min_speed"`
	FilamentCoolingFinalSpeed        string `json:"filament_cooling_final_speed"`
	FilamentCoolingInitialSpeed      string `json:"filament_cooling_initial_initial"`
	FilamentCoolingMoves             string `json:"filament_cooling_moves"`
	FilamentCost                     string `json:"filament_cost"`
	FilamentDensity                  string `json:"filament_density"`
	FilamentDeretractionSpeed        string `json:"filament_deretraction_speed"`
	FilamentDiameter                 string `json:"filament_diameter"`
	FilamentEndGcode                 string `json:"filament_end_gcode"`
	FilamentFlowRatio                string `json:"filament_flow_ratio"`
	FilamentIsSupport                string `json:"filament_is_support"`
	FilamentLoadTime                 string `json:"filament_load_time"`
	FilamentLoadingSpeed             string `json:"filament_loading_speed"`
	FilamentLoadingSpeedStart        string `json:"filament_loading_speed_start"`
	FilamentMaxVolumetricSpeed       string `json:"filament_max_volumetric_speed"`
	FilamentMinimalPurgeOnWipeTower  string `json:"filament_minimal_purge_on_wipe_tower"`
	FilamentMultitoolRamming         string `json:"filament_multitool_ramming"`
	FilamentMultitoolRammingFlow     string `json:"filament_multitool_ramming_flow"`
	FilamentMultitoolRammingVolume   string `json:"filament_multitool_ramming_volume"`
	FilamentNotes                    string `json:"filament_notes"` // This is the stringified JSON from slicer
	FilamentRammingParameters        string `json:"filament_ramming_parameters"`
	FilamentRetractBeforeWipe        string `json:"filament_retract_before_wipe"`
	FilamentRetractLiftAbove         string `json:"filament_retract_lift_above"`
	FilamentRetractLiftBelow         string `json:"filament_retract_lift_below"`
	FilamentRetractLiftEnforce       string `json:"filament_retract_lift_enforce"`
	FilamentRetractRestartExtra      string `json:"filament_retract_restart_extra"`
	FilamentRetractWhenChangingLayer string `json:"filament_retract_when_changing_layer"`
	FilamentRetractionLength         string `json:"filament_retraction_length"`
	FilamentRetractionMinimumTravel  string `json:"filament_retraction_minimum_travel"`
	FilamentRetractionSpeed          string `json:"filament_retraction_speed"`
	FilamentShrink                   string `json:"filament_shrink"`
	FilamentShrinkageCompensationZ   string `json:"filament_shrinkage_compensation_z"`
	FilamentSoluble                  string `json:"filament_soluble"`
	FilamentStartGcode               string `json:"filament_start_gcode"`
	FilamentToolchangeDelay          string `json:"filament_toolchange_delay"`
	FilamentType                     string `json:"filament_type"` // This will be the Creality-specific type
	FilamentUnloadTime               string `json:"filament_unload_time"`
	FilamentUnloadingSpeed           string `json:"filament_unloading_speed"`
	FilamentUnloadingSpeedStart      string `json:"filament_unloading_speed_start"`
	FilamentVendor                   string `json:"filament_vendor"` // This will be the Creality-specific vendor
	FilamentWipe                     string `json:"filament_wipe"`
	FilamentWipeDistance             string `json:"filament_wipe_distance"`
	FilamentZHop                     string `json:"filament_z_hop"`
	FilamentZHopTypes                string `json:"filament_z_hop_types"`
	FullFanSpeedLayer                string `json:"full_fan_speed_layer"`
	HotPlateTemp                     string `json:"hot_plate_temp"`
	HotPlateTempInitialLayer         string `json:"hot_plate_temp_initial_layer"`
	Inherits                         string `json:"inherits"`
	MaterialFlowDependentTemperature string `json:"material_flow_dependent_temperature"`
	MaterialFlowTempGraph            string `json:"material_flow_temp_graph"`
	NozzleTemperature                string `json:"nozzle_temperature"`
	NozzleTemperatureInitialLayer    string `json:"nozzle_temperature_initial_layer"`
	NozzleTemperatureRangeHigh       string `json:"nozzle_temperature_range_high"`
	NozzleTemperatureRangeLow        string `json:"nozzle_temperature_range_low"`
	OverhangFanSpeed                 string `json:"overhang_fan_speed"`
	OverhangFanThreshold             string `json:"overhang_fan_threshold"`
	PressureAdvance                  string `json:"pressure_advance"`
	ReduceFanStopStartFreq           string `json:"reduce_fan_stop_start_freq"`
	RequiredNozzleHRC                string `json:"required_nozzle_HRC"`
	SlowDownForLayerCooling          string `json:"slow_down_for_layer_cooling"`
	WarmupStartLayer                 string `json:"warmup_start_layer"`
	SlowDownLayerTime                string `json:"slow_down_layer_time"`
	SlowDownMinSpeed                 string `json:"slow_down_min_speed"`
	SupportMaterialInterfaceFanSpeed string `json:"support_material_interface_fan_speed"`
	TemperatureVitrification         string `json:"temperature_vitrification"`
	TexturedPlateTemp                string `json:"textured_plate_temp"`
	TexturedPlateTempInitialLayer    string `json:"textured_plate_temp_initial_layer"`
}

// BaseInfo holds the basic identifying information for a filament.
type BaseInfo struct {
	ID            string   `json:"id"`
	Brand         string   `json:"brand"`
	Name          string   `json:"name"`
	MaterialType  string   `json:"meterialType"` // Note: "meterialType" as in source JSON
	Colors        []string `json:"colors"`
	Density       float64  `json:"density"`
	Diameter      string   `json:"diameter"`
	CostPerMeter  int      `json:"costPerMeter"`
	WeightPerMeter int      `json:"weightPerMeter"`
	Rank          int      `json:"rank"`
	MinTemp       int      `json:"minTemp"`
	MaxTemp       int      `json:"maxTemp"`
	IsSoluble     bool     `json:"isSoluble"`
	IsSupport     bool     `json:"isSuppoert"` // Note: "isSuppoert" as in source JSON
	ShrinkageRate int      `json:"shrinkageRate"`
	SofteningTemp int      `json:"softeningTemp"`
	DryingTemp    int      `json:"dryingTemp"`
	DryingTime    int      `json:"dryingTime"`
}

// MaterialOptions represents the structure of material_option.json.
// The inner values are newline-separated strings of names.
type MaterialOptions map[string]map[string]string

// LoadDefaultDatabase reads the default material_database.json file.
// This function is kept but will not be used by main.go due to embedding.
func LoadDefaultDatabase(filePath string) (*MaterialDatabase, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read default material database file %s: %w", filePath, err)
	}

	var db MaterialDatabase
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, fmt.Errorf("failed to parse default material database JSON: %w", err)
	}
	return &db, nil
}

// LoadDefaultDatabaseFromBytes loads the material database from a byte slice.
func LoadDefaultDatabaseFromBytes(data []byte) (*MaterialDatabase, error) {
	var db MaterialDatabase
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, fmt.Errorf("failed to parse material database from bytes: %w", err)
	}
	return &db, nil
}

// LoadDefaultOptions reads the default material_option.json file.
// This function is kept but will not be used by main.go due to embedding.
func LoadDefaultOptions(filePath string) (MaterialOptions, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read default material options file %s: %w", filePath, err)
	}

	var options MaterialOptions
	if err := json.Unmarshal(data, &options); err != nil {
		return nil, fmt.Errorf("failed to parse default material options JSON: %w", err)
	}
	return options, nil
}

// LoadDefaultOptionsFromBytes loads the material options from a byte slice.
func LoadDefaultOptionsFromBytes(data []byte) (MaterialOptions, error) {
	var options MaterialOptions
	if err := json.Unmarshal(data, &options); err != nil {
		return nil, fmt.Errorf("failed to parse material options from bytes: %w", err)
	}
	return options, nil
}

// ConvertToCrealityFormat converts a normalized slicer profile into a CrealityFilamentData structure.
func ConvertToCrealityFormat(slicerProfileData map[string]string, notes *profiles.FilamentNotes) (*FilamentProfileEntry, error) {
	// Marshal the original filamentNotes struct back into a JSON string, then escape it.
	notesBytes, err := json.Marshal(notes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filament notes: %w", err)
	}
	quotedNotes := strconv.Quote(string(notesBytes)) // Add quotes around the JSON string

	newEntry := &FilamentProfileEntry{
		EngineVersion:  slicerProfileData["version"],
		PrinterIntName: "F008",          // Hardcoded from original JS
		NozzleDiameter: []string{"0.4"}, // Hardcoded from original JS
		KVParam:        KVParam{},
		Base: BaseInfo{
			ID:           notes.ID,
			Brand:        notes.Vendor,
			Name:         notes.Name,
			MaterialType: notes.Type,
			Colors:       []string{"#ffffff"}, // Default color, can be enhanced
			Diameter:     slicerProfileData["filament_diameter"],
		},
	}

	assignKVParam := func(key string) string {
		if val, ok := slicerProfileData[key]; ok {
			return val
		}
		return ""
	}

	newEntry.KVParam = KVParam{
		ActivateAirFiltration:            assignKVParam("activate_air_filtration"),
		ActivateChamberTempControl:       assignKVParam("activate_chamber_temp_control"),
		AdditionalCoolingFanSpeed:        assignKVParam("additional_cooling_fan_speed"),
		ChamberTemperature:               assignKVParam("chamber_temperature"),
		CloseFanTheFirstXLayers:          assignKVParam("close_fan_the_first_x_layers"),
		CompatiblePrinters:               assignKVParam("compatible_printers"),
		CompatiblePrintersCondition:      assignKVParam("compatible_printers_condition"),
		CompatiblePrints:                 assignKVParam("compatible_prints"),
		CompatiblePrintsCondition:        assignKVParam("compatible_prints_condition"),
		CompletePrintExhaustFanSpeed:     assignKVParam("complete_print_exhaust_fan_speed"),
		CoolCdsFanStartAtHeight:          assignKVParam("cool_cds_fan_start_at_height"),
		CoolPlateTemp:                    assignKVParam("cool_plate_temp"),
		CoolPlateTempInitialLayer:        assignKVParam("cool_plate_temp_initial_layer"),
		CoolSpecialCdsFanSpeed:           assignKVParam("cool_special_cds_fan_speed"),
		DefaultFilamentColour:            assignKVParam("default_filament_colour"),
		DuringPrintExhaustFanSpeed:       assignKVParam("during_print_exhaust_fan_speed"),
		EnableOverhangBridgeFan:          assignKVParam("enable_overhang_bridge_fan"),
		EnablePressureAdvance:            assignKVParam("enable_pressure_advance"),
		EnableSpecialAreaAdditionalCoolingFan: assignKVParam("enable_special_area_additional_cooling_fan"),
		EngPlateTemp:                     assignKVParam("eng_plate_temp"),
		EngPlateTempInitialLayer:         assignKVParam("eng_plate_temp_initial_layer"),
		EpoxyResinPlateTemp:              assignKVParam("epoxy_resin_plate_temp"),
		EpoxyResinPlateTempInitialLayer:  assignKVParam("epoxy_resin_plate_initial_layer"),
		FanCoolingLayerTime:              assignKVParam("fan_cooling_layer_time"),
		FanMaxSpeed:                      assignKVParam("fan_max_speed"),
		FanMinSpeed:                      assignKVParam("fan_min_speed"),
		FilamentCoolingFinalSpeed:        assignKVParam("filament_cooling_final_speed"),
		FilamentCoolingInitialSpeed:      assignKVParam("filament_cooling_initial_initial"),
		FilamentCoolingMoves:             assignKVParam("filament_cooling_moves"),
		FilamentCost:                     assignKVParam("filament_cost"),
		FilamentDensity:                  assignKVParam("filament_density"),
		FilamentDeretractionSpeed:        assignKVParam("filament_deretraction_speed"),
		FilamentDiameter:                 assignKVParam("filament_diameter"),
		FilamentEndGcode:                 assignKVParam("filament_end_gcode"),
		FilamentFlowRatio:                assignKVParam("filament_flow_ratio"),
		FilamentIsSupport:                assignKVParam("filament_is_support"),
		FilamentLoadTime:                 assignKVParam("filament_load_time"),
		FilamentLoadingSpeed:             assignKVParam("filament_loading_speed"),
		FilamentLoadingSpeedStart:        assignKVParam("filament_loading_speed_start"),
		FilamentMaxVolumetricSpeed:       assignKVParam("filament_max_volumetric_speed"),
		FilamentMinimalPurgeOnWipeTower:  assignKVParam("filament_minimal_purge_on_wipe_tower"),
		FilamentMultitoolRamming:         assignKVParam("filament_multitool_ramming"),
		FilamentMultitoolRammingFlow:     assignKVParam("filament_multitool_ramming_flow"),
		FilamentMultitoolRammingVolume:   assignKVParam("filament_multitool_ramming_volume"),
		FilamentNotes:                    quotedNotes,
		FilamentRammingParameters:        assignKVParam("filament_ramming_parameters"),
		FilamentRetractBeforeWipe:        assignKVParam("filament_retract_before_wipe"),
		FilamentRetractLiftAbove:         assignKVParam("filament_retract_lift_above"),
		FilamentRetractLiftBelow:         assignKVParam("filament_retract_lift_below"),
		FilamentRetractLiftEnforce:       assignKVParam("filament_retract_lift_enforce"),
		FilamentRetractRestartExtra:      assignKVParam("filament_retract_restart_extra"),
		FilamentRetractWhenChangingLayer: assignKVParam("filament_retract_when_changing_layer"),
		FilamentRetractionLength:         assignKVParam("filament_retraction_length"),
		FilamentRetractionMinimumTravel:  assignKVParam("filament_retraction_minimum_travel"),
		FilamentRetractionSpeed:          assignKVParam("filament_retraction_speed"),
		FilamentShrink:                   assignKVParam("filament_shrink"),
		FilamentShrinkageCompensationZ:   assignKVParam("filament_shrinkage_compensation_z"),
		FilamentSoluble:                  assignKVParam("filament_soluble"),
		FilamentStartGcode:               assignKVParam("filament_start_gcode"),
		FilamentToolchangeDelay:          assignKVParam("filament_toolchange_delay"),
		FilamentType:                     assignKVParam("filament_type"),
		FilamentUnloadTime:               assignKVParam("filament_unload_time"),
		FilamentUnloadingSpeed:           assignKVParam("filament_unloading_speed"),
		FilamentUnloadingSpeedStart:      assignKVParam("filament_unloading_speed_start"),
		FilamentVendor:                   assignKVParam("filament_vendor"),
		FilamentWipe:                     assignKVParam("filament_wipe"),
		FilamentWipeDistance:             assignKVParam("filament_wipe_distance"),
		FilamentZHop:                     assignKVParam("filament_z_hop"),
		FilamentZHopTypes:                assignKVParam("filament_z_hop_types"),
		FullFanSpeedLayer:                assignKVParam("full_fan_speed_layer"),
		HotPlateTemp:                     assignKVParam("hot_plate_temp"),
		HotPlateTempInitialLayer:         assignKVParam("hot_plate_temp_initial_layer"),
		Inherits:                         assignKVParam("inherits"),
		MaterialFlowDependentTemperature: assignKVParam("material_flow_dependent_temperature"),
		MaterialFlowTempGraph:            assignKVParam("material_flow_temp_graph"),
		NozzleTemperature:                assignKVParam("nozzle_temperature"),
		NozzleTemperatureInitialLayer:    assignKVParam("nozzle_temperature_initial_layer"),
		NozzleTemperatureRangeHigh:       assignKVParam("nozzle_temperature_range_high"),
		NozzleTemperatureRangeLow:        assignKVParam("nozzle_temperature_range_low"),
		OverhangFanSpeed:                 assignKVParam("overhang_fan_speed"),
		OverhangFanThreshold:             assignKVParam("overhang_fan_threshold"),
		PressureAdvance:                  assignKVParam("pressure_advance"),
		ReduceFanStopStartFreq:           assignKVParam("reduce_fan_stop_start_freq"),
		RequiredNozzleHRC:                assignKVParam("required_nozzle_HRC"),
		SlowDownForLayerCooling:          assignKVParam("slow_down_for_layer_cooling"),
		WarmupStartLayer:                 assignKVParam("warmup_start_layer"),
		SlowDownLayerTime:                assignKVParam("slow_down_layer_time"),
		SlowDownMinSpeed:                 assignKVParam("slow_down_min_speed"),
		SupportMaterialInterfaceFanSpeed: assignKVParam("support_material_interface_fan_speed"),
		TemperatureVitrification:         assignKVParam("temperature_vitrification"),
		TexturedPlateTemp:                assignKVParam("textured_plate_temp"),
		TexturedPlateTempInitialLayer:    assignKVParam("textured_plate_temp_initial_layer"),
	}

	if notes != nil {
		if notes.Type != "" {
			newEntry.KVParam.FilamentType = notes.Type
		}
		if notes.Vendor != "" {
			newEntry.KVParam.FilamentVendor = notes.Vendor
		}
	} else {
		newEntry.KVParam.FilamentType = assignKVParam("filament_type")
		newEntry.KVParam.FilamentVendor = assignKVParam("filament_vendor")
	}

	newEntry.Base.IsSoluble = assignKVParam("filament_soluble") == "1"
	newEntry.Base.IsSupport = assignKVParam("filament_is_support") == "1"

	if density, err := strconv.ParseFloat(assignKVParam("filament_density"), 64); err == nil {
		newEntry.Base.Density = density
	}
	if cost, err := strconv.Atoi(assignKVParam("filament_cost")); err == nil {
		newEntry.Base.CostPerMeter = cost
	}

	if minTemp, err := strconv.Atoi(assignKVParam("nozzle_temperature_range_low")); err == nil {
		newEntry.Base.MinTemp = minTemp
	}
	if maxTemp, err := strconv.Atoi(assignKVParam("nozzle_temperature_range_high")); err == nil {
		newEntry.Base.MaxTemp = maxTemp
	}
	if shrinkageRate, err := strconv.Atoi(strings.TrimSuffix(assignKVParam("filament_shrink"), "%")); err == nil {
		newEntry.Base.ShrinkageRate = shrinkageRate
	}
	if softeningTemp, err := strconv.Atoi(assignKVParam("temperature_vitrification")); err == nil {
		newEntry.Base.SofteningTemp = softeningTemp
	}
	if dryingTemp, err := strconv.Atoi(assignKVParam("hot_plate_temp")); err == nil {
		newEntry.Base.DryingTemp = dryingTemp
	}
	if dryingTime, err := strconv.Atoi(assignKVParam("slow_down_layer_time")); err == nil {
		newEntry.Base.DryingTime = dryingTime
	}

	return newEntry, nil
}

// AddProfileToDatabase adds a new filament profile entry to the MaterialDatabase.
// It replaces an existing entry if one with the same ID is found, otherwise appends.
func AddProfileToDatabase(db *MaterialDatabase, newProfile *FilamentProfileEntry, version string) {
	found := false
	for i, entry := range db.Result.List {
		if entry.Base.ID == newProfile.Base.ID {
			db.Result.List[i] = *newProfile // Replace existing
			found = true
			break
		}
	}
	if !found {
		db.Result.List = append(db.Result.List, *newProfile) // Add new
		db.Result.Count = len(db.Result.List) // Update count
	}
	db.Result.Version = version // Update the database version
}

// UpdateOptions updates the MaterialOptions with a new filament entry.
// It handles adding new vendors, new filament types for existing vendors,
// and appending new names to existing vendor/type combinations.
func UpdateOptions(options MaterialOptions, notes *profiles.FilamentNotes) {
	if options == nil {
		return // Should not happen if LoadDefaultOptions is called
	}

	vendor := notes.Vendor
	if vendor == "" { // Fallback if notes.Vendor is empty
		vendor = "Generic" // Or some other default
	}

	filamentType := notes.Type
	if filamentType == "" { // Fallback if notes.Type is empty
		filamentType = "PLA" // Or some other default
	}

	name := notes.Name
	if name == "" { // Fallback if notes.Name is empty
		name = "Custom Filament" // Or some other default
	}

	if _, ok := options[vendor]; !ok {
		options[vendor] = make(map[string]string)
	}

	if existingNames, ok := options[vendor][filamentType]; ok {
		names := strings.Split(existingNames, "\n")
		for _, n := range names {
			if n == name {
				return // Name already exists, do nothing
			}
		}
		options[vendor][filamentType] = existingNames + "\n" + name
	} else {
		options[vendor][filamentType] = name
	}
}

// Removed SaveDatabase function - no longer persisting to local disk.
/*
func SaveDatabase(db *MaterialDatabase, filePath string) error {
	data, err := json.MarshalIndent(db, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal material database JSON: %w", err)
	}
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write material database file %s: %w", filePath, err)
	}
	return nil
}
*/

// Removed SaveOptions function - no longer persisting to local disk.
/*
func SaveOptions(options MaterialOptions, filePath string) error {
	data, err := json.MarshalIndent(options, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal material options JSON: %w", err)
	}
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write material options file %s: %w", filePath, err)
	}
	return nil
}
*/

// MarshalDatabase converts a MaterialDatabase struct to its JSON byte representation.
func MarshalDatabase(db *MaterialDatabase) ([]byte, error) {
	data, err := json.MarshalIndent(db, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal material database to bytes: %w", err)
	}
	return data, nil
}

// MarshalOptions converts MaterialOptions to its JSON byte representation.
func MarshalOptions(options MaterialOptions) ([]byte, error) {
	data, err := json.MarshalIndent(options, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal material options to bytes: %w", err)
	}
	return data, nil
}

