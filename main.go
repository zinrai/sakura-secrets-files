package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		configPath string
		zone       string
		dryRun     bool
	)

	flag.StringVar(&configPath, "config", "", "Path to configuration file (required)")
	flag.StringVar(&zone, "zone", "is1a", "Sakura Cloud zone (e.g., is1a, tk1a)")
	flag.BoolVar(&dryRun, "dry-run", false, "Show what would be done without actually doing it")
	flag.Parse()

	if configPath == "" {
		fmt.Fprintln(os.Stderr, "Error: -config flag is required")
		flag.Usage()
		os.Exit(1)
	}

	if err := run(configPath, zone, dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(configPath, zone string, dryRun bool) error {
	// Load configuration
	fmt.Printf("Loading configuration from %s...\n", configPath)

	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	fmt.Printf("Configuration loaded: %d secret(s) to pull\n", len(cfg.Secrets))

	// Create API client
	client, err := NewClientFromEnv(zone)
	if err != nil {
		return err
	}

	fmt.Printf("Using zone: %s\n", client.Zone)
	fmt.Printf("Using vault: %s\n", cfg.Vault.ID)

	// Fetch and write each secret
	for i, secret := range cfg.Secrets {
		fmt.Printf("[%d/%d] Pulling secret '%s'...\n", i+1, len(cfg.Secrets), secret.Name)

		value, err := client.GetSecret(cfg.Vault.ID, secret.Name, secret.Version)
		if err != nil {
			return fmt.Errorf("failed to pull secret '%s': %w", secret.Name, err)
		}

		if dryRun {
			fmt.Printf("[OK] [DRY-RUN] Would write: %s -> %s\n", secret.Name, secret.Dest)
			continue
		}

		if err := WriteSecretToFile(secret.Dest, value); err != nil {
			return fmt.Errorf("failed to write secret '%s' to %s: %w", secret.Name, secret.Dest, err)
		}

		fmt.Printf("[OK] %s -> %s\n", secret.Name, secret.Dest)
	}

	if dryRun {
		fmt.Println("\nDry-run completed successfully (no files were written)")
	} else {
		fmt.Printf("\nSuccessfully pulled and wrote %d secret(s)\n", len(cfg.Secrets))
	}

	return nil
}
