package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	connAttempts = 10
	connTimeout  = time.Second
	maxPoolSize  = 1
)

type DB struct {
	connAttempts int
	connTimeout  time.Duration
	maxPoolSize  int32

	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

func New(url string) (*DB, error) {
	db := &DB{
		connAttempts: connAttempts,
		connTimeout:  connTimeout,
		maxPoolSize:  maxPoolSize,
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	db.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig.MaxConns = db.maxPoolSize

	for db.connAttempts > 0 {
		db.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		} else {
			log.Printf("error db connecting: %s", err)
		}
		log.Printf("connection attempts of PostgreSQL: %d", db.connAttempts)
		time.Sleep(db.connTimeout)
		db.connAttempts--
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (p *DB) Close() {
	if p.Pool == nil {
		p.Pool.Close()
	}
}
