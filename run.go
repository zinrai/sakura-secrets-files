package main

import (
	"fmt"
	"os"
)

type plannedSecret struct {
	name string
	dest string
}

func run(configPath, baseDir string) error {
	absBase, err := resolveBaseDir(baseDir)
	if err != nil {
		return err
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	// Up front, not inside the fetch loop: a manifest that would write outside
	// base-dir must fail before a single secret leaves the vault
	plan := make([]plannedSecret, 0, len(cfg.Secrets))
	for _, secret := range cfg.Secrets {
		dest, err := resolveDest(absBase, secret.Dest)
		if err != nil {
			return err
		}
		plan = append(plan, plannedSecret{name: secret.Name, dest: dest})
	}

	client, err := NewClientFromEnv(cfg.Vault.Zone)
	if err != nil {
		return err
	}

	return materialize(client, cfg.Vault.ID, plan)
}

func materialize(client *SakuraClient, vaultID string, plan []plannedSecret) error {
	for _, secret := range plan {
		value, err := client.GetSecret(vaultID, secret.name)
		if err != nil {
			return fmt.Errorf("failed to get secret '%s': %w", secret.name, err)
		}

		if err := WriteSecretToFile(secret.dest, value); err != nil {
			return fmt.Errorf("failed to write secret '%s' to %s: %w", secret.name, secret.dest, err)
		}

		fmt.Fprintf(os.Stderr, "[write] %s -> %s\n", secret.name, secret.dest)
	}

	return nil
}
