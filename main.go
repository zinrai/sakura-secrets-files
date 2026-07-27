package main

import (
	"flag"
	"fmt"
	"os"
)

// Injected at build time by goreleaser via -ldflags -X
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var (
		configPath  string
		dryRun      bool
		showVersion bool
	)

	flag.StringVar(&configPath, "config", "", "Path to configuration file (required)")
	flag.BoolVar(&dryRun, "dry-run", false, "Show what would be done without actually doing it")
	flag.BoolVar(&showVersion, "version", false, "Print version")
	flag.Parse()

	if showVersion {
		fmt.Printf("sakura-secrets-pull version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built: %s\n", date)
		return
	}

	if configPath == "" {
		fmt.Fprintln(os.Stderr, "Error: -config flag is required")
		flag.Usage()
		os.Exit(1)
	}

	if err := run(configPath, dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(configPath string, dryRun bool) error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	client, err := NewClientFromEnv(cfg.Vault.Zone)
	if err != nil {
		return err
	}

	for _, secret := range cfg.Secrets {
		value, err := client.GetSecret(cfg.Vault.ID, secret.Name, secret.Version)
		if err != nil {
			return fmt.Errorf("failed to pull secret '%s': %w", secret.Name, err)
		}

		if dryRun {
			status, err := dryRunStatus(secret.Dest, value)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", secret.Dest, err)
			}
			fmt.Fprintf(os.Stderr, "[DRY-RUN] %s -> %s (%s)\n", secret.Name, secret.Dest, status)
			continue
		}

		if err := WriteSecretToFile(secret.Dest, value); err != nil {
			return fmt.Errorf("failed to write secret '%s' to %s: %w", secret.Name, secret.Dest, err)
		}

		fmt.Fprintf(os.Stderr, "[OK] %s -> %s\n", secret.Name, secret.Dest)
	}

	if dryRun {
		fmt.Fprintf(os.Stderr, "dry-run: %d secret(s), no files written\n", len(cfg.Secrets))
	} else {
		fmt.Fprintf(os.Stderr, "pulled %d secret(s)\n", len(cfg.Secrets))
	}

	return nil
}

// dryRunStatus reports what writing value to dest would do.
func dryRunStatus(dest, value string) (string, error) {
	current, err := os.ReadFile(dest)
	if err != nil {
		if os.IsNotExist(err) {
			return "create", nil
		}
		return "", err
	}
	if string(current) == value {
		return "unchanged", nil
	}
	return "update", nil
}
