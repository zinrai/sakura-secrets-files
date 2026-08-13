package main

import (
	"flag"
	"fmt"
	"os"
)

// runPresent executes the present subcommand
func runPresent(args []string) error {
	fs := flag.NewFlagSet("present", flag.ExitOnError)
	configPath := fs.String("config", "", "Path to configuration file (required)")
	dryRun := fs.Bool("dry-run", false, "Report decisions without writing files")

	fs.Parse(args)

	if *configPath == "" {
		fs.Usage()
		return fmt.Errorf("-config is required")
	}

	cfg, err := LoadConfig(*configPath)
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
			return fmt.Errorf("failed to get secret '%s': %w", secret.Name, err)
		}

		status, err := fileStatus(secret.Dest, value)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", secret.Dest, err)
		}

		if status != "unchanged" && !*dryRun {
			if err := WriteSecretToFile(secret.Dest, value); err != nil {
				return fmt.Errorf("failed to write secret '%s' to %s: %w", secret.Name, secret.Dest, err)
			}
		}

		fmt.Fprintf(os.Stderr, "[%s] %s -> %s\n", status, secret.Name, secret.Dest)
	}

	if *dryRun {
		fmt.Fprintln(os.Stderr, "dry-run: no files written")
	}

	return nil
}

// fileStatus reports what writing value to dest would do.
func fileStatus(dest, value string) (string, error) {
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
