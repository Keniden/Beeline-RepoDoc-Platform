package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/beeline/repodoc/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg configs.StorageConfig) (*Store, error) {
	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	return &Store{Pool: pool}, nil
}

func (s *Store) Close() {
	if s.Pool != nil {
		s.Pool.Close()
	}
}

func (s *Store) SaveRepo(ctx context.Context, repo string) error {
	_, err := s.Pool.Exec(ctx, "INSERT INTO repos (id, name, url, status, created_at, updated_at) VALUES (gen_random_uuid(), $1, $2, 'pending', NOW(), NOW())", repo, repo)
	return err
}

func (s *Store) RunMigrations(ctx context.Context) error {
	const (
		migrationVersion = "001_init"
		advisoryLockKey  = int64(83942011)
	)

	migrationPath := os.Getenv("REPDOC_MIGRATION_FILE")
	if migrationPath == "" {
		if _, err := os.Stat("/migrations/001_init.sql"); err == nil {
			migrationPath = "/migrations/001_init.sql"
		} else {
			migrationPath = "migrations/001_init.sql"
		}
	}

	content, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("read migration %s: %w", migrationPath, err)
	}

	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", advisoryLockKey); err != nil {
		return fmt.Errorf("advisory lock: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	var alreadyApplied bool
	if err := tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)", migrationVersion).Scan(&alreadyApplied); err != nil {
		return fmt.Errorf("check schema_migrations: %w", err)
	}
	if alreadyApplied {
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit tx: %w", err)
		}
		return nil
	}

	// If schema already exists (for example, created manually), mark migration as applied.
	var reposRegclass *string
	if err := tx.QueryRow(ctx, "SELECT to_regclass('public.repos')").Scan(&reposRegclass); err != nil {
		return fmt.Errorf("check existing schema: %w", err)
	}
	if reposRegclass == nil {
		if _, err := tx.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("run migration SQL: %w", err)
		}
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO schema_migrations(version)
		VALUES($1)
		ON CONFLICT (version) DO NOTHING
	`, migrationVersion); err != nil {
		return fmt.Errorf("record migration: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
