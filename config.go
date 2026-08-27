package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Vault   VaultConfig    `yaml:"vault"`
	Secrets []SecretConfig `yaml:"secrets"`
}

type VaultConfig struct {
	ID   string `yaml:"id"`
	Zone string `yaml:"zone"`
}

type SecretConfig struct {
	Name string `yaml:"name"`
	Dest string `yaml:"dest"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Strict, so that a key this tool has stopped honoring fails loudly rather
	// than leaving the caller with a setting that no longer does anything
	var cfg Config
	if err := yaml.UnmarshalWithOptions(data, &cfg, yaml.Strict()); err != nil {
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

// Not MkdirAll: the caller is expected to have prepared this directory, and
// creating a missing one would put secrets somewhere nobody chose
func resolveBaseDir(baseDir string) (string, error) {
	abs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve base-dir: %w", err)
	}

	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("base-dir is not usable: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("base-dir is not a directory: %s", abs)
	}

	return abs, nil
}

// Lexical, not filepath.EvalSymlinks: base-dir belongs to the caller, so a
// symlink planted inside it is not a threat this can meaningfully defend against
func resolveDest(baseDir, dest string) (string, error) {
	if filepath.IsAbs(dest) {
		return "", fmt.Errorf("dest must be relative to base-dir: %s", dest)
	}

	resolved := filepath.Join(baseDir, dest)

	rel, err := filepath.Rel(baseDir, resolved)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("dest escapes base-dir: %s", dest)
	}
	if rel == "." {
		return "", fmt.Errorf("dest must name a file, not base-dir itself: %s", dest)
	}

	return resolved, nil
}
