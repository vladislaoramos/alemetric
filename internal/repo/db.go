package repo

import (
	"context"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/pkg/postgres"
)

type PostgresRepo struct {
	*postgres.DB
}

func NewPostgresRepo(pg *postgres.DB) *PostgresRepo {
	return &PostgresRepo{pg}
}

func (r *PostgresRepo) GetMetricsNames() []string {
	return nil
}

func (r *PostgresRepo) GetMetrics(_ string) (entity.Metrics, error) {
	return entity.Metrics{}, nil
}

func (r *PostgresRepo) StoreMetrics(_ entity.Metrics) error {
	return nil
}

func (r *PostgresRepo) StoreAll() error {
	return nil
}

func (r *PostgresRepo) Upload() error {
	return nil
}

func (r *PostgresRepo) Ping(ctx context.Context) error {
	return r.Pool.Ping(ctx)
}
