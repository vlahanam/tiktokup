package database

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const createMigrationsTableSQL = `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version TEXT PRIMARY KEY,
	applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`

// ListMigrationFiles returns sorted *.sql migration filenames from fsys.
func ListMigrationFiles(fsys fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("read migrations directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)
	return files, nil
}

// Migrate applies all pending SQL migrations transactionally and idempotently.
func Migrate(ctx context.Context, pool *pgxpool.Pool, fsys fs.FS, logger *slog.Logger) (int, error) {
	if _, err := pool.Exec(ctx, createMigrationsTableSQL); err != nil {
		return 0, fmt.Errorf("ensure schema_migrations table: %w", err)
	}

	files, err := ListMigrationFiles(fsys)
	if err != nil {
		return 0, err
	}

	appliedCount := 0
	for _, filename := range files {
		var exists bool
		err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", filename).Scan(&exists)
		if err != nil {
			return appliedCount, fmt.Errorf("check migration %s: %w", filename, err)
		}
		if exists {
			continue
		}

		sqlBytes, err := fs.ReadFile(fsys, filename)
		if err != nil {
			return appliedCount, fmt.Errorf("read migration %s: %w", filename, err)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return appliedCount, fmt.Errorf("begin migration tx %s: %w", filename, err)
		}

		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
			_ = tx.Rollback(ctx)
			return appliedCount, fmt.Errorf("execute migration %s: %w", filename, err)
		}

		if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", filename); err != nil {
			_ = tx.Rollback(ctx)
			return appliedCount, fmt.Errorf("record migration %s: %w", filename, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return appliedCount, fmt.Errorf("commit migration %s: %w", filename, err)
		}

		appliedCount++
		logger.Info("postgres_migration_applied", slog.String("version", filename))
	}

	return appliedCount, nil
}
