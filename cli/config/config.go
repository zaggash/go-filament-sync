package config

import (
	"flag"
	"log"
	"os"
)

// ToolConfig holds the application-wide configuration parameters from command-line flags.
type ToolConfig struct {
	PrinterIP string
	User      string
	Password  string
	Slicer    string // "orca" or "creality"
	UserID    string // User ID for slicer profile paths
}

// LoadConfig parses command-line arguments and returns a populated ToolConfig.
// It performs basic validation and exits if essential arguments are missing.
func LoadConfig() *ToolConfig {
	// Initialize logging (can be done here or in main)
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Define command-line flags
	printerIP := flag.String("printer-ip", "", "IP address of the Creality printer (required)")
	user := flag.String("user", "root", "Username for SSH connection to printer")
	password := flag.String("password", "creality_2024", "Password for SSH connection to printer")
	slicerType := flag.String("slicer", "orca", "Specify the slicer type: 'orca' or 'creality'")
	userID := flag.String("userid", "default", "Specify the user ID for the slicer profile folder")

	// Parse the command-line flags
	flag.Parse()

	cfg := &ToolConfig{
		PrinterIP: *printerIP,
		User:      *user,
		Password:  *password,
		Slicer:    *slicerType,
		UserID:    *userID,
	}

	// Basic validation for critical configuration (printer IP is now required)
	if cfg.PrinterIP == "" {
		log.Fatal("Error: Printer IP address is required. Use --printer-ip flag.")
	}

	log.Printf("Tool Config: PrinterIP=%s, User=%s, Slicer=%s, UserID=%s",
		cfg.PrinterIP, cfg.User, cfg.Slicer, cfg.UserID)

	return cfg
}


