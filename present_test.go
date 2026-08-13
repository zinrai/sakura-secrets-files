package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileStatus(t *testing.T) {
	dir := t.TempDir()

	same := filepath.Join(dir, "same")
	if err := os.WriteFile(same, []byte("value"), 0600); err != nil {
		t.Fatal(err)
	}
	differs := filepath.Join(dir, "differs")
	if err := os.WriteFile(differs, []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		dest  string
		value string
		want  string
	}{
		{"missing file", filepath.Join(dir, "missing"), "value", "create"},
		{"same content", same, "value", "unchanged"},
		{"different content", differs, "value", "update"},
	}

	for _, tt := range tests {
		got, err := fileStatus(tt.dest, tt.value)
		if err != nil {
			t.Errorf("%s: fileStatus() error = %v", tt.name, err)
			continue
		}
		if got != tt.want {
			t.Errorf("%s: fileStatus() = %q, want %q", tt.name, got, tt.want)
		}
	}
}
