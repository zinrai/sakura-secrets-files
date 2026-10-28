package main

import "testing"

func TestPlanSecretsRejectsTwoSecretsWritingToOneDest(t *testing.T) {
	secrets := []secretConfig{
		{Name: "first", Dest: "database/db_password"},
		{Name: "second", Dest: "database/./db_password"},
	}

	if plan, err := planSecrets("/tmp/tmp.abc123", secrets); err == nil {
		t.Errorf("planSecrets() = %v, want error", plan)
	}
}
