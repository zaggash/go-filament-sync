package config

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// ToolConfig holds the application-wide configuration parameters from command-line flags.
type ToolConfig struct {
	PrinterIP   string
	User        string
	Password    string
	ProfilePath string
}

// LoadConfig parses command-line arguments and returns a populated ToolConfig.
// It performs basic validation and exits if essential arguments are missing.
func LoadConfig() *ToolConfig {
	// Initialize logging (can be done here or in main)
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Define command-line flags
	profilePath := flag.String("profile-path", "", "Path to slicer filament profile directory (required)")
	printerIP := flag.String("printer-ip", "", "IP address of the Creality printer (required)")
	user := flag.String("user", "root", "Username for SSH connection to printer")
	password := flag.String("password", "creality_2024", "Password for SSH connection to printer")

	// Install custom usage handler with migration note BEFORE parsing
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of filament-sync-tool:\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nMigration note: --userid, --flatpak, and --slicer flags have been removed.\n")
		fmt.Fprintf(os.Stderr, "Use --profile-path with the explicit path to your filament profile directory.\n\n")
		fmt.Fprintf(os.Stderr, "Example paths:\n")
		fmt.Fprintf(os.Stderr, "  OrcaSlicer (Linux):    ~/.config/OrcaSlicer/user/default/filament/base\n")
		fmt.Fprintf(os.Stderr, "  OrcaSlicer (macOS):    ~/Library/Application Support/OrcaSlicer/user/default/filament/base\n")
		fmt.Fprintf(os.Stderr, "  OrcaSlicer (Windows):  %%APPDATA%%\\OrcaSlicer\\user\\default\\filament\\base\n")
		fmt.Fprintf(os.Stderr, "  Creality Print (Linux):    ~/.config/Creality/Creality Print/6.0/user/default/filament/base   (replace 6.0 with your installed version)\n")
		fmt.Fprintf(os.Stderr, "  Creality Print (macOS):    ~/Library/Application Support/Creality/Creality Print/6.0/user/default/filament/base   (replace 6.0 with your installed version)\n")
		fmt.Fprintf(os.Stderr, "  Creality Print (Windows):  %%APPDATA%%\\Creality\\Creality Print\\6.0\\user\\default\\filament\\base   (replace 6.0 with your installed version)\n")
	}

	// Parse the command-line flags
	flag.Parse()

	// --profile-path is required; exit 2 so slicers can detect misconfiguration
	if *profilePath == "" {
		fmt.Fprintf(os.Stderr, "Error: --profile-path is required\n\n")
		flag.Usage()
		os.Exit(2)
	}

	// Verify the path exists
	stat, err := os.Stat(*profilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: profile path does not exist: %s\n", *profilePath)
		os.Exit(1)
	}

	// Verify the path is a directory, not a file
	if !stat.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: profile path is not a directory: %s\n", *profilePath)
		os.Exit(1)
	}

	// Validate required --printer-ip
	if *printerIP == "" {
		fmt.Fprintf(os.Stderr, "Error: --printer-ip is required\n\n")
		flag.Usage()
		os.Exit(2)
	}

	log.Printf("Tool Config: PrinterIP=%s, User=%s, ProfilePath=%s", *printerIP, *user, *profilePath)

	return &ToolConfig{
		PrinterIP:   *printerIP,
		User:        *user,
		Password:    *password,
		ProfilePath: *profilePath,
	}
}
