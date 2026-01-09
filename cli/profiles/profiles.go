package profiles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
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
	Name          string                 `json:"name"`
	Version       string                 `json:"version"`
	FilamentNotes []string               `json:"filament_notes"` // This contains the embedded JSON string
	RawData       map[string]interface{} `json:"-"`              // Store raw data for dynamic access
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
func NormalizeSlicerProfile(slicerProfile *SlicerFilamentProfile) (map[string]string, *FilamentNotes, error) {
	normalizedData := make(map[string]string)
	var notes *FilamentNotes

	// Process raw data to flatten array fields and convert to string
	for key, value := range slicerProfile.RawData {
		if key == "name" || key == "version" || key == "filament_notes" {
			continue
		}

		if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
			if strVal, ok := arr[0].(string); ok {
				normalizedData[key] = strVal
			} else if numVal, ok := arr[0].(float64); ok {
				normalizedData[key] = fmt.Sprintf("%v", numVal)
			} else if boolVal, ok := value.([]interface{})[0].(bool); ok {
				normalizedData[key] = fmt.Sprintf("%v", boolVal)
			} else {
				normalizedData[key] = fmt.Sprintf("%v", arr[0])
			}
		} else {
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

	normalizedData["name"] = slicerProfile.Name
	normalizedData["version"] = slicerProfile.Version

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

	if notes != nil {
		normalizedData["filament_notes.id"] = notes.ID
		normalizedData["filament_notes.vendor"] = notes.Vendor
		normalizedData["filament_notes.type"] = notes.Type
		normalizedData["filament_notes.name"] = notes.Name
	}

	return normalizedData, notes, nil
}

// GetSlicerProfileDir determines the correct slicer profile directory based on OS, slicer type, and the flatpak flag.
// It dynamically searches for the latest version folder for Creality Print.
func GetSlicerProfileDir(slicerType, userID string, isFlatpak bool) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Determine Base Directory based on OS
	var osBaseDir string
	switch runtime.GOOS {
	case "darwin":
		osBaseDir = filepath.Join(homeDir, "Library", "Application Support")
	case "linux":
		if isFlatpak {
			osBaseDir = homeDir
		} else {
			osBaseDir = filepath.Join(homeDir, ".config")
		}
	case "windows":
		osBaseDir = filepath.Join(homeDir, "AppData", "Roaming")
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Handle OrcaSlicer (Static Path)
	if slicerType == "orca" {
		var relPath string
		segment := filepath.Join("OrcaSlicer", "user", userID, "filament", "base")

		if isFlatpak && runtime.GOOS == "linux" {
			// Flatpak specific path for Orca
			relPath = filepath.Join(".var", "app", "io.github.softfever.OrcaSlicer", "config", segment)
		} else {
			relPath = segment
		}
		
		fullPath := filepath.Join(osBaseDir, relPath)
		return checkPath(fullPath)
	}

	// Handle Creality Print (Dynamic Version path)
	if slicerType == "creality" {
		var searchDir string
		
		if isFlatpak && runtime.GOOS == "linux" {
			searchDir = filepath.Join(osBaseDir, ".var", "app", "io.github.crealityofficial.CrealityPrint", "config", "Creality", "Creality Print")
		} else {
			searchDir = filepath.Join(osBaseDir, "Creality", "Creality Print")
		}

		// Read the directory to find version folders (e.g., "6.0", "7.0")
		entries, err := os.ReadDir(searchDir)
		if err != nil {
			// If the base folder is missing, provide a clear error
			return "", fmt.Errorf("could not find Creality Print installation directory at %s: %w", searchDir, err)
		}

		var versions []string
		for _, e := range entries {
			if e.IsDir() {
				versions = append(versions, e.Name())
			}
		}

		if len(versions) == 0 {
			return "", fmt.Errorf("no version folders found in %s", searchDir)
		}

		// Sort versions using Natural Sort (so 10.0 comes after 9.0)
		sort.Slice(versions, func(i, j int) bool {
			return isVersionLess(versions[i], versions[j])
		})

		// Pick the highest version
		latestVersion := versions[len(versions)-1]
		
		// Construct the final full path
		fullPath := filepath.Join(searchDir, latestVersion, "user", userID, "filament", "base")
		
		return checkPath(fullPath)
	}

	return "", fmt.Errorf("unsupported slicer type: %s", slicerType)
}

// LoadCustomProfiles reads JSON files from the given directory and filters them.
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
					unquotedNoteStr := strings.Trim(noteStr, `"`)

					if err := json.Unmarshal([]byte(unquotedNoteStr), &tempNotes); err == nil && tempNotes.ID != "" {
						profilePaths = append(profilePaths, filePath)
					} else {
						log.Printf("Ignoring %s: filament_notes invalid or missing 'id'. Inner content: %s", file.Name(), unquotedNoteStr)
					}
				} else {
					log.Printf("Ignoring %s: filament_notes empty/invalid", file.Name())
				}
			} else {
				log.Printf("Ignoring %s: filament_notes missing/invalid format", file.Name())
			}
		} else {
			log.Printf("Ignoring %s: missing required 'filament_notes' field", file.Name())
		}
	}

	return profilePaths, nil
}

// checkPath verifies if the calculated path actually exists on the disk.
func checkPath(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("slicer profile directory does not exist: %s", path)
	}
	return path, nil
}

// isVersionLess compares two version strings (like "6.0" and "10.1") naturally.
// It splits by "." and compares numeric segments.
func isVersionLess(v1, v2 string) bool {
	s1 := strings.Split(v1, ".")
	s2 := strings.Split(v2, ".")

	for i := 0; i < len(s1) && i < len(s2); i++ {
		n1, err1 := strconv.Atoi(s1[i])
		n2, err2 := strconv.Atoi(s2[i])

		// If both segments are numbers, compare numerically
		if err1 == nil && err2 == nil {
			if n1 != n2 {
				return n1 < n2
			}
			continue
		}
		// Fallback to string comparison for non-numeric segments
		if s1[i] != s2[i] {
			return s1[i] < s2[i]
		}
	}
	return len(s1) < len(s2)
}
