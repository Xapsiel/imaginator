package repository

import (
	"context"
	"errors"
	"fmt"

	"imageBot/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(config config.DatabaseConfig) (*pgxpool.Pool, error) {
	PgUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", config.User, config.Password, config.Host, config.Port, config.Name, config.Sslmode)
	err := Migrate(PgUrl)
	if err != nil {
		return nil, fmt.Errorf("migration error: %s", err.Error())
	}
	poolConfig, err := pgxpool.ParseConfig(PgUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}
	poolConfig.MaxConns = config.MaxConnections
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create pool to database: %s", err.Error())
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %s", err.Error())
	}
	return pool, nil
}
func Migrate(pgURL string) error {
	m, err := migrate.New("file://migrations", pgURL)
	if err != nil {
		return err
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}
