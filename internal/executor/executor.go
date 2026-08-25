package executor

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type Executor struct {
	db *sql.DB
}

func New(cfg Config) (*Executor, error) {
	db, err := sql.Open("pgx", cfg.DSN)

	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	return &Executor{db: db}, nil
}

func (e *Executor) Ping(ctx context.Context) error {
	return e.db.PingContext(ctx)
}

func (e *Executor) Close() error {
	return e.db.Close()
}
