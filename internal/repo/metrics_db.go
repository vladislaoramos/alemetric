package repo

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
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
	res := make([]string, 0)
	q, args, _ := r.Builder.
		Select("name").
		From("metrics").
		ToSql()

	_ = pgxscan.Select(context.TODO(), r.Pool, &res, q, args)

	return res
}

func (r *PostgresRepo) GetMetrics(name string) (entity.Metrics, error) {
	q, args, err := r.Builder.
		Select("name", "mtype", "delta", "value", "hash").
		From("metrics").
		Where(squirrel.Eq{"name": name}).
		ToSql()

	if err != nil {
		return entity.Metrics{}, err
	}

	dst := make([]entity.Metrics, 0)
	if err = pgxscan.Select(context.TODO(), r.Pool, &dst, q, args...); err != nil {
		return entity.Metrics{}, err
	}

	if len(dst) == 0 {
		return entity.Metrics{}, ErrNotFound
	}

	return dst[0], nil
}

func (r *PostgresRepo) StoreMetrics(metrics entity.Metrics) error {
	insertQuery, insertArgs, err := r.Builder.
		Insert("metrics").
		Columns("name", "mtype", "delta", "value", "hash").
		Values(metrics.ID, metrics.MType, metrics.Delta, metrics.Value, metrics.Hash).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.Pool.Exec(context.TODO(), insertQuery, insertArgs...)
	if err == nil {
		return nil
	}

	updateQuery, updateArgs, err := r.Builder.
		Update("metrics").
		Set("delta", metrics.Delta).
		Set("value", metrics.Value).
		Where(squirrel.Eq{"name": metrics.ID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.Pool.Exec(context.TODO(), updateQuery, updateArgs...)
	if err != nil {
		return err
	}

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
