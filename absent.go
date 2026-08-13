package main

import (
	"flag"
	"fmt"
	"os"
)

// runAbsent executes the absent subcommand. It needs no credentials and no
// network, because removing the declared paths is a purely local operation.
func runAbsent(args []string) error {
	fs := flag.NewFlagSet("absent", flag.ExitOnError)
	configPath := fs.String("config", "", "Path to configuration file (required)")
	dryRun := fs.Bool("dry-run", false, "Report what would be removed without removing")

	fs.Parse(args)

	if *configPath == "" {
		fs.Usage()
		return fmt.Errorf("-config is required")
	}

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		return err
	}

	for _, secret := range cfg.Secrets {
		_, err := os.Stat(secret.Dest)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "[absent] %s\n", secret.Dest)
			continue
		}
		if err != nil {
			return err
		}

		if *dryRun {
			fmt.Fprintf(os.Stderr, "[DRY-RUN] %s (remove)\n", secret.Dest)
			continue
		}

		if err := os.Remove(secret.Dest); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "[removed] %s\n", secret.Dest)
	}

	if *dryRun {
		fmt.Fprintln(os.Stderr, "dry-run: no files removed")
	}

	return nil
}
