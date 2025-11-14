package postgres

import (
    "context"
    "fmt"

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
