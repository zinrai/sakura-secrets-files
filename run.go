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

	cfg, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	// Up front, not inside the fetch loop: a manifest that would write outside
	// base-dir must fail before a single secret leaves the vault
	plan, err := planSecrets(absBase, cfg.Secrets)
	if err != nil {
		return err
	}

	client, err := newClientFromEnv(cfg.Vault.Zone)
	if err != nil {
		return err
	}

	return materialize(client, cfg.Vault.ID, plan)
}

func planSecrets(baseDir string, secrets []secretConfig) ([]plannedSecret, error) {
	plan := make([]plannedSecret, 0, len(secrets))
	planned := make(map[string]string, len(secrets))
	for _, secret := range secrets {
		dest, err := resolveDest(baseDir, secret.Dest)
		if err != nil {
			return nil, err
		}
		if other, ok := planned[dest]; ok {
			return nil, fmt.Errorf("secrets '%s' and '%s' both write to %s", other, secret.Name, dest)
		}
		planned[dest] = secret.Name
		plan = append(plan, plannedSecret{name: secret.Name, dest: dest})
	}

	return plan, nil
}

func materialize(client *sakuraClient, vaultID string, plan []plannedSecret) error {
	for _, secret := range plan {
		value, err := client.getSecret(vaultID, secret.name)
		if err != nil {
			return fmt.Errorf("failed to get secret '%s': %w", secret.name, err)
		}

		if err := writeSecretToFile(secret.dest, value); err != nil {
			return fmt.Errorf("failed to write secret '%s' to %s: %w", secret.name, secret.dest, err)
		}

		fmt.Fprintf(os.Stderr, "[write] %s -> %s\n", secret.name, secret.dest)
	}

	return nil
}
