-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TYPE metric_types AS ENUM (
  'counter',
  'gauge'
);

CREATE TABLE IF NOT EXISTS metrics(
    id serial PRIMARY KEY,
    name VARCHAR(255),
    mtype metric_types,
    delta BIGINT,
    value DOUBLE PRECISION,
    hash VARCHAR(255),
    CONSTRAINT unique_name_idx UNIQUE (name)
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE metrics;
DROP TYPE metric_types;
-- +goose StatementEnd