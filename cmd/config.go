package cmd

import (
	"flag"
	"fmt"
	"os"

	"stash/internal/config"
)

// RunConfig handles the config subcommand
func RunConfig(args []string) error {
	fs := flag.NewFlagSet("config", flag.ExitOnError)
	dirFlag := fs.String("dir", "", "Target directory for notes")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if *dirFlag == "" {
		fmt.Println("Usage: stash config --dir <path>")
		fmt.Println("\nOptions:")
		fmt.Println("  --dir string    Set the target directory for saved notes")
		return nil
	}

	if err := config.SetTargetDir(*dirFlag); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("Configuration saved. Notes will be saved to: %s\n", cfg.TargetDir)
	return nil
}

// ShowConfigHelp shows help for the config command
func ShowConfigHelp() {
	fmt.Println("Usage: stash config --dir <path>")
	fmt.Println("\nSet the target directory for saved notes.")
	fmt.Println("\nOptions:")
	fmt.Println("  --dir string    Target directory path (supports ~ for home directory)")
	fmt.Println("\nExamples:")
	fmt.Println("  stash config --dir ~/notes")
	fmt.Println("  stash config --dir /home/user/Documents/stash")
}

// CheckConfig checks if config exists and shows appropriate message
func CheckConfig() error {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Stash is not configured yet.")
		fmt.Println("\nPlease run: stash config --dir <path>")
		fmt.Println("\nExample:")
		fmt.Println("  stash config --dir ~/notes")
		os.Exit(1)
	}

	// Check if target directory exists
	if _, err := os.Stat(cfg.TargetDir); os.IsNotExist(err) {
		fmt.Printf("Error: Target directory does not exist: %s\n", cfg.TargetDir)
		fmt.Println("\nPlease create the directory or update the config:")
		fmt.Println("  stash config --dir <path>")
		os.Exit(1)
	}

	return nil
}

// GetTargetDir returns the target directory from config
func GetTargetDir() (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	return cfg.TargetDir, nil
}
