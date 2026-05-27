package persistence

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

func RunMigrations(ctx context.Context, db *PostgresDB, migrationsDir string) error {
	_, err := db.Pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		filename TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`)
	if err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	for _, f := range upFiles {
		var count int
		if err := db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE filename = $1", f).Scan(&count); err != nil {
			return fmt.Errorf("check migration %s: %w", f, err)
		}
		if count > 0 {
			log.Printf("migration already applied, skipping: %s", f)
			continue
		}

		path := filepath.Join(migrationsDir, f)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f, err)
		}

		sql := string(content)
		if _, err := db.Pool.Exec(ctx, sql); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "42P07" {
				log.Printf("migration already applied (detected), recording: %s", f)
			} else {
				return fmt.Errorf("execute migration %s: %w", f, err)
			}
		} else {
			log.Printf("migration applied: %s", f)
		}

		if _, err := db.Pool.Exec(ctx, "INSERT INTO schema_migrations (filename) VALUES ($1)", f); err != nil {
			return fmt.Errorf("record migration %s: %w", f, err)
		}
	}

	return nil
}
