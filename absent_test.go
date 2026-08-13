package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRunAbsent(t *testing.T) {
	dir := t.TempDir()

	existing := filepath.Join(dir, "out", "a.txt")
	missing := filepath.Join(dir, "out", "b.txt")
	if err := os.MkdirAll(filepath.Dir(existing), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(existing, []byte("secret"), 0600); err != nil {
		t.Fatal(err)
	}

	config := filepath.Join(dir, "secrets.yaml")
	manifest := fmt.Sprintf(`vault:
  id: "123456789012"
secrets:
  - name: a
    dest: %s
  - name: b
    dest: %s
`, existing, missing)
	if err := os.WriteFile(config, []byte(manifest), 0600); err != nil {
		t.Fatal(err)
	}

	if err := runAbsent([]string{"-config", config}); err != nil {
		t.Fatalf("runAbsent() error = %v", err)
	}
	if _, err := os.Stat(existing); !os.IsNotExist(err) {
		t.Errorf("existing dest was not removed")
	}

	// Running again with everything already absent must still succeed.
	if err := runAbsent([]string{"-config", config}); err != nil {
		t.Errorf("second runAbsent() error = %v", err)
	}
}

func TestRunAbsentDryRun(t *testing.T) {
	dir := t.TempDir()

	existing := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(existing, []byte("secret"), 0600); err != nil {
		t.Fatal(err)
	}

	config := filepath.Join(dir, "secrets.yaml")
	manifest := fmt.Sprintf(`vault:
  id: "123456789012"
secrets:
  - name: a
    dest: %s
`, existing)
	if err := os.WriteFile(config, []byte(manifest), 0600); err != nil {
		t.Fatal(err)
	}

	if err := runAbsent([]string{"-config", config, "-dry-run"}); err != nil {
		t.Fatalf("runAbsent() error = %v", err)
	}
	if _, err := os.Stat(existing); err != nil {
		t.Errorf("dry-run removed the file")
	}
}
