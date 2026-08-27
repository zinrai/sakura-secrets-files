package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteSecretToFileWritesTheValueOwnerOnly(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "database", "db_password")

	if err := WriteSecretToFile(dest, "s3cret"); err != nil {
		t.Fatalf("WriteSecretToFile() error = %v", err)
	}

	content, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("failed to read written secret: %v", err)
	}
	if string(content) != "s3cret" {
		t.Errorf("content = %q, want %q", content, "s3cret")
	}

	info, err := os.Stat(dest)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("file permissions = %v, want %v", info.Mode().Perm(), os.FileMode(0600))
	}

	parent, err := os.Stat(filepath.Dir(dest))
	if err != nil {
		t.Fatal(err)
	}
	if parent.Mode().Perm() != 0700 {
		t.Errorf("created directory permissions = %v, want %v", parent.Mode().Perm(), os.FileMode(0700))
	}
}

func TestWriteSecretToFileNarrowsAWiderExistingFile(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "db_password")
	if err := os.WriteFile(dest, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := WriteSecretToFile(dest, "s3cret"); err != nil {
		t.Fatalf("WriteSecretToFile() error = %v", err)
	}

	info, err := os.Stat(dest)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("file permissions = %v, want %v", info.Mode().Perm(), os.FileMode(0600))
	}

	entries, err := os.ReadDir(filepath.Dir(dest))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("%d entries left in the destination directory, want 1", len(entries))
	}
}
