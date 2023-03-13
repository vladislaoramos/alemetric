package repo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/pkg/postgres"
)

type PostgresRepo struct {
	*postgres.DB
}

func NewPostgresRepo(pg *postgres.DB) (*PostgresRepo, error) {
	return &PostgresRepo{pg}, nil
}

func (r *PostgresRepo) GetMetricsNames(ctx context.Context) []string {
	res := make([]string, 0)
	pgxscan.Select(ctx, r.Pool, &res, "select name from metrics;")

	return res
}

func (r *PostgresRepo) GetMetrics(ctx context.Context, name string) (entity.Metrics, error) {
	q, args, err := r.Builder.
		Select(
			"name",
			"mtype",
			"delta",
			"value",
			"hash").
		From("metrics").
		Where(sq.Eq{"name": name}).
		ToSql()
	if err != nil {
		return entity.Metrics{}, fmt.Errorf("builder error getting metrics from db: %w", err)
	}

	dst := make([]entity.Metrics, 0)
	if err = pgxscan.Select(ctx, r.Pool, &dst, q, args...); err != nil {
		return entity.Metrics{}, fmt.Errorf("error selecting metrics from db: %w", err)
	}

	if len(dst) == 0 {
		return entity.Metrics{}, ErrNotFound
	}

	return dst[0], nil
}

func (r *PostgresRepo) StoreMetrics(ctx context.Context, metrics entity.Metrics) error {
	updateQuery, updateArgs, err := r.Builder.
		Update("metrics").
		Set("delta", metrics.Delta).
		Set("value", metrics.Value).
		Set("hash", metrics.Hash).
		Where(sq.Eq{"name": metrics.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("builder error storing metrics: %w", err)
	}

	insertQuery, insertArgs, err := r.Builder.
		Insert("metrics").
		Columns(
			"name",
			"mtype",
			"delta",
			"value",
			"hash").
		Values(
			metrics.ID,
			metrics.MType,
			metrics.Delta,
			metrics.Value,
			metrics.Hash).
		ToSql()
	if err != nil {
		return fmt.Errorf("error inserting metrics into db: %w", err)
	}

	_, err = r.Pool.Exec(ctx, insertQuery, insertArgs...)
	if err != nil {
		_, err = r.Pool.Exec(ctx, updateQuery, updateArgs...)
		if err != nil {
			return fmt.Errorf("error executing insert query: %w", err)
		}
	}

	return nil
}

func (r *PostgresRepo) StoreAll() error {
	return nil
}

func (r *PostgresRepo) Upload(context.Context) error {
	return nil
}

func (r *PostgresRepo) Ping(ctx context.Context) error {
	return r.Pool.Ping(ctx)
}
