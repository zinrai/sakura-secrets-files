package main

import "testing"

func TestResolveDest(t *testing.T) {
	base := "/tmp/tmp.abc123"

	tests := []struct {
		name    string
		dest    string
		want    string
		wantErr bool
	}{
		{"a nested dest lands below base-dir", "database/files/db_password", "/tmp/tmp.abc123/database/files/db_password", false},
		{"dot dot that stays below base-dir is accepted", "database/../db_password", "/tmp/tmp.abc123/db_password", false},
		{"an absolute dest is rejected", "/etc/passwd", "", true},
		{"a dest that climbs out of base-dir is rejected", "../db_password", "", true},
		{"a dest that climbs out through a subdirectory is rejected", "database/../../db_password", "", true},
		{"a dest naming base-dir itself is rejected", ".", "", true},
	}

	for _, tt := range tests {
		got, err := resolveDest(base, tt.dest)
		if tt.wantErr {
			if err == nil {
				t.Errorf("%s: resolveDest() = %q, want error", tt.name, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: resolveDest() error = %v", tt.name, err)
			continue
		}
		if got != tt.want {
			t.Errorf("%s: resolveDest() = %q, want %q", tt.name, got, tt.want)
		}
	}
}
