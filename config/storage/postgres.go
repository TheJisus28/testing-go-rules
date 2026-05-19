package storage

import (
	"context"
	"embed"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/donca/user-crud/config/generals/logger"
	"github.com/donca/user-crud/pkg/kit/enums"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type PostgresPool struct {
	Pool *pgxpool.Pool
}

func NewPostgresPool() (*PostgresPool, error) {
	url := os.Getenv(enums.DatabaseURL)
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("postgres pool: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping: %w", err)
	}
	p := &PostgresPool{Pool: pool}
	if err := p.runMigrations(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}
	logger.Get().Info().Msg("storage: postgres connected and migrated")
	return p, nil
}

func (p *PostgresPool) runMigrations(ctx context.Context) error {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	for _, name := range files {
		body, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if _, err := p.Pool.Exec(ctx, string(body)); err != nil {
			return fmt.Errorf("exec migration %s: %w", name, err)
		}
		logger.Get().Debug().Str("file", name).Msg("storage: migration applied")
	}
	return nil
}
