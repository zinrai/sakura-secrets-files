package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// Config represents the entire configuration file
type Config struct {
	Vault   VaultConfig    `yaml:"vault"`
	Secrets []SecretConfig `yaml:"secrets"`
}

// VaultConfig represents the vault configuration
type VaultConfig struct {
	ID   string `yaml:"id"`
	Zone string `yaml:"zone"`
}

// SecretConfig represents a single secret configuration
type SecretConfig struct {
	Name    string `yaml:"name"`
	Dest    string `yaml:"dest"`
	Version int    `yaml:"version,omitempty"`
}

// LoadConfig loads and parses the YAML configuration file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.Vault.Zone == "" {
		cfg.Vault.Zone = "is1a"
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Vault.ID == "" {
		return fmt.Errorf("vault.id is required")
	}

	if len(c.Secrets) == 0 {
		return fmt.Errorf("at least one secret must be defined")
	}

	for i, secret := range c.Secrets {
		if secret.Name == "" {
			return fmt.Errorf("secrets[%d].name is required", i)
		}
		if secret.Dest == "" {
			return fmt.Errorf("secrets[%d].dest is required", i)
		}
	}

	return nil
}
