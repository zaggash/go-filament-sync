package profiles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// FilamentNotes represents the JSON structure expected within the "filament_notes" field of the slicer profile.
type FilamentNotes struct {
	ID     string `json:"id"`
	Vendor string `json:"vendor"`
	Type   string `json:"type"`
	Name   string `json:"name"`
}

// SlicerFilamentProfile represents the structure of a filament profile JSON from OrcaSlicer/Creality Print.
type SlicerFilamentProfile struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	FilamentNotes []string `json:"filament_notes"` // This contains the embedded JSON string
	RawData map[string]interface{} `json:"-"` // Store raw data for dynamic access
}

// UnmarshalJSON custom unmarshaler to capture all raw data.
func (s *SlicerFilamentProfile) UnmarshalJSON(data []byte) error {
	type Alias SlicerFilamentProfile
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.RawData); err != nil {
		return err
	}
	return nil
}

// ReadSlicerProfile reads a JSON file and unmarshals it into a SlicerFilamentProfile.
func ReadSlicerProfile(filePath string) (*SlicerFilamentProfile, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var profile SlicerFilamentProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from %s: %w", filePath, err)
	}
	return &profile, nil
}

// NormalizeSlicerProfile processes the raw SlicerFilamentProfile.
// It extracts single values from arrays and parses the embedded filament_notes.
// Returns a map of string values and the parsed FilamentNotes struct.
func NormalizeSlicerProfile(slicerProfile *SlicerFilamentProfile) (map[string]string, *FilamentNotes, error) {
	normalizedData := make(map[string]string)
	var notes *FilamentNotes

	// Process raw data to flatten array fields and convert to string
	for key, value := range slicerProfile.RawData {
		if key == "name" || key == "version" || key == "filament_notes" {
			continue // Skip known fields handled explicitly or where raw type is used
		}

		// Handle values that are arrays (common in slicer profiles)
		if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
			if strVal, ok := arr[0].(string); ok {
				normalizedData[key] = strVal
			} else if numVal, ok := arr[0].(float64); ok {
				normalizedData[key] = fmt.Sprintf("%v", numVal)
			} else if boolVal, ok := arr[0].(bool); ok {
				normalizedData[key] = fmt.Sprintf("%v", boolVal)
			} else {
				// For complex types like "material_flow_temp_graph", keep as string representation
				normalizedData[key] = fmt.Sprintf("%v", arr[0])
			}
		} else { // Handle values that are not arrays
			if strVal, ok := value.(string); ok {
				normalizedData[key] = strVal
			} else if numVal, ok := value.(float64); ok {
				normalizedData[key] = fmt.Sprintf("%v", numVal)
			} else if boolVal, ok := value.(bool); ok {
				normalizedData[key] = fmt.Sprintf("%v", boolVal)
			} else {
				normalizedData[key] = fmt.Sprintf("%v", value)
			}
		}
	}

	// Manually add the "name" and "version" from the SlicerFilamentProfile struct itself
	normalizedData["name"] = slicerProfile.Name
	normalizedData["version"] = slicerProfile.Version

	// Parse filament_notes from its string content
	if len(slicerProfile.FilamentNotes) > 0 && slicerProfile.FilamentNotes[0] != "" {
		var parsedNotes FilamentNotes
		noteString := strings.Trim(slicerProfile.FilamentNotes[0], "\"")
		err := json.Unmarshal([]byte(noteString), &parsedNotes)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse filament_notes JSON: %w, raw note: %s", err, noteString)
		}
		notes = &parsedNotes
	} else {
		return nil, nil, fmt.Errorf("filament_notes is missing or empty in the profile")
	}

	// Add the individual notes fields to normalizedData for easier access in conversion (though not directly used in creality conversion now)
	if notes != nil {
		normalizedData["filament_notes.id"] = notes.ID
		normalizedData["filament_notes.vendor"] = notes.Vendor
		normalizedData["filament_notes.type"] = notes.Type
		normalizedData["filament_notes.name"] = notes.Name
	}

	return normalizedData, notes, nil
}

// GetSlicerProfileDir determines the correct slicer profile directory based on OS and slicer type.
func GetSlicerProfileDir(slicerType, userID string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	baseDir := ""
	profilePath := ""

	switch runtime.GOOS {
	case "darwin": // macOS
		baseDir = filepath.Join(homeDir, "Library", "Application Support")
		if slicerType == "orca" {
			profilePath = filepath.Join(baseDir, "OrcaSlicer", "user", userID, "filament", "base")
		} else if slicerType == "creality" {
			profilePath = filepath.Join(baseDir, "Creality", "Creality Print", "6.0", "user", userID, "filament", "base")
		} else {
			return "", fmt.Errorf("unsupported slicer type: %s. Must be 'orca' or 'creality'", slicerType)
		}
	case "linux":
		baseDir = filepath.Join(homeDir, ".config")
		if slicerType == "orca" {
			profilePath = filepath.Join(baseDir, "OrcaSlicer", "user", userID, "filament") // Corrected path based on user feedback
		} else if slicerType == "creality" {
			profilePath = filepath.Join(baseDir, "Creality", "Creality Print", "6.0", "user", userID, "filament", "base")
		} else {
			return "", fmt.Errorf("unsupported slicer type: %s. Must be 'orca' or 'creality'", slicerType)
		}
	case "windows":
		baseDir = filepath.Join(homeDir, "AppData", "Roaming")
		if slicerType == "orca" {
			profilePath = filepath.Join(baseDir, "OrcaSlicer", "user", userID, "filament", "base")
		} else if slicerType == "creality" {
			profilePath = filepath.Join(baseDir, "Creality", "Creality Print", "6.0", "user", userID, "filament", "base")
		} else {
			return "", fmt.Errorf("unsupported slicer type: %s. Must be 'orca' or 'creality'", slicerType)
		}
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("slicer profile directory does not exist: %s", profilePath)
	}

	return profilePath, nil
}

// LoadCustomProfiles reads JSON files from the given directory and filters them.
// It checks if the "filament_notes" field is present and non-empty as a JSON string.
func LoadCustomProfiles(dir string) ([]string, error) {
	var profilePaths []string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("Skipping %s: failed to read file: %v", file.Name(), err)
			continue
		}

		var rawProfile map[string]interface{}
		if err := json.Unmarshal(data, &rawProfile); err != nil {
			log.Printf("Skipping %s: failed to unmarshal JSON: %v", file.Name(), err)
			continue
		}

		if notesVal, ok := rawProfile["filament_notes"]; ok {
			if notesArr, isArray := notesVal.([]interface{}); isArray && len(notesArr) > 0 {
				if noteStr, isString := notesArr[0].(string); isString && strings.TrimSpace(noteStr) != "" && strings.Contains(noteStr, `"id":`) {
					var tempNotes FilamentNotes
					if err := json.Unmarshal([]byte(strings.Trim(noteStr, `"`)), &tempNotes); err == nil && tempNotes.ID != "" {
						profilePaths = append(profilePaths, filePath)
					} else {
						log.Printf("Ignoring %s: filament_notes is not a valid JSON string or missing 'id'", file.Name())
					}
				} else {
					log.Printf("Ignoring %s: filament_notes is empty or not a string array", file.Name())
				}
			} else {
				log.Printf("Ignoring %s: filament_notes field is missing or not in expected array format", file.Name())
			}
		} else {
			log.Printf("Ignoring %s: missing required 'filament_notes' field", file.Name())
		}
	}

	return profilePaths, nil
}


