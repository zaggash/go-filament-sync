package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"filament-sync-tool/cli/config"   // Import our new config package
	"filament-sync-tool/cli/creality" // Import our new creality package
	"filament-sync-tool/cli/profiles"  // Import our new profiles package
	"filament-sync-tool/cli/scp"      // Import our new scp package
)

// Global variable to hold parsed config, populated by config.LoadConfig()
var appConfig *config.ToolConfig

// Global variables for in-memory databases
var (
	materialDB      *creality.MaterialDatabase
	materialOptions creality.MaterialOptions
)

//go:embed data/*
var embeddedData embed.FS // Embed the 'data' directory into the binary

// init function runs automatically before main()
func init() {
	// Initialize logging
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Load configuration from command-line arguments
	appConfig = config.LoadConfig()

	// Read default JSONs from embedded filesystem
	sourceDBData, err := embeddedData.ReadFile("data/material_database.json") // embed.FS uses ReadFile directly
	if err != nil {
		log.Fatalf("Failed to read embedded material_database.json: %v", err)
	}
	sourceOptData, err := embeddedData.ReadFile("data/material_option.json") // embed.FS uses ReadFile directly
	if err != nil {
		log.Fatalf("Failed to read embedded material_option.json: %v", err)
	}

	// Load initial database and options into memory from embedded data
	materialDB, err = creality.LoadDefaultDatabaseFromBytes(sourceDBData)
	if err != nil {
		log.Fatalf("Failed to load material database from embedded data: %v", err)
	}
	log.Println("Material database loaded from embedded data successfully.")

	materialOptions, err = creality.LoadDefaultOptionsFromBytes(sourceOptData)
	if err != nil {
		log.Fatalf("Failed to load material options from embedded data: %v", err)
	}
	log.Println("Material options loaded from embedded data successfully.")
}

func main() {
	// printerTargetDir is the remote path on the printer
	printerTargetDir := "/mnt/UDISK/creality/userdata/box"

	// Discover and load custom slicer profiles
	profileDir, err := profiles.GetSlicerProfileDir(appConfig.Slicer, appConfig.UserID, appConfig.Flatpak)
	if err != nil {
		log.Fatalf("Error determining slicer profile directory: %v", err)
	}
	log.Printf("Scanning for %s profiles in: %s (Flatpak mode: %v)", appConfig.Slicer, profileDir, appConfig.Flatpak)

	slicerProfilePaths, err := profiles.LoadCustomProfiles(profileDir)
	if err != nil {
		log.Fatalf("Error loading custom profiles from %s: %v", profileDir, err)
	}

	if len(slicerProfilePaths) == 0 {
		log.Println("No custom filament profiles found with required 'filament_notes'. Skipping synchronization.")
		log.Println("Please ensure your custom profiles have 'filament_notes' as described in the original README:")
		log.Println("https://github.com/HurricanePrint/Filament-Sync#creating-custom-filament-presets")
		return // Exit if no profiles to sync
	}

	log.Printf("Found %d custom profiles. Processing and preparing for transfer...", len(slicerProfilePaths))

	// Process each custom profile
	for _, path := range slicerProfilePaths {
		log.Printf("Processing profile: %s", path)
		slicerProfile, err := profiles.ReadSlicerProfile(path)
		if err != nil {
			log.Printf("Skipping profile %s due to read error: %v", path, err)
			continue
		}

		normalizedData, filamentNotes, err := profiles.NormalizeSlicerProfile(slicerProfile)
		if err != nil {
			log.Printf("Skipping profile %s due to normalization error: %v", path, err)
			continue
		}

		crealityEntry, err := creality.ConvertToCrealityFormat(normalizedData, filamentNotes)
		if err != nil {
			log.Printf("Skipping profile %s due to conversion error: %v", path, err)
			continue
		}

		// Update in-memory databases with the new/updated entry
		currentTimestamp := fmt.Sprintf("%d", time.Now().Unix())
		creality.AddProfileToDatabase(materialDB, crealityEntry, currentTimestamp)
		creality.UpdateOptions(materialOptions, filamentNotes)
	}

	// --- Prepare Data for SCP (from memory) ---
	updatedDBBytes, err := creality.MarshalDatabase(materialDB)
	if err != nil {
		log.Fatalf("Failed to marshal updated material database to bytes: %v", err)
	}

	updatedOptBytes, err := creality.MarshalOptions(materialOptions)
	if err != nil {
		log.Fatalf("Failed to marshal updated material options to bytes: %v", err)
	}

	// --- SCP Transfer to Printer ---
	log.Println("Initiating SCP transfer to printer (from memory)...")
	scpClient, err := scp.NewSCPClient(appConfig.PrinterIP, appConfig.User, appConfig.Password)
	if err != nil {
		log.Fatalf("Failed to initialize SCP client: %v", err)
	}

	// Establish and defer close SSH connection once for transfer operations
	if err := scpClient.Connect(); err != nil {
		log.Fatalf("Failed to establish SSH connection to printer: %v", err)
	}
	defer scpClient.Close()

	// Check if the remote directory exists
	_, err = scpClient.CheckRemoteDirectory(printerTargetDir)
	if err != nil {
		log.Fatalf("Error checking remote directory %s: %v", printerTargetDir, err)
	}

	// Upload material_database.json directly from bytes
	dbReader := bytes.NewReader(updatedDBBytes)
	dbFileName := "material_database.json" // Hardcoded filename for remote upload
	dbFileSize := int64(len(updatedDBBytes))
	dbFileMode := os.FileMode(0644) // Default permissions for the file on the printer

	remoteDBPath := filepath.Join(printerTargetDir, dbFileName)
	if err := scpClient.UploadFile(dbReader, remoteDBPath, dbFileName, dbFileSize, dbFileMode); err != nil {
		log.Fatalf("Failed to upload %s to printer: %v", dbFileName, err)
	}
	log.Printf("Uploaded %s to printer.", dbFileName)

	// Upload material_option.json directly from bytes
	optReader := bytes.NewReader(updatedOptBytes)
	optFileName := "material_option.json" // Hardcoded filename for remote upload
	optFileSize := int64(len(updatedOptBytes))
	optFileMode := os.FileMode(0644) // Default permissions for the file on the printer

	remoteOptionsPath := filepath.Join(printerTargetDir, optFileName)
	if err := scpClient.UploadFile(optReader, remoteOptionsPath, optFileName, optFileSize, optFileMode); err != nil {
		log.Fatalf("Failed to upload %s to printer: %v", optFileName, err)
	}
	log.Printf("Uploaded %s to printer.", optFileName)

	log.Println("Filament profiles synchronized successfully with the printer!")
}

