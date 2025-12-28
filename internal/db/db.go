package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/backoff"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PGConnect struct {
	db      *sql.DB
	backoff *backoff.LinearBackoff
}

func New(cfg Config, backoff *backoff.LinearBackoff) (conn *PGConnect, err error) {
	dsn, err := cfg.ToDSN()

	if err != nil {
		return nil, err
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("db not ok, %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("migration driver not ok, %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("migration not ok, %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to migrate, %w", err)
	}

	return &PGConnect{
		db:      db,
		backoff: backoff,
	}, nil

}

func (pg *PGConnect) Ping(ctx context.Context) error {
	fn := func(ctx context.Context) error {
		return pg.db.PingContext(ctx)
	}
	backoff := pg.backoff.WithRetry()

	return backoff(fn)(ctx)
}

func (pg *PGConnect) QueryWithOneResult(
	ctx context.Context, dst any, query string, args ...any,
) error {
	fn := func(ctx context.Context) error {
		row := pg.db.QueryRowContext(ctx, query, args...)

		if err := row.Scan(dst); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return errors.Join(err, errScanRow)
		}
		return nil
	}
	backoff := pg.backoff.WithRetry()

	return backoff(fn)(ctx)
}

func (pg *PGConnect) QueryWithOneResultJSON(
	ctx context.Context, dst any, query string, args ...any,
) error {

	var res []byte

	if err := pg.QueryWithOneResult(ctx, &res, query, args...); err != nil {
		return err
	}

	if len(res) == 0 {
		return errNoData
	}

	if err := json.Unmarshal(res, dst); err != nil {
		return errors.Join(err, errUnmarshal)
	}
	return nil
}

func (pg *PGConnect) Close() error {
	return pg.db.Close()
}
