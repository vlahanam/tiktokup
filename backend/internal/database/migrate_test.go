package database

import (
	"io/fs"
	"strings"
	"testing"

	"tiktokup/backend/migrations"
)

func TestEmbeddedMigrationsPresent(t *testing.T) {
	files, err := ListMigrationFiles(migrations.Files)
	if err != nil {
		t.Fatalf("failed to list embedded migrations: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("expected at least one embedded SQL migration file")
	}
	if files[0] != "000001_init.sql" {
		t.Fatalf("expected first migration to be 000001_init.sql, got %q", files[0])
	}

	for _, name := range files {
		raw, err := fs.ReadFile(migrations.Files, name)
		if err != nil {
			t.Fatalf("failed reading migration %s: %v", name, err)
		}
		if strings.TrimSpace(string(raw)) == "" {
			t.Fatalf("migration %s is empty", name)
		}
	}
}
